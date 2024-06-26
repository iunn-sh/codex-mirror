package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	// TODO: error handling

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Default level: info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	fileUrl := "https://law.moj.gov.tw/api/data/chlaw.json.zip"
	rawdir := filepath.Join(".", "raw")
	err := os.MkdirAll(rawdir, os.ModePerm)
	// rawdir, err := filepath.Abs("./raw")
	if err != nil {
		log.Error().Err(err).Send()
	}
	zippath := filepath.Join(rawdir, "ChLaw.json.zip")
	depotdir := filepath.Join(".", "depot")
	err = os.MkdirAll(depotdir, os.ModePerm)
	// depotdir, err := filepath.Abs("./depot")
	if err != nil {
		log.Error().Err(err).Send()
	}

	err = Cleanup(rawdir)
	if err != nil {
		log.Error().Err(err).Send()
	}
	err = Cleanup(depotdir)
	if err != nil {
		log.Error().Err(err).Send()
	}

	err = Download(zippath, fileUrl)
	if err != nil {
		log.Error().Err(err).Send()
	}

	err = Unzip(zippath, rawdir)
	if err != nil {
		log.Error().Err(err).Send()
	}

	// TODO: deal with newline within 'ArticleContent'
	// TODO: deal with '（刪除）' of ArticleContent (or not)
	err = ParseAndSplit(filepath.Join(rawdir, "ChLaw.json"), depotdir)
	if err != nil {
		log.Error().Err(err).Send()
	}

	tmplfile := "law.tmpl"
	mddir, err := filepath.Abs("./docs/")
	if err != nil {
		log.Error().Err(err).Send()
	}
	err = Cleanup(mddir)
	if err != nil {
		log.Error().Err(err).Send()
	}

	readmefile, err := filepath.Abs("README.md")
	if err != nil {
		log.Error().Err(err).Send()
	}
	err = CopyFile(readmefile, mddir+"/index.md")
	if err != nil {
		log.Error().Err(err).Send()
	}

	// TODO: 中華民國刑法 includes 編/章 -> might need extra template to reflect that
	// TODO: read list from config file -> share list with mkdocs is even better
	counter := 0
	for _, p := range GetFileList(depotdir, ".json") {
		err = JsonToMarkdown(p, tmplfile, mddir)
		if err != nil {
			log.Error().Err(err).Send()
		}
		counter++
	}
	log.Info().Int("Processed .md count", counter).Send()
}

func Cleanup(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	log.Info().Str("Remove files in", dir).Send()

	return nil
}

func Download(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	log.Info().Str("Downloaded from", url).Str("to", filepath).Send()

	return err
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}
	log.Info().Str("Unzipped from", src).Str("to", dest).Send()

	return nil
}

func ParseAndSplit(srcfile string, destdir string) error {
	fileContent, err := os.Open(srcfile)
	if err != nil {
		return err
	}
	defer fileContent.Close()

	byteResult, _ := io.ReadAll(fileContent)
	// Cleanup data https://stackoverflow.com/questions/31398044/got-error-invalid-character-%C3%AF-looking-for-beginning-of-value-from-json-unmar
	byteResult = bytes.TrimPrefix(byteResult, []byte("\xef\xbb\xbf"))
	var codex Codex
	if err := json.Unmarshal(byteResult, &codex); err != nil {
		return err
	}
	log.Info().Str("Codex UpdateDate", codex.UpdateDate).Send()

	counterEnacted := 0
	counterRepealed := 0
	for _, p := range codex.Laws {
		fo, _ := json.MarshalIndent(p, "", " ")
		if "廢" != p.LawAbandonNote {
			shortLawName := TrimLawName(p.LawName)
			_ = os.WriteFile(filepath.Join(destdir, shortLawName+".json"), fo, 0644)
			counterEnacted++
			log.Debug().Str("Enacted law", p.LawName).Send()
		} else {
			counterRepealed++
			log.Debug().Str("Repealed law (skip)", p.LawName).Send()
		}
	}
	log.Info().Int("Enacted law count", counterEnacted).Int("Repealed law count", counterRepealed).Send()
	log.Info().Str("Parsed and splitted enacted laws from", srcfile).Str("to", destdir).Send()

	return nil
}

func CopyFile(src, dst string) error {
	buf := make([]byte, 1024)

	fin, err := os.Open(src)
	if err != nil {
		return err
	}

	defer fin.Close()

	fout, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer fout.Close()

	for {
		n, err := fin.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := fout.Write(buf[:n]); err != nil {
			return err
		}
	}
	log.Info().Str("Copied file from", src).Str("to", dst).Send()

	return nil
}

func GetFileList(dir, ext string) []string {
	var a []string
	filepath.WalkDir(dir, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	log.Info().Int("File list length", len(a)).Send()

	return a
}

func JsonToMarkdown(jsonfile string, tmplfile string, destdir string) error {
	lawFile, err := os.Open(jsonfile)
	if err != nil {
		return err
	}
	defer lawFile.Close()

	byteResult, _ := io.ReadAll(lawFile)
	var law Law
	if err := json.Unmarshal(byteResult, &law); err != nil {
		return err
	}
	log.Debug().Str("Processed from .json", law.LawName).Send()
	shortLawName := TrimLawName(law.LawName)
	f, err := os.Create(filepath.Join(destdir, shortLawName+".md"))
	if err != nil {
		return err
	}
	defer f.Close()

	// Execute the template to the file
	tmpl, err := template.ParseFiles(tmplfile)
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, law)
	if err != nil {
		return err
	}
	log.Debug().Str("Processed to .md", law.LawName).Send()

	return nil
}

func TrimLawName(lawName string) string {
	before, _, found := strings.Cut(lawName, "（")

	shortname := lawName
	if found {
		log.Debug().Str("Original LawName", lawName).Str("Trimmed", before).Send()
		shortname = before
	}

	return shortname
}

package main

import (
	"os"
	"io"
	"net/http"
	"encoding/json"
	"fmt"
    "log"
    "archive/zip"
    "path/filepath"
    "strings"
    "bytes"
)

type LawArticle struct {
	ArticleType    string `json:"ArticleType"`
	ArticleNo      string `json:"ArticleNo"`
	ArticleContent string `json:"ArticleContent"`
}

type LawAttachement struct {
	FileName string `json:"FileName"`
	FileUrl  string `json:"FileURL"`
}

type Law struct {
	LawLevel         string           `json:"LawLevel"`
	LawName          string           `json:"LawName"`
	LawUrl           string           `json:"LawURL"`
	LawCategory      string           `json:"LawCategory"`
	LawModifiedDate  string           `json:"LawModifiedDate"`
	LawEffectiveDate string           `json:"LawEffectiveDate"`
	LawEffectiveNote string           `json:"LawEffectiveNote"`
	LawAbandonNote   string           `json:"LawAbandonNote"`
	LawHasEngVersion string           `json:"LawHasEngVersion"`
	EngLawName       string           `json:"EngLawName"`
	LawAttachements  []LawAttachement `json:"LawAttachements"`
	LawHistories     string           `json:"LawHistories"`
	LawForeword      string           `json:"LawForeword"`
	LawArticles      []LawArticle     `json:"LawArticles"`
}

type Codex struct {
	UpdateDate string `json:"UpdateDate"`
	Laws       []Law  `json:"Laws"`
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
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

	return nil
}

func ParseAndSplit(srcfile string, destdir string) {
	fileContent, err := os.Open(srcfile)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(srcfile + " is opened successfully...")
	defer fileContent.Close()

	byteResult, _ := io.ReadAll(fileContent)
	// Cleanup data https://stackoverflow.com/questions/31398044/got-error-invalid-character-%C3%AF-looking-for-beginning-of-value-from-json-unmar
	byteResult = bytes.TrimPrefix(byteResult, []byte("\xef\xbb\xbf"))
	var codex Codex
	if err := json.Unmarshal(byteResult, &codex); err != nil {
		panic(err)
	}
	fmt.Println(codex.UpdateDate)
	for _, p := range codex.Laws {
		fo, _ := json.MarshalIndent(p, "", " ")
		_ = os.WriteFile(filepath.Join(destdir ,p.LawName+".json"), fo, 0644)
		fmt.Println(p.LawName + " is extracted...")
	}
}

func main() {
	// TODO: error handling

	fileUrl := "https://law.moj.gov.tw/api/Ch/Law/JSON"
	err := Download("./depot/ChLaw.json.zip", fileUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Downloaded: " + fileUrl)

    err = Unzip("./depot/ChLaw.json.zip", "./depot")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Unzipped: " + "ChLaw.json.zip")

    files, err := os.ReadDir("./depot/")
	if err != nil {
		log.Fatal(err)
	}
	dirpath, err := filepath.Abs("./depot/")
	for _, file := range files {
        fmt.Println(filepath.Join(dirpath ,file.Name()))
	}

	ParseAndSplit(filepath.Join(dirpath, "ChLaw.json"), dirpath)
}

package main

type Codex struct {
	UpdateDate string `json:"UpdateDate"`
	Laws       []Law  `json:"Laws"`
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

type LawAttachement struct {
	FileName string `json:"FileName"`
	FileUrl  string `json:"FileURL"`
}

type LawArticle struct {
	ArticleType    string `json:"ArticleType"`
	ArticleNo      string `json:"ArticleNo"`
	ArticleContent string `json:"ArticleContent"`
}

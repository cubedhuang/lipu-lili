package models

type WordResources struct {
	SonaPona          *string `json:"sona_pona,omitempty"`
	LipamankaSemantic *string `json:"lipamanka_semantic,omitempty"`
}

type WordRepresentations struct {
	SitelenEmosi   *string  `json:"sitelen_emosi,omitempty"`
	SitelenJelo    []string `json:"sitelen_jelo,omitempty"`
	Ligatures      []string `json:"ligatures,omitempty"`
	SitelenSitelen *string  `json:"sitelen_sitelen,omitempty"`
	Ucsur          *string  `json:"ucsur,omitempty"`
}

type WordEtymology struct {
	Word string `json:"word"`
	Alt  string `json:"alt"`
}

type WordAudio struct {
	Author string `json:"author"`
	Link   string `json:"link"`
}

type WordPuVerbatim struct {
	En string `json:"en"`
	Fr string `json:"fr"`
	De string `json:"de"`
	Eo string `json:"eo"`
}

type WordEtymologyTranslation struct {
	Definition string `json:"definition"`
	Language   string `json:"language"`
}

type WordTranslation struct {
	Commentary  string                     `json:"commentary"`
	Definition  string                     `json:"definition"`
	Etymology   []WordEtymologyTranslation `json:"etymology"`
	SpEtymology string                     `json:"sp_etymology"`
}

type WordData struct {
	Id                   string                     `json:"id"`
	AuthorVerbatim       string                     `json:"author_verbatim"`
	AuthorVerbatimSource string                     `json:"author_verbatim_source"`
	Book                 string                     `json:"book"`
	CoinedEra            string                     `json:"coined_era"`
	CoinedYear           string                     `json:"coined_year"`
	Creator              []string                   `json:"creator"`
	KuData               map[string]int             `json:"ku_data,omitempty"`
	SeeAlso              []string                   `json:"see_also"`
	Resources            *WordResources             `json:"resources,omitempty"`
	Representations      *WordRepresentations       `json:"representations,omitempty"`
	SourceLanguage       string                     `json:"source_language"`
	UsageCategory        string                     `json:"usage_category"`
	Word                 string                     `json:"word"`
	Deprecated           bool                       `json:"deprecated"`
	Etymology            []WordEtymology            `json:"etymology"`
	Audio                []WordAudio                `json:"audio"`
	PuVerbatim           *WordPuVerbatim            `json:"pu_verbatim,omitempty"`
	Usage                map[string]int             `json:"usage"`
	Translations         map[string]WordTranslation `json:"translations"`
}

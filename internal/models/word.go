package models

type WordResources struct {
	SonaPona          *string `json:"sona_pona,omitempty"`
	LipamankaSemantic *string `json:"lipamanka_semantic,omitempty"`
}

type Representations struct {
	SitelenEmosi   *string  `json:"sitelen_emosi,omitempty"`
	SitelenJelo    []string `json:"sitelen_jelo,omitempty"`
	Ligatures      []string `json:"ligatures,omitempty"`
	SitelenSitelen *string  `json:"sitelen_sitelen,omitempty"`
	Ucsur          *string  `json:"ucsur,omitempty"`
}

type Etymology struct {
	Word string `json:"word"`
	Alt  string `json:"alt"`
}

type Audio struct {
	Author string `json:"author"`
	Link   string `json:"link"`
}

type PuVerbatim struct {
	En string `json:"en"`
	Fr string `json:"fr"`
	De string `json:"de"`
	Eo string `json:"eo"`
}

type EtymologyTranslation struct {
	Definition string `json:"definition"`
	Language   string `json:"language"`
}

type Translation struct {
	Commentary  string                 `json:"commentary"`
	Definition  string                 `json:"definition"`
	Etymology   []EtymologyTranslation `json:"etymology"`
	SpEtymology string                 `json:"sp_etymology"`
}

type WordData struct {
	Id                   string                 `json:"id"`
	AuthorVerbatim       string                 `json:"author_verbatim"`
	AuthorVerbatimSource string                 `json:"author_verbatim_source"`
	Book                 string                 `json:"book"`
	CoinedEra            string                 `json:"coined_era"`
	CoinedYear           string                 `json:"coined_year"`
	Creator              []string               `json:"creator"`
	KuData               map[string]int         `json:"ku_data,omitempty"`
	SeeAlso              []string               `json:"see_also"`
	Resources            *WordResources         `json:"resources,omitempty"`
	Representations      *Representations       `json:"representations,omitempty"`
	SourceLanguage       string                 `json:"source_language"`
	UsageCategory        string                 `json:"usage_category"`
	Word                 string                 `json:"word"`
	Deprecated           bool                   `json:"deprecated"`
	Etymology            []Etymology            `json:"etymology"`
	Audio                []Audio                `json:"audio"`
	PuVerbatim           *PuVerbatim            `json:"pu_verbatim,omitempty"`
	Usage                map[string]int         `json:"usage"`
	Translations         map[string]Translation `json:"translations"`
}

type LinkuData map[string]WordData

var UsageCategories = map[string]int{
	"core":     4,
	"common":   3,
	"uncommon": 2,
	"obscure":  1,
	"sandbox":  0,
}

// Helper for rendering pu data in templates.
type PuData struct {
	PartOfSpeech string
	Definition   string
}

// Helpers for rendering etymology data in templates.
type EtymologyData struct {
	Source  string
	Entries []EtymologyEntry
}

type EtymologyEntry struct {
	Word       string
	Alt        string
	Definition string
	Language   string
}

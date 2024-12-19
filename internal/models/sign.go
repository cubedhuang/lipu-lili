package models

type SignEtymology struct {
	Language string  `json:"language"`
	Sign     *string `json:"sign,omitempty"`
}

type SignWriting struct {
	Fsw string `json:"fsw"`
	Swu string `json:"swu"`
}

type SignVideo struct {
	Gif *string `json:"gif,omitempty"`
	Mp4 *string `json:"mp4,omitempty"`
}

type SignParameters struct {
	Handshape   *string `json:"handshape,omitempty"`
	Movement    *string `json:"movement,omitempty"`
	Placement   *string `json:"placement,omitempty"`
	Orientation *string `json:"orientation,omitempty"`
}

type SignTranslation struct {
	Parameters SignParameters `json:"parameters"`
	Icons      string         `json:"icons"`
}

type SignData struct {
	Definition   string                     `json:"definition"`
	Id           string                     `json:"id"`
	IsTwoHanded  bool                       `json:"is_two_handed"`
	NewGloss     string                     `json:"new_gloss"`
	OldGloss     string                     `json:"old_gloss"`
	Etymology    []SignEtymology            `json:"etymology"`
	SignWriting  SignWriting                `json:"signwriting"`
	Video        SignVideo                  `json:"video"`
	Translations map[string]SignTranslation `json:"translations"`
}

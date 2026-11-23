package tcgdex

import (
	"net/http"
)

type TCGDex struct {
	client  *http.Client
	baseURL string
}

type TCGDexSerieBrief struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	LogoURL *string `json:"logo,omitempty"`
}

type TCGDexSetBrief struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	LogoURL   *string `json:"logo"`
	SymbolURL *string `json:"symbol"`
	CardCount struct {
		Total    int `json:"total"`
		Official int `json:"official"`
	} `json:"cardCount"`
}

type TCGDexCardBrief struct {
	ID       string  `json:"id"`
	LocalID  string  `json:"localId"`
	Name     string  `json:"name"`
	ImageURL *string `json:"image"`
}

type TCGDexSet struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ReleaseDate string  `json:"releaseDate"`
	LogoURL     *string `json:"logo,omitempty"`
	SymbolURL   *string `json:"symbol,omitempty"`
	TCGOnline   *string `json:"tcgOnline,omitempty"`
	Legal       struct {
		Standard bool `json:"standard"`
		Expanded bool `json:"expanded"`
	} `json:"legal"`
	CardCount struct {
		Total    int `json:"total"`
		Official int `json:"official"`
		Reverse  int `json:"reverse"`
		Holo     int `json:"holo"`
		FirstEd  int `json:"firstEd"`
	} `json:"cardCount"`
	Boosters []TCGDexBooster   `json:"boosters,omitempty"`
	Cards    []TCGDexCardBrief `json:"cards"`
	Serie    TCGDexSerieBrief  `json:"serie"`
}

type TCGDexBooster struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	LogoURL      *string `json:"logo,omitempty"`
	ArtworkFront *string `json:"artwork_front,omitempty"`
	ArtworkBack  *string `json:"artwork_back,omitempty"`
}

type TCGDexCardCommon struct {
	ID          string         `json:"id"`
	LocalID     string         `json:"localId"`
	Name        string         `json:"name"`
	Category    string         `json:"category"`
	UpdatedAt   string         `json:"updated"`
	ImageURL    *string        `json:"image,omitempty"`
	Illustrator *string        `json:"illustrator,omitempty"`
	Rarity      *string        `json:"rarity,omitempty"`
	SetBrief    TCGDexSetBrief `json:"set"`
	Variants    struct {
		Normal  bool `json:"normal"`
		Reverse bool `json:"reverse"`
		Holo    bool `json:"holo"`
		FirstEd bool `json:"firstEdition"`
	} `json:"variants"`
}

type TCGDexCardPokemon struct { // what a trashy interface
	TCGDexCardCommon

	Hp          *int      `json:"hp,omitempty"`
	EvolveFrom  *string   `json:"evolveFrom,omitempty"`
	Description *string   `json:"description,omitempty"`
	Level       *string   `json:"level,omitempty"`
	Stage       *string   `json:"stage,omitempty"`
	Suffix      *string   `json:"suffix,omitempty"`
	DexId       *[]int    `json:"dexId,omitempty"`
	Types       *[]string `json:"types,omitempty"`
	Item        *struct {
		Name   string `json:"name"`
		Effect string `json:"effect"`
	} `json:"item,omitempty"`
}

type TCGDexCardTrainer struct {
	TCGDexCardCommon

	Effect      string `json:"effect"`
	TrainerType string `json:"trainerType"`
}

type TCGDexCardEnergy struct {
	TCGDexCardCommon

	Effect     string `json:"effect"`
	EnergyType string `json:"energyType"`
}

package pokemon

type Stat struct {
	BaseSt int `json:"base_stat"`
	StName struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type Sprites struct {
	FrontDflt string `json:"front_default"`
}

type TypeInfo struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type RawMove struct {
	Move struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"move"`
}

type Pokemon struct {
	Name    string     `json:"name"`
	Stats   []Stat     `json:"stats"`
	Sprites Sprites    `json:"sprites"`
	Types   []TypeInfo `json:"types"`
	Moves   []RawMove  `json:"moves"`
}

type Move struct {
	Name        string `json:"name"`
	Power       int    `json:"power"`
	StaminaCost int    `json:"stamina_cost"`
	Type        string `json:"attack_type"`
}

type Card struct {
	Name    string
	HP      int
	HPMax   int // original max HP
	Stamina int
	Defense int
	Attack  int
	Moves   []Move
	Types   []string
	Sprite  string
}

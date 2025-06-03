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

type Pokemon struct {
	Name    string  `json:"name"`
	Stats   []Stat  `json:"stats"`
	Sprites Sprites `json:"sprites"`
	Type    []TypeInfo `json:"types"`
}

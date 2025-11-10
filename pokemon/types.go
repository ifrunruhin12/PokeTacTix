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
	Name        string
	HP          int
	HPMax       int
	Stamina     int
	Defense     int
	Attack      int
	Speed       int
	Moves       []Move
	Types       []string
	Sprite      string
	Level       int
	XP          int
	IsLegendary bool
	IsMythical  bool
}

// GetCurrentStats calculates current stats based on level for Card
func (c *Card) GetCurrentStats() CardStats {
	levelMultiplier := float64(c.Level - 1)
	
	hp := int(float64(c.HPMax) * (1.0 + levelMultiplier*0.03))
	attack := int(float64(c.Attack) * (1.0 + levelMultiplier*0.02))
	defense := int(float64(c.Defense) * (1.0 + levelMultiplier*0.02))
	speed := int(float64(c.Speed) * (1.0 + levelMultiplier*0.01))
	stamina := speed * 2
	
	return CardStats{
		HP:      hp,
		Attack:  attack,
		Defense: defense,
		Speed:   speed,
		Stamina: stamina,
	}
}

// CardStats represents computed stats based on level
type CardStats struct {
	HP      int
	Attack  int
	Defense int
	Speed   int
	Stamina int
}

package main

type User struct {
	Name        string
	Username    string
	Password    string
	Role        string
	Active      bool
	LastSeen    int64
	PowerPoints bool
	Picked      [][]float64
}

type SpellSetup struct {
	UserId          float64
	SpellsPerDay    []int
	RemainingSpells []int
	PreparedSpells  [][]float64
}

type PowerPointsSetup struct {
	UserId               float64
	TotalPowerPoints     int
	RemainingPowerPoints int
}

type Spell struct {
	Area                 string
	CastingTime          string
	DescriptionHtml      string
	Descriptors          string
	Duration             string
	Effect               string
	Name                 string
	Page                 int
	Rulebook             string
	SavingThrow          string
	School               string
	SpellRange           string
	SpellResistance      string
	Subschool            string
	Target               string
	Components           string
	ArcaneFocusComponent bool
	CorruptComponent     bool
	DivineFocusComponent bool
	MaterialComponent    bool
	MetaBreathComponent  bool
	SomaticComponent     bool
	TrueNameComponent    bool
	VerbalComponent      bool
	XPComponent          bool
	Custom               bool
	Public               bool
	UserId               float64
}

package main

type User struct {
	Id          string     `json:"id,omitempy"`
	Name        string     `json:"name,omitempy"`
	Username    string     `json:"username,omitempy" auth:"username"`
	Password    string     `json:"password,omitempy" auth:"password"`
	Role        string     `json:"role,omitempy"`
	Active      bool       `json:"active" auth:"active"`
	LastSeen    int64      `json:"lastSeen,omitempy"`
	PowerPoints bool       `json:"powerPoints"`
	Picked      [][]string `json:"picked"`
	Theme       string     `json:"theme,omitempy"`
}

type SpellSetup struct {
	Id              string     `json:"id"`
	UserId          string     `json:"userId,omitempty"`
	SpellsPerDay    []int      `json:"spellsPerDay"`
	RemainingSpells []int      `json:"remainingSpells"`
	PreparedSpells  [][]string `json:"preparedSpells"`
}

type PowerPointsSetup struct {
	Id                   string `json:"id"`
	UserId               string `json:"userId"`
	TotalPowerPoints     int    `json:"totalPowerPoints"`
	RemainingPowerPoints int    `json:"remainingPowerPoints"`
}

type Spell struct {
	Id                   string `json:"id"`
	Area                 string `json:"area,omitempty"`
	CastingTime          string `json:"castingTime,omitempty"`
	DescriptionHtml      string `json:"descriptionHtml,omitempty"`
	Descriptors          string `json:"descriptors,omitempty"`
	Duration             string `json:"duration,omitempty"`
	Effect               string `json:"effect,omitempty"`
	Name                 string `json:"name,omitempty"`
	Page                 int    `json:"page,omitempty"`
	Rulebook             string `json:"rulebook,omitempty"`
	SavingThrow          string `json:"savingThrow,omitempty"`
	School               string `json:"school,omitempty"`
	SpellRange           string `json:"spellRange,omitempty"`
	SpellResistance      string `json:"spellResistance,omitempty"`
	Subschool            string `json:"subschool,omitempty"`
	Target               string `json:"target,omitempty"`
	Components           string `json:"components,omitempty"`
	ArcaneFocusComponent bool   `json:"arcaneFocusComponent"`
	CorruptComponent     bool   `json:"corruptComponent"`
	DivineFocusComponent bool   `json:"divineFocusComponent"`
	MaterialComponent    bool   `json:"materialComponent"`
	MetaBreathComponent  bool   `json:"metaBreathComponent"`
	SomaticComponent     bool   `json:"somaticComponent"`
	TrueNameComponent    bool   `json:"trueNameComponent"`
	VerbalComponent      bool   `json:"verbalComponent"`
	XPComponent          bool   `json:"xpComponent"`
	Displays             string `json:"displays,omitempty"`
	AuditoryDisplay      bool   `json:"auditoryDisplay"`
	MaterialDisplay      bool   `json:"materialDisplay"`
	MentalDisplay        bool   `json:"mentalDisplay"`
	OlfactoryDisplay     bool   `json:"olfactoryDisplay"`
	VisualDisplay        bool   `json:"visualDisplay"`
	Custom               bool   `json:"custom"`
	Public               bool   `json:"public"`
	UserId               string `json:"userId,omitempty"`
}

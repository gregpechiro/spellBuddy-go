package main

import (
	"fmt"
	"strconv"
)

func testSpells(count int) {
	for i := 1; i <= count; i++ {
		var spell Spell
		id := strconv.Itoa(i)
		if !db.Get("spell", id, &spell) {
			fmt.Printf("Could not retrieve spell with id %s\n\n", id)
		}
	}
}

func move() int {

	moveUsers()

	moveSpellSetup()

	movePowerPointsSetup()

	i := moveSpells()

	fmt.Println("move complete")
	return i
}

func moveUsers() {
	docs := db2.GetAll("user")
	fmt.Printf("moving %d users...\n", len(docs))
	for _, doc := range docs {
		id := strconv.FormatFloat(doc.Id, 'f', -1, 64)
		user := User{
			Id:          id,
			Name:        doc.Data["Name"].(string),
			Username:    doc.Data["Username"].(string),
			Password:    doc.Data["Password"].(string),
			Role:        doc.Data["Role"].(string),
			Active:      doc.Data["Active"].(bool),
			LastSeen:    int64(doc.Data["LastSeen"].(float64)),
			PowerPoints: doc.Data["PowerPoints"].(bool),
			Theme:       doc.Data["Theme"].(string),
		}
		for _, fl := range doc.Data["Picked"].([]interface{}) {
			var sl []string
			for _, f := range fl.([]interface{}) {
				sl = append(sl, strconv.FormatFloat(f.(float64), 'f', -1, 64))
			}
			user.Picked = append(user.Picked, sl)
		}
		db.Set("user", id, user)
	}
}

func moveSpellSetup() {
	docs := db2.GetAll("spell-setup")
	fmt.Printf("moving %d spell-setup...\n", len(docs))
	for _, doc := range docs {
		id := strconv.FormatFloat(doc.Id, 'f', -1, 64)
		ss := SpellSetup{
			Id:     id,
			UserId: strconv.FormatFloat(doc.Data["UserId"].(float64), 'f', -1, 64),
		}
		for _, i := range doc.Data["SpellsPerDay"].([]interface{}) {
			ss.SpellsPerDay = append(ss.SpellsPerDay, int(i.(float64)))
		}

		for _, i := range doc.Data["RemainingSpells"].([]interface{}) {
			ss.RemainingSpells = append(ss.RemainingSpells, int(i.(float64)))
		}

		for _, fl := range doc.Data["PreparedSpells"].([]interface{}) {
			var sl []string
			if fl != nil {
				for _, f := range fl.([]interface{}) {
					sl = append(sl, strconv.FormatFloat(f.(float64), 'f', -1, 64))
				}
			}
			ss.PreparedSpells = append(ss.PreparedSpells, sl)
		}
		db.Set("spell-setup", id, ss)
	}
}

func movePowerPointsSetup() {
	docs := db2.GetAll("pp-setup")
	fmt.Printf("moving %d pp-setup...\n", len(docs))
	for _, doc := range docs {
		id := strconv.FormatFloat(doc.Id, 'f', -1, 64)
		pp := PowerPointsSetup{
			Id:                   id,
			UserId:               strconv.FormatFloat(doc.Data["UserId"].(float64), 'f', -1, 64),
			TotalPowerPoints:     int(doc.Data["TotalPowerPoints"].(float64)),
			RemainingPowerPoints: int(doc.Data["RemainingPowerPoints"].(float64)),
		}
		db.Set("pp-setup", id, pp)
	}
}

func moveSpells() int {
	docs := db2.GetAll("spell")
	fmt.Printf("moving %d spells...\n", len(docs))
	i := 0
	for _, doc := range docs {
		id := strconv.FormatFloat(doc.Id, 'f', -1, 64)
		spell := Spell{
			Id:                   id,
			Area:                 doc.Data["Area"].(string),
			CastingTime:          doc.Data["CastingTime"].(string),
			DescriptionHtml:      doc.Data["DescriptionHtml"].(string),
			Descriptors:          doc.Data["Descriptors"].(string),
			Duration:             doc.Data["Duration"].(string),
			Effect:               doc.Data["Effect"].(string),
			Name:                 doc.Data["Name"].(string),
			Page:                 int(doc.Data["Page"].(float64)),
			Rulebook:             doc.Data["Rulebook"].(string),
			SavingThrow:          doc.Data["SavingThrow"].(string),
			School:               doc.Data["School"].(string),
			SpellRange:           doc.Data["SpellRange"].(string),
			SpellResistance:      doc.Data["SpellResistance"].(string),
			Subschool:            doc.Data["Subschool"].(string),
			Target:               doc.Data["Target"].(string),
			Components:           doc.Data["Components"].(string),
			ArcaneFocusComponent: doc.Data["ArcaneFocusComponent"].(bool),
			CorruptComponent:     doc.Data["CorruptComponent"].(bool),
			DivineFocusComponent: doc.Data["DivineFocusComponent"].(bool),
			MaterialComponent:    doc.Data["MaterialComponent"].(bool),
			MetaBreathComponent:  doc.Data["MetaBreathComponent"].(bool),
			SomaticComponent:     doc.Data["SomaticComponent"].(bool),
			TrueNameComponent:    doc.Data["TrueNameComponent"].(bool),
			VerbalComponent:      doc.Data["VerbalComponent"].(bool),
			XPComponent:          doc.Data["XPComponent"].(bool),
			Displays:             doc.Data["Displays"].(string),
			AuditoryDisplay:      doc.Data["AuditoryDisplay"].(bool),
			MaterialDisplay:      doc.Data["MaterialDisplay"].(bool),
			MentalDisplay:        doc.Data["MentalDisplay"].(bool),
			OlfactoryDisplay:     doc.Data["OlfactoryDisplay"].(bool),
			VisualDisplay:        doc.Data["VisualDisplay"].(bool),
			Custom:               doc.Data["Custom"].(bool),
			Public:               doc.Data["Public"].(bool),
			UserId:               strconv.FormatFloat(doc.Data["UserId"].(float64), 'f', -1, 64),
		}
		if id == "610" {
			spell.DescriptionHtml = "\u003cp\u003eYou bring forth the subject\u0026#39;s inner sins and crimes, causing them to manifest in its appearance and aura.\u003cbr /\u003eYour successful touch attack leaves a mystical mark upon the subject.\u003cbr /\u003eAfter a number of rounds equal to your divine caster level, the subject is entitled to a Will save.\u003cbr /\u003eSuccess ends the spell at that point, but failure renders the mark of sin permanent.\u003cbr /\u003eThough the mark is invisible, all living creatures can sense its presence and are repulsed by it.\u003cbr /\u003eThus, they begin their initial interactions with the subject one step nearer to a hostile attitude than they normally would, unless they already know the subject personally.\u003cbr /\u003eFurthermore, the subject takes a -10 circumstance penalty on all Diplomacy checks designed to change the attitudes of others.\u003cbr /\u003e(See Diplomacy, PH 71).\u003cbr /\u003eIn addition, the subject takes a -4 penalty to a specific ability score based on your deity, as given in the table for the divine retribution spell (page 119).\u003cbr /\u003eThis penalty cannot be removed in any way as long as the mark of sin remains, if you do not worship a deity, you must choose one whose alignment is within one step of your own when you cast this spell for the first time.\u003cbr /\u003eThis choice is for the purpose of this effect only, and you cannot subsequently change it unless your alignment shifts in such a way that your previous choice is no longer applicable.\u003cbr /\u003eA mark of sin cannot be dispelled, but it can be removed with a break enchantment, limited wish,miracle, remove curse,\u003cbr /\u003eor wish spell.\u003cbr /\u003eRemove curse works only if its caster level is equal to or higher than that of the mark of sin.\u003c/p\u003e"
		}
		if id == "2900" {
			spell.DescriptionHtml = "This spell makes certain other spells permanent. Depending on the spell, you must be of a minimum caster level and must expend a number of XP.<br><br>You can make the following spells permanent in regard to yourself.<br><span style=\"font-family:monospace\"><br><b>Spell</b>                   <b>Minimum Caster Level</b>      <b>XP Cost</b><br>Arcane sight	                11th              1,500 XP<br>Comprehend languages            9th                 500 XP<br>Darkvision                      10th              1,000 XP<br>Detect magic	                9th                 500 XP<br>Read magic                      9th                 500 XP<br>See invisibility                10th              1,000 XP<br>Tongues	                        11th              1,500 XP<br></span><br>You cast the desired spell and then follow it with the permanency spell. You cannot cast these spells on other creatures. This application of permanency can be dispelled only by a caster of higher level than you were when you cast the spell.<br><br>In addition to personal use, permanency can be used to make the following spells permanent on yourself, another creature, or an object (as appropriate).<br><span style=\"font-family:monospace\"><br><b>Spell</b>                   <b>Minimum Caster Level</b>      <b>XP Cost</b><br>Enlarge person                  9th                 500 XP<br>Magic fang                      9th                 500 XP<br>Magic fang, greater             11th              1,500 XP<br>Rary's telepathic bond*         13th              2,500 XP<br>Reduce person                   9th                 500 XP<br>Resistance                      9th                 500 XP<br></span><br>*Only bonds two creatures per casting of permanency.<br><br>Additionally, the following spells can be cast upon objects or areas only and rendered permanent.<br><span style=\"font-family:monospace\"><br><b>Spell</b>                   <b>Minimum Caster Level</b>      <b>XP Cost</b><br>Alarm                           9th                 500 XP<br>Animate objects                 14th              3,000 XP<br>Dancing lights                  9th                 500 XP<br>Ghost sound                     9th                 500 XP<br>Gust of wind                    11th              1,500 XP<br>Invisibility                    10th              1,000 XP<br>Magic mouth                     10th              1,000 XP<br>Mordenkainen's private sanctum	13th              2,500 XP<br>Phase door                      15th              3,500 XP<br>Prismatic sphere                17th              4,500 XP<br>Prismatic wall                  16th              4,000 XP<br>Shrink item                     11th              1,500 XP<br>Solid fog                       12th              2,000 XP<br>Stinking cloud                  11th              1,500 XP<br>Symbol of death                 16th              4,000 XP<br>Symbol of fear                  14th              3,000 XP<br>Symbol of insanity              16th              4,000 XP<br>Symbol of pain                  13th              2,500 XP<br>Symbol of persuasion            14th              3,000 XP<br>Symbol of sleep                 16th              4,000 XP<br>Symbol of stunning              15th              3,500 XP<br>Symbol of weakness              15th              3,500 XP<br>Teleportation circle            17th              4,500 XP<br>Wall of fire                    12th              2,000 XP<br>Wall of force                   13th              2,500 XP<br>Web                             10th              1,000 XP<br></span><br>Spells cast on other creatures, objects, or locations (not on you) are vulnerable to dispel magic as normal.<br><br>The DM may allow other selected spells to be made permanent. Researching this possible application of a spell costs as much time and money as independently researching the selected spell (see the Dungeon Master's Guide for details). If the DM has already determined that the application is not possible, the research automatically fails. Note that you never learn what is possible except by the success or failure of your research.<br><br>XP Cost: See tables above."
		}
		if db.Set("spell", id, spell) {
			i++
		}
	}
	return i
}

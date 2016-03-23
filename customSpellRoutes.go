package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cagnosolutions/web"
)

var addSpell = web.Route{"GET", "/add/spell", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	tmpl.Render(w, r, "addSpell.tmpl", web.Model{
		"user": db.Get("user", userId),
	})
}}

var editSpell = web.Route{"GET", "/edit/spell/:id", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	spellId := ParseId(r.FormValue(":id"))
	if spellId < 1 {
		web.SetErrorRedirect(w, r, "/setup?cat=userc", "Invalid Spell")
		return
	}
	var spell Spell
	db.Get("spell", spellId).As(&spell)
	if spell.UserId != userId || !spell.Custom {
		web.SetErrorRedirect(w, r, "/setup?cat=userc", "Cannot edit spell")
		return
	}
	tmpl.Render(w, r, "addSpell.tmpl", web.Model{
		"spell":   spell,
		"spellId": spellId,
		"user":    db.Get("user", userId),
	})
	return
}}

var saveSpell = web.Route{"POST", "/save/spell", func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	spellId := ParseId(r.FormValue("id"))
	userId := ParseId(web.GetSess(r, "id"))
	var spell Spell
	if spellId > 0 {
		db.Get("spell", spellId).As(&spell)
	}

	FormToStruct(&spell, r.Form, "")
	c := []string{}
	if spell.ArcaneFocusComponent {
		c = append(c, "AF")
	}
	if spell.CorruptComponent {
		c = append(c, "C")
	}
	if spell.DivineFocusComponent {
		c = append(c, "DF")
	}
	if spell.MaterialComponent {
		c = append(c, "M")
	}
	if spell.MetaBreathComponent {
		c = append(c, "MB")
	}
	if spell.SomaticComponent {
		c = append(c, "S")
	}
	if spell.TrueNameComponent {
		c = append(c, "TN")
	}
	if spell.VerbalComponent {
		c = append(c, "V")
	}
	if spell.XPComponent {
		c = append(c, "XP")
	}
	spell.Components = strings.Join(c, " ")

	d := []string{}
	if spell.AuditoryDisplay {
		d = append(d, "Auditory")
	}
	if spell.MaterialDisplay {
		d = append(d, "Material")
	}
	if spell.MentalDisplay {
		d = append(d, "Mental")
	}
	if spell.OlfactoryDisplay {
		d = append(d, "Olfactory")
	}
	if spell.VisualDisplay {
		d = append(d, "Visual")
	}
	spell.Displays = strings.Join(d, " ")

	if spellId != 0 {
		if spell.UserId != userId {
			web.SetErrorRedirect(w, r, "/setup/?cat=userc", "Cannot save spell")
			fmt.Printf("Spell: %v\n", spell)
			fmt.Printf("UserId: %v\n", userId)
			return
		}
		db.Set("spell", spellId, spell)
	} else {
		spell.UserId = userId
		spell.Custom = true
		db.Add("spell", spell)
	}
	web.SetSuccessRedirect(w, r, "/setup?cat=userc", "Successfully saved spell")
	return
}}

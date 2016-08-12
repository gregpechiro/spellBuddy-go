package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cagnosolutions/web"
)

var addSpell = web.Route{"GET", "/add/spell", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id").(string)
	var user User
	db.Get("user", userId, &user)
	tmpl.Render(w, r, "addSpell.tmpl", web.Model{
		"user": user,
	})
}}

var editSpell = web.Route{"GET", "/edit/spell/:id", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id").(string)
	spellId := r.FormValue(":id")
	// if spellId < 1 {
	// 	web.SetErrorRedirect(w, r, "/setup?cat=userc", "Invalid Spell")
	// 	return
	// }
	var spell Spell
	//db2.Get("spell", spellId).As(&spell)

	if !db.Get("spell", spellId, &spell) {
		web.SetErrorRedirect(w, r, "/setup?cat=userc", "Invalid Spell")
		return
	}

	if spell.UserId != userId || !spell.Custom {
		web.SetErrorRedirect(w, r, "/setup?cat=userc", "Cannot edit spell")
		return
	}
	var user User
	db.Get("user", userId, &user)
	tmpl.Render(w, r, "addSpell.tmpl", web.Model{
		"spell":   spell,
		"spellId": spellId,
		//"user":    db2.Get("user", userId),
		"user": user,
	})
	return
}}

var saveSpell = web.Route{"POST", "/save/spell", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id").(string)
	spellId := r.FormValue("id")
	var spell Spell
	db.Get("spell", spellId, &spell)
	//if spellId != "" {
	//	db2.Get("spell", spellId).As(&spell)
	//}

	if spell.Id == "" {
		spell.Id = genId()
		spell.UserId = userId
		spell.Custom = true
	} else if spell.UserId != userId {
		web.SetErrorRedirect(w, r, "/setup/?cat=userc", "Cannot save spell")
		fmt.Printf("Spell: %v\n", spell)
		fmt.Printf("UserId: %v\n", userId)
		return
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
	spell.Displays = strings.Join(d, ", ")

	// if spellId != 0 {
	// 	if spell.UserId != userId {
	// 		web.SetErrorRedirect(w, r, "/setup/?cat=userc", "Cannot save spell")
	// 		fmt.Printf("Spell: %v\n", spell)
	// 		fmt.Printf("UserId: %v\n", userId)
	// 		return
	// 	}
	// 	db2.Set("spell", spellId, spell)
	// } else {
	// 	spell.UserId = userId
	// 	spell.Custom = true
	// 	db2.Add("spell", spell)
	// }
	db.Set("spell", spell.Id, spell)

	web.SetSuccessRedirect(w, r, "/setup?cat=userc", "Successfully saved spell")
	return
}}

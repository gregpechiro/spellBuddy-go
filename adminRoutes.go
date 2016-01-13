package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cagnosolutions/dbdb"
	"github.com/cagnosolutions/web"
)

var adminHome = web.Route{"GET", "/admin", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "admin.tmpl", web.Model{
		"users": db.GetAll("user"),
	})
}}

var adminUser = web.Route{"GET", "/admin/:id", func(w http.ResponseWriter, r *http.Request) {
	id := ParseId(r.FormValue(":id"))
	if id == 0 {
		http.Redirect(w, r, "/admin", 303)
		return
	}
	tmpl.Render(w, r, "admin.tmpl", web.Model{
		"users": db.GetAll("user"),
		"user":  db.Get("user", id),
	})
}}

var addUser = web.Route{"POST", "/addUser", func(w http.ResponseWriter, r *http.Request) {
	docs := db.Query("user", dbdb.Eq{"username", r.FormValue("username")})
	if len(docs) > 0 {
		web.SetErrorRedirect(w, r, "/admin", "Error adding user. Username already exists")
		return
	}
	user := User{
		Name:        r.FormValue("name"),
		Username:    r.FormValue("username"),
		Password:    r.FormValue("password"),
		Role:        r.FormValue("role"),
		Active:      GetBool(r.FormValue("active")),
		PowerPoints: GetBool(r.FormValue("powerPoints")),
		Picked: [][]float64{
			make([]float64, 0),
			make([]float64, 0),
			make([]float64, 0),
			make([]float64, 0),
			make([]float64, 0),
			make([]float64, 0),
			make([]float64, 0),
			make([]float64, 0),
			make([]float64, 0),
			make([]float64, 0),
		},
	}
	userId := db.Add("user", user)
	if user.PowerPoints {
		ppSetup := PowerPointsSetup{
			UserId:               userId,
			TotalPowerPoints:     0,
			RemainingPowerPoints: 0,
		}
		db.Add("pp-setup", ppSetup)
	} else {
		spellSetup := SpellSetup{
			UserId:          userId,
			SpellsPerDay:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			RemainingSpells: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			PreparedSpells:  make([][]float64, 10),
		}
		db.Add("spell-setup", spellSetup)
	}

	web.SetSuccessRedirect(w, r, "/admin", "Successfully added user")
	return
}}

var saveUser = web.Route{"POST", "/saveUser/:id", func(w http.ResponseWriter, r *http.Request) {
	id := ParseId(r.FormValue(":id"))
	if id == 0 {
		web.SetErrorRedirect(w, r, "/admin", "Error saving user")
		return
	}
	docs := db.Query("user", dbdb.Eq{"Username", r.FormValue("username")}, dbdb.Ne{"::Id::", id})
	if len(docs) > 0 {
		web.SetErrorRedirect(w, r, "/admin", "Error adding user. Username already exists")
		return
	}
	var user User
	db.Get("user", id).As(&user)
	user.Name = r.FormValue("name")
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Role = r.FormValue("role")
	user.Active = GetBool(r.FormValue("active"))
	user.PowerPoints = GetBool(r.FormValue("powerPoints"))
	db.Set("user", id, user)
	web.SetSuccessRedirect(w, r, "/admin", "Successfully saved user")
	return
}}

var delUser = web.Route{"POST", "/delUser/:id", func(w http.ResponseWriter, r *http.Request) {
	id := ParseId(r.FormValue(":id"))
	if id == 0 {
		web.SetSuccessRedirect(w, r, "/admin", "Error saving user")
		return
	}
	db.Del("user", id)
	doc := db.Query("spell-setup", dbdb.Eq{"UserId", id}).One()
	if doc.Id != 0 {
		db.Del("spell-setup", doc.Id)
	}
	web.SetSuccessRedirect(w, r, "/admin", "Successfully deleted user")
	return
}}

var modifiySpells = web.Route{"GET", "/mod/spell", func(w http.ResponseWriter, r *http.Request) {
	docs := db.GetAll("spells")
	for _, doc := range docs {
		spell := Spell{
			Area:            Ternary(doc.Data["area"] == nil, "", doc.Data["area"]).(string),
			CastingTime:     Ternary(doc.Data["casting_time"] == nil, "", doc.Data["casting_time"]).(string),
			DescriptionHtml: Ternary(doc.Data["description_html"] == nil, "", doc.Data["description_html"]).(string),
			Descriptors:     Ternary(doc.Data["descriptors"] == nil, "", doc.Data["descriptors"]).(string),
			Duration:        Ternary(doc.Data["duration"] == nil, "", doc.Data["duration"]).(string),
			Effect:          Ternary(doc.Data["effect"] == nil, "", doc.Data["effect"]).(string),
			Name:            Ternary(doc.Data["name"] == nil, "", doc.Data["name"]).(string),
			Page:            int(Ternary(doc.Data["page"] == nil, float64(0), doc.Data["page"]).(float64)),
			Rulebook:        Ternary(doc.Data["rulebook"] == nil, "", doc.Data["rulebook"]).(string),
			SavingThrow:     Ternary(doc.Data["saving_throw"] == nil, "", doc.Data["saving_throw"]).(string),
			School:          Ternary(doc.Data["school"] == nil, "", doc.Data["school"]).(string),
			SpellRange:      Ternary(doc.Data["spell_range"] == nil, "", doc.Data["spell_range"]).(string),
			SpellResistance: Ternary(doc.Data["spell_resistance"] == nil, "", doc.Data["spell_resistance"]).(string),
			Subschool:       Ternary(doc.Data["subschool"] == nil, "", doc.Data["subschool"]).(string),
			Target:          Ternary(doc.Data["target"] == nil, "", doc.Data["target"]).(string),
		}
		c := []string{}
		if doc.Data["arcane_focus_component"].(string) == "1" {
			c = append(c, "AF")
			spell.ArcaneFocusComponent = true
		} else {
			spell.ArcaneFocusComponent = false
		}
		if doc.Data["corrupt_component"].(string) == "1" {
			c = append(c, "C")
			spell.CorruptComponent = true
		} else {
			spell.CorruptComponent = false
		}
		if doc.Data["divine_focus_component"].(string) == "1" {
			c = append(c, "DF")
			spell.DivineFocusComponent = true
		} else {
			spell.DivineFocusComponent = false
		}
		if doc.Data["material_component"].(string) == "1" {
			c = append(c, "M")
			spell.MaterialComponent = true
		} else {
			spell.MaterialComponent = false
		}
		if doc.Data["meta_breath_component"].(string) == "1" {
			c = append(c, "MB")
			spell.MetaBreathComponent = true
		} else {
			spell.MetaBreathComponent = false
		}
		if doc.Data["somatic_component"].(string) == "1" {
			c = append(c, "S")
			spell.SomaticComponent = true
		} else {
			spell.SomaticComponent = false
		}
		if doc.Data["true_name_component"].(string) == "1" {
			c = append(c, "TN")
			spell.TrueNameComponent = true
		} else {
			spell.TrueNameComponent = false
		}
		if doc.Data["verbal_component"].(string) == "1" {
			c = append(c, "V")
			spell.VerbalComponent = true
		} else {
			spell.VerbalComponent = false
		}
		if doc.Data["xp_component"].(string) == "1" {
			c = append(c, "XP")
			spell.XPComponent = true
		} else {
			spell.XPComponent = false
		}
		spell.Components = strings.Join(c, " ")
		db.Add("spell", spell)
	}
	fmt.Fprintf(w, "Finished Modifying spell. Hopefully it worked...")
}}

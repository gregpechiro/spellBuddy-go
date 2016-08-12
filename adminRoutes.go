package main

import (
	"net/http"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var adminHome = web.Route{"GET", "/admin", func(w http.ResponseWriter, r *http.Request) {
	var users []User
	db.All("user", &users)
	tmpl.Render(w, r, "admin.tmpl", web.Model{
		//"users": db2.GetAll("user"),
		"users": users,
	})
}}

var adminUser = web.Route{"GET", "/admin/:id", func(w http.ResponseWriter, r *http.Request) {
	//id := ParseId(r.FormValue(":id"))
	//if id == 0 {
	//	http.Redirect(w, r, "/admin", 303)
	//	return
	//}

	var user User
	if !db.Get("user", r.FormValue(":id"), &user) {
		web.SetErrorRedirect(w, r, "/admin", "Error finding user")
		return
	}
	var users []User
	db.All("user", &users)
	tmpl.Render(w, r, "admin.tmpl", web.Model{
		//"users": db2.GetAll("user"),
		//"user":  db2.Get("user", id),
		"users": users,
		"user":  user,
	})
}}

var addUser = web.Route{"POST", "/addUser", func(w http.ResponseWriter, r *http.Request) {

	// docs := db2.Query("user", dbdbMod.Eq{"username", r.FormValue("username")})
	// if len(docs) > 0 {
	// 	web.SetErrorRedirect(w, r, "/admin", "Error adding user. Username already exists")
	// 	return
	// }

	var users []User
	db.TestQuery("user", &users, adb.Eq("username", r.FormValue("username")))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/admin", "Error saving user. Username already exists")
		return
	}

	userId := genId()
	user := User{
		Id:          userId,
		Name:        r.FormValue("name"),
		Username:    r.FormValue("username"),
		Password:    r.FormValue("password"),
		Role:        r.FormValue("role"),
		Active:      GetBool(r.FormValue("active")),
		PowerPoints: GetBool(r.FormValue("powerPoints")),
		Picked: [][]string{
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
		},
	}
	//userId := db2.Add("user", user)
	db.Add("user", userId, user)
	setupId := genId()
	if user.PowerPoints {
		ppSetup := PowerPointsSetup{
			Id:                   setupId,
			UserId:               userId,
			TotalPowerPoints:     0,
			RemainingPowerPoints: 0,
		}
		//db2.Add("pp-setup", ppSetup)
		db.Add("pp-setup", setupId, ppSetup)
	} else {
		spellSetup := SpellSetup{
			Id:              setupId,
			UserId:          userId,
			SpellsPerDay:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			RemainingSpells: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			PreparedSpells: [][]string{
				make([]string, 0),
				make([]string, 0),
				make([]string, 0),
				make([]string, 0),
				make([]string, 0),
				make([]string, 0),
				make([]string, 0),
				make([]string, 0),
				make([]string, 0),
				make([]string, 0),
			},
		}
		//db2.Add("spell-setup", spellSetup)
		db.Add("spell-setup", setupId, spellSetup)
	}

	web.SetSuccessRedirect(w, r, "/admin", "Successfully added user")
	return
}}

var saveUser = web.Route{"POST", "/saveUser/:id", func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue(":id")
	if id == "" {
		web.SetErrorRedirect(w, r, "/admin", "Error saving user")
		return
	}

	// docs := db2.Query("user", dbdbMod.Eq{"Username", r.FormValue("username")}, dbdbMod.Ne{"::Id::", id})
	// if len(docs) > 0 {
	// 	web.SetErrorRedirect(w, r, "/admin", "Error adding user. Username already exists")
	// 	return
	// }

	var users []User
	db.TestQuery("company", &users, adb.Eq("username", r.FormValue("username")))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/admin", "Error saving company. Email is already registered")
		return
	}

	var user User
	//db2.Get("user", id).As(&user)

	db.Get("user", "id", &user)

	user.Name = r.FormValue("name")
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Role = r.FormValue("role")
	user.Active = GetBool(r.FormValue("active"))
	user.PowerPoints = GetBool(r.FormValue("powerPoints"))
	//db2.Set("user", id, user)
	db.Set("user", id, user)
	web.SetSuccessRedirect(w, r, "/admin", "Successfully saved user")
	return
}}

var delUser = web.Route{"POST", "/delUser/:id", func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue(":id")
	if id == "" {
		web.SetSuccessRedirect(w, r, "/admin", "Error deleting user")
		return
	}
	//db2.Del("user", id)
	var user User
	db.Get("user", id, &user)

	// doc := db2.Query("spell-setup", dbdbMod.Eq{"UserId", id}).One()
	// if doc.Id != 0 {
	// 	db2.Del("spell-setup", doc.Id)
	// }

	if user.PowerPoints {
		var powerPointsSetup PowerPointsSetup
		if db.TestQueryOne("pp-setup", &powerPointsSetup, adb.Eq("userId", `"`+id+`"`)) {
			db.Del("pp-setup", powerPointsSetup.Id)
		}
	} else {
		var spellSetup SpellSetup
		if db.TestQueryOne("spell-setup", &spellSetup, adb.Eq("userId", `"`+id+`"`)) {
			db.Del("spell-setup", spellSetup.Id)
		}
	}

	db.Del("user", id)

	web.SetSuccessRedirect(w, r, "/admin", "Successfully deleted user")
	return
}}

// var modifiySpells = web.Route{"GET", "/mod/spell", func(w http.ResponseWriter, r *http.Request) {
// 	docs := db2.GetAll("spell")
// 	for _, doc := range docs {
// 		spell := Spell{
// 			Area:            Ternary(doc.Data["area"] == nil, "", doc.Data["area"]).(string),
// 			CastingTime:     Ternary(doc.Data["casting_time"] == nil, "", doc.Data["casting_time"]).(string),
// 			DescriptionHtml: Ternary(doc.Data["description_html"] == nil, "", doc.Data["description_html"]).(string),
// 			Descriptors:     Ternary(doc.Data["descriptors"] == nil, "", doc.Data["descriptors"]).(string),
// 			Duration:        Ternary(doc.Data["duration"] == nil, "", doc.Data["duration"]).(string),
// 			Effect:          Ternary(doc.Data["effect"] == nil, "", doc.Data["effect"]).(string),
// 			Name:            Ternary(doc.Data["name"] == nil, "", doc.Data["name"]).(string),
// 			Page:            int(Ternary(doc.Data["page"] == nil, float64(0), doc.Data["page"]).(float64)),
// 			Rulebook:        Ternary(doc.Data["rulebook"] == nil, "", doc.Data["rulebook"]).(string),
// 			SavingThrow:     Ternary(doc.Data["saving_throw"] == nil, "", doc.Data["saving_throw"]).(string),
// 			School:          Ternary(doc.Data["school"] == nil, "", doc.Data["school"]).(string),
// 			SpellRange:      Ternary(doc.Data["spell_range"] == nil, "", doc.Data["spell_range"]).(string),
// 			SpellResistance: Ternary(doc.Data["spell_resistance"] == nil, "", doc.Data["spell_resistance"]).(string),
// 			Subschool:       Ternary(doc.Data["subschool"] == nil, "", doc.Data["subschool"]).(string),
// 			Target:          Ternary(doc.Data["target"] == nil, "", doc.Data["target"]).(string),
// 		}
// 		c := []string{}
// 		if doc.Data["arcane_focus_component"].(string) == "1" {
// 			c = append(c, "AF")
// 			spell.ArcaneFocusComponent = true
// 		} else {
// 			spell.ArcaneFocusComponent = false
// 		}
// 		if doc.Data["corrupt_component"].(string) == "1" {
// 			c = append(c, "C")
// 			spell.CorruptComponent = true
// 		} else {
// 			spell.CorruptComponent = false
// 		}
// 		if doc.Data["divine_focus_component"].(string) == "1" {
// 			c = append(c, "DF")
// 			spell.DivineFocusComponent = true
// 		} else {
// 			spell.DivineFocusComponent = false
// 		}
// 		if doc.Data["material_component"].(string) == "1" {
// 			c = append(c, "M")
// 			spell.MaterialComponent = true
// 		} else {
// 			spell.MaterialComponent = false
// 		}
// 		if doc.Data["meta_breath_component"].(string) == "1" {
// 			c = append(c, "MB")
// 			spell.MetaBreathComponent = true
// 		} else {
// 			spell.MetaBreathComponent = false
// 		}
// 		if doc.Data["somatic_component"].(string) == "1" {
// 			c = append(c, "S")
// 			spell.SomaticComponent = true
// 		} else {
// 			spell.SomaticComponent = false
// 		}
// 		if doc.Data["true_name_component"].(string) == "1" {
// 			c = append(c, "TN")
// 			spell.TrueNameComponent = true
// 		} else {
// 			spell.TrueNameComponent = false
// 		}
// 		if doc.Data["verbal_component"].(string) == "1" {
// 			c = append(c, "V")
// 			spell.VerbalComponent = true
// 		} else {
// 			spell.VerbalComponent = false
// 		}
// 		if doc.Data["xp_component"].(string) == "1" {
// 			c = append(c, "XP")
// 			spell.XPComponent = true
// 		} else {
// 			spell.XPComponent = false
// 		}
// 		spell.Components = strings.Join(c, " ")
// 		db2.Set("spell", doc.Id, spell)
// 	}
// 	fmt.Fprintf(w, "Finished Modifying spell. Hopefully it worked...")
// }}
//
// var updateSpells = web.Route{"GET", "/update/spells", func(w http.ResponseWriter, r *http.Request) {
// 	docs := db2.GetAll("spell")
// 	for _, doc := range docs {
// 		var spell Spell
// 		doc.As(&spell)
// 		db2.Set("spell", doc.Id, spell)
// 	}
// 	web.SetSuccessRedirect(w, r, "/admin", "Successfully updated spells")
// 	return
// }}

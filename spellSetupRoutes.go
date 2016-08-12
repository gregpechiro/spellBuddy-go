package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var setup = web.Route{"GET", "/setup", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id").(string)
	var user User
	//db2.Get("user", userId).As(&user)
	db.Get("user", userId, &user)
	var spellCat string
	var spells []Spell
	cat := r.FormValue("cat")

	if cat == "allc" {
		//spells = db2.Query("spell", dbdbMod.Eq{"Custom", true}, dbdbMod.Eq{"Public", true})
		db.TestQuery("spell", &spells, adb.Eq("custom", "true"), adb.Eq("public", "true"))
		spellCat = "allC"
	} else if cat == "userc" {
		//spells = db2.Query("spell", dbdbMod.Eq{"UserId", userId})
		db.TestQuery("spell", &spells, adb.Eq("userId", `"`+userId+`"`))
		spellCat = "userC"
	} else {
		//spells = db2.Query("spell", dbdbMod.Eq{"Custom", false})
		db.TestQuery("spell", &spells, adb.Eq("custom", "false"))
		spellCat = "dndtool"
	}
	var setup map[string]interface{}
	if user.PowerPoints {
		//setup = db2.Query("pp-setup", dbdbMod.Eq{"UserId", userId}).One()
		db.TestQueryOne("pp-setup", &setup, adb.Eq("userId", `"`+userId+`"`))
	} else {
		//setup = db2.Query("spell-setup", dbdbMod.Eq{"UserId", userId}).One()
		db.TestQueryOne("spell-setup", &setup, adb.Eq("userId", `"`+userId+`"`))

	}
	tmpl.Render(w, r, "setup.tmpl", web.Model{
		"picked":   getPickedNames(user.Picked),
		"user":     user,
		"setup":    setup,
		"spells":   spells,
		"spellCat": spellCat,
	})
}}

var spellsPerDay = web.Route{"POST", "/user/spd", func(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userId")
	setupId := r.FormValue("setupId")
	var spellSetup SpellSetup
	//db2.Get("spell-setup", setupId).As(&spellSetup)
	db.Get("spell-setup", setupId, &spellSetup)
	if userId != spellSetup.UserId {
		web.SetErrorRedirect(w, r, "/setup", "Error updating spells per day")
		return
	}
	for i := range spellSetup.SpellsPerDay {
		spd, _ := strconv.Atoi(r.FormValue(fmt.Sprintf("level-%d", i)))
		spellSetup.SpellsPerDay[i] = spd
	}
	//db2.Set("spell-setup", setupId, spellSetup)
	db.Set("spell-setup", setupId, spellSetup)
	web.SetSuccessRedirect(w, r, "/setup", "Successfully updated spells per day")
	return
}}

var pp = web.Route{"POST", "/user/pp", func(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userId")
	setupId := r.FormValue("setupId")
	var ppSetup PowerPointsSetup
	//db2.Get("pp-setup", setupId).As(&ppSetup)
	db.Get("pp-setup", setupId, &ppSetup)
	if userId != ppSetup.UserId {
		web.SetErrorRedirect(w, r, "/setup", "Error updating total power points")
		return
	}
	totalPP, _ := strconv.Atoi(r.FormValue("totalPP"))
	ppSetup.TotalPowerPoints = totalPP
	//db2.Set("pp-setup", setupId, ppSetup)
	db.Set("pp-setup", setupId, ppSetup)
	web.SetSuccessRedirect(w, r, "/setup", "Successfully updated total power points")
	return
}}

var addSpellToUser = web.Route{"POST", "/user/addSpell", func(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	userId := r.FormValue("userId")
	spellId := r.FormValue("spellId")
	spellLvl, err := strconv.ParseInt(r.FormValue("spellLvl"), 10, 64)
	if err != nil {
		response["success"] = false
		response["msg"] = "Error adding spell"
		ajaxResponse(w, response)
		return
	}
	if spellLvl > 9 {
		response["success"] = false
		response["msg"] = "Spell level must be between 0 and 9"
		ajaxResponse(w, response)
		return
	}
	if spellLvl < 0 {
		response["success"] = false
		response["msg"] = "Spell level cannot be negative"
		ajaxResponse(w, response)
		return
	}

	var user User
	//db2.Get("user", userId).As(&user)
	db.Get("user", userId, &user)
	if user.Picked == nil {
		user.Picked = make([][]string, 10)
	}
	pickedLvl := user.Picked[spellLvl]
	pickedLvl = append(pickedLvl, spellId)
	user.Picked[spellLvl] = pickedLvl
	//db2.Set("user", userId, user)
	db.Set("user", userId, user)

	response["success"] = true
	response["msg"] = "Successfully added spell"
	response["picked"] = getPickedNames(user.Picked)
	ajaxResponse(w, response)
	return
}}

var delSpellFromUser = web.Route{"POST", "/user/delSpell", func(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	userId := r.FormValue("userId")
	spellLvl, err := strconv.Atoi(r.FormValue("spellLvl"))
	if err != nil {
		log.Printf("delSpellFfromUser >> spellLvl >> Atoi(): %v\n", err)
		response["success"] = false
		response["msg"] = "Error deleting spell"
		ajaxResponse(w, response)
		return
	}
	idx, err := strconv.Atoi(r.FormValue("idx"))
	if err != nil {
		log.Printf("delSpellFfromUser >> idx >> Atoi(): %v\n", err)
		response["success"] = false
		response["msg"] = "Error deleting spell"
		ajaxResponse(w, response)
		return
	}
	var user User
	//db2.Get("user", userId).As(&user)
	db.Get("user", userId, &user)
	if user.Picked == nil || len(user.Picked) < spellLvl || len(user.Picked[spellLvl]) < idx {
		log.Printf("delSpellFfromUser >> user.Picked size\n")
		response["success"] = false
		response["msg"] = "Error deleting spell"
		ajaxResponse(w, response)
		return
	}
	pickedLvl := user.Picked[spellLvl]
	pickedLvl = append(pickedLvl[:idx], pickedLvl[idx+1:]...)
	user.Picked[spellLvl] = pickedLvl
	//db2.Set("user", userId, user)
	db.Set("user", userId, user)
	if !user.PowerPoints {
		var spellSetup SpellSetup
		//doc := db2.Query("spell-setup", dbdbMod.Eq{"UserId", userId}).One()
		//doc.As(&spellSetup)
		db.TestQueryOne("spell-setup", &spellSetup, adb.Eq("userId", `"`+userId+`"`))
		spellSetup.PreparedSpells[spellLvl] = make([]string, 0)
		//db2.Set("spell-setup", doc.Id, spellSetup)
		db.Set("spell-setup", spellSetup.Id, spellSetup)
	}
	response["success"] = true
	response["msg"] = "Successfully deleted spell"
	response["picked"] = getPickedNames(user.Picked)
	ajaxResponse(w, response)
	return
}}

var changeLvl = web.Route{"POST", "/user/changeLvl", func(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	userId := r.FormValue("userId")
	spellLvl, err := strconv.Atoi(r.FormValue("spellLvl"))
	if err != nil {
		log.Printf("delSpellFfromUser >> spellLvl >> Atoi(): %v\n", err)
		response["success"] = false
		response["msg"] = "Error changing spell level"
		ajaxResponse(w, response)
		return
	}
	idx, err := strconv.Atoi(r.FormValue("idx"))
	if err != nil {
		log.Printf("delSpellFfromUser >> idx >> Atoi(): %v\n", err)
		response["success"] = false
		response["msg"] = "Error changing spell level"
		ajaxResponse(w, response)
		return
	}
	newLvl, err := strconv.Atoi(r.FormValue("newLvl"))
	if err != nil {
		log.Printf("delSpellFfromUser >> idx >> Atoi(): %v\n", err)
		response["success"] = false
		response["msg"] = "Error changing spell level"
		ajaxResponse(w, response)
		return
	}
	var user User
	//db2.Get("user", userId).As(&user)
	db.Get("user", userId, &user)
	if user.Picked == nil || len(user.Picked) < spellLvl || len(user.Picked) < newLvl || len(user.Picked[spellLvl]) < idx {
		log.Printf("delSpellFfromUser >> user.Picked size\n")
		response["success"] = false
		response["msg"] = "Error changing spell level"
		ajaxResponse(w, response)
		return
	}
	spellId := user.Picked[spellLvl][idx]

	oldPicked := user.Picked[spellLvl]
	oldPicked = append(oldPicked[:idx], oldPicked[idx+1:]...)
	user.Picked[spellLvl] = oldPicked

	newPicked := user.Picked[newLvl]
	newPicked = append(newPicked, spellId)
	user.Picked[newLvl] = newPicked
	//db2.Set("user", userId, user)
	db.Set("user", userId, user)
	if !user.PowerPoints {
		var spellSetup SpellSetup
		//doc := db2.Query("spell-setup", dbdbMod.Eq{"UserId", userId}).One()
		//doc.As(&spellSetup)
		db.TestQueryOne("spell-setup", &spellSetup, adb.Eq("userId", `"`+userId+`"`))

		spellSetup.PreparedSpells[spellLvl] = make([]string, 0)
		//db2.Set("spell-setup", doc.Id, spellSetup)
		db.Set("spell-setup", spellSetup.Id, spellSetup)
	}
	response["success"] = true
	response["msg"] = "Successfully added spell"
	response["picked"] = getPickedNames(user.Picked)
	ajaxResponse(w, response)
	return
}}

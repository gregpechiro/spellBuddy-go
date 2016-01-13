package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cagnosolutions/dbdb"
	"github.com/cagnosolutions/web"
)

var setup = web.Route{"GET", "/setup", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	var user User
	db.Get("user", userId).As(&user)
	var spellCat string
	var spells dbdb.DocSorted
	cat := r.FormValue("cat")

	if cat == "allc" {
		spells = db.Query("spell", dbdb.Eq{"Custom", true}, dbdb.Eq{"Public", true})
		spellCat = "allC"
	} else if cat == "userc" {
		spells = db.Query("spell", dbdb.Eq{"UserId", userId})
		spellCat = "userC"
	} else {
		spells = db.Query("spell", dbdb.Eq{"Custom", false})
		spellCat = "dndtool"
	}
	var setup *dbdb.Doc
	if user.PowerPoints {
		setup = db.Query("pp-setup", dbdb.Eq{"UserId", userId}).One()
	} else {
		setup = db.Query("spell-setup", dbdb.Eq{"UserId", userId}).One()
	}
	tmpl.Render(w, r, "setup.tmpl", web.Model{
		"picked":   getPickedNames(user.Picked),
		"user":     db.Get("user", userId),
		"setup":    setup,
		"spells":   spells,
		"spellCat": spellCat,
	})
}}

var spellsPerDay = web.Route{"POST", "/user/spd", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(r.FormValue("userId"))
	setupId := ParseId(r.FormValue("setupId"))
	var spellSetup SpellSetup
	db.Get("spell-setup", setupId).As(&spellSetup)
	if userId != spellSetup.UserId {
		web.SetErrorRedirect(w, r, "/setup", "Error updating spells per day")
		return
	}
	for i := range spellSetup.SpellsPerDay {
		spd, _ := strconv.Atoi(r.FormValue(fmt.Sprintf("level-%d", i)))
		spellSetup.SpellsPerDay[i] = spd
	}
	db.Set("spell-setup", setupId, spellSetup)
	web.SetSuccessRedirect(w, r, "/setup", "Successfully updated spells per day")
	return
}}

var pp = web.Route{"POST", "/user/pp", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(r.FormValue("userId"))
	setupId := ParseId(r.FormValue("setupId"))
	var ppSetup PowerPointsSetup
	db.Get("pp-setup", setupId).As(&ppSetup)
	if userId != ppSetup.UserId {
		web.SetErrorRedirect(w, r, "/setup", "Error updating total power points")
		return
	}
	totalPP, _ := strconv.Atoi(r.FormValue("totalPP"))
	ppSetup.TotalPowerPoints = totalPP
	db.Set("pp-setup", setupId, ppSetup)
	web.SetSuccessRedirect(w, r, "/setup", "Successfully updated total power points")
	return
}}

var addSpellToUser = web.Route{"POST", "/user/addSpell", func(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	userId := ParseId(r.FormValue("userId"))
	spellId := ParseId(r.FormValue("spellId"))
	spellLvl, err := strconv.ParseInt(r.FormValue("spellLvl"), 10, 64)
	if err != nil || spellLvl > 9 || spellLvl < 0 {
		response["success"] = false
		response["msg"] = "Error adding spell"
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, "%s", b)
		return
	}
	var user User
	db.Get("user", userId).As(&user)
	if user.Picked == nil {
		user.Picked = make([][]float64, 0)
	}
	pickedLvl := user.Picked[spellLvl]
	pickedLvl = append(pickedLvl, spellId)
	user.Picked[spellLvl] = pickedLvl
	db.Set("user", userId, user)

	response["success"] = true
	response["msg"] = "Successfully added spell"
	response["picked"] = getPickedNames(user.Picked)
	b, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", b)
	return
}}

var delSpellFromUser = web.Route{"POST", "/user/delSpell", func(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	userId := ParseId(r.FormValue("userId"))
	spellLvl, err := strconv.Atoi(r.FormValue("spellLvl"))
	if err != nil {
		log.Printf("delSpellFfromUser >> spellLvl >> Atoi(): %v\n", err)
		response["success"] = false
		response["msg"] = "Error deleting spell"
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, "%s", b)
		return
	}
	idx, err := strconv.Atoi(r.FormValue("idx"))
	if err != nil {
		log.Printf("delSpellFfromUser >> idx >> Atoi(): %v\n", err)
		response["success"] = false
		response["msg"] = "Error deleting spell"
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, "%s", b)
		return
	}
	var user User
	db.Get("user", userId).As(&user)
	if user.Picked == nil || len(user.Picked) < spellLvl || len(user.Picked[spellLvl]) < idx {
		log.Printf("delSpellFfromUser >> user.Picked size\n")
		response["success"] = false
		response["msg"] = "Error deleting spell"
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, "%s", b)
		return
	}
	pickedLvl := user.Picked[spellLvl]
	pickedLvl = append(pickedLvl[:idx], pickedLvl[idx+1:]...)
	user.Picked[spellLvl] = pickedLvl
	db.Set("user", userId, user)
	response["success"] = true
	response["msg"] = "Successfully added spell"
	response["picked"] = getPickedNames(user.Picked)
	b, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", b)
	return
}}

var changeLvl = web.Route{"POST", "/user/changeLvl", func(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	userId := ParseId(r.FormValue("userId"))
	spellLvl, err := strconv.Atoi(r.FormValue("spellLvl"))
	if err != nil {
		log.Printf("delSpellFfromUser >> spellLvl >> Atoi(): %v\n", err)
		response["success"] = false
		response["msg"] = "Error changing spell level"
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, "%s", b)
		return
	}
	idx, err := strconv.Atoi(r.FormValue("idx"))
	if err != nil {
		log.Printf("delSpellFfromUser >> idx >> Atoi(): %v\n", err)
		response["success"] = false
		response["msg"] = "Error changing spell level"
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, "%s", b)
		return
	}
	newLvl, err := strconv.Atoi(r.FormValue("newLvl"))
	if err != nil {
		log.Printf("delSpellFfromUser >> idx >> Atoi(): %v\n", err)
		response["success"] = false
		response["msg"] = "Error changing spell level"
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, "%s", b)
		return
	}
	var user User
	db.Get("user", userId).As(&user)
	if user.Picked == nil || len(user.Picked) < spellLvl || len(user.Picked) < newLvl || len(user.Picked[spellLvl]) < idx {
		log.Printf("delSpellFfromUser >> user.Picked size\n")
		response["success"] = false
		response["msg"] = "Error changing spell level"
		b, _ := json.Marshal(response)
		fmt.Fprintf(w, "%s", b)
		return
	}
	spellId := user.Picked[spellLvl][idx]

	oldPicked := user.Picked[spellLvl]
	oldPicked = append(oldPicked[:idx], oldPicked[idx+1:]...)
	user.Picked[spellLvl] = oldPicked

	newPicked := user.Picked[newLvl]
	newPicked = append(newPicked, spellId)
	user.Picked[newLvl] = newPicked
	db.Set("user", userId, user)
	response["success"] = true
	response["msg"] = "Successfully added spell"
	response["picked"] = getPickedNames(user.Picked)
	b, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", b)
	return
}}

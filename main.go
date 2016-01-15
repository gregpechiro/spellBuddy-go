package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cagnosolutions/dbdb"
	"github.com/cagnosolutions/web"
)

var mux = web.NewMux()
var db = dbdb.NewDataStore()
var tmpl *web.TmplCache

func init() {

	go dbdb.Serve(db, ":9999", "spell-buddy")
	web.Funcs["title"] = strings.Title
	web.Funcs["json"] = func(v interface{}) string {
		b, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
		}
		return string(b)
	}
	web.Funcs["isIn"] = func(src []interface{}, target interface{}) bool {
		for _, i := range src {
			if target == i {
				return true
			}
		}
		return false
	}
	web.Funcs["lenEq"] = func(src []interface{}, target interface{}) bool {
		switch target.(type) {
		case int:
			return len(src) == target.(int)
		case int8:
			return len(src) == int(target.(int8))
		case int16:
			return len(src) == int(target.(int16))
		case int32:
			return len(src) == int(target.(int32))
		case int64:
			return len(src) == int(target.(int64))
		case uint:
			return len(src) == int(target.(uint))
		case uint8:
			return len(src) == int(target.(uint8))
		case uint16:
			return len(src) == int(target.(uint16))
		case uint32:
			return len(src) == int(target.(uint32))
		case uint64:
			return len(src) == int(target.(uint64))
		case float32:
			return len(src) == int(target.(float32))
		case float64:
			return len(src) == int(target.(float64))
		}
		return false
	}
	tmpl = web.NewTmplCache()
}

func main() {
	db.AddStore("spell")
	db.AddStore("user")
	db.AddStore("spell-setup")
	db.AddStore("pp-setup")

	// auth routed
	mux.AddRoutes(login, loginPost, logout)

	// admin routes
	mux.AddSecureRoutes(ADMIN, adminHome, adminUser, addUser, saveUser, delUser, modifiySpells)

	// spell routes
	mux.AddSecureRoutes(USER, home, setup, addSpell, saveSpell, editSpell)

	// user account routes
	mux.AddSecureRoutes(USER, account, saveAccount, saveTheme)

	// main app routes
	mux.AddSecureRoutes(USER, addSpellToUser, delSpellFromUser, changeLvl, spellsPerDay, rest, cast)
	mux.AddSecureRoutes(USER, pp, ppRest, ppCast)

	http.ListenAndServe(":8080", mux)
}

var login = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "login.tmpl", web.Model{})
}}

var loginPost = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("username") == "admin" && r.FormValue("password") == "admin" {
		web.Login(w, r, "ADMIN")
		web.SetSuccessRedirect(w, r, "/admin", "Welcome in memory admin")
		return
	}
	docs := db.Query("user", dbdb.Eq{"Username", r.FormValue("username")}, dbdb.Eq{"Password", r.FormValue("password")}, dbdb.Eq{"Active", true})

	if len(docs) == 1 {
		var user User
		docs[0].As(&user)
		sess := web.Login(w, r, user.Role)
		sess["id"] = docs[0].Id

		sess["username"] = user.Username

		web.PutMultiSess(w, r, sess)

		user.LastSeen = time.Now().Unix()
		db.Set("user", docs[0].Id, user)
		web.SetSuccessRedirect(w, r, "/home", "Welcome "+user.Username)
		return
	}
	web.SetErrorRedirect(w, r, "/", "Incorrect username or password")
	return
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w, r)
	web.SetSuccessRedirect(w, r, "/", "See you next time")
}}

var home = web.Route{"GET", "/home", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	var user User
	uDoc := db.Get("user", userId)
	uDoc.As(&user)
	if user.PowerPoints {
		tmpl.Render(w, r, "pp-home.tmpl", web.Model{
			"user":   uDoc,
			"setup":  db.Query("pp-setup", dbdb.Eq{"UserId", userId}).One(),
			"picked": getPicked(user.Picked),
		})
		return
	}
	tmpl.Render(w, r, "home.tmpl", web.Model{
		"user":   uDoc,
		"setup":  db.Query("spell-setup", dbdb.Eq{"UserId", userId}).One(),
		"picked": getPicked(user.Picked),
	})
	return
}}

var rest = web.Route{"POST", "/rest", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	setupId := ParseId(r.FormValue("setupId"))
	var prepared [][]float64
	err := json.Unmarshal([]byte(r.FormValue("prepared")), &prepared)
	var spellSetup SpellSetup
	db.Get("spell-setup", setupId).As(&spellSetup)
	if spellSetup.UserId != userId || err != nil {
		fmt.Println("userId: ", userId)
		fmt.Println("spellId: ", setupId)
		web.SetErrorRedirect(w, r, "/home", "Error Resting")
		return
	}
	copy(spellSetup.RemainingSpells, spellSetup.SpellsPerDay)
	copy(spellSetup.PreparedSpells, prepared)
	db.Set("spell-setup", setupId, spellSetup)
	http.Redirect(w, r, "/home", 303)
	return
}}

var ppRest = web.Route{"POST", "/pp-rest", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	setupId := ParseId(r.FormValue("setupId"))
	var ppSetup PowerPointsSetup
	db.Get("pp-setup", setupId).As(&ppSetup)
	if ppSetup.UserId != userId {
		fmt.Println("userId: ", userId)
		fmt.Println("spellId: ", setupId)
		web.SetErrorRedirect(w, r, "/home", "Error Resting")
		return
	}
	ppSetup.RemainingPowerPoints = ppSetup.TotalPowerPoints
	db.Set("pp-setup", setupId, ppSetup)
	http.Redirect(w, r, "/home", 303)
	return
}}

var cast = web.Route{"POST", "/cast", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	setupId := ParseId(r.FormValue("setupId"))
	level, err := strconv.Atoi(r.FormValue("level"))
	var spellSetup SpellSetup
	db.Get("spell-setup", setupId).As(&spellSetup)
	if spellSetup.UserId != userId || err != nil || level < 0 || level > 9 {
		fmt.Println("userId: ", userId)
		fmt.Println("spellId: ", setupId)
		web.SetErrorRedirect(w, r, "/home", "Error casting")
		return
	}
	spellSetup.RemainingSpells[level]--
	db.Set("spell-setup", setupId, spellSetup)
	http.Redirect(w, r, "/home", 303)
	return

}}

var ppCast = web.Route{"POST", "/pp-cast", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	setupId := ParseId(r.FormValue("setupId"))
	pp, err := strconv.Atoi(r.FormValue("pp"))
	var ppSetup PowerPointsSetup
	db.Get("pp-setup", setupId).As(&ppSetup)
	if ppSetup.UserId != userId || err != nil || pp < 0 || pp > ppSetup.TotalPowerPoints {
		log.Printf("ppCast() >> strconv.Atoi(): %v\n", err)
		fmt.Println("userId: ", userId)
		fmt.Println("setup userId: ", ppSetup.UserId)
		web.SetErrorRedirect(w, r, "/home", "Error casting")
		return
	}
	ppSetup.RemainingPowerPoints -= pp
	db.Set("pp-setup", setupId, ppSetup)
	http.Redirect(w, r, "/home", 303)
	return

}}

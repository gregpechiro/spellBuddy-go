package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var mux = web.NewMux()

//var db2 = dbdbMod.NewDataStore()
var db = adb.NewDB()
var tmpl *web.TmplCache

func init() {
	//go dbdbMod.Serve(db2, ":9999", "spell-buddy")
	//web.NewCookieSalt()
	web.SESSDUR = time.Minute * 60 * 3
	web.Funcs["add"] = func(i, j int) int {
		return i + j
	}
	web.Funcs["title"] = strings.Title
	web.Funcs["json"] = func(v interface{}) string {
		b, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
		}
		return string(b)
	}
	web.Funcs["isIn"] = isIn
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
	//db2.AddStore("spell")
	//db2.AddStore("user")
	//db2.AddStore("spell-setup")
	//db2.AddStore("pp-setup")

	db.AddStore("spell")
	db.AddStore("user")
	db.AddStore("spell-setup")
	db.AddStore("pp-setup")

	//i := move()
	//testSpells(i)

	// auth routed
	mux.AddRoutes(login, loginPost, logout)

	// admin routes
	mux.AddSecureRoutes(ADMIN, adminHome, adminUser, addUser, saveUser, delUser)

	// custom spell routes
	mux.AddSecureRoutes(USER, addSpell, saveSpell, editSpell)

	// user account routes
	mux.AddSecureRoutes(USER, account, saveAccount, saveTheme)

	// spell setup routes
	mux.AddSecureRoutes(USER, setup, addSpellToUser, delSpellFromUser, changeLvl, spellsPerDay)

	// main app routes
	mux.AddSecureRoutes(USER, home, rest, cast, pp, ppRest, ppCast)

	http.ListenAndServe(":8080", mux)
}

var login = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "login.tmpl", web.Model{})
}}

var loginPost = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "admin" && password == "admin" {
		web.Login(w, r, "ADMIN")
		web.SetSuccessRedirect(w, r, "/admin", "Welcome in memory admin")
		return
	}

	var user User
	if !db.Auth("user", username, password, &user) {
		web.SetErrorRedirect(w, r, "/", "Incorrect username or password")
		return
	}
	sess := web.Login(w, r, user.Role)
	sess["id"] = user.Id
	sess["username"] = user.Username
	web.PutMultiSess(w, r, sess)
	user.LastSeen = time.Now().Unix()
	db.Set("user", user.Id, user)
	web.SetSuccessRedirect(w, r, "/home", "Welcome "+user.Username)
	return

	// docs := db2.Query("user", dbdbMod.Eq{"Username", r.FormValue("username")}, dbdbMod.Eq{"Password", r.FormValue("password")}, dbdbMod.Eq{"Active", true})
	//
	// if len(docs) == 1 {
	// 	var user User
	// 	docs[0].As(&user)
	// 	sess := web.Login(w, r, user.Role)
	// 	sess["id"] = docs[0].Id
	//
	// 	sess["username"] = user.Username
	//
	// 	web.PutMultiSess(w, r, sess)
	//
	// 	user.LastSeen = time.Now().Unix()
	// 	db2.Set("user", docs[0].Id, user)
	// 	web.SetSuccessRedirect(w, r, "/home", "Welcome "+user.Username)
	// 	return
	// }
	// web.SetErrorRedirect(w, r, "/", "Incorrect username or password")
	// return
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w)
	web.SetSuccessRedirect(w, r, "/", "See you next time")
}}

var home = web.Route{"GET", "/home", func(w http.ResponseWriter, r *http.Request) {

	// userId := ParseId(web.GetSess(r, "id"))
	// var user User
	// uDoc := db2.Get("user", userId)
	// uDoc.As(&user)
	// if user.PowerPoints {
	// 	tmpl.Render(w, r, "pp-home.tmpl", web.Model{
	// 		"user":   uDoc,
	// 		"setup":  db2.Query("pp-setup", dbdbMod.Eq{"UserId", userId}).One(),
	// 		"picked": getPicked(user.Picked),
	// 	})
	// 	return
	// }
	// tmpl.Render(w, r, "home.tmpl", web.Model{
	// 	"user":   uDoc,
	// 	"setup":  db2.Query("spell-setup", dbdbMod.Eq{"UserId", userId}).One(),
	// 	"picked": getPicked(user.Picked),
	// })
	// return

	var user User
	db.Get("user", web.GetSess(r, "id").(string), &user)
	if user.PowerPoints {
		var powerPointsSetup PowerPointsSetup
		db.TestQueryOne("pp-setup", &powerPointsSetup, adb.Eq("userId", `"`+user.Id+`"`))
		tmpl.Render(w, r, "pp-home.tmpl", web.Model{
			"user":   user,
			"setup":  powerPointsSetup,
			"picked": getPickedPP(user.Picked),
		})
		return
	}
	var spellSetup SpellSetup
	db.TestQueryOne("spell-setup", &spellSetup, adb.Eq("userId", `"`+user.Id+`"`))
	tmpl.Render(w, r, "home.tmpl", web.Model{
		"user":   user,
		"setup":  spellSetup,
		"picked": getPicked(user.Picked, spellSetup.PreparedSpells),
	})
	return

}}

var rest = web.Route{"POST", "/rest", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id")
	setupId := r.FormValue("setupId")
	var prepared [][]string
	err := json.Unmarshal([]byte(r.FormValue("prepared")), &prepared)
	var spellSetup SpellSetup
	//db2.Get("spell-setup", setupId).As(&spellSetup)
	db.Get("spell-setup", setupId, &spellSetup)
	if spellSetup.UserId != userId || err != nil {
		fmt.Println("userId: ", userId)
		fmt.Println("setupId: ", setupId)
		fmt.Println(err)
		web.SetErrorRedirect(w, r, "/home", "Error Resting")
		return
	}
	copy(spellSetup.RemainingSpells, spellSetup.SpellsPerDay)
	copy(spellSetup.PreparedSpells, prepared)
	//db2.Set("spell-setup", setupId, spellSetup)
	db.Set("spell-setup", setupId, spellSetup)
	http.Redirect(w, r, "/home", 303)
	return
}}

var ppRest = web.Route{"POST", "/pp-rest", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id")
	setupId := r.FormValue("setupId")
	var powerPointSetup PowerPointsSetup
	//db2.Get("pp-setup", setupId).As(&ppSetup)
	db.Get("pp-setup", setupId, &powerPointSetup)
	if powerPointSetup.UserId != userId {
		fmt.Println("userId: ", userId)
		fmt.Println("spellId: ", setupId)
		web.SetErrorRedirect(w, r, "/home", "Error Resting")
		return
	}
	powerPointSetup.RemainingPowerPoints = powerPointSetup.TotalPowerPoints
	//db2.Set("pp-setup", setupId, ppSetup)
	db.Set("pp-setup", setupId, powerPointSetup)
	http.Redirect(w, r, "/home", 303)
	return
}}

var cast = web.Route{"POST", "/cast", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id")
	setupId := r.FormValue("setupId")
	level, err := strconv.Atoi(r.FormValue("level"))
	var spellSetup SpellSetup
	//db2.Get("spell-setup", setupId).As(&spellSetup)
	db.Get("spell-setup", setupId, &spellSetup)
	if spellSetup.UserId != userId || err != nil || level < 0 || level > 9 {
		fmt.Println("userId: ", userId)
		fmt.Println("spellId: ", setupId)
		web.SetErrorRedirect(w, r, "/home", "Error casting")
		return
	}
	spellSetup.RemainingSpells[level]--
	//db2.Set("spell-setup", setupId, spellSetup)
	db.Set("spell-setup", setupId, spellSetup)
	http.Redirect(w, r, "/home", 303)
	return

}}

var ppCast = web.Route{"POST", "/pp-cast", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id")
	setupId := r.FormValue("setupId")
	pp, err := strconv.Atoi(r.FormValue("pp"))
	var powerPointsSetup PowerPointsSetup
	//db2.Get("pp-setup", setupId).As(&ppSetup)
	db.Get("pp-setup", setupId, &powerPointsSetup)
	//fmt.Println())
	if err != nil {
		log.Printf("ppCast() >> strconv.Atoi(): %v\n", err)
		web.SetErrorRedirect(w, r, "/home", "Error casting")
		return
	}
	if powerPointsSetup.UserId != userId {
		fmt.Println("userId: ", userId)
		fmt.Println("setup userId: ", powerPointsSetup.UserId)
		web.SetErrorRedirect(w, r, "/home", "Error casting")
	}
	if pp > powerPointsSetup.RemainingPowerPoints {
		web.SetErrorRedirect(w, r, "/home", "You do not have enough power points left to cast that!")
		return
	}
	if pp < 0 {
		web.SetErrorRedirect(w, r, "/home", "You cannot cast negative power points")
		return
	}
	db.Get("pp-setup", setupId, &powerPointsSetup)
	powerPointsSetup.RemainingPowerPoints -= pp
	//db2.Set("pp-setup", setupId, ppSetup)
	db.Set("pp-setup", setupId, powerPointsSetup)
	http.Redirect(w, r, "/home", 303)
	return

}}

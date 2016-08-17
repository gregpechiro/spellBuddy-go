package main

import (
	"net/http"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var account = web.Route{"GET", "/account", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id").(string)
	var user User
	db.Get("user", userId, &user)
	tmpl.Render(w, r, "account.tmpl", web.Model{
		//"user": db2.Get("user", userId),
		"user": user,
	})
}}

var saveAccount = web.Route{"POST", "/account", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id").(string)

	// docs := db2.Query("user", dbdbMod.Eq{"Username", r.FormValue("username")}, dbdbMod.Ne{"::Id::", userId})
	// if len(docs) > 0 {
	// 	web.SetErrorRedirect(w, r, "/account", "Error adding user. Username already exists")
	// 	return
	// }

	var users []User
	db.TestQuery("user", &users, adb.Eq("username", r.FormValue("username")), adb.Ne("id", `"`+userId+`"`))
	if len(users) > 0 {
		web.SetErrorRedirect(w, r, "/account", "Error saving user. Username already exists")
		return
	}

	r.ParseForm()
	var user User
	//db2.Get("user", userId).As(&user)
	db.Get("user", userId, &user)
	FormToStruct(&user, r.Form, "")
	//db2.Set("user", userId, user)
	db.Set("user", userId, user)
	web.SetSuccessRedirect(w, r, "/account", "Successfully updated account")
	return
}}

var saveTheme = web.Route{"POST", "/theme", func(w http.ResponseWriter, r *http.Request) {
	userId := web.GetSess(r, "id").(string)
	var user User
	//db2.Get("user", userId).As(&user)
	db.Get("user", userId, &user)
	user.Theme = r.FormValue("theme")
	//db2.Set("user", userId, user)
	db.Set("user", userId, user)
	//web.SetSuccessRedirect(w, r, "/account", "Successfully updated theme")
	http.Redirect(w, r, "/account", 303)
	return
}}

package main

import (
	"net/http"

	"github.com/cagnosolutions/dbdb"
	"github.com/cagnosolutions/web"
)

var account = web.Route{"GET", "/account", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	tmpl.Render(w, r, "account.tmpl", web.Model{
		"user": db.Get("user", userId),
	})
}}

var saveAccount = web.Route{"POST", "/account", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	docs := db.Query("user", dbdb.Eq{"Username", r.FormValue("username")}, dbdb.Ne{"::Id::", userId})
	if len(docs) > 0 {
		web.SetErrorRedirect(w, r, "/account", "Error adding user. Username already exists")
		return
	}
	r.ParseForm()
	var user User
	db.Get("user", userId).As(&user)
	FormToStruct(&user, r.Form, "")
	db.Set("user", userId, user)
	web.SetSuccessRedirect(w, r, "/account", "Successfully updated account")
	return
}}

var saveTheme = web.Route{"POST", "/theme", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	var user User
	db.Get("user", userId).As(&user)
	user.Theme = r.FormValue("theme")
	db.Set("user", userId, user)
	web.SetSuccessRedirect(w, r, "/account", "Successfully updated theme")
	return
}}

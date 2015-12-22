package main

import "github.com/cagnosolutions/web"

var ADMIN = web.Auth{
	Roles:    []string{"ADMIN"},
	Redirect: "/login",
	Msg:      "You are not authorized",
}

var USER = web.Auth{
	Roles:    []string{"USER"},
	Redirect: "/login",
	Msg:      "You must be logged in",
}

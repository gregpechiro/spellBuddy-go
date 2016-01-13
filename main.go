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

	mux.AddRoutes(login, loginPost, logout)
	mux.AddSecureRoutes(ADMIN, adminHome, adminUser, addUser, saveUser, delUser, modifiySpells)
	mux.AddSecureRoutes(USER, home, setup, addSpell, saveSpell, editSpell)
	mux.AddSecureRoutes(USER, addSpellToUser, delSpellFromUser, changeLvl, spellsPerDay, rest, cast)
	mux.AddSecureRoutes(USER, pp, ppRest, ppCast)

	http.ListenAndServe(":8080", mux)
}

var login = web.Route{"GET", "/login", func(w http.ResponseWriter, r *http.Request) {
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
		web.SetSuccessRedirect(w, r, "/", "Welcome "+user.Username)
		return
	}
	web.SetErrorRedirect(w, r, "/login", "Incorrect username or password")
	return
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w, r)
	web.SetSuccessRedirect(w, r, "/login", "See you next time")
}}

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

var home = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
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
		web.SetErrorRedirect(w, r, "/", "Error Resting")
		return
	}
	copy(spellSetup.RemainingSpells, spellSetup.SpellsPerDay)
	copy(spellSetup.PreparedSpells, prepared)
	db.Set("spell-setup", setupId, spellSetup)
	http.Redirect(w, r, "/", 303)
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
		web.SetErrorRedirect(w, r, "/", "Error Resting")
		return
	}
	ppSetup.RemainingPowerPoints = ppSetup.TotalPowerPoints
	db.Set("pp-setup", setupId, ppSetup)
	http.Redirect(w, r, "/", 303)
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
		web.SetErrorRedirect(w, r, "/", "Error casting")
		return
	}
	spellSetup.RemainingSpells[level]--
	db.Set("spell-setup", setupId, spellSetup)
	http.Redirect(w, r, "/", 303)
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
		web.SetErrorRedirect(w, r, "/", "Error casting")
		return
	}
	ppSetup.RemainingPowerPoints -= pp
	db.Set("pp-setup", setupId, ppSetup)
	http.Redirect(w, r, "/", 303)
	return

}}

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

var addSpell = web.Route{"GET", "/add/spell", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "addSpell.tmpl", web.Model{})
}}

var editSpell = web.Route{"GET", "/edit/spell/:id", func(w http.ResponseWriter, r *http.Request) {
	userId := ParseId(web.GetSess(r, "id"))
	spellId := ParseId(r.FormValue(":id"))
	if spellId < 1 {
		web.SetErrorRedirect(w, r, "/setup?userc", "Invalid Spell")
		return
	}
	var spell Spell
	db.Get("spell", spellId).As(&spell)
	if spell.UserId != userId || !spell.Custom {
		web.SetErrorRedirect(w, r, "/setup?userc", "Cannot edit spell")
		return
	}
	tmpl.Render(w, r, "addSpell.tmpl", web.Model{
		"spell":   spell,
		"spellId": spellId,
	})
	return
}}

var saveSpell = web.Route{"POST", "/save/spell", func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	spellId := ParseId(r.FormValue("id"))
	userId := ParseId(web.GetSess(r, "id"))
	var spell Spell
	db.Get("spell", spellId).As(&spell)
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
	if spellId != 0 {
		if spell.UserId != userId {
			web.SetErrorRedirect(w, r, "/setup/?cat=userc", "Cannot save spell")
			fmt.Printf("Spell: %v\n", spell)
			fmt.Printf("UserId: %v\n", userId)
			return
		}
		db.Set("spell", spellId, spell)
	} else {
		spell.UserId = userId
		spell.Custom = true
		db.Add("spell", spell)
	}
	web.SetSuccessRedirect(w, r, "/setup?cat=userc", "Successfully saved spell")
	return
}}

func Ternary(comp bool, v, w interface{}) interface{} {
	if comp {
		return v
	}
	return w
}

func getPicked(userP [][]float64) [][]interface{} {
	var picked [][]interface{}
	if userP != nil {
		pickedLvl := []interface{}{}
		for _, lvl := range userP {
			pickedLvl = []interface{}{}
			if len(lvl) > 0 {
				for _, spellId := range lvl {
					pickedLvl = append(pickedLvl, db.Get("spell", spellId))
				}
			}
			picked = append(picked, pickedLvl)
		}
	}
	return picked
}

func getPickedNames(userP [][]float64) [][]string {
	var picked [][]string
	if userP != nil {
		pickedLvl := []string{}
		for _, lvl := range userP {
			pickedLvl = []string{}
			if len(lvl) > 0 {
				for _, spellId := range lvl {
					pickedLvl = append(pickedLvl, db.Get("spell", spellId).Data["Name"].(string))
				}
			}
			picked = append(picked, pickedLvl)
		}
	}
	return picked
}

func GetBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func ParseId(v interface{}) float64 {
	var id float64
	var err error
	switch v.(type) {
	case string:
		id, err = strconv.ParseFloat(v.(string), 64)
		if err != nil {
			log.Printf("ParseId() >> strconv.ParseFloat(): ", err)
		}
	case uint64:
		id = float64(v.(uint64))
	case float64:
		id = v.(float64)
	}
	return id
}

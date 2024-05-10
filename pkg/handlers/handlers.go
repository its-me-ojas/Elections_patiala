package handlers

import (
	"Elections_Patiala/pkg/db"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"text/template"
)

var store = sessions.NewCookieStore([]byte("2eb7ddef6411bca4205d73dbfcbf9115fcf2ec43"))

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	tmpl.Execute(w, map[string]interface{}{})
}

func HandlerAdminLogin(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Error getting session", err)
		return
	}

	if session.Values["authenticated"] == true {
		http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
		return
	}

	var errorMsg string
	if val, ok := session.Values["error"].(string); ok {
		errorMsg = val
		delete(session.Values, "error")
		session.Save(r, w)
	}

	tmpl := template.Must(template.ParseFiles("web/templates/login.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Error": errorMsg,
	})

}

func HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	usertype := r.FormValue("usertype")

	id, err := db.AuthenticateAdmin(username, password, usertype)
	if err != nil {
		log.Println("Authentication error: ", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}
	session, _ := store.Get(r, "session-name")
	if id == "" || id == "0" {
		session.Values["error"] = "Incorrect username or password"
		session.Save(r, w)
		http.Redirect(w, r, "/admin/login", http.StatusExpectationFailed)
		return
	} else {
		session.Values["authenticated"] = true
		session.Values["usertype"] = usertype
		if usertype == "aro" {
			session.Values["cid"] = id
		} else {
			session.Values["bid"] = id
		}
		session.Save(r, w)
		http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
	}
}

func HanldeAdminDashboard(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["authenticated"] != true {
		http.Redirect(w, r, "/admin/login", http.StatusUnauthorized)
		return
	}

	if session.Values["usertype"] == "blo" {
		bid, ok := session.Values["bid"].(string)
		if !ok {
			http.Error(w, "Invalid Session Data", http.StatusInternalServerError)
		}

	}

}

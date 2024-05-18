package handlers

import (
	"Elections_Patiala/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("2eb7ddef6411bca4205d73dbfcbf9115fcf2ec43"))

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Error getting session", err)
		return
	}
	if session.Values["authenticated"] == true {
		http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/admin/login", http.StatusFound)

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
	contact := r.FormValue("contact")
	password := r.FormValue("password")

	userInfo, err := db.AuthenticateAdmin(contact, password)
	if err != nil {
		log.Println("Authentication error: ", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	if userInfo == nil {
		session, _ := store.Get(r, "session-name")
		session.Values["error"] = "Incorrect username or password"
		session.Save(r, w)
		http.Redirect(w, r, "/admin/login", http.StatusExpectationFailed)
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["usertype"] = userInfo["usertype"]

	if userInfo["usertype"] == "aro" {
		session.Values["cid"] = userInfo["cid"]
	} else {
		session.Values["cid"] = userInfo["cid"]
		session.Values["bid"] = userInfo["bid"]
	}

	session.Save(r, w)
	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
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
		liveTraffic, err := db.GetLiveTrafficByBoothID(bid)
		if err != nil {
			log.Printf("Failed to retrieve polling station data: %v", err)
			http.Error(w, "Failed to retrieve polling station data", http.StatusInternalServerError)
			return
		}
		tmpl := template.Must(template.ParseFiles("web/templates/adminBLO.html"))
		err = tmpl.Execute(w, liveTraffic)
		if err != nil {
			log.Println("Failed to render the template")
		}
	} // else {
	//	cid, ok := session.Values["cid"].(string)
	//	if !ok {
	//		http.Error(w, "Invalid Session Data", http.StatusInternalServerError)
	//	}
	//	liveTraffic, err := db.GetLiveTrafficByConstituencyID(cid)
	//	if err != nil {
	//		log.Printf("Failed to retrieve polling station data: %v", err)
	//		http.Error(w, "Failed to retrieve polling station data", http.StatusInternalServerError)
	//		return
	//	}
	//	tmpl := template.Must(template.ParseFiles("web/templates/adminARO.html"))
	//	err = tmpl.Execute(w, liveTraffic)
	//	if err != nil {
	//		log.Println("Failed to render the template")
	//	}
	//}

}

func HanldeLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Error getting session", err)
		return
	}
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

func CounterUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get the counter value from the form
	peopleInQueueStr := r.FormValue("counter")

	// Convert string to integer
	peopleInQueue, err := strconv.Atoi(peopleInQueueStr)
	if err != nil {
		log.Println("Error converting string to integer: ", err)
		http.Error(w, "Invalid counter value", http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["authenticated"] != true {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bid, ok := session.Values["bid"].(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	db.UpdateQueue(bid, peopleInQueue)
	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)

}

//func VoterReqUpdate(w http.ResponseWriter, r *http.Request) {
//	// Parse the form data
//	err := r.ParseForm()
//	if err != nil {
//		http.Error(w, "Invalid form data", http.StatusBadRequest)
//		return
//	}
//
//	// Get the counter value from the form
//	voterReqStr := r.FormValue("voterReq")
//
//	// Convert string to integer
//	voterReq, err := strconv.Atoi(voterReqStr)
//	if err != nil {
//		log.Println("Error converting string to integer: ", err)
//		http.Error(w, "Invalid counter value", http.StatusBadRequest)
//		return
//	}
//
//	session, err := store.Get(r, "session-name")
//	if err != nil || session.Values["authenticated"] != true {
//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
//		return
//	}
//
//	bid, ok := session.Values["bid"].(string)
//	if !ok {
//		http.Error(w, "Invalid session data", http.StatusInternalServerError)
//		return
//	}
//
//	db.UpdateVoterReq(bid, voterReq)
//	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
//
//}

func HandleGetAllVoters(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["authenticated"] != true {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	cid, ok := session.Values["cid"].(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}
	voters, err := db.GetAllVoters(cid)
	if err != nil {
		http.Error(w, "Failed to get voters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(voters)
}

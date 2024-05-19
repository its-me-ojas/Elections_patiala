package handlers

import (
	"Elections_Patiala/pkg/db"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"text/template"
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
	fmt.Println("Handling Auth")
	contact := r.FormValue("contact")
	password := r.FormValue("password")

	userInfo, err := db.AuthenticateAdmin(contact, password)
	if err != nil {
		log.Println("Authentication error: ", err)
		userInfo = nil
	}

	if userInfo == nil {
		session, _ := store.Get(r, "session-name")
		session.Values["error"] = "Incorrect username or password"
		session.Save(r, w)
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["error"] = ""
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

func HandleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Dashboard Opening")

	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["authenticated"] != true {
		http.Redirect(w, r, "/admin/login", http.StatusUnauthorized)
		return
	}
	if session.Values["usertype"] == "blo" {
		cid, ok := session.Values["cid"].(string)
		bid, ok := session.Values["bid"].(string)
		if !ok {
			http.Error(w, "Invalid Session Data", http.StatusInternalServerError)
		}

		booth, err := db.GetBooth(cid, bid)
		display_data, err := db.GetDisplayData(cid, bid)
		voter_req_data, err := db.GetAllVoters(cid)
		if err != nil {
			log.Printf("Failed to retrieve polling station data: %v", err)
			http.Error(w, "Failed to retrieve polling station data", http.StatusInternalServerError)
			return
		}
		tmpl := template.Must(template.ParseFiles("web/templates/adminBLO.html"))
		err = tmpl.Execute(w, map[string]interface{}{
			"Booth":        booth,
			"DisplayData":  display_data,
			"VoterReqData": voter_req_data,
		})
		if err != nil {
			log.Println("Failed to render the template")
		}

	} else if session.Values["usertype"] == "ps" {
		cid, ok := session.Values["cid"].(string)
		bid, ok := session.Values["bid"].(string)
		if !ok {
			http.Error(w, "Invalid Session Data", http.StatusInternalServerError)
		}

		booth, err := db.GetBooth(cid, bid)
		display_data, err := db.GetDisplayData(cid, bid)
		if err != nil {
			log.Printf("Failed to retrieve polling station data: %v", err)
			http.Error(w, "Failed to retrieve polling station data", http.StatusInternalServerError)
			return
		}
		tmpl := template.Must(template.ParseFiles("web/templates/adminPS.html"))
		err = tmpl.Execute(w, map[string]interface{}{
			"Booth":       booth,
			"DisplayData": display_data,
		})
		if err != nil {
			log.Println("Failed to render the template")
		}

	} else {
		cid, ok := session.Values["cid"].(string)
		if !ok {
			http.Error(w, "Invalid Session Data", http.StatusInternalServerError)
		}
		voterReqData, err := db.GetAllVoters(cid)
		if err != nil {
			log.Printf("Failed to retrieve polling station data: %v", err)
			http.Error(w, "Failed to retrieve polling station data", http.StatusInternalServerError)
			return
		}
		tmpl := template.Must(template.ParseFiles("web/templates/adminARO.html"))
		err = tmpl.Execute(w, map[string]interface{}{
			"VoterReqData": voterReqData,
		})
		if err != nil {
			log.Println("Failed to render the template")
		}
	}

}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Error getting session", err)
		return
	}
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

func HandleCounterUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get the counter value from the form
	peopleInQueue := r.FormValue("counter")

	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["authenticated"] != true {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	cid, ok := session.Values["cid"].(string)
	bid, ok := session.Values["bid"].(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	db.UpdateQueue(cid, bid, peopleInQueue)
	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)

}

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
	// bid, ok := session.Values["bid"].(string)
	// if !ok {
	// 	http.Error(w, "Invalid session data", http.StatusInternalServerError)
	// 	return
	// }
	voters, err := db.GetAllVoters(cid)
	if err != nil {
		http.Error(w, "Failed to get voters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(voters)
}

func HandleGetQueue(w http.ResponseWriter, r *http.Request) {
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
	cid, ok := session.Values["cid"].(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}
	queue, err := db.GetQueue(cid, bid)
	if err != nil {
		http.Error(w, "Failed to get queue", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queue)
}

func HandleGetBoothData(w http.ResponseWriter, r *http.Request) {
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
	bid, ok := session.Values["bid"].(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}
	booth, err := db.GetBooth(cid, bid)
	if err != nil {
		http.Error(w, "Failed to get booth data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booth)
}

func HandleVoterReqStatus(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	objectID := r.FormValue("objectID")

	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["authenticated"] != true {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println("First")
	err = db.UpdateVoterRequest(objectID)
	if err != nil {
		http.Error(w, "Failed to get voter request", http.StatusInternalServerError)
		return
	}
	fmt.Println("Second")

	http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
}

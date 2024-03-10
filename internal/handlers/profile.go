package handlers

import (
	"errors"
	"ipwhitelister/internal/database"
	"ipwhitelister/internal/web"
	"log"
	"net/http"
	"strings"
)

type ProfilePageData struct {
	Email            string
	CurrentIP        string
	LastAssociatedIP string
	PendingIP        string
	IsAuthenticated  bool
	CommonPageData
}

// profileHandler handles the HTTP request for the profile page - renders page with profile and IPs data.
func (app *App) profileHandler(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value(emailKey).(string)
	if !ok {
		log.Println("Error getting email from context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := web.Template("layout.html", "profile.html")
	if err != nil {
		log.Println("Error parsing profile template: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	profile, err := database.GetProfileByEmail(app.DB, email)
	if err != nil {
		log.Println("Error getting profile: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	currentIP, err := app.getCurrentIP(r)
	if err != nil {
		log.Println("Error getting current IP")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := ProfilePageData{
		Email:            profile.Email,
		CurrentIP:        currentIP,
		LastAssociatedIP: profile.AssociatedIP,
		PendingIP:        profile.PendingIP,
		IsAuthenticated:  true,
		CommonPageData:   app.getCommonPageData(r),
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing profile template: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// getCurrentIP retrieves the client's IP address from the HTTP request.
// If the `RealIPHeader` configuration option is set, it uses the value of that header.
// Otherwise, it falls back to using the `RemoteAddr` field from the request.
// The IP address is extracted by splitting the value at the first colon (if present).
// Returns the IP address as a string and an error if the IP address could not be obtained.
func (app *App) getCurrentIP(r *http.Request) (string, error) {
	var ip string
	if app.CFG.RealIPHeader == "" {
		ip = r.RemoteAddr
	} else {
		ip = r.Header.Get(app.CFG.RealIPHeader)
		ip = strings.Split(ip, ":")[0]
	}

	if ip == "" {
		return "", errors.New("could not get IP address")
	}

	return ip, nil
}

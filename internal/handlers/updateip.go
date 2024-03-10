package handlers

import (
	"ipwhitelister/internal/cloudflare"
	"ipwhitelister/internal/database"
	"log"
	"net/http"
)

// updateIPHandler handles the HTTP request for updating the IP address.
// It expects a "new_ip" value in the request form and the user's email
// address in the request context. It updates the IP address in the local
// database and synchronizes it with Cloudflare. If the update is successful,
// it triggers a page reload.
func (app *App) updateIPHandler(w http.ResponseWriter, r *http.Request) {
	newIP := r.FormValue("new_ip")
	email := r.Context().Value(emailKey).(string)

	profile, err := database.GetProfileByEmail(app.DB, email)
	if err != nil {
		log.Println("Error getting profile: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// first step: verify user has no pending IP, if they do - reassing that in local DB so last
	// associated IP gets syncronised with cloudflare and then updates the pending IP to be ""
	// Fetch the specific rule from Cloudflare
	rule, err := app.CFM.FetchLockdownRule(app.CFG.CloudflareRuleID)
	if err != nil {
		log.Println("Error fetching lockdown rule: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fix previous pending IP if was not synced. Or set newIP as pending.
	if profile.PendingIP != "" {
		for _, config := range rule.Configurations {
			if config.Target == "ip" && config.Value == profile.PendingIP {
				// Confirm the IP update in the local database
				// override last associated IP with last pending IP
				// pending will become the new IP anyway
				profile.AssociatedIP = profile.PendingIP
				break
			}
		}
	}

	if profile.PendingIP != newIP {
		profile.PendingIP = newIP
	} else {
		profile.PendingIP = ""
	}

	err = database.UpdateProfile(app.DB, profile)
	if err != nil {
		log.Println("Error updating profile: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if profile.PendingIP == "" {
		w.Header().Set("HX-Redirect", "/profile")
		return
	}

	// lookup for last associated which is 100% synced if it was in Cloudflare's records
	// if IP is not found in Cloudflare's records, add it to the rule and update it
	isFound := false
	for _, config := range rule.Configurations {
		if config.Target == "ip" && config.Value == profile.AssociatedIP {
			config.Value = profile.PendingIP
			isFound = true
			break
		}
	}

	if !isFound {
		rule.Configurations = append(rule.Configurations, &cloudflare.Configuration{
			Target: "ip",
			Value:  profile.AssociatedIP,
		})
	}

	err = app.CFM.UpdateLockdownRule(rule)
	if err != nil {
		log.Println("Error updating lockdown rule: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Cofirm the IP update in the local database
	profile.AssociatedIP = profile.PendingIP
	profile.PendingIP = ""
	err = database.UpdateProfile(app.DB, profile)
	if err != nil {
		log.Println("Error updating profile: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// HTMX trigger entire page reload
	w.Header().Set("HX-Redirect", "/profile")
}

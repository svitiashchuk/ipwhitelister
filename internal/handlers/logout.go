package handlers

import (
	"log"
	"net/http"
)

// logoutHandler handles the logout functionality for the application.
// It invalidates the session token, clears the session cookie in the user's browser,
// and redirects the user to the login page or home page.
func (app *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Error getting session token from cookie: ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Invalidate the session token
	app.SM.DeleteSessionData(cookie.Value)
	// Proceed to clear the session cookie in the user's browser
	app.clearSessionCookie(w)

	http.Redirect(w, r, "/logged-out", http.StatusSeeOther)
}

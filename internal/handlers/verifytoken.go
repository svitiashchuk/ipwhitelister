package handlers

import (
	"ipwhitelister/internal/database"
	"ipwhitelister/internal/session"
	"net/http"
)

// verifyTokenHandler verifies the validity of a token and handles the corresponding logic.
// It checks if the token is valid by calling the CheckToken function from the database package.
// If the token is valid, it generates a session token using the GenerateSessionToken function from the session package,
// stores the session data using the StoreSessionData method of the App's SM field,
// sets a session cookie with the generated session token,
// and sets the "HX-Redirect" header to "/profile" to redirect the client to the profile page.
// If the token is invalid, it returns a "Invalid token" error with a status code of 400 (Bad Request).
func (app *App) verifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenIsValid, _ := database.CheckToken(app.DB, r.FormValue("email"), r.FormValue("token"))

	if tokenIsValid {
		sessionToken, err := session.GenerateSessionToken()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		app.SM.StoreSessionData(sessionToken, r.FormValue("email"))

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			HttpOnly: true,
		})

		w.Header().Set("HX-Redirect", "/profile")
		return
	} else {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}
}

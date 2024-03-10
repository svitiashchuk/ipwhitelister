package handlers

import (
	"context"
	"database/sql"
	"ipwhitelister/internal/cloudflare"
	"ipwhitelister/internal/config"
	"ipwhitelister/internal/email"
	"ipwhitelister/internal/session"
	"ipwhitelister/internal/web"
	"net/http"
)

// App holds the application's shared resources
type App struct {
	CFG    *config.Config
	DB     *sql.DB
	SM     *session.SessionManager
	CFM    *cloudflare.CloudflareManager
	Mailer *email.Mailer
}

// CommonPageData holds common data that is passed to all page templates.
type CommonPageData struct {
	IsAuthenticated bool
	CSRFToken       string
}

// SetupServer configures the HTTP server by setting up the routes and handlers.
func (app *App) SetupServer() {
	mux := http.NewServeMux()

	mux.Handle("/static/", web.StaticHandler())
	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/login", app.loginHandler)
	mux.HandleFunc("/verify-token", app.verifyTokenHandler)
	mux.HandleFunc("/profile", app.validateSessionMiddleware(app.profileHandler))
	mux.HandleFunc("/update-ip", app.validateSessionMiddleware(app.updateIPHandler))
	mux.HandleFunc("/logout", app.logoutHandler)
	mux.HandleFunc("/logged-out", app.goodbyeHandler)

	http.Handle("/", mux)
}

// validateSession checks the request for a session token cookie and validates it.
// If the session is valid, the user's email is returned along with true.
// If the session is invalid, an empty string and false are returned.
func (app *App) validateSession(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", false
	}

	email, exists := app.SM.RetrieveSessionData(cookie.Value)

	return email, exists
}

// getCommonPageData returns a CommonPageData struct with the IsAuthenticated field set to true if the user is logged in.
func (app *App) getCommonPageData(r *http.Request) CommonPageData {
	_, loggedIn := app.validateSession(r)
	return CommonPageData{
		IsAuthenticated: loggedIn,
	}
}

type contextKey string

const emailKey contextKey = "email"

func (app *App) validateSessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, loggedIn := app.validateSession(r)
		if !loggedIn {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), emailKey, email)

		// User is authenticated; proceed with the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// clearSessionCookie deletes the session cookie by setting it to expire immediately.
func (app *App) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

package handlers

import (
	"fmt"
	"ipwhitelister/internal/auth"
	"ipwhitelister/internal/database"
	"ipwhitelister/internal/web"
	"log"
	"net/http"
)

type LoginPageData struct {
	Email string
	CommonPageData
}

// loginHandler handles the login requests.
// For GET requests, it renders the login page.
// For POST requests, it processes the login request, generates a login token,
// stores it in the database, and sends a verification email to the user.
// If any error occurs during the process, it returns an Internal Server Error response.
func (app *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, err := web.Template("layout.html", "login.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, app.getCommonPageData(r))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	case "POST":
		// Process the login request
		r.ParseForm()
		email := r.FormValue("email")

		isExists, err := database.EmailExists(app.DB, email)
		if err != nil {
			log.Println("Error checking if email exists: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !isExists {
			fmt.Fprintf(w, "Email %s does not exist in the system", email)
			return
		}

		token, err := auth.GenerateToken()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		database.StoreLoginToken(app.DB, email, token, auth.TokenValidityDuration)
		err = app.Mailer.SendTokenEmail(email, token)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tmpl, err := web.Template("fragments/verify_token_form.html")
		if err != nil {
			log.Println("Error checking if email exists: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, LoginPageData{Email: email, CommonPageData: app.getCommonPageData(r)})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

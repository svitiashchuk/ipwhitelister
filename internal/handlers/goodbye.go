package handlers

import (
	"ipwhitelister/internal/web"
	"log"
	"net/http"
)

func (app *App) goodbyeHandler(w http.ResponseWriter, r *http.Request) {
	_, loggedIn := app.validateSession(r)
	if loggedIn {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	tmpl, err := web.Template("layout.html", "logged_out.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, app.getCommonPageData(r))
	if err != nil {
		log.Println("Error executing logged-out template: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

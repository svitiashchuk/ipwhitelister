package handlers

import (
	"ipwhitelister/internal/web"
	"log"
	"net/http"
)

func (app *App) homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := web.Template("layout.html", "home.html")
	if err != nil {
		log.Println("Error parsing home template: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, app.getCommonPageData(r))
	if err != nil {
		log.Println("Error executing home template: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

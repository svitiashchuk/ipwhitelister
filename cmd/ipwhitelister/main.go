package main

import (
	"fmt"
	"ipwhitelister/internal/cloudflare"
	"ipwhitelister/internal/config"
	"ipwhitelister/internal/database"
	"ipwhitelister/internal/email"
	"ipwhitelister/internal/handlers"
	"ipwhitelister/internal/session"
	"log"
	"net/http"
)

func main() {
	app := &handlers.App{
		CFG: config.New(),
		SM:  session.NewSessionManager(),
	}

	app.DB = database.InitDB(app.CFG.DBPath)
	if err := app.DB.Ping(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	app.Mailer = email.NewMailer(app.CFG.MandrillAPIKey, app.CFG.SenderEmail)

	app.CFM = cloudflare.NewCloudflareManager(
		app.CFG.CloudflareAPIToken,
		app.CFG.CloudflareZoneID,
		app.DB,
	)

	app.SetupServer()
	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

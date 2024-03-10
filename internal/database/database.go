// Database operations
package database

import (
	"database/sql"
	"ipwhitelister/internal/model"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the database and creates the profiles table if it does not exist.
func InitDB(filePath string) *sql.DB {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS profiles (
        "email" TEXT NOT NULL UNIQUE,
        "associated_ip" TEXT NOT NULL,
        "pending_ip" TEXT NOT NULL,
        PRIMARY KEY("email")
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	createLoginTokensTableSQL := `CREATE TABLE IF NOT EXISTS login_tokens (
        "email" TEXT NOT NULL UNIQUE,
        "token" TEXT NOT NULL,
        "expiration" DATETIME NOT NULL,
        PRIMARY KEY("email")
    );`
	_, err = db.Exec(createLoginTokensTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// CreateProfile inserts a new profile into the database.
func CreateProfile(db *sql.DB, profile model.Profile) error {
	query := `INSERT INTO profiles(email, associated_ip, pending_ip) VALUES (?, ?, ?)`
	_, err := db.Exec(query, profile.Email, "", "")
	return err
}

// GetProfileByEmail retrieves a profile from the database based on the provided email.
func GetProfileByEmail(db *sql.DB, email string) (model.Profile, error) {
	query := `SELECT email, associated_ip, pending_ip FROM profiles WHERE email = ?`
	var profile model.Profile
	err := db.QueryRow(query, email).Scan(&profile.Email, &profile.AssociatedIP, &profile.PendingIP)
	if err != nil {
		return model.Profile{}, err
	}
	return profile, nil
}

// UpdateProfile updates the profile in the database with the provided associated IP and pending IP for the given email.
func UpdateProfile(db *sql.DB, profile model.Profile) error {
	query := `UPDATE profiles SET associated_ip = ?, pending_ip = ? WHERE email = ?`
	_, err := db.Exec(query, profile.AssociatedIP, profile.PendingIP, profile.Email)
	return err
}

// DeleteProfile deletes a profile from the database based on the provided email.
func DeleteProfile(db *sql.DB, email string) error {
	query := `DELETE FROM profiles WHERE email = ?`
	_, err := db.Exec(query, email)
	return err
}

// EmailExists checks if a given email exists in the database.
func EmailExists(db *sql.DB, email string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT exists (SELECT 1 FROM profiles WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, err
}

// StoreLoginToken saves the login token and its expiration time in the database.
func StoreLoginToken(db *sql.DB, email, token string, validityDur time.Duration) error {
	expirationTime := time.Now().Add(validityDur)
	_, err := db.Exec(`INSERT INTO login_tokens (email, token, expiration) VALUES (?, ?, ?)
                       ON CONFLICT(email) DO UPDATE SET token = excluded.token, expiration = excluded.expiration`,
		email, token, expirationTime)
	return err
}

// CheckToken checks if a given token is valid for the email.
func CheckToken(db *sql.DB, email, token string) (bool, error) {
	var expiration time.Time
	err := db.QueryRow("SELECT expiration FROM login_tokens WHERE email = ? AND token = ?", email, token).Scan(&expiration)
	if err != nil {
		return false, err
	}
	return time.Now().Before(expiration), nil
}

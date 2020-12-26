package db

import (
	"database/sql"
	"errors"
	"os"

	// the linter complains about the blank identifier in the import block if I don't comment here
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// Create creates a new sqlite database for credentials with name at path.
func Create(path string) (*sql.DB, error) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return nil, err
	}
	if err == nil {
		return nil, errors.New("cannot create db, file already exists")
	}
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	// Create a simple table to store our credential hashes
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS creds (id INTEGER PRIMARY KEY, hash TEXT type UNIQUE, description TEXT)")
	statement.Exec()
	return database, nil
}

// AddCred takes a credential in plaintext, hashes it, and then stores it in the database provided.
func AddCred(db *sql.DB, cred, description string) error {
	credHash, err := HashCred(cred)
	if err != nil {
		return err
	}
	statement, err := db.Prepare("INSERT INTO creds (hash, description) VALUES (?, ?)")
	_, err = statement.Exec(credHash, description)
	return err
}

// HashCred hashes a credential for safe storage in the database
func HashCred(cred string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(cred), 14)
	return string(bytes), err
}

// CheckCredHash checks if the credential provided matches the hash provided
func CheckCredHash(cred, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(cred))
	return err == nil
}

// CredExists returns true if a cred exists in the database.
func CredExists(db *sql.DB, credHash string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT exists (SELECT hash FROM creds WHERE hash = ?)", credHash).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

// GetCredDescription returns a description of a given credential hash
func GetCredDescription(db *sql.DB, credHash string) (string, error) {
	rows, err := db.Query("SELECT description FROM creds WHERE hash = ?", credHash)
	if err != nil {
		return "", err
	}
	var description string
	for rows.Next() {
		rows.Scan(&description)
	}
	return description, nil
}

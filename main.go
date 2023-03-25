package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Define a struct to hold the form data
type Shift struct {
	Name      string `db:"name"`
	ShiftDate string `db:"shift_date"`
	ShiftType string `db:"shift_type"`
	Task      string `db:"task"`
	TaskType  string `db:"task_type"`
	Hours     int    `db:"hours"`
	Minutes   int    `db:"minutes"`
}

type SearchResult struct {
	ShiftTasks []Shift
	TotalHours float64
}

func main() {

	if len(os.Args) <= 1 {
		log.Fatal("Usage: ./shiftform <dbPath>")
	}

	dbPath := os.Args[1]
	// Open the database
	db, err := sqlx.Open("sqlite3", filepath.Join(dbPath, "shifts.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createDDL(db)

	// Set up the HTTP server
	http.HandleFunc("/", handleAddTask(db))

	http.HandleFunc("/search", searchHandler())

	http.HandleFunc("/searchresults", searchResults(db))
	http.HandleFunc("/backuptable", backupTable(db))
	http.HandleFunc("/migratetable", migrateTable(db))
	http.HandleFunc("/backupdatabase", backupDatabase(dbPath))
	http.HandleFunc("/downloadtasks", downloadTasks(db))
	http.HandleFunc("/landing", landingPageHandler())

	// Start the server
	log.Fatal(http.ListenAndServe(":8085", nil))
}

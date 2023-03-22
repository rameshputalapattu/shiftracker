package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Define a struct to hold the form data
type Shift struct {
	Name      string `db:"name"`
	ShiftDate string `db:"shift_date"`
	ShiftType string `db:"shift_type"`
	Task      string `db:"task"`
	Hours     int    `db:"hours"`
}

func main() {
	// Open the database
	db, err := sqlx.Open("sqlite3", "shifts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createDDL(db)

	// Set up the HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Parse the form data
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Create a new Shift struct from the form data
			shift := Shift{
				Name:      r.FormValue("name"),
				ShiftDate: r.FormValue("shift_date"),
				ShiftType: r.FormValue("shift_type"),
				Task:      r.FormValue("task"),
				Hours: func() int {
					hours, _ := strconv.Atoi(r.FormValue("hours"))
					return hours
				}(),
			}

			// Insert the shift data into the database
			_, err = db.NamedExec(`INSERT INTO shifts (name, shift_date, shift_type, task, hours) VALUES (:name, :shift_date, :shift_type, :task, :hours)`, shift)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Redirect the user back to the form page
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Render the form page
		t, err := template.New("form").Parse(formHTML)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {

		// Render the form page
		t, err := template.New("search_form").Parse(searchForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/searchresults", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {

			// Parse the form data
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			name := r.FormValue("name")
			shiftDate := r.FormValue("shift_date")

			shiftTasks := []Shift{}

			err = db.Select(&shiftTasks,
				"select name,shift_type,shift_date,task,hours from shifts where name like $1 and shift_date = $2",
				"%"+name+"%",
				shiftDate,
			)

			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			t := template.Must(template.New("shiftTable").Parse(searchresults))

			err = t.Execute(w, shiftTasks)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return

		}

		http.Redirect(w, r, "/search", http.StatusSeeOther)

	})

	// Start the server
	log.Fatal(http.ListenAndServe(":8085", nil))
}

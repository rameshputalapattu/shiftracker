package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
)

func handleAddTask(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func searchHandler() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

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

	}

}

func searchResults(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
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

	}
}

func backupTable(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		db.MustExec(`
		drop table if exists shifts_bkp;
		CREATE table shifts_bkp as select * from shifts;
		
		
		`)

		fmt.Fprintln(w, "backed up the table successfully")

	}

}

func migrateTable(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db.MustExec(`
		drop table if exists shifts;
		CREATE TABLE IF NOT EXISTS shifts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			shift_date TEXT,
			shift_type TEXT,
			task TEXT,
			task_type TEXT,

			hours INTEGER default 0,
			minutes INTEGER default 0,
			created_timestamp TIMESTAMP default CURRENT_TIMESTAMP
		);
		insert into shifts(name,shift_date,shift_type,task,hours)
		select name,shift_date,shift_type,task,hours from shifts_bkp;
`)

		fmt.Fprintln(w, "migration completed successfully")

	}
}

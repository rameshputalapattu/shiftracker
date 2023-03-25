package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"embed"

	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
)

//go:embed shiftform.html
var shiftformPage embed.FS

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
				TaskType:  r.FormValue("task_type"),

				Hours: func() int {
					hours, _ := strconv.Atoi(r.FormValue("hours"))
					return hours
				}(),
				Minutes: func() int {
					minutes, _ := strconv.Atoi(r.FormValue("minutes"))
					return minutes

				}(),
			}

			// Insert the shift data into the database
			_, err = db.NamedExec(`INSERT INTO shifts (name, shift_date, shift_type, task_type,task, hours,minutes) 
			VALUES (:name, :shift_date, :shift_type,:task_type, :task, :hours,:minutes)`, shift)
			if err != nil {

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Redirect the user back to the form page
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Render the form page
		t, err := template.ParseFS(shiftformPage, "shiftform.html")
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

//go:embed searchform.html
var searchformPage embed.FS

func searchHandler() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// Render the form page

		t, err := template.ParseFS(searchformPage, "searchform.html")
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

//go:embed searchresults.html
var searchResultsPage embed.FS

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
				"select name,shift_type,shift_date,task,task_type,hours,minutes from shifts where name like $1 and shift_date = $2",
				"%"+name+"%",
				shiftDate,
			)

			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			th, err := totalHours(db, name, shiftDate)

			if err != nil {

				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			searchResults := SearchResult{
				ShiftTasks: shiftTasks,
				TotalHours: th,
			}

			t := template.Must(template.ParseFS(searchResultsPage, "searchresults.html"))

			err = t.Execute(w, searchResults)
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

func backupDatabase(dbPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tgt, err := os.Create(filepath.Join(dbPath, "shifts_bkp.db"))

		if err != nil {
			log.Println(err)
			return
		}

		defer tgt.Close()

		src, err := os.Open(filepath.Join(dbPath, "shifts.db"))
		if err != nil {
			log.Println(err)
			return
		}
		defer src.Close()
		_, err = io.Copy(tgt, src)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(w, "backup successful")

	}

}

func totalHours(db *sqlx.DB, name string, shiftDate string) (float64, error) {
	var totalHours float64
	query := `
	SELECT ROUND(sum(hours + (minutes/60.0)), 2) AS total_hours FROM shifts
	where name like $1 and shift_date = $2
	`

	err := db.Get(&totalHours,
		query,
		"%"+name+"%",
		shiftDate,
	)

	return totalHours, err

}

func downloadTasks(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// Retrieve all the rows from the "shifts" table
		shifts := []Shift{}
		err := db.SelectContext(context.Background(), &shifts, "SELECT name, shift_date, shift_type, task_type, task, hours, minutes FROM shifts")
		if err != nil {

			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Create a new Excel file
		file := excelize.NewFile()

		// Add a new sheet
		sheetName := "Sheet1"
		_, err = file.NewSheet(sheetName)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set the headers for the sheet

		sytleBoldID, _ := file.NewStyle(&excelize.Style{Font: &excelize.Font{
			Bold: true,
		}})
		headers := []string{"Name", "Shift Date", "Shift Type", "Task Type", "Task", "Hours", "Minutes"}
		for i, header := range headers {

			cellName, _ := excelize.CoordinatesToCellName(i+1, 1)

			file.SetCellValue(sheetName, cellName, header)
			file.SetCellStyle(sheetName, cellName, cellName, sytleBoldID)
		}

		// Add the rows to the sheet
		for i, shift := range shifts {
			row := []interface{}{
				shift.Name,
				shift.ShiftDate,
				shift.ShiftType,
				shift.TaskType,
				shift.Task,
				shift.Hours,
				shift.Minutes,
			}
			rowIndex := i + 2
			for j, cellValue := range row {

				cellName, _ := excelize.CoordinatesToCellName(j+1, rowIndex)
				file.SetCellValue(sheetName, cellName, cellValue)
			}
		}

		// Set the content type header for the response
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		// Set the content disposition header to force a download
		w.Header().Set("Content-Disposition", "attachment; filename=\"shifts.xlsx\"")

		// Write the Excel file to the response writer
		err = file.Write(w)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

}

//go:embed landing.html
var landingPage embed.FS

func landingPageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the contents of the embedded HTML file
		htmlContent, err := landingPage.ReadFile("landing.html")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		// Set the response header
		w.Header().Set("Content-Type", "text/html")

		// Write the HTML content to the response
		_, err = w.Write(htmlContent)
		if err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}
	}
}

package main

import "github.com/jmoiron/sqlx"

func createDDL(db *sqlx.DB) {

	// Create the shifts table if it doesn't exist
	db.MustExec(`CREATE TABLE IF NOT EXISTS shifts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        shift_date TEXT,
        shift_type TEXT,
        task TEXT,
        hours INTEGER,
		created_timestamp TIMESTAMP default CURRENT_TIMESTAMP
    )`)

}

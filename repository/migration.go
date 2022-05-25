package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

type Mahasiswa struct {
	ID       int    `db:"id"`
	NIM      string `db:"nim"`
	Name     string `db:"name"`
	Gender   string `db:"gender"`
	Religion string `db:"religion"`
}

// Migrate digunakan untuk melakukan migrasi database dengan data yang dibutuhkan
func Migrate() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./mini_project.db")
	if err != nil {
		panic(err)
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		email VARCHAR (100) NOT NULL
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
			INSERT INTO 
			    users (name, username, password, email) 
			VALUES 
			    ("Dion Maulana", "Noiid", "12345", "dion@gmail.com"),
				("Akbar Hasan", "Akbar", "098765", "akbar30@gmail.com");`)

	if err != nil {
		return nil, err
	}

	sqlStmt = `CREATE TABLE IF NOT EXISTS mahasiswas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nim VARCHAR (30) NOT NULL,
		name TEXT NOT NULL,
		gender VARCHAR(20) NOT NULL,
		religion TEXT NOT NULL
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
			INSERT INTO 
			    mahasiswas (nim, name, gender, religion) 
			VALUES 
			    ("123", "Iqbal Kabul", "Laki-laki", "Islam"),
				("345", "Albertus Wibisono", "Laki-laki","Kristen"),
				("678", "Ni Wayan Suyani", "Perempuan","Buddha"),
				("990", "Nur Aini", "Perempuan", "Islam");`)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Rollback(db *sql.DB) {
	sqlStmt := `DROP TABLE users;`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}

	sqlStmt = `DROP TABLE mahasiswas;`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}
}

// Jalankan main untuk melakukan migrasi database
func main() {
	db, err := Migrate()
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", user)
	}

	rows, err = db.Query("SELECT * FROM mahasiswas")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var mhs Mahasiswa
		err = rows.Scan(&mhs.ID, &mhs.NIM, &mhs.Name, &mhs.Gender, &mhs.Religion)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", mhs)
	}

	// Use Rollback() to rollback the changes
	//Rollback(db)
}

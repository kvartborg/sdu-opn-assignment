package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Person struct {
	PersonID            int
	Firstname, Lastname string
}

func env(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func setup(db *sql.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS persons; CREATE TABLE persons (
			PersonID serial NOT NULL,
			Firstname varchar(255) NOT NULL,
			Lastname varchar(255) NOT NULL,
			PRIMARY KEY (PersonID)
	)`)

	fmt.Println(err)

	if err != nil {
		time.Sleep(1 * time.Second)
		setup(db)
	}
}

func setupCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func main() {
	db, _ := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=root password=root dbname=root sslmode=disable", env("DB_HOST", "localhost")),
	)

	setup(db)

	http.HandleFunc("/getPersons", func(w http.ResponseWriter, r *http.Request) {
		persons := []Person{}
		rows, err := db.Query("SELECT * FROM persons")

		if err != nil {
			http.Error(w, "Failed to fetch persons", 500)
			return
		}

		for rows.Next() {
			var id int
			var firstname, lastname string
			rows.Scan(&id, &firstname, &lastname)
			persons = append(persons, Person{id, firstname, lastname})
		}

		b, err := json.Marshal(persons)

		if err != nil {
			http.Error(w, "Failed at encoding json", 500)
			return
		}

		setupCors(&w)
		w.Write(b)
	})

	http.HandleFunc("/insertPerson", func(w http.ResponseWriter, r *http.Request) {
		db.QueryRow(
			"INSERT INTO persons (firstname, lastname) VALUES($1,$2)",
			r.PostFormValue("firstname"),
			r.PostFormValue("lastname"),
		)

		setupCors(&w)
		http.Redirect(w, r, "/select.html", http.StatusSeeOther)
	})

	fs := http.FileServer(http.Dir(env("HTTP_VIEWS", "views")))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}

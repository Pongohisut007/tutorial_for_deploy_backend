package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connectdb() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/world")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil

}

type USER struct {
	ID   int    `json:"id"`
	NAME string `json:"name"`
	AGE  int    `json:"age"`
}

func sayhi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "users\n")

	rows, err := db.Query("SELECT id, name, age FROM user")
	if err != nil {
		http.Error(w, "can't pull user from db", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var users []USER
	for rows.Next() {
		var u USER
		if err := rows.Scan(&u.ID, &u.NAME, &u.AGE); err != nil {
			http.Error(w, "can't read data from db", http.StatusInternalServerError)
			log.Println("Scan error:", err)
			return
		}
		users = append(users, u)
	}

	// set response header เป็น JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "can't encode data to json", http.StatusInternalServerError)
		log.Println("JSON encode error:", err)
		return
	}

}

type Car struct {
	ID    int    `json:id`
	NAME  string `json:name`
	COUNT int    `json:count`
}

func car(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "cars")

	rows, err := db.Query("SELECT id, name, count FROM car")
	if err != nil {
		http.Error(w, "can't pull car from db", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var c Car
		if err := rows.Scan(&c.ID, &c.NAME, &c.COUNT); err != nil {
			http.Error(w, "can't read car from db", http.StatusInternalServerError)
			log.Println("Scan error:", err)
			return
		}
		cars = append(cars, c)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cars); err != nil {
		http.Error(w, "can't encode data to json", http.StatusInternalServerError)
		log.Println("JSON encode error:", err)
	}
}

type Hero struct {
	ID     int    `json:id"`
	NAME   string `json:name`
	TOWN   string `json:town`
	WEAPON string `json:weapon`
}

func hero(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hero")

	rows, err := db.Query("SELECT id, name, town, weapon FROM hero")
	if err != nil {
		http.Error(w, "can't pull hero from db", http.StatusInternalServerError)
		log.Println("Query error", err)
	}
	defer rows.Close()

	var heros []Hero
	for rows.Next() {
		var h Hero
		if err := rows.Scan(&h.ID, &h.NAME, &h.TOWN, &h.WEAPON); err != nil {
			http.Error(w, "can't read hero from db", http.StatusInternalServerError)
			log.Println("Scan error: ", err)
		}
		heros = append(heros, h)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(heros); err != nil {
		http.Error(w, "can't encode data to json", http.StatusInternalServerError)
		log.Println("JSON encode error:", err)
	}
}

func main() {
	var err error
	db, err = connectdb()
	if err != nil {
		log.Fatal("can't connect to db")
	}
	log.Println("conect db success")
	defer db.Close()

	http.HandleFunc("/hi", sayhi)
	http.HandleFunc("/car", car)
	http.HandleFunc("/hero", hero)

	s := &http.Server{
		Addr: ":8000",
	}
	log.Fatal(s.ListenAndServe())

}

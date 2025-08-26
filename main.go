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
	fmt.Fprint(w, "hi\n")

	rows, err := db.Query("SELECT id, name, age FROM nongao")
	if err != nil {
		http.Error(w, "can't pull data from db", http.StatusInternalServerError)
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
		
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "can't encode data to json", http.StatusInternalServerError)
		log.Println("JSON encode error:", err)
	}
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

	s := &http.Server{
		Addr: ":8000",
	}
	log.Fatal(s.ListenAndServe())

}

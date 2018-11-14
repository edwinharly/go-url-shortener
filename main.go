package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/teris-io/shortid"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	db *sql.DB
)

type ShortenedURLResponse struct {
	Original  string `json:"original_url"`
	Shortened string `json:"code"`
}

func main() {
	fmt.Println("URL Shortener with Go")

	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	checkErr(err)
	shortid.SetDefault(sid)

	db = openDb()
	defer db.Close()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", rootHandler)      // set router
	http.HandleFunc("/login", login)       // set router
	http.HandleFunc("/register", register) // set router
	http.HandleFunc("/new", shortener)

	port := os.Getenv("PORT")

	log.Fatal(
		http.ListenAndServe(":"+port, nil), // set listening port
	)
}

func openDb() *sql.DB {
	dbinfo := "root:edwinharly@tcp(127.0.0.1:3306)/url_shortener"
	db, err := sql.Open("mysql", dbinfo)
	checkErr(err)
	return db
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/"):]

	if shortURL == "" {
		t, _ := template.ParseFiles("pages/index.html")
		t.Execute(w, nil)
	} else {
		dest := queryShortenedURL(shortURL)
		if dest != "" {
			http.Redirect(w, r, dest, http.StatusSeeOther)
		}
	}
}

func queryShortenedURL(code string) string {
	var (
		dest string
	)
	rows, err := db.Query("select original_url from shortened where shortened_url = ?", code)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&dest)
		checkErr(err)
	}
	err = rows.Err()
	checkErr(err)

	return dest
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		fmt.Println("username:", r.FormValue("username"))
		fmt.Println("password:", r.FormValue("password"))
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("register.gtpl")
		t.Execute(w, nil)
	} else {
		fmt.Println("username:", r.FormValue("username"))
		fmt.Println("password:", r.FormValue("password"))
	}
}

func shortener(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)

	if r.Method == "POST" {
		originalURL := r.FormValue("url")
		short, err := shortid.Generate()
		userID := 0

		stmt, err := db.Prepare("INSERT INTO shortened(original_url, shortened_url, user_id) VALUES(?, ?, ?)")
		checkErr(err)
		defer stmt.Close()

		_, err = stmt.Exec(originalURL, short, userID)
		checkErr(err)

		successResponse := &ShortenedURLResponse{
			Original:  originalURL,
			Shortened: short}
		successResponseJSON, _ := json.Marshal(successResponse)

		w.Write(successResponseJSON)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method Not Allowed"))
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

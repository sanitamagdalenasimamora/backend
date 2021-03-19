package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	//
)

var db *sql.DB
var err error

type user struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func connect_db() {
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1)/doox")
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln()
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	var message = "Berhasil melakukan pendaftaran"
	w.Write([]byte(message))

}

func main() {
	connect_db()
	// routes()
	http.HandleFunc("/", daftar)
	http.HandleFunc("/login", login)

	defer db.Close()
	fmt.Println("Server running on port:8080")
	http.ListenAndServe(":8080", nil)
}

func checkErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {
		fmt.Println(r.Host + r.URL.Path)
		http.Redirect(w, r, r.Host+r.URL.Path, 301)
		return false
	}
	return true
}

func User(username string) user {
	var users = user{}
	err = db.QueryRow(`
	SELECT id,
	email,
	username,
	password
	FROM users WHERE username=?
	`, username).
		Scan(
			&users.ID,
			&users.Email,
			&users.Username,
			&users.Password,
		)
	return users
}

func daftar(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "src/App.vue")
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	users := User(username)

	if (user{}) == users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if len(hashedPassword) != 0 && checkErr(w, r, err) {
			stmt, err := db.Prepare("INSERT INTO users SET  email=?, username=?, password=?")
			if err == nil {
				_, err := stmt.Exec(&email, &username, &hashedPassword)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}

				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

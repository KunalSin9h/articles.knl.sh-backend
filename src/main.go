package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/KunalSin9h/go/src/database"
	"github.com/KunalSin9h/go/src/server"
)

var db = database.CreateDB("sqlite3", os.Getenv("DB"))

/*
Handler Functions
*/

// Home the home page
func Home(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("src/home.html")
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	err = t.Execute(res, nil)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func composeArticle(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	password := req.FormValue("password")
	if password != os.Getenv("PASSWORD") {
		http.Error(res, err.Error(), http.StatusForbidden)
	}

	t, errNew := template.ParseFiles("src/compose.html")

	if errNew != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	errNew = t.Execute(res, nil)

	if errNew != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func addArticle(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	title := req.FormValue("title")
	slug := req.FormValue("slug")
	description := req.FormValue("description")
	date := req.FormValue("date")
	md := req.FormValue("md")

	if title == "" ||
		slug == "" ||
		date == "" ||
		md == "" {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	database.InsertArticle(db, title, slug, description, date, md)

	http.Redirect(res, req, "/", http.StatusSeeOther)

}

func getArticles(res http.ResponseWriter, req *http.Request) {
	var articles []database.Article = database.GetArticles(db)
	res.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(articles)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}
	_, err2 := res.Write(jsonResponse)
	if err2 != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}

}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/compose-article/", composeArticle)
	http.HandleFunc("/add-article/", addArticle)
	http.HandleFunc("/get-articles/", getArticles)
	serverConfig := server.Server{
		Port:    5000,
		Timeout: 3,
	}
	serverConfig.StartServer()
}

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/KunalSin9h/go/src/database"
)

var db = database.CreateDB("sqlite3", os.Getenv("DB"))

/*
Handler Functions
*/

// Enable Cors
func enableCors(res *http.ResponseWriter) {
	// Todo * -> our origin for more security
	(*res).Header().Set("Access-Control-Allow-Origin", "*")
}

// Home the home page
func Home(res http.ResponseWriter, _ *http.Request) {
	enableCors(&res)
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
	enableCors(&res)
	err := req.ParseForm()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	password := strings.TrimSpace(req.FormValue("password"))
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
	enableCors(&res)
	err := req.ParseForm()

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	title := strings.TrimSpace(req.FormValue("title"))
	slug := strings.TrimSpace(req.FormValue("slug"))
	description := strings.TrimSpace(req.FormValue("description"))
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

func getArticles(res http.ResponseWriter, _ *http.Request) {
	enableCors(&res)
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

func getSingleArticle(res http.ResponseWriter, req *http.Request) {
	enableCors(&res)
	slug := req.URL.Query().Get("slug")
	article := database.GetSingleArticle(slug, db)
	res.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(article)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}
	_, err2 := res.Write(jsonResponse)
	if err2 != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}
}

func getAllArticlesMeta(res http.ResponseWriter, req *http.Request) {
	enableCors(&res)
	var articlesMeta []database.ArticleMeta = database.GetArticlesMeta(db)
	res.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(articlesMeta)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}
	_, err2 := res.Write(jsonResponse)
	if err2 != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
	}
}

func getSingleArticleMeta(res http.ResponseWriter, req *http.Request) {
	enableCors(&res)
	slug := req.URL.Query().Get("slug")
	var articleMeta database.ArticleMeta = database.GetSingleArticleMeta(slug, db)
	res.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(articleMeta)
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
	http.HandleFunc("/get-article/", getSingleArticle)
	http.HandleFunc("/get-articles-meta/", getAllArticlesMeta)
	http.HandleFunc("/get-article-meta/", getSingleArticleMeta)
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5005"
	}
	log.Printf("Server running at port: %s", PORT)
	http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil)
}

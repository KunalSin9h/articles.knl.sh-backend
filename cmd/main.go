package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/KunalSin9h/go/cmd/database"
)

var (
	DB       *sql.DB
	DB_PATH  = os.Getenv("DB")
	PORT     = os.Getenv("PORT")
	PASSWORD = os.Getenv("PASSWORD")
)

func init() {

	if PORT == "" {
		log.Println("Use default port: 5000")
		PORT = "5000"
	}

	if DB_PATH == "" {
		log.Println("Use default database: ./database/dev.db")
		DB_PATH = "./database/dev.db"
	}

	if PASSWORD == "" {
		log.Println("Use default password: 1234")
		PASSWORD = "1234"
	}

	if _, err := os.Stat(DB_PATH); err != nil {
		// sqlite3 database file is not present
		folders := filepath.Dir(DB_PATH)
		err = os.MkdirAll(folders, os.ModePerm)
		if err != nil {
			log.Fatalf("Does not able to create folders for database file, try creating the folder manually: %s", err.Error())
		}
		_, err = os.Create(DB_PATH)
		if err != nil {
			log.Fatalf("Does not able to create sqlite3 db file, try creating the file manually: %s", err.Error())
		}
	}

	DB = database.CreateDB("sqlite3", DB_PATH)

	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS articles (
		id varchar(16) primary key,
		title varchar(255) not null,
		slug varchar(255) not null,
		description text,
		date varchar(10) not null,
		views INTEGER DEFAULT 0,
		md text not null
	);
	`)
	if err != nil {
		log.Fatalf("Does not able to create table in the database file, try creating the table manually: %s", err.Error())
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

	log.Printf("Server running at port: %s", PORT)
	http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil)
}

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
	t, err := template.ParseFiles("cmd/home.html")
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
	if password != PASSWORD {
		http.Error(res, err.Error(), http.StatusForbidden)
	}

	t, errNew := template.ParseFiles("cmd/compose.html")

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

	err = database.InsertArticle(DB, title, slug, description, date, md)

	if err != nil {
		http.Error(res, "Something went wrong", http.StatusInternalServerError)
	}

	http.Redirect(res, req, "/get-article?slug="+slug, http.StatusSeeOther)
}

func getArticles(res http.ResponseWriter, _ *http.Request) {
	enableCors(&res)
	var articles []database.Article = database.GetArticles(DB)
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
	article := database.GetSingleArticle(slug, DB)

	defer database.IncrementViews(article.Views, slug, DB)

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
	var articlesMeta []database.ArticleMeta = database.GetArticlesMeta(DB)
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
	var articleMeta database.ArticleMeta = database.GetSingleArticleMeta(slug, DB)
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

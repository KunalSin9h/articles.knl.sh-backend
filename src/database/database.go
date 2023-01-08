package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Article struct {
	Id, Title, Slug, Description, Date, Md string
}

func CreateDB(driver string, source string) *sql.DB {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func CreateTable(db *sql.DB, table string) {
	_, err := db.Exec(table)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func InsertArticle(db *sql.DB, title, slug, description, date, md string) {
	id := slug
	query := "insert into articles (id, title, slug, description, date, md) values (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, id, title, slug, description, date, md)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GetArticles(db *sql.DB) []Article {
	row, err := db.Query("select * from articles")
	if err != nil {
		log.Fatal(err.Error())
	}

	var articles []Article

	for row.Next() {
		art := Article{}
		row.Scan(&art.Id, &art.Title, &art.Slug, &art.Description, &art.Date, &art.Md)
		articles = append(articles, art)
	}
	return articles
}

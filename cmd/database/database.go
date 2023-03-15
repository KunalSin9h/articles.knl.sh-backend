package database

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Article struct {
	Id, Title, Slug, Description, Date, Md string
}

type ArticleMeta struct {
	Id, Title, Slug, Description, Date string
}

func CreateDB(driver string, source string) *sql.DB {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func InsertArticle(db *sql.DB, title, slug, description, date, md string) {
	id := uuid.New()
	query := "insert into articles (id, title, slug, description, date, md) values (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, id, title, slug, description, date, md)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func GetArticles(db *sql.DB) []Article {
	row, err := db.Query("select * from articles")
	if err != nil {
		return []Article{}
	}

	var articles []Article

	for row.Next() {
		art := Article{}
		row.Scan(&art.Id, &art.Title, &art.Slug, &art.Description, &art.Date, &art.Md)
		articles = append(articles, art)
	}
	return articles
}

func GetSingleArticle(slug string, db *sql.DB) Article {
	row := db.QueryRow("select * from articles where slug = ?", slug)

	art := Article{}
	row.Scan(&art.Id, &art.Title, &art.Slug, &art.Description, &art.Date, &art.Md)

	return art
}

func GetArticlesMeta(db *sql.DB) []ArticleMeta {
	row, err := db.Query("select id, title, slug, description, date from articles")

	if err != nil {
		return []ArticleMeta{}
	}

	var artMeta []ArticleMeta

	for row.Next() {
		meta := ArticleMeta{}
		row.Scan(&meta.Id, &meta.Title, &meta.Slug, &meta.Description, &meta.Date)

		artMeta = append(artMeta, meta)
	}

	return artMeta
}

func GetSingleArticleMeta(slug string, db *sql.DB) ArticleMeta {
	row := db.QueryRow("select id, title, slug, description, date from articles where slug = ?", slug)

	var artMeta = ArticleMeta{}

	row.Scan(&artMeta.Id, &artMeta.Title, &artMeta.Slug, &artMeta.Description, &artMeta.Date)

	return artMeta
}

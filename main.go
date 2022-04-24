package main

import (
	"database/sql"
	"encoding/json"
	"go-service/newsfeed"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Content string
}

func main() {
	db, err := sql.Open("sqlite3", "./newsfeed.db")
	if err != nil {
		log.Fatal(err)
	}
	feed := newsfeed.NewFeed(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json", "text/xml"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		items := feed.Get()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
		// w.Write([]byte("Hello World!"))
	})

	r.Post("/post", func(w http.ResponseWriter, r *http.Request) {
		var p Post
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		item := newsfeed.Item{
			Content: p.Content,
		}
		rowid, err := feed.Add(item)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(rowid)
		// fmt.Fprintf(w, "Post: %+v", p)
	})

	address := ":3000"
	log.Println("Starting server on ", address)
	http.ListenAndServe(address, r)
}

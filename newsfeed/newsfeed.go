package newsfeed

import (
	"database/sql"
	"log"
)

type Feed struct {
	DB *sql.DB
}

func NewFeed(db *sql.DB) *Feed {
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "newsfeed" (
		"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"contents" TEXT
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()
	return &Feed{
		DB: db,
	}
}

func (feed *Feed) GetAll() []Item {
	items := []Item{}
	rows, err := feed.DB.Query(`
	SELECT * FROM newsfeed
	`)
	if err != nil {
		log.Fatal(err)
	}
	var id int
	var contents string
	for rows.Next() {
		rows.Scan(&id, &contents)
		item := Item{
			ID:       id,
			Contents: contents,
		}
		items = append(items, item)
	}
	return items
}

func (feed *Feed) Get(rowid int64) (Item, error) {
	item := Item{}
	var id int
	var contents string
	err := feed.DB.QueryRow(`
	SELECT * FROM newsfeed
	WHERE ID=?`, rowid).Scan(&id, &contents)
	if err != nil {
		return item, err
	}
	item = Item{
		ID:       id,
		Contents: contents,
	}
	return item, nil
}

func (feed *Feed) Add(item Item) (int64, error) {
	stmt, err := feed.DB.Prepare(`
INSERT INTO newsfeed (contents) values (?)
`)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(item.Contents)
	if err != nil {
		return 0, err
	}
	rowid, err := result.LastInsertId()
	if err != nil {
		return rowid, err
	}
	return rowid, nil
}

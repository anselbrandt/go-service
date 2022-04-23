```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "elon.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare(`
		SELECT timestamp, tweets.text
		FROM tweets
		INNER JOIN tweetsearch ON tweetsearch.tweetsrowid = tweets.rowid
		WHERE tweetsearch.text MATCH ?;`)
	if err != nil {
		log.Fatal(err)
	}

	var timestamp, text string
	rows, err = stmt.Query("dogecoin").Scan(&timestamp, &text)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Printf("%s %s\n", timestamp, text)
}

// go build --tags "fts5" .
```

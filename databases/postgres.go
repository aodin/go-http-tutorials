package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Item struct {
	Id   int64
	Name string
}

func (item *Item) String() string {
	return fmt.Sprintf("%d: %s", item.Id, item.Name)
}

// Create table syntax must be run separately
var CREATE_TABLE_ITEMS = `CREATE TABLE "items" (
	"id" SERIAL PRIMARY KEY,
	"name" varchar(256) NOT NULL
);`

func main() {
	// Connect to the database
	db, err := sql.Open("postgres", "host=localhost port=5432 dbname=gotest user=postgres password=gotest")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a new item with an auto-increment id
	item := &Item{Name: "Bass-o-matic"}
	err = db.QueryRow(`INSERT INTO "items" ("name") VALUES ($1) RETURNING "id"`, &item.Name).Scan(&item.Id)
	if err != nil {
		panic(err)
	}
	log.Print(item)

	// Select an item
	var bass Item
	err = db.QueryRow(`SELECT "id", "name" FROM "items" WHERE "id" = $1`, &item.Id).Scan(&bass.Id, &bass.Name)
	if err != nil {
		panic(err)
	}
	log.Println(&bass)

	// Add another item without returning the id
	_, err = db.Exec(`INSERT INTO "items" (name) VALUES ($1)`, "Jam Hawkers")
	if err != nil {
		panic(err)
	}

	// Update an item
	_, err = db.Exec(`UPDATE "items" SET "name" = $1 WHERE "id" = $2`, "Bass-o-matic 1976", item.Id)
	if err != nil {
		panic(err)
	}

	// Select multiple items
	rows, err := db.Query(`SELECT "id", "name" FROM "items"`)
	if err != nil {
		panic(err)
	}

	items := make([]*Item, 0)
	for rows.Next() {
		var i Item
		if err = rows.Scan(&i.Id, &i.Name); err != nil {
			panic(err)
		}
		items = append(items, &i)
	}
	log.Println(items)

	// Delete all
	_, err = db.Exec(`DELETE FROM "items"`)
	if err != nil {
		panic(err)
	}
}

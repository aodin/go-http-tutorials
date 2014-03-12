package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

type Item struct {
	Id   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string
}

func (item *Item) String() string {
	return item.Name
}

func main() {
	// Connect to a local instance
	session, err := mgo.Dial(`localhost:27017`)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional - switches the session to a monotonic behavior.
	// See more: http://godoc.org/labix.org/v2/mgo#Session.SetMode
	session.SetMode(mgo.Monotonic, true)

	// Connect to the collection "items" on database "test"
	// This operation is lightweight and does not require network activity
	c := session.DB("test").C("items")

	// Create an item and insert
	cowbell := &Item{Id: bson.NewObjectId(), Name: "More Cowbell"}
	err = c.Insert(cowbell)
	if err != nil {
		panic(err)
	}

	// Get the item by name
	byName := Item{}
	err = c.Find(bson.M{"name": "More Cowbell"}).One(&byName)
	if err != nil {
		panic(err)
	}
	log.Println(byName)

	// Get the item by id
	byId := Item{}
	err = c.FindId(cowbell.Id).One(&byId)
	if err != nil {
		panic(err)
	}
	log.Println(byId)

	// Add another item
	c.Insert(&Item{Id: bson.NewObjectId(), Name: "Don't Fear the Reaper"})

	// Get all items
	var items []*Item
	iter := c.Find(nil).Iter()
	err = iter.All(&items)
	if err != nil {
		panic(err)
	}
	log.Println(items)

	// Delete one item by id
	err = c.RemoveId(cowbell.Id)
	if err != nil {
		panic(err)
	}

	// Delete all items
	err = c.Remove(nil)
	if err != nil {
		panic(err)
	}
}

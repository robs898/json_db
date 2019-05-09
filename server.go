package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Birthday struct {
	FirstName, LastName string
	Day, Month, Year int
}

type DB struct {
	User		string
	Birthdays	[]Birthday
}

func (db *DB) save() error {
	filename := "data/" + db.User + ".json"
	json, _ := json.Marshal(db)
	return ioutil.WriteFile(filename, json, 0600)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	//body := r.FormValue("body")
	decoder := json.NewDecoder(r.Body)
	var birthday Birthday
	err := decoder.Decode(&birthday)
	if err != nil {
		panic(err)
	}
	birthdays := []Birthday
	append.(birthdays, birthday)
	db := &DB{User: title, Birthdays: birthdays)}
	err := db.save()
	//if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
	//http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var validPath = regexp.MustCompile("^/(save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

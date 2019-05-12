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

func view(user string) ([]byte, error) {
	filename := "data/" + user + ".json"
	log.Println(filename)
	return ioutil.ReadFile(filename)
}

func (db *DB) save() error {
	filename := "data/" + db.User + ".json"
	json, _ := json.Marshal(db)
	return ioutil.WriteFile(filename, json, 0600)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) DB {
	byte_file, err := view(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var db DB
	json.Unmarshal(byte_file, &db)
	http.Return
	return
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	decoder := json.NewDecoder(r.Body)
	var birthday Birthday
	err := decoder.Decode(&birthday)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var birthdays []Birthday
	birthdays = append(birthdays, birthday)
	log.Println(birthdays)
	db := DB{
		User: title,
		Birthdays: birthdays,
	}
	err2 := db.save()
	if err2 != nil {
		panic(err)
	}
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
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/view/", ServeFile(w http.ResponseWriter, *http.Request, string))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

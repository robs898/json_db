package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"io/ioutil"
)

type Birthday struct {
	Name string
}

type Birthdays []Birthday

func getUser(w http.ResponseWriter, r *http.Request) string {
	var validPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.Error(w, "Invalid user", 400)
		return ""
	}
	return m[1]
}

func getDb(w http.ResponseWriter, user string, ignoreError bool) Birthday {
	var db Birthday
	filename := "data/" + user + ".json"
	byteFile, err := ioutil.ReadFile(filename)
	if err != nil {
		if ! ignoreError {
			http.Error(w, "User db not found", 404)
		}
		return db
	}
	err2 := json.Unmarshal(byteFile, &db)
	if err2 != nil {
		http.Error(w, "User db invalid json", 500)
		return db
	}
	return db
}

func parseData() {
	decoder := json.NewDecoder(r.Body)
	        var birthday Birthday
		        err := decoder.Decode(&birthday)
}

func createDb(bday Birthday, user string) {
	filename := "data/" + user + ".json"
	json, _ := json.Marshal(bday)
	return ioutil.WriteFile(filename, json, 0600)
}

func mainHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := getUser(w, r)
		log.Println("Current user:", user)
		bday := getDb(w, user, true)
		log.Println(bday)

	} else if r.Method == "GET" {
		user := getUser(w, r)
		log.Println(r.Method)
		bday := getDb(w, user, false)
		log.Println(bday)
	}
	//birthdays := Birthdays{
	//	Birthday{Name: "dave"},
	//	Birthday{Name: "bob"},
	//}
	//json.NewEncoder(w).Encode(birthdays)
}

func main() {
	http.HandleFunc("/", mainHandle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

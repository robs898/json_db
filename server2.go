package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Birthday struct {
	Name string
}

type Birthdays []Birthday

func getUser(w http.ResponseWriter, r *http.Request) string {
	var validPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		return ""
	}
	return m[1]
}

func getDB(w http.ResponseWriter, user string) (Birthdays, error) {
	var db Birthdays
	filename := "data/" + user + ".json"
	byteFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("readfile fail")
		return db, err
	}
	err2 := json.Unmarshal(byteFile, &db)
	if err2 != nil {
		log.Println("marshal fail")
		return db, err2
	}
	return db, nil
}

func parseData(w http.ResponseWriter, r *http.Request) (Birthday, error) {
	decoder := json.NewDecoder(r.Body)
	var birthday Birthday
	err := decoder.Decode(&birthday)
	if err != nil {
		return birthday, err
	}
	return birthday, nil
}

func writeDB(bdays Birthdays, user string) error {
	filename := "data/" + user + ".json"
	json, _ := json.Marshal(bdays)
	return ioutil.WriteFile(filename, json, 0600)
}

func mainHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO", r.Method, r.URL.Path)
	if r.Method == "POST" {
		user := getUser(w, r)
		db, err := getDB(w, user)
		if err != nil {
			log.Println("INFO no existing db for user", user)
		}
		log.Println("existing db", db)
		data, err2 := parseData(w, r)
		if err2 != nil {
			log.Println("ERROR invalid JSON", data)
			http.Error(w, "invalid JSON", 400)
		}
		log.Println("new data", data)
		bdays := append(db, data)
		log.Println("all bdays", bdays)
		err3 := writeDB(bdays, user)
		if err3 != nil {
			log.Println("ERROR failed to write bdays to file", bdays)
			http.Error(w, "failed to write db", 500)
		} else {
			json.NewEncoder(w).Encode(bdays)
		}
	} else if r.Method == "GET" {
		user := getUser(w, r)
		if user == "" {
			log.Println("ERROR invalid user", user)
			http.Error(w, "invalid user", 400)
		} else {
			db, err := getDB(w, user)
			if err != nil {
				log.Println("ERROR no db found for", user)
				http.Error(w, "no found db for user", 404)
			} else {
				json.NewEncoder(w).Encode(db)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", mainHandle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

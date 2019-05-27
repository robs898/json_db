package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

type Birthday struct {
	FirstName string
	LastName  string
	Day       int
	Month     int
	Year      int
}

type Birthdays []Birthday

func getUser(w http.ResponseWriter, r *http.Request) string {
	var validPath = regexp.MustCompile("^/(api/)?([a-zA-Z0-9]+)$")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		return ""
	}
	return m[2]
}

func getDB(w http.ResponseWriter, user string) (Birthdays, error) {
	var db Birthdays
	filename := user + ".json"
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
	filename := user + ".json"
	json, _ := json.Marshal(bdays)
	return ioutil.WriteFile(filename, json, 0600)
}

func apiHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO API", r.Method, r.URL.Path)
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
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

func htmlHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO HTML", r.Method, r.URL.Path)
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	var templates = template.Must(template.ParseFiles("user.html"))
	if r.Method == "GET" {
		user := getUser(w, r)
		if user == "" {
			http.ServeFile(w, r, "index.html")
		} else {
			db, err := getDB(w, user)
			if err != nil {
				log.Println("ERROR no db found for", user)
				http.Error(w, "no found db for user", 404)
			} else {
				err := templates.ExecuteTemplate(w, "user.html", db)
				if err != nil {
					log.Println("ERROR failed to render template", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", htmlHandle)
	mux.HandleFunc("/api/", apiHandle)
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	srv := &http.Server{
		Addr:         ":8443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

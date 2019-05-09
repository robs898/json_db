package main

import (
	"encoding/json"
	"io/ioutil"
)

type Birthday struct {
	Day, Month, Year int
}

type BirthdayEntry struct {
	FirstName, LastName string
	Birthday            Birthday
}

func main() {
	data := BirthdayEntry{
		FirstName: "Mark",
		LastName:  "Jones",
		Birthday: Birthday{
			Day:   1,
			Month: 1,
			Year:  1911,
		},
	}

	file, _ := json.Marshal(data)

	_ = ioutil.WriteFile("test.json", file, 0644)
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDay  int    `json:"birth_date"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
}
type Users struct {
	Users []User `json:"users"`
}

type Location struct {
	Id       int    `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance int    `json:"distance"`
}

type Locations struct {
	Location []Location `json:"locations"`
}

type Visit struct {
	Id       int `json:"id"`
	Location int `json:"location"`
	User     int `json:"user"`
	Visited  int `json:"visited_at"`
	Mark     int `json:"mark"`
}

type Visits struct {
	Visit []Visit `json:"visits"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/{entity}/{id}/", updateByID).Methods("POST")
	router.HandleFunc("/{entity}/{id}/", getByID).Methods("GET")
	router.HandleFunc("/user/{id}/visits/", getUserVisitsByID).Methods("GET")
	router.HandleFunc("/location/{id}/avg/", getAVGLocations).Methods("GET")
	http.ListenAndServe(":8080", router)
}

func updateByID(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	strID := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(strID)
	fileName := mux.Vars(r)["entity"]

	switch fileName {
	case "user":
		fileName = "data/users_1.json"
		file, _ := ioutil.ReadFile(fileName)
		users := Users{}
		_ = json.Unmarshal([]byte(file), &users)

		u := User{}
		_ = decoder.Decode(&u)

		for i := range users.Users {
			if users.Users[i].Id == id {
				if u.BirthDay != 0 {
					users.Users[i].BirthDay = u.BirthDay
				}
				if u.Email != "" {
					users.Users[i].Email = u.Email
				}
				if u.FirstName != "" {
					users.Users[i].FirstName = u.FirstName
				}
				if u.Gender != "" {
					users.Users[i].Gender = u.Gender
				}
				if u.LastName != "" {
					users.Users[i].LastName = u.LastName
				}
				break
			}
		}

		data, _ := json.Marshal(users)
		files, _ := os.Create(fileName)

		defer files.Close()
		files.Write(data)

	case "locations":
		fileName = "data/locations_1.json"
		file, _ := ioutil.ReadFile(fileName)
		locations := Locations{}

		l := Location{}
		_ = decoder.Decode(&l)

		_ = json.Unmarshal([]byte(file), &locations)

		for i := range locations.Location {
			if locations.Location[i].Id == id {
				if l.City != "" {
					locations.Location[i].City = l.City
				}
				if l.Country != "" {
					locations.Location[i].Country = l.Country
				}
				if l.Place != "" {
					locations.Location[i].Place = l.Place
				}
				if l.Distance != 0 {
					locations.Location[i].Distance = l.Distance
				}
				break
			}
		}

		data, _ := json.Marshal(locations)
		files, _ := os.Create(fileName)

		defer files.Close()
		files.Write(data)

	case "visit":
		fileName = "data/visits_1.json"
		file, _ := ioutil.ReadFile(fileName)
		visits := Visits{}

		v := Visit{}
		_ = decoder.Decode(&v)

		_ = json.Unmarshal([]byte(file), &visits)

		for i := range visits.Visit {
			if visits.Visit[i].Id == id {
				if v.Location != 0 {
					visits.Visit[i].Location = v.Location
				}
				if v.Mark != 0 {
					visits.Visit[i].Mark = v.Mark
				}
				if v.User != 0 {
					visits.Visit[i].User = v.User
				}
				if v.Visited != 0 {
					visits.Visit[i].Visited = v.Visited
				}
				break
			}
		}

		data, _ := json.Marshal(visits)
		files, _ := os.Create(fileName)

		defer files.Close()
		files.Write(data)
	}
}

func getByID(w http.ResponseWriter, r *http.Request) {
	strID := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(strID)
	fileName := mux.Vars(r)["entity"]

	User := User{}
	Location := Location{}
	Visit := Visit{}

	switch fileName {
	case "user":
		fileName = "data/users_1.json"
		file, _ := ioutil.ReadFile(fileName)
		users := Users{}
		_ = json.Unmarshal([]byte(file), &users)

		for _, user := range users.Users {
			if user.Id == id {
				User = user
				break
			}
		}
		data, _ := json.Marshal(User)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	case "location":
		fileName = "data/locations_1.json"
		file, _ := ioutil.ReadFile(fileName)
		locations := Locations{}

		_ = json.Unmarshal([]byte(file), &locations)

		for _, location := range locations.Location {
			if location.Id == id {
				Location = location
				break
			}
		}
		data, _ := json.Marshal(Location)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	case "visit":
		fileName = "data/visits_1.json"
		visits := Visits{}

		file, _ := ioutil.ReadFile(fileName)
		_ = json.Unmarshal([]byte(file), &visits)

		for _, visit := range visits.Visit {
			if visit.Id == id {
				Visit = visit

			}
		}
		data, _ := json.Marshal(Visit)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func getUserVisitsByID(w http.ResponseWriter, r *http.Request) {
	strID := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(strID)

	visits := Visits{}
	v := Visits{}

	file, _ := ioutil.ReadFile("data/visits_1.json")
	_ = json.Unmarshal([]byte(file), &visits)

	for _, visit := range visits.Visit {
		if visit.User == id {
			v.Visit = append(v.Visit, visit)
		}
	}

	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func getAVGLocations(w http.ResponseWriter, r *http.Request) {
	strID := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(strID)

	var count, value float32

	file, _ := ioutil.ReadFile("data/visits_1.json")
	visits := Visits{}
	_ = json.Unmarshal([]byte(file), &visits)

	for _, visit := range visits.Visit {
		if visit.Location == id {
			count++
			value += float32(visit.Mark)
		}
	}

	avg := make(map[string]float32)

	avg["avg"] = value / count

	data, _ := json.Marshal(avg)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

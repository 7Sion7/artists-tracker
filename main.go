package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

var Rude_Info INFO

var DT Date_Loc

var artistPage *AP

var errs *error

var tpl = template.Must(template.ParseGlob("Templates/*.html"))

type INFO []struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Date_Loc struct {
	ID int `json:"id"`

	DatesLocations map[string][]string `json:"datesLocations"`
}

type AP struct {
	Image           string
	Name            string
	Members         []string
	CreationDate    int
	FirstAlbum      string
	DatesNLocations map[string][]string
}

func main() {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		errs = &err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&Rude_Info)
	if err != nil {
		errs = &err
	}
	http.HandleFunc("/", This)
	http.ListenAndServe(":6789", nil)
}

func This(w http.ResponseWriter, r *http.Request) {
	if errs != nil {
		fmt.Fprint(w, "<html><body><h1>500 INTERNAL SERVER ERROR</h1></body></html>")
		return
	}
	if r.URL.Path == "/" {
		err := tpl.ExecuteTemplate(w, "index.html", Rude_Info)
		if err != nil {
			fmt.Fprint(w, "<html><body><h1>500 INTERNAL SERVER ERROR</h1></body></html>")
			return
		}
	} else {
		index := Checker(r)
		if index == 53 {
			fmt.Fprint(w, "<html><body><h1>404 SERVER NOT FOUND</h1></body></html>")
			return
		} else if index == 54 {
			fmt.Fprint(w, "<html><body><h1>500 INTERNAL SERVER ERROR</h1></body></html>")
			return
		}
		err := tpl.ExecuteTemplate(w, "artist.html", artistPage)
		if err != nil {
			fmt.Fprint(w, "<html><body><h1>500 INTERNAL SERVER ERROR</h1></body></html>")
			return
		}
	}
}

func Checker(r *http.Request) int {
	NAMEID := r.URL.Path[1:]
	numeric := regexp.MustCompile(`\d`).MatchString(NAMEID)
	if !numeric {
		for i := 0; i < len(Rude_Info); i++ {
			if NAMEID == Rude_Info[i].Name {
				la := i + 1
				id := strconv.Itoa(la)
				err := ArtistPage(id, i)
				if err != nil {
					return 54
				}
				return i
			}
		}
	} else {
		id, _ := strconv.Atoi(NAMEID)
		if id >= 1 || id <= 52 {
			ind := id - 1
			err := ArtistPage(NAMEID, ind)
			if err != nil {
				return 54
			}
			return id - 1
		}
		return 53
	}
	return 53
}

func ArtistPage(ID string, index int) error {
	err := DL(ID)
	if err != nil {
		return err
	}
	artistPage = &AP{
		Image:           Rude_Info[index].Image,
		Name:            Rude_Info[index].Name,
		Members:         Rude_Info[index].Members,
		CreationDate:    Rude_Info[index].CreationDate,
		FirstAlbum:      Rude_Info[index].FirstAlbum,
		DatesNLocations: DT.DatesLocations,
	}
	return nil
}

func DL(ID string) error {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + ID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&DT)
	if err != nil {
		return err
	}
	return nil
}

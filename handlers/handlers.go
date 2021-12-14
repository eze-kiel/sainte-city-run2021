package handlers

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type participant struct {
	name string
	bib  string
	sex  string
	time time.Duration
}

type Response struct {
	Name string
	Bib  string
	Sex  string
	Time string
	Rank string
}

func HandleFunc() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/rank", rankPage).Methods("POST")

	return r
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/home.html")
	if err != nil {
		log.Fatalf("Can not parse home page : %v", err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Can not execute templates for home page : %v", err)
	}
}

func rankPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/rank.html")
	if err != nil {
		log.Fatalf("Can not parse rank page : %v", err)
	}
	bibFromForm := r.FormValue("bib")

	fd, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	var s = bufio.NewScanner(fd)
	var all = []participant{}
	for s.Scan() {
		line := s.Text()
		if !strings.Contains(line, "°") {
			continue
		}
		all = append(all, parseParticipantInfos(line))
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].time < all[j].time
	})
	var resp Response
	for i, p := range all {
		bibWanted, err := strconv.Atoi(bibFromForm)
		if err != nil {
			if resp == (Response{}) {
				resp.Name = "n/a"
				resp.Bib = "n/a"
				resp.Sex = "n/a"
				resp.Time = "n/a"
				resp.Rank = "n/a"
			}
			break
		}
		bib, err := strconv.Atoi(p.bib)
		if err != nil {
			panic(err)
		}
		if bibWanted == bib {
			showParticipantInfos(p, i+1, len(all))
			resp.Name = p.name
			resp.Bib = p.bib
			resp.Sex = p.sex
			resp.Time = p.time.String()
			resp.Rank = strconv.Itoa(i+1) + "/" + strconv.Itoa(len(all))
			break
		}
	}

	if resp == (Response{}) {
		resp.Name = "n/a"
		resp.Bib = "n/a"
		resp.Sex = "n/a"
		resp.Time = "n/a"
		resp.Rank = "n/a"
	}

	err = tmpl.Execute(w, resp)
	if err != nil {
		log.Fatalf("Can not execute templates for rank page : %v", err)
	}
}

func parseParticipantInfos(line string) participant {
	var p participant
	var err error
	parts := strings.Split(line, " ")
	for i, v := range parts {
		if v == "N°" {
			p.bib = parts[i+1]
			p.sex = parts[i+2]
			if strings.Count(parts[i+3], ":") == 2 {
				half := strings.Replace(parts[i+3], ":", "h", 1)
				p.time, err = time.ParseDuration(strings.ReplaceAll(half, ":", "m") + "s")
				if err != nil {
					panic(err)
				}
			} else {
				p.time, err = time.ParseDuration(strings.ReplaceAll(parts[i+3], ":", "m") + "s")
				if err != nil {
					panic(err)
				}
			}
			break
		}
		if p.name == "" {
			p.name = parts[i]
		} else {
			p.name = p.name + " " + parts[i]
		}
	}
	return p
}

func showParticipantInfos(p participant, rank, tot int) {
	fmt.Printf(`
Name: %s
Bib:  %s
Sex:  %s
Time: %s
Rank: %d/%d
	`, p.name, p.bib, p.sex, p.time, rank, tot)
}

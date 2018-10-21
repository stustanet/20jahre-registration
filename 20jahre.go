package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	minAdults   = 1
	maxAdults   = 5
	minChildren = 0
	maxChildren = 5
)

var (
	templates map[string]*template.Template
	db        *sql.DB
)

func loadTemplates() {
	templates = make(map[string]*template.Template)

	base, err := template.ParseFiles("tpl/base.html")
	if err != nil {
		log.Fatal(err)
	}

	for _, name := range []string{"form", "thanks"} {
		templates[name], err = base.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[name] = template.Must(templates[name].ParseFiles("tpl/" + name + ".html"))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := templates["form"].ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println(err)
		}
	case "POST":
		r.ParseForm()
		f := r.Form

		valid := true

		var name string
		var attend bool
		if len(f["name"]) == 1 && len(f["attend"]) == 1 {
			name = f["name"][0]
			if len(name) < 6 || len(name) > 60 {
				valid = false
			}
			attend = (f["attend"][0] == "1")
		} else {
			valid = false
		}

		if !valid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var comment string
		if len(f["comment"]) == 1 {
			comment = f["comment"][0]
		}

		var adults int
		var adultsVeg int
		var children int
		var childrenVeg int
		if attend {
			if len(f["adults"]) == 1 && len(f["adults_veg"]) == 1 && len(f["children"]) == 1 && len(f["children_veg"]) == 1 {
				var err error
				adults, err = strconv.Atoi(f["adults"][0])
				if err != nil {
					valid = false
				}

				adultsVeg, err = strconv.Atoi(f["adults_veg"][0])
				if err != nil {
					valid = false
				}

				children, err = strconv.Atoi(f["children"][0])
				if err != nil {
					valid = false
				}

				childrenVeg, err = strconv.Atoi(f["children_veg"][0])
				if err != nil {
					valid = false
				}

				if valid && (adults < minAdults || adults > maxAdults || adultsVeg < 0 ||
					children < minChildren || children > maxChildren || childrenVeg < 0 ||
					adultsVeg > adults || childrenVeg > children) {
					valid = false
				}
			} else {
				valid = false
			}

			if !valid {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}

		_, err := db.Exec("INSERT INTO registration(name, attending, adults, adults_veg, children, children_veg, comment) VALUES ($1, $2, $3, $4, $5, $6, $7)", name, attend, adults, adultsVeg, children, childrenVeg, comment)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = templates["thanks"].ExecuteTemplate(w, "base", nil)
		if err != nil {
			log.Println(err)
		}

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	var err error
	db, err = sql.Open("postgres", "user=20jahre dbname=20jahre password=0KS_(Z5")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db.Ping())

	loadTemplates()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

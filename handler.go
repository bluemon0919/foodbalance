package main

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func Sample(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("input.html"))
	tpl.Execute(w, nil)
}

// Astois is equivalent to Atoi for slice.
func Astois(ss []string) ([]int, error) {
	var is []int
	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		is = append(is, i)
	}
	return is, nil
}

func SamplePost(w http.ResponseWriter, r *http.Request) {
	// HTTPメソッドをチェック（POSTのみ許可）
	if r.Method != http.MethodPost {
		return
	}
	r.ParseForm()

	name := r.Form["Name"][0]
	if len(name) == 0 {
		return
	}
	timeZone, err := strconv.Atoi(r.Form["TimeZone"][0])
	if err != nil {
		timeZone = 0
	}
	is, err := Astois(r.Form["Group"])
	if err != nil {
		return
	}
	t := time.Now()
	ts := t.Format(dateFormat)
	m := Create(ts, name, timeZone, is[0], is[1], is[2], is[3], is[4])
	Put(m)
}

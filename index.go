package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg", newsAggHandler)
	http.ListenAndServe(":8000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Go is really really neat</h1>")
}

func newsAggHandler(w http.ResponseWriter, r *http.Request) {

	var s SitemapIndex
	var n News
	newsMap := make(map[string]NewsMap)

	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	xml.Unmarshal(bytes, &s)

	for _, Location := range s.Locations {
		resp, err := http.Get(strings.TrimSpace(Location))
		if err != nil {
			fmt.Println("Response :", resp)
			fmt.Println("Error :", err)
			return
		}
		bytes, _ := ioutil.ReadAll(resp.Body)

		resp.Body.Close()

		xml.Unmarshal(bytes, &n)
		for idx, _ := range n.Titles {
			newsMap[n.Titles[idx]] = NewsMap{n.Keywords[idx], n.Locations[idx]}
		}
		resp.Body.Close()
	}

	p := newsAggPage{Title: "Latest News", News: newsMap}
	t, _ := template.ParseFiles("basictemplate.html")
	t.Execute(w, p)
}

type newsAggPage struct {
	Title string
	News  map[string]NewsMap
}

type SitemapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles    []string `xml:"url>news>title"`
	Keywords  []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type NewsMap struct {
	Keyword  string
	Location string
}

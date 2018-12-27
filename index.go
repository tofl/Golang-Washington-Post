package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg", newsAggHandler)
	http.ListenAndServe(":8000", nil)
}

func newsRoutine(c chan News, Location string) {
	defer wg.Done()
	var n News

	resp, err := http.Get(strings.TrimSpace(Location))
	if err != nil {
		fmt.Println("Response :", resp)
		fmt.Println("Error :", err)
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &n)
	resp.Body.Close()

	c <- n
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Go is really really neat</h1>")
}

func newsAggHandler(w http.ResponseWriter, r *http.Request) {

	var s SitemapIndex
	newsMap := make(map[string]NewsMap)

	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	xml.Unmarshal(bytes, &s)

	queue := make(chan News, 30)

	for _, Location := range s.Locations {
		wg.Add(1)
		go newsRoutine(queue, Location)
	}

	wg.Wait()
	close(queue)

	for elem := range queue {
		for idx, _ := range elem.Keywords {
			newsMap[elem.Titles[idx]] = NewsMap{elem.Keywords[idx], elem.Locations[idx]}
		}
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

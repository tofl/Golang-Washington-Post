package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/xml"
)

func main() {
	var s SitemapIndex
	var n News

	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	xml.Unmarshal(bytes, &s)
	
	for _, Location := range s.Locations {
		resp, _ := http.Get(Location)
		bytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		xml.Unmarshal(bytes, &n)
	}
}

type SitemapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/xml"
	"strings"
)

func main() {
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
	
	
	for idx, data := range newsMap {
		fmt.Println("\n\n\n", idx)
		fmt.Println("\n", data.Keyword)
		fmt.Println("\n", data.Location)
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

type NewsMap struct {
	Keyword string
	Location string
}


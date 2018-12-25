package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/xml"
)

func main() {
	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var s SitemapIndex
	xml.Unmarshal(bytes, &s)
	
	for _, Location := range s.Locations {
		fmt.Printf("%s", Location)
	}
}

type SitemapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}


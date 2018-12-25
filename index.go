package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	//"encoding/xml"
)

func main() {
	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	stringBody := string(bytes)
	resp.Body.Close()
	fmt.Println(stringBody)
}


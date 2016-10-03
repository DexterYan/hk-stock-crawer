package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	type Stock struct {
		Data string `xml:"data_sets>data_set>csv>data"`
	}
	url := "http://www.sl886.com/post/flashchart?c=01538"
	resp, err := http.Get(url)
	if err != nil {

	}

	s := Stock{Data: ""}
	body, _ := ioutil.ReadAll(resp.Body)
	err = xml.Unmarshal(body, &s)
	if err != nil {

	}
	fmt.Print(string(s.Data))
	fmt.Println("vim-go")
}

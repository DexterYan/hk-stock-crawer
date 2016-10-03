package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// type Stock struct {
	// 	Data string `xml:"data_sets>data_set>csv>data"`
	// }
	url := "http://money18.on.cc/js/daily/hk/quote/01538_d.js?t=1475503151919"
	resp, err := http.Get(url)
	if err != nil {

	}

	// s := Stock{Data: ""}
	// "http://money18.on.cc/chartdata/full/price/01538_price_full.txt"
	// "http://money18.on.cc/info/liveinfo_quote.html?symbol=01538"
	body, _ := ioutil.ReadAll(resp.Body)
	// err = xml.Unmarshal(body, &s)
	// if err != nil {
	//
	// }
	fmt.Print(string(body))
	fmt.Println("vim-go")
}

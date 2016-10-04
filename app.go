package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/axgle/mahonia"
)

func main() {
	type Stock struct {
		NameEn      string  `json:"name"`
		NameCh      string  `json:"nameChi"`
		LastClose   string  `json:"preCPrice"`
		LotSize     string  `json:"lotSize"`
		MthHigh     string  `json:"mthHigh"`
		MthLow      string  `json:"mthLow"`
		Wk52High    string  `json:"wk52High"`
		Wk52Low     string  `json:"wk52Low"`
		Ma10        float32 `json:"ma10"`
		Ma20        float32 `json:"ma20"`
		Ma50        float32 `json:"ma50"`
		Rsi10       float32 `json:"rsi10"`
		Rsi14       float32 `json:"rsi14"`
		Rsi20       float32 `json:"rsi20"`
		Dividend    string  `json:"dividend"`
		Eps         string  `json:"eps"`
		ParentType  string  `json:"parentType"`
		IssuedShare string  `json:"issuedShare"`
	}
	url := "http://money18.on.cc/js/daily/hk/quote/01538_d.js?t=1475503151919"
	resp, err := http.Get(url)
	if err != nil {

	}

	// s := Stock{Data: ""}
	// "http://money18.on.cc/chartdata/full/price/01538_price_full.txt"
	// "http://money18.on.cc/info/liveinfo_quote.html?symbol=01538"
	// "http://money18.on.cc/js/daily/short_put/short_put_00700.js?t=201695"
	body, _ := ioutil.ReadAll(resp.Body)
	enc := mahonia.NewDecoder("big5")
	safe := enc.ConvertString(string(body))
	reg, err := regexp.Compile(".*=")
	safe = reg.ReplaceAllString(safe, "")
	if err != nil {
		log.Fatal(err)
	}
	data := Stock{}
	if err := json.Unmarshal([]byte(safe), &data); err != nil {
		panic(err)
	}
	fmt.Println(data)
}

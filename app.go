package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/axgle/mahonia"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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

func getStockCurrentSummary(code string, time int64) string {
	url := "http://money18.on.cc/js/daily/hk/quote/" + code + "_d.js?t=" + strconv.FormatInt(time, 10)
	fmt.Print(url)
	resp, err := http.Get(url)
	checkError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	decode := mahonia.NewDecoder("big5")
	decodeBody := decode.ConvertString(string(body))
	reg, err := regexp.Compile(".*=")
	checkError(err)
	decodeBody = reg.ReplaceAllString(decodeBody, "")
	stock := Stock{}
	err = json.Unmarshal([]byte(decodeBody), &stock)
	checkError(err)
	return stock.NameCh
}

func main() {
	fmt.Print(getStockCurrentSummary("00700", time.Now().Unix()))

}

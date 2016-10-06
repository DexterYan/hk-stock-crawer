package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

type StockMonth struct {
	Code  string        `json:"StockNo"`
	Trade []interface{} `json:"Trade"`
}

func getStockCurrentSummary(code string, time int64) interface{} {
	var stock Stock
	url := "http://money18.on.cc/js/daily/hk/quote/" + code + "_d.js?t=" + strconv.FormatInt(time, 10)
	fmt.Print(url)
	resp, err := http.Get(url)
	checkError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	decode := mahonia.NewDecoder("big5")
	decodeBody := decode.ConvertString(string(body))
	reg := regexp.MustCompile(".*=")
	decodeBody = reg.ReplaceAllString(decodeBody, "")
	dec := json.NewDecoder(strings.NewReader(decodeBody))
	for {
		var s Stock
		if err := dec.Decode(&s); err == io.EOF {
			break
		} else {
			checkError(err)
		}
		stock = s
	}
	return stock
}

func getStockMonthSummary(code string) interface{} {
	var stockMonth StockMonth
	url := "http://money18.on.cc/js/daily/short_put/short_put_" + code + ".js"
	resp, err := http.Get(url)
	checkError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	reg := regexp.MustCompile(".*=")
	reg1 := regexp.MustCompile(";")
	decodeBody := reg.ReplaceAllString(string(body), "")
	decodeBody = reg1.ReplaceAllString(decodeBody, "")
	dec := json.NewDecoder(strings.NewReader(decodeBody))
	for {
		var s StockMonth
		if err := dec.Decode(&s); err == io.EOF {
			break
		} else {
			checkError(err)
		}
		stockMonth = s
	}
	return stockMonth
}

func getStockList() []string {
	url := "http://money18.on.cc/js/daily/hk/stocklist/stockList_secCode.js"
	resp, err := http.Get(url)
	checkError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	reg := regexp.MustCompile("[0-9]{5}")
	stockList := reg.FindAllString(string(body), -1)
	return stockList
}

func main() {
	fmt.Print(getStockCurrentSummary("00700", time.Now().Unix()))
	fmt.Print(getStockMonthSummary("00700"))
	//fmt.Print(getStockList())
}

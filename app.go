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
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/axgle/mahonia"
)

func checkError(err error) {
	if err != nil {
		if err.Error() != "not found" {
			log.Fatal(err)
		}
	}
}

func startOfDay(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, date.Location())
}

func endOfDay(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 23, 59, 0, 0, date.Location())
}

type Stock struct {
	Code       string
	NameEn     string    `json:"name"`
	NameCh     string    `json:"nameChi"`
	LastClose  float32   `json:"preCPrice,string"`
	LotSize    float32   `json:"lotSize,string"`
	MthHigh    float32   `json:"mthHigh,string"`
	MthLow     float32   `json:"mthLow,string"`
	Wk52High   float32   `json:"wk52High,string"`
	Wk52Low    float32   `json:"wk52Low,string"`
	Ma10       float32   `json:"ma10"`
	Ma20       float32   `json:"ma20"`
	Ma50       float32   `json:"ma50"`
	Rsi10      float32   `json:"rsi10"`
	Rsi14      float32   `json:"rsi14"`
	Rsi20      float32   `json:"rsi20"`
	Dividend   float32   `json:"dividend,string"`
	Eps        float32   `json:"eps,string"`
	ParentType string    `json:"parentType"`
	Timestamp  time.Time `bson:"timestamp"`
}

type StockMonth struct {
	Code  string        `json:"StockNo"`
	Trade []interface{} `json:"Trade"`
}

func getStockCurrentSummary(code string, date int64) Stock {
	var stock Stock
	url := "http://money18.on.cc/js/daily/hk/quote/" + code + "_d.js?t=" + strconv.FormatInt(date, 10)
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
		s.Code = code
		s.Timestamp = time.Now()
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

func readStockList(stockList []string) <-chan string {
	out := make(chan string)
	go func() {
		for _, n := range stockList {
			out <- n
		}
		close(out)
	}()
	return out
}

func saveStockToDB(stock Stock) {
	fmt.Println(stock.Code)
	result := Stock{}
	session, err := mgo.Dial("localhost:27017")
	checkError(err)
	defer session.Close()
	c := session.DB("test").C("stock")
	date := time.Now()
	year, month, day := date.Date()
	query := bson.M{"code": stock.Code, "timestamp": bson.M{
		"$gte": time.Date(year, month, day, 0, 0, 0, 0, date.Location()),
		"$lt":  time.Date(year, month, day+1, 0, 0, 0, 0, date.Location()),
	}}
	err = c.Find(query).One(&result)
	if (Stock{} == result) {
		c.Insert(stock)
	} else {
		// fmt.Print(result)
	}
	checkError(err)
}

func main() {
	var wg sync.WaitGroup
	stock := readStockList(getStockList())
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			for s := range stock {
				saveStockToDB(getStockCurrentSummary(s, time.Now().Unix()))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	// fmt.Print(getStockMonthSummary("00700"))
	//fmt.Print(getStockList())
}

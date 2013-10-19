package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"github.com/dhruvbird/cowsay.go"
)

var quotes []string
var visits int

func readQuotes(file string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	qParts := bytes.Split(b, []byte("\n"))
	for _, line := range qParts {
	    if len(line) > 0 {
	        quotes = append(quotes, string(line))
	    }
	}
}

var index []byte

func readIndex(file string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	index = b
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", index)
	logVisit(r)
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
	// fmt.Printf("%s\n", r.Form)
	quoteStr := getRandQuote()
	qFormats, ok := r.Form["format"]
	if !ok || len(qFormats) == 0 {
	   qFormats = append(qFormats, "text")
	}
	qFormat := qFormats[0]
	switch qFormat {
	    case "cowsay":
	        quoteStr = cowsay.Format(quoteStr)
  	    default:
		quoteStr = fmt.Sprintf("\"%s\"", quoteStr)
	}

	fmt.Fprintf(w, "%s\n", quoteStr)
	logVisit(r)
}

func logVisit(r *http.Request) {
	visits = visits + 1
	fmt.Printf("Visit by: %s\n", r.RemoteAddr)
}

func visitsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Number of visits: %d\n", visits)
}

func getRandQuote() string {
	return quotes[rand.Intn(len(quotes))]
}

func main() {
	visits = 0
	rand.Seed(time.Now().UTC().UnixNano())

	listenPort := flag.Int("port", 80,
		"The HTTP port to listen on (default: 80)")

	flag.Parse()

	readQuotes("quotes.txt")
	readIndex("index.html")
	http.HandleFunc("/", root)
	http.HandleFunc("/quote", quoteHandler)
	http.HandleFunc("/visits", visitsHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(*listenPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

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
		quotes = append(quotes, string(line))
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
	fmt.Fprintf(w, "%s\n", getRandQuote())
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

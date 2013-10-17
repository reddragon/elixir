package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var quotes []string
var visits int

func readQuotes(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	quotes = make([]string, 0)

	r := bufio.NewReader(f)
	for {
		q, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		quotes = append(quotes, q)
	}
}

var index []byte

func readIndex(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	index = make([]byte, 1024)

	for {
		n, err := f.Read(index)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
	}
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
	readQuotes("quotes.txt")
	readIndex("index.html")
	http.HandleFunc("/", root)
	http.HandleFunc("/quote", quoteHandler)
	http.HandleFunc("/visits", visitsHandler)
	http.ListenAndServe(":80", nil)
}

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/dhruvbird/go-cowsay"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
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

const layout = "02/01/2006 15:04:05"

func logVisit(r *http.Request) {
	visits = visits + 1
	fmt.Printf("%s %s\n", time.Now().Format(layout), r.RemoteAddr)
}

func visitsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Number of visits: %d\n", visits)
}

func getRandQuote() string {
	return quotes[rand.Intn(len(quotes))]
}

// The idea is to be able to reload the quotes file without having to shutdown
// the server. I'm using the mtime in the file stats, to figure out if a change
// has happened since we last loaded it. This could have been done better with
// inotify(2) system call, but alas it is not supported on Darwin
func fileChangeListener() {
	fi, err := os.Lstat("quotes.txt")
	if err != nil {
		panic(err)
	}
	mTime := fi.ModTime()

	for {
		time.Sleep(1 * time.Second)
		fi, err := os.Lstat("quotes.txt")
		if err != nil {
			// TODO
			// Log a warning that something went wrong.
			continue
		}
		newMTime := fi.ModTime()
		if newMTime.After(mTime) {
			fmt.Println("Reloading the quotes")
			quotes = quotes[len(quotes):]
			readQuotes("quotes.txt")
			mTime = newMTime
		}
	}
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
	go fileChangeListener()
	err := http.ListenAndServe(":"+strconv.Itoa(*listenPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

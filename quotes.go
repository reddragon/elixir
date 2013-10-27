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
	"strings"
	"time"
)

var quotes []string
var visits int

func readQuotes(file string) []string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	qParts := bytes.Split(b, []byte("\n"))
	quotesList := make([]string, 0)
	for _, line := range qParts {
		if len(line) > 0 {
			quotesList = append(quotesList, string(line))
		}
	}
	return quotesList
}

var index string

func readIndex(file string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	index = string(b)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	quoteEndpoint := strings.Split(r.RequestURI[1:], "?")[0]
	if quoteEndpoint == "" {
		fmt.Fprintf(w, "%s\n", index)
	} else {
		serveRandQuote(quoteEndpoint, w, r)
	}
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

func serveRandQuote(quoteEndpoint string, w http.ResponseWriter, r *http.Request) {
	quoteStr := ""
	if quoteMap[quoteEndpoint] == nil {
		http.NotFound(w, r)
		return
	} else {
		quoteStr = quoteMap[quoteEndpoint][rand.Intn(len(quoteMap[quoteEndpoint]))]
	}

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
}

var mtimeMap map[string]time.Time

func maintainQuotes(c chan bool) {
	fileSet = make(map[string]bool)
	quoteMap = make(map[string][]string)
	mtimeMap = make(map[string]time.Time)

	firstPassDone := false

	for {
		for _, fileName := range getQuotesFiles() {
			if !fileSet[fileName] {
				fmt.Println("Found a new quotes file", fileName)
				quoteEndpoint := getEndpoint(fileName)
				quoteMap[quoteEndpoint] = readQuotes(fileName)
				mtimeMap[fileName] = getMTime(fileName)
				fmt.Println("Loaded quotes from", fileName)
				fileSet[fileName] = true
			}
		}
		if !firstPassDone {
			firstPassDone = true
			// Signal to the main thread that it is safe to serve traffic now
			c <- true
		}

		for fileName, _ := range fileSet {
			fi, err := os.Lstat(fileName)
			if err != nil {
				fmt.Println("Removing", fileName, "from the list, since it is no longer available.")
				delete(fileSet, fileName)
				delete(quoteMap, getEndpoint(fileName))
				continue
			}
			if fi.ModTime().After(mtimeMap[fileName]) {
				fmt.Println("Reloading quotes from", fileName)
				quoteEndpoint := fileName[:len(fileName)-len(".quotes")]
				quoteMap[quoteEndpoint] = readQuotes(fileName)
				mtimeMap[fileName] = fi.ModTime()
			}
		}
		time.Sleep(1 * time.Second)
	}
}

var fileSet map[string]bool
var quoteMap map[string][]string

func getEndpoint(fileName string) string {
	return fileName[:len(fileName)-len(".quotes")]
}

func getQuotesFiles() []string {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	quoteFiles := make([]string, 0)
	for _, file := range files {
		if fileName := file.Name(); strings.HasSuffix(fileName, ".quotes") {
			quoteFiles = append(quoteFiles, fileName)
		}
	}
	return quoteFiles
}

func getMTime(fileName string) time.Time {
	fi, err := os.Lstat(fileName)
	if err != nil {
		panic(err)
	}
	return fi.ModTime()
}

func main() {
	visits = 0
	rand.Seed(time.Now().UTC().UnixNano())

	listenPort := flag.Int("port", 80,
		"The HTTP port to listen on (default: 80)")

	flag.Parse()

	readIndex("index.html")
	http.HandleFunc("/", handler)
	http.HandleFunc("/visits", visitsHandler)
	c := make(chan bool)
	fmt.Println("Starting the quotes maintenance goroutine")
	go maintainQuotes(c)
	fmt.Println("Waiting for the signal to serve traffic")
	safe := <-c
	if safe {
		fmt.Println("It is safe to serve traffic now")
	} else {
		panic("Something bad happened. Dying")
	}

	err := http.ListenAndServe(":"+strconv.Itoa(*listenPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

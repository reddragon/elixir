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
	"strings"
)

var quotes []string
var visits int

func readQuotes(file string) ([]string){
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

// Periodically polling to check the mtime.
// Sad that inotify isn't present on Darwin.
func fileChangeListener() {
	mtimeMap := make(map[string]time.Time)
	for _, fileName := range fileList {
		fi, err := os.Lstat(fileName)
		if err != nil {
			panic(err)
		}
		mtimeMap[fileName] = fi.ModTime()
	}

	for {
		time.Sleep(1 * time.Second)
		for _, fileName := range fileList {
			fi, err := os.Lstat(fileName)
			if err != nil {
				// TODO
				// Log a warning that something went wrong.
				continue
			}
			if fi.ModTime().After(mtimeMap[fileName]) {
				fmt.Println("Reloading quotes from", fileName)
				quoteEndpoint := fileName[:len(fileName) - len(".quotes")]
				quoteMap[quoteEndpoint] = readQuotes(fileName)
				mtimeMap[fileName] = fi.ModTime()
			}
		}
	}
}

var fileList []string
var quoteMap map[string][]string
func loadQuotes() {
	tmpFiles, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	quoteMap = make(map[string][]string)
	for _, file := range tmpFiles {
		if fileName := file.Name(); strings.HasSuffix(fileName, ".quotes") {
			quoteEndpoint := fileName[:len(fileName) - len(".quotes")]
			fileList = append(fileList, fileName)
			quoteMap[quoteEndpoint] = readQuotes(fileName)
			fmt.Println("Loaded quotes from", fileName)
		}
	}
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
	loadQuotes()
	go fileChangeListener()
	err := http.ListenAndServe(":"+strconv.Itoa(*listenPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

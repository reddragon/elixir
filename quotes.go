/*
Copyright (c) 2013, Gaurav Menghani
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
    notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
    notice, this list of conditions and the following disclaimer in the
    documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package elixir

import (
	"bytes"
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

var index string

func readIndex(file string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	index = string(b)
}

// The handler for the quotes and the index page
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

// This method will serve a random quote given the endpoint.
func serveRandQuote(quoteEndpoint string, w http.ResponseWriter, r *http.Request) {
	quoteStr := ""
	if quoteMap[quoteEndpoint] == nil {
		http.NotFound(w, r)
		return
	} else {
		quoteStr = quoteMap[quoteEndpoint][rand.Intn(len(quoteMap[quoteEndpoint]))]
	}

	qFormats, ok := r.Form["f"]
	if !ok || len(qFormats) == 0 {
		qFormats = append(qFormats, "text")
	}
	qFormat := qFormats[0]
	switch qFormat {
	case "cowsay":
		quoteStr = cowsay.Format(quoteStr)
	default:
		quoteStr = fmt.Sprintf("%s", quoteStr)
	}
	fmt.Fprintf(w, "%s\n", quoteStr)
}

const layout = "02/01/2006 15:04:05"

var visits int

// Minimally log a user's visit
func logVisit(r *http.Request) {
	visits = visits + 1
	fmt.Printf("%s %s\n", time.Now().Format(layout), r.RemoteAddr)
}

func visitsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Number of visits: %d\n", visits)
}

var fileSet map[string]bool
var quoteMap map[string][]string
var mtimeMap map[string]time.Time

// Get all the .quotes files in the CWD
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

// Helper method to convert the file name to the API endpoint
// For instance, if the name of the quotes file is "lotr.quotes", the quote
// server will serve quotes at "/lotr" and this method will return "lotr"
func getEndpoint(fileName string) string {
	return fileName[:len(fileName)-len(".quotes")]
}

// Get the time at which a file was last modified.
func getMTime(fileName string) time.Time {
	fi, err := os.Lstat(fileName)
	if err != nil {
		panic(err)
	}
	return fi.ModTime()
}

// Read and return the quotes from a file
func readQuotes(file string) []string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	qParts := bytes.Split(b, []byte("\n%"))
	quotesList := make([]string, 0)
	for _, line := range qParts {
		if len(line) > 0 {
			quotesList = append(quotesList, string(line))
		}
	}
	return quotesList
}

// This is the goroutine which will run forever, loading in the quotes from the
// quotes files present in your CWD, deleting them when they get removed and
// adding them when they get created, without you needing to have to restart
// the server.
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

// The method that starts up the server, given the port where to listen on.
func Start(listenPort int) {
	visits = 0
	rand.Seed(time.Now().UTC().UnixNano())
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

	err := http.ListenAndServe(":"+strconv.Itoa(listenPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

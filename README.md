Elixir
======

A server which returns random hilarious quotes from movies or TV serials that you like. An example server is hosted on 
andazapnapna.com/quote, which currently only serves quotes from one movie of my choice, 
'<a href="http://en.wikipedia.org/wiki/Andaz_Apna_Apna" target="_blank">Andaz Apna Apna</a>'. 

But elixir is capable of serving quotes from different sources at different points. What's more, is that you can 
add/delete/modify the quotes files at any point of time. All you need to do is, create a file containing the quotes that 
you want to get served, one quote on one line, and the file should have a '.quotes' extension. If the filename was 
"foo.quotes", it will get served at "/foo". You can also modify the index.html page, and it will be auto-reloaded.

Installing and Using
====================

1. Install ```go```: On ubuntu, this is ```sudo apt-get install golang```.

2. Get the package by doing ```go get github.com/reddragon/elixir```.

3. Get cowsay-go by doing ```go get github.com/dhruvbird/go-cowsay```

4. You can take a look at the example provided in the ```examples``` directory on how to use elixir. 


Running
=======

1. To run the example server, switch to its directory and do: ```go build example.go```

2. Run the server by: ```./example -port=8080```, or whatever port you like. elixir logs the diagnostic messages to 
```stderr``` and visits to ```stdout```. So, feel free to direct those two streams to any log file that you feel 
appropriate.

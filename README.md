Elixir
======

A server which returns random hilarious quotes from movies or TV serials that you like. An example server is hosted on andazapnapna.com/quote, which currently only serves quotes from one movie of my choice, 'Andaz Apna Apna'. 

But elixir is capable of serving quotes from different sources at different points. What's more, is that you can add/delete/modify the quotes files at any point of time. All you need to do is, create a file containing the quotes that you want to get served, one quote on one line, and the file should have a '.quotes' extension. If the filename was "foo.quotes", it will get served at "/foo". 

Installing
==========

1. Install ```go```: On ubuntu, this is ```sudo apt-get install golang```

2. Get cowsay-go by doing ```go get github.com/dhruvbird/go-cowsay```

3. Build the package: ```go build quotes.go```

4. Build the example server: ```go build example.go```

5. Run the server by: ```./example -port=8080```

Running
=======

1. ./quotes -port=80

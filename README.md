Elixir
======

Elixir is a simple to use, fortune-cookie server which returns a random quote from a quote database. 

Some features:
* Capable of serving quotes from multiple quote databases.
* A quote database is simply a '.quotes' file (say, ```foo.quotes```), which is placed in the directory where you launched the server, and has one quote per line. The quotes from the ```foo.quotes``` database would be served at ```/foo``` endpoint, and similarily the quotes from 'bar.quotes' would be served at ```/bar```. Look at the ```examples``` dir for an example.
* Can detect when you add / delete / make changes to any of the quotes files. No need to restart the server.
* You can also modify the index.html page, which is served at "/".
* You can also serve the quotes in [cowsay](http://en.wikipedia.org/wiki/Cowsay) format, wherein, the quotes would be returned with ASCII art of a cow trying to say them. For example,

```
 _________________________________________
/ Walter White: If you don’t know who I   \
| am, maybe your best course would be to  |
\ tread lightly.                          /
 -----------------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
```

An example server is hosted at [randquotes.com](http://randquotes.com), which serves quotes form the movies and tv-serials that I like. 


Installing and Using
====================

1. Install ```go```: On ubuntu, this is ```sudo apt-get install golang```.

2. Get the package by doing ```go get github.com/reddragon/elixir```.

3. You can take a look at the example provided in the ```examples``` directory on how to use elixir. 


Running
=======

1. To run the example server, switch to its directory and build it. For example, if you would like to build the randquotes server, do ```cd examples/randquotes.com``` and ```go build randquotes.go```

2. Run the server by: ```./randquotes -port=8080```, or whatever port you like. elixir logs the diagnostic messages to 
```stderr``` and visits to ```stdout```. So, feel free to direct those two streams to any log file that you feel 
appropriate.

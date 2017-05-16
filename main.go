package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"bitbucket.org/iulianclita/amazonsurfer/crawler"
	"github.com/gorilla/websocket"
)

// Template container
var tpl *template.Template

// Create the web crawler. It will be shared accross all calls to the server
// This is why is mandatory to have only one opened session
var crw = &crawler.Crawler{
	Timeout: 10,
}

// These are constants related to websockets buffer sizes
const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// Upgrader takes care of upgrading a standard HTTP conection to a websockt conection persistent connection
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// Set up some options before starting the program
func init() {
	log.SetOutput(os.Stdout)
	tpl = template.Must(template.ParseFiles(filepath.Join("templates", "index.gohtml")))
}

// We do not have a favicon so send 404 response
func favicon(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// Index loads the home page and fills the parsed template with all the necessary data
func index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Host       string
		Categories map[uint8]string
		Relax      int
	}{
		Host:       r.Host,
		Categories: crawler.GetCategories(),
		Relax:      crawler.GetRelax(),
	}

	tpl.Execute(w, data)
}

// search attaches useful data from the request to the existing web crawler
// The crawler stores that data into its options property
func search(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Fatal("Method should be POST")
	}
	if err := r.ParseForm(); err != nil {
		log.Fatal("Error parsing the form")
	}
	err := crw.MapOptions(r)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Fatal(err)
		return
	}
	io.WriteString(w, "ok")
}

// Here is the core processing where the lookup is made
// We start by upgrading our connection to websockets
// After we launch the crawler in the background to search for products
// We wait for products to be sent in the main goroutine and flush them in the frontend
func start(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Fatal("Socker error:", err)
	}
	// This channel will receive the products from the crawler
	prods := make(chan crawler.Product)
	// Run the crawler in the background
	go crw.Run(conn, prods)
	// Wait for incoming products
	for p := range prods {
		if err := conn.WriteJSON(p); err != nil {
			log.Println("Send error:", err)
		}
	}
}

// stop closses the current websockets connection
func stop(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Fatal("Method should be POST")
	}
	crw.Stop()
}

// Main function starts the program
// We can use a custom port if the default one is already taken
// The default port is 1234 so in order to access out server we must visit http://localhost:1234
func main() {
	port := flag.String("port", "1234", "Port where the server should listen")
	flag.Parse()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/search", search)
	http.HandleFunc("/start", start)
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/", index)
	log.Printf("Listening on port %s\n", *port)
	addr := ":" + *port
	log.Fatal(http.ListenAndServe(addr, nil))
}

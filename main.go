package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"bitbucket.org/achillezz/amazonsurfer/crawler"
	"github.com/gorilla/websocket"
)

var tpl *template.Template

var crw = &crawler.Crawler{
	Prods: make(chan crawler.Product),
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func init() {
	log.SetOutput(os.Stdout)
	tpl = template.Must(template.ParseFiles(filepath.Join("templates", "index.gohtml")))
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Host       string
		Categories map[uint8]string
	}{
		Host:       r.Host,
		Categories: crawler.GetCategories(),
	}

	tpl.Execute(w, data)
}

func search(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Fatal("Method should be POST")
	}
	err := crw.MapOptions(r)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	io.WriteString(w, "ok")
}

func process(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Fatal("Socker error:", err)
	}

	go crw.Run()

	for p := range crw.Prods {
		if err := conn.WriteJSON(p); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	port := flag.String("port", "1234", "Port where the server should listen")
	flag.Parse()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/search", search)
	http.HandleFunc("/process", process)
	http.HandleFunc("/", index)
	log.Printf("Listening on port %s\n", *port)
	addr := ":" + *port
	log.Fatal(http.ListenAndServe(addr, nil))
}

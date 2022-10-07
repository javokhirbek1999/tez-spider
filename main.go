package main

import (
	"html/template"
	"log"
	"net/http"
	"runtime"

	"github.com/javokhirbek1999/tez-spider/handlers"
)

type Context struct {
	Songs []handlers.Song
}

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {

	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", fs)
	// http.Handle("/static/*", http.StripPrefix("/static/", fs))
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		templ, err := template.ParseFiles("./client/index.html")

		if err != nil {
			log.Fatal(err)
		}

		// query := r.Form.Get("search")

		songs := handlers.Crawler("enrique iglesias")

		// fmt.Fprintf(w, fmt.Sprintf("%v", songs))

		templ.Execute(w, Context{Songs: songs})
		// err := templ.Execute(w, songs)

		// if err != nil {
		// log.Fatal(err)
		// }
	})

	http.ListenAndServe(":9090", nil)

}

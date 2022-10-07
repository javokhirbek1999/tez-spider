package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/javokhirbek1999/tez-spider/handlers"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type Response struct {
	Data       []handlers.Song
	Total      int64
	StatusCode int
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query().Get("query")

		songs := handlers.Crawler(query)
		encoder := json.NewEncoder(w)

		totalSongs := len(songs)

		if totalSongs == 0 {
			encoder.Encode(Response{
				Data:       nil,
				Total:      0,
				StatusCode: http.StatusBadRequest,
			})
		} else {
			encoder.Encode(Response{
				Data:       songs,
				Total:      int64(totalSongs),
				StatusCode: http.StatusOK,
			})
		}
	})

	go func() {
		err := http.ListenAndServe(":9090", nil)

		if err != nil {
			log.Fatalf("Server Error: %v", err)
			return
		}
	}()

	sig := make(chan os.Signal)

	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, os.Kill)

	signalMessage := <-sig

	log.Fatalf("Server stopped: %v\n", signalMessage)

}

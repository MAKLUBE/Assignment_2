package api

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	data    = make(map[string]string)
	mu      sync.Mutex
	logChan = make(chan string, 100)
)

var templates = template.Must(template.ParseFiles("index.html", "logs.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		value := r.FormValue("value")

		if id != "" && value != "" {
			data[id] = value
			logChan <- fmt.Sprintf("Post Request done/ addId %v", id)
		}
	}

	templates.ExecuteTemplate(w, "index.html", data)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if r.Method == http.MethodPost {
		go func() {
			time.Sleep(1 * time.Second)
			mu.Lock()
			datalen := len(data)
			data = make(map[string]string)
			logChan <- fmt.Sprintf("DELETE called/ Date deleted %v items", datalen)
			defer mu.Unlock()
		}()
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case logMsg := <-logChan:
		templates.ExecuteTemplate(w, "logs.html", logMsg)
	default:
		templates.ExecuteTemplate(w, "logs.html", "no new logs")
	}

}

func startLogger(ctx context.Context) {
	for {
		select {
		case logMsg := <-logChan:
			log.Println(logMsg)
		case <-ctx.Done():
			log.Println("shutting down logger")
			return
		}
	}
}

package main

import (
	"flag"
	"github.com/yvdinesh/go-blueprints/chat-server/chat"
	"github.com/yvdinesh/go-blueprints/chat-server/trace"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filepath string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(t.filepath))
	})
	t.templ.Execute(w, req)
}

var (
	addr = flag.String("addr", ":8080", "")
)

func main() {
	flag.Parse()
	tracer := trace.NewTracer(os.Stdout)
	r := chat.NewRoom(chat.WithTracer(tracer))
	http.Handle("/", &templateHandler{filepath: "templates/client.html"})
	http.Handle("/room", r)
	go r.Run()
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("ListenAndServer: %v", err)
	}
}

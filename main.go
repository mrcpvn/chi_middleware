package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pressly/chi"

	"golang.org/x/net/context"
)

var (
	initOnce  sync.Once
	chirouter *chi.Mux
)

type ContextHandler interface {
	ServeHTTPC(context.Context, http.ResponseWriter, *http.Request)
}

type TimerMiddleware struct {
	Handler ContextHandler
}

type ChiHandler struct{}

func (ch *ChiHandler) ServeHTTPC(c context.Context, w http.ResponseWriter, r *http.Request) {
	chirouter.ServeHTTPC(c, w, r)
}

func indexRouter() chi.Router {
	r := chi.NewRouter()
	r.Handle("/hello/:name", chi.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		name := chi.URLParam(c, "name")
		w.Write([]byte(fmt.Sprintf("hello %v :)", name)))
		times, _ := c.Value("timing").(map[string]time.Duration)
		times["time"] = time.Now().Sub(start)
		fmt.Printf("%v\n", c)
	}))
	r.Handle("/bye/:name", chi.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		time.Sleep(500 * time.Millisecond)
		name := chi.URLParam(c, "name")
		w.Write([]byte(fmt.Sprintf("bye %v :(", name)))
		times, _ := c.Value("timing").(map[string]time.Duration)
		times["time"] = time.Now().Sub(start)
		fmt.Printf("%v\n", c)
	}))
	return r
}

func main() {
	initOnce.Do(func() {
		chirouter = chi.NewRouter()
		chirouter.Mount("/timed", indexRouter())
	})

	host := os.Args[1]
	mux := http.NewServeMux()

	mux.Handle("/", http.Handler(&TimerMiddleware{
		Handler: &ChiHandler{},
	}))

	log.Printf("Server started at %v; pid = %v\n", fmt.Sprintf("%s", host), os.Getpid())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", host), mux))
}

func (tm *TimerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := context.Background()
	c = context.WithValue(c, "timing", make(map[string]time.Duration))
	fmt.Printf("ctx berfore route: %v\n", c)
	tm.Handler.ServeHTTPC(c, w, r)
	times, _ := c.Value("timing").(map[string]time.Duration)
	w.Write([]byte(fmt.Sprintf("\nroute time = %v", times["time"])))
	fmt.Printf("out ctx: %v\n", c)

}

package inv

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Addrport string // config option
	Pubdir   string // config option

	*mux.Router
	*log.Logger
}

var (
	r *mux.Router
)

/*
   Use the runtime/trace and net/http/pprof packages
*/

// HTTPServer will register routes, open connection and listen for incoming
func StartServer(addrport string, done chan<- bool) {
	if addrport == "" {
		log.Fatalf("Server must have a port to listen with")
	}

	r = mux.NewRouter()
	r.HandleFunc("/crawl/{url}", HandleCrawl)
	r.Handle("/", http.FileServer(http.Dir("./pub")))

	log.Infoln("listening on ", addrport)
	err := http.ListenAndServe(addrport, r)
	if err != nil {
		log.Errorf("Server terminated error %v", err)
	}
	done <- true
}
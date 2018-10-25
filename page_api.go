package moni

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func PageListHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, pages)
}

func PageIdHandler(w http.ResponseWriter, r *http.Request) {

	// Get a couple vars ready for later
	var out interface{}
	var err error

	// Get the url from the request and extract the storage
	// index (name) of the corresponding object.
	url := urlFromRequest(r)
	name := NameFromURL(url)
	out = "No Output"

	if page := pages.Get(url); page != nil {
		switch r.Method {
		case "GET":
			out = page
		case "PUT", "POST":
			log.Infoln("overwriting ", name)
			out = "done"
		case "DELETE":
			delete(pages, name)
			out = "done"
		}
	} else {
		switch r.Method {
		case "GET", "DELETE":
			// Nothing to get or delete
			err = errors.New("object not found " + name)
		case "PUT", "POST":
			st := GetStorage()
			if _, err := st.StoreObject(name, page); err != nil {
				err = fmt.Errorf("page %s error %v", url, err)
			} else {
				out = `{"msg": "done"}`
			}
		}
	}
	if err != nil {
		JSONError(w, err)
	}
	writeJSON(w, out)
}
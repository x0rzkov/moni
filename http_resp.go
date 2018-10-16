package main

import (
	"encoding/json"
	"net/http"

	log "github.com/mgutz/logxi/v1"
)

/*
// Write the response as HTML
func writeHTML(w http.ResponseWriter, title string, body string) {
	if Config.Hasfiles {
		t, err := template.ParseFiles("tmpl/index.html")
		if err != nil {
			HTMLError(w, err.Error())
			return
		}
		p := NewPage(title, body)
		if err = t.Execute(w, p); err != nil {
			HTMLError(w, err.Error())
		}
	}
}
*/
// Write the response as JSON
func writeJSON(w http.ResponseWriter, obj interface{}) {
	var jbytes []byte
	var err error

	// Set response content type to json
	if jbytes, err = json.Marshal(obj); err != nil {
		http.Error(w, err.Error(), 500) // TODO: replace 500 with http.XXX
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jbytes)
	if err != nil {
		JSONError(w, err.Error())
	}
}

func JSONError(w http.ResponseWriter, err string) {
	http.Error(w, err, 500)
	log.Error(err)
}

func HTMLError(w http.ResponseWriter, err string) {
	http.Error(w, err, 500)
	log.Error(err)
}
package isoserv_handlers

import (
	"github.com/rkbalgi/isosim/iso"
	"log"
	"net/http"
	"path/filepath"
)

func AddAll() {

	addIsoServerHandlers()
	addIsoServerSaveDefHandler()
	fetchDefHandler()
	startServerHandler()
	addGetActiveServersHandler()
	stopServerHandler()

}

func addIsoServerHandlers() {

	log.Print("Adding ISO server handler .. ")
	http.HandleFunc("/iso/v0/server", func(rw http.ResponseWriter, req *http.Request) {

		pattern := "/iso/v0/server"
		if iso.DebugEnabled {
			log.Printf("Pattern: %s . Requested URI = %s", pattern, req.RequestURI)
		}

		file := filepath.Join(iso.HtmlDir, "iso_server.html")
		if iso.DebugEnabled {
			log.Print("Serving file = " + file)
		}
		http.ServeFile(rw, req, file)

	})

}

func sendError(rw http.ResponseWriter, errorMsg string) {
	if iso.DebugEnabled {
		log.Print("Sending error = " + errorMsg)
	}
	rw.Header().Set("X-IsoSim-ErrorText", errorMsg)
	rw.WriteHeader(http.StatusBadRequest)
	_, _ = rw.Write([]byte(errorMsg))

}

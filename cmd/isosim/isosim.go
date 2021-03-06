package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"isosim/iso"
	"isosim/iso/server"
	"isosim/web/handlers"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

var version = "v0.6.x"

//v0.1 - Initial version
//v0.2 - ISO server development (08/31/2016)
//v0.5 - Support for embedded/nested fields and logging via sirupsen/logrus
//v0.6 - react front and multiple other changes

func main() {

	isDebugEnabled := flag.Bool("debug-enabled", true, "true if debug logging should be enabled.")
	flag.StringVar(&iso.HTMLDir, "html-dir", "", "Directory that contains any HTML's and js/css files etc.")

	specsDir := flag.String("specs-dir", "", "The directory containing the ISO spec definition files.")
	httpPort := flag.Int("http-port", 8080, "Http port to listen on.")
	dataDir := flag.String("data-dir", "", "Directory to store messages (data sets). This is a required field.")

	flag.Parse()

	if *isDebugEnabled {
		log.SetLevel(log.DebugLevel)
		log.Infoln("debug level logging is enabled.")
	}

	//log.SetFormatter(&log.TextFormatter{ForceColors: true, DisableColors: false})

	if *dataDir == "" || *specsDir == "" || iso.HTMLDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := server.Init(*dataDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	//read all the specs from the spec file
	err = iso.ReadSpecs(*specsDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	//check if all the required HTML files are available
	if err = handlers.Init(iso.HTMLDir); err != nil {
		log.Fatal(err.Error())
	}

	log.Infoln("Starting ISO WebSim ", "Version = "+version)

	tlsEnabled := os.Getenv("TLS_ENABLED")
	if tlsEnabled == "true" {
		certFile := os.Getenv("TLS_CERT_FILE")
		keyFile := os.Getenv("TLS_KEY_FILE")

		log.Infof("Using certificate file - %s, key file: %s", certFile, keyFile)
		log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(*httpPort), certFile, keyFile, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*httpPort), nil))
	}

	log.Infof("ISO WebSim started!")

}

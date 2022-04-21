package main

// modified from
// https://github.com/jmalloc/echo-server
import (
	"bytes"
	"log"
	"net/http"
	"os"
	"time"
)

type ServerTime struct {
	UtcTime time.Time
	Time    time.Time
	TZ      string
	Offset  int
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("proxy server received request")
	// check for request Header and forward it
	reqIdHeaderKey := http.CanonicalHeaderKey("x-request-id")
	originalVal, incomingOk := r.Header[reqIdHeaderKey]
	if incomingOk {
		w.Header().Add(reqIdHeaderKey, originalVal[0])
	}

	defer r.Body.Close()

	log.Printf("--------  %s | %s %s\n", r.RemoteAddr, r.Method, r.URL)

	log.Printf("Headers\n")
	//Iterate over all header fields
	for k, v := range r.Header {
		log.Printf("%q : %q\n", k, v)
	}

	buf := &bytes.Buffer{}
	buf.ReadFrom(r.Body) // nolint:errcheck

	if buf.Len() != 0 {
		bodyStr := buf.String() // nolint:errcheck
		log.Println(bodyStr)
	}

	// sendServerHostnameString := os.Getenv("SEND_SERVER_HOSTNAME")
	// if v := r.Header.Get("X-Send-Server-Hostname"); v != "" {
	// 	sendServerHostnameString = v
	// }

	// Return a 200 Status code
	return
}

func main() {
	// log.SetOutput(os.Stdout)
	log.Println("proxy container config")
	log.Printf("SERVERURL: %s", os.Getenv("SERVERURL"))

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

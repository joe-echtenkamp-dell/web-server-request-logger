package main

import (
	"encoding/json"
	"fmt"
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

type Result struct {
	ServerTime time.Time
	ServerTZ   string

	ClientTime time.Time
	ClientTZ   string

	Pass bool
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("proxy server received request")
	// check for request Header and forward it
	reqIdHeaderKey := http.CanonicalHeaderKey("x-request-id")
	originalVal, incomingOk := r.Header[reqIdHeaderKey]

	serverUrl := os.Getenv("SERVERURL")

	// call the server
	req, err := http.NewRequest("GET", serverUrl, nil)
	if err != nil {
		log.Println("Failed to create new request.")
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if incomingOk {
		req.Header.Set(reqIdHeaderKey, originalVal[0])
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to make request")
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	var serverresult ServerTime
	json.NewDecoder(resp.Body).Decode(&serverresult)

	dt := time.Now()
	zone_name, _ := dt.Zone()

	res := Result{
		ServerTime: serverresult.Time,
		ServerTZ:   serverresult.TZ,
		ClientTime: dt,
		ClientTZ:   zone_name,
		Pass:       true,
	}

	// check that the server returned the correct request id
	returnVal, proxyOk := resp.Header[reqIdHeaderKey]

	if proxyOk {
		w.Header().Add(reqIdHeaderKey, returnVal[0])
	}

	if incomingOk && proxyOk {
		// both are set but are they the same?
		if returnVal[0] != originalVal[0] {
			log.Println("Server returned different X-Request-Id then the request")
			log.Println(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
	} else if incomingOk {
		if !proxyOk {
			// if we have an incoming, but not a proxy, then we have a problem
			log.Println("Server lost request id")
			log.Println(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
	}

	str, _ := json.Marshal(res)
	fmt.Fprint(w, string(str))
}

func main() {
	// log.SetOutput(os.Stdout)
	log.Println("proxy container config")
	log.Printf("SERVERURL: %s", os.Getenv("SERVERURL"))

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

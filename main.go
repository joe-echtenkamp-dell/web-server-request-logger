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
	// check for request Header and forward it
	reqIdHeaderKey := http.CanonicalHeaderKey("x-request-id")
	originalVal, ok := r.Header[reqIdHeaderKey]

	serverUrl := os.Getenv("SERVERURL")

	// call the server
	req, err := http.NewRequest("GET", serverUrl, nil)
	if err != nil {
		log.Fatal("Failed to create new request.")
		log.Fatal(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if ok {
		req.Header.Set(reqIdHeaderKey, originalVal[0])
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to make request")
		log.Fatal(err.Error())
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
	returnVal, ok := resp.Header[reqIdHeaderKey]

	if ok {
		w.Header().Add(reqIdHeaderKey, returnVal[0])
	}

	if returnVal[0] != originalVal[0] {
		log.Fatal("Server returned different X-Request-Id then the request")
		log.Fatal(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	str, _ := json.Marshal(res)
	fmt.Fprintf(w, string(str))
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

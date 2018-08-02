package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

// server is the address to send "Up" messages too
var server string

// port is the TCP port that either ther server will listen to, or be connected to on.
var port int

// servicename is the hostname of the service that should loop back to myself
var servicename string

// when set to TRUE, enables starting as the listener
var listener bool

//
var stopwatch time.Time

func main() {

	log.Infoln("AM I UP?")

	server := flag.String("server", "", "Server address that is listening for an UP message")
	port := flag.Int("port", 8080, "Server port that will be used")
	listener := flag.Bool("listen", false, "Listen for UP messages")

	flag.Parse()
	// build URL from server and port
	url := fmt.Sprintf("%s:%d", *server, *port)
	if *listener == false {
		if *server == "" {
			log.Fatalln("No Server was specified, unable to send UP messages")
		}
		if servicename == "" {
			log.Infoln("No Service name specified, so sending a single UP message")
			oneShot(fmt.Sprintf("%s/up", url))
		} else {
			url := fmt.Sprintf("%s:%d", servicename, *port)
			log.Infof("Attempting to connect to %s ", url)

			// spin out the listener to a go routine, so that the service has an
			// endpoint that will be reached by the service VIP
			go func() {
				http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
				http.HandleFunc("/service", service)
			}()
		}

		log.Warnln("Now sleeping for an hour ZZzzzz.")
		time.Sleep(time.Hour)
	} else {
		log.Infof("Starting listener on port [%d]", *port)
		log.Infof("Stop watch started at %v, restart it by using the /stopwatch endpoint", time.Now())
		stopwatch = time.Now()

		http.HandleFunc("/", root)
		http.HandleFunc("/stopwatch", startStopwatch)
		http.HandleFunc("/up", upPost)
		http.HandleFunc("/serviceup", servicePost)
		http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	}
}

// oneShot will attempt "once" to POST back to the server
func oneShot(url string) {
	log.Infof("Sending UP message to listener [%s]", url)
	hostnameString, err := os.Hostname()
	if err != nil {
		log.Fatalln("Error retrieveing contaier hostname")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(hostnameString)))
	if err != nil {
		log.Fatalln("Error creating new HTTP request")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error during POST to server: %s", url)
	}

	if resp.StatusCode != 200 {
		log.Fatalln("Error speaking with server, ensure the correct server is listening")
	}
	log.Infof("Succesfully informed [%s] that I AM UP", url)
}

// serviceLoop will continually hit the service endpoint until the VIP responds, it will then
// hit the listner /service endpoint
func serviceLoop(serviceurl, listenerurl string) {
	log.Infof("Beginning attempt to connect to service [%s]", serviceurl)

	for {
		req, err := http.NewRequest("GET", serviceurl+"/service", nil)
		if err != nil {
			log.Fatalln("Error creating new HTTP request")
		}
		client := &http.Client{}
		client.Timeout = time.Second // this might be too much, too little
		resp, err := client.Do(req)
		if err != nil {
			log.Infoln("Service endpoint doesn't appear available")
		} else {

			if resp.StatusCode != 200 {
				log.Fatalln("Error speaking with server, ensure the correct server is listening")
			} else {
				break
			}
		}
	}
	log.Infof("Succesfully connected to [%s]", serviceurl)
	oneShot(fmt.Sprintf("%s/serviceurl", listenerurl))
}

func root(w http.ResponseWriter, r *http.Request) {
	log.Printf("HTTP Request to / from  %s", r.RemoteAddr)
	io.WriteString(w, "<h1>Hello African Bank!</h1>")
}

func service(w http.ResponseWriter, r *http.Request) {
	log.Printf("Service responding to  %s", r.RemoteAddr)
	io.WriteString(w, "")
}

// listner - endpoints

// Captures a container being up
func startStopwatch(w http.ResponseWriter, r *http.Request) {
	log.Printf("Stopwatch start request from %s", r.RemoteAddr)
	stopwatch = time.Now()
	io.WriteString(w, "")

}

// Captures a container being up
func upPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("UP Message from %s after %v", r.RemoteAddr, time.Since(stopwatch))
	io.WriteString(w, "")
}

// Captures the service responding
func servicePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("Service Message from %s after %v", r.RemoteAddr, time.Since(stopwatch))
	io.WriteString(w, "")
}

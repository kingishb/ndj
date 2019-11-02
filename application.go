package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type result struct {
	line []byte // some arbitrary json
	err  error
}

func streamDecode(url string, c chan result) {
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		c <- result{[]byte{}, err}
		close(c)
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		c <- result{[]byte{}, fmt.Errorf("http status code at %s non-200", url)}
		close(c)
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	// read open bracket
	_, err = dec.Token()
	if err != nil {
		c <- result{[]byte{}, err}
		close(c)
		return

	}

	// while the array contains values, json decode and emit to the
	// channel
	for dec.More() {
		var m map[string]interface{}
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			c <- result{[]byte{}, err}
			close(c)
			return

		}
		j, err := json.Marshal(m)
		if err != nil {

			c <- result{[]byte{}, err}
			close(c)
			return
		}
		c <- result{j, nil}

	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		c <- result{[]byte{}, err}
		close(c)
		return

	}

	// close the channel when finished
	close(c)

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	f, _ := os.Create("/var/log/golang/golang-server.log")
	defer f.Close()
	log.SetOutput(f)
	const indexPage = "public/index.html"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if buf, err := ioutil.ReadAll(r.Body); err == nil {
				log.Printf("Received message: %s\n", string(buf))

			}

		} else {
			log.Printf("Serving %s to %s...\n", indexPage, r.RemoteAddr)
			http.ServeFile(w, r, indexPage)

		}

	})

	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")

		}
		w.Header().Set("content-type", "application/json")
		results := make(chan result)
		url := r.URL.Query().Get("url")
		go streamDecode(url, results)
		for x := range results {
			if x.err != nil {
				http.Error(w, "err: "+x.err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "%s\n", x.line)
			flusher.Flush()
		}

	})

	log.Printf("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}

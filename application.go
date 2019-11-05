package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type result struct {
	line []byte
	err  error
}

type config struct {
	head   int
	sample int
	nth    int
	filter string
}

func parseConfig(u *url.URL) config {
	head, err := strconv.Atoi(u.Query().Get("head"))
	if err != nil {
		head = 0
	}
	sample, err := strconv.Atoi(u.Query().Get("sample"))
	if err != nil {
		sample = 0
	}
	nth, err := strconv.Atoi(u.Query().Get("nth"))
	if err != nil {
		nth = 0
	}
	filter := u.Query().Get("filter")

	return config{
		head:   head,
		sample: sample,
		nth:    nth,
		filter: filter,
	}

}

func streamDecode(u string, conf config, c chan result) {
	resp, err := http.Get(u)
	if err != nil {
		c <- result{[]byte{}, err}
		close(c)
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		c <- result{[]byte{}, fmt.Errorf("http status code at %s non-200", u)}
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

	parse := func(dec *json.Decoder, filter string, send bool) {

		var m map[string]interface{}
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			c <- result{[]byte{}, err}
			close(c)
			return

		}
		if !send {
			return
		}
		j, err := json.Marshal(m)
		if err != nil {
			c <- result{[]byte{}, err}
			close(c)
			return
		}
		if filter == "" {
			c <- result{j, nil}
		} else {
			if bytes.Contains(bytes.ToLower(j), bytes.ToLower([]byte(filter))) {
				c <- result{j, nil}
			}
			return
		}
	}

	// while the array contains values, json decode and emit to the
	// channel
	iter := 0
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for dec.More() {
		// head the first n values
		if conf.head > 0 {
			if iter == conf.head {
				break
			} else {
				parse(dec, "", true)
				iter++
				continue

			}
			// sample at n %
		} else if conf.sample > 0 {
			flip := r1.Intn(100)
			if conf.sample > flip {
				parse(dec, "", true)
				iter++
				continue
			} else {
				// parse the value but drop the values
				parse(dec, "", false)
				iter++
				continue
			}
			// return the nth record
		} else if conf.nth > 0 {
			if iter == conf.nth {
				parse(dec, "", true)
				iter++
				break
			} else {
				parse(dec, "", false)
				iter++
				continue
			}
			// filter values with substring s
		} else if conf.filter != "" {
			parse(dec, conf.filter, true)
			iter++

		} else {
			parse(dec, "", true)
			iter++
		}

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

	const indexPage = "public/index.html"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		u := r.URL.Query().Get("url")
		if u == "" {
			http.ServeFile(w, r, indexPage)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		results := make(chan result)
		conf := parseConfig(r.URL)
		go streamDecode(u, conf, results)
		w.Header().Set("content-type", "application/json")

		// begin json array
		fmt.Fprint(w, "[")
		iter := 0
		for x := range results {
			// add comma and newlines after first record forward
			if iter != 0 {
				fmt.Fprint(w, ",\n")
			}
			// handle errors
			if x.err != nil {
				http.Error(w, "\nerr: "+x.err.Error(), http.StatusBadRequest)
				return
			}
			// print line to http flusher
			fmt.Fprintf(w, "%s", x.line)
			flusher.Flush()
			iter++
		}

		// end json array
		fmt.Fprintln(w, "]")

	})

	log.Printf("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}

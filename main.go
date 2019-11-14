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

	"github.com/savaki/jq"
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
	jq     string
}

func validateConfig(c config) error {
	// initialize all the flags as zero
	f := 0
	head := 0
	sample := 0
	nth := 0
	// set them to 1 if they're set
	if c.filter != "" {
		f = 1
	}
	if c.head > 0 {
		head = 1
	}
	if c.sample > 0 {
		sample = 1
	}
	if c.nth > 0 {
		nth = 1
	}

	// add them up -- anything greater than 1 has too many config flags
	if head+sample+nth+f > 1 {
		return fmt.Errorf("error: too many configuration values set in %v", c)
	}
	return nil

}

func parseConfig(u *url.URL) (config, error) {
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
	op := u.Query().Get("jq")
	if op == "" {
		op = "."
	}

	conf := config{
		head:   head,
		sample: sample,
		nth:    nth,
		filter: filter,
		jq:     op,
	}

	log.Println(conf)

	err = validateConfig(conf)
	if err != nil {
		return config{}, err
	}
	return conf, nil

}

func jqfilter(b []byte, s string) ([]byte, error) {

	op, err := jq.Parse(s)
	if err != nil {
		return []byte{}, err
	}
	value, err := op.Apply(b)
	if err != nil {
		return []byte{}, err
	}
	return value, nil
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

	parse := func(dec *json.Decoder, filter string, op string, send bool) {

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

		p, err := jqfilter(j, op)
		if err != nil {
			c <- result{[]byte{}, err}
			close(c)
			return
		}

		if filter == "" {
			c <- result{p, nil}
		} else {
			if bytes.Contains(bytes.ToLower(j), bytes.ToLower([]byte(filter))) {
				c <- result{p, nil}
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
				parse(dec, "", conf.jq, true)
				iter++
				continue

			}
			// sample at n %
		} else if conf.sample > 0 {
			flip := r1.Intn(100)
			if conf.sample > flip {
				parse(dec, "", conf.jq, true)
				iter++
				continue
			} else {
				// parse the value but drop the values
				parse(dec, "", conf.jq, false)
				iter++
				continue
			}
			// return the nth record
		} else if conf.nth > 0 {
			if iter == conf.nth {
				parse(dec, "", conf.jq, true)
				iter++
				break
			} else {
				parse(dec, "", conf.jq, false)
				iter++
				continue
			}
			// filter values with substring s
		} else if conf.filter != "" {
			parse(dec, conf.filter, conf.jq, true)
			iter++

		} else {
			parse(dec, "", conf.jq, true)
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
		port = "8080"
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
			return
		}
		results := make(chan result)
		conf, err := parseConfig(r.URL)
		if err != nil {
			http.Error(w, "err: too many parameters", http.StatusBadRequest)
			return

		}
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

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	echo(w, r)
	return
}

func formatRequest(r *http.Request) (s string) {
	path := r.URL.String()
	method := r.Method
	// status
	s = fmt.Sprintf("%s %s HTTP/1.1\r\n", method, path)
	// headers
	for k, vs := range r.Header {
		v := strings.Join(vs, ", ")
		s += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	s += "\r\n"
	// body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s += "[Read body error!]"
	} else {
		s += string(body)
	}
	return
}

func echo(w http.ResponseWriter, r *http.Request) {
	s := formatRequest(r)
	s += "\n"
	fmt.Print(s)
	fmt.Fprintf(w, s)
	return
}

func main() {
	mux := &MyMux{}
	args := os.Args
	listen := ":9000"
	if len(args) > 1 {
		listen = args[1]
		if !strings.Contains(listen, ":") {
			listen = ":" + listen
		}
	}

	fmt.Printf("Serving echo HTTPServer on %s ...\n\n", listen)
	err := http.ListenAndServe(listen, mux)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/bmizerany/pat"
)

func homeView(w http.ResponseWriter, r *http.Request) {
	// http.Error(w, http.StatusText(404), 404)
	fmt.Fprint(w, "Hello World!")
}

func pingHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func runcmd(script string, arg ...string) (message string) {
	cmd := exec.Command(script, arg...)
	fmt.Println(cmd.Args)

	out := new(bytes.Buffer)
	cmd.Stdout = out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		message = "error"
	}
	message = out.String()
	fmt.Println(message)
	return message
}

func commitHandle(repo string, w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "commit")
	switch repo {
	case "web":
		go func() {
			script, _ := config.Get("hooks").Get("web").Get("cmd").StringArray()
			runcmd("/bin/bash", script...)
		}()
	case "frontend":
		go func() {
			script, _ := config.Get("hooks").Get("frontend").Get("cmd").StringArray()
			runcmd("/bin/bash", script...)
		}()
	}
}

func hooksView(w http.ResponseWriter, r *http.Request) {
	// http.Error(w, http.StatusText(404), 404)
	repo := r.URL.Query().Get(":repo")

	fmt.Printf("Request headers: %s\n", r.Header)
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("Request body: %s\n", body)

	switch event := r.Header.Get("X-Github-Event"); event {
	case "ping":
		pingHandle(w, r)
	case "commit_comment":
	case "push":
		commitHandle(repo, w, r)
		fmt.Fprintf(w, "Received %s %s event", repo, event)
	default:
		fmt.Fprint(w, "Hello!")
	}
}

func logView(w http.ResponseWriter, r *http.Request) {
	path := config.Get("log").MustString()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Fprintf(w, string(content))
	}
}

type viewFunc func(http.ResponseWriter, *http.Request)

func basicAuth(f viewFunc, user, passwd []byte) viewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		basicAuthPrefix := "Basic "

		// 获取 request header
		auth := r.Header.Get("Authorization")
		// 如果是 http basic auth
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// 解码认证信息
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(basicAuthPrefix):],
			)
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 && bytes.Equal(pair[0], user) &&
					bytes.Equal(pair[1], passwd) {
					// 执行被装饰的函数
					f(w, r)
					return
				}
			}
		}

		// 认证失败，提示 401 Unauthorized
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		// 401 状态码
		w.WriteHeader(http.StatusUnauthorized)
	}
}

var config = simplejson.New()

func main() {
	path := os.Getenv("CONFIG")
	if path == "" {
		path = "config.json"
	}
	content, err := ioutil.ReadFile(path)
	config, _ = simplejson.NewJson(content)
	listen := config.Get("listen").MustString()

	user := []byte(os.Getenv("USER"))
	pass := []byte(os.Getenv("PASSWD"))

	router := pat.New()
	router.Get("/", http.HandlerFunc(homeView))
	router.Get("/logs", http.HandlerFunc(basicAuth(logView, user, pass)))
	router.Post("/hooks/:repo", http.HandlerFunc(hooksView))

	http.Handle("/", router)
	log.Printf("listen %s\n", listen)
	err = http.ListenAndServe(listen, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

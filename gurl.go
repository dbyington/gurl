package main

import (
	"bytes"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	FollowRedirectSecureEnv = "GURL_FOLLOW_REDIRECT_SECURE"
	FollowRedirectEnv       = "GURL_FOLLOW_REDIRECT"
	//AWSkey                  = "AWS_ACCESS_KEY"
	//AWSsecretKey            = "AWS_SECRET_KEY"
)

var logit = log.New(os.Stderr, "gurl: ", log.LstdFlags|log.Lshortfile)

type CommandFlags struct {
	User                  string   `short:"u" long:"user" description:"HTTP Basic Auth username"`
	Header                []string `short:"H" long:"header" description:"Additional request headers"`
	Verbose               []bool   `short:"v" long:"verbose" description:"Print verbose logging, incremental"`
	Debug                 bool     `short:"d" long:"debug" description:"Print debug logging"`
	Method                string   `short:"M" long:"method" default:"GET" description:"The HTTP(S) method verb"`
	FollowRedirects       bool     `short:"f" long:"follow" description:"follow 3XX redirects"`
	FollowSecureRedirects bool     `long:"followSecure" description:"follow 3XX redirects over HTTPS"`
	Body                  string   `short:"b" long:"body" description:"request body"`
	Args                  struct {
		URL  string
		Rest []string
	} `positional-args:"yes" required:"yes"`
}

var cmdLineFlags CommandFlags

func main() {

	var flagsParser = flags.NewParser(&cmdLineFlags, flags.Default)
	flagsParser.Usage = "options"
	if _, err := flagsParser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			flagsParser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
	url, err := urlFromString(cmdLineFlags.Args.URL)
	if err != nil {
		logit.Printf("Bad URL provided; %s: %s\n", cmdLineFlags.Args.URL, err.Error())
		os.Exit(1)
	}

	headers := &http.Header{}
	if len(cmdLineFlags.Header) > 0 {
		for _, header := range cmdLineFlags.Header {
			if strings.ContainsRune(header, ':') {
				kv := strings.Split(header, ":")
				headers.Add(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]))
			}
		}
	}

	method := cmdLineFlags.Method

	var body io.ReadCloser
	if len(cmdLineFlags.Body) > 0 {
		logit.Printf("Body: %v", cmdLineFlags.Body)
		body = ioutil.NopCloser(bytes.NewBuffer([]byte(cmdLineFlags.Body)))
		method = http.MethodPost
	}
	request := &http.Request{
		Method: method,
		Header: *headers,
		URL:    url,
		Body:   body,
	}

	client := &http.Client{CheckRedirect: clientRedirects}
	response, err := client.Do(request)
	if err != nil {
		logit.Printf("Failed to %s; %s: %s\n", method, cmdLineFlags.Args.URL, err.Error())
		os.Exit(1)
	}

	respBody, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		logit.Printf("Failed to read response respBody: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(respBody))
}

func urlFromString(s string) (*url.URL, error) {

	u, err := url.Parse(s)
	if err != nil {
		logit.Printf("Bad URL provided; %s: %s\n", s, err.Error())
		os.Exit(1)
	}

	if u.Scheme == "" {
		u.Scheme = "http"
		u, _ = u.Parse(u.String())
	}

	if u.Host == "" {
		logit.Printf("URL must include a valid host: %s\n", u.Host)
		os.Exit(1)

	} else {
		_, err := net.LookupHost(u.Host)
		if err != nil {
			logit.Printf("Failed: %s\n", err.Error())
			os.Exit(1)
		}
	}
	return u, nil
}

var clientRedirects = func(req *http.Request, via []*http.Request) error {
	var followSecure bool
	var follow bool

	if followSecureENV, ok := os.LookupEnv(FollowRedirectSecureEnv); ok {
		if ok, _ := strconv.ParseBool(followSecureENV); ok {
			followSecure = true
		}
	}
	if cmdLineFlags.FollowSecureRedirects {
		followSecure = true
	}

	if followENV, ok := os.LookupEnv(FollowRedirectEnv); ok {
		if ok, _ := strconv.ParseBool(followENV); ok {
			follow = true
		}
	}
	if cmdLineFlags.FollowRedirects {
		follow = true
	}

	if req.URL.Scheme == "https" {
		if followSecure {
			return nil
		}
	} else if req.URL.Scheme == "http" {
		if follow {
			return nil
		}
	}

	return http.ErrUseLastResponse
}

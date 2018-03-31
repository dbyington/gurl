package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

var mylog = log.New(os.Stderr, "app: ", log.LstdFlags|log.Lshortfile)

func ValidateURL(u *url.URL) (*url.URL, error) {

	if u.Scheme == "" {
		u.Scheme = "http"
		u, _ = u.Parse(u.String())
	}

	if u.Host == "" {
		mylog.Printf("URL must include a valid host: %s\n", u.Host)
		os.Exit(1)

	} else {
		addr, err := net.LookupHost(u.Host)
		if err != nil {
			mylog.Println(u.Host, addr)
			return &url.URL{}, fmt.Errorf("Failed: %s", err.Error())
		}
	}
	return u, nil
}

type CommandFlags struct {
	User    string   `short:"u" long:"user" description:"HTTP Basic Auth username"`
	Header  []string `short:"H" long:"header" description:"Additional request headers"`
	Verbose []bool   `short:"v" long:"verbose" description:"Print verbose logging, incremental"`
	Debug   bool     `short:"d" long:"debug" description:"Print debug logging"`
	Args    struct {
		URL  string
		Rest []string
	} `positional-args:"yes" required:"yes"`
}

func main() {
	var commandFlags CommandFlags
	flagsParser := flags.NewParser(&commandFlags, flags.Default)
	flagsParser.Usage = "options"
	if _, err := flagsParser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			flagsParser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

	if len(commandFlags.Args.Rest) > 0 {
		fmt.Println("Extra args:")
		for _, arg := range commandFlags.Args.Rest {
			fmt.Println(arg)
		}
	}

	url, err := url.Parse(commandFlags.Args.URL)
	if err != nil {
		fmt.Errorf("Bad URL provided; %s: %s\n", commandFlags.Args.URL, err.Error())
		os.Exit(1)
	}

	url, err = ValidateURL(url)
	if err != nil {
		mylog.Printf("URL failed validation: %s", err.Error())
		os.Exit(1)
	}
	response, err := http.Get(url.String())
	if err != nil {
		fmt.Errorf("Failed to GET; %s: %s\n", commandFlags.Args.URL, err.Error())
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		fmt.Errorf("Failed to read response body: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(body))
}

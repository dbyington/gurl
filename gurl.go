package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"net/http"
	"net/url"
)

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

type Usage interface {
	Usage() string
}

func (c CommandFlags) Usage() string {
	return "usage"
}
func main() {
	var commandFlags CommandFlags
	flagsParser := flags.NewParser(&commandFlags, flags.Default)
	flagsParser.Usage = "options"
	flagsParser.Parse()

	if len(commandFlags.Args.Rest) > 0 {
		fmt.Println("Extra args:")
		for _, arg := range commandFlags.Args.Rest {
			fmt.Println(arg)
		}
	}

	url, err := url.Parse(commandFlags.Args.URL)
	if err != nil {
		fmt.Errorf("Bad URL provided; %s: %s\n", commandFlags.Args.URL, err.Error())
	}

	response, err := http.Get(url.String())
	if err != nil {
		fmt.Errorf("Failed to GET; %s: %s\n", commandFlags.Args.URL, err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		fmt.Errorf("Failed to read response body: %s\n", err.Error())
	}

	fmt.Println(string(body))
}

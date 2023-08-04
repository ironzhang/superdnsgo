package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ironzhang/tlog"
	"github.com/ironzhang/tlog/iface"

	"github.com/ironzhang/superdnsgo"
	"github.com/ironzhang/superdnsgo/pkg/model"
)

type options struct {
	Count    int
	Tags     string
	LogLevel string
}

func (p *options) Setup() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: superdns-lookup [OPTIONS] DOMAINS\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		fmt.Fprintf(flag.CommandLine.Output(), `Example: superdns-lookup -tags="X-Lane-Cluster=sim001,X-Base-Cluster=sim000" www.superdns.com`)
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
	}

	flag.IntVar(&p.Count, "n", 1, "loop count")
	flag.StringVar(&p.Tags, "tags", "", "route tags")
	flag.StringVar(&p.LogLevel, "log-level", "fatal", "log level")
	flag.Parse()

	if flag.NArg() <= 0 {
		flag.Usage()
		os.Exit(0)
	}
}

func setLogLevel(s string) {
	gsl, ok := tlog.GetLogger().(iface.GetSetLevel)
	if ok {
		lv, _ := iface.StringToLevel(s)
		gsl.SetLevel(lv)
	}
}

func parseTags(s string) (map[string]string, error) {
	if s == "" {
		return nil, nil
	}

	m := make(map[string]string)
	tags := strings.Split(s, ",")
	for _, tag := range tags {
		keyvalues := strings.Split(tag, "=")
		if len(keyvalues) != 2 {
			return nil, fmt.Errorf("%s is an invalid tag", tag)
		}
		m[keyvalues[0]] = keyvalues[1]
	}
	return m, nil
}

func printError(domain string, err error) {
	fmt.Printf("domain: %s\n", domain)
	fmt.Printf("error: %q\n", err)
	fmt.Printf("\n")
}

func printEndpoint(domain string, cluster string, endpoint model.Endpoint) {
	data, _ := json.Marshal(endpoint)
	fmt.Printf("domain: %s\n", domain)
	fmt.Printf("cluster: %s\n", cluster)
	fmt.Printf("endpoint: %s\n", data)
	fmt.Printf("\n")
}

func main() {
	var opts options
	opts.Setup()
	setLogLevel(opts.LogLevel)

	err := superdnsgo.AutoSetup()
	if err != nil {
		fmt.Printf("superdnsgo auto setup: %v\n", err)
		return
	}

	tags, err := parseTags(opts.Tags)
	if err != nil {
		fmt.Printf("parse tags: %v\n", err)
		return
	}

	for i := 0; i < opts.Count; i++ {
		for _, domain := range flag.Args() {
			endpoint, cluster, err := superdnsgo.LookupEndpoint(context.Background(), domain, tags)
			if err != nil {
				printError(domain, err)
			} else {
				printEndpoint(domain, cluster, endpoint)
			}
		}
	}
}

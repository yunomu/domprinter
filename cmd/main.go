package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/yunomu/domprinter"
)

var (
	url     = flag.String("url", "", "URL")
	outfile = flag.String("o", "output", "Outfile")
)

func init() {
	flag.Parse()
	log.SetOutput(os.Stderr)
}

func main() {
	if *url == "" {
		log.Fatalf("url is required")
	}

	f, err := os.OpenFile(*outfile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("OpenFile: %v", err)
	}
	defer f.Close()

	ctx := context.Background()

	p := domprinter.New()
	if err := p.PrintWithError(ctx, *url, f, os.Stderr); err != nil {
		log.Fatalf("dump: %v", err)
	}
}

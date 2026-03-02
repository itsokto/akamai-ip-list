package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var targets = []struct {
	Name string
	ASNs []string
}{
	{"akamai", []string{"AS-AKAMAI"}},
	{"alibaba", []string{"AS37963", "AS45102", "AS24429"}},
	{"cognosphere", []string{"AS135377"}},
}

func main() {
	outputDir := flag.String("output", "output", "output directory")
	noV4 := flag.Bool("no-v4", false, "skip IPv4")
	noV6 := flag.Bool("no-v6", false, "skip IPv6")
	noAggregate := flag.Bool("no-aggregate", false, "disable prefix aggregation")
	sources := flag.String("S", "", "IRR sources (passed to bgpq4 -S)")
	host := flag.String("h", "", "IRR server (passed to bgpq4 -h)")
	flag.Parse()

	log.SetFlags(0)

	var extra []string
	if *sources != "" {
		extra = append(extra, "-S", *sources)
	}
	if *host != "" {
		extra = append(extra, "-h", *host)
	}
	if !*noAggregate {
		extra = append(extra, "-A")
	}

	srsDir := filepath.Join(*outputDir, "srs")
	plainDir := filepath.Join(*outputDir, "plain")
	if err := os.MkdirAll(srsDir, 0755); err != nil {
		log.Fatalf("create srs dir: %v", err)
	}
	if err := os.MkdirAll(plainDir, 0755); err != nil {
		log.Fatalf("create plain dir: %v", err)
	}

	for _, t := range targets {
		fmt.Fprintf(os.Stderr, "\n=== %s (%s) ===\n", t.Name, strings.Join(t.ASNs, " "))

		prefixes, err := queryPrefixes(t.ASNs, extra, *noV4, *noV6)
		if err != nil {
			log.Fatalf("%s: %v", t.Name, err)
		}

		plainPath := filepath.Join(plainDir, t.Name+".txt")
		if err := writePlain(plainPath, prefixes); err != nil {
			log.Fatalf("%s: write plain: %v", t.Name, err)
		}

		srsPath := filepath.Join(srsDir, t.Name+".srs")
		if err := writeSRS(srsPath, prefixes); err != nil {
			log.Fatalf("%s: write srs: %v", t.Name, err)
		}

		fmt.Fprintf(os.Stderr, "  Wrote %s (%d prefixes)\n", plainPath, len(prefixes))
		fmt.Fprintf(os.Stderr, "  Wrote %s\n", srsPath)
	}
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type bgpq4Result struct {
	Prefixes []struct {
		Prefix string `json:"prefix"`
	} `json:"prefixes"`
}

type ruleSetCompat struct {
	Version int    `json:"version"`
	Rules   []rule `json:"rules"`
}

type rule struct {
	IPCIDR []string `json:"ip_cidr"`
}

func queryBGPQ4(args []string) ([]string, error) {
	cmd := exec.Command("bgpq4", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("bgpq4 %s: %w", strings.Join(args, " "), err)
	}
	var result bgpq4Result
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("parse bgpq4 output: %w", err)
	}
	prefixes := make([]string, len(result.Prefixes))
	for i, p := range result.Prefixes {
		prefixes[i] = p.Prefix
	}
	return prefixes, nil
}

func main() {
	output := flag.String("o", "rule-set.json", "output file path")
	asSet := flag.String("as", "AS-AKAMAI", "AS-SET or ASN to query")
	noV4 := flag.Bool("no-v4", false, "skip IPv4")
	noV6 := flag.Bool("no-v6", false, "skip IPv6")
	noAggregate := flag.Bool("no-aggregate", false, "disable prefix aggregation")
	sources := flag.String("S", "", "IRR sources (passed to bgpq4 -S)")
	host := flag.String("h", "", "IRR server (passed to bgpq4 -h)")
	version := flag.Int("version", 2, "rule-set version")
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

	var all []string

	if !*noV4 {
		args := append([]string{"-j", "-l", "prefixes"}, extra...)
		args = append(args, *asSet)
		fmt.Fprintf(os.Stderr, "Querying IPv4: bgpq4 %s\n", strings.Join(args, " "))
		v4, err := queryBGPQ4(args)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stderr, "  %d IPv4 prefixes\n", len(v4))
		all = append(all, v4...)
	}

	if !*noV6 {
		args := append([]string{"-j", "-6", "-l", "prefixes"}, extra...)
		args = append(args, *asSet)
		fmt.Fprintf(os.Stderr, "Querying IPv6: bgpq4 %s\n", strings.Join(args, " "))
		v6, err := queryBGPQ4(args)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stderr, "  %d IPv6 prefixes\n", len(v6))
		all = append(all, v6...)
	}

	data, err := json.MarshalIndent(ruleSetCompat{
		Version: *version,
		Rules:   []rule{{IPCIDR: all}},
	}, "", "  ")
	if err != nil {
		log.Fatalf("json: %v", err)
	}

	if err := os.WriteFile(*output, data, 0644); err != nil {
		log.Fatalf("write: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Wrote %s (%d prefixes)\n", *output, len(all))
}

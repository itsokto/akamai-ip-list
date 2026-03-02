package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type bgpq4Result struct {
	Prefixes []struct {
		Prefix string `json:"prefix"`
	} `json:"prefixes"`
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

func queryPrefixes(asns []string, extra []string, noV4, noV6 bool) ([]string, error) {
	var all []string

	if !noV4 {
		args := append([]string{"-j", "-l", "prefixes"}, extra...)
		args = append(args, asns...)
		fmt.Fprintf(os.Stderr, "Querying IPv4: bgpq4 %s\n", strings.Join(args, " "))
		v4, err := queryBGPQ4(args)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "  %d IPv4 prefixes\n", len(v4))
		all = append(all, v4...)
	}

	if !noV6 {
		args := append([]string{"-j", "-6", "-l", "prefixes"}, extra...)
		args = append(args, asns...)
		fmt.Fprintf(os.Stderr, "Querying IPv6: bgpq4 %s\n", strings.Join(args, " "))
		v6, err := queryBGPQ4(args)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "  %d IPv6 prefixes\n", len(v6))
		all = append(all, v6...)
	}

	return all, nil
}

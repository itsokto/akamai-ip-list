package main

import (
	"os"
	"strings"

	"github.com/sagernet/sing-box/common/srs"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func writePlain(path string, prefixes []string) error {
	return os.WriteFile(path, []byte(strings.Join(prefixes, "\n")+"\n"), 0644)
}

const srsVersion = 2

func writeSRS(path string, prefixes []string) error {
	var headlessRule option.DefaultHeadlessRule
	headlessRule.IPCIDR = prefixes

	plainRuleSet := option.PlainRuleSet{
		Rules: []option.HeadlessRule{
			{
				Type:           C.RuleTypeDefault,
				DefaultOptions: headlessRule,
			},
		},
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return srs.Write(f, plainRuleSet, srsVersion)
}

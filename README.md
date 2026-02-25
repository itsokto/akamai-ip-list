Generate [sing-box](https://sing-box.sagernet.org/) rule-set files from IRR data using [bgpq4](https://github.com/bgp/bgpq4).

## Local usage

```bash
# requires: bgpq4
go run . -o akamai-ip.json

# compile to binary (requires: sing-box)
sing-box rule-set compile --output akamai-ip.srs akamai-ip.json
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-o` | `rule-set.json` | Output file path |
| `-as` | `AS-AKAMAI` | AS-SET or ASN to query |
| `-S` | | IRR sources (passed to bgpq4 `-S`) |
| `-h` | | IRR server (passed to bgpq4 `-h`) |
| `-no-v4` | `false` | Skip IPv4 |
| `-no-v6` | `false` | Skip IPv6 |
| `-no-aggregate` | `false` | Disable bgpq4 prefix aggregation (`-A`) |
| `-version` | `2` | Rule-set version |

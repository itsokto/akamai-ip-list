# cloud-geoip

IP prefix lists for cloud/CDN providers, generated from IRR data using [bgpq4](https://github.com/bgp/bgpq4).

## Supported providers

- Akamai
- Alibaba Cloud

## Downloads

Prebuilt files are published to the [`release`](../../tree/release) branch and updated daily.

| Format | Path |
|--------|------|
| Plain text | `plain/<name>.txt` |
| [sing-box](https://sing-box.sagernet.org/) SRS | `srs/<name>.srs` |

### sing-box usage

```json
{
  "rule_set": [
    {
      "tag": "akamai",
      "type": "remote",
      "format": "binary",
      "url": "https://raw.githubusercontent.com/<owner>/cloud-geoip/release/srs/akamai.srs"
    }
  ]
}
```

## Local usage

```bash
# requires: bgpq4
go run . -output out
```

This generates `out/plain/*.txt` and `out/srs/*.srs`.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-output` | `output` | Output directory |
| `-S` | | IRR sources (passed to bgpq4 `-S`) |
| `-h` | | IRR server (passed to bgpq4 `-h`) |
| `-no-v4` | `false` | Skip IPv4 |
| `-no-v6` | `false` | Skip IPv6 |
| `-no-aggregate` | `false` | Disable bgpq4 prefix aggregation (`-A`) |

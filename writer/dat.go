package writer

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/xtls/xray-core/app/router"
	"google.golang.org/protobuf/proto"
)

var _ Writer = (*DatWriter)(nil)

type DatWriter struct {
	Filename string
}

func (w *DatWriter) Write(baseDir string, entries []Entry) error {
	geoips := make([]*router.GeoIP, 0, len(entries))
	for _, e := range entries {
		cidrs, err := parseCIDRs(e.Prefixes)
		if err != nil {
			return fmt.Errorf("%s: %w", e.Name, err)
		}
		geoips = append(geoips, &router.GeoIP{
			CountryCode: strings.ToUpper(e.Name),
			Cidr:        cidrs,
		})
	}

	data, err := proto.Marshal(&router.GeoIPList{Entry: geoips})
	if err != nil {
		return err
	}

	name := w.Filename
	if name == "" {
		name = "geoip.dat"
	}
	return os.WriteFile(filepath.Join(baseDir, name), data, 0644)
}

func parseCIDRs(prefixes []string) ([]*router.CIDR, error) {
	cidrs := make([]*router.CIDR, 0, len(prefixes))
	for _, prefix := range prefixes {
		_, ipNet, err := net.ParseCIDR(prefix)
		if err != nil {
			return nil, fmt.Errorf("parse CIDR %q: %w", prefix, err)
		}
		ones, _ := ipNet.Mask.Size()

		ipBytes := ipNet.IP
		if v4 := ipBytes.To4(); v4 != nil {
			ipBytes = v4
		}

		cidrs = append(cidrs, &router.CIDR{
			Ip:     ipBytes,
			Prefix: uint32(ones),
		})
	}
	return cidrs, nil
}

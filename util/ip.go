package util

import (
	"crypto/sha256"
	"fmt"
	"net"
	"sort"
	"strings"

	"stellar.af/netbox-to-nfa/cidr"
	"stellar.af/netbox-to-nfa/types"
)

// GetIPRange finds the start and end IP address from a CIDR range.
func GetIPRange(subnet string) (net.IP, net.IP) {
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		panic(err)
	}
	first, last := cidr.AddressRange(ipnet)
	return first, last
}

// GetPrefixGroupSha creates a SHA256 hash from a tenant name and its (sorted) prefixes. This is
// used to determine if an NFA filter entry needs to be updated or added.
func GetPrefixGroupSha(pg types.PrefixGroup) string {
	d := []string{pg.Tenant}
	for _, p := range pg.Prefixes {
		d = append(d, p.Prefix)
	}
	sort.Strings(d)
	s := strings.Join(d, ",")
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

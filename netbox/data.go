package netbox

import (
	"fmt"
	"math"

	"github.com/tidwall/gjson"
	"stellar.af/netbox-to-nfa/types"
	"stellar.af/netbox-to-nfa/util"
)

func processPrefixes(results gjson.Result, prefixes types.PrefixMap) *types.PrefixMap {
	results.ForEach(func(_, r gjson.Result) bool {
		t := r.Get("tenant.name").String()
		if _, exists := prefixes[t]; !exists {
			prefixes[t] = []string{}
		}
		return true
	})

	results.ForEach(func(_, r gjson.Result) bool {
		t := r.Get("tenant.name").String()
		p := r.Get("prefix").String()
		prefixes[t] = append(prefixes[t], p)
		return true
	})
	return &prefixes
}

func NFAPrefixes() (prefixes types.PrefixMap) {
	prefixes = make(types.PrefixMap)
	role := util.GetEnv("NETBOX_NFA_ROLE")

	filter := make(map[string]interface{})
	filter["role"] = role
	filter["status"] = "active"

	d, err := NetboxRequest("GET", "/api/ipam/prefixes", filter)
	util.Check("Error getting Netbox prefixes: ", err)

	results := d.Get("results")
	count := d.Get("count").Float()
	runs := int(math.Round(count / 50))
	next := d.Get("next").Value()

	if next != nil {
		for run := 2; run <= runs; run++ {
			u := fmt.Sprintf("%s", next)
			nd, err := NetboxRawRequest("GET", u)
			util.Check("Error getting page %d of Netbox prefixes: ", err, run)
			processPrefixes(nd.Get("results"), prefixes)
			next = nd.Get("next").Value()
		}
	}

	processPrefixes(results, prefixes)

	return prefixes
}

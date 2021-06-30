package main

import (
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
	"stellar.af/netbox-to-nfa/netbox"
	"stellar.af/netbox-to-nfa/nfa"
	"stellar.af/netbox-to-nfa/types"
	"stellar.af/netbox-to-nfa/util"
)

// findPrefixes gets all Netbox prefixes and compares them to NFA filters to determine which
// prefixes should be synchronized with NFA.
func findPrefixes() []types.PrefixGroup {
	filters, err := nfa.GetFilters()
	util.Check("Error getting NFA filters", err)

	allPrefixes := []types.PrefixGroup{}
	nbp := netbox.NFAPrefixes()

	tenants := make([]string, 0, len(nbp))
	for t := range nbp {
		tenants = append(tenants, t)
	}
	sort.Strings(tenants)

	for _, t := range tenants {
		prefix := types.PrefixGroup{Tenant: t, Prefixes: []types.Prefix{}}
		nbprefix := nbp[t]
		for _, i := range nbprefix {
			start, end := util.GetIPRange(i)
			p := types.Prefix{Start: start.String(), End: end.String(), Prefix: i}
			prefix.Prefixes = append(prefix.Prefixes, p)
		}
		prefix.Sha = util.GetPrefixGroupSha(prefix)
		allPrefixes = append(allPrefixes, prefix)
	}

	prefixes := []types.PrefixGroup{}

	for _, pg := range allPrefixes {
		var matchingFilter *nfa.NFAFilter
		matchingFilter = nil
		for _, f := range filters {
			if f.Description == pg.Sha {
				matchingFilter = &f
				break
			}
		}
		if matchingFilter == nil {
			prefixes = append(prefixes, pg)
		}
	}

	return prefixes
}

// SyncPrefixes adds prefixes from Netbox as filters to NFA.
func SyncPrefixes() ([]types.PrefixGroup, error) {
	prefixes := findPrefixes()
	for _, p := range prefixes {
		// Create an NFA filter for this prefix.
		_, err := nfa.NewFilter(p)
		if err != nil {
			log.Error(fmt.Sprintf("Error creating filter for tenant '%s': %s", p.Tenant, err.Error()))
			return []types.PrefixGroup{}, err
		}
		log.Debug(fmt.Sprintf("Synchronized prefixes for tenant '%s'", p.Tenant))
	}
	log.Info(fmt.Sprintf("Synchronized %d tenant prefixes", len(prefixes)))
	return prefixes, nil
}

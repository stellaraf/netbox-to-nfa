package nfa

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"stellar.af/netbox-to-nfa/types"
	"stellar.af/netbox-to-nfa/util"
)

func createName(tenant string) (name string) {
	name = tenant
	if len(tenant) > 45 {
		name = fmt.Sprintf("%s...", tenant[0:41])
	}
	return
}

// GetFilters gets all NFA filters. Because of how stupid the NFA API is, the `parameters` and
// `filters` objects inside each filter are JSON strings, not actual JSON, so they need to be
// individually unmarshalled.
func GetFilters() ([]NFAFilter, error) {
	raw, err := NFARequest("GET", "/api/filters", emptyMap, nil)
	if err != nil {
		return []NFAFilter{}, err
	}
	var filters []NFAFilter
	raw.Get("data").ForEach(func(_, value gjson.Result) bool {
		params := gjson.Parse(value.Get("parameters").String())

		var p NFAParameter
		json.Unmarshal([]byte(params.Raw), &p)

		filterItems := gjson.Parse(params.Get("filters").String())
		var fi NFAFilterItem
		json.Unmarshal([]byte(filterItems.Raw), &fi)

		r := []byte(value.Raw)
		var f NFAFilter
		json.Unmarshal(r, &f)
		f.Parameters = p
		f.Parameters.Filters = fi

		filters = append(filters, f)
		return true
	})
	return filters, nil
}

func GetExcluded() (excluded []string) {
	raw := util.GetEnv("NB2NFA_EXCLUDED_RANGES")
	values := strings.Split(raw, ",")
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
		_, net, err := net.ParseCIDR(values[i])
		util.Check("Error parsing value '%s' from NB2NFA_EXCLUDED_RANGES environment variable to a valid network", err, values[i])
		s, e := util.GetIPRange(net.String())
		excluded = append(excluded, fmt.Sprintf("%s-%s", s, e))
	}
	return excluded
}

func buildRules(pg types.PrefixGroup) string {
	excluded := GetExcluded()

	// Create an AND rule that includes the prefix as the source, but excludes the excluded ranges
	// as destinations.
	srcInc := NFARule{ComparisonOperator: "eq", Key: "src-addr", Value: []string{}}
	srcExc := NFARule{ComparisonOperator: "noteq", Key: "dst-addr", Value: excluded}

	// Create an AND rule that includes the prefix as the destination, but excludes the excluded
	// ranges as sources.
	dstInc := NFARule{ComparisonOperator: "eq", Key: "dst-addr", Value: []string{}}
	dstExc := NFARule{ComparisonOperator: "noteq", Key: "src-addr", Value: excluded}

	for _, p := range pg.Prefixes {
		ipr := fmt.Sprintf("%s-%s", p.Start, p.End)
		srcInc.Value = append(srcInc.Value, ipr)
		dstInc.Value = append(dstInc.Value, ipr)
	}

	srcGroup := NFAFilterItem{Condition: "and", Rules: []NFARule{srcInc, srcExc}}
	dstGroup := NFAFilterItem{Condition: "and", Rules: []NFARule{dstInc, dstExc}}

	rules := NFAFilterGroup{Condition: "or", Rules: []NFAFilterItem{srcGroup, dstGroup}}

	rJson, err := json.Marshal(rules)
	if err != nil {
		panic(err)
	}
	return string(rJson)
}

// buildfilter builds an NFA filter based on a tenant prefix group.
func buildfilter(pg types.PrefixGroup) []byte {
	name := createName(pg.Tenant)
	rules := buildRules(pg)

	rParam := cleanSprintf(`
	{
		"aggregateBy": "octets",
		"aggregateFunction": "sum",
		"aggregateColumn": "octets",
		"asPathPrependRemove": 1,
		"limit": 10,
		"orderby": "octets",
		"pageSize": 300,
		"filters": %v,
		"groupby": ["ts", "ip-version"],
		"order": "descending",
		"rateUnit": "seconds"
	}`, rules)

	pb, err := json.Marshal(rParam)
	if err != nil {
		panic(err)
	}
	param := string(pb)

	newFilter := cleanSprintf(`
	{
		"name": "%s",
		"shared": true,
		"description": "%s",
		"report": "flows",
		"parameters": %v
	}`, name, pg.Sha, param)

	return []byte(newFilter)
}

// DeleteFilter deletes one NFA filter object by ID.
func DeleteFilter(id int) (gjson.Result, error) {
	return NFARequest("DELETE", fmt.Sprintf("/api/filters/%d", id), emptyMap, nil)
}

// PurgeFilters deletes all NFA filters managed by netbox-to-nfa (inferred by the title and
// existence of a description).
func PurgeFilters() int {
	count := 0
	allFilters, err := GetFilters()
	if err != nil {
		return count
	}

	for _, f := range allFilters {
		desc := &f.Description

		if desc != nil && util.IsProbablySHA256(*desc) {
			_, err := NFARequest("DELETE", fmt.Sprintf("/api/filters/%d", f.Id), emptyMap, nil)
			util.Check("Error deleting filter %d (%s)", err, f.Id, f.Name)
		}
		count++
	}
	return count
}

// NewFilter creates a new NFA filter.
func NewFilter(pg types.PrefixGroup) (gjson.Result, error) {
	name := createName(pg.Tenant)
	nf := buildfilter(pg)
	allFilters, err := GetFilters()
	if err != nil {
		return emptyResult, err
	}
	for _, f := range allFilters {
		if name == f.Name {
			_, err := DeleteFilter(f.Id)
			util.Check("Error deleting filter %d (%s)", err, f.Id, f.Name)
		}
	}
	log.Debug(fmt.Sprintf("Constructed filter for tenant '%s':\n%v\n", pg.Tenant, string(nf)))
	return NFARequest("POST", "/api/filters", emptyMap, &nf)
}

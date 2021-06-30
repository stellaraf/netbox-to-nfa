package types

// Prefix represents a single tenant prefix.
type Prefix struct {
	Prefix string
	Start  string
	End    string
}

// PrefixGroup represents a tenant and all its associated prefixes.
type PrefixGroup struct {
	Tenant   string
	Prefixes []Prefix
	Sha      string
}

// QueryParams represents URL query parameters as a map.
type QueryParams map[string]interface{}

// PrefixMap is a mapping of tenant name → array of associated prefixes.
type PrefixMap map[string][]string

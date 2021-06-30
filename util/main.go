package util

import (
	"fmt"
	"net/url"
	"strings"

	"stellar.af/netbox-to-nfa/types"
)

// SplitRemoveEmpty splits a string by a separator and removes any empty strings from the result.
func SplitRemoveEmpty(s string, sep string) (a []string) {
	parts := strings.Split(s, sep)
	for _, p := range parts {
		if p != "" {
			a = append(a, p)
		}
	}
	return a
}

// BuildUrl creates a "guaranteed correct" URL from a base URL and path.
func BuildUrl(base string, path string, query types.QueryParams) (u *url.URL) {
	// Create a URL object from the base URL.
	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}
	// Create an array of all path elements.
	pathParts := SplitRemoveEmpty(path, "/")
	// Ensure the final URL has a trailing slash.
	// pathParts = append(pathParts, "/")
	// Override the URL object's path property with the cleaned path.
	u.Path = strings.Join(pathParts, "/")
	// Add query params.
	q := url.Values{}
	for k, v := range query {
		q.Set(k, fmt.Sprintf("%s", v))
	}
	u.RawQuery = q.Encode()
	return u
}

package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mkideal/cli"
	"stellar.af/netbox-to-nfa/netbox"
	"stellar.af/netbox-to-nfa/nfa"
	"stellar.af/netbox-to-nfa/util"
)

type argT struct {
	cli.Helper
}

type cfgT struct {
	cli.Helper
	ShowPasswords bool `cli:"p,show-passwords" usage:"Show unobscured passwords in output"`
}

var rootCmd = &cli.Command{
	Desc: "Synchronize Netbox Prefixes with Noction NFA",
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		return nil
	},
}

var cfgCmd = &cli.Command{
	Name: "config",
	Desc: "Get validated configuration variables",
	Argv: func() interface{} { return new(cfgT) },
	Fn: func(ctx *cli.Context) error {
		args := ctx.Argv().(*cfgT)
		c := ctx.Color()
		vars := make(map[string]string)
		plainVars := []string{
			"NETBOX_URL",
			"NETBOX_TOKEN",
			"NETBOX_NFA_ROLE",
			"NFA_URL",
			"NFA_USERNAME",
		}
		secureVars := []string{
			"NETBOX_TOKEN",
			"NFA_PASSWORD",
		}
		excluded := nfa.GetExcluded()
		for _, k := range plainVars {
			v := util.GetEnv(k)
			vars[c.Magenta(k)] = c.Green(c.Bold(v))
		}
		for _, k := range secureVars {
			v := util.GetEnv(k)
			if args.ShowPasswords {
				vars[c.Magenta(k)] = c.Red(c.Bold(v))
			} else {
				o := strings.Repeat("*", len(v))
				vars[c.Magenta(k)] = c.Blue(c.Bold(o))
			}
		}
		excludedV := ""
		for i, e := range excluded {
			l := len(excluded) - 1
			pl := len("NB2NFA_EXCLUDED_RANGES: ")
			p := strings.Repeat(" ", pl)
			if i == 0 {
				excludedV += fmt.Sprintf("%s\n", c.Yellow(c.Bold(e)))
			} else if i == l {
				excludedV += fmt.Sprintf("\t%s%s", p, c.Yellow(c.Bold(e)))
			} else {
				excludedV += fmt.Sprintf("\t%s%s\n", p, c.Yellow(c.Bold(e)))
			}
		}
		vars[c.Magenta("NB2NFA_EXCLUDED_RANGES")] = excludedV
		ctx.String("\n")
		for k, v := range vars {
			ctx.String("\t%s: %s\n", k, v)
		}
		ctx.String("\n")
		return nil
	},
}

var purgeCmd = &cli.Command{
	Name: "purge",
	Desc: "Purge all NFA Filters Managed by netbox-to-nfa",
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		count := nfa.PurgeFilters()
		ctx.String("Purged %d NFA filters", count)
		return nil
	},
}

var syncCmd = &cli.Command{
	Name: "sync",
	Desc: "Run synchronization",
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		u, err := SyncPrefixes()
		if err != nil {
			return err
		}
		ctx.String("Synchronized %d tenant prefixes", len(u))
		return nil
	},
}

var prefixesCmd = &cli.Command{
	Name: "prefixes",
	Desc: "List prefixes from NetBox that should be synced to NFA",
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		c := ctx.Color()
		prefixes := netbox.NFAPrefixes()
		var keys []string
		for t := range prefixes {
			keys = append(keys, t)
		}
		sort.Strings(keys)
		for _, t := range keys {
			tp := prefixes[t]
			ctx.String("%s\n", c.Blue(c.Bold(c.Underline(t))))
			for _, p := range tp {
				if strings.Contains(p, ":") {
					ctx.String("  %s\n", c.Green(c.Bold(p)))
				} else {
					ctx.String("  %s\n", c.Red(c.Bold(p)))
				}
			}
		}
		return nil
	},
}

var filtersCmd = &cli.Command{
	Name: "filters",
	Desc: "List all NFA filters",
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
		filters, err := nfa.GetFilters()
		if err != nil {
			return err
		}
		for _, f := range filters {
			ctx.String("%s\n", util.PrettyStruct(f))
		}
		return nil
	},
}

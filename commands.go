package main

import (
	"sort"
	"strings"

	"github.com/mkideal/cli"
	"stellar.af/netbox-to-nfa/netbox"
	"stellar.af/netbox-to-nfa/nfa"
	"stellar.af/netbox-to-nfa/util"
)

type argT struct {
	Help bool `cli:"h,help" usage:"show help"`
}

var help = cli.HelpCommand("display help information")

func (argv *argT) AutoHelp() bool {
	return argv.Help
}

var rootCmd = &cli.Command{
	Desc: "Synchronize Netbox Prefixes with Noction NFA",
	Argv: func() interface{} { return new(argT) },
	Fn: func(ctx *cli.Context) error {
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

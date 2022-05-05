<div align="center">
  <br/>
  <img src="https://res.cloudinary.com/stellaraf/image/upload/v1604277355/stellar-logo-gradient.svg" width=300 />
  <br/>
  <h3>netbox-to-nfa</h3>
  <br/>
  <a href="https://github.com/stellaraf/netbox-to-nfa/actions?query=workflow%3Arelease">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/stellaraf/netbox-to-nfa/release?color=9100fa&style=for-the-badge">
  </a>
  <br/>
  
  Synchronize [NetBox](https://github.com/netbox-community/netbox) prefixes with Noction NFA Query Filters

</div>

## Usage

### Download the latest [release](https://github.com/stellaraf/netbox-to-nfa/releases/latest)

There are multiple builds of the release, for different CPU architectures/platforms:

There are multiple builds of the release, for different CPU architectures/platforms. Download and unpack the release for your platform:

```shell
wget <release url>
tar xvfz <release file> nb2nfa
```

### Run the binary

```console
$ ./nb2nfa --help

nb2nfa 0.2.0
  Synchronize Netbox Prefixes with Noction NFA

Options:

  -h, --help   display help information

Commands:

  help       display help information
  purge      Purge all NFA Filters Managed by netbox-to-nfa
  sync       Run synchronization
  prefixes   List prefixes from NetBox that should be synced to NFA
  filters    List all NFA filters
  config     Get validated configuration variables

```

### Environment Variables

All of the below environment variables are required for netbox-to-nfa to run.

| Name                     | Description                                                                                      |
| :----------------------- | :----------------------------------------------------------------------------------------------- |
| `NETBOX_URL`             | NetBox URL, e.g. `https://netbox.example.com`                                                    |
| `NETBOX_TOKEN`           | NetBox API Token                                                                                 |
| `NETBOX_NFA_ROLE`        | NetBox prefix role. A prefix must be assigned this role for it to be picked up.                  |
| `NFA_URL`                | NFA URL, e.g. `https://nfa.example.com`                                                          |
| `NFA_USERNAME`           | NFA admin username                                                                               |
| `NB2NFA_EXCLUDED_RANGES` | Comma-separated list of prefixes for which traffic to/from will be excluded from each NFA filter |

#### `NB2NFA_EXCLUDED_RANGES` Detail

For example, if you had a prefix from NetBox of `203.0.113.0/24`, and set:

```bash
export NB2NFA_EXCLUDED_RANGES="192.0.2.0/24,2001:db8::/32"
```

The resulting NFA filter logic would look something like:

- Anything **from** `203.0.113.0/24`, unless the destination is `192.0.2.0/24` or `2001:db8::/32`
- Anything **to** `203.0.113.0/24`, unless the source is `192.0.2.0/24` or `2001:db8::/32`

## Creating a New Release

This project uses [GoReleaser](https://goreleaser.com/) to manage releases. After completing code changes and committing them via Git, be sure to tag the release before pushing:

```
git tag <release>
```

Once a new tag is pushed, GoReleaser will automagically create a new build & release.

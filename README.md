[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=jimmystewpot_dns-preload&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=jimmystewpot_dns-preload) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=jimmystewpot_dns-preload&metric=coverage)](https://sonarcloud.io/summary/new_code?id=jimmystewpot_dns-preload) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=jimmystewpot_dns-preload&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=jimmystewpot_dns-preload) [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=jimmystewpot_dns-preload&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=jimmystewpot_dns-preload)
# dns-preload
A simple go based dns cache preloader.

## Why?

With an intermittent satellite connection when the connection is lost the local DNS cache is flushed. When the connection is restablished the high latency on the link causes a slow down for all connections being reestablished. This was an easy way to update the cache asynchronously when the connection is restored which improves usability of the whole system.

It also works well on trains with crappy reception.
## Usage

Add an @reboot line to your users crontab so when the host reboots it will run the tool to pre-populate the dns server
with the entries that you wish to use.

`dns-preload all --config-file==dns-preload.yaml --full --quiet --server=::1`

Can be added to crontab as a user `crontab -e`

`@reboot $HOME/dns-preload all --config-file=dns-preload.yaml --full --quiet --server=::1`

replace $HOME with where you have placed the executable.
### Configuration

An example configuration file can be found at `example-config.yaml` in the root of the repository.

### Building

#### Local System

```make dns-preload```

#### Linux ARM64

```make linux-arm64```

#### Linux ARM32

```make linux-arm32```

#### Linux x86-64

```make linux-x64```



## Help

```
dns-preload --help
Usage: dns-preload <command>

Preload a series of Domain Names into a DNS server from a yaml configuration

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  all      preload all of the following types from the configuration file
  cname    preload only the cname entries from the configuration file
  hosts    preload only the hosts entries from the configuration file, this does an A and AAAA lookup
  mx       preload only the mx entries from the configuration file
  ns       preload only the ns entries from the configuration file
  txt      preload only the txt entries from the configuration file

Run "./dns-preload <command> --help" for more information on a command.
```

all of the commands above have the same flags.

```
dns-preload all --help
Usage: dns-preload all --config-file=STRING

preload all of the following types from the configuration file

Flags:
  -h, --help                  Show context-sensitive help.

      --config-file=STRING    The configuration file to read the domain list to query from
      --server="localhost"    The server to query to seed the domain list into
      --port="53"             The port to query for on the DNS server
      --workers=5             The number of concurrent goroutines used to query the DNS server
      --quiet                 Suppress the preload response output to console
      --full                  For record types that return a Hostname ensure that these are resolved
      --timeout=30s           The timeout for DNS queries to succeed
      --delay=0s              How long to wait until the queries are executed
```
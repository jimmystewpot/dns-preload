# dns-preload
A simple go based dns cache preloader.

## usage

Add an @reboot line to your users crontab so when the host reboots it will run the tool to pre-populate the dns server
with the entries that you wish to use.

## why?

900ms latency on a slow satellite connection was very frustrating, improving the lookup times by pre-fetching all of the NS and other recorrds.
plus some % of needing something to do while on a very slow crappy internet connection in a remote area on a train.



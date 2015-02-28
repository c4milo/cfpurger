# cfpurger
Checks a given URL for changes in its HTML document. If it finds changes, it then
proceeds to purge Cloudflare's cache.

## Usage
```shell
$ ./cfpurger
Usage of ./cfpurger:
  -cftkn="": Cloudflare API token
  -dryrun=false: Simulates a purging without hitting Cloudflare.
  -email="": Cloudflare account email
  -interval=15: The time in seconds to check for changes.
  -url="": The url to watch for changes
  -version=false: Prints version
```

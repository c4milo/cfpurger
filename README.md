# cfpurger
A 2 hour project that checks a given URL for changes in its HTML document to then purge Cloudflare's cache.

## Usage
```shell
$ ./cfpurger
Usage of ./cfpurger:
  -cfsite="": The name of the site to purge in Cloudflare
  -cftkn="": Cloudflare API token
  -dryrun=false: Simulates a purging without hitting Cloudflare.
  -email="": Cloudflare account email
  -interval=15: The time in seconds to check for changes.
  -url="": The url to watch for changes
  -version=false: Prints version
```

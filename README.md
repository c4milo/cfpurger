# cfpurger
Small tool that checks a given URL for changes in its HTML document, if it finds any, it then purges Cloudflare's cache for the given site.

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

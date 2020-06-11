# SW5CacheWarm
SW5CacheWarm is a small application designed to accelerate cache warming for Shopware 5 shop systems.
It is built as an alternative to `bin/console sw:warm:http:cache --product`.  

In my tests, it needed 5 minutes for warming 10000 articles, while `bin/console sw:warm:http:cache --product` needed 9
minutes.

Note that it only uses the URLs from the `s_core_rewrite_urls` table.
## Installation
You need to have the [Go compiler](https://golang.org/doc/install) installed. If you do not have root
access, you can extract the compiler to `/home/youruser/go1.14` and add `/home/youruser/go1.14/go/bin` to your PATH.

After ensuring the compiler works, simply execute 
```
go get github.com/pkramme/sw5cachewarm
```

The binary is now under `/home/youruser/go/bin/`.

You can also compile on another computer and then upload the resulting binary to the desired location.
## Usage
### Options
```text
  -basepath string
        Shop Basepath
  -dbaddr string
        Shopware Database Host
  -dbname string
        Shopware Database Name
  -dbpass string
        Shopware Database Password
  -dbuser string
        Shopware Database User
  -parallel int
        Number of articles to warm at once (default 4)
  -ratelimit
        Reduces the rate when 503 Service Unavailable is returned by the server (default true)
  -subshopid int
        Subshop ID (default 1)
```
### Example
```text
sw5cachewarm -basepath http://192.168.178.86/ -dbaddr localhost:3306 -dbname shopware -dbuser shopware -dbpass somepass
```

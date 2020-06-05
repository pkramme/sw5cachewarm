# SW5CacheWarm
SW5CacheWarm is a small application designed to accelerate cache warming for Shopware 5 shop systems.
It is built as an alternative to `bin/console sw:warm:http:cache --product`.  

Note that it only uses the URLs from the `s_core_rewrite_urls` table.
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
  -queuedepth int
        Queue depth (default 100)
```
### Example
```text
sw5cachewarm -basepath http://192.168.178.86/ -dbaddr localhost:3306 -dbname shopware -dbuser shopware -dbpass somepass
```
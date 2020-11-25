package main

import (
	"database/sql"
	"flag"
	"github.com/go-sql-driver/mysql"
	"github.com/schollz/progressbar/v3"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	dbuser := flag.String("dbuser", "", "Shopware Database User")
	dbpass := flag.String("dbpass", "", "Shopware Database Password")
	dbname := flag.String("dbname", "", "Shopware Database Name")
	dbaddr := flag.String("dbaddr", "", "Shopware Database Host")
	parallel := flag.Int("parallel", 4, "Number of articles to warm at once")
	basepath := flag.String("basepath", "", "Shop Basepath")
	ratelimit := flag.Bool("ratelimit", true, "Reduces the rate when 503 Service Unavailable is returned by the server")
	subshopid := flag.Int("subshopid", 1, "Subshop ID")

	flag.Parse()

	dbconf := mysql.NewConfig()
	dbconf.User = *dbuser
	dbconf.Passwd = *dbpass
	dbconf.DBName = *dbname
	dbconf.Net = "tcp"
	dbconf.Addr = *dbaddr
	db, err := sql.Open("mysql", dbconf.FormatDSN())
	if err != nil {
		log.Fatalln(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(4)

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Successfully connected to database...")
	log.Println("Gathering SEO URLs...")

	queue := make(chan string, 100)

	var wg sync.WaitGroup

	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go func() {
			var delay time.Duration
			for url := range queue {
				resp, err := http.Get(*basepath + url)
				if err != nil {
					log.Println(err)
					continue
				}
				if *ratelimit {
					if resp.StatusCode == http.StatusServiceUnavailable {
						delay += 10 * time.Millisecond
						log.Println("Server returned 503, adding 10ms delay. Delay is now", delay.Round(time.Millisecond).String())
					}
					time.Sleep(delay)
				}

			}
			wg.Done()
		}()
	}

	var urls []string

	rows, err := db.Query("SELECT path FROM s_core_rewrite_urls WHERE main = 1 AND subshopID = ?", *subshopid)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		if err != nil {
			log.Fatalln(err)

		}
		urls = append(urls, path)
	}

	db.Close()

	bar := progressbar.Default(int64(len(urls)))

	for _, u := range urls {
		bar.Add(1)
		queue <- u
	}

	close(queue)
	wg.Wait()
}

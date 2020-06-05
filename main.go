package main

import (
	"database/sql"
	"flag"
	"github.com/go-sql-driver/mysql"
	"github.com/schollz/progressbar/v3"
	"log"
	"net/http"
	"sync"
)

func main() {
	dbuser := flag.String("dbuser", "", "Shopware Database User")
	dbpass := flag.String("dbpass", "", "Shopware Database Password")
	dbname := flag.String("dbname", "", "Shopware Database Name")
	dbaddr := flag.String("dbaddr", "", "Shopware Database Host")
	queuedepth := flag.Int("queuedepth", 100, "Queue depth")
	basepath := flag.String("basepath", "", "Shop Basepath")

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
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Successfully connected to database...")
	log.Println("Gathering SEO URLs...")

	queue := make(chan string, *queuedepth)

	var wg sync.WaitGroup

	for i := 0; i < *queuedepth; i++ {
		wg.Add(1)
		go func() {
			for url := range queue {
				http.Get(*basepath + url)
			}
			wg.Done()
		}()
	}

	var numberofurls int
	err = db.QueryRow("SELECT COUNT(path) FROM s_core_rewrite_urls").Scan(&numberofurls)
	if err != nil {
		log.Println(err)
	}

	bar := progressbar.Default(int64(numberofurls))

	rows, err := db.Query("SELECT path FROM s_core_rewrite_urls")
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
		bar.Add(1)
		queue <- path
	}
	close(queue)
	wg.Wait()
}

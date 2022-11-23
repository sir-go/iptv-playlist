package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"

	"iptv-playlist/internal/store"
)

func main() {
	fCfgPath := flag.String("c", "config.yml", "path to conf file")
	flag.Parse()

	confBytes, err := ioutil.ReadFile(*fCfgPath)
	if err != nil {
		log.Panic(err)
	}
	cfg, err := LoadConfig(confBytes)
	if err != nil {
		log.Panic(err)
	}

	s := store.NewStore(
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.Dbname)

	tM3U := template.Must(template.ParseFiles(cfg.Playlist.Templates.M3U))
	tXSPF := template.Must(template.ParseFiles(cfg.Playlist.Templates.XSPF))
	r := gin.New()

	r.GET("/pl.m3u", handle("audio/x-mpegurl", s, tM3U, cfg.Playlist.UnicastUrl))
	r.GET("/px.pl.m3u", handle("audio/x-mpegurl", s, tM3U, ""))
	r.GET("/pl.xspf", handle("application/xspf+xml", s, tXSPF, cfg.Playlist.UnicastUrl))
	r.GET("/px.pl.xspf", handle("application/xspf+xml", s, tXSPF, ""))

	srv := &http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           r,
		Addr:              fmt.Sprintf(":%d", cfg.Service.Port),
	}

	log.Printf("run server on http://localhost:%d", cfg.Service.Port)
	log.Fatal(srv.ListenAndServe())
}

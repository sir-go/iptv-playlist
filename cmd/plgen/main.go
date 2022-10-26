package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
)

var (
	dbConn *sql.DB
)

func handleErr(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func handleOk(w http.ResponseWriter, data *string) {
	w.Header().Set("Content-Length", fmt.Sprint(len(*data)))
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(*data)); err != nil {
		log.Panic(err)
	}
}

func getPlaylist() (*[]PlaylistRecord, error) {
	//goland:noinspection ALL
	rows, err := dbConn.Query(`
	select
		categories.pos as c_pos,
		categories.title as c_title,
		pos_in_category,
		ch_num,
		name,
		name_tr,
		url
	from playlist
	join categories on categories.id = playlist.category_id
	order by ch_num`)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Panic(err)
		}
	}()

	var (
		plRec  PlaylistRecord
		result []PlaylistRecord
	)

	for rows.Next() {
		err = rows.Scan(&plRec.CatPos, &plRec.CatTitle, &plRec.PosInCat,
			&plRec.ChNum, &plRec.Name, &plRec.NameTr, &plRec.Url)
		if err != nil {
			return nil, err
		}

		result = append(result, plRec)
	}

	return &result, nil
}

func m3u(cfg *Cfg, multicast bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		playlist, err := getPlaylist()
		if err != nil {
			handleErr(w)
		}

		data := "#EXTM3U m3uautoload=1 cache=500 deinterlace=1 aspect-ratio=none\n"

		for _, prec := range *playlist {
			data += fmt.Sprintf("#EXTINF:-1 tvg-logo=\"http://%s/static/logo/%s.png\" group-title=\"%s\",%s\n",
				r.Host, prec.Name, prec.CatTitle, prec.Name)

			if multicast {
				data += fmt.Sprintf("%s\n", prec.Url)
			} else {
				data += fmt.Sprintf("%s/%s\n", cfg.UnicastUrl, prec.NameTr)
			}
		}

		w.Header().Set("Content-type", "audio/x-mpegurl")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handleOk(w, &data)
	}
}

func xspf(cfg *Cfg, multicast bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		playlist, err := getPlaylist()
		if err != nil {
			handleErr(w)
		}

		data := `<?xml version="1.0" encoding="UTF-8"?>
		<playlist version="1" xmlns="http://xspf.org/ns/0/" xmlns:vlc="http://www.videolan.org/vlc/playlist/ns/0/">
		<trackList>
		`

		var chLocation string

		for _, prec := range *playlist {

			if multicast {
				chLocation = prec.Url
			} else {
				chLocation = fmt.Sprintf("%s/%s", cfg.UnicastUrl, prec.NameTr)
			}

			data += fmt.Sprintf(`
	    <track>
			<title>%s</title>
			<location>%s</location>
			<image>http://%s/static/logo/%s.png</image>
			<extension application="http://www.videolan.org/vlc/playlist/0">
                <vlc:id>%d</vlc:id>
            </extension>
        </track>
		`, prec.Name, chLocation, r.Host, prec.Name, prec.ChNum)
		}

		data += `
		</trackList>
	`
		data += `
		<extension application="http://www.videolan.org/vlc/playlist/0">
	`

		pl := *playlist
		sort.Slice(pl[:], func(i, j int) bool {
			if pl[i].CatPos == pl[j].CatPos {

				if pl[i].PosInCat == pl[j].PosInCat {
					return pl[i].ChNum < pl[j].ChNum
				}
				return pl[i].PosInCat < pl[j].PosInCat

			}
			return pl[i].CatPos < pl[j].CatPos
		})

		currentCat := ""
		for _, prec := range pl {
			if currentCat == "" {
				data += fmt.Sprintf(`
			<vlc:node title="%s">
			`, prec.CatTitle)
				currentCat = prec.CatTitle
			}
			if prec.CatTitle == currentCat {
				data += fmt.Sprintf(`
			   <vlc:item tid="%d"/>`, prec.ChNum)
			} else {
				data += fmt.Sprintf(`
			</vlc:node>
            <vlc:node title="%s">
			  <vlc:item tid="%d"/>`, prec.CatTitle, prec.ChNum)
				currentCat = prec.CatTitle
			}
		}

		data += `
			</vlc:node>
		</extension>
	</playlist>
	`

		w.Header().Set("Content-type", "application/xspf+xml")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handleOk(w, &data)
	}
}

func main() {

	fCfgPath := flag.String("c", "conf.json", "path to conf file")
	flag.Parse()

	cfg, err := LoadConfig(*fCfgPath)
	if err != nil {
		panic(err)
	}

	dbConn, err = Connect(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Panic(err)
		}
	}()
	dbConn.SetMaxIdleConns(100)

	http.HandleFunc("/pl.m3u", m3u(cfg, true))
	http.HandleFunc("/px.pl.m3u", m3u(cfg, false))

	http.HandleFunc("/pl.xspf", xspf(cfg, true))
	http.HandleFunc("/px.pl.xspf", xspf(cfg, false))

	log.Printf("run server on http://localhost:%d", cfg.Service.Port)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Service.Port), nil))
}

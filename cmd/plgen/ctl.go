package main

import (
	"bytes"
	"log"
	"net/http"
	"sort"
	"text/template"

	"github.com/gin-gonic/gin"

	"iptv-playlist/internal/store"
)

type CatMap map[string][]store.Record

type Playlist struct {
	Host       string
	Records    []store.Record
	UnicastUrl string
}

func SortInCat(recs []store.Record) []store.Record {
	pl := recs[:]
	sort.Slice(pl[:], func(i, j int) bool {
		if pl[i].CatPos == pl[j].CatPos {
			if pl[i].PosInCat == pl[j].PosInCat {
				return pl[i].ChNum < pl[j].ChNum
			}
			return pl[i].PosInCat < pl[j].PosInCat
		}
		return pl[i].CatPos < pl[j].CatPos
	})
	return pl
}

func (pl Playlist) RecordsByCategories() CatMap {
	res := make(CatMap)
	for _, rec := range pl.Records {
		_r := rec
		if _, ok := res[_r.CatTitle]; !ok {
			res[_r.CatTitle] = make([]store.Record, 0)
		}
		res[rec.CatTitle] = append(res[rec.CatTitle], rec)
	}
	for cat, recs := range res {
		res[cat] = SortInCat(recs)
	}
	return res
}

func handle(mime string, s store.IStore, t *template.Template, unicastUrl string,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error

		log.Println("new request", ctx.Request.RequestURI, ctx.Request.RemoteAddr)

		plData := Playlist{Host: ctx.Request.Host, UnicastUrl: unicastUrl}
		plData.Records, err = s.GetAll()
		if err != nil {
			panic(err)
		}

		var b []byte
		buf := bytes.NewBuffer(b)

		if err = t.Execute(buf, plData); err != nil {
			panic(err)
		}

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Data(http.StatusOK, mime, buf.Bytes())
	}
}

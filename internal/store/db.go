package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type IStore interface {
	GetAll() ([]Record, error)
}

type Db struct {
	conn *sql.DB
	dsn  string
}

func formatDsnMySQL(host string, port int, username, password, dbName string) string {
	var suffix string

	if host == "localhost" {
		suffix = "/" + dbName
	} else {
		suffix = fmt.Sprintf("tcp(%s:%d)/%s", host, port, dbName)
	}

	return fmt.Sprintf("%s:%s@%s", username, password, suffix)
}

func NewStore(host string, port int, username, password, dbName string) IStore {
	return &Db{conn: new(sql.DB), dsn: formatDsnMySQL(host, port, username, password, dbName)}
}

func (db *Db) connect() (err error) {
	db.conn, err = sql.Open("mysql", db.dsn)
	if err != nil {
		return
	}
	return db.conn.Ping()
}

func (db *Db) GetAll() (recs []Record, err error) {
	if err = db.connect(); err != nil {
		return nil, err
	}
	defer func() {
		if err := db.conn.Close(); err != nil {
			panic(err)
		}
	}()

	//goland:noinspection ALL
	rows, err := db.conn.Query(`
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

	plRec := new(Record)

	for rows.Next() {
		err = rows.Scan(&plRec.CatPos, &plRec.CatTitle, &plRec.PosInCat,
			&plRec.ChNum, &plRec.Name, &plRec.NameTr, &plRec.Url)
		if err != nil {
			return nil, err
		}
		recs = append(recs, *plRec)
	}

	return recs, nil
}

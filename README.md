## IPTV playlists generator

> m3u: http://iptv.ttnet.ru/pl.m3u
>
> xspf: http://iptv.ttnet.ru/pl.xspf

### How it works

- the service has two endpoints for `m3u` and `xspf` formats
- both handlers get data from MySQL database
- each handler generates and returns a playlist file in the certain format

logo images store at the same host in `static/logo/%s.png` files and updates by an operator manually via FTP
___
### Configuration

Configuration file `conf.json` (can be set in the `-c` running option) contains:

- `service` - running service parameters
- `database` - DB connection parameters
- `unicast_url` - URL of IPTV proxy unicast service

___
### Build and run
```bash
go mod download
go build -o plgen ./cmd/plgen
./plgen -c conf.json
```

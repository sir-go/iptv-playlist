# IPTV playlists generator
[![Go](https://github.com/sir-go/iptv-playlist/actions/workflows/go.yml/badge.svg)](https://github.com/sir-go/iptv-playlist/actions/workflows/go.yml)

> m3u unicast: http://iptv.ttnet.ru/pl.m3u
> m3u multicast: http://iptv.ttnet.ru/pl.m3u
>
> xspf unicast: http://iptv.ttnet.ru/pl.xspf
> xspf multicast: http://iptv.ttnet.ru/pl.xspf

## How it works
- the service has four endpoints for `m3u` and `xspf` formats for `unicast` and `multicast` streaming servers
- data stores in the MySQL database
- each handler generates and returns a playlist file in the certain format

logo images store at the same host in `static/logo/%s.png` files and updates by an operator manually via FTP

## Configuration
Configuration file `config.yml` (can be set in the `-c` running option) contains:
```yaml
service:
  port: 8081         # backend running port

db:                  # database connection
  host: iptv-db      # db server host
  port: 3306         # port
  user: root         # username (root in the test composer)
  password:          # db password
  dbname: iptv       # db name

playlist:                                   # playlist generation options
  templates:
    m3u: /opt/templates/m3u.tmpl            # m3u template
    xspf: /opt/templates/xspf.tmpl          # xspf template
  unicast_url: http://localhost:8000        # streaming url for unicast playlist

```

## Test
Run `docker compose` with testing `mysql` before
```bash
docker compose up -d
go test -v ./cmd/plgen
gosec ./...
```

## Docker
```bash
docker build -t iptv-back .
docker run --rm -it -p 8081:8081 --name iptv-back -v ${PWD}/config.yml:/opt/config.yml:ro iptv-back:latest
```

## Standalone
```bash
go mod download
go build -o plgen ./cmd/plgen

./plgen -c config.yml
```

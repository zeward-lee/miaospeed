# Building this fork

Personal build of [AirportR/miaospeed](https://github.com/AirportR/miaospeed) 4.7.8.
It is identical to upstream except for the private adapter names added to
`interfaces.AllProxyTypes`, and two `replace` lines in `go.mod` that point the
hihomo core and its shadowsocks2 fork at tagged builds.

The Go module path is intentionally left as `github.com/airportr/miaospeed` so the
imports do not have to be rewritten; the code is built from this repository, only
the module identity is unchanged.

## Dependencies

```bash
export GO111MODULE=on
export HTTPS_PROXY=http://127.0.0.1:7801
export HTTP_PROXY=http://127.0.0.1:7801
export GOPRIVATE='github.com/zeward-lee/*'
export GONOSUMDB='github.com/zeward-lee/*'
export GOSUMDB=off
```

`proxy.golang.org` does not index the two forks, so they are fetched directly
over HTTPS. `GOPRIVATE` also keeps `sum.golang.org` from rejecting them.

The replaces in `go.mod` are pinned by tag, not by a local path:

```
replace github.com/metacubex/mihomo        => github.com/zeward-lee/Clash.Meta v1.19.34-feebox.1
replace github.com/metacubex/sing-shadowsocks2 => github.com/zeward-lee/sing-shadowsocks2 v0.2.8-feebox.1
```

## Build-only inputs (never committed)

Two embedded paths are excluded by `.gitignore` because they carry secrets. Both
must exist before compiling, otherwise `go:embed` fails with
`pattern ...: no matching files found`.

1. `utils/embeded/BUILDTOKEN.key` - the build token, one line, taken from the
   Koipy config field `buildtoken`. It is baked into `hashMiaoSpeed` and must
   match what the server and clients expect; a different value makes every
   request fail challenge validation. Delete the file right after `go build`.
2. `preconfigs/embeded/tls/` - one or more `<name>.crt` / `<name>.key` pairs
   served as the embedded TLS certificate list (`miaospeed_tls`, `miaoko` win
   over other names). **Do not reuse the production pair** for a throwaway
   local build; a self-signed pair is enough to compile and to test.

```bash
mkdir -p preconfigs/embeded/tls
printf '%s\n' "$BUILDTOKEN" > utils/embeded/BUILDTOKEN.key
# self-signed placeholder, only for local/CI verification builds
openssl req -x509 -newkey rsa:2048 -nodes -days 1 \
  -keyout preconfigs/embeded/tls/local_dev.key \
  -out preconfigs/embeded/tls/local_dev.crt \
  -subj "/CN=localhost"
```

## Cross-compile (linux/arm64)

```bash
LDF="-X main.VERSION=4.7.8-feebox.1 \
     -X main.COMPILATIONTIME=$(date -u +%Y%m%dT%H%M%SZ) \
     -X main.COMMIT=$(git rev-parse --short HEAD) \
     -X main.MihomoVersion=v1.19.34-feebox.1 -w -s -buildid="

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOARM= GOAMD64= \
  go build -trimpath -ldflags "$LDF" -o bin/miaospeed-linux-arm64 .
```

Then remove the injected secrets from the working tree:

```bash
rm -f utils/embeded/BUILDTOKEN.key preconfigs/embeded/tls/*.crt preconfigs/embeded/tls/*.key
```

## Deploy / rollback (router)

The binary lives at `/mnt/sda1/mmdb/run/miaospeed` on the router and is started
by `/etc/init.d/miaospeed`. Copy the new binary as a new file and swap it in,
keeping the old one around for rollback:

```bash
# deploy
scp bin/miaospeed-linux-arm64 root@ROUTER:/mnt/sda1/mmdb/run/miaospeed.new
ssh root@ROUTER 'cd /mnt/sda1/mmdb/run && cp -a miaospeed "miaospeed.bak.$(date +%Y%m%d%H%M%S)" \
  && chmod 755 miaospeed.new && mv miaospeed.new miaospeed \
  && /etc/init.d/miaospeed restart'

# rollback: pick the backup you want and put it back
ssh root@ROUTER 'cd /mnt/sda1/mmdb/run && mv miaospeed.bak.<stamp> miaospeed \
  && /etc/init.d/miaospeed restart'
```

## Reproducibility notes

- `-trimpath` drops local absolute paths from the binary. Without it, two
  machines produce different hashes from identical sources.
- `-buildid=` (empty) removes the content-derived Go build ID, so the hash
  depends only on inputs, not on the order the toolchain linked objects.
- `main.COMPILATIONTIME` is a real timestamp; embed a fixed value if you need
  byte-identical output across time.
- The embedded build token and TLS bytes are part of the binary. A production
  build (real token + real certificate pair) therefore never has the same
  sha256 as a verification build that used placeholders.

To check what a binary was actually built from:

```bash
go version -m /mnt/sda1/mmdb/run/miaospeed
```

The `dep` lines show the `=>` targets; a tag target means the binary came from
this repository's default build, a filesystem path means it came from a local
checkout.

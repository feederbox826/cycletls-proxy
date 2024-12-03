# cycletls-proxy

A lightweight http (MITM) proxy with [CycleTLS](https://github.com/Danny-Dasilva/CycleTLS) built in.

Usage:
```sh
$ curl -pkx http://localhost:8080 https://tls.peet.ws/api/all
```

See `docker-compose.yml` for docker

Because the proxy needs MITM to decode http responses, `-k` is needed to allow the unsigned MITM certificate and `-p` allows explicit tunneling via CONNECT.

## Environment Variables
```ini
TLS_JA3 = 771,4865-... # unhashed JA3 string
TLS_UA = Mozilla/5.0... # forced useragent
TLS_EXPOSE_UA = FALSE # use client useragent instead of TLS_UA
TLS_PROXY_ADDR = :8080 # listen address and port
```
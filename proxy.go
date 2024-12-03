package main

import (
	"io/ioutil"
	"flag"
	"log"
	"net"
	"net/http"
	"regexp"
	"os"
	"github.com/elazarl/goproxy"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func lookupEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// set cycleclient via env
	var ENV_JA3 = lookupEnv("TLS_JA3", "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0")
	var ENV_UA = lookupEnv("TLS_UA", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0")
	// enable cycleclient
	tlsClient := cycletls.Init()
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).
		HandleConnect(goproxy.AlwaysMitm)
	// enable curl -p for all hosts on port 80
	proxy.OnRequest().
		HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
		defer func() {
			if e := recover(); e != nil {
				ctx.Logf("error connecting to remote: %v", e)
				client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
			}
			client.Close()
		}()
		bodyData, _ := ioutil.ReadAll(req.Body)
		var headers = make(map[string]string)
		for k, v := range req.Header {
			headers[k] = v[0]
		}
		// override User-Agent UA
		headers["User-Agent"] = ENV_UA
		// log
		log.Printf("Request: %s %s %s\n", req.Method, req.URL.String(), headers)
		response, err := tlsClient.Do(req.URL.String(), cycletls.Options{
			Ja3: ENV_JA3,
			UserAgent: ENV_UA,
			Headers: headers,
			Body: string(bodyData),
		}, req.Method)
		orPanic(err)
		client.Write([]byte(response.Body))
	})
	verbose := flag.Bool("v", true, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"github.com/elazarl/goproxy"
)

var ENV_JA3 = lookupEnv("TLS_JA3", "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-34-18-51-43-13-45-28-27-65037,4588-29-23-24-25-256-257,0")
var ENV_UA = lookupEnv("TLS_UA", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:134.0) Gecko/20100101 Firefox/134.0")

func orPanic(err error) {
	if err != nil {
		log.Fatal(err) // Using log.Fatal for better error handling
	}
}

func lookupEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Custom struct to implement both io.Reader and io.ReadCloser
type ReadCloser struct {
	*strings.Reader
}

// Implement Close() method to satisfy io.ReadCloser
func (rc *ReadCloser) Close() error {
	// For strings.Reader, Close can be a no-op
	return nil
}

func tlsTripper(req *http.Request, ctx *goproxy.ProxyCtx) (resp *http.Response, err error) {
	// Set cycleclient via environment

	client := cycletls.Init()

	bodyData, err := ioutil.ReadAll(req.Body)
	orPanic(err) // Handle read error
	var headers = make(map[string]string)
	for k, v := range req.Header {
		headers[k] = v[0]
	}

	response, err := client.Do(req.URL.String(), cycletls.Options{
		Ja3:       ENV_JA3,
		UserAgent: ENV_UA,
		Headers:   headers,
		Body:      string(bodyData),
	}, req.Method)
	orPanic(err) // Handle request failure

	var Header = make(http.Header)
	for k, v := range response.Headers {
		Header.Add(k, v)
	}

	BodyReader := &ReadCloser{strings.NewReader(response.Body)}

	// Prepare the response for the proxy
	httpResponse := &http.Response{
		Request:    req,
		StatusCode: response.Status,
		Body:       BodyReader,
		Header:     Header,
	}
	ctx.UserData = httpResponse
	return httpResponse, nil
}

func main() {
	// Initialize proxy
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().
		HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.RoundTripper = goproxy.RoundTripperFunc(tlsTripper)
		return req, nil
	})

	addr := lookupEnv("TLS_PROXY_ADDR", ":8080")
	log.Fatal(http.ListenAndServe(addr, proxy)) // Start proxy server
}

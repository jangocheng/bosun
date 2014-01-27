package relay

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/StackExchange/tsaf/search"
)

func RelayHTTP(addr, dest string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handle(dest, w, r)
	})
	log.Println("OpenTSDB relay listening on:", addr)
	log.Println("OpenTSDB destination:", dest)
	return http.ListenAndServe(addr, mux)
}

func handle(dest string, w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	search.HTTPExtract(body)
	r.Body.Close()
	durl := url.URL{
		Scheme: "http",
		Host:   dest,
	}
	durl.Path = r.URL.Path
	durl.RawQuery = r.URL.RawQuery
	durl.Fragment = r.URL.Fragment
	req, err := http.NewRequest(r.Method, durl.String(), bytes.NewReader(body))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header = r.Header
	req.TransferEncoding = append(req.TransferEncoding, "identity")
	req.ContentLength = int64(len(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	w.Write(b)
}

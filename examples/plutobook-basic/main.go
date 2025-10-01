package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"

	"github.com/maitredede/puregolibs/plutobook"
)

//go:embed content.html
var kHTMLContent string

var (
	useGoHttp     bool
	verifySSLPeer bool
)

func main() {
	flag.BoolVar(&useGoHttp, "use-go-http", true, "use go http (instead of libcurl)")
	flag.BoolVar(&verifySSLPeer, "verify-ssl-peer", true, "verify ssl peer")
	flag.Parse()
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// handle Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	version := plutobook.Version()
	buildinfo := plutobook.BuildInfo()
	fmt.Printf("plutobook version: %s\n%s\n", version, buildinfo)

	book, err := plutobook.NewBook(plutobook.PageSizeA4, plutobook.PageMarginsNarrow, plutobook.MediaTypePrint)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer book.Close()

	if useGoHttp {
		// create a custom http fetcher
		cookies, err := cookiejar.New(nil)
		if err != nil {
			fmt.Println(fmt.Errorf("cookie jar create failed: %w", err))
			os.Exit(1)
		}
		transport := http.DefaultTransport.(*http.Transport).Clone()
		if !verifySSLPeer {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		httpClient := &http.Client{
			Jar:       cookies,
			Transport: transport,
		}

		var customFetcher plutobook.CustomResourceFetcher = buildResourceFetcher(ctx, httpClient)
		if err := book.SetCustomResourceFetcher(customFetcher); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		plutobook.SetSSLVerifyPeer(verifySSLPeer)
	}

	if err := book.LoadHTML(kHTMLContent, "", "", ""); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := book.WriteToPDF("hello.pdf"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildResourceFetcher(ctx context.Context, httpClient *http.Client) plutobook.CustomResourceFetcher {
	return func(url string) (*plutobook.ResourceData, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("new request failed: %w", err)
		}

		res, err := httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request execution failed: %w", err)
		}
		defer res.Body.Close()

		slog.Debug(fmt.Sprintf("go: loadResource url=%s status=%s", url, res.Status))
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unhandled status : %s", res.Status)
		}

		data := plutobook.ResourceData{}
		headerCT := res.Header.Values("Content-Type")
		if len(headerCT) > 0 {
			//TODO get content-type and text encoding
			data.Mime = headerCT[0]
		}
		slog.Debug(fmt.Sprintf("go: loadResource url=%s mime=%s textEncoding=%s", url, data.Mime, data.TextEncoding))
		bin, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("body read failed: %w", err)
		}
		data.Bin = bin
		return &data, nil
	}
}

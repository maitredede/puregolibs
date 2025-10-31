package plutobook

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"github.com/maitredede/puregolibs/tools/strings"
)

type resourceDataPtr unsafe.Pointer

type ResourceData struct {
	Mime         string
	Bin          []byte
	TextEncoding string
}

type CustomResourceFetcher func(url string) (*ResourceData, error)

func customResourceFetcherCallback(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	argArr := unsafe.Slice(args, cif.NArgs)
	if cif.NArgs < 2 {
		slog.Error("cif.NArgs < 2")
		*(*resourceDataPtr)(ret) = nil
		return 0
	}
	closureArgPtr := argArr[0]
	urlArgPtr := argArr[1]

	closurePtr := *(*unsafe.Pointer)(closureArgPtr)
	urlPtr := *(**byte)(urlArgPtr)

	book := (*Book)(closurePtr)
	url := strings.GoString(urlPtr)

	slog.Info(fmt.Sprintf("callback: loading url=%s", url))

	data, err := book.fetcher(url)
	if err != nil {
		slog.Error(err.Error())
		*(*resourceDataPtr)(ret) = nil
		return 0
	}

	cData := binPtr(&data.Bin[0])
	cLen := uint32(len(data.Bin))
	cMime := strings.CString(data.Mime)
	cEncoding := strings.CString(data.TextEncoding)
	resourcePtr := libResourceDataCreate(cData, cLen, stringPtr(cMime), stringPtr(cEncoding))
	*(*resourceDataPtr)(ret) = resourcePtr

	return 0
}

func DefaultHttpLoader(url string) (*ResourceData, error) {
	slog.Debug(fmt.Sprintf("go: loadResource url=%s", url))

	cookies, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookie jar create failed: %w", err)
	}
	client := &http.Client{
		Jar: cookies,
	}
	ctx := context.TODO()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request failed: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}
	defer res.Body.Close()

	slog.Debug(fmt.Sprintf("go: loadResource url=%s status=%s", url, res.Status))
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unhandled status : %s", res.Status)
	}

	data := ResourceData{}
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

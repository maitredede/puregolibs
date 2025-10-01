package plutobook

type ResourceData struct {
	Mime         string
	Bin          []byte
	TextEncoding string
}

type CustomResourceFetcher func(url string) (*ResourceData, error)

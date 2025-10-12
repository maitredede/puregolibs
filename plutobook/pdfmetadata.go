package plutobook

// PdfMetadata Defines different metadata fields for a PDF document
type PdfMetadata int32

const (
	PdfMetadataTitle PdfMetadata = iota
	PdfMetadataAuthor
	PdfMetadataSubject
	PdfMetadataKeywords
	PdfMetadataCreator
	PdfMetadataCreationDate
	PdfMetadataModificationDate
)

package sane

type Info = SANE_Int

const (
	InfoInexact       Info = 1
	InfoReloadOptions Info = 2
	InfoReloadParams  Info = 4
)

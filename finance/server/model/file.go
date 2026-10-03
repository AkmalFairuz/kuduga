package model

import "io"

type File struct {
	Name  string
	Bytes io.Reader
}

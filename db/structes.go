package db

import "os"

type DB struct {
	FileDataBase *os.File
	seek         int64
	errSeek      error
}

type indexTable struct {
	Name  [50]byte
	Index int64
	LenOfDataTable int64
}

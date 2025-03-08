package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

type PageReader struct {
	readInternal func(pageNumber int, dbName string) []byte
	conId        string
	dbName       string
}

func NewReader(conId string) *PageReader {
	return &PageReader{
		readInternal: readInternal,
		conId:        conId,
		dbName:       dbName,
	}
}

func (reader PageReader) readFromMemory(pageNumber int) PageParsed {
	for _, v := range softWritePages {
		if v.pageNumber == pageNumber {
			return v
		}
	}
	return parseReadPage(reader.readDbPage(pageNumber), pageNumber)
}

func (reader PageReader) readDbPage(pageNumber int) []byte {
	lockMutex.RLock()
	sharedLock[reader.conId] = struct{}{}

	page := reader.readInternal(pageNumber, reader.dbName)

	delete(sharedLock, reader.conId)
	lockMutex.RUnlock()
	return page
}

func readInternal(pageNumber int, dbName string) []byte {
	// err, ok := locks[LockTypeExclusive]
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(pwd + "/" + dbName + ".db"); os.IsNotExist(err) {
		return []byte{}
	}
	fd, err := os.Open(pwd + "/" + dbName + ".db")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	fd.Seek(int64(pageNumber)*int64(PageSize), 0)

	buff := make([]byte, PageSize)
	n, err := io.ReadFull(fd, buff)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading file:", err)
		return nil
	}

	return buff[:n]
}

func (reader PageReader) readInternalFileInfo() fs.FileInfo {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(pwd + "/" + reader.dbName + ".db"); os.IsNotExist(err) {
		return nil
	}
	fd, err := os.Open(pwd + "/" + reader.dbName + ".db")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	fileInfo, err := fd.Stat()
	if err != nil || fileInfo == nil {

		panic("Error while getting information about file")
	}

	return fileInfo
}

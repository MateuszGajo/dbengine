package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

type PageReader struct {
	readInternal func(pageNumber int) ([]byte, fs.FileInfo)
	conId        string
}

func NewReader(conId string) *PageReader {
	return &PageReader{
		readInternal: readInternal,
		conId:        conId,
	}
}

func (reader PageReader) readDbPage(pageNumber int) ([]byte, fs.FileInfo) {
	lockMutex.RLock()
	sharedLock[reader.conId] = struct{}{}

	page, info := reader.readInternal(pageNumber)

	delete(sharedLock, reader.conId)
	lockMutex.RUnlock()
	return page, info
}

func readInternal(pageNumber int) ([]byte, fs.FileInfo) {
	// err, ok := locks[LockTypeExclusive]
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(pwd + "/aa.db"); os.IsNotExist(err) {
		return []byte{}, nil
	}
	fd, err := os.Open(pwd + "/aa.db")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	fd.Seek(int64(pageNumber)*int64(PageSize), 0)

	buff := make([]byte, PageSize)
	n, err := io.ReadFull(fd, buff)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading file:", err)
		return nil, nil
	}

	fileInfo, err := fd.Stat()
	if err != nil || fileInfo == nil {

		panic("Error while getting information about file")
	}
	fmt.Println("size")
	fmt.Println(fileInfo.Size())

	return buff[:n], fileInfo
}

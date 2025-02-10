package main

import (
	"fmt"
	"os"
	"time"
)

type WriterStruct struct {
	retry          int
	writeToFileRaw func(data []byte, page int) error
}

// LEts think how to structure this

func NewWriter() *WriterStruct {
	return &WriterStruct{retry: 0, writeToFileRaw: writeToFileRaw}
}

func (writer WriterStruct) writeToFile(data []byte, page int, firstPage PageParsed, conId string) {

	if page == 0 {
		return
	}
	firstPage.dbHeader.fileChangeCounter++

	assembledPage := assembleDbPage(firstPage)
	writer.WriteToFileWithRetry(assembledPage, 0, conId)

}

func (writer *WriterStruct) WriteToFileWithRetry(data []byte, page int, conId string) {
	lockMutex.Lock()

	lockTypeExclusive = &conId
	err := writer.writeToFileRaw(data, page)
	if err != nil {
		if writer.retry == 3 {
			fmt.Println(err)
			panic("Hardware problem, can't write to disk")
		}
		time.Sleep(1 * time.Millisecond)
		writer.retry++
		writer.WriteToFileWithRetry(data, page, conId)
	}

	lockTypeExclusive = nil
	lockMutex.Unlock()

}

func writeToFileRaw(data []byte, page int) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(pwd+"/aa.db", os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}
	f.Seek(int64(page*PageSize), 0)

	_, err = f.Write(data)

	return err
}

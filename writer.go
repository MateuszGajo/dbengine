package main

import (
	"fmt"
	"os"
	"time"
)

type WriterStruct struct {
	retry          int
	writeToFileRaw func(data []byte, page int, dbName string) error
	dbName         string
}

func NewWriter() *WriterStruct {
	return &WriterStruct{
		retry:          0,
		writeToFileRaw: writeToFileRaw,
		dbName:         dbName,
	}
}

var softWritePages map[int]PageParsed = make(map[int]PageParsed)

// Action plan:
// TOOD: handle retry logic for old data
var tree bool = false

// WE update pages all over the code, so this information its important while saving
// we update some header information while saving file
// how to deal with header???

// What about use only header data in server struct?
// we update header all over the places, but how we attach to the first pagE?
// first page will have header and header size,
// lets try use it header insted of first page

func (writer WriterStruct) softwiteToFile(data PageParsed, page int, dbHeader *DbHeader) {
	// WE need update header here!!!
	dbHeader.fileChangeCounter++
	dbHeader.versionValidForNumber++
	if page == dbHeader.dbSizeInPages {
		dbHeader.dbSizeInPages++
	} else if page > dbHeader.dbSizeInPages {
		panic("don't leave empty space")
	}

	if page != data.pageNumber {
		panic("should never happend, page cant be different than data.pageNumber")
	}

	softWritePages[page] = data

	//Write header

	zeroPage, ok := softWritePages[0]

	if ok {
		zeroPage.dbHeader = *dbHeader
		softWritePages[0] = zeroPage
	} else {
		panic("should never happend, we shoul always keep zero page in memory")
	}

}

func (writer WriterStruct) flushPages(conId string) {

	for _, v := range softWritePages {
		fmt.Println("save page", v.pageNumber)
		if v.pageNumber == 0 {
			fmt.Println("show me page 0")
			fmt.Printf("%+v", v)
		}
		writer.writeToFile(assembleDbPage(v), v.pageNumber, conId)
		if v.isOverflow {
			fmt.Println("page number", v.pageNumber)
			fmt.Println("is overflow", !v.isSpace())
			panic("can't save overflow page")
		}
	}
}

func (writer WriterStruct) writeToFile(data []byte, page int, conId string) {

	_, exists := softWritePages[page]

	if exists {
		delete(softWritePages, page)
	}

	writer.WriteToFileWithRetry(data, page, conId)
	if page == 0 {
		return
	}

}

func (writer *WriterStruct) WriteToFileWithRetry(data []byte, page int, conId string) {
	lockMutex.Lock()

	lockTypeExclusive = &conId
	err := writer.writeToFileRaw(data, page, writer.dbName)
	if err != nil {
		if writer.retry == 3 {
			panic("Hardware problem, can't write to disk")
		}
		time.Sleep(1 * time.Millisecond)
		writer.retry++
		writer.WriteToFileWithRetry(data, page, conId)
	}

	lockTypeExclusive = nil
	lockMutex.Unlock()

}

func writeToFileRaw(data []byte, page int, dbName string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(pwd+"/"+dbName+".db", os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}
	f.Seek(int64(page*PageSize), 0)

	_, err = f.Write(data)

	return err
}

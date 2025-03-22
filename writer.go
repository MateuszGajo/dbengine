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
	softWritePages = make(map[int]PageParsed)
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

func (writer WriterStruct) softwiteToFile(data *PageParsed, page int, firstPage *PageParsed) {
	fmt.Println("soft write writer", page)
	if data.pageNumber == 0 {
		if !tree && data.btreeType == 0x05 {
			fmt.Println("change tree!!")
			tree = true
		}
		if tree && data.btreeType != 0x05 {
			// problem is with insert, is the tree type changed??
			fmt.Println("override tree????")
			fmt.Println(data.btreeType)
		}
	}
	fmt.Println("content area len", len(data.cellAreaParsed))
	firstPage.dbHeader.fileChangeCounter++
	firstPage.dbHeader.versionValidForNumber++
	if page == firstPage.dbHeader.dbSizeInPages {
		firstPage.dbHeader.dbSizeInPages++
	} else if page > firstPage.dbHeader.dbSizeInPages {
		fmt.Println(page)
		fmt.Println("is greater than number of paghes")
		fmt.Println(firstPage.dbHeader.dbSizeInPages)
		panic("don't leave empty space")
	}
	if !data.isSpace() {
		fmt.Println("st overflow to true??")
		data.isOverflow = true
	} else {
		data.isOverflow = false
	}

	softWritePages[page] = *data
}

func (writer WriterStruct) flushPages(conId string, firstPage *PageParsed) {

	for _, v := range softWritePages {
		fmt.Println("flush ``````")
		fmt.Println("flush ``````")
		fmt.Println("flush ``````")
		fmt.Printf("%+v", v)
		fmt.Println("flush ``````")
		writer.writeToFile(assembleDbPage(v), v.pageNumber, conId, firstPage)
		if v.isOverflow {
			panic("can't save overflow page")
		}
	}
	softWritePages = map[int]PageParsed{}
}

func (writer WriterStruct) writeToFile(data []byte, page int, conId string, firstPage *PageParsed) {
	if page == 0 {
		fmt.Println("-------------------")
		fmt.Println("-------------------")
		fmt.Println("---write to file 0---")
		fmt.Println("-------------------")
	}
	_, exists := softWritePages[page]

	if exists {
		delete(softWritePages, page)
	}

	fmt.Println(page)
	writer.WriteToFileWithRetry(data, page, conId)
	fmt.Println("hmm??")
	if page == 0 {
		return
	}

	firstPage.dbHeader.fileChangeCounter++
	firstPage.dbHeader.versionValidForNumber++
	if page == firstPage.dbHeader.dbSizeInPages {
		firstPage.dbHeader.dbSizeInPages++
	} else if page > firstPage.dbHeader.dbSizeInPages {
		fmt.Println(page)
		fmt.Println("is greater than number of paghes")
		fmt.Println(firstPage.dbHeader.dbSizeInPages)
		panic("don't leave empty space")
	}
	//TODO fix it

	// assembledPage := assembleDbPage(*firstPage)

	// writer.WriteToFileWithRetry(assembledPage, 0, conId)

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

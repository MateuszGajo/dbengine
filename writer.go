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

var memoryPages map[int]PageParsed = make(map[int]PageParsed)

// // let think again about attaching this header
// // Myabe server should flush data from time to time and only there it should be this header
// // I would like to remove header from this func entirly
// // 1. What about updating file change for counter? we can do it while saing
// // 1.1 Follow up: so what about all the incosistencies?
// // 1.1 I think we flush after evry query, if eighter everything went ok, or error and retry
// // 1.2 Same for version valid for number,
// // So what about check with dbSize in page, we can do it a flush phase

// func (writer WriterStruct) softwiteToFile(data PageParsed, page int, dbHeader *DbHeader) {
// 	// WE need update header here!!!
// 	// dbHeader.fileChangeCounter++
// 	// dbHeader.versionValidForNumber++
// 	// if page == dbHeader.dbSizeInPages {
// 	// 	dbHeader.dbSizeInPages++
// 	// } else if page > dbHeader.dbSizeInPages {
// 	// 	panic("don't leave empty space")
// 	// }

// 	if page != data.pageNumber {
// 		panic("should never happend, page cant be different than data.pageNumber")
// 	}

// 	softWritePages[page] = data

// 	//Write header

// 	// zeroPage, ok := softWritePages[0]

// 	// if ok {
// 	// zeroPage.dbHeader = *dbHeader
// 	// softWritePages[0] = zeroPage
// 	// }
// 	// else {
// 	// 	panic("should never happend, we shoul always keep zero page in memory")
// 	// }
// }

// maybe lets softwrite files add to server, then pass them to flush command
// flush command check current headers and makes decision

func (writer WriterStruct) flushPages(conId string, dbHeader *DbHeader, softWritePages map[int]PageParsed) {

	dbHeader.fileChangeCounter++
	dbHeader.versionValidForNumber++

	for _, v := range softWritePages {
		fmt.Println("save page", v.pageNumber)
		if v.pageNumber == dbHeader.dbSizeInPages {
			// dbHeader.dbSizeInPages++
		} else if v.pageNumber > dbHeader.dbSizeInPages {
			panic("don't leave empty space")
		}
		writer.writeToFile(assembleDbPage(v), v.pageNumber, conId)
		if v.isOverflow {
			panic("can't save overflow page" + fmt.Sprint(v.pageNumber))
		}
	}
}

func (writer WriterStruct) writeToFile(data []byte, page int, conId string) {

	// _, exists := softWritePages[page]

	// if exists {
	// 	delete(softWritePages, page)
	// }

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

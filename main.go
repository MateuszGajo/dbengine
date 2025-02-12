package main

import (
	"fmt"
	"sync"
)

// Database Header Format
// Offset	Size	Description
// 0	16	The header string: "SQLite format 3\000"
// 16	2	The database page size in bytes. Must be a power of two between 512 and 32768 inclusive, or the value 1 representing a page size of 65536.
// 18	1	File format write version. 1 for legacy; 2 for WAL.
// 19	1	File format read version. 1 for legacy; 2 for WAL.
// 20	1	Bytes of unused "reserved" space at the end of each page. Usually 0.
// 21	1	Maximum embedded payload fraction. Must be 64.
// 22	1	Minimum embedded payload fraction. Must be 32.
// 23	1	Leaf payload fraction. Must be 32.
// 24	4	File change counter.
// 28	4	Size of the database file in pages. The "in-header database size".
// 32	4	Page number of the first freelist trunk page.
// 36	4	Total number of freelist pages.
// 40	4	The schema cookie.
// 44	4	The schema format number. Supported schema formats are 1, 2, 3, and 4.
// 48	4	Default page cache size.
// 52	4	The page number of the largest root b-tree page when in auto-vacuum or incremental-vacuum modes, or zero otherwise.
// 56	4	The database text encoding. A value of 1 means UTF-8. A value of 2 means UTF-16le. A value of 3 means UTF-16be.
// 60	4	The "user version" as read and set by the user_version pragma.
// 64	4	True (non-zero) for incremental-vacuum mode. False (zero) otherwise.
// 68	4	The "Application ID" set by PRAGMA application_id.
// 72	20	Reserved for expansion. Must be zero.

// 92	4	The version-valid-for number.
// 96	4	SQLITE_VERSION_NUMBER

var PageSize = 4096
var LegacyFileWriteFormat = 1
var LegacyFileReadFormat = 1

var (
	lockTypeExclusive *string = nil
	lockMutex         sync.RWMutex
	sharedLock        map[string]struct{} = map[string]struct{}{}
)

// Maybe lets create diffrent struct as READER/WRITER don't know

type DbInfo struct {
	pageNumber int
}

type DbHeader struct {
	headerString               []byte
	databasePageSize           int
	databaseFileWriteVersion   []byte
	databaseFileReadVersion    []byte
	reservedBytesSpace         []byte
	maxEmbeddedPayloadFraction []byte
	minEmbeddedPayloadFraction []byte
	leafPayloadFraction        []byte
	fileChangeCounter          int
	dbSizeInPages              int
	firstFreeListTrunkPage     []byte
	totalNumberOfFreeListPages []byte
	schemaCookie               int
	schemaFormatNumber         []byte
	defaultPageCacheSize       []byte
	largestBTreePage           []byte
	databaseEncoding           []byte
	userVersion                []byte
	incrementalVacuumMode      []byte
	applicationId              []byte
	reservedForExpansion       []byte
	versionValidForNumber      int
	sqlVersionNumber           []byte
}

// ````````````````````````````
// ````````````````````````````
// ````````````````````````````
// ````````TODO!!!!!!!``````
// ````````````````````````````
// ````````````````````````````
// ````````````````````````````
// 1. Error handling
// 2. Write some e2e test
// 2. Add data to have multiple pages, multiple pages for schema too, implement this

func exectueCommand(input string, pNumber int) {
	fmt.Println("run generic parser")
	_, parsedQuery := genericParser(input)

	server := ServerStruct{
		pageSize: PageSize,
		conId:    pseudo_uuid(),
		dbInfo:   DbInfo{},
	}

	data, fileInfo := NewReader(server.conId).readDbPage(0)

	if fileInfo != nil {
		server.dbInfo.pageNumber = int(fileInfo.Size() / int64(PageSize))
	} else {
		server.dbInfo.pageNumber = 0
	}

	parsedData := parseReadPage(data, 0)

	fmt.Println("first page readed")
	fmt.Printf("%+v", parsedData)

	server.firstPage = parsedData

	err := server.handleActionType(parsedQuery, input, parsedData)

	if err != nil {
		fmt.Println(err)
	}
}

type ServerStruct struct {
	firstPage PageParsed
	pageSize  int
	conId     string
	dbInfo    DbInfo
}

func main() {

	input := "CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)"
	exectueCommand(input, 0)
	// writeExtraPageTMP()

	input = "INSERT INTO user (name) values('Alice')"
	exectueCommand(input, 2)

	//

}

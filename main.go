package main

import (
	"fmt"
	"io/fs"
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

type BtreeType int

const (
	TableBtreeLeafCell     BtreeType = 0x0d
	TableBtreeInteriorCell BtreeType = 0x05
	IndexBtreeLeafCell     BtreeType = 0x0a
	IndexBtreeInteriorCell BtreeType = 0x02
)

type QueryType string

const (
	QueryTypeTable QueryType = "table"
	QueryTypeIndex QueryType = "index"
)

type PageParseColumn struct {
	columnType   string
	columnLength int
	columnValue  []byte
}

type LastPageParseLatestRow struct {
	rowId   int
	data    []byte
	columns []PageParseColumn
}

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

type PageParsed struct {
	dbHeader             DbHeader // only for first page
	dbHeaderSize         int
	btreeType            int
	freeBlock            int
	numberofCells        int
	startCellContentArea int
	framgenetedArea      int
	rightMostpointer     []byte
	pointers             []byte
	cellArea             []byte
	latesRow             *LastPageParseLatestRow
	dbInfo               DbInfo
}

type SQLQueryActionType string

const (
	SqlQueryCreateActionType SQLQueryActionType = "create"
	SqlQueryInsertActionType SQLQueryActionType = "insert"

	//TODO: fill them
)

var validQueryActionTypes = map[SQLQueryActionType]struct{}{
	SqlQueryCreateActionType: {},
	SqlQueryInsertActionType: {},
}

type SQLCreateQueryObjectType string

const (
	SqlQueryDatabaseObjectType SQLCreateQueryObjectType = "database"
	SqlQueryTableObjectType    SQLCreateQueryObjectType = "table"
	SqlQueryIndexObjectType    SQLCreateQueryObjectType = "index"
	//TODO fill them
)

var validCreateQueryObjectTypes = map[SQLCreateQueryObjectType]struct{}{
	SqlQueryDatabaseObjectType: {},
	SqlQueryTableObjectType:    {},
	SqlQueryIndexObjectType:    {},
}

var validInsertObjectType = "into"

type SQLQueryColumnConstrains struct {
	columnName string
	columnType string
	constrains []string
}

type UserData struct {
	input     string
	queryData []string
	sqlType   SQLQueryActionType
}

type CreateCell struct {
	dataLength int
	data       []byte
}

type SQLQueryColumnAttribute string

const (
	SQLQueryColumnAttributePrimaryKey SQLQueryColumnAttribute = "PRIMARY KEY"
	SQLQueryColumnAttributeForeignKey SQLQueryColumnAttribute = "FOREIGN KEY"
	SQLQueryColumnAttributeUniuq      SQLQueryColumnAttribute = "UNIQUE"
	SQLQueryColumnAttributeNotNull    SQLQueryColumnAttribute = "NOT NULL"
	SQLQueryColumnAttributeIndex      SQLQueryColumnAttribute = "INDEX"
	// TODO: fill it
)

var sqlQueryAllowedColumnAttributes = map[SQLQueryColumnAttribute]struct{}{
	SQLQueryColumnAttributeForeignKey: {},
	SQLQueryColumnAttributePrimaryKey: {},
	SQLQueryColumnAttributeUniuq:      {},
	SQLQueryColumnAttributeNotNull:    {},
	SQLQueryColumnAttributeIndex:      {},
	// TODO: fill it after const
}

type SQLQueryColumnType string

const (
	SQLQueryColumnTypeInteger SQLQueryColumnType = "INTEGER"
	SQLQueryColumnTypeText    SQLQueryColumnType = "TEXT"
	// TODO: fill it
)

var sqlQueryAllowedColumnType = map[SQLQueryColumnType]struct{}{
	SQLQueryColumnTypeInteger: {},
	SQLQueryColumnTypeText:    {},
}

type CreateActionQueryData struct {
	action     SQLQueryActionType
	objectType SQLCreateQueryObjectType
	entityName string
	columns    []SQLQueryColumnConstrains
	rawQuery   string
}

// WE need to have some locking, lets do for bow exclusive lock only

// Exclusive lock,
// cache is invalidated, so no one can read from page
//generally lock are on db file not on cache used to retrive pages faster

// As first implementation we can do exclusive lock, and then implement more granuall ones

// WE we need to implement
// if transaction is writing its block all other transaction
// if we have two insert transaction and second one has stall data and tries to write stall data, it should be detected at write, and retry implement to get current data and write again

// Action plan:
// 3 while inserting check for lock, in case data has changed need to retry logic

func exectueCommand(input string, pNumber int) {
	// res := parseStartQuery(input)
	fmt.Println("run generic parser")
	_, parsedQuery := genericParser(input)

	server := ServerStruct{
		pageSize:     PageSize,
		conId:        pseudo_uuid(),
		readInternal: readInternal,
	}

	// We always need to rage page 0
	data, fileInfo := NewReader(server.conId).readDbPage(0)

	parsedData := parseReadPage(data, 0, fileInfo)

	server.firstPage = parsedData

	// fmt.Println("parsed last page")
	// fmt.Printf("%+v", parsedData)

	server.handleActionType(parsedQuery, input, parsedData)

}

type ServerStruct struct {
	firstPage    PageParsed
	pageSize     int
	conId        string
	readInternal func(pageNumber int) ([]byte, fs.FileInfo)
}

func main() {

	input := "CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)"
	exectueCommand(input, 0)
	// writeExtraPageTMP()

	// input = "INSERT INTO user (name) values('Alice')"
	// exectueCommand(input, 2)

	//

}

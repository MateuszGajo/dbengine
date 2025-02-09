package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/fs"
	"os"
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

func intToBinary(val int, size int) []byte {
	resBinary := make([]byte, size)
	if size == 1 {
		resBinary = []byte{byte(val)}
	} else if size == 2 {
		binary.BigEndian.PutUint16(resBinary, uint16(val))
	} else if size == 4 {
		binary.BigEndian.PutUint32(resBinary, uint32(val))
	} else {
		panic("doesn support this size")
	}

	return resBinary
}

var PageSize = 4096
var LegacyFileWriteFormat = 1
var LegacyFileReadFormat = 1

func BtreeHeaderValue() []byte {

	// 1	The one-byte flag at offset 0 indicating the b-tree page type.

	// A value of 2 (0x02) means the page is an interior index b-tree page.
	// A value of 5 (0x05) means the page is an interior table b-tree page.
	// A value of 10 (0x0a) means the page is a leaf index b-tree page.
	// A value of 13 (0x0d) means the page is a leaf table b-tree page.

	// Any other value for the b-tree page type is an error.

	bTreePageType := intToBinary(0x0d, 1)
	// 	The two-byte integer at offset 1 gives the start of the first freeblock on the page, or is zero if there are no freeblocks.
	firstFreeBlockOnPage := intToBinary(0, 2)
	// 	The two-byte integer at offset 3 gives the number of cells on the page.

	numberOfCells := intToBinary(1, 2)
	// we added one table so its 1
	// 5	2	The two-byte integer at offset 5 designates the start of the cell content area. A zero value for this integer is interpreted as 65536.
	//If a page contains no cells (which is only possible for a root page of a table that contains no rows) then the offset to the cell content area will equal the page size minus the bytes of reserved space. If the database uses a 65536-byte page size and the reserved space is zero (the usual value for reserved space) then the cell content offset of an empty page wants to be 65536. However, that integer is too large to be stored in a 2-byte unsigned integer, so a value of 0 is used in its place
	//
	startCell := 4096 - 8 - 2
	startCellContentArea := intToBinary(startCell, 2)
	// 	The one-byte integer at offset 7 gives the number of fragmented free bytes within the cell content area.
	framgentedFreeBytesWithingCellContentArea := intToBinary(0, 1)
	//The four-byte page number at offset 8 is the right-most pointer. This value appears in the header of interior b-tree pages only and is omitted from all other pages.
	// right most pointer points to address fro mthe left side u see on xxd
	// 00000000: 5351 4c69 7465 2066 6f72 6d61 7420 3300  SQLite format 3.
	// 00000010: 1000 0101 0040 2020 0000 0003 0000 0003  .....@  ........
	// 00000020: 0000 0000 0000 0000 0000 0002 0000 0004  ................
	// 00000030: 0000 0000 0000 0000 0000 0001 0000 0000  ................
	// 00000040: 0000 0000 0000 0000 0000 0000 0000 0000  ................
	// 00000050: 0000 0000 0000 0000 0000 0000 0000 0003  ................
	// 00000060: 002e 7689 0d00 0000 020f 6600 0fb6 0f66  ..v.......f....f
	// .... more bytes
	// 00000f60: 0000 0000 0000 4e02 0617 1717 017d 7461  ......N......}ta
	// 00000f70: 626c 6575 7365 7232 7573 6572 3203 4352  bleuser2user2.CR
	// 00000f80: 4541 5445 2054 4142 4c45 2075 7365 7232  EATE TABLE user2
	// 00000f90: 2869 6420 696e 7465 6765 7220 7072 696d  (id integer prim
	// 00000fa0: 6172 7920 6b65 792c 206e 616d 6532 3232  ary key, name222
	// 00000fb0: 2074 6578 7429 4801 0617 1515 0175 7461   text)H......uta
	// 00000fc0: 626c 6575 7365 7275 7365 7202 4352 4541  bleuseruser.CREA
	// 00000fd0: 5445 2054 4142 4c45 2075 7365 7228 6964  TE TABLE user(id
	// 00000fe0: 2069 6e74 6567 6572 2070 7269 6d61 7279   integer primary
	// 00000ff0: 206b 6579 2c20 6e61 6d65 2074 6578 7429   key, name text)

	// There is rightpointer at 0f66, at is pointing to line 0f60 and then add 6 bytes so it after 00, so starting from 4e which is length of cell,
	// there is also pointer 0fb6, that is point to 0fb0 and then add 6 bytes, so it stars from 48 which is cell length

	cellAreaPointerFirst := intToBinary(startCell, 2)
	cellAreaPointerTwo := intToBinary(0, 2)
	data := []byte{}
	data = append(data, bTreePageType...)
	data = append(data, firstFreeBlockOnPage...)
	data = append(data, numberOfCells...)
	data = append(data, startCellContentArea...)
	data = append(data, framgentedFreeBytesWithingCellContentArea...)
	data = append(data, cellAreaPointerFirst...)
	data = append(data, cellAreaPointerTwo...)

	return data
}

// Maybe lets create diffrent struct as READER/WRITER don't know
func (serverData *ServerStruct) writeToFile(data []byte, page int, firstPage PageParsed, conId string) {
	lockMutex.Lock()
	lockTypeExclusive = &conId
	serverData.writeToFileRaw(data, page)
	if page == 0 {
		return
	}
	firstPage.dbHeader.fileChangeCounter++

	assembledPage := assembleDbPage(firstPage)
	serverData.writeToFileRaw(assembledPage, 0)
	lockTypeExclusive = nil
	lockMutex.Unlock()

}

func writeToFileRaw(data []byte, page int) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(pwd+"/aa.db", os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		panic(err)
	}
	f.Seek(int64(page*PageSize), 0)

	_, err = f.Write(data)

	if err != nil {
		panic(err)
	}
}

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
	// input      string
	action     SQLQueryActionType
	objectType SQLCreateQueryObjectType
	entityName string
	columns    []SQLQueryColumnConstrains
	rawQuery   string
	// queryData  []string
}

func parseStartQuery(input string) []string {
	res := []string{}
	start := 0
	for i := 0; i < len(input); i++ {
		if input[i] == 32 {
			res = append(res, input[start:i])
			start = i + 1
		}

		if len(res) == 2 {
			break
		}
	}

	if len(res) < 2 {
		panic("query should have at lest 2 items")
	}

	return res
}

// WE need to have some locking, lets do for bow exclusive lock only

// type LockType string

// const (
// 	LockTypeExclusive = "ExclusiveLockType"
// )

// // sqlite stores lock in db i think
// var locks = map[LockType]*string{
// 	LockTypeExclusive: nil,
// } //

// Shared lock allow multiple select to read data but block writer, write can use WAL to save data without interuption

// Exclusive lock,
// cache is invalidated, so no one can read from page
//generally lock are on db file not on cache used to retrive pages faster

// As first implementation we can do exclusive lock, and then implement more granuall ones

// WE we need to implement
// if transaction is writing its block all other transaction
// if we have two insert transaction and second one has stall data and tries to write stall data, it should be detected at write, and retry implement to get current data and write again

// Action plan:
// 1. Add/remove exlucsive lock by process
// 2 While reading data check for lock
// 3 while inserting check for lock, in case data has changed need to retry logic

func pseudo_uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}

func exectueCommand(input string, pNumber int) {
	// res := parseStartQuery(input)
	fmt.Println("run generic parser")
	_, parsedQuery := genericParser(input)

	server := ServerStruct{
		pageSize:       PageSize,
		conId:          pseudo_uuid(),
		readInternal:   readInternal,
		writeToFileRaw: writeToFileRaw,
	}

	// We always need to rage page 0
	data, fileInfo := server.readDbPage(0)

	parsedData := parseReadPage(data, 0, fileInfo)

	server.firstPage = parsedData

	// fmt.Println("parsed last page")
	// fmt.Printf("%+v", parsedData)

	server.handleActionType(parsedQuery, input, parsedData)

}

type ServerStruct struct {
	firstPage      PageParsed
	pageSize       int
	conId          string
	readInternal   func(pageNumber int) ([]byte, fs.FileInfo)
	writeToFileRaw func(data []byte, page int)
}

func main() {

	input := "CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)"
	exectueCommand(input, 0)
	// writeExtraPageTMP()

	// input = "INSERT INTO user (name) values('Alice')"
	// exectueCommand(input, 2)

	//

}

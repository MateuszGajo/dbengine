package main

import (
	"encoding/binary"
	"fmt"
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

func header() []byte {
	headerString := []byte("SQLite format 3\000")
	pageSize := intToBinary(PageSize, 2)
	writeFileVersion := intToBinary(LegacyFileWriteFormat, 1)
	readFileVersion := intToBinary(LegacyFileReadFormat, 1)
	// SQLite has the ability to set aside a small number of extra bytes at the end of every page for use by extensions. These extra bytes are used, for example, by the SQLite Encryption Extension to store a nonce and/or cryptographic checksum associated with each page. The "reserved space" size in the 1-byte integer at offset 20 is the number of bytes of space at the end of each page to reserve for extensions. This value is usually 0. The value can be odd.
	reservedByte := intToBinary(0, 1)
	maxEmbededPayloadFranction := intToBinary(64, 1)
	minEmbededPayloadFranction := intToBinary(32, 1)
	//The maximum and minimum embedded payload fractions and the leaf payload fraction values must be 64, 32, and 32. These values were originally intended to be tunable parameters that could be used to modify the storage format of the b-tree algorithm. However, that functionality is not supported and there are no current plans to add support in the future. Hence, these three bytes are fixed at the values specified.
	leafPayloadFranction := intToBinary(32, 1)
	// The file change counter is a 4-byte big-endian integer at offset 24 that is incremented whenever the database file is unlocked after having been modified. When two or more processes are reading the same database file, each process can detect database changes from other processes by monitoring the change counter. A process will normally want to flush its database page cache when another process modified the database, since the cache has become stale. The file change counter facilitates this.
	fileChangeCounter := intToBinary(1, 4)
	sizeOfDataBaseInPages := intToBinary(2, 4)
	// 1 page for schemas, 2 page for values
	// END of 0000001
	//The trunk pages are the primary pages in the freelist structure. Each trunk page can list multiple leaf pages (individual unused pages) or point to other trunk pages.
	//When data is deleted from the database, the pages that were used to store that data aren't immediately erased or discarded. Instead, they are added to the freelist so that they can be reused later.
	pageNumberFirstFreeListTrunk := intToBinary(0, 4)
	totalNumerOfFreeListPages := intToBinary(0, 4)

	//The schema cookie is a 4-byte big-endian integer at offset 40 that is incremented whenever the database schema changes.
	schemaCookie := intToBinary(1, 4)
	// The schema format number is a 4-byte big-endian integer at offset 44. The schema format number is similar to the file format read and write version numbers at offsets 18 and 19 except that the schema format number refers to the high-level SQL formatting rather than the low-level b-tree formatting. Four schema format numbers are currently defined:

	//     Format 1 is understood by all versions of SQLite back to version 3.0.0 (2004-06-18).
	//     Format 2 adds the ability of rows within the same table to have a varying number of columns, in order to support the ALTER TABLE ... ADD COLUMN functionality. Support for reading and writing format 2 was added in SQLite version 3.1.3 on 2005-02-20.
	//     Format 3 adds the ability of extra columns added by ALTER TABLE ... ADD COLUMN to have non-NULL default values. This capability was added in SQLite version 3.1.4 on 2005-03-11.
	//     Format 4 causes SQLite to respect the DESC keyword on index declarations. (The DESC keyword is ignored in indexes for formats 1, 2, and 3.) Format 4 also adds two new boolean record type values (serial types 8 and 9). Support for format 4 was added in SQLite 3.3.0 on 2006-01-10.

	// New database files created by SQLite use format 4 by default. The legacy_file_format pragma can be used to cause SQLite to create new database files using format 1. The format version number can be made to default to 1 instead of 4 by setting SQLITE_DEFAULT_FILE_FORMAT=1 at compile-time.

	// If the database is completely empty, if it has no schema, then the schema format number can be zero.
	schemaFormatNumber := intToBinary(4, 4)
	// end of 0000002

	// The 4-byte big-endian signed integer at offset 48 is the suggested cache size in pages for the database file. The value is a suggestion only and SQLite is under no obligation to honor it. The absolute value of the integer is used as the suggested size. The suggested cache size can be set using the default_cache_size pragma.
	defaultPageSize := intToBinary(0, 4)
	// The page number of the largest root b-tree page when in auto-vacuum or incremental-vacuum modes, or zero otherwise
	// The two 4-byte big-endian integers at offsets 52 and 64 are used to manage the auto_vacuum and incremental_vacuum modes. If the integer at offset 52 is zero then pointer-map (ptrmap) pages are omitted from the database file and neither auto_vacuum nor incremental_vacuum are supported. If the integer at offset 52 is non-zero then it is the page number of the largest root page in the database file, the database file will contain ptrmap pages, and the mode must be either auto_vacuum or incremental_vacuum. In this latter case, the integer at offset 64 is true for incremental_vacuum and false for auto_vacuum. If the integer at offset 52 is zero then the integer at offset 64 must also be zero.
	pageNumberOfLargestBtreeToAutovacuum := intToBinary(0, 4)
	//The database text encoding. A value of 1 means UTF-8. A value of 2 means UTF-16le. A value of 3 means UTF-16be.
	databaseTextEncoding := intToBinary(1, 4)
	// The "user version" as read and set by the user_version pragma.
	userVersionNumber := intToBinary(0, 4)
	// end of  end of 0000003
	//// 64	4	True (non-zero) for incremental-vacuum mode. False (zero) otherwise.
	incrementalVacuumMode := intToBinary(0, 4)
	// The "Application ID" set by PRAGMA application_id.
	// The 4-byte big-endian integer at offset 68 is an "Application ID" that can be set by the PRAGMA application_id command in order to identify the database as belonging to or associated with a particular application. The application ID is intended for database files used as an application file-format. The application ID can be used by utilities such as file(1) to determine the specific file type rather than just reporting "SQLite3 Database". A list of assigned application IDs can be seen by consulting the magic.txt file in the SQLite source repository.
	applicationId := intToBinary(0, 4)
	// Reserved for expansion. Must be zero.
	reservedForExpansion := make([]byte, 20)
	// The 4-byte big-endian integer at offset 96 stores the SQLITE_VERSION_NUMBER value for the SQLite library that most recently modified the database file. The 4-byte big-endian integer at offset 92 is the value of the change counter when the version number was stored. The integer at offset 92 indicates which transaction the version number is valid for and is sometimes called the "version-valid-for number".
	// The change counter is incremented each time the database schema is modified, such as when:

	// A table is created or dropped.
	// An index is added or removed.
	// Other schema-altering operations occur.
	// we added one table
	versionValidForNumber := intToBinary(2, 4)
	// end of 0000005
	versionNumber := intToBinary(3045001, 4)

	data := headerString
	data = append(data, pageSize...)
	data = append(data, writeFileVersion...)
	data = append(data, readFileVersion...)
	data = append(data, reservedByte...)
	data = append(data, maxEmbededPayloadFranction...)
	data = append(data, minEmbededPayloadFranction...)
	data = append(data, leafPayloadFranction...)
	data = append(data, fileChangeCounter...)
	data = append(data, sizeOfDataBaseInPages...)
	data = append(data, pageNumberFirstFreeListTrunk...)
	data = append(data, totalNumerOfFreeListPages...)
	data = append(data, schemaCookie...)
	data = append(data, schemaFormatNumber...)
	data = append(data, defaultPageSize...)
	data = append(data, pageNumberOfLargestBtreeToAutovacuum...)
	data = append(data, databaseTextEncoding...)
	data = append(data, userVersionNumber...)
	data = append(data, incrementalVacuumMode...)
	data = append(data, applicationId...)
	data = append(data, reservedForExpansion...)
	data = append(data, versionValidForNumber...)
	data = append(data, versionNumber...)

	return data
}

func BtreeHeaderSchema() []byte {

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
	startCell := 4096 - 73 - 2
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

	rightMostPointerFirstPart := intToBinary(startCell, 2)
	rightMostPointerSecondPart := intToBinary(0, 2)
	data := []byte{}
	data = append(data, bTreePageType...)
	data = append(data, firstFreeBlockOnPage...)
	data = append(data, numberOfCells...)
	data = append(data, startCellContentArea...)
	data = append(data, framgentedFreeBytesWithingCellContentArea...)
	data = append(data, rightMostPointerFirstPart...)
	data = append(data, rightMostPointerSecondPart...)

	return data
}

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

	rightMostPointerFirstPart := intToBinary(startCell, 2)
	rightMostPointerSecondPart := intToBinary(0, 2)
	data := []byte{}
	data = append(data, bTreePageType...)
	data = append(data, firstFreeBlockOnPage...)
	data = append(data, numberOfCells...)
	data = append(data, startCellContentArea...)
	data = append(data, framgentedFreeBytesWithingCellContentArea...)
	data = append(data, rightMostPointerFirstPart...)
	data = append(data, rightMostPointerSecondPart...)

	return data
}

func writeToFile(data []byte) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(pwd+"/aa.db", data, 0644)

	if err != nil {
		panic(err)
	}

}

func appendQuery() []byte {
	// table := "create table user(id integer primary key, name text)"
	// table := "CREATE TABLE user (id INTEGER PRIMARY KEY, name TEXT)"

	// insertValueId := nil

	cellLength := 73                          // 2 bytes, i think can be extended to more than 2???
	emptyspace := 4096 - cellLength - 112 - 2 // page size - cell length - 112 headers, 2 bytes for byte for length and row id
	emptySpaceBytes := make([]byte, emptyspace)
	rowId := 1         // first row 2 byes
	headerLength := 6  // number of column + current field 2 bytes
	tableType1 := 23   // for table, 23-17/2 = 5 bytes and type text for 'table'
	tableType2 := 21   // for table, 21-15/2 = 4 bytes and type text for 'user'
	tableType3 := 21   // for table, 21-15/2 = 4 bytes and type text for 'user'
	columnType1 := 01  // integer for primary key
	columnType2 := 119 // text 117-13/2 = 52 this is for CREATE TABLE user(id INTEGER PRIMARY KEY, name TEXT)

	tableSqlpart1 := "tableuseruser"
	tableSqlpart2Id := 2
	tableSqlpart3 := "CREATE TABLE user (id INTEGER PRIMARY KEY, name TEXT)"
	// insert := "insert into user(name) values('Alice')"
	// insertValueName := "Alice"

	data := []byte{}
	data = append(data, emptySpaceBytes...)
	data = append(data, intToBinary(cellLength, 1)...)
	data = append(data, intToBinary(rowId, 1)...)
	data = append(data, intToBinary(headerLength, 1)...)
	data = append(data, intToBinary(tableType1, 1)...)
	data = append(data, intToBinary(tableType2, 1)...)
	data = append(data, intToBinary(tableType3, 1)...)
	data = append(data, intToBinary(columnType1, 1)...)
	data = append(data, intToBinary(columnType2, 1)...)
	data = append(data, []byte(tableSqlpart1)...)
	data = append(data, intToBinary(tableSqlpart2Id, 1)...)
	data = append(data, []byte(tableSqlpart3)...)

	return data
}

func appendValues() []byte {
	// table := "create table user(id integer primary key, name text)"
	// table := "CREATE TABLE user (id INTEGER PRIMARY KEY, name TEXT)"

	// insertValueId := nil

	cellLength := 8                          // 2 bytes, i think can be extended to more than 2???
	emptyspace := 4096 - cellLength - 12 - 2 // page size - cell length - 112 headers, 2 bytes for byte for length and row id
	emptySpaceBytes := make([]byte, emptyspace)
	rowId := 1        // first row 2 byes
	headerLength := 3 // number of column + current field 2 bytes
	columnType1 := 0  // id is null
	columnType2 := 23 // for table, 23-15/2 = 5 bytes and type text for 'alice'

	value := "Alice"
	// insert := "insert into user(name) values('Alice')"
	// insertValueName := "Alice"

	data := []byte{}
	data = append(data, emptySpaceBytes...)
	data = append(data, intToBinary(cellLength, 1)...)
	data = append(data, intToBinary(rowId, 1)...)
	data = append(data, intToBinary(headerLength, 1)...)
	data = append(data, intToBinary(columnType1, 1)...)
	data = append(data, intToBinary(columnType2, 1)...)
	data = append(data, []byte(value)...)

	return data
}

func main() {
	fmt.Println("Let's start with db")

	header := header()
	btreeHeaderSchema := BtreeHeaderSchema()
	cell := appendQuery()
	btreeHeaderValue := BtreeHeaderValue()
	value := appendValues()
	allData := header
	allData = append(allData, btreeHeaderSchema...)
	allData = append(allData, cell...)
	allData = append(allData, btreeHeaderValue...)
	allData = append(allData, value...)
	writeToFile(allData)

}

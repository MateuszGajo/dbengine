package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
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

type BtreeType int

const (
	TableBtreeLeafCell     BtreeType = 0x0d
	TableBtreeInteriorCell BtreeType = 0x5
	IndexBtreeLeafCell     BtreeType = 0x0a
	IndexBtreeInteriorCell BtreeType = 0x02
)

func IndexInteriorBtreeHeaderValue() []byte {

	// 1	The one-byte flag at offset 0 indicating the b-tree page type.

	// A value of 2 (0x02) means the page is an interior index b-tree page.
	// A value of 5 (0x05) means the page is an interior table b-tree page.
	// A value of 10 (0x0a) means the page is a leaf index b-tree page.
	// A value of 13 (0x0d) means the page is a leaf table b-tree page.

	// Any other value for the b-tree page type is an error.

	headerLength := 12 // contains also right-most pointer

	bTreePageType := intToBinary(int(IndexBtreeInteriorCell), 1)
	// 	The two-byte integer at offset 1 gives the start of the first freeblock on the page, or is zero if there are no freeblocks.
	firstFreeBlockOnPage := intToBinary(0, 2)
	// 	The two-byte integer at offset 3 gives the number of cells on the page.

	numberOfCells := intToBinary(1, 2)
	// we added one table so its 1
	// 5	2	The two-byte integer at offset 5 designates the start of the cell content area. A zero value for this integer is interpreted as 65536.
	//If a page contains no cells (which is only possible for a root page of a table that contains no rows) then the offset to the cell content area will equal the page size minus the bytes of reserved space. If the database uses a 65536-byte page size and the reserved space is zero (the usual value for reserved space) then the cell content offset of an empty page wants to be 65536. However, that integer is too large to be stored in a 2-byte unsigned integer, so a value of 0 is used in its place
	//
	startCell := 4096 - headerLength - 2
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

	rightMostPointer := intToBinary(0, 4)
	cellAreaPointerFirst := intToBinary(startCell, 2)

	data := []byte{}
	data = append(data, bTreePageType...)
	data = append(data, firstFreeBlockOnPage...)
	data = append(data, numberOfCells...)
	data = append(data, startCellContentArea...)
	data = append(data, framgentedFreeBytesWithingCellContentArea...)
	data = append(data, rightMostPointer...)

	// cellAreaPointerTwo := intToBinary(0, 2)

	data = append(data, cellAreaPointerFirst...)
	// data = append(data, cellAreaPointerTwo...)

	return data
}
func interiorBtreeHeaderValue() []byte {

	// 1	The one-byte flag at offset 0 indicating the b-tree page type.

	// A value of 2 (0x02) means the page is an interior index b-tree page.
	// A value of 5 (0x05) means the page is an interior table b-tree page.
	// A value of 10 (0x0a) means the page is a leaf index b-tree page.
	// A value of 13 (0x0d) means the page is a leaf table b-tree page.

	// Any other value for the b-tree page type is an error.

	headerLength := 8 // contains also right-most pointer

	bTreePageType := intToBinary(int(IndexBtreeInteriorCell), 1)
	// 	The two-byte integer at offset 1 gives the start of the first freeblock on the page, or is zero if there are no freeblocks.
	firstFreeBlockOnPage := intToBinary(0, 2)
	// 	The two-byte integer at offset 3 gives the number of cells on the page.

	numberOfCells := intToBinary(1, 2)
	// we added one table so its 1
	// 5	2	The two-byte integer at offset 5 designates the start of the cell content area. A zero value for this integer is interpreted as 65536.
	//If a page contains no cells (which is only possible for a root page of a table that contains no rows) then the offset to the cell content area will equal the page size minus the bytes of reserved space. If the database uses a 65536-byte page size and the reserved space is zero (the usual value for reserved space) then the cell content offset of an empty page wants to be 65536. However, that integer is too large to be stored in a 2-byte unsigned integer, so a value of 0 is used in its place
	//
	startCell := 4096 - headerLength - 2
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

	data := []byte{}
	data = append(data, bTreePageType...)
	data = append(data, firstFreeBlockOnPage...)
	data = append(data, numberOfCells...)
	data = append(data, startCellContentArea...)
	data = append(data, framgentedFreeBytesWithingCellContentArea...)

	// cellAreaPointerTwo := intToBinary(0, 2)

	data = append(data, cellAreaPointerFirst...)
	// data = append(data, cellAreaPointerTwo...)

	return data
}

func interiroBtreeValues() []byte {
	// 2 bytes, i think can be extended to more than 2???

	// rowId := 1        // first row 2 byes
	headerLength := 3 // number of column + current field 2 bytes
	columnType1 := 23 // for table, 23-15/2 = 5 bytes and type text for 'alice'
	columnType2 := 1  // id is null

	value1 := "Alice"
	value2 := 1
	// insert := "insert into user(name) values('Alice')"
	// insertValueName := "Alice"
	cellLength := 3 + len(value1) + 1           // header length + x2 column type     +length of values, +
	emptyspace := 4096 - cellLength - 8 - 2 - 1 // page size - cell length - 8 headers, 2 for one byte cell pointer, 1 bytes for byte for length
	emptySpaceBytes := make([]byte, emptyspace)

	data := []byte{}
	data = append(data, emptySpaceBytes...)
	data = append(data, intToBinary(cellLength, 1)...)
	data = append(data, intToBinary(headerLength, 1)...)
	data = append(data, intToBinary(columnType1, 1)...)
	data = append(data, intToBinary(columnType2, 1)...)
	data = append(data, []byte(value1)...)
	data = append(data, intToBinary(value2, 1)...)

	return data
}

func leafIndexBtreeHeaderValue() []byte {

	// 1	The one-byte flag at offset 0 indicating the b-tree page type.

	// A value of 2 (0x02) means the page is an interior index b-tree page.
	// A value of 5 (0x05) means the page is an interior table b-tree page.
	// A value of 10 (0x0a) means the page is a leaf index b-tree page.
	// A value of 13 (0x0d) means the page is a leaf table b-tree page.

	// Any other value for the b-tree page type is an error.

	headerLength := 8 // contains also right-most pointer

	bTreePageType := intToBinary(int(IndexBtreeLeafCell), 1)
	// 	The two-byte integer at offset 1 gives the start of the first freeblock on the page, or is zero if there are no freeblocks.
	firstFreeBlockOnPage := intToBinary(0, 2)
	// 	The two-byte integer at offset 3 gives the number of cells on the page.

	numberOfCells := intToBinary(1, 2)
	// we added one table so its 1
	// 5	2	The two-byte integer at offset 5 designates the start of the cell content area. A zero value for this integer is interpreted as 65536.
	//If a page contains no cells (which is only possible for a root page of a table that contains no rows) then the offset to the cell content area will equal the page size minus the bytes of reserved space. If the database uses a 65536-byte page size and the reserved space is zero (the usual value for reserved space) then the cell content offset of an empty page wants to be 65536. However, that integer is too large to be stored in a 2-byte unsigned integer, so a value of 0 is used in its place
	//
	startCell := 4096 - headerLength - 2
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

	data := []byte{}
	data = append(data, bTreePageType...)
	data = append(data, firstFreeBlockOnPage...)
	data = append(data, numberOfCells...)
	data = append(data, startCellContentArea...)
	data = append(data, framgentedFreeBytesWithingCellContentArea...)

	// cellAreaPointerTwo := intToBinary(0, 2)

	data = append(data, cellAreaPointerFirst...)
	// data = append(data, cellAreaPointerTwo...)

	return data
}

func leafIndexBtreeValues() []byte {
	// 2 bytes, i think can be extended to more than 2???

	// rowId := 1        // first row 2 byes
	headerLength := 3 // number of column + current field 2 bytes
	columnType1 := 23 // for table, 23-15/2 = 5 bytes and type text for 'alice'
	columnType2 := 1  // id is null

	value1 := "Alice"
	value2 := 1
	// insert := "insert into user(name) values('Alice')"
	// insertValueName := "Alice"
	cellLength := 3 + len(value1) + 1           // header length + x2 column type     +length of values, +
	emptyspace := 4096 - cellLength - 8 - 2 - 1 // page size - cell length - 8 headers, 2 for one byte cell pointer, 1 bytes for byte for length
	emptySpaceBytes := make([]byte, emptyspace)

	data := []byte{}
	data = append(data, emptySpaceBytes...)
	data = append(data, intToBinary(cellLength, 1)...)
	data = append(data, intToBinary(headerLength, 1)...)
	data = append(data, intToBinary(columnType1, 1)...)
	data = append(data, intToBinary(columnType2, 1)...)
	data = append(data, []byte(value1)...)
	data = append(data, intToBinary(value2, 1)...)

	return data
}

var input = ""

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

type QueryType string

const (
	QueryTypeTable QueryType = "table"
	QueryTypeIndex QueryType = "index"
)

// func uint16ToByte(val uint16) []byte {
// 	result := make([]byte, )
// 	binary.BigEndian.PutUint16(result, val)

// 	return result
// }

func calculateTextLength(value string) []byte {

	stringLen := 2*len(value) + 13

	if stringLen <= 255 {
		return []byte{byte(uint8(stringLen))}
	} else {
		//TODO: implement this
		panic("implement calculate text length")
	}
}

func createCell(latestRow LastPageParseLatestRow, values ...interface{}) (int, []byte) {
	var columnValues []byte = []byte{}
	var columnLength []byte = []byte{}
	var schemaRowId = latestRow.rowId

	schemaRowId++

	for _, v := range values {
		switch v.(type) {
		case int:
			value := v.(int)
			if value > 255 {
				panic("need to handle this type laters")
			}
			columnValues = append(columnValues, byte(uint16(value)))
			columnLength = append(columnLength, byte(1))
		case string:
			value := v.(string)
			columnValues = append(columnValues, []byte(value)...)
			columnLength = append(columnLength, calculateTextLength(value)...)
		default:
			panic("unssporrted cell type")
		}
	}

	headerLength := len(columnLength) + 1 // 5 column + 1 for current byte
	rowId := schemaRowId                  // first row 2 byes

	row := []byte{byte(rowId)}
	row = append(row, byte(headerLength))
	row = append(row, columnLength...)
	row = append(row, columnValues...)

	rowLength := len(row) - 1 // we don't count row id i guess

	result := []byte{byte(rowLength)}
	result = append(result, row...)

	return rowLength, result

}

func createSchemaCell(userval string, queryParts []string, latestRow LastPageParseLatestRow) (int, []byte) {
	// TODO: Should be read from previous row
	var internalIndexForSchema = 1

	internalIndexForSchema++

	if len(queryParts) < 3 {
		panic("query parts less than 3")
	}

	if len(queryParts) < 4 {
		queryParts = append(queryParts, queryParts[2])
	}

	columnOneValue := queryParts[1]
	columnTwoValue := queryParts[2]
	columnThreeValue := queryParts[3]
	columnFourValue := internalIndexForSchema
	columnFifthValue := userval

	return createCell(latestRow, columnOneValue, columnTwoValue, columnThreeValue, columnFourValue, columnFifthValue)

}

func uint16toByte(val uint16) []byte {
	result := make([]byte, 2)

	binary.BigEndian.PutUint16(result, val)

	return result
}

func BtreeHeaderSchema(rowPointerAdded []byte, rowAdded int, rowsAddedLength int, parsedData LastPageParsed) []byte {
	//This should be read from the page
	var currentNumberOfCell = parsedData.numberofCells
	var currentCellStart = parsedData.startCellContentArea

	currentNumberOfCell += rowAdded

	currentCellStart -= rowsAddedLength

	bTreePageType := intToBinary(int(TableBtreeLeafCell), 1)
	firstFreeBlockOnPage := intToBinary(0, 2)
	numberOfCells := intToBinary(currentNumberOfCell, 2)

	startCellContentArea := intToBinary(currentCellStart, 2)
	framgentedFreeBytesWithingCellContentArea := intToBinary(0, 1)

	data := []byte{}
	data = append(data, bTreePageType...)
	data = append(data, firstFreeBlockOnPage...)
	data = append(data, numberOfCells...)
	data = append(data, startCellContentArea...)
	data = append(data, framgentedFreeBytesWithingCellContentArea...)
	data = append(data, rowPointerAdded...)

	return data
}

func parseQuery(input string) {

}

func readDbPage(pageNumber int) []byte {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(pwd + "/aa.db"); os.IsNotExist(err) {
		return []byte{}
	}
	fd, err := os.Open(pwd + "/aa.db") // Adjust path as needed
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	fd.Seek(int64(pageNumber)*int64(PageSize), 0)

	buff := make([]byte, PageSize) // Create a buffer for 1024 bytes
	n, err := io.ReadFull(fd, buff)
	if err != nil && err != io.EOF { // Handle EOF and other errors
		fmt.Println("Error reading file:", err)
		return nil
	}

	fmt.Printf("Number of bytes read: %d\n", n)
	fmt.Printf("Data: %v\n", buff[:n]) // Only print the bytes read

	return buff[:n] // Return the bytes actually read
}

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

type LastPageParsed struct {
	btreeType            int
	freeBlock            int
	numberofCells        int
	startCellContentArea int
	framgenetedArea      int
	rightMostpointer     []byte
	pointers             [][]byte
	cellArea             []byte
	latesRow             LastPageParseLatestRow
}

func parseReadPage(data []byte, isFirstPage bool) LastPageParsed {
	if len(data) == 0 {
		return LastPageParsed{
			// btreeType:            int(TableBtreeLeafCell),
			numberofCells:        0,
			startCellContentArea: PageSize,
			cellArea:             []byte{},
			pointers:             [][]byte{},
			latesRow: LastPageParseLatestRow{
				rowId:   0,
				data:    []byte{},
				columns: []PageParseColumn{},
			},
		}
	}
	dataToParse := data
	if isFirstPage {
		//Skip header for now
		dataToParse = dataToParse[100:]
	}

	btreeType := dataToParse[0]
	isPointerValue := false
	switch BtreeType(btreeType) {
	case TableBtreeInteriorCell, IndexBtreeInteriorCell:
		isPointerValue = true

	}
	freebBlocks := dataToParse[1:3]
	if freebBlocks[0] != 0 {
		fmt.Println(freebBlocks)
		panic("implement free blocks more than 0 cell")
	}
	freeBlocksInt := int(freebBlocks[1])
	numberofCells := dataToParse[3:5]
	if numberofCells[0] != 0 {
		fmt.Println(numberofCells)
		panic("implement number of cell more than 0 cell")
	}
	numberofCellsInt := int(numberofCells[1])
	startCellContentArea := dataToParse[5:7]
	// if startCellContentArea[0] != 0 {
	// 	fmt.Println("start cell content area")
	// 	fmt.Println(startCellContentArea)
	// 	panic("implement startCellContentArea more than 0")
	// }
	startCellContentAreaInt := binary.BigEndian.Uint16(startCellContentArea)
	startCellContentAreaBigEndian := binary.BigEndian.Uint16(startCellContentArea)
	fragmenetedArea := dataToParse[7]
	var rightMostPointer []byte
	if isPointerValue {
		rightMostPointer = dataToParse[8:12]
		dataToParse = dataToParse[12:]
	} else {
		dataToParse = dataToParse[8:]
	}

	var pointers [][]byte

	for {
		pointer := dataToParse[:2]
		if pointer[0] == 0 && pointer[1] == 0 {
			break
		}
		dataToParse = dataToParse[2:]
		pointers = append(pointers, pointer)
	}

	if len(data) < int(startCellContentAreaBigEndian) {
		panic("data length is lesser than start of cell content area")
	}

	cellAreaContent := data[startCellContentAreaBigEndian:]
	var latestRowLengthArr []byte
	for i := 0; i < 9; i++ {
		latestRowLengthArr = append(latestRowLengthArr, cellAreaContent[i])
		if cellAreaContent[i] < 127 {
			break
		}
	}
	if len(latestRowLengthArr) > 1 {
		panic("Need to be handled later")
	}

	latestRowLength := int(latestRowLengthArr[0]) + 2 // 1 bytes for length, 1 bytes for row id
	if len(cellAreaContent) < int(latestRowLength) {
		panic("cellAreaContent length is lesser than start of cell content area, row length%")
	}
	latestRow := cellAreaContent[:latestRowLength]
	latestRowId := latestRow[1]
	latestRowheaderLength := latestRow[2]
	latestRowHeaders := latestRow[3 : 3-1+int(latestRowheaderLength)] // 3 - 1 (-1 because of header length contains itself)
	latestRowValues := latestRow[3-1+int(latestRowheaderLength):]
	var latestRowColumns []PageParseColumn
	for _, v := range latestRowHeaders {
		if int(v) > 127 {
			panic("handle case that we have multiple bytes")
		}
		if int(v) < 10 {
			//int
			column := PageParseColumn{
				columnType:   string(strconv.Itoa(int(v))),
				columnLength: 1,
				columnValue:  []byte{latestRowValues[0]},
			}
			latestRowColumns = append(latestRowColumns, column)
			latestRowValues = latestRowValues[1:]
			continue
		}
		if int(v) >= 10 && int(v) < 12 {
			panic("reserved values, shouldnt be used")
		}
		if int(v)%2 == 0 {
			//blob
			panic("implement hadnling blobs")

		} else {

			//string
			length := (int(v) - 13) / 2
			value := latestRowValues[:length]
			column := PageParseColumn{
				columnType:   "13",
				columnLength: length,
				columnValue:  value,
			}
			latestRowColumns = append(latestRowColumns, column)
			latestRowValues = latestRowValues[length:]
			continue
		}
		panic("should never enter this state in parsing")

	}
	fmt.Println("read what we have parsed")
	fmt.Printf("Header: Btree type: %v, freeBlocks: %v, number of cells: %v, fragmentedArea: %v, rightMostPOinter: %v \n", btreeType, freebBlocks, numberofCells, fragmenetedArea, rightMostPointer)
	fmt.Printf("pointer %v \n", pointers)
	fmt.Printf("Cell Area, latest Row length: %v, row: %v \n", latestRow, latestRow)
	fmt.Printf("latest row id: %v \n", latestRowId)
	fmt.Printf("Cell all: %v", cellAreaContent)

	return LastPageParsed{
		btreeType:            int(btreeType),
		numberofCells:        numberofCellsInt,
		startCellContentArea: int(startCellContentAreaInt),
		rightMostpointer:     rightMostPointer,
		cellArea:             cellAreaContent,
		framgenetedArea:      int(fragmenetedArea),
		freeBlock:            int(freeBlocksInt),
		pointers:             pointers,
		latesRow: LastPageParseLatestRow{
			rowId:   int(latestRowId),
			data:    latestRow,
			columns: latestRowColumns,
		},
	}
}

func main() {

	// "insert into user(name) values('alice')"
	// "select name from user where name='432423'"
	data := readDbPage(0)
	parsedData := parseReadPage(data, true)

	fmt.Printf("%+v", parsedData)

	// page := createSchema(true, parsedData, userData)

	// writeToFile(page)

}

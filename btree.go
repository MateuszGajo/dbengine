package main

import (
	"fmt"
	"os"
	"reflect"
)

// 1. We need sturcutre to have children ranges, and if it index, or if it a cell?

// inserting and balancing??
// We have spae
// 1. We start at root if there is a space we add there

// We dont have space on root page
// 1. root doesnt have space, we overflow and create a new page, special case for root
// 2. move all the data to new page, it overflow but we can reuse balacning log same as for leaf node
// 2. move all data to memory to page not be in overflow state
// 3. Calculate how many pages we need?
// 3.1 left to right, so we can first three put in first page and fourth enrry goes to parent, dummer look, we leave one page empty

/// OVERFLOW state * assuming we can only fit 3 pages in every node
///                                     +--------+ Root Node
///                  					|1,2,3,4 |
///                                     +--------+

///                                     +--------+ Root Node
///                   +-----------------|   4   | --------------------
///                  /                  +--------+
///                 /                                             \
///            +-------+                                     +----------+ empty
///       +----|  1,2,3  |----+                 +------------| 			|------------+
///      /     +-------+     \                 /             +----------+             \
///     /          |          \               /                /      \                \

// 3.2 Second look it goes right to left, readjust the entires so the tree is balances, move 3 as root, move 4 to righT???, no we have leaf baias distribution

// leaf bias distribution

type Cell struct {
	size       int
	pageNumber int
	rowId      int
	data       []byte
}

type CreateCell struct {
	dataLength int
	data       []byte
	rowId      int
}

type BtreeType int

const (
	TableBtreeLeafCell     BtreeType = 0x0d
	TableBtreeInteriorCell BtreeType = 0x05
	IndexBtreeLeafCell     BtreeType = 0x0a
	IndexBtreeInteriorCell BtreeType = 0x02
)

// fix this size later
var usableSpacePerPage = PageSize - 12 //

func EncodeVarint(n uint64) []byte {
	var groups []byte

	// Always append at least one group.
	groups = append(groups, byte(n&0x7F))
	n >>= 7

	// Extract 7-bit groups.
	for n > 0 {
		groups = append(groups, byte(n&0x7F))
		n >>= 7
	}

	// Reverse the groups to convert from little-endian (LSB first)
	// to big-endian (most-significant group first).
	for i, j := 0, len(groups)-1; i < j; i, j = i+1, j-1 {
		groups[i], groups[j] = groups[j], groups[i]
	}

	// For every byte except the last, set the continuation flag.
	for i := 0; i < len(groups)-1; i++ {
		groups[i] |= 0x80
	}

	return groups
}

// DecodeVarint decodes a byte slice encoded in SQLite's varint format
// back into a uint64. It also returns the number of bytes read.
func DecodeVarint(data []byte) uint64 {
	var n uint64
	var bytesRead int

	for i, b := range data {
		// For the first 8 bytes, each contributes 7 bits.
		if i < 8 {
			n = (n << 7) | uint64(b&0x7F)
			bytesRead++
			// If the continuation flag is not set, we're done.
			if b&0x80 == 0 && i == 7 {
				return n
			}
		} else {
			// The 9th byte contains 8 bits of data.
			n = (n << 8) | uint64(b)
			bytesRead++
			return n
		}
	}
	return n
}

func calculateTextLength(value string) []byte {

	stringLen := 2*len(value) + 13

	return EncodeVarint(uint64(stringLen))

	// if stringLen <= 127 {
	// 	return []byte{byte(uint8(stringLen))}
	// } else {
	// 	//TODO: implement this
	// 	panic("implement calculate text length")
	// }
}

func createCell(btreeType BtreeType, rowId int, values ...interface{}) CreateCell {
	if btreeType == TableBtreeLeafCell {
		var columnValues []byte = []byte{}
		var columnLength []byte = []byte{}
		// rowId++

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
			case nil:
				columnValues = append(columnValues, []byte{}...)
				columnLength = append(columnLength, 0)
			default:
				fmt.Fprintln(os.Stdout, []any{v}...)
				fmt.Println(reflect.TypeOf(v))
				panic("unssporrted cell type")
			}
		}

		headerLength := len(columnLength) + 1 // 5 column + 1 for current byte

		row := []byte{byte(rowId)}
		row = append(row, byte(headerLength))
		row = append(row, columnLength...)
		row = append(row, columnValues...)

		rowLength := len(row) - 1 // we don't count row id

		result := []byte{byte(rowLength)}
		result = append(result, row...)

		return CreateCell{
			dataLength: len(result),
			data:       result,
			rowId:      rowId,
		}
	}

	panic("not handle create cell")

}

func header() DbHeader {
	headerString := []byte("SQLite format 3\000")
	pageSize := PageSize
	writeFileVersion := intToBinary(LegacyFileWriteFormat, 1)
	readFileVersion := intToBinary(LegacyFileReadFormat, 1)
	// SQLite has the ability to set aside a small number of extra bytes at the end of every page for use by extensions. These extra bytes are used, for example, by the SQLite Encryption Extension to store a nonce and/or cryptographic checksum associated with each page. The "reserved space" size in the 1-byte integer at offset 20 is the number of bytes of space at the end of each page to reserve for extensions. This value is usually 0. The value can be odd.
	reservedByte := intToBinary(0, 1)
	maxEmbededPayloadFranction := intToBinary(64, 1)
	minEmbededPayloadFranction := intToBinary(32, 1)
	//The maximum and minimum embedded payload fractions and the leaf payload fraction values must be 64, 32, and 32. These values were originally intended to be tunable parameters that could be used to modify the storage format of the b-tree algorithm. However, that functionality is not supported and there are no current plans to add support in the future. Hence, these three bytes are fixed at the values specified.
	leafPayloadFranction := intToBinary(32, 1)
	// The file change counter is a 4-byte big-endian integer at offset 24 that is incremented whenever the database file is unlocked after having been modified. When two or more processes are reading the same database file, each process can detect database changes from other processes by monitoring the change counter. A process will normally want to flush its database page cache when another process modified the database, since the cache has become stale. The file change counter facilitates this.
	fileChangeCounter := 0
	sizeOfDataBaseInPages := 1
	// 1 page for schemas, 2 page for values
	// END of 0000001
	//The trunk pages are the primary pages in the freelist structure. Each trunk page can list multiple leaf pages (individual unused pages) or point to other trunk pages.
	//When data is deleted from the database, the pages that were used to store that data aren't immediately erased or discarded. Instead, they are added to the freelist so that they can be reused later.
	pageNumberFirstFreeListTrunk := intToBinary(0, 4)
	totalNumerOfFreeListPages := intToBinary(0, 4)

	//The schema cookie is a 4-byte big-endian integer at offset 40 that is incremented whenever the database schema changes.
	schemaCookie := 0
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
	// TODO: this need to be updated????
	versionValidForNumber := 0
	// end of 0000005
	versionNumber := intToBinary(3045001, 4)

	return DbHeader{
		headerString:               headerString,
		databasePageSize:           pageSize,
		databaseFileWriteVersion:   writeFileVersion,
		databaseFileReadVersion:    readFileVersion,
		reservedBytesSpace:         reservedByte,
		maxEmbeddedPayloadFraction: maxEmbededPayloadFranction,
		minEmbeddedPayloadFraction: minEmbededPayloadFranction,
		leafPayloadFraction:        leafPayloadFranction,
		fileChangeCounter:          fileChangeCounter,
		dbSizeInPages:              sizeOfDataBaseInPages,
		firstFreeListTrunkPage:     pageNumberFirstFreeListTrunk,
		totalNumberOfFreeListPages: totalNumerOfFreeListPages,
		schemaCookie:               schemaCookie,
		schemaFormatNumber:         schemaFormatNumber,
		defaultPageCacheSize:       defaultPageSize,
		largestBTreePage:           pageNumberOfLargestBtreeToAutovacuum,
		databaseEncoding:           databaseTextEncoding,
		userVersion:                userVersionNumber,
		incrementalVacuumMode:      incrementalVacuumMode,
		applicationId:              applicationId,
		reservedForExpansion:       reservedForExpansion,
		versionValidForNumber:      versionValidForNumber,
		sqlVersionNumber:           versionNumber,
	}

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

// sibling redistribution

// WE are missing divider entry, so 3 pages gotes to left page, then fourth one is the root

// add new page 12
// start

///                                +-------+
///                        +-------| 4,8   |-------+                           X
///                      /        +-------+        \
///                      /             |             \
///   			  +----------+  +----------+  +----------+
///   			  | 1,2,3    |  | 5, 6, 7  |  | 9, 10, 11|
/// 			  +----------+  +----------+  +----------+

// phse one
///                                +-------+
///                        +-------| 4,8,12 |-------+                           X
///                      /        +-------+        \
///                      /             |             \
///   			  +----------+  +----------+  +----------+
///   			  | 1,2,3    |  | 5, 6, 7  |  | 9, 10, 11|
/// 			  +----------+  +----------+  +----------+

/// phase two create empty page
//|                               +-------+
///                        +-------| 4,8,12 |-------+    -------                        X
///                      /        +-------+        \			\
///                      /             |             \			\
///   			  +----------+  +----------+  +----------+     +---------+
///   			  | 1,2,3    |  | 5, 6, 7  |  | 9, 10, 11|     |         |
/// 			  +----------+  +----------+  +----------+     +---------+

/// phase  three: second distribution loop readjust entry
//|                               +-------+
///                        +-------| 4,8,12 |-------+    -------                        X
///                      /        +-------+        \			\
///                      /             |             \			\
///   			  +----------+  +----------+  +----------+     +---------+
///   			  | 1,2,3    |  | 5, 6, 7  |  | 9, 10, 11|     |         |
/// 			  +----------+  +----------+  +----------+     +---------+

/// end
//|                               +-------+
///                        +-------| 4,7,10 |-------+    -------                        X
///                      /        +-------+        \			\
///                      /             |             \			\
///   			  +----------+  +----------+  +----------+     +---------+
///   			  | 1,2,3    |  | 5,6      |  | 8,9      |     |  11,12   |
/// 			  +----------+  +----------+  +----------+     +---------+

// func getSiblings() ([][]Cell, []Cell) {
// 	// Get only like 1 left 1 right sibiling, or sometime 2 left sibling
// 	divider := []Cell{{size: 1, number: 4}} // 2 divider, betwen first and secopd siblign, betwen second and third
// 	siblings := [][]Cell{[]Cell{{size: 1, number: 1}, {size: 1, number: 2}, {size: 1, number: 3}}, []Cell{{size: 1, number: 5}, {size: 1, number: 6}, {size: 1, number: 7}, {size: 1, number: 8}}}
// 	usableSpacePerPage = 3
// 	return siblings, divider
// }

func (parentPage PageParsed) findSiblings(currentNode PageParsed) (*PageParsed, *PageParsed) {

	if len(parentPage.cellAreaParsed) == 0 && parentPage.pageNumber == 0 {
		return nil, nil
	}
	if len(currentNode.cellAreaParsed) == 0 {
		panic("should never occur, cell area parsed empty, find siblings 1")
	}
	currnetPageLatestRow := parseCellArea(currentNode.cellAreaParsed[0], BtreeType(currentNode.btreeType))
	rowIdAsEntry := currnetPageLatestRow.rowId
	pageNumber := currentNode.pageNumber

	if len(parentPage.cellAreaParsed) == 0 {
		panic("should never occur, cell area parsed empty, find siblings 2")
	}

	reader := NewReader("")
	cellParsed := parseCellArea(parentPage.cellAreaParsed[0], BtreeType(parentPage.btreeType))

	rightRowId := cellParsed.rowId

	var leftPageNumber int
	var rightPageNumber int

	if rowIdAsEntry > rightRowId {

		panic("should never happend, rowIdASEntry > rightRowId in find sibling")
	}

	for i := 0; i < len(parentPage.cellAreaParsed); i++ {
		if i < len(parentPage.cellAreaParsed)-1 {
			cellParsed := parseCellArea(parentPage.cellAreaParsed[i+1], BtreeType(parentPage.btreeType))
			leftPageNumber = cellParsed.pageNumber
		}
		if i > 0 {
			cellParsed := parseCellArea(parentPage.cellAreaParsed[i-1], BtreeType(parentPage.btreeType))
			rightPageNumber = cellParsed.pageNumber
		}

		v := parentPage.cellAreaParsed[i]

		cellParsed := parseCellArea(v, BtreeType(parentPage.btreeType))
		rowId := cellParsed.rowId

		if rowIdAsEntry == rowId && cellParsed.pageNumber == pageNumber {
			var leftSiblings *PageParsed
			var rightSiblings *PageParsed
			if i < len(parentPage.cellAreaParsed)-1 {

				page := reader.readFromMemory(leftPageNumber)
				leftSiblings = &page
			}
			if i > 0 {

				page := reader.readFromMemory(rightPageNumber)
				rightSiblings = &page
			}
			// we're on the items we looking siblings for
			return leftSiblings, rightSiblings
		} else if rowIdAsEntry == rowId {
			panic("terrible wrong, we can't have same rowid on diffrent pages")
		}

	}
	panic("should never enter here")
}

type BtreeStruct struct {
	softPages map[int]PageParsed
}

func (btree *BtreeStruct) softWrite(page PageParsed) {
	btree.softPages[page.pageNumber] = page
}

func (btree *BtreeStruct) balancingForNode(node *PageParsed, parents []*PageParsed, header *DbHeader) {

	siblings := []PageParsed{}

	isRoot := len(parents) == 0

	if !node.isOverflow {
		return
	}
	fmt.Println("is overflow??", node.pageNumber)
	fmt.Println("is overflow??", node.pageNumber)
	fmt.Println("is overflow??", node.pageNumber)

	var parent PageParsed

	if isRoot && node.isOverflow {

		cellArea := []byte{}
		cell := parseCellArea(node.cellAreaParsed[0], BtreeType(node.btreeType))

		cellArea = append(cellArea, intToBinary(header.dbSizeInPages, 4)...)
		cellArea = append(cellArea, intToBinary(cell.rowId, 2)...)
		parsedCellArea := append([][]byte{}, cellArea)

		var headerForRootPage *DbHeader
		if node.pageNumber == 0 {
			headerForRootPage = &node.dbHeader
		}

		rootPage := CreateNewPage(TableBtreeInteriorCell, parsedCellArea, node.pageNumber, headerForRootPage)

		node.dbHeaderSize = 0
		node.dbHeader = DbHeader{}
		parents = append(parents, &rootPage)
		node.pageNumber = header.dbSizeInPages
		if node.isSpace() {
			node.isOverflow = false
		}

		fmt.Println("new root??")

		header.dbSizeInPages++

		if node.pageNumber == 2 {
			fmt.Println("save page 2")
			fmt.Printf("%+v\n", node.cellAreaParsed)
		}

		if rootPage.pageNumber == 2 {
			fmt.Println("save page 2")
			fmt.Printf("%+v\n", node.cellAreaParsed)
		}

		btree.softWrite(*node)
		btree.softWrite(rootPage)
	}

	usableSpacePerPage = PageSize - node.btreePageHeaderSize - node.dbHeaderSize
	fmt.Println("calculate usable")
	fmt.Println("btree page btree header size", node.btreePageHeaderSize)
	fmt.Println("btree page header size", node.dbHeaderSize)
	fmt.Println("btree page pointers?", node.pointers)
	fmt.Println("calculate usable")

	if len(parents) > 0 {

		parent = *parents[len(parents)-1]
		parents = parents[:len(parents)-1]
	}

	leftSibling, rightSibling := parent.findSiblings(*node)

	if leftSibling != nil {
		siblings = append(siblings, *leftSibling)
	}

	siblings = append(siblings, *node)

	if rightSibling != nil {
		siblings = append(siblings, *rightSibling)
	}

	cellToDistribute := []Cell{}
	startIndex := PageSize
	endIndex := 0

	for i, v := range siblings {
		// need to start from last, because first item its saved at the end of page
		for j, _ := range v.cellAreaParsed {
			index := len(v.cellAreaParsed) - 1 - j
			cell := parseCellArea(v.cellAreaParsed[index], BtreeType(v.btreeType))
			// +2 for pointers
			cellToDistribute = append(cellToDistribute, Cell{size: len(v.cellAreaParsed[index]) + 2, pageNumber: cell.pageNumber, rowId: cell.rowId, data: v.cellAreaParsed[index]})
		}

		if i < len(siblings) {
			_, newStartIndex, newEndIndex := parent.getDivider(v.pageNumber)
			if startIndex > newStartIndex {
				startIndex = newStartIndex
			}
			if newEndIndex > endIndex {
				endIndex = newEndIndex
			}

		}
	}

	totalSizeInEachPage, numberOfCellPerPage := leaf_bias(cellToDistribute)
	fmt.Println("number of cell per page?")
	fmt.Println(numberOfCellPerPage)

	totalSizeInEachPage, numberOfCellPerPage = accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute, *node)

	fmt.Println("number of cell per page?")
	fmt.Println(numberOfCellPerPage)

	for len(siblings) < len(numberOfCellPerPage) {
		fmt.Println("new page??")
		newPage := CreateNewPage(BtreeType(node.btreeType), [][]byte{}, header.assignNewPage(), nil)
		siblings = append(siblings, newPage)
	}
	for len(numberOfCellPerPage) < len(siblings) {
		siblings = siblings[:len(siblings)-1]
	}
	deivider, pages := redistribution(numberOfCellPerPage, cellToDistribute, siblings)

	for i, v := range pages {
		rowData := []byte{}
		//again start from the last
		// fmt.Println("what are we saving here?")
		// fmt.Println(siblings[i].pageNumber)
		for j, _ := range v {
			index := len(v) - 1 - j

			rowData = append(rowData, v[index].data...)

		}

		fmt.Println("data to save???", v)
		siblings[i].updateCells(dbReadparseCellArea(byte(siblings[i].btreeType), rowData))

		btree.softWrite(siblings[i])
		fmt.Println("save page sibling???")
		fmt.Println("save page sibling???")
		fmt.Println("save page sibling???", siblings[i].pageNumber)
		fmt.Printf("%+v\n", node.cellAreaParsed)
		fmt.Println("save page sibling???")
		fmt.Println("save page sibling???")
	}

	fmt.Println("parents??", parents)
	fmt.Println("parents??", deivider)

	modifyDivider(&parent, deivider, startIndex, endIndex, header, parents)

	fmt.Println("parent??")
	fmt.Printf("%+v", parent)

	for _, v := range parents {
		fmt.Println("save page?? in balancing?", v.pageNumber)
		// if(v.isDirty) {
		btree.softWrite(*v)
		// }
	}
	fmt.Println("save page?? in balancing", parent.pageNumber)
	if parent.pageNumber == 2 {
		fmt.Println("save page 2")
		fmt.Println(parent.cellAreaParsed)
	}
	btree.softWrite(parent)

	btree.balancingForNode(&parent, parents, header)

}

func leaf_bias(cells []Cell) ([]int, []int) {

	totalSizeInEachPage := []int{0}
	numberOfCellPerPage := []int{0}

	fmt.Println("show data leaf bias")
	fmt.Printf("%+v", cells)
	fmt.Println(usableSpacePerPage)
	fmt.Println("show data leaf bias")

	for _, v := range cells {
		i := len(totalSizeInEachPage) - 1

		if totalSizeInEachPage[i]+v.size <= usableSpacePerPage {
			totalSizeInEachPage[i] += v.size
			numberOfCellPerPage[i]++
		} else {
			totalSizeInEachPage = append(totalSizeInEachPage, v.size)
			numberOfCellPerPage = append(numberOfCellPerPage, 1)
		}
	}

	return totalSizeInEachPage, numberOfCellPerPage
}

func accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage []int, cellToDistribute []Cell, node PageParsed) ([]int, []int) {
	divCell := len(cellToDistribute) - numberOfCellPerPage[len(numberOfCellPerPage)-1] - 1

	if len(numberOfCellPerPage) >= 2 {

		for i := len(totalSizeInEachPage) - 1; i > 0; i-- {
			for totalSizeInEachPage[i] <= ((usableSpacePerPage / 2) - cellToDistribute[0].size) {
				totalSizeInEachPage[i] += cellToDistribute[divCell].size
				numberOfCellPerPage[i]++

				numberOfCellPerPage[i-1]--
				if node.isLeaf {
					totalSizeInEachPage[i-1] -= cellToDistribute[divCell].size
				} else {
					totalSizeInEachPage[i-1] -= cellToDistribute[divCell-1].size
				}

				divCell--
			}
		}

		if totalSizeInEachPage[0] < usableSpacePerPage/2 {
			numberOfCellPerPage[0] += 1
			numberOfCellPerPage[1] -= 1
		}
	}

	return totalSizeInEachPage, numberOfCellPerPage
}

// TODO:
// leaf this redistrubtion
// work on saving

func redistribution(numberOfCellPerPage []int, cellToDistribute []Cell, siblingsLength []PageParsed) ([]Cell, [][]Cell) {
	dividers := []Cell{}
	pages := make([][]Cell, len(numberOfCellPerPage))
	pageIndex := -1
	cellIndex := -1
	for i, v := range numberOfCellPerPage {
		pageIndex++
		for range v {
			if cellIndex < len(cellToDistribute) {
				cellIndex++
				pages[pageIndex] = append(pages[pageIndex], cellToDistribute[cellIndex])

			}
		}

		if i < len(siblingsLength) {

			if cellIndex >= len(cellToDistribute) {
				panic("should never occur, cell index >= than distribute")
			}
			dividers = append(dividers, Cell{
				rowId:      cellToDistribute[cellIndex].rowId,
				pageNumber: siblingsLength[i].pageNumber,
				data:       cellToDistribute[cellIndex].data,
				size:       cellToDistribute[cellIndex].size,
			})

		}
	}

	return dividers, pages
}

// inserts a new key into key updates the value if keys already exists

// entries are always inserted at leaft nodes. Internal nodes and the root can only grow in size when leaft nodes overflow and siblings cant take any load to kee pthe leaves balanced, causing a split

// lets walk through an exmaple
// suppose we have the following tree of orde=r, which means each node hold at maxiumum 3 keys and 4 children

/// ```text
///                             PAGE 0 (ROOT)
///                              +-------+
///                          +---|  3,6  |---+
///                         /    +-------+    \
///                        /         |         \
///                   +-------+  +-------+  +-------+
///                   |  1,2  |  |  4,5  |  |  7,8  |
///                   +-------+  +-------+  +-------+
///                     PAGE 1     PAGE 2     PAGE 3
/// ```

// lets try insert key 9, the inserion alogirthm will call find to the page and index, when new key should be added, it will simply insert the key
/// ```text
///                             PAGE 0 (ROOT)
///                              +-------+
///                          +---|  3,6  |---+
///                         /    +-------+    \
///                        /         |         \
///                   +-------+  +-------+  +---------+
///                   |  1,2  |  |  4,5  |  |  7,8,9  |
///                   +-------+  +-------+  +---------+
///                     PAGE 1     PAGE 2     PAGE 3
/// ```

// page 3 has maxiumum number of key if we try to insert now key number 10

///
/// ```text
///                             PAGE 0 (ROOT)
///                              +-------+
///                          +---|  3,6  |---+
///                         /    +-------+    \
///                        /         |         \
///                   +-------+  +-------+  +----------+
///                   |  1,2  |  |  4,5  |  | 7,8,9,10 |
///                   +-------+  +-------+  +----------+
///                     PAGE 1     PAGE 2      PAGE 3
/// ```
/// at the end we need to call balcning algorithm

// Test insert function!!!

// cell
// insert a new value, where there is a rowid??

func (page *PageParsed) isSpace() bool {
	return (page.cellAreaSize + page.btreePageHeaderSize + page.dbHeaderSize + len(page.pointers)) <= PageSize
}

func (page *PageParsed) insertData(data CreateCell, header *DbHeader, parents []*PageParsed) {

	if len(parents) > 0 {
		parent := parents[len(parents)-1]
		parents = parents[:len(parents)-1]
		divider, startIndex, endIndex := parent.getDivider(page.pageNumber)
		if divider.rowId == data.rowId {
			panic("that shouldn't never happen, we can't insert exisiting id")
		}
		if divider.rowId < data.rowId && parent != nil {
			cell := Cell{
				pageNumber: page.pageNumber,
				rowId:      data.rowId,
			}

			modifyDivider(parent, []Cell{cell}, startIndex, endIndex, header, parents)
		} else {
			panic("should never happen, insert data")
		}
	}

	cellParsedData := [][]byte{data.data}
	cellParsedData = append(cellParsedData, page.cellAreaParsed...)

	page.updateCells(cellParsedData)

}

func (page *PageParsed) updateParsedCells(data CreateCell, index int) {
	newParsedCells := make([][]byte, len(page.cellAreaParsed))
	for i := range page.cellAreaParsed {
		copy(newParsedCells[i], page.cellAreaParsed[i])
	}
	newParsedCells[index] = data.data

	page.updateCells(newParsedCells)
}

func (btree *BtreeStruct) insert(rowId int, cell CreateCell, header *DbHeader, startPageNumber *int) PageParsed {

	start := 0
	if startPageNumber != nil {
		start = *startPageNumber
	}

	ok, index, node, parents := search(start, rowId, []*PageParsed{})

	if ok {
		node.updateParsedCells(cell, index)
		return node
	} else {
		node.insertData(cell, header, parents)
	}

	btree.softWrite(node)
	fmt.Println("parents?```", parents)
	allPages := []*PageParsed{&node}
	allPages = append(allPages, parents...)
	for _, v := range allPages {
		// if v.isDirty {
		fmt.Println("save page??", v.pageNumber)
		fmt.Println("save page?/", v.pageNumber)
		fmt.Printf("%+v \n", v.cellAreaParsed)
		btree.softPages[v.pageNumber] = *v
		// }
	}

	btree.balancingForNode(&node, parents, header)
	fmt.Println("save page?/", node.pageNumber)

	// btree.softPages[node.pageNumber] = node

	return node
}

// Search algorithm
// 1. Read the subtree node into memory
// 2. Run a binary search on the entries to find the given key.
// 3. IF successful, return the result
// 4. If not, the binary search result will tell us which child to pick for

//  find key 9 in this tree, located at page 5

///
/// ```text
///                             PAGE 0
///                           +--------+
///                   +-------|   11   |-------+
///                  /        +--------+        \
///                 /                            \
///            +-------+ PAGE 1              +--------+
///       +----|  4,8  |----+                |   14   |
///      /     +-------+     \               +--------+
///     /          |          \               /      \
/// +-------+  +-------+  +-------+     +-------+  +-------+
/// | 1,2,3 |  | 5,6,7 |  | 9,10  |     | 12,13 |  | 15,16 |
/// +-------+  +-------+  +-------+     +-------+  +-------+
///  PAGE 3     PAGE 4      PAGE 5

// first iteration
// 1. read page 0 into memory
// 2. binary search on page result in err(not found)
// 3. read index 0 using page.child and recurse into the result

// second iteraion
// 1. read page 1 into memroy
// 2. binary search result in err
// 3. read child pointer at index 3, and recurse again

// final iteration
// 1. read page 5 into memory
// 2. binary search reuslt it ok
// 3. done, return result

func search(pageNumber int, entry int, parents []*PageParsed) (bool, int, PageParsed, []*PageParsed) {
	reader := NewReader("")
	page := reader.readFromMemory(pageNumber)

	ok, newPageNumber, index := binarySearch(page, pageNumber, entry)

	if !ok && page.isLeaf {
		return false, index, page, parents
	}
	if ok {
		return ok, index, page, parents
	}
	parents = append(parents, &page)
	return search(newPageNumber, entry, parents)
}

// WE add new entry?
// we add it to last height to 0x0d tree
// then we need to balance tree
// let say pointer has change, how we approach it?
// we redistribut load as we start
///                                +-------+
///                        +-------| 4,8,12 |-------+                           X
///                      /        +-------+        \
///                      /             |             \
///   			  +----------+  +----------+  +----------+
///   			  | 1,2,3    |  | 5, 6, 7  |  | 9, 10, 11|
/// 			  +----------+  +----------+  +----------+

///                        +-------| 4,7,10, .... |-------+    -------            -------            X
///                      /        +----------------+ \			\					\
///                      /             |             \			\					\
///   			  +----------+  +----------+  +----------+     +---------+			........
///   			  | 1,2,3    |  | 5,6      |  | 8,9      |     |  11,12   |
/// 			  +----------+  +----------+  +----------+     +---------+
// last height are leaft page that stores only value, so on the 0x05 iterior page we change from 4,8,12 to 4,7, 12
// WE need to have somehow index, or space or something wheer to replace this pointers
// we should remove all three 4,8,12 annd the paste 4,7,12, (could be more items)
// These pointer are basically cellArea, but we do have a index, so we can use it for now
// separate inserting from updating pointers
// pointers requires rowId, these value we move around should have rowid
//

func parseCellArea(data []byte, btreeType BtreeType) Cell {
	if len(data) == 0 {
		panic("data is empty")
	}
	var rowId int
	var pageNumber int
	if btreeType == TableBtreeInteriorCell {
		//TODO:  i think we need to fix it
		pageNumberInt64 := DecodeVarint(data[:4])
		pageNumber = int(pageNumberInt64)
		rowId = int(DecodeVarint(data[4:6]))
	} else if btreeType == TableBtreeLeafCell {
		if data[0] > 127 || data[1] > 127 {
			panic("implement this")
		}
		rowId = int(data[1])
	} else {
		panic("tree don't implemented")
	}
	return Cell{
		rowId:      rowId,
		size:       len(data),
		data:       data,
		pageNumber: pageNumber,
	}
}

func binarySearch(page PageParsed, pageNumber int, rowIdAsEntry int) (bool, int, int) {

	if len(page.cellAreaParsed) == 0 && page.btreeType == int(TableBtreeLeafCell) {
		fmt.Println("pointer 0.1")
		return false, pageNumber, 0
	} else if len(page.cellAreaParsed) == 0 {
		panic("should never occur, i guess")
	}

	cellParsed := parseCellArea(page.cellAreaParsed[0], BtreeType(page.btreeType))

	rightRowId := cellParsed.rowId

	rightPointerPage := int(DecodeVarint(page.rightMostpointer))

	if page.btreeType == int(TableBtreeInteriorCell) {
		if rowIdAsEntry > rightRowId && rightPointerPage != 0 {
			return false, rightPointerPage, 0
		}
	} else if page.btreeType == int(TableBtreeLeafCell) {
		if rowIdAsEntry > rightRowId {
			return false, page.pageNumber, 0
		}
	}

	fmt.Println("pointer 2")

	var leftRowId int

	for i := 0; i < len(page.cellAreaParsed); i++ {
		if i < len(page.cellAreaParsed)-1 {
			cellParsed := parseCellArea(page.cellAreaParsed[i+1], BtreeType(page.btreeType))
			leftRowId = cellParsed.rowId
		}

		v := page.cellAreaParsed[i]

		cellParsed := parseCellArea(v, BtreeType(page.btreeType))
		rowId := cellParsed.rowId

		if page.isLeaf && rowId == rowIdAsEntry {
			return true, page.pageNumber, i
		}

		// we are on last page
		if i == len(page.cellAreaParsed)-1 && rowIdAsEntry <= rowId {
			// go to current page (as this is last page)
			return false, cellParsed.pageNumber, i
		} else if i == len(page.cellAreaParsed)-1 {
			panic("should never happen, binary search")
		}

		// entry row id is smaller than current row id but grater than elft one, we need to go to the page
		if rowIdAsEntry <= rowId && rowIdAsEntry > leftRowId {
			// go to current page
			return false, cellParsed.pageNumber, i

		}

	}

	panic("should never occur, binary search")
}

package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
)

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

type Divider struct {
	page  int
	rowid int
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
	cellAreaParsed       [][]byte
	latesRow             *LastPageParseLatestRow
	isOverflow           bool
	leftSibling          *int
	rightSiblisng        *int
	divider              []Divider
	pageNumber           int
	isLeaf               bool
}

func getDivider(pageNumber int) (Divider, int, int) {
	reader := NewReader("")
	pageRaw := reader.readDbPage(pageNumber)
	page := parseReadPage(pageRaw, pageNumber)
	// page := getNode(pageNumber)

	if page.btreeType != int(TableBtreeInteriorCell) {
		panic("onyl divider for interior cell tree")
	}
	i := 0
	for i < len(page.cellArea) {
		pointer := page.cellArea[i : i+6]
		pointerPageNumber := binary.BigEndian.Uint32(pointer[:4])

		if int(pointerPageNumber) == pageNumber {
			return Divider{
				page:  pageNumber,
				rowid: int(binary.BigEndian.Uint16(pointer[4:])),
			}, i, i + 6
		}

		i += 6
	}

	panic("get, didnt find divider")
}

// focus on this test it etc

func updateDivider(pageNumber int, cells []Cell, startIndex, endIndex int, firstPage *PageParsed) {
	reader := NewReader("")
	pageRaw := reader.readDbPage(pageNumber)
	page := parseReadPage(pageRaw, pageNumber)
	// page := getNode(pageNumber)

	if page.btreeType != int(TableBtreeInteriorCell) {
		panic("onyl divider for interior cell tree")
	}
	contentAreaFirst := page.cellArea[:startIndex]
	contentAreaSecond := page.cellArea[:endIndex]
	for _, v := range cells {
		newPointer := intToBinary(v.pageNumber, 4)
		newPointer = append(newPointer, intToBinary(v.rowId, 2)...)
		contentAreaFirst = append(contentAreaFirst, newPointer...)
	}
	contentAreaFirst = append(contentAreaFirst, contentAreaSecond...)
	if len(contentAreaFirst) > 12 {
		page.isOverflow = true
	}
	page.cellArea = contentAreaFirst
	page.cellAreaParsed = [][]byte{}

	for len(contentAreaFirst) > 0 {
		page.cellAreaParsed = append(page.cellAreaParsed, contentAreaFirst[:6])
		contentAreaFirst = contentAreaFirst[6:]
	}
	page.startCellContentArea = PageSize - len(page.cellArea)
	writer := NewWriter()

	writer.softwiteToFile(page, pageNumber, firstPage)

	if pageNumber == 0 {
		firstPage = &page
	}
	// check if cell area overflow page
	// writer.writeToFile(assembleDbPage(page), pageNumber, "", firstPage)
	// saveNode(pageNumber, page)

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

func parseDbHeader(data []byte) DbHeader {
	if len(data) != 100 {
		panic("header should have 100 bytes")
	}
	headerString := data[:16]
	databasePageSize := data[16:18]
	databaseFileWriteVersion := data[18:19]
	databaseFileReadVersion := data[19:20]
	reservedBytesSpace := data[20:21]
	maxEmbeddedPayloadFraction := data[21:22]
	minEmbeddedPayloadFraction := data[22:23]
	leafPayloadFraction := data[23:24]
	fileChangeCounter := data[24:28]
	dbSizeInPages := data[28:32]
	firstFreeListTrunkPage := data[32:36]
	totalNumberOfFreeListPages := data[36:40]
	schemaCookie := data[40:44]
	schemaFormatNumber := data[44:48]
	defaultPageCacheSize := data[48:52]
	largestBTreePage := data[52:56]
	databaseEncoding := data[56:60]
	userVersion := data[60:64]
	incrementalVacuumMode := data[64:68]
	applicationId := data[68:72]
	reservedForExpansion := data[72:92]
	versionValidForNumber := data[92:96]
	sqlVersionNumber := data[96:100]

	return DbHeader{
		headerString:               headerString,
		databasePageSize:           int(binary.BigEndian.Uint16(databasePageSize)),
		databaseFileWriteVersion:   databaseFileWriteVersion,
		databaseFileReadVersion:    databaseFileReadVersion,
		reservedBytesSpace:         reservedBytesSpace,
		maxEmbeddedPayloadFraction: maxEmbeddedPayloadFraction,
		minEmbeddedPayloadFraction: minEmbeddedPayloadFraction,
		leafPayloadFraction:        leafPayloadFraction,
		fileChangeCounter:          int(binary.BigEndian.Uint32(fileChangeCounter)),
		dbSizeInPages:              int(binary.BigEndian.Uint32(dbSizeInPages)),
		firstFreeListTrunkPage:     firstFreeListTrunkPage,
		totalNumberOfFreeListPages: totalNumberOfFreeListPages,
		schemaCookie:               int(binary.BigEndian.Uint32(schemaCookie)),
		schemaFormatNumber:         schemaFormatNumber,
		defaultPageCacheSize:       defaultPageCacheSize,
		largestBTreePage:           largestBTreePage,
		databaseEncoding:           databaseEncoding,
		userVersion:                userVersion,
		incrementalVacuumMode:      incrementalVacuumMode,
		applicationId:              applicationId,
		reservedForExpansion:       reservedForExpansion,
		versionValidForNumber:      int(binary.BigEndian.Uint32(versionValidForNumber)),
		sqlVersionNumber:           sqlVersionNumber,
	}
}

func parseDbPageColumn(rowHeader []byte, rowValues []byte) []PageParseColumn {
	var rowColumn []PageParseColumn
	for _, v := range rowHeader {
		if int(v) > 127 {
			panic("handle case that we have multiple bytes")
		}
		if int(v) == 0 {
			column := PageParseColumn{
				columnType:   string(strconv.Itoa(int(v))),
				columnLength: 0,
				columnValue:  []byte{0},
			}
			rowColumn = append(rowColumn, column)
			continue

		}
		if int(v) < 10 {
			//int
			column := PageParseColumn{
				columnType:   string(strconv.Itoa(int(v))),
				columnLength: 1,
				columnValue:  []byte{rowValues[0]},
			}
			rowColumn = append(rowColumn, column)
			rowValues = rowValues[1:]
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

			if length > len(rowValues) {
				panic("there is not enough data")
			}
			value := rowValues[:length]
			column := PageParseColumn{
				columnType:   "13",
				columnLength: length,
				columnValue:  value,
			}
			rowColumn = append(rowColumn, column)
			rowValues = rowValues[length:]
			continue
		}
		panic("should never enter this state in parsing")

	}
	return rowColumn
}

func parseReadPage(data []byte, dbPage int) PageParsed {
	if dbPage == 0 && len(data) == 0 {
		fmt.Println("here??")
		return PageParsed{
			dbHeader:             header(),
			dbHeaderSize:         100,
			numberofCells:        0,
			startCellContentArea: PageSize,
			cellArea:             []byte{},
			pointers:             []byte{},
			latesRow: &LastPageParseLatestRow{
				rowId:   0,
				data:    []byte{},
				columns: []PageParseColumn{},
			},
		}
	}

	if len(data) != PageSize {
		panic("invalid page size, expected" + strconv.Itoa(PageSize) + " got: " + strconv.Itoa(len(data)))
	}
	dataToParse := data
	var dbHeader DbHeader
	var dbHeaderSize int = 0
	if dbPage == 0 {
		//Skip header for now
		if !reflect.DeepEqual(dataToParse[:16], []byte("SQLite format 3\000")) {
			fmt.Println(dataToParse[:16])
			fmt.Println(string(dataToParse[:16]))
			panic("header on page 0 should start with sqlite format....")
		}
		dataToParse = dataToParse[100:]
		dbHeader = parseDbHeader(data[:100])
		dbHeaderSize = 100
	}

	btreeType := dataToParse[0]
	isPointerValue := btreeType == 0x05

	if btreeType != 0x05 && btreeType != 0x0d {
		panic("implement this btree types")
	}
	isLeaf := false
	if btreeType == byte(TableBtreeLeafCell) {
		isLeaf = true
	}
	switch BtreeType(btreeType) {
	case TableBtreeInteriorCell, IndexBtreeInteriorCell:
		isPointerValue = true

	}

	freebBlocks := dataToParse[1:3]
	if freebBlocks[0] != 0 {
		panic("implement free blocks more than 0 cell")
	}
	freeBlocksInt := int(freebBlocks[1])
	numberofCells := dataToParse[3:5]
	if numberofCells[0] != 0 {
		panic("implement number of cell more than 0 cell")
	}

	numberofCellsInt := int(numberofCells[1])
	startCellContentArea := dataToParse[5:7]

	startCellContentAreaInt := binary.BigEndian.Uint16(startCellContentArea)
	startCellContentAreaBigEndian := binary.BigEndian.Uint16(startCellContentArea)
	fragmenetedArea := dataToParse[7]
	var rightMostPointer []byte = []byte{}
	if isPointerValue {
		fmt.Println("is right pointer")
		rightMostPointer = dataToParse[8:12]
		dataToParse = dataToParse[12:]
	} else {
		dataToParse = dataToParse[8:]
	}

	var pointers []byte

	for i := 0; i < numberofCellsInt; i++ {
		if len(dataToParse) < 2 {
			panic("should never happend this")
		}
		// fmt.Println("pointers???")
		pointer := dataToParse[:2]
		if pointer[0] == 0 && pointer[1] == 0 {
			break
		}
		dataToParse = dataToParse[2:]
		pointers = append(pointers, pointer...)
	}
	fmt.Println("after pointerS???")
	if len(data) < int(startCellContentAreaBigEndian) {
		panic("data length is lesser than start of cell content area")
	}
	cellAreaContent := []byte{}
	if startCellContentAreaBigEndian != 0 {
		// 0 means PageSize, no data
		cellAreaContent = data[startCellContentAreaBigEndian:]
	}

	latestRowHeaders := []byte{}
	latestRowValues := []byte{}
	var latestRow LastPageParseLatestRow

	cellAreaContentTmp := cellAreaContent
	cellAreaParsed := [][]byte{}

	for len(cellAreaContentTmp) > 0 {
		cellAreaParsed = append(cellAreaParsed, cellAreaContentTmp[:6])
		cellAreaContentTmp = cellAreaContentTmp[6:]
	}

	if len(cellAreaContent) > 0 && btreeType == 0x13 {
		latestRowLength := int(cellAreaContent[0]) + 2
		var latestRowLengthArr []byte
		for i := 0; i < latestRowLength; i++ {
			latestRowLengthArr = append(latestRowLengthArr, cellAreaContent[i])
			if cellAreaContent[i] < 127 {
				break
			}
		}
		if len(latestRowLengthArr) > 1 {
			panic("Need to be handled later")
		}

		// latestRowLength := int(latestRowLengthArr[0]) + 2 // 1 bytes for length, 1 bytes for row id
		if len(cellAreaContent) < int(latestRowLength) {
			panic("cellAreaContent length is lesser than start of cell content area, row length%")
		}
		latestRowRaw := cellAreaContent[:latestRowLength]
		latestRowId := latestRowRaw[1]
		fmt.Println("hello here??")

		latestRowheaderLength := latestRowRaw[2]
		fmt.Println("hello here2??")
		latestRowHeaders = latestRowRaw[3 : 3-1+int(latestRowheaderLength)] // 3 - 1 (-1 because of header length contains itself)
		latestRowValues = latestRowRaw[3-1+int(latestRowheaderLength):]
		latestRowColumns := parseDbPageColumn(latestRowHeaders, latestRowValues)
		latestRow = LastPageParseLatestRow{
			rowId:   int(latestRowId),
			data:    latestRowRaw,
			columns: latestRowColumns,
		}
	}
	isOverflow := false
	if dbPage == 1 {
		isOverflow = true
	}

	return PageParsed{
		dbHeader:             dbHeader,
		dbHeaderSize:         dbHeaderSize,
		pageNumber:           dbPage,
		isLeaf:               isLeaf,
		btreeType:            int(btreeType),
		numberofCells:        numberofCellsInt,
		startCellContentArea: int(startCellContentAreaInt),
		rightMostpointer:     rightMostPointer,
		cellArea:             cellAreaContent,
		isOverflow:           isOverflow,
		cellAreaParsed:       cellAreaParsed,
		framgenetedArea:      int(fragmenetedArea),
		freeBlock:            int(freeBlocksInt),
		pointers:             pointers,
		latesRow:             &latestRow,
	}
}

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
	rowId int
}

type PageParsed struct {
	dbHeader             DbHeader // only for first page
	dbHeaderSize         int
	btreePageHeaderSize  int
	btreeType            int
	freeBlock            int
	numberofCells        int
	startCellContentArea int
	framgenetedArea      int
	rightMostpointer     []byte
	pointers             []byte
	cellArea             []byte
	cellAreaParsed       [][]byte
	cellAreaSize         int
	latesRow             *LastPageParseLatestRow
	isOverflow           bool
	leftSibling          *int
	rightSiblisng        *int
	divider              []Divider
	pageNumber           int
	isLeaf               bool
}

func (parentPage PageParsed) getDivider(pageNumber int) (Divider, int, int) {

	if parentPage.pageNumber == 0 && parentPage.numberofCells == 0 {
		fmt.Println("zero condition???!!!")
		fmt.Println(parentPage.cellArea)
		fmt.Println(parentPage.pageNumber)
		fmt.Println(parentPage.numberofCells)
		fmt.Println("zero condition???!!!")
		fmt.Println("zero condition???!!!")
		return Divider{}, 0, 0
	}

	// page := getNode(pageNumber)

	cellAreaTmp := parentPage.cellArea
	// cellAreaTmp = append(cellAreaTmp, []byte{0, 0}...)
	// cellAreaTmp = append(cellAreaTmp, parentPage.cellArea...)

	if parentPage.btreeType != int(TableBtreeInteriorCell) {
		fmt.Println("page number", pageNumber)
		panic("onyl divider for interior cell tree, get divider")
	}
	i := 0
	fmt.Println("get divider for page number??", pageNumber)
	for i < len(cellAreaTmp) {
		pointer := parentPage.cellArea[i : i+6]
		pointerPageNumber := binary.BigEndian.Uint32(pointer[:4])
		rowId := binary.BigEndian.Uint16(pointer[4:6])

		if int(pointerPageNumber) == pageNumber {
			return Divider{
				page:  pageNumber,
				rowId: int(rowId),
				// rowid: int(binary.BigEndian.Uint16(pointer[4:])),
			}, i, i + 6
		}

		i += 6
	}

	panic("get, didnt find divider")
}

// i think we only need to update the higest rowId as we references this
func updateDivider(page *PageParsed, parents []*PageParsed, mostRightRowId int, pageNumber int, firstPage *PageParsed) {
	fmt.Println("hello update divider???", page.pageNumber)
	fmt.Println("looking for page?", pageNumber)
	for i, v := range page.cellAreaParsed {
		pageNumberCell := binary.BigEndian.Uint32(v[:4])
		rowId := binary.BigEndian.Uint16(v[4:6])
		if pageNumberCell == uint32(pageNumber) && mostRightRowId == int(rowId) {
			return
		} else if pageNumberCell == uint32(pageNumber) {
			newCellArea := intToBinary(pageNumber, 4)
			newCellArea = append(newCellArea, intToBinary(mostRightRowId, 2)...)
			page.cellAreaParsed[i] = newCellArea
			break
		} else if i == len(page.cellAreaParsed)-1 {
			panic("last page, didn't find what we looking for")
		}
	}

	cellArea := []byte{}
	pointers := []byte{}
	lastPointer := PageSize

	for _, v := range page.cellAreaParsed {
		cellArea = append(cellArea, v...)
		lastPointer -= len(v)
		pointers = append(pointers, intToBinary(lastPointer, 2)...)
	}

	page.pointers = pointers
	page.cellArea = cellArea

	writer := NewWriter()
	writer.softwiteToFile(page, page.pageNumber, firstPage)

	if len(parents) == 0 {
		return
	}

	if len(page.cellAreaParsed) == 0 {
		panic("cell area parsed 0 in update dividr")
	}

	lastCell := page.cellAreaParsed[0]
	// mostRightPageNumber := int(binary.BigEndian.Uint16(lastCell[:4]))
	mostRightRowId = int(binary.BigEndian.Uint16(lastCell[4:6]))

	parent := parents[len(parents)-1]
	parents = parents[:len(parents)-1]

	updateDivider(parent, parents, mostRightRowId, page.pageNumber, firstPage)

}

// we need update recursively parent too

func modifyDivider(page *PageParsed, cells []Cell, startIndex, endIndex int, firstPage *PageParsed, parents []*PageParsed) {

	fmt.Println("what divider are we adding?")
	fmt.Println(cells)
	fmt.Println("at position", startIndex, endIndex)

	if page.btreeType != int(TableBtreeInteriorCell) {
		panic("onyl divider for interior cell tree, update divider")
	}
	if len(cells) == 0 {
		panic("cells are empty")
	}
	cellArea := page.cellArea

	contentAreaFirst := append([]byte{}, cellArea[:startIndex]...)
	contentAreaSecond := append([]byte{}, cellArea[endIndex:]...)

	for i := len(cells) - 1; i >= 0; i-- {
		newPointer := intToBinary(cells[i].pageNumber, 4)
		newPointer = append(newPointer, intToBinary(cells[i].rowId, 2)...)
		contentAreaFirst = append(contentAreaFirst, newPointer...)
	}

	contentAreaFirst = append(contentAreaFirst, contentAreaSecond...)
	if len(contentAreaFirst) < 6 {
		fmt.Println(contentAreaFirst)
		panic("content area is small")
	}

	page.cellArea = contentAreaFirst
	cellAreaParsed := dbReadparseCellArea(byte(page.btreeType), contentAreaFirst)
	pointers := []byte{}
	lastPointer := PageSize
	for _, v := range cellAreaParsed {
		lastPointer -= len(v)
		pointers = append(pointers, intToBinary(lastPointer, 2)...)
	}
	fmt.Println("cell area parsed??", cellAreaParsed)
	fmt.Println("right most page??", cellAreaParsed[0][:4])
	page.rightMostpointer = cellAreaParsed[0][:4]
	page.cellAreaParsed = cellAreaParsed
	page.pointers = pointers

	fmt.Println(page.cellArea)
	fmt.Println(page.cellAreaParsed)

	page.numberofCells = len(cellAreaParsed)
	page.startCellContentArea = PageSize - len(contentAreaFirst)
	page.cellAreaSize = len(contentAreaFirst)

	writer := NewWriter()
	writer.softwiteToFile(page, page.pageNumber, firstPage)

	if page.pageNumber == 0 {
		firstPage = page
	}

	if len(parents) == 0 {
		return
	}

	lastCell := page.cellAreaParsed[0]
	mostRightRowId := int(binary.BigEndian.Uint16(lastCell[4:6]))
	parent := parents[len(parents)-1]
	parents = parents[:len(parents)-1]

	updateDivider(parent, parents, mostRightRowId, page.pageNumber, firstPage)

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

func dbReadparseCellArea(btreeType byte, cellAreaContentTmp []byte) [][]byte {
	cellAreaParsed := [][]byte{}
	for len(cellAreaContentTmp) > 0 {
		if btreeType == byte(TableBtreeInteriorCell) {
			cellAreaParsed = append(cellAreaParsed, cellAreaContentTmp[:6])
			cellAreaContentTmp = cellAreaContentTmp[6:]
		} else if btreeType == byte(TableBtreeLeafCell) {
			if cellAreaContentTmp[0] > 127 {
				//length
				panic("implement this later1")
			}
			if cellAreaContentTmp[1] > 127 {
				panic("implement this later11")
			}
			length := int(cellAreaContentTmp[0]) + 2 // byte for length, and rowid

			cellAreaParsed = append(cellAreaParsed, cellAreaContentTmp[:length])
			cellAreaContentTmp = cellAreaContentTmp[length:]
		}
	}

	return cellAreaParsed
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
			cellAreaSize:         0,
			latesRow: &LastPageParseLatestRow{
				rowId:   0,
				data:    []byte{},
				columns: []PageParseColumn{},
			},
		}
	}

	if len(data) != PageSize {
		fmt.Printf("\n page number: %v", dbPage)
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
	btreePageHeaderSize := 8
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
		btreePageHeaderSize = 12

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

	cellAreaParsed = dbReadparseCellArea(btreeType, cellAreaContentTmp)

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

	return PageParsed{
		dbHeader:             dbHeader,
		dbHeaderSize:         dbHeaderSize,
		btreePageHeaderSize:  btreePageHeaderSize,
		pageNumber:           dbPage,
		isLeaf:               isLeaf,
		btreeType:            int(btreeType),
		numberofCells:        numberofCellsInt,
		startCellContentArea: int(startCellContentAreaInt),
		rightMostpointer:     rightMostPointer,
		cellAreaSize:         len(cellAreaContent),
		cellArea:             cellAreaContent,
		isOverflow:           false,
		cellAreaParsed:       cellAreaParsed,
		framgenetedArea:      int(fragmenetedArea),
		freeBlock:            int(freeBlocksInt),
		pointers:             pointers,
		latesRow:             &latestRow,
	}
}

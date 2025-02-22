package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
)

func TestCreateCellWithSingleValue(t *testing.T) {
	btreeType := TableBtreeLeafCell
	page := PageParsed{
		latesRow: &LastPageParseLatestRow{
			rowId: 1,
		},
	}
	val := "Alice"
	res := createCell(btreeType, &page, val)

	if res.dataLength != 7 {
		t.Errorf("Expect cell data length to be 7 (9 total - 1 bytes for row id  -1 bytes for length byte), got: %v", res.dataLength)
	}

	if res.data[0] != 7 {
		t.Errorf("Expect cell length to be 7 (9 total - 1 bytes for row id  -1 bytes for length byte) got: %v", res.data[0])
	}

	if res.data[1] != 2 {
		t.Errorf("Expect row id to be increment +1 from previous row (1): %v", res.data[1])
	}

	if res.data[2] != 2 {
		t.Errorf("Expect header length to be 2 (1 for alice value, 1 for itself) got: %v", res.data[2])
	}

	if res.data[3] != 23 {
		t.Errorf("Expect column type to be 23 ((23-13)/2 =5 length), 13 is type of text, we got : %v", res.data[3])
	}

	if !reflect.DeepEqual(res.data[4:], []byte(val)) {
		t.Errorf("Expect rest of cell to be value Alice, we got: %v", res.data[4:])
	}
}

func TestCreateCellWithMultipleValues(t *testing.T) {
	btreeType := TableBtreeLeafCell
	page := PageParsed{
		latesRow: &LastPageParseLatestRow{
			rowId: 1,
		},
	}
	val2 := 12     //1
	val3 := "test" //4
	res := createCell(btreeType, &page, nil, val2, val3)

	fmt.Println(res)

	if res.dataLength != 9 {
		t.Errorf("Expect cell data length to be 11 (11 total - 1 bytes for row id  -1 bytes for length byte), got: %v", res.dataLength)
	}

	if res.data[0] != 9 {
		t.Errorf("Expect cell length to be 7 (9 total - 1 bytes for row id  -1 bytes for length byte) got: %v", res.data[0])
	}

	if res.data[1] != 2 {
		t.Errorf("Expect row id to be increment +1 from previous row (1): %v", res.data[1])
	}

	if res.data[2] != 4 {
		t.Errorf("Expect header length to be 2 (1 for alice value, 1 for itself) got: %v", res.data[2])
	}

	if res.data[3] != 0 {
		t.Errorf("Exoected column type for type null should be 0 we got : %v", res.data[3])
	}

	if res.data[4] != 1 {
		t.Errorf("Expect column type for small int to be 1 we got : %v", res.data[4])
	}
	if res.data[5] != 21 {
		t.Errorf("Expect column type for test to be 23 ((23-13)/2 =4 length), 13 is type of text, we got : %v", res.data[5])
	}

	if res.data[6] != 12 {
		t.Errorf("Expect column type to be 23 ((23-13)/2 =4 length), 13 is type of text, we got : %v", res.data[6])
	}

	if !reflect.DeepEqual(res.data[7:], []byte(val3)) {
		t.Errorf("Expect rest of cell to be value Alice, we got: %v", res.data[7:])
	}
}

func TestCreateHeader(t *testing.T) {
	btreeType := TableBtreeLeafCell
	page := PageParsed{
		latesRow: &LastPageParseLatestRow{
			rowId: 1,
		},
	}
	val2 := 12     //1
	val3 := "test" //4
	cell := createCell(btreeType, &page, nil, val2, val3)

	btreePage := updateBtreeHeaderLeafTable(cell, &page)

	fmt.Println("btree page")
	fmt.Println(btreePage)

	if len(btreePage) != 10 {
		t.Errorf("Header should have 10 bytes, we got: %v", len(btreePage))
	}

	if btreePage[0] != 0x0d {
		t.Errorf("Expected header type to be: %v, insted we got: %v", 0x0d, btreePage[0])
	}

	if btreePage[1] != 0 || btreePage[2] != 0 {
		t.Errorf("Expected free block to be: %v, insted we got bytes: %v %v", 0, btreePage[1], btreePage[2])
	}

	if binary.BigEndian.Uint16(btreePage[3:5]) != 1 {
		t.Errorf("Expected number of cell to be: %v, insted we got: %v", 1, binary.BigEndian.Uint16(btreePage[3:5]))
	}

	if binary.BigEndian.Uint16(btreePage[5:7]) != uint16(PageSize-cell.dataLength) {
		t.Errorf("Expected start content area to be: %v, insted we got : %v", uint16(PageSize-cell.dataLength), binary.BigEndian.Uint16(btreePage[5:7]))
	}

	if btreePage[7] != 0 {
		t.Errorf("Expected fragmeneted free bytes to be %v, insted we got : %v", 0, btreePage[7])
	}

	if binary.BigEndian.Uint16(btreePage[8:10]) != uint16(PageSize-cell.dataLength) {
		t.Errorf("Expected cell's pointer to be: %v, insted we got : %v", uint16(PageSize-cell.dataLength), binary.BigEndian.Uint16(btreePage[8:10]))
	}

}

func TestCreateHeaderWithPage(t *testing.T) {
	btreeType := TableBtreeLeafCell
	page := PageParsed{
		latesRow: &LastPageParseLatestRow{
			rowId: 1,
		},
	}

	cellArea := []byte{0, 1, 2, 3, 4, 5}
	parsedData := PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		cellArea:             cellArea,
		startCellContentArea: PageSize - len(cellArea),
		numberofCells:        1,
		pointers:             intToBinary(PageSize-len(cellArea), 2),
	}
	val2 := 12     //1
	val3 := "test" //4
	cell := createCell(btreeType, &page, nil, val2, val3)
	btreeHeader := updateBtreeHeaderLeafTable(cell, &parsedData)

	if len(btreeHeader) != 12 {
		t.Errorf("Header should have 10 bytes, we got: %v", len(btreeHeader))
	}

	if btreeHeader[0] != 0x0d {
		t.Errorf("Expected header type to be: %v, insted we got: %v", 0x0d, btreeHeader[0])
	}

	if btreeHeader[1] != 0 || btreeHeader[2] != 0 {
		t.Errorf("Expected free block to be: %v, insted we got bytes: %v %v", 0, btreeHeader[1], btreeHeader[2])
	}

	if binary.BigEndian.Uint16(btreeHeader[3:5]) != 2 {
		t.Errorf("Expected number of cell to be: %v, insted we got: %v", 2, binary.BigEndian.Uint16(btreeHeader[3:5]))
	}

	if binary.BigEndian.Uint16(btreeHeader[5:7]) != uint16(PageSize-len(cellArea)-cell.dataLength) {
		t.Errorf("Expected start content area to be: %v, insted we got : %v", uint16(PageSize-len(cellArea)-cell.dataLength), binary.BigEndian.Uint16(btreeHeader[5:7]))
	}

	if btreeHeader[7] != 0 {
		t.Errorf("Expected fragmeneted free bytes to be %v, insted we got : %v", 0, btreeHeader[7])
	}

	if !reflect.DeepEqual(btreeHeader[8:10], intToBinary(PageSize-len(cellArea), 2)) {
		t.Errorf("Expected cell's pointer to be: %v, insted we got : %v", PageSize-len(cellArea), binary.BigEndian.Uint16(btreeHeader[8:10]))
	}

	if binary.BigEndian.Uint16(btreeHeader[10:12]) != uint16(PageSize-len(cellArea)-cell.dataLength) {
		t.Errorf("Expected cell's pointer to be: %v, insted we got : %v", uint16(PageSize-cell.dataLength), binary.BigEndian.Uint16(btreeHeader[8:10]))
	}

}

func TestUpdatePage(t *testing.T) {

	parsedData := PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		cellArea:             []byte{},
		startCellContentArea: 0,
		numberofCells:        1,
		pointers:             []byte{},
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		latesRow: &LastPageParseLatestRow{
			rowId: 1,
		},
	}

	firstPage := PageParsed{
		dbHeader: DbHeader{
			dbSizeInPages: 1,
		},
	}

	server := ServerStruct{
		firstPage: firstPage,
	}

	// val2 := 12     //1
	// val3 := "test" //4
	cellAlice := createCell(TableBtreeLeafCell, &parsedData, "Alice")
	cellBob := createCell(TableBtreeLeafCell, &parsedData, "Bob")

	parsedData.cellArea = cellAlice.data
	parsedData.startCellContentArea = PageSize - cellAlice.dataLength
	parsedData.pointers = append(parsedData.pointers, intToBinary(PageSize-cellAlice.dataLength, 2)...)

	reader := NewReader("conId")
	writer := NewWriter()
	updatedPage := server.updatePageRoot(&parsedData, 1, *reader, *writer, "Bob")

	aliceLength := 9
	bobLength := 7

	alicePointer := intToBinary(PageSize-aliceLength, 2)
	BobPointer := intToBinary(PageSize-aliceLength-bobLength, 2)

	if len(updatedPage) != 4096 {
		t.Errorf("Header should have 4097 bytes, we got: %v", len(updatedPage))
	}

	if updatedPage[0] != byte(TableBtreeLeafCell) {
		t.Errorf("Expected btree type to be %v, got: %v", TableBtreeLeafCell, updatedPage[0])
	}

	if !reflect.DeepEqual(updatedPage[1:3], []byte{0, 0}) {
		t.Errorf("Expeected free block to be empty")
	}

	if binary.BigEndian.Uint16(updatedPage[3:5]) != 2 {
		t.Errorf("Expeected data to have two cells")
	}

	expectedStartCellArea := uint16(PageSize) - uint16(aliceLength) - uint16(bobLength)

	if binary.BigEndian.Uint16(updatedPage[5:7]) != expectedStartCellArea {
		t.Errorf("Expeected cell area start to be: %v, got: %v", expectedStartCellArea, binary.BigEndian.Uint16(updatedPage[5:7]))
	}

	if updatedPage[7] != 0 {
		t.Error("fragmenet erra shoul be 0")
	}

	if !reflect.DeepEqual(updatedPage[8:10], alicePointer) {
		t.Errorf("wront pointer to Alice value, expected: %v, got :%v", updatedPage[8:10], alicePointer)
	}

	if !reflect.DeepEqual(updatedPage[10:12], BobPointer) {
		t.Errorf("wront pointer to Bob value, expected: %v, got :%v", updatedPage[10:12], BobPointer)
	}

	emptyBytes := make([]byte, expectedStartCellArea-12)

	if !reflect.DeepEqual(updatedPage[12:expectedStartCellArea], emptyBytes) {
		t.Errorf("expected bytes to be empty")
	}

	expectedCellContent := cellBob.data
	expectedCellContent = append(expectedCellContent, cellAlice.data...)

	if !reflect.DeepEqual(updatedPage[expectedStartCellArea:], expectedCellContent) {
		t.Errorf("Exepected cell content area to be: %v, got: %v", expectedCellContent, updatedPage[expectedStartCellArea:])
	}

}

func TestFindPageToInsertOnePage(t *testing.T) {
	page := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		latesRow:             &LastPageParseLatestRow{},
	}

	firstPage := PageParsed{
		dbHeader: DbHeader{
			dbSizeInPages: 1,
		},
	}

	server := ServerStruct{
		firstPage: firstPage,
	}

	dataAssembled := assembleDbPage(page)

	writer := NewWriter()

	writer.writeToFile(dataAssembled, 0, "fds", &firstPage)

	reader := NewReader("conId")

	_, pageNumber := server.findPageToInsertData(*reader, 0)

	if pageNumber != 0 {
		t.Errorf("Expected page to be 0, got :%v", pageNumber)
	}
}

func TestFindPageToInsertNestedPage(t *testing.T) {
	firstPage := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 1},
		latesRow:             &LastPageParseLatestRow{},
	}

	secondPage := PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		latesRow:             &LastPageParseLatestRow{},
	}

	server := ServerStruct{
		firstPage: firstPage,
	}

	firstPageAssembled := assembleDbPage(firstPage)

	writer := NewWriter()

	writer.writeToFile(firstPageAssembled, 0, "fds", &firstPage)

	secondPageAssembled := assembleDbPage(secondPage)

	writer.writeToFile(secondPageAssembled, 1, "fds", &firstPage)

	reader := NewReader("conId")

	_, pageNumber := server.findPageToInsertData(*reader, 0)

	if pageNumber != 1 {
		t.Errorf("Expected page to be 1, got :%v", pageNumber)
	}
}

// func (server *ServerStruct) createNewLeafPage(pageNumber int) *PageParsed {
// 	btreeHeader := updateBtreeHeaderLeafTable(CreateCell{}, &PageParsed{})

// 	allCells := []byte{}
// 	// Zeros data
// 	zerosLength := PageSize - len(btreeHeader) - len(allCells)
// 	zerosSpace := make([]byte, zerosLength)
// 	//General method to save the daata to disk i guess
// 	dataToSave := []byte{}

// 	dataToSave = append(dataToSave, btreeHeader...)
// 	dataToSave = append(dataToSave, zerosSpace...)
// 	dataToSave = append(dataToSave, allCells...)

// 	parsedData := parseReadPage(dataToSave, pageNumber)
// 	return &parsedData
// }

func TestCreateNewLeafPage(t *testing.T) {
	clearDbFile("test")
	firstPage := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 1},
		latesRow:             &LastPageParseLatestRow{},
	}

	server := ServerStruct{
		firstPage: firstPage,
	}

	leaftPage := server.createNewLeafPage()

	if leaftPage.btreeType != int(TableBtreeLeafCell) {
		t.Error("incorect btree type")
	}

	if len(leaftPage.cellArea) > 0 {
		t.Error("Cell area should be empty")
	}

	if len(leaftPage.rightMostpointer) > 0 {
		t.Error("There should be no right pointer")
	}

	if leaftPage.numberofCells != 0 {
		t.Error("number of cell shoud be 0")
	}

	if len(leaftPage.pointers) > 0 {
		t.Error("there should be no pointer")
	}

	if leaftPage.framgenetedArea != 0 {
		t.Errorf("should be no fragmeneted area")
	}

}

func TestCreateNewinteriorPage(t *testing.T) {
	clearDbFile("test")
	firstPage := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 1},
		latesRow:             &LastPageParseLatestRow{},
	}

	server := ServerStruct{
		firstPage: firstPage,
	}
	rightPointer := []byte{0, 0, 0, 1}

	interiorPage := server.createNewInteriorPage(rightPointer)

	if interiorPage.btreeType != int(TableBtreeInteriorCell) {
		t.Error("incorect btree type")
	}

	if len(interiorPage.cellArea) > 0 {
		t.Error("Cell area should be empty")
	}

	if !reflect.DeepEqual(interiorPage.rightMostpointer, rightPointer) {
		t.Errorf("Pointer should be %v , got: %v", rightPointer, interiorPage.rightMostpointer)
	}

	if interiorPage.numberofCells != 0 {
		t.Error("number of cell shoud be 0")
	}

	if len(interiorPage.pointers) > 0 {
		t.Error("there should be no pointer")
	}

	if interiorPage.framgenetedArea != 0 {
		t.Errorf("should be no fragmeneted area")
	}

}

func TestInsertPointerIntoInteriorPage(t *testing.T) {
	clearDbFile("test")
	page := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 1},
		latesRow:             &LastPageParseLatestRow{},
	}

	server := ServerStruct{
		firstPage: page,
	}
	rightPointer := []byte{0, 0, 0, 2}
	rowId := 2
	expectedCellAreaContent := page.rightMostpointer
	expectedCellAreaContent = append(expectedCellAreaContent, intToBinary(rowId, 2)...)

	err := server.insertPointerIntoInteriorPage(rightPointer, rowId, &page)

	if err != nil {
		t.Errorf("should add pointer without any errors, got %v", err)
	}

	if page.btreeType != int(TableBtreeInteriorCell) {
		t.Error("incorect btree type")
	}

	if !reflect.DeepEqual(expectedCellAreaContent, page.cellArea) {
		t.Errorf("Cell area should be: %v, got: %v", expectedCellAreaContent, page.cellArea)
	}

	if !reflect.DeepEqual(page.rightMostpointer, rightPointer) {
		t.Errorf("Pointer should be %v , got: %v", rightPointer, page.rightMostpointer)
	}

	if page.numberofCells != 1 {
		t.Error("number of cell shoud be 1")
	}

	if len(page.pointers) != 2 {
		t.Error("there should be 1 pointer")
	}

	if page.framgenetedArea != 0 {
		t.Errorf("should be no fragmeneted area")
	}
}

func TestInsertMultiplePointerIntoInteriorPage(t *testing.T) {
	clearDbFile("test")
	page := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 1},
		latesRow:             &LastPageParseLatestRow{},
	}

	server := ServerStruct{
		firstPage: page,
	}
	rightPointers := [][]byte{{0, 0, 0, 2}, {0, 0, 0, 3}, {0, 0, 0, 4}}
	rowIds := []int{5, 15, 20}
	expectedCellAreaContent := page.rightMostpointer
	expectedCellAreaContent = append(expectedCellAreaContent, intToBinary(rowIds[0], 2)...)
	expectedPointers := []byte{}

	var err error
	var interiorPage *PageParsed = &page

	for i, v := range rightPointers {
		if i < len(rightPointers)-1 {
			expectedCellAreaContent = append(expectedCellAreaContent, v...)
			expectedCellAreaContent = append(expectedCellAreaContent, intToBinary(rowIds[i+1], 2)...)
		}

		// if i != len(rightPointers)-1 {
		expectedPointers = append(expectedPointers, intToBinary(PageSize-(i+1)*6, 2)...)
		// }

		err = server.insertPointerIntoInteriorPage(v, rowIds[i], interiorPage)

		if err != nil {
			t.Errorf("Err while adding pointer :%v", err)
		}

	}

	if err != nil {
		t.Errorf("should add pointer without any errors, got %v", err)
	}

	if interiorPage.btreeType != int(TableBtreeInteriorCell) {
		t.Error("incorect btree type")
	}

	if !reflect.DeepEqual(expectedCellAreaContent, interiorPage.cellArea) {
		t.Errorf("Cell area should be: %v, got: %v", expectedCellAreaContent, interiorPage.cellArea)
	}

	if !reflect.DeepEqual(interiorPage.rightMostpointer, rightPointers[2]) {
		t.Errorf("Pointer should be %v , got: %v", rightPointers[2], interiorPage.rightMostpointer)
	}

	if interiorPage.numberofCells != len(rightPointers) {
		t.Errorf("number of cell shoud be :%v, insted we got: %v", len(rightPointers), interiorPage.numberofCells)
	}

	if !reflect.DeepEqual(interiorPage.pointers, expectedPointers) {
		t.Errorf("Expected pointer to be: %v, got: %v", expectedPointers, interiorPage.pointers)
	}

	if interiorPage.framgenetedArea != 0 {
		t.Errorf("should be no fragmeneted area")
	}
}

func TestInsertFindPointerInteriorPageWithAvilableSpaceNoPage(t *testing.T) {
	clearDbFile("test")
	firstPage := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{0, 0, 0, 1},
		latesRow:             &LastPageParseLatestRow{},
	}

	server := ServerStruct{
		firstPage: firstPage,
		reader:    NewReader("aa"),
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(firstPage), 0, "fds", &firstPage)

	foundPage, _ := server.findPointerInteriorPageWithAvilableSpace(0)

	if foundPage != nil {
		t.Errorf("found page should be nil, insted we got: %v", foundPage)
	}
}

func TestInsertFindPointerInteriorPageWithAvilableSpaceFirstPage(t *testing.T) {
	clearDbFile("test")
	firstPage := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 1},
		latesRow:             &LastPageParseLatestRow{},
	}

	secondPage := PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		latesRow:             &LastPageParseLatestRow{},
	}

	server := ServerStruct{
		firstPage: firstPage,
		reader:    NewReader("aa"),
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(firstPage), 0, "fds", &firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 1, "fds", &firstPage)

	foundPage, pageNumber := server.findPointerInteriorPageWithAvilableSpace(0)

	if foundPage == nil {
		t.Error("found page should not be nill")
	}

	if pageNumber != 0 {
		t.Errorf("should return first page, insted it returned: %v", pageNumber)
	}
}

func TestInsertFindPointerInteriorPageWithAvilableSpaceNestedPaged(t *testing.T) {
	firstPage := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 1},
		latesRow:             &LastPageParseLatestRow{},
	}

	secondPage := PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 2},
		latesRow:             &LastPageParseLatestRow{},
	}

	thirdPage := PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		latesRow:             &LastPageParseLatestRow{},
	}

	server := ServerStruct{
		firstPage: firstPage,
		reader:    NewReader("aa"),
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(firstPage), 0, "fds", &firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 1, "fds", &firstPage)
	writer.writeToFile(assembleDbPage(thirdPage), 2, "fds", &firstPage)

	foundPage, pageNumber := server.findPointerInteriorPageWithAvilableSpace(0)

	if foundPage == nil {
		t.Error("found page should not be nill")
	}

	if pageNumber != 1 {
		t.Errorf("should return second page, insted it returned: %v", pageNumber)
	}
}

func TestInsertFindPointerInteriorPageWithAvilableSpaceNestedPageFullCapacity(t *testing.T) {
	clearDbFile("test")
	firstPage := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 1},
		latesRow:             &LastPageParseLatestRow{},
	}

	// whoel page 4096 bytes, 4x poiters = 8 bytes, 4088, 12 header length = 4076, cell area 4076/4 = 1019 to take full page
	cellArea := make([]byte, 1019)
	for i := 0; i < 1019; i++ {
		cellArea[i] = byte(i)
	}

	numberOfCell := 4

	allCellArea := []byte{}
	pointers := []byte{}
	for i := 0; i < numberOfCell; i++ {
		allCellArea = append(allCellArea, cellArea...)

		pointers = append(pointers, intToBinary(PageSize-len(allCellArea)*(i+1), 2)...)
	}

	secondPage := PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		cellArea:             allCellArea,
		startCellContentArea: PageSize - len(cellArea)*numberOfCell,
		numberofCells:        numberOfCell,
		pointers:             pointers,
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 2},
		latesRow:             &LastPageParseLatestRow{},
	}

	thirdPage := PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		latesRow:             &LastPageParseLatestRow{},
	}

	server := ServerStruct{
		firstPage: firstPage,
		reader:    NewReader("aa"),
	}

	// fmt.Println("Are we here??")
	// fmt.Println(len(assembleDbPage(secondPage)))
	// fmt.Println("Are we here??")
	writer := NewWriter()

	writer.writeToFile(assembleDbPage(firstPage), 0, "fds", &firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 1, "fds", &firstPage)
	writer.writeToFile(assembleDbPage(thirdPage), 2, "fds", &firstPage)

	foundPage, pageNumber := server.findPointerInteriorPageWithAvilableSpace(0)

	if foundPage == nil {
		t.Error("found page should not be nill")
	}

	if pageNumber != 0 {
		t.Errorf("should return first page, insted it returned: %v", pageNumber)
	}
}

func TestInsertSchemaOnlyOneLeafPage(t *testing.T) {
	clearDbFile("test")
	firstPage := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		latesRow:             &LastPageParseLatestRow{},
	}

	server := ServerStruct{
		firstPage: firstPage,
		reader:    NewReader("aa"),
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(firstPage), 0, "fds", &firstPage)

	_, parsedQuery := genericParser("CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")

	createSqlQueryData := parseCreateTableQuery(parsedQuery, "CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")
	pointer := 0

	insertedPage, pageNumber := server.insertSchema(string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)

	if pageNumber != 0 {
		t.Errorf("Expected page number to be 0, we got:%v", pageNumber)
	}

	if insertedPage.numberofCells != 1 {
		t.Errorf("Expected nuber of cells on the page to be 1, we got: %v", insertedPage.numberofCells)
	}

}

// func (server *ServerStruct) findPointerInteriorPageWithAvilableSpace(reader PageReader, pageNumber int) (*PageParsed, int) {
// 	page := reader.readDbPage(pageNumber)
// 	parsedPage := parseReadPage(page, pageNumber)
// 	if parsedPage.btreeType != int(TableBtreeLeafCell) {
// 		return nil, 0
// 	}

// 	btreeLeftCellHeaderLength := 12
// 	spaceAvilable := parsedPage.startCellContentArea - btreeLeftCellHeaderLength - len(parsedPage.pointers) - parsedPage.dbHeaderSize
// 	newPointerLength := 2
// 	cellDataLength := 6
// 	newCellSpace := cellDataLength + newPointerLength
// 	if spaceAvilable >= newCellSpace {
// 		return &parsedPage, pageNumber
// 	}

// 	newPageNumber := binary.BigEndian.Uint32(parsedPage.rightMostpointer)

// 	return server.findPointerInteriorPageWithAvilableSpace(reader, int(newPageNumber))
// }

func TestUpdatePageMultiple(t *testing.T) {
	//START HERE,FIX MOVING PAGE FROM 0X0D TO 0X05
	//START HERE,FIX MOVING PAGE FROM 0X0D TO 0X05
	//START HERE,FIX MOVING PAGE FROM 0X0D TO 0X05
	//START HERE,FIX MOVING PAGE FROM 0X0D TO 0X05
	//START HERE,FIX MOVING PAGE FROM 0X0D TO 0X05
	//START HERE,FIX MOVING PAGE FROM 0X0D TO 0X05
	clearDbFile("test")
	dbName = "test"

	parsedData := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		latesRow:             &LastPageParseLatestRow{},
	}

	firstPage := PageParsed{
		dbHeader: DbHeader{
			dbSizeInPages: 1,
		},
	}

	server := ServerStruct{
		firstPage: firstPage,
	}

	reader := NewReader("conId")
	writer := NewWriter()
	page := assembleDbPage(parsedData)
	parsedPage := parseReadPage(page, 0)
	for i := 1; i < 53; i++ {
		fmt.Printf("\niteration :%v", i)
		table := fmt.Sprintf("user%v", i)
		query := fmt.Sprintf("CREATE TABLE user%v (id INTEGER PRIMARY KEY,name TEXT)", i)

		page = server.updatePageRoot(&parsedPage, 0, *reader, *writer, string(SqlQueryTableObjectType), table, table, i, query)
		// fmt.Println("did we return it")

		// if i == 49 {
		// 	fmt.Println(page)
		// }

		parsedPage = parseReadPage(page, 0)
	}

}

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

	foundPage, _, _ := server.findLastPointerInteriorPageWithAvilableSpace(0)

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

	foundPage, pageNumber, _ := server.findLastPointerInteriorPageWithAvilableSpace(0)

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

	foundPage, pageNumber, _ := server.findLastPointerInteriorPageWithAvilableSpace(0)

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

	foundPage, pageNumber, _ := server.findLastPointerInteriorPageWithAvilableSpace(0)

	if foundPage == nil {
		t.Error("found page should not be nill")
	}

	if pageNumber != 0 {
		t.Errorf("should return first page, insted it returned: %v", pageNumber)
	}
}

func TestInsertData(t *testing.T) {
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

	value := "Alice"

	cell := createCell(TableBtreeLeafCell, nil, value)

	server.insertData(cell, &firstPage)

	if firstPage.numberofCells != 1 {
		t.Errorf("expected to have one cell, insted we got: %v", firstPage.numberofCells)
	}

	if !reflect.DeepEqual(firstPage.cellArea, cell.data) {
		t.Errorf("expected cell area to be: %v, got: %v", cell.data, firstPage.cellArea)
	}

	pointer := intToBinary(PageSize-cell.dataLength, 2)

	if !reflect.DeepEqual(firstPage.pointers, pointer) {
		t.Errorf("expected one pointer: %v, insted we got: %v", pointer, firstPage.pointers)
	}
}

func TestInsertMultipleData(t *testing.T) {
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

	value := "Alice"

	cell := createCell(TableBtreeLeafCell, nil, value)

	server.insertData(cell, &firstPage)

	if firstPage.numberofCells != 1 {
		t.Errorf("expected to have one cell, insted we got: %v", firstPage.numberofCells)
	}
	expectedCellArea := cell.data

	if !reflect.DeepEqual(firstPage.cellArea, expectedCellArea) {
		t.Errorf("expected cell area to be: %v, got: %v", expectedCellArea, firstPage.cellArea)
	}

	pointerOne := PageSize - cell.dataLength
	pointers := intToBinary(pointerOne, 2)

	if !reflect.DeepEqual(firstPage.pointers, pointers) {
		t.Errorf("expected one pointer: %v, insted we got: %v", pointers, firstPage.pointers)
	}

	if firstPage.startCellContentArea != pointerOne {
		t.Errorf("expected start content area to be: %v, got: %v", pointerOne, firstPage.startCellContentArea)
	}

	server.insertData(cell, &firstPage)

	if firstPage.numberofCells != 2 {
		t.Errorf("expected to have two cell, insted we got: %v", firstPage.numberofCells)
	}

	expectedCellArea = append(expectedCellArea, cell.data...)

	if !reflect.DeepEqual(firstPage.cellArea, expectedCellArea) {
		t.Errorf("expected cell area to be: %v, got: %v", cell.data, firstPage.cellArea)
	}

	pointerTwo := PageSize - (cell.dataLength * 2)
	pointers = append(pointers, intToBinary(pointerTwo, 2)...)

	if !reflect.DeepEqual(firstPage.pointers, pointers) {
		t.Errorf("expected one pointer: %v, insted we got: %v", pointers, firstPage.pointers)
	}

	if firstPage.startCellContentArea != pointerTwo {
		t.Errorf("expected start content area to be: %v, got: %v", pointerTwo, firstPage.startCellContentArea)
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

	cell := createCell(TableBtreeLeafCell, nil, string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)

	insertedPage, pageNumber := server.insertSchema(string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)

	firstPageRead := server.reader.readDbPage(0)
	firstPageReadParsed := parseReadPage(firstPageRead, 0)

	if !reflect.DeepEqual(firstPageReadParsed, *insertedPage) {
		t.Errorf("should save same page as one as it returned, expected: %v, got: %v", firstPageReadParsed, insertedPage)
	}

	if pageNumber != 0 {
		t.Errorf("Expected page number to be 0, we got:%v", pageNumber)
	}

	if insertedPage.numberofCells != 1 {
		t.Errorf("Expected nuber of cells on the page to be 1, we got: %v", insertedPage.numberofCells)
	}

	if !reflect.DeepEqual(insertedPage.cellArea, cell.data) {
		t.Errorf("Expected cell area to be: %v, got: %v", cell.data, insertedPage.cellArea)
	}

	if !reflect.DeepEqual(insertedPage.pointers, intToBinary(PageSize-cell.dataLength, 2)) {
		t.Errorf("Expected one pointer with value: %v, got: %v", intToBinary(PageSize-cell.dataLength, 2), insertedPage.pointers)
	}

}

func TestInsertSchemaOnlyOneInteriorAndLeafPage(t *testing.T) {
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

	_, parsedQuery := genericParser("CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")

	createSqlQueryData := parseCreateTableQuery(parsedQuery, "CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")
	pointer := 0

	cell := createCell(TableBtreeLeafCell, nil, string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)
	insertedPage, pageNumber := server.insertSchema(string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)

	secondPageRead := server.reader.readDbPage(1)
	secondPageReadParsed := parseReadPage(secondPageRead, 1)

	if !reflect.DeepEqual(secondPageReadParsed, *insertedPage) {
		t.Errorf("should save same page as one as it returned, expected: %v, got: %v", secondPageReadParsed, insertedPage)
	}

	if pageNumber != 1 {
		t.Errorf("Expected page number to be %v, we got:%v", 1, pageNumber)
	}

	if insertedPage.numberofCells != 1 {
		t.Errorf("Expected nuber of cells on the page to be 1, we got: %v", insertedPage.numberofCells)
	}

	if !reflect.DeepEqual(insertedPage.cellArea, cell.data) {
		t.Errorf("Expected cell area to be: %v, got: %v", cell.data, insertedPage.cellArea)
	}

	if !reflect.DeepEqual(insertedPage.pointers, intToBinary(PageSize-cell.dataLength, 2)) {
		t.Errorf("Expected one pointer with value: %v, got: %v", intToBinary(PageSize-cell.dataLength, 2), insertedPage.pointers)
	}
}

func createAFullLeafPage() PageParsed {
	page := PageParsed{
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
		firstPage: PageParsed{dbHeader: DbHeader{}},
		reader:    NewReader("aa"),
	}

	cellArea := make([]byte, 1014)
	for i := 0; i < 1014; i++ {
		cellArea[i] = byte('a')
	}

	for i := 0; i < 4; i++ {
		cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: i}}, string(cellArea))
		server.insertData(cell, &page)
	}

	return page

}

func createAFullLeafPageWitHeader() PageParsed {
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
	server := ServerStruct{
		firstPage: PageParsed{dbHeader: DbHeader{}},
		reader:    NewReader("aa"),
	}

	cellArea := make([]byte, 989)
	for i := 0; i < 989; i++ {
		cellArea[i] = byte('a')
	}

	for i := 0; i < 4; i++ {
		cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: i}}, string(cellArea))
		server.insertData(cell, &page)
	}

	return page

}

func createAFullInteriorPageWitHeader() PageParsed {
	page := PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{},
		latesRow:             &LastPageParseLatestRow{},
	}
	server := ServerStruct{
		firstPage: PageParsed{dbHeader: DbHeader{}},
		reader:    NewReader("aa"),
	}

	cellArea := make([]byte, 989)
	for i := 0; i < 989; i++ {
		cellArea[i] = byte('a')
	}

	for i := 0; i < 4; i++ {
		cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: i}}, string(cellArea))
		server.insertData(cell, &page)
	}

	return page
}

func createAFullInteriorPage() PageParsed {
	page := PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		cellArea:             []byte{},
		startCellContentArea: PageSize,
		numberofCells:        0,
		pointers:             []byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{},
		latesRow:             &LastPageParseLatestRow{},
	}
	server := ServerStruct{
		firstPage: PageParsed{dbHeader: DbHeader{}},
		reader:    NewReader("aa"),
	}

	cellArea := make([]byte, 1014)
	for i := 0; i < 1014; i++ {
		cellArea[i] = byte('b')
	}

	for i := 0; i < 4; i++ {
		cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: i}}, string(cellArea))
		server.insertData(cell, &page)
	}

	return page

}

func TestInsertSchemaOnlyOneInteriorAndFullLeafPage(t *testing.T) {
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

	fmt.Println("hello first page??")
	fmt.Printf("%+v", firstPage)

	server := ServerStruct{
		firstPage: firstPage,
		reader:    NewReader("aa"),
	}

	secondPage := createAFullLeafPage()

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(firstPage), 0, "fds", &firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 1, "fds", &firstPage)

	_, parsedQuery := genericParser("CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")

	createSqlQueryData := parseCreateTableQuery(parsedQuery, "CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")
	pointer := 0

	cell := createCell(TableBtreeLeafCell, nil, string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)
	insertedPage, pageNumber := server.insertSchema(string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)

	if pageNumber != 2 {
		t.Errorf("Expected page number to be %v, we got:%v", 2, pageNumber)
	}

	if insertedPage.numberofCells != 1 {
		t.Errorf("Expected nuber of cells on the page to be 1, we got: %v", insertedPage.numberofCells)
	}

	if !reflect.DeepEqual(insertedPage.cellArea, cell.data) {
		t.Errorf("Expected cell area to be: %v, got: %v", cell.data, insertedPage.cellArea)
	}

	if !reflect.DeepEqual(insertedPage.pointers, intToBinary(PageSize-cell.dataLength, 2)) {
		t.Errorf("Expected one pointer with value: %v, got: %v", intToBinary(PageSize-cell.dataLength, 2), insertedPage.pointers)
	}

	firstPageRead := server.reader.readDbPage(0)
	firstPageReadParsed := parseReadPage(firstPageRead, 0)

	fmt.Println("parsed new page")
	fmt.Printf("%+v", firstPageReadParsed)

	if !reflect.DeepEqual(firstPageReadParsed.rightMostpointer, []byte{0, 0, 0, 2}) {
		t.Errorf("expected right most pointer now to be to second page, insted we got: %v", firstPageReadParsed.rightMostpointer)
	}

	expectedCellContentArea := firstPage.rightMostpointer
	expectedCellContentArea = append(expectedCellContentArea, []byte{0, 0}...)

	if !reflect.DeepEqual(firstPageReadParsed.cellArea, expectedCellContentArea) {
		t.Errorf("expected cell content area to be: %v, got: %v", expectedCellContentArea, firstPageReadParsed.cellArea)
	}

	if firstPageReadParsed.numberofCells != 1 {
		t.Errorf("expected cell number to be: %v, got: %v", 1, firstPageReadParsed.numberofCells)
	}

	if !reflect.DeepEqual(firstPageReadParsed.pointers, intToBinary(PageSize-6, 2)) {
		t.Errorf("expected cell number to be: %v, got: %v", 1, firstPageReadParsed.numberofCells)
	}

	if firstPageReadParsed.startCellContentArea != PageSize-6 {
		t.Errorf("expected start cell content area to be: %v, got: %v", PageSize-6, firstPageReadParsed.startCellContentArea)
	}
}

// TODO
func TestInsertSchemaOneFullLeafPageSpaceAvialableAfterHeaderStrip(t *testing.T) {
	clearDbFile("test")
	firstPage := createAFullLeafPageWitHeader()

	server := ServerStruct{
		firstPage: firstPage,
		reader:    NewReader("aa"),
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(firstPage), 0, "fds", &firstPage)

	fmt.Println("checkpoint 1")

	_, parsedQuery := genericParser("CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")

	createSqlQueryData := parseCreateTableQuery(parsedQuery, "CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")
	pointer := 0

	cell := createCell(TableBtreeLeafCell, nil, string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)
	fmt.Println("checkpoint 2")
	insertedPage, pageNumber := server.insertSchema(string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)

	if pageNumber != 1 {
		t.Errorf("Expected page number to be %v, we got:%v", 1, pageNumber)
	}

	if insertedPage.numberofCells != firstPage.numberofCells+1 {
		t.Errorf("Expected nuber of cells on the page to be %v, we got: %v", firstPage.numberofCells+1, insertedPage.numberofCells)
	}

	expectedCellContentArea := cell.data
	expectedCellContentArea = append(expectedCellContentArea, firstPage.cellArea...)

	if !reflect.DeepEqual(insertedPage.cellArea, expectedCellContentArea) {
		t.Errorf("Expected cell area to be: %v, got: %v", expectedCellContentArea, insertedPage.cellArea)
	}

	expectedPointers := firstPage.pointers
	expectedPointers = append(expectedPointers, intToBinary(firstPage.startCellContentArea-cell.dataLength, 2)...)

	if !reflect.DeepEqual(insertedPage.pointers, expectedPointers) {
		t.Errorf("Expected one pointer with value: %v, got: %v", expectedPointers, insertedPage.pointers)
	}

	firstPageRead := server.reader.readDbPage(0)
	firstPageReadParsed := parseReadPage(firstPageRead, 0)

	if firstPageReadParsed.btreeType != int(TableBtreeInteriorCell) {
		t.Errorf("Expected btree type to be interior cell, insted we got: %v", firstPageReadParsed.btreeType)
	}

	expectedHeader := firstPage.dbHeader
	expectedHeader.dbSizeInPages++
	expectedHeader.fileChangeCounter++
	expectedHeader.versionValidForNumber++

	if !reflect.DeepEqual(firstPageReadParsed.dbHeader, expectedHeader) {
		t.Errorf("expected header: %v, \ngot: %v", expectedHeader, firstPageReadParsed.dbHeader)
	}

	if firstPageReadParsed.numberofCells != 0 {
		t.Errorf("expectd number of cells to be: %v, got: %v", 0, firstPageReadParsed.numberofCells)
	}

	if reflect.DeepEqual(firstPageReadParsed.rightMostpointer, []byte{0, 0, 0, 1}) {
		t.Errorf("expectd right most pointer to be: %v, got: %v", []byte{0, 0, 0, 1}, firstPageReadParsed.rightMostpointer)
	}

}

func TestInsertSchemaFullLeafPageAndFullInterior(t *testing.T) {
	clearDbFile("test")
	firstPage := createAFullInteriorPageWitHeader()
	firstPage.rightMostpointer = []byte{0, 0, 0, 1}
	secondPage := createAFullLeafPage()

	fmt.Printf("\n???D?FSD?\n")
	fmt.Println(secondPage.startCellContentArea)

	server := ServerStruct{
		firstPage: firstPage,
		reader:    NewReader("aa"),
		writer:    NewWriter(),
	}
	// we have 0x0d, and 0x05, both full, so now we need to create new 0x0d, new 0x05, and new root page 0x05 so 3 pages

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(firstPage), 0, "fds", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 1, "fds", &server.firstPage)

	_, parsedQuery := genericParser("CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")

	createSqlQueryData := parseCreateTableQuery(parsedQuery, "CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")
	pointer := 0

	cell := createCell(TableBtreeLeafCell, nil, string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)
	fmt.Println("checkpoint 2")
	insertedPage, pageNumber := server.insertSchema(string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)

	if pageNumber != 2 {
		t.Errorf("Expected page number to be: %v instead we got: %v", 2, pageNumber)
	}

	if insertedPage.numberofCells != 1 {
		t.Errorf("expectd number of cells to be: %v, got: %v", 1, insertedPage.numberofCells)
	}

	if !reflect.DeepEqual(insertedPage.cellArea, cell.data) {
		t.Errorf("Expected cell area to be: %v, got: %v", cell.data, insertedPage.cellArea)
	}

	if !reflect.DeepEqual(insertedPage.pointers, intToBinary(PageSize-cell.dataLength, 2)) {
		t.Errorf("expected one pointer: %v, got: %v", intToBinary(PageSize-cell.dataLength, 2), insertedPage.pointers)
	}

	firstPageRead := server.reader.readDbPage(0)
	firstPageReadParsed := parseReadPage(firstPageRead, 0)

	secondPageRead := server.reader.readDbPage(1)
	secondPageReadParsed := parseReadPage(secondPageRead, 1)

	fourthPageRead := server.reader.readDbPage(3)
	fourthPageReadParsed := parseReadPage(fourthPageRead, 3)

	if firstPageReadParsed.numberofCells != 0 {
		t.Errorf("expected to have new interior page with only one cell, insted we got: %v", firstPageReadParsed.numberofCells)
	}

	if firstPageReadParsed.btreeType != int(TableBtreeInteriorCell) {
		t.Errorf("expected btree to be %v, got: %v", TableBtreeInteriorCell, firstPageReadParsed.btreeType)
	}

	if !reflect.DeepEqual(firstPageReadParsed.rightMostpointer, []byte{0, 0, 0, 3}) {
		t.Errorf("Expected right most pointer to be: %v, insted we got: %v", []byte{0, 0, 0, 3}, firstPageReadParsed.rightMostpointer)
	}

	firstPageStrippedHeader := firstPage
	firstPageStrippedHeader.dbHeader = DbHeader{}
	firstPageStrippedHeader.dbHeaderSize = 0

	if !reflect.DeepEqual(fourthPageReadParsed, firstPageStrippedHeader) {
		t.Errorf("expected last page to be move moved expected: %v, got :%v", firstPageStrippedHeader, fourthPageReadParsed)
	}

	fmt.Println(secondPageReadParsed.rightMostpointer == nil)
	fmt.Println(secondPage.rightMostpointer == nil)

	if !reflect.DeepEqual(secondPageReadParsed, secondPage) {
		t.Errorf("expected leaf page to stay as it was expected: %v, got :%v", secondPage.rightMostpointer, secondPageReadParsed.rightMostpointer)
	}

}

func TestMultipleFullLeafPage(t *testing.T) {
	clearDbFile("test")
	firstPage := createAFullInteriorPageWitHeader()
	firstPage.rightMostpointer = []byte{0, 0, 0, 1}
	secondPage := createAFullInteriorPage()
	secondPage.rightMostpointer = []byte{0, 0, 0, 2}
	thirdPage := createAFullLeafPage()

	server := ServerStruct{
		firstPage: firstPage,
		reader:    NewReader("aa"),
		writer:    NewWriter(),
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(firstPage), 0, "fds", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 1, "fds", &server.firstPage)
	writer.writeToFile(assembleDbPage(thirdPage), 2, "fds", &server.firstPage)

	_, parsedQuery := genericParser("CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")
	createSqlQueryData := parseCreateTableQuery(parsedQuery, "CREATE TABLE user (id INTEGER PRIMARY KEY,name TEXT)")
	pointer := 0

	cell := createCell(TableBtreeLeafCell, nil, string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)
	fmt.Println("checkpoint 2")
	insertedPage, pageNumber := server.insertSchema(string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointer, createSqlQueryData.rawQuery)
	fmt.Println("test??")
	if pageNumber != 3 {
		t.Errorf("Should append data to exisitng page: %v, instead we got :%v", 3, pageNumber)
	}

	if insertedPage.numberofCells != 1 {
		t.Errorf("expected number of cell to be: %v, we got: %v", 1, insertedPage.numberofCells)
	}

	if !reflect.DeepEqual(cell.data, insertedPage.cellArea) {
		t.Errorf("expected cell area to be: %v, instead we got: %v", cell.data, insertedPage.cellArea)
	}

	firstPageRead := server.reader.readDbPage(0)
	firstPageReadParsed := parseReadPage(firstPageRead, 0)

	secondPageRead := server.reader.readDbPage(1)
	secondPageReadParsed := parseReadPage(secondPageRead, 1)

	fourthPageRead := server.reader.readDbPage(3)
	fourthPageReadParsed := parseReadPage(fourthPageRead, 3)

	fifthPageRead := server.reader.readDbPage(4)
	fifthPageReadParsed := parseReadPage(fifthPageRead, 4)

	// There is problem with saving last page

	sixthPageRead := server.reader.readDbPage(5)
	sixthPageReadParsed := parseReadPage(sixthPageRead, 5)

	// fmt.Println("fifth")
	// fmt.Printf("%+v", sixthPageReadParsed)

	// fmt.Println("six page")
	// fmt.Printf("%+v", sixthPageReadParsed)

	if fourthPageReadParsed.btreeType != int(TableBtreeLeafCell) {
		t.Errorf("Expected fourth page header to be: %v, insted we got: %v", TableBtreeLeafCell, fourthPageReadParsed.btreeType)
	}

	if fifthPageReadParsed.btreeType != int(TableBtreeInteriorCell) {
		t.Errorf("Expected fourth page header to be: %v, insted we got: %v", TableBtreeInteriorCell, fifthPageReadParsed.btreeType)
	}

	if sixthPageReadParsed.btreeType != int(TableBtreeInteriorCell) {
		t.Errorf("Expected fourth page header to be: %v, insted we got: %v", TableBtreeInteriorCell, sixthPageReadParsed.btreeType)
	}

	// fmt.Println("parsed")
	// fmt.Printf("%+v", thirdPageReadParsed)

	if firstPageReadParsed.btreeType != int(TableBtreeInteriorCell) {
		t.Errorf("Expected header to be: %v, got: %v", TableBtreeInteriorCell, firstPageReadParsed.btreeType)
	}

	if firstPageReadParsed.numberofCells != 0 {
		t.Errorf("Expected %v cells, got: %v", 0, firstPageReadParsed.numberofCells)
	}

	if !reflect.DeepEqual(firstPageReadParsed.rightMostpointer, []byte{0, 0, 0, 5}) {
		t.Errorf("expected pointer to be: %v, got: %v", []byte{0, 0, 0, 5}, firstPageReadParsed.rightMostpointer)
	}

	if fourthPageReadParsed.numberofCells != 1 {
		t.Errorf("expected fourth page to have only number of cells: %v, we got: %v", 1, fourthPageReadParsed.numberofCells)
	}

	if !reflect.DeepEqual(fourthPageReadParsed.cellArea, cell.data) {
		t.Errorf("expected fourth page to have only number of cells: %v, we got: %v", 1, fourthPageReadParsed.numberofCells)
	}

	fmt.Println(firstPageReadParsed.rightMostpointer)
	fmt.Println(secondPageReadParsed.rightMostpointer)
	fmt.Println(fifthPageReadParsed.rightMostpointer)
	fmt.Println(sixthPageReadParsed.rightMostpointer)
}

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

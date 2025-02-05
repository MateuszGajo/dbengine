package main

import (
	"encoding/binary"
	"fmt"
	"io/fs"
	"reflect"
	"testing"
	"time"
)

func TestCreateCellWithSingleValue(t *testing.T) {
	btreeType := TableBtreeLeafCell
	page := LastPageParsed{
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
	page := LastPageParsed{
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
	page := LastPageParsed{
		latesRow: &LastPageParseLatestRow{
			rowId: 1,
		},
	}
	val2 := 12     //1
	val3 := "test" //4
	cell := createCell(btreeType, &page, nil, val2, val3)
	btreeHeader := BtreeHeaderSchema(TableBtreeLeafCell, cell, nil)

	if len(btreeHeader) != 10 {
		t.Errorf("Header should have 10 bytes, we got: %v", len(btreeHeader))
	}

	if btreeHeader[0] != 0x0d {
		t.Errorf("Expected header type to be: %v, insted we got: %v", 0x0d, btreeHeader[0])
	}

	if btreeHeader[1] != 0 || btreeHeader[2] != 0 {
		t.Errorf("Expected free block to be: %v, insted we got bytes: %v %v", 0, btreeHeader[1], btreeHeader[2])
	}

	if binary.BigEndian.Uint16(btreeHeader[3:5]) != 1 {
		t.Errorf("Expected number of cell to be: %v, insted we got: %v", 1, binary.BigEndian.Uint16(btreeHeader[3:5]))
	}

	if binary.BigEndian.Uint16(btreeHeader[5:7]) != uint16(PageSize-cell.dataLength-2) {
		t.Errorf("Expected start content area to be: %v, insted we got : %v", uint16(PageSize-cell.dataLength-2), binary.BigEndian.Uint16(btreeHeader[5:7]))
	}

	if btreeHeader[7] != 0 {
		t.Errorf("Expected fragmeneted free bytes to be %v, insted we got : %v", 0, btreeHeader[7])
	}

	if binary.BigEndian.Uint16(btreeHeader[8:10]) != uint16(PageSize-cell.dataLength-2) {
		t.Errorf("Expected cell's pointer to be: %v, insted we got : %v", uint16(PageSize-cell.dataLength-2), binary.BigEndian.Uint16(btreeHeader[8:10]))
	}

}

type MockFileInfo struct {
	NameVal    string
	SizeVal    int64
	ModeVal    fs.FileMode
	ModTimeVal time.Time
	IsDirVal   bool
}

func (m MockFileInfo) Name() string       { return m.NameVal }
func (m MockFileInfo) Size() int64        { return m.SizeVal }
func (m MockFileInfo) Mode() fs.FileMode  { return m.ModeVal }
func (m MockFileInfo) ModTime() time.Time { return m.ModTimeVal }
func (m MockFileInfo) IsDir() bool        { return m.IsDirVal }
func (m MockFileInfo) Sys() any           { return nil }

func TestParseDbPageWithOnlyHeader(t *testing.T) {
	btreeHeader := BtreeHeaderSchema(TableBtreeLeafCell, CreateCell{dataLength: 0}, nil)
	zeros := make([]byte, PageSize-len(btreeHeader))
	btreeHeader = append(btreeHeader, zeros...)
	res := parseReadPage(btreeHeader, false, MockFileInfo{SizeVal: 10})

	if res.btreeType != int(TableBtreeLeafCell) {
		t.Errorf("Expected: %v tree type, insted we got: %v", TableBtreeLeafCell, res.btreeType)
	}

	if res.framgenetedArea != 0 {
		t.Errorf("Expected fragmeneted area to be: %v instead we got: %v", 0, res.framgenetedArea)
	}

	if res.freeBlock != 0 {
		t.Errorf("Expected first free block address to be: %v, insted we got: %v", 0, res.freeBlock)
	}

	if res.numberofCells != 0 {
		t.Errorf("Expected numbrs of cell to be: %v, insted we got: %v", 0, res.numberofCells)
	}

	if res.startCellContentArea != PageSize {
		t.Errorf("Expected start cell of content area to be: %v, insted we got: %v", PageSize, res.startCellContentArea)
	}

}

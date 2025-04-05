package main

import (
	"fmt"
	"reflect"
	"testing"
)

func cleanUp() {
	PageSize = 4096
}

func TestCreateNewPageEmptyNoHeader(t *testing.T) {
	t.Cleanup(cleanUp)
	page := CreateNewPage(TableBtreeLeafCell, [][]byte{}, 1, nil)

	if page.btreeType != int(TableBtreeLeafCell) {
		t.Errorf("expected header to be %v, instead we got: %v", TableBtreeLeafCell, page.btreeType)
	}

	if page.btreePageHeaderSize != 8 {
		t.Errorf("expected btree page header size to be: %v, got: %v", 8, page.btreePageHeaderSize)
	}

	if !reflect.DeepEqual(page.cellArea, []byte{}) {
		t.Errorf("expected cell area to be: %v, got: %v", []byte{}, page.cellArea)
	}

	if !reflect.DeepEqual(page.cellAreaParsed, [][]byte{}) {
		t.Errorf("expected parsed cell area to be: %v, got: %v", [][]byte{}, page.cellAreaParsed)
	}

	if !reflect.DeepEqual(page.pointers, []byte{}) {
		t.Errorf("expected pointer to be: %v, got: %v", []byte{}, page.pointers)
	}

	if !reflect.DeepEqual(page.rightMostpointer, []byte{}) {
		t.Errorf("expected right most pointer to be: %v, got: %v", []byte{}, page.rightMostpointer)
	}

	if page.cellAreaSize != 0 {
		t.Errorf("expected cell area size to be:%v got: %v", 0, page.cellArea)
	}

	if page.dbHeaderSize != 0 {
		t.Errorf("expected db header size to be: %v, got: %v", 0, page.dbHeaderSize)
	}

	if page.framgenetedArea != 0 {
		t.Errorf("expected fragmented area to be: %v, got: %v", 0, page.framgenetedArea)
	}

	if page.freeBlock != 0 {
		t.Errorf("expected free area to be: %v, got: %v", 0, page.freeBlock)
	}

	if page.isLeaf != true {
		t.Errorf("expected page is leaf prop to be: %v, got: %v", true, page.isLeaf)
	}

	if page.latesRow != nil {
		t.Errorf("expected latestRow to be: %v, instead we got: %v", nil, page.latesRow)
	}

	if page.isOverflow != false {
		t.Errorf("expected is overflow to be: %v, got: %v", false, page.isOverflow)
	}

	if page.numberofCells != 0 {
		t.Errorf("expected number of cells to be: %v, got: %v", 0, page.numberofCells)
	}

	if page.pageNumber != 1 {
		t.Errorf("expected page number to be: %v, got :%v", 1, page.pageNumber)
	}
}

func TestCreateNewPageNoHeader(t *testing.T) {
	t.Cleanup(cleanUp)
	expectedCellArea := []byte{}
	expectedCellAreaStart := PageSize
	expectedPointers := []byte{}
	expectedNumberOfCells := 0
	expectedCellAreaSize := 0
	var expectedLatestRowParsed *LastPageParseLatestRow
	iterations := 5
	for i := iterations; i >= 0; i-- {
		cell := createCell(TableBtreeLeafCell, i, fmt.Sprintf("Alice%v", i))
		expectedCellArea = append(expectedCellArea, cell.data...)
		expectedCellAreaStart -= cell.dataLength
		expectedPointers = append(expectedPointers, intToBinary(expectedCellAreaStart, 2)...)
		expectedNumberOfCells++
		expectedCellAreaSize += cell.dataLength

		if i == iterations {
			expectedLatestRowParsed = &LastPageParseLatestRow{
				rowId: i,
				data:  cell.data,
			}
		}
	}
	expectedCellAreaParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), expectedCellArea)
	page := CreateNewPage(TableBtreeLeafCell, expectedCellAreaParsed, 1, nil)

	if page.btreeType != int(TableBtreeLeafCell) {
		t.Errorf("expected header to be %v, instead we got: %v", TableBtreeLeafCell, page.btreeType)
	}

	if page.btreePageHeaderSize != 8 {
		t.Errorf("expected btree page header size to be: %v, got: %v", 8, page.btreePageHeaderSize)
	}

	if !reflect.DeepEqual(page.cellArea, expectedCellArea) {
		t.Errorf("expected cell area to be: %v, got: %v", expectedCellArea, page.cellArea)
	}

	if !reflect.DeepEqual(page.cellAreaParsed, expectedCellAreaParsed) {
		t.Errorf("expected parsed cell area to be: %v, got: %v", expectedCellAreaParsed, page.cellAreaParsed)
	}

	if !reflect.DeepEqual(page.pointers, expectedPointers) {
		t.Errorf("expected pointer to be: %v, got: %v", expectedPointers, page.pointers)
	}

	if !reflect.DeepEqual(page.rightMostpointer, []byte{}) {
		t.Errorf("expected right most pointer to be: %v, got: %v", []byte{}, page.rightMostpointer)
	}

	if page.cellAreaSize != expectedCellAreaSize {
		t.Errorf("expected cell area size to be:%v got: %v", expectedCellAreaSize, page.cellArea)
	}

	if page.dbHeaderSize != 0 {
		t.Errorf("expected db header size to be: %v, got: %v", 0, page.dbHeaderSize)
	}

	if page.framgenetedArea != 0 {
		t.Errorf("expected fragmented area to be: %v, got: %v", 0, page.framgenetedArea)
	}

	if page.freeBlock != 0 {
		t.Errorf("expected free area to be: %v, got: %v", 0, page.freeBlock)
	}

	if page.isLeaf != true {
		t.Errorf("expected page is leaf prop to be: %v, got: %v", true, page.isLeaf)
	}

	if !reflect.DeepEqual(page.latesRow, expectedLatestRowParsed) {
		t.Errorf("expected latestRow to be: %v, instead we got: %v", expectedLatestRowParsed, page.latesRow)
	}

	if page.isOverflow != false {
		t.Errorf("expected is overflow to be: %v, got: %v", false, page.isOverflow)
	}

	if page.numberofCells != expectedNumberOfCells {
		t.Errorf("expected number of cells to be: %v, got: %v", expectedNumberOfCells, page.numberofCells)
	}

	if page.pageNumber != 1 {
		t.Errorf("expected page number to be: %v, got :%v", 1, page.pageNumber)
	}
}

func TestCreateNewPageWithHeader(t *testing.T) {
	t.Cleanup(cleanUp)
	expectedCellArea := []byte{}
	expectedCellAreaStart := PageSize
	expectedPointers := []byte{}
	expectedNumberOfCells := 0
	expectedCellAreaSize := 0
	expectedRightMostPointer := []byte{}
	var expectedLatestRowParsed *LastPageParseLatestRow
	iterations := 5
	expectedBtreeType := TableBtreeInteriorCell
	for i := iterations; i >= 0; i-- {
		cell := []byte{0, 0, 0, byte(i), 0, byte(i * 5)}
		expectedCellArea = append(expectedCellArea, cell...)
		expectedCellAreaStart -= len(cell)
		expectedPointers = append(expectedPointers, intToBinary(expectedCellAreaStart, 2)...)
		expectedNumberOfCells++
		expectedCellAreaSize += len(cell)

		if i == iterations {
			expectedLatestRowParsed = &LastPageParseLatestRow{
				rowId: i * 5,
				data:  cell,
			}
			expectedRightMostPointer = []byte{0, 0, 0, byte(i)}
		}
	}
	expectedCellAreaParsed := dbReadparseCellArea(byte(expectedBtreeType), expectedCellArea)
	page := CreateNewPage(expectedBtreeType, expectedCellAreaParsed, 0, &DbHeader{})

	if page.btreeType != int(expectedBtreeType) {
		t.Errorf("expected header to be %v, instead we got: %v", expectedBtreeType, page.btreeType)
	}

	if page.btreePageHeaderSize != 12 {
		t.Errorf("expected btree page header size to be: %v, got: %v", 12, page.btreePageHeaderSize)
	}

	if !reflect.DeepEqual(page.cellArea, expectedCellArea) {
		t.Errorf("expected cell area to be: %v, got: %v", expectedCellArea, page.cellArea)
	}

	if !reflect.DeepEqual(page.cellAreaParsed, expectedCellAreaParsed) {
		t.Errorf("expected parsed cell area to be: %v, got: %v", expectedCellAreaParsed, page.cellAreaParsed)
	}

	if !reflect.DeepEqual(page.pointers, expectedPointers) {
		t.Errorf("expected pointer to be: %v, got: %v", expectedPointers, page.pointers)
	}

	if !reflect.DeepEqual(page.rightMostpointer, expectedRightMostPointer) {
		t.Errorf("expected right most pointer to be: %v, got: %v", expectedRightMostPointer, page.rightMostpointer)
	}

	if page.cellAreaSize != expectedCellAreaSize {
		t.Errorf("expected cell area size to be:%v got: %v", expectedCellAreaSize, page.cellArea)
	}

	if page.dbHeaderSize != 100 {
		t.Errorf("expected db header size to be: %v, got: %v", 100, page.dbHeaderSize)
	}

	if page.framgenetedArea != 0 {
		t.Errorf("expected fragmented area to be: %v, got: %v", 0, page.framgenetedArea)
	}

	if page.freeBlock != 0 {
		t.Errorf("expected free area to be: %v, got: %v", 0, page.freeBlock)
	}

	if page.isLeaf != false {
		t.Errorf("expected page is leaf prop to be: %v, got: %v", false, page.isLeaf)
	}

	if !reflect.DeepEqual(page.latesRow, expectedLatestRowParsed) {
		t.Errorf("expected latestRow to be: %v, instead we got: %v", expectedLatestRowParsed, page.latesRow)
	}

	if page.isOverflow != false {
		t.Errorf("expected is overflow to be: %v, got: %v", false, page.isOverflow)
	}

	if page.numberofCells != expectedNumberOfCells {
		t.Errorf("expected number of cells to be: %v, got: %v", expectedNumberOfCells, page.numberofCells)
	}

	if page.pageNumber != 0 {
		t.Errorf("expected page number to be: %v, got :%v", 0, page.pageNumber)
	}
}

func TestOverflowPage(t *testing.T) {
	PageSize = 20
	t.Cleanup(cleanUp)
	expectedCellArea := []byte{}
	iterations := 5
	for i := iterations; i >= 0; i-- {
		cell := createCell(TableBtreeLeafCell, i, fmt.Sprintf("Alice%v", i))
		expectedCellArea = append(expectedCellArea, cell.data...)

	}
	expectedCellAreaParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), expectedCellArea)
	page := CreateNewPage(TableBtreeLeafCell, expectedCellAreaParsed, 1, nil)

	if page.isOverflow != true {
		t.Errorf("expected overflow value to be %v,  got: %v", true, page.isOverflow)
	}

}

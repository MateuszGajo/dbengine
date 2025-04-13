package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func clearDbFile(fileName string) {
	memoryPages = make(map[int]PageParsed)
	dbName = fileName
	path, err := os.Getwd()
	PageSize = 4096
	if err != nil {
		panic(err)
	}
	os.Remove(path + "/" + fileName + ".db")

}

func TestShouldFindResultSingleInteriorPlusLeafPage(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:         header(),
		dbHeaderSize:     100,
		pageNumber:       0,
		cellAreaParsed:   [][]byte{{0, 0, 0, 2, 0, 7}, {0, 0, 0, 1, 0, 4}},
		btreeType:        int(TableBtreeInteriorCell),
		rightMostpointer: []byte{0, 0, 0, 3},
		cellArea:         []byte{0, 0, 0, 2, 0, 7, 0, 0, 0, 1, 0, 4},

		startCellContentArea: PageSize - 12,
		isOverflow:           false,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 4, 0, 0, 0, 0}, []byte{4, 3, 0, 0, 0, 0}, []byte{4, 2, 0, 0, 0, 0}, []byte{4, 1, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 7, 0, 0, 0, 0}, []byte{4, 6, 0, 0, 0, 0}, []byte{4, 5, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (3 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 12, 0, 0, 0, 0}, []byte{4, 11, 0, 0, 0, 0}, []byte{4, 10, 0, 0, 0, 0}, []byte{4, 9, 0, 0, 0, 0}, []byte{4, 8, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (5 * 6),
		isOverflow:           true,

		isLeaf: true,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage

	found, _, page, parents := search(0, 7, []*PageParsed{})

	if !found {
		t.Errorf("expected to find rowid 7")
	}

	if page.pageNumber != 2 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 2, page.pageNumber)
	}
	if len(parents) != 1 {
		t.Errorf("Expected to be only 1 parent1, instead we got: %v", len(parents))
	}

	if !reflect.DeepEqual(parents[0].pageNumber, 0) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[0].pageNumber, 0)
	}
	found, _, page, parents = search(0, 12, []*PageParsed{})

	if !found {
		t.Errorf("expected to find result rowid 12")
	}

	if page.pageNumber != 3 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 3, page.pageNumber)
	}
	if len(parents) != 1 {
		t.Errorf("Expected to be only 1 parent1, instead we got: %v", len(parents))
	}

	if !reflect.DeepEqual(parents[0].pageNumber, 0) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[0].pageNumber, 0)
	}

	found, _, page, parents = search(0, 4, []*PageParsed{})

	if !found {
		t.Errorf("expected to find result rowid 4")
	}
	if len(parents) != 1 {
		t.Errorf("Expected to be only 1 parent1, instead we got: %v", len(parents))
	}

	if !reflect.DeepEqual(parents[0].pageNumber, 0) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[0].pageNumber, 0)
	}

	if page.pageNumber != 1 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 1, page.pageNumber)
	}

	found, _, page, parents = search(0, 1, []*PageParsed{})

	if !found {
		t.Errorf("expected to find result rowid 1")
	}

	if page.pageNumber != 1 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 1, page.pageNumber)
	}
	if len(parents) != 1 {
		t.Errorf("Expected to be only 1 parent1, instead we got: %v", len(parents))
	}

	if !reflect.DeepEqual(parents[0].pageNumber, 0) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[0].pageNumber, 0)
	}
}

func TestShouldNotFindResultSingleInteriorPlusLeafPage(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:         header(),
		dbHeaderSize:     100,
		pageNumber:       0,
		cellAreaParsed:   [][]byte{{0, 0, 0, 2, 0, 7}, {0, 0, 0, 1, 0, 4}},
		btreeType:        int(TableBtreeInteriorCell),
		rightMostpointer: []byte{0, 0, 0, 3},
		cellArea:         []byte{0, 0, 0, 2, 0, 7, 0, 0, 0, 1, 0, 4},

		startCellContentArea: PageSize - 12,
		isOverflow:           false,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 4, 0, 0, 0, 0}, []byte{4, 3, 0, 0, 0, 0}, []byte{4, 2, 0, 0, 0, 0}, []byte{4, 1, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 7, 0, 0, 0, 0}, []byte{4, 6, 0, 0, 0, 0}, []byte{4, 5, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (3 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 12, 0, 0, 0, 0}, []byte{4, 11, 0, 0, 0, 0}, []byte{4, 10, 0, 0, 0, 0}, []byte{4, 9, 0, 0, 0, 0}, []byte{4, 8, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (5 * 6),
		isOverflow:           true,

		isLeaf: true,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage

	found, _, page, parents := search(0, 13, []*PageParsed{})

	if found {
		t.Errorf("Shouldn't find rowid 13 in any of pages")
	}

	if page.pageNumber != 3 {
		t.Errorf("expected to insert new value at page 3, so be return value, insted we got: %v", page.pageNumber)
	}

	if len(parents) != 1 {
		t.Errorf("Expected to be only 1 parent1, instead we got: %v", len(parents))
	}

	if !reflect.DeepEqual(parents[0].pageNumber, 0) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[0].pageNumber, 0)
	}

}

func TestShouldFindResultMultipleInteriorPlusLeafPage(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:         header(),
		dbHeaderSize:     100,
		pageNumber:       0,
		cellAreaParsed:   [][]byte{{0, 0, 0, 6, 0, 7}},
		btreeType:        int(TableBtreeInteriorCell),
		rightMostpointer: []byte{0, 0, 0, 7},
		cellArea:         []byte{0, 0, 0, 6, 0, 7},

		startCellContentArea: PageSize - 6,
		isOverflow:           false,
	}

	var sixthPage = PageParsed{
		dbHeader:         DbHeader{},
		dbHeaderSize:     0,
		pageNumber:       6,
		cellAreaParsed:   [][]byte{{0, 0, 0, 1, 0, 4}},
		btreeType:        int(TableBtreeInteriorCell),
		rightMostpointer: []byte{0, 0, 0, 2},
		cellArea:         []byte{0, 0, 0, 1, 0, 4},

		startCellContentArea: PageSize - 6,
		isOverflow:           false,
	}

	var seventhPage = PageParsed{
		dbHeader:         DbHeader{},
		dbHeaderSize:     0,
		pageNumber:       7,
		cellAreaParsed:   [][]byte{{0, 0, 0, 3, 0, 12}},
		btreeType:        int(TableBtreeInteriorCell),
		rightMostpointer: []byte{0, 0, 0, 4},
		cellArea:         []byte{0, 0, 0, 3, 0, 12},

		startCellContentArea: PageSize - 6,
		isOverflow:           false,
	}

	// 1-4, 5-7, 8-12, 16-13

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 4, 0, 0, 0, 0}, []byte{4, 3, 0, 0, 0, 0}, []byte{4, 2, 0, 0, 0, 0}, []byte{4, 1, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 7, 0, 0, 0, 0}, []byte{4, 6, 0, 0, 0, 0}, []byte{4, 5, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (3 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 12, 0, 0, 0, 0}, []byte{4, 11, 0, 0, 0, 0}, []byte{4, 10, 0, 0, 0, 0}, []byte{4, 9, 0, 0, 0, 0}, []byte{4, 8, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (5 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var fourthPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           4,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 16, 0, 0, 0, 0}, []byte{4, 15, 0, 0, 0, 0}, []byte{4, 14, 0, 0, 0, 0}, []byte{4, 13, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,

		isLeaf: true,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage
	memoryPages[fourthPage.pageNumber] = fourthPage
	memoryPages[5] = CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, 5, nil)
	memoryPages[sixthPage.pageNumber] = sixthPage
	memoryPages[seventhPage.pageNumber] = seventhPage

	found, _, pageFound, parents := search(0, 7, []*PageParsed{})

	if !found {
		t.Errorf("expected to find rowid 7")
	}

	if pageFound.pageNumber != 2 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 2, pageFound.pageNumber)
	}

	if len(parents) != 2 {
		t.Errorf("Expected to be only 2 parents, instead we got: %v", len(parents))
	}

	if !reflect.DeepEqual(parents[0].pageNumber, 0) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[0].pageNumber, 0)
	}
	if !reflect.DeepEqual(parents[1].pageNumber, 6) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[1].pageNumber, 6)
	}

	found, _, pageFound, parents = search(0, 12, []*PageParsed{})

	if !found {
		t.Errorf("expected to find result rowid 12")
	}

	if pageFound.pageNumber != 3 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 3, pageFound.pageNumber)
	}
	if len(parents) != 2 {
		t.Errorf("Expected to be only 2 parents, instead we got: %v", len(parents))
	}

	if !reflect.DeepEqual(parents[0].pageNumber, 0) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[0].pageNumber, 0)
	}
	if !reflect.DeepEqual(parents[1].pageNumber, 7) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[1].pageNumber, 7)
	}
	found, _, pageFound, parents = search(0, 4, []*PageParsed{})

	if !found {
		t.Errorf("expected to find result rowid 4")
	}

	if pageFound.pageNumber != 1 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 1, pageFound.pageNumber)
	}

	if len(parents) != 2 {
		t.Errorf("Expected to be only 2 parents, instead we got: %v", len(parents))
	}

	if !reflect.DeepEqual(parents[0].pageNumber, 0) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[0].pageNumber, 0)
	}
	if !reflect.DeepEqual(parents[1].pageNumber, 6) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[1].pageNumber, 6)
	}
	found, _, pageFound, parents = search(0, 1, []*PageParsed{})

	if !found {
		t.Errorf("expected to find result rowid 1")
	}

	if pageFound.pageNumber != 1 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 1, pageFound.pageNumber)
	}
	if len(parents) != 2 {
		t.Errorf("Expected to be only 2 parents, instead we got: %v", len(parents))
	}

	if !reflect.DeepEqual(parents[0].pageNumber, 0) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[0].pageNumber, 0)
	}
	if !reflect.DeepEqual(parents[1].pageNumber, 6) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[1].pageNumber, 6)
	}
	found, _, pageFound, parents = search(0, 16, []*PageParsed{})

	if !found {
		t.Errorf("expected to find result rowid 16")
	}

	if pageFound.pageNumber != 4 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 4, pageFound.pageNumber)
	}
	if len(parents) != 2 {
		t.Errorf("Expected to be only 2 parents, instead we got: %v", len(parents))
	}

	if !reflect.DeepEqual(parents[0].pageNumber, 0) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[0].pageNumber, 0)
	}
	if !reflect.DeepEqual(parents[1].pageNumber, 7) {
		t.Errorf("expect parent to be page: %v, instead we got: %v", parents[1].pageNumber, 7)
	}
}

func TestUpdateDividerSameAmount(t *testing.T) {
	clearDbFile("test")
	cellAreaParsed := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 8}, {0, 0, 0, 1, 0, 5}}
	cellArea := []byte{0, 0, 0, 4, 0, 16, 0, 0, 0, 3, 0, 12, 0, 0, 0, 2, 0, 8, 0, 0, 0, 1, 0, 5}
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		numberofCells:        len(cellAreaParsed),
		cellAreaSize:         len(cellArea),
		cellAreaParsed:       cellAreaParsed,
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 4},
		cellArea:             cellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
	}

	cells := []Cell{{rowId: 9, pageNumber: 2}, {rowId: 13, pageNumber: 3}}

	modifyDivider(&zeroPage, cells, 6, 6*3, &zeroPage.dbHeader, []*PageParsed{})

	cellAreaParsedExpected := cellAreaParsed
	cellAreaParsedExpected[1] = []byte{0, 0, 0, 3, 0, 13}
	cellAreaParsedExpected[2] = []byte{0, 0, 0, 2, 0, 9}

	if !reflect.DeepEqual(cellAreaParsedExpected, zeroPage.cellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", cellAreaParsedExpected, zeroPage.cellAreaParsed)
	}

	if zeroPage.cellAreaSize != len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("cell area size should be: %v, got: %v", len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPage.cellAreaSize)
	}

	if zeroPage.numberofCells != len(cellAreaParsedExpected) {
		t.Errorf("number of cell should be: %v, got: %v", len(cellAreaParsedExpected), zeroPage.numberofCells)
	}

	if zeroPage.startCellContentArea != PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("start of cell content should be: %v, got: %v", PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPage.startCellContentArea)
	}

	expectedRightMostPointer := []byte{0, 0, 0, 4}

	if !reflect.DeepEqual(zeroPage.rightMostpointer, expectedRightMostPointer) {
		t.Errorf("expected right most pointer to be: %v, instead we got: %v", expectedRightMostPointer, zeroPage.rightMostpointer)
	}
}

func TestUpdateDividerLessAmount(t *testing.T) {
	clearDbFile("test")
	cellAreaParsed := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 8}, {0, 0, 0, 1, 0, 5}}
	cellArea := []byte{0, 0, 0, 4, 0, 16, 0, 0, 0, 3, 0, 12, 0, 0, 0, 2, 0, 8, 0, 0, 0, 1, 0, 5}
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		cellAreaParsed:       cellAreaParsed,
		btreeType:            int(TableBtreeInteriorCell),
		numberofCells:        len(cellAreaParsed),
		cellAreaSize:         len(cellArea),
		rightMostpointer:     []byte{0, 0, 0, 4},
		cellArea:             cellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
	}

	cells := []Cell{{rowId: 9, pageNumber: 2}}

	modifyDivider(&zeroPage, cells, 6, 6*3, &zeroPage.dbHeader, []*PageParsed{})

	cellAreaParsedExpected := [][]byte{cellAreaParsed[0]}
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 2, 0, 9})
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[3])

	if !reflect.DeepEqual(cellAreaParsedExpected, zeroPage.cellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", cellAreaParsedExpected, zeroPage.cellAreaParsed)
	}

	if zeroPage.cellAreaSize != len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("cell area size should be: %v, got: %v", len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPage.cellAreaSize)
	}

	if zeroPage.numberofCells != len(cellAreaParsedExpected) {
		t.Errorf("number of cell should be: %v, got: %v", len(cellAreaParsedExpected), zeroPage.numberofCells)
	}

	if zeroPage.startCellContentArea != PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("start of cell content should be: %v, got: %v", PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPage.startCellContentArea)
	}

	expectedRightMostPointer := []byte{0, 0, 0, 4}

	if !reflect.DeepEqual(zeroPage.rightMostpointer, expectedRightMostPointer) {
		t.Errorf("expected right most pointer to be: %v, instead we got: %v", expectedRightMostPointer, zeroPage.rightMostpointer)
	}
}

func TestUpdateDividerGreaterAmount(t *testing.T) {
	clearDbFile("test")
	cellAreaParsed := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 8}, {0, 0, 0, 1, 0, 5}}
	cellArea := []byte{0, 0, 0, 4, 0, 16, 0, 0, 0, 3, 0, 12, 0, 0, 0, 2, 0, 8, 0, 0, 0, 1, 0, 5}
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		cellAreaParsed:       cellAreaParsed,
		btreeType:            int(TableBtreeInteriorCell),
		numberofCells:        len(cellAreaParsed),
		cellAreaSize:         len(cellArea),
		rightMostpointer:     []byte{0, 0, 0, 4},
		cellArea:             cellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
	}

	cells := []Cell{{rowId: 9, pageNumber: 2}, {rowId: 13, pageNumber: 3}, {rowId: 14, pageNumber: 5}}

	modifyDivider(&zeroPage, cells, 6, 6*3, &zeroPage.dbHeader, []*PageParsed{})

	cellAreaParsedExpected := [][]byte{cellAreaParsed[0]}
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 5, 0, 14})
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 3, 0, 13})
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 2, 0, 9})
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[3])

	if !reflect.DeepEqual(cellAreaParsedExpected, zeroPage.cellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", cellAreaParsedExpected, zeroPage.cellAreaParsed)
	}

	if zeroPage.cellAreaSize != len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("cell area size should be: %v, got: %v", len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPage.cellAreaSize)
	}

	if zeroPage.numberofCells != len(cellAreaParsedExpected) {
		t.Errorf("number of cell should be: %v, got: %v", len(cellAreaParsedExpected), zeroPage.numberofCells)
	}

	if zeroPage.startCellContentArea != PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("start of cell content should be: %v, got: %v", PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPage.startCellContentArea)
	}
}

func TestUpdateTestRightMostPointerWhenUpdatePointerWithNewPage(t *testing.T) {
	clearDbFile("test")
	cellAreaParsed := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 8}, {0, 0, 0, 1, 0, 5}}
	cellArea := []byte{0, 0, 0, 4, 0, 16, 0, 0, 0, 3, 0, 12, 0, 0, 0, 2, 0, 8, 0, 0, 0, 1, 0, 5}
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		cellAreaParsed:       cellAreaParsed,
		btreeType:            int(TableBtreeInteriorCell),
		numberofCells:        len(cellAreaParsed),
		cellAreaSize:         len(cellArea),
		rightMostpointer:     []byte{0, 0, 0, 4},
		cellArea:             cellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
	}

	cells := []Cell{{rowId: 18, pageNumber: 5}}

	modifyDivider(&zeroPage, cells, 0, 6, &zeroPage.dbHeader, []*PageParsed{})

	cellAreaParsedExpected := [][]byte{[]byte{0, 0, 0, 5, 0, 18}}
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[1])
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[2])
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[3])

	if !reflect.DeepEqual(cellAreaParsedExpected, zeroPage.cellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", cellAreaParsedExpected, zeroPage.cellAreaParsed)
	}

	if zeroPage.cellAreaSize != len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("cell area size should be: %v, got: %v", len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPage.cellAreaSize)
	}

	if zeroPage.numberofCells != len(cellAreaParsedExpected) {
		t.Errorf("number of cell should be: %v, got: %v", len(cellAreaParsedExpected), zeroPage.numberofCells)
	}

	if zeroPage.startCellContentArea != PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("start of cell content should be: %v, got: %v", PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPage.startCellContentArea)
	}

	expectedRightMostPointer := []byte{0, 0, 0, 5}

	if !reflect.DeepEqual(zeroPage.rightMostpointer, expectedRightMostPointer) {
		t.Errorf("expected right most pointer to be: %v, instead we got: %v", expectedRightMostPointer, zeroPage.rightMostpointer)
	}
}

func utilsTestContent(t *testing.T, expectedCellAreaParsed [][]byte, expectedBtreeType BtreeType, expectedPageNumber int, page PageParsed) {
	cellArea := []byte{}
	expectedPointers := []byte{}
	startContent := PageSize
	expecteRightMostPointer := []byte{}
	isLeaf := true
	for _, v := range expectedCellAreaParsed {
		cellArea = append(cellArea, v...)
		startContent -= len(v)
		expectedPointers = append(expectedPointers, intToBinary(startContent, 2)...)
	}
	if page.numberofCells != len(expectedCellAreaParsed) {
		t.Errorf("Expected cell number to be: %v, got: %v", len(expectedCellAreaParsed), page.numberofCells)
	}

	if !reflect.DeepEqual(expectedCellAreaParsed, page.cellAreaParsed) {
		t.Errorf("Expected cell area parsed to be: %v, got: %v", expectedCellAreaParsed, page.cellAreaParsed)
	}

	if !reflect.DeepEqual(cellArea, page.cellArea) {
		t.Errorf("Expected cell area to be: %v, got: %v", cellArea, page.cellArea)
	}
	if !reflect.DeepEqual(expectedPointers, page.pointers) {
		t.Errorf("Expected pointters to be: %v, got: %v", expectedPointers, page.pointers)
	}

	if page.startCellContentArea != startContent {
		t.Errorf("Expected start cell area to be: %v, got: %v", startContent, page.cellArea)
	}

	if page.cellAreaSize != len(cellArea) {
		t.Errorf("Expected cell area size to be: %v, got: %v", len(cellArea), page.cellAreaSize)
	}

	if expectedBtreeType == TableBtreeInteriorCell {
		if len(expectedCellAreaParsed) > 0 {
			expecteRightMostPointer = expectedCellAreaParsed[0][:4]
		}
		isLeaf = false
	}

	if page.btreeType != int(expectedBtreeType) {
		t.Errorf("Expected btree type to be: %v, got :%v", expectedBtreeType, page.btreeType)
	}

	if page.pageNumber != expectedPageNumber {
		t.Errorf("Expected page number to be: %v, got: %v", expectedPageNumber, page.pageNumber)
	}

	if page.isLeaf != isLeaf {
		t.Errorf("Expected value is leaf to be: %v, got :%v", isLeaf, page.isLeaf)
	}

	if !reflect.DeepEqual(page.rightMostpointer, expecteRightMostPointer) {
		t.Errorf("expected right most pointer to be: %v, got: %v", expecteRightMostPointer, page.rightMostpointer)
	}

	expectedIsOverflow := !page.isSpace()

	if page.isOverflow != expectedIsOverflow {
		t.Errorf("Expected is overflow value to be: %v, got: %v", expectedIsOverflow, page.isOverflow)
	}

}

func TestUpdateDividerWithParentUpdate(t *testing.T) {
	zeroPagecellAreaParsed := [][]byte{{0, 0, 0, 1, 0, 16}}
	zeroPagecellArea := []byte{0, 0, 0, 1, 0, 16}
	zeroPagePointers := append([]byte{}, intToBinary(PageSize-len(zeroPagecellArea), 2)...)
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		cellAreaParsed:       zeroPagecellAreaParsed,
		btreeType:            int(TableBtreeInteriorCell),
		numberofCells:        len(zeroPagecellAreaParsed),
		cellAreaSize:         len(zeroPagecellArea),
		rightMostpointer:     []byte{0, 0, 0, 4},
		cellArea:             zeroPagecellArea,
		pointers:             zeroPagePointers,
		isLeaf:               false,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
	}
	clearDbFile("test")
	firstPagecellAreaParsed := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 8}, {0, 0, 0, 8, 0, 5}}
	firstPagecellArea := []byte{0, 0, 0, 4, 0, 16, 0, 0, 0, 3, 0, 12, 0, 0, 0, 2, 0, 8, 0, 0, 0, 8, 0, 5}
	var firstPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           1,
		cellAreaParsed:       firstPagecellAreaParsed,
		btreeType:            int(TableBtreeInteriorCell),
		numberofCells:        len(firstPagecellAreaParsed),
		cellAreaSize:         len(firstPagecellArea),
		rightMostpointer:     []byte{0, 0, 0, 4},
		cellArea:             firstPagecellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
	}

	cells := []Cell{{rowId: 18, pageNumber: 5}}

	modifyDivider(&firstPage, cells, 0, 6, &zeroPage.dbHeader, []*PageParsed{&zeroPage})

	utilsTestContent(t, [][]byte{{0, 0, 0, 1, 0, 18}}, TableBtreeInteriorCell, 0, zeroPage)

}

func TestUpdateDividerWithParentandGrandParentUpdate(t *testing.T) {
	clearDbFile("test")
	zeroPagecellAreaParsed := [][]byte{{0, 0, 0, 1, 0, 16}}
	zeroPagecellArea := []byte{0, 0, 0, 1, 0, 16}
	zeroPagePointers := append([]byte{}, intToBinary(PageSize-len(zeroPagecellArea), 2)...)
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		cellAreaParsed:       zeroPagecellAreaParsed,
		btreeType:            int(TableBtreeInteriorCell),
		numberofCells:        len(zeroPagecellAreaParsed),
		cellAreaSize:         len(zeroPagecellArea),
		rightMostpointer:     []byte{0, 0, 0, 4},
		cellArea:             zeroPagecellArea,
		pointers:             zeroPagePointers,
		isLeaf:               false,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
	}

	firstPagecellAreaParsed := [][]byte{{0, 0, 0, 2, 0, 16}}
	firstPagecellArea := []byte{0, 0, 0, 2, 0, 16}
	firstPagePointers := append([]byte{}, intToBinary(PageSize-len(firstPagecellArea), 2)...)
	var firstPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           1,
		cellAreaParsed:       firstPagecellAreaParsed,
		btreeType:            int(TableBtreeInteriorCell),
		numberofCells:        len(firstPagecellAreaParsed),
		cellAreaSize:         len(firstPagecellArea),
		rightMostpointer:     []byte{0, 0, 0, 4},
		cellArea:             firstPagecellArea,
		pointers:             firstPagePointers,
		isLeaf:               false,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
	}

	secondPagecellAreaParsed := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 9, 0, 8}, {0, 0, 0, 8, 0, 5}}
	secondPagecellArea := []byte{0, 0, 0, 4, 0, 16, 0, 0, 0, 3, 0, 12, 0, 0, 0, 9, 0, 8, 0, 0, 0, 8, 0, 5}
	var secondPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           2,
		cellAreaParsed:       secondPagecellAreaParsed,
		btreeType:            int(TableBtreeInteriorCell),
		numberofCells:        len(secondPagecellAreaParsed),
		cellAreaSize:         len(secondPagecellArea),
		rightMostpointer:     []byte{0, 0, 0, 4},
		cellArea:             secondPagecellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
	}

	cells := []Cell{{rowId: 18, pageNumber: 5}}

	modifyDivider(&secondPage, cells, 0, 6, &zeroPage.dbHeader, []*PageParsed{&zeroPage, &firstPage})

	utilsTestContent(t, [][]byte{{0, 0, 0, 1, 0, 18}}, TableBtreeInteriorCell, 0, zeroPage)
	utilsTestContent(t, [][]byte{{0, 0, 0, 2, 0, 18}}, TableBtreeInteriorCell, 1, firstPage)

}

func TestFindSiblingOnlyRightSiblingsAsLastPointer(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		btreePageHeaderSize:  12,
		pageNumber:           0,
		cellAreaParsed:       [][]byte{{0, 0, 0, 2, 0, 2}, {0, 0, 0, 1, 0, 1}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 2},
		cellArea:             []byte{0, 0, 0, 2, 0, 2, 0, 0, 0, 1, 0, 1},
		cellAreaSize:         12,
		numberofCells:        2,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,

		isLeaf: false,
	}

	cell := creareARowItem(100, 2)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		cellArea:             cell.data,
		startCellContentArea: PageSize - cell.dataLength,
		cellAreaSize:         cell.dataLength,
		numberofCells:        1,
		isOverflow:           false,

		isLeaf: true,
	}

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		isOverflow:           false,

		isLeaf: true,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage

	leftSibling, rightSibling := zeroPage.findSiblings(firstPage)

	if leftSibling != nil {
		t.Errorf("Expected left sibling to be nil, got: %v", leftSibling)
	}

	if rightSibling == nil || rightSibling.pageNumber != 2 {
		t.Errorf("expected right sibling to be page number 2")
	}
}

func TestFindSiblingOnlyLeftSiblingsAsFirst(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		btreePageHeaderSize:  12,
		pageNumber:           0,
		cellAreaParsed:       [][]byte{{0, 0, 0, 2, 0, 2}, {0, 0, 0, 1, 0, 1}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 2},
		cellArea:             []byte{0, 0, 0, 2, 0, 2, 0, 0, 0, 1, 0, 1},
		cellAreaSize:         12,
		numberofCells:        2,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,

		isLeaf: false,
	}

	cell := creareARowItem(100, 3)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		numberofCells:        1,
		isOverflow:           false,

		isLeaf: true,
	}

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		cellArea:             cell.data,
		startCellContentArea: PageSize - cell.dataLength,
		cellAreaSize:         cell.dataLength,
		numberofCells:        1,
		isOverflow:           false,

		isLeaf: true,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage

	leftSibling, rightSibling := zeroPage.findSiblings(secondPage)

	if rightSibling != nil {
		t.Errorf("Expected left sibling to be nil, got: %v", rightSibling)
	}

	if leftSibling == nil || leftSibling.pageNumber != 1 {
		t.Errorf("expected right sibling to be page number 1")
	}
}

func TestFindSiblingBothSiblings(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		btreePageHeaderSize:  12,
		pageNumber:           0,
		cellAreaParsed:       [][]byte{{0, 0, 0, 3, 0, 3}, {0, 0, 0, 2, 0, 2}, {0, 0, 0, 1, 0, 1}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 3},
		cellArea:             []byte{0, 0, 0, 3, 0, 3, 0, 0, 0, 2, 0, 2, 0, 0, 0, 1, 0, 1},
		cellAreaSize:         12,
		numberofCells:        2,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,

		isLeaf: false,
	}

	cell := creareARowItem(100, 3)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		numberofCells:        1,
		isOverflow:           false,

		isLeaf: true,
	}

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		cellArea:             cell.data,
		startCellContentArea: PageSize - cell.dataLength,
		cellAreaSize:         cell.dataLength,
		numberofCells:        1,
		isOverflow:           false,

		isLeaf: true,
	}

	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		numberofCells:        1,
		isOverflow:           false,

		isLeaf: true,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage

	leftSibling, rightSibling := zeroPage.findSiblings(secondPage)

	if rightSibling == nil || rightSibling.pageNumber != 3 {
		t.Errorf("Expected right siblings to be page number 3")
	}

	if leftSibling == nil || leftSibling.pageNumber != 1 {
		t.Errorf("expected right sibling to be page number 1")
	}
}

func TestFindSiblingNoSiblings(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		btreePageHeaderSize:  12,
		pageNumber:           0,
		cellAreaParsed:       [][]byte{{0, 0, 0, 1, 0, 1}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 1},
		cellArea:             []byte{0, 0, 0, 1, 0, 1},
		cellAreaSize:         12,
		numberofCells:        2,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,

		isLeaf: false,
	}

	cell := creareARowItem(100, 2)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		cellArea:             cell.data,
		startCellContentArea: PageSize - cell.dataLength,
		cellAreaSize:         cell.dataLength,
		numberofCells:        1,
		isOverflow:           false,

		isLeaf: true,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage

	leftSibling, rightSibling := zeroPage.findSiblings(firstPage)

	if leftSibling != nil {
		t.Errorf("Expected left sibling to be nil, got: %v", leftSibling)
	}

	if rightSibling != nil {
		t.Errorf("Expected right sibling to be nil, got: %v", leftSibling)
	}
}

func TestUpdateData(t *testing.T) {
	clearDbFile("test")

	cell := creareARowItem(100, 1)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		cellArea:             cell.data,
		startCellContentArea: PageSize - cell.dataLength,
		pointers:             intToBinary(PageSize-cell.dataLength, 2),
		cellAreaSize:         cell.dataLength,
		numberofCells:        1,
		isOverflow:           false,

		isLeaf: true,
	}

	newCell := createCell(TableBtreeLeafCell, 0, "Alice")
	firstPage.updateParsedCells(newCell, 0)

	if !reflect.DeepEqual(firstPage.cellAreaParsed[0], newCell.data) {
		t.Errorf("Expected cell area parsed to be: %v, instead we got: %v", newCell.data, firstPage.cellAreaParsed[0])
	}

	if firstPage.numberofCells != 1 {
		t.Errorf("Number of cell should stay as it was (1) instead we got :%v", firstPage.numberofCells)
	}

	if firstPage.startCellContentArea != PageSize-newCell.dataLength {
		t.Errorf("cell area should start at: %v, instead we got: %v", PageSize-newCell.dataLength, firstPage.startCellContentArea)
	}

	if !reflect.DeepEqual(firstPage.cellArea, firstPage.cellAreaParsed[0]) {
		t.Errorf("Expected cell area to be: %v, got: %v", firstPage.cellAreaParsed[0], firstPage.cellArea)
	}

	if firstPage.cellAreaSize != newCell.dataLength {
		t.Errorf("Expected cell area length to be: %v, instead we got: %v", newCell.dataLength, firstPage.cellAreaSize)
	}

	if !reflect.DeepEqual(firstPage.pointers, intToBinary(PageSize-newCell.dataLength, 2)) {
		t.Errorf("Expected pointers to be: %v, instead got: %v", intToBinary(PageSize-newCell.dataLength, 2), firstPage.pointers)
	}
}

func TestInsertNewData(t *testing.T) {
	clearDbFile("test")

	cell := creareARowItem(100, 1)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		numberofCells:        1,
		cellAreaParsed:       [][]byte{{0, 0, 0, 1, 0, 0}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{},
		cellArea:             []byte{0, 0, 0, 1, 0, 0},
		isLeaf:               true,
		startCellContentArea: PageSize - 6,
		isOverflow:           true,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		cellArea:             cell.data,
		startCellContentArea: PageSize - cell.dataLength,
		pointers:             intToBinary(PageSize-cell.dataLength, 2),
		cellAreaSize:         cell.dataLength,
		numberofCells:        1,
		isOverflow:           false,

		isLeaf: true,
	}

	newCell := createCell(TableBtreeLeafCell, 1, "Alice")

	firstPage.insertData(newCell, &zeroPage.dbHeader, []*PageParsed{&zeroPage})

	expectedCellArea := append([]byte{}, newCell.data...)
	expectedCellArea = append(expectedCellArea, cell.data...)

	parsedCellArea := dbReadparseCellArea(byte(TableBtreeLeafCell), expectedCellArea)

	if firstPage.numberofCells != 2 {
		t.Errorf("expected first page to have 2 cells, instead we got: %v", firstPage.numberofCells)
	}

	if !reflect.DeepEqual(firstPage.cellAreaParsed, parsedCellArea) {
		t.Errorf("Expected cell area parsed to be: %v, instead we got: %v", parsedCellArea, firstPage.cellAreaParsed)
	}

	if firstPage.startCellContentArea != PageSize-newCell.dataLength-cell.dataLength {
		t.Errorf("cell area should start at: %v, instead we got: %v", PageSize-newCell.dataLength-cell.dataLength, firstPage.startCellContentArea)
	}

	if !reflect.DeepEqual(firstPage.cellArea, expectedCellArea) {
		t.Errorf("Expected cell area to be: %v, got: %v", expectedCellArea, firstPage.cellArea)
	}

	expectedPointer := append([]byte{}, intToBinary(PageSize-newCell.dataLength, 2)...)
	expectedPointer = append(expectedPointer, intToBinary(PageSize-cell.dataLength-newCell.dataLength, 2)...)

	if !reflect.DeepEqual(firstPage.pointers, expectedPointer) {
		t.Errorf("Expected pointers to be: %v, instead got: %v", expectedPointer, firstPage.pointers)
	}
}

func TestBalancingSplitRootIntTwoChildren(t *testing.T) {
	clearDbFile("test")
	cellAreaParsed := [][]byte{{7, 4, 2, 23, 65, 108, 105, 99, 101}, {7, 3, 2, 23, 65, 108, 105, 99, 101}, {7, 2, 2, 23, 65, 108, 105, 99, 101}, {7, 1, 2, 23, 65, 108, 105, 99, 101}}
	cellArea := []byte{7, 4, 2, 23, 65, 108, 105, 99, 101, 7, 3, 2, 23, 65, 108, 105, 99, 101, 7, 2, 2, 23, 65, 108, 105, 99, 101, 7, 1, 2, 23, 65, 108, 105, 99, 101}
	var zeroPage = PageParsed{
		dbHeader: DbHeader{
			dbSizeInPages: 2,
		},
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		btreePageHeaderSize:  8,
		numberofCells:        len(cellAreaParsed),
		cellAreaParsed:       cellAreaParsed,
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		cellArea:             cellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 9*4,
		isOverflow:           true,
	}

	PageSize = 9*2 + 2*2 + 8

	btree := BtreeStruct{
		softPages: map[int]PageParsed{},
	}

	btree.balancingForNode(&firstPage, []*PageParsed{}, &zeroPage.dbHeader)

	zeroPageSaved := btree.softPages[1]
	firstPageSaved := btree.softPages[2]
	secondPageSaved := btree.softPages[3]

	utilsTestContent(t, [][]byte{{0, 0, 0, 3, 0, 4}, {0, 0, 0, 2, 0, 2}}, TableBtreeInteriorCell, 1, zeroPageSaved)

	firstPageExpectedCellAreaParsed := append([][]byte{}, cellAreaParsed[2])
	firstPageExpectedCellAreaParsed = append(firstPageExpectedCellAreaParsed, cellAreaParsed[3])
	utilsTestContent(t, firstPageExpectedCellAreaParsed, TableBtreeLeafCell, 2, firstPageSaved)

	secondPageExpectedCellAreaParsed := append([][]byte{}, cellAreaParsed[0])
	secondPageExpectedCellAreaParsed = append(secondPageExpectedCellAreaParsed, cellAreaParsed[1])

	utilsTestContent(t, secondPageExpectedCellAreaParsed, TableBtreeLeafCell, 3, secondPageSaved)

}

func TestBalancingSplitRootIntOneChild(t *testing.T) {
	clearDbFile("test")
	cellAreaParsed := [][]byte{{7, 4, 2, 23, 65, 108, 105, 99, 101}, {7, 3, 2, 23, 65, 108, 105, 99, 101}, {7, 2, 2, 23, 65, 108, 105, 99, 101}, {7, 1, 2, 23, 65, 108, 105, 99, 101}}
	cellArea := []byte{7, 4, 2, 23, 65, 108, 105, 99, 101, 7, 3, 2, 23, 65, 108, 105, 99, 101, 7, 2, 2, 23, 65, 108, 105, 99, 101, 7, 1, 2, 23, 65, 108, 105, 99, 101}
	pointers := append([]byte{}, intToBinary(PageSize-9, 2)...)
	pointers = append(pointers, intToBinary(PageSize-9*2, 2)...)
	pointers = append(pointers, intToBinary(PageSize-9*3, 2)...)
	pointers = append(pointers, intToBinary(PageSize-9*4, 2)...)
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		numberofCells:        len(cellAreaParsed),
		cellAreaParsed:       cellAreaParsed,
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		cellArea:             cellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 9*4,
		isOverflow:           true,
		cellAreaSize:         len(cellArea),
		pointers:             pointers,
	}

	server := ServerStruct{
		header: zeroPage.dbHeader,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage

	btree := BtreeStruct{
		softPages: map[int]PageParsed{},
	}
	btree.balancingForNode(&zeroPage, []*PageParsed{}, &server.header)

	zeroPageSaved := btree.softPages[0]
	firstPageSaved := btree.softPages[1]

	if server.header.dbSizeInPages != 2 {
		t.Errorf("Expected number of pages to be 2, got: %v", server.header.dbSizeInPages)
	}

	utilsTestContent(t, [][]byte{{0, 0, 0, 1, 0, 4}}, TableBtreeInteriorCell, 0, zeroPageSaved)
	utilsTestContent(t, cellAreaParsed, TableBtreeLeafCell, 1, firstPageSaved)

	if !reflect.DeepEqual(firstPageSaved.cellAreaParsed, cellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", cellAreaParsed, firstPageSaved.cellAreaParsed)
	}

	if firstPageSaved.numberofCells != 4 {
		t.Errorf("expected first page to have 4 cells, instead we have: %v", firstPageSaved.numberofCells)
	}

	if !reflect.DeepEqual(firstPageSaved.cellArea, cellArea) {
		t.Errorf("Expected first page cell area to be: %v, instead we got: %v", firstPageSaved.cellArea, cellArea)
	}

}

var PointerToDataLength = 2

type BtreeHeaderSize int

var (
	BtreeHeaderSizeLeafTable BtreeHeaderSize = 8
)

// // work on this
func TestBalancingSplitOneLeaftPageIntoTwo(t *testing.T) {
	clearDbFile("test")
	cellAreaParsed := [][]byte{{7, 4, 2, 23, 65, 108, 105, 99, 101}, {7, 3, 2, 23, 65, 108, 105, 99, 101}, {7, 2, 2, 23, 65, 108, 105, 99, 101}, {7, 1, 2, 23, 65, 108, 105, 99, 101}}
	PageSize = len(cellAreaParsed[0])*2 + PointerToDataLength*2 + int(BtreeHeaderSizeLeafTable)

	server := ServerStruct{
		header: DbHeader{},
	}

	CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{}, server.header.assignNewPage(), &server.header)

	firstPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 4}}, server.header.assignNewPage(), nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), cellAreaParsed, server.header.assignNewPage(), nil)

	btree := BtreeStruct{
		softPages: map[int]PageParsed{},
	}
	btree.balancingForNode(&secondPage, []*PageParsed{&firstPage}, &server.header)

	zeroPageSaved := btree.softPages[1]
	firstPageSaved := btree.softPages[2]
	secondPageSaved := btree.softPages[3]

	if server.header.dbSizeInPages != 4 {
		t.Errorf("Expected number of pages to be 4, got: %v", server.header.dbSizeInPages)
	}

	utilsTestContent(t, [][]byte{{0, 0, 0, 3, 0, 4}, {0, 0, 0, 2, 0, 2}}, TableBtreeInteriorCell, 1, zeroPageSaved)

	firstPageExpectedCellAreaParsed := append([][]byte{}, cellAreaParsed[2])
	firstPageExpectedCellAreaParsed = append(firstPageExpectedCellAreaParsed, cellAreaParsed[3])
	utilsTestContent(t, firstPageExpectedCellAreaParsed, TableBtreeLeafCell, 2, firstPageSaved)

	secondPageExpectedCellAreaParsed := append([][]byte{}, cellAreaParsed[0])
	secondPageExpectedCellAreaParsed = append(secondPageExpectedCellAreaParsed, cellAreaParsed[1])
	utilsTestContent(t, secondPageExpectedCellAreaParsed, TableBtreeLeafCell, 3, secondPageSaved)
}

func TestBinarySearchInInteriorForNewValueHighestRowId(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		btreePageHeaderSize:  12,
		pageNumber:           0,
		cellAreaParsed:       [][]byte{{0, 0, 0, 1, 0, 0}, {0, 0, 0, 2, 0, 3}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 2},
		cellArea:             []byte{0, 0, 0, 1, 0, 0, 0, 0, 0, 2, 0, 3},
		cellAreaSize:         12,
		numberofCells:        1,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,

		isLeaf: false,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage

	found, pageNumber, _ := binarySearch(zeroPage, 0, 4)

	if found {
		t.Errorf("Expected found var to be false, as we inserting new value")
	}

	if pageNumber != 2 {
		t.Errorf("expected new value to be inserted into page 2, instead we got: %v", pageNumber)
	}

}

func TestBinarySearchInLeafForExisitngValue(t *testing.T) {
	clearDbFile("test")
	cell1 := createCell(TableBtreeLeafCell, 0, "alice")
	cell2 := createCell(TableBtreeLeafCell, 1, "bob")
	cell3 := createCell(TableBtreeLeafCell, 2, "tom")

	cellArea := append([]byte{}, cell3.data...)
	cellArea = append(cellArea, cell2.data...)
	cellArea = append(cellArea, cell1.data...)
	cellAreaParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cellArea)

	var secondPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		btreePageHeaderSize:  8,
		pageNumber:           0,
		cellAreaParsed:       cellAreaParsed,
		btreeType:            int(TableBtreeLeafCell),
		cellArea:             cellArea,
		cellAreaSize:         3 * cell1.dataLength,
		numberofCells:        3,
		startCellContentArea: PageSize - 3*cell1.dataLength,
		isOverflow:           false,

		isLeaf: true,
	}

	memoryPages[secondPage.pageNumber] = secondPage

	found, pageNumber, cellAreaParsedIndex := binarySearch(secondPage, 0, 2)

	if !found {
		t.Errorf("Expected value to be found, as its already exists")
	}

	if pageNumber != 0 {
		t.Errorf("expected new value to be update on page 0, instead we got: %v", pageNumber)
	}

	if cellAreaParsedIndex != 0 {
		t.Errorf("Expected index on parsed cell area to be 0, instead we got: %v", cellAreaParsedIndex)
	}

}

func TestInsert(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		cellAreaParsed:       [][]byte{{0, 0, 0, 6, 0, 7}, {0, 0, 0, 7, 0, 16}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 7},
		cellArea:             []byte{0, 0, 0, 6, 0, 7, 0, 0, 0, 7, 0, 16},
		numberofCells:        2,
		cellAreaSize:         12,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,
	}

	server := ServerStruct{
		header: zeroPage.dbHeader,
	}

	var sixthPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           6,
		cellAreaParsed:       [][]byte{{0, 0, 0, 1, 0, 4}, {0, 0, 0, 2, 0, 7}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 2},
		cellArea:             []byte{0, 0, 0, 1, 0, 4, 0, 0, 0, 2, 0, 7},
		numberofCells:        2,
		cellAreaSize:         12,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,
	}

	var seventhPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           7,
		cellAreaParsed:       [][]byte{{0, 0, 0, 3, 0, 12}, {0, 0, 0, 4, 0, 16}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 4},
		cellArea:             []byte{0, 0, 0, 3, 0, 12, 0, 0, 0, 4, 0, 16},
		numberofCells:        2,
		cellAreaSize:         12,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 4, 0, 0, 0, 0}, []byte{4, 3, 0, 0, 0, 0}, []byte{4, 2, 0, 0, 0, 0}, []byte{4, 1, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 7, 0, 0, 0, 0}, []byte{4, 6, 0, 0, 0, 0}, []byte{4, 5, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (3 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 12, 0, 0, 0, 0}, []byte{4, 11, 0, 0, 0, 0}, []byte{4, 10, 0, 0, 0, 0}, []byte{4, 9, 0, 0, 0, 0}, []byte{4, 8, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (5 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var fourthPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           4,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 16, 0, 0, 0, 0}, []byte{4, 15, 0, 0, 0, 0}, []byte{4, 14, 0, 0, 0, 0}, []byte{4, 13, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,

		isLeaf: true,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage
	memoryPages[fourthPage.pageNumber] = fourthPage
	memoryPages[5] = CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, 5, nil)
	memoryPages[sixthPage.pageNumber] = sixthPage
	memoryPages[seventhPage.pageNumber] = seventhPage

	btree := BtreeStruct{
		softPages: map[int]PageParsed{},
	}
	cell := createCell(TableBtreeLeafCell, 17)
	node := btree.insert(17, cell, &server.header, nil)

	if node.pageNumber != 4 {
		t.Errorf("Insert values should be in page: %v, instead we got: %v", 4, node.pageNumber)
	}

	if len(node.cellAreaParsed) != 5 {
		t.Errorf("expected cell area to have 5 elements, instead we got: %v", len(node.cellAreaParsed))
	}

	if !reflect.DeepEqual(node.cellAreaParsed[0], cell.data) {
		t.Errorf("expected cell area start with newly added cell: %v, instead we got: %v", cell.data, node.cellAreaParsed[0])
	}

}

// noe lets remove hardcoded data
func TestInsertToExisting(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:         header(),
		dbHeaderSize:     100,
		pageNumber:       0,
		cellAreaParsed:   [][]byte{{0, 0, 0, 6, 0, 7}},
		btreeType:        int(TableBtreeInteriorCell),
		rightMostpointer: []byte{0, 0, 0, 7},
		cellArea:         []byte{0, 0, 0, 6, 0, 7},

		startCellContentArea: PageSize - 6,
		isOverflow:           false,
	}

	server := ServerStruct{
		header: zeroPage.dbHeader,
	}

	var sixthPage = PageParsed{
		dbHeader:         DbHeader{},
		dbHeaderSize:     0,
		pageNumber:       6,
		cellAreaParsed:   [][]byte{{0, 0, 0, 1, 0, 4}},
		btreeType:        int(TableBtreeInteriorCell),
		rightMostpointer: []byte{0, 0, 0, 2},
		cellArea:         []byte{0, 0, 0, 1, 0, 4},

		startCellContentArea: PageSize - 6,
		isOverflow:           false,
	}

	var seventhPage = PageParsed{
		dbHeader:         DbHeader{},
		dbHeaderSize:     0,
		pageNumber:       7,
		cellAreaParsed:   [][]byte{{0, 0, 0, 3, 0, 12}},
		btreeType:        int(TableBtreeInteriorCell),
		rightMostpointer: []byte{0, 0, 0, 4},
		cellArea:         []byte{0, 0, 0, 3, 0, 12},

		startCellContentArea: PageSize - 6,
		isOverflow:           false,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 4, 0, 0, 0, 0}, []byte{4, 3, 0, 0, 0, 0}, []byte{4, 2, 0, 0, 0, 0}, []byte{4, 1, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 7, 0, 0, 0, 0}, []byte{4, 6, 0, 0, 0, 0}, []byte{4, 5, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (3 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 12, 0, 0, 0, 0}, []byte{4, 11, 0, 0, 0, 0}, []byte{4, 10, 0, 0, 0, 0}, []byte{4, 9, 0, 0, 0, 0}, []byte{4, 8, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (5 * 6),
		isOverflow:           true,

		isLeaf: true,
	}
	var fourthPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           4,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 16, 0, 0, 0, 0}, []byte{4, 15, 0, 0, 0, 0}, []byte{4, 14, 0, 0, 0, 0}, []byte{4, 13, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,

		isLeaf: true,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage
	memoryPages[fourthPage.pageNumber] = fourthPage
	memoryPages[5] = CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, 5, nil)
	memoryPages[sixthPage.pageNumber] = sixthPage
	memoryPages[seventhPage.pageNumber] = seventhPage

	btree := BtreeStruct{
		softPages: map[int]PageParsed{},
	}
	cell := createCell(TableBtreeLeafCell, 4, "Alice")
	node := btree.insert(4, cell, &server.header, nil)

	if node.pageNumber != 1 {
		t.Errorf("Insert values should be in page: %v, instead we got: %v", 4, node.pageNumber)
	}

	if len(node.cellAreaParsed) != 4 {
		t.Errorf("expected cell area to have 5 elements, instead we got: %v", len(node.cellAreaParsed))
	}

	if !reflect.DeepEqual(node.cellAreaParsed[0], cell.data) {
		t.Errorf("expected cell area start with newly added cell: %v, instead we got: %v", cell.data, node.cellAreaParsed[0])
	}

}

// //finish this one

// right most pointer points to right page with another values 0x05 or 0x0d
// cell  area contains all right most page 0, 0, 0, 2, 0, , but additonal has rowid, not like right most page
func TestInsertOneRecord(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:            header(),
		dbHeaderSize:        100,
		btreePageHeaderSize: 12,
		pageNumber:          0,
		cellAreaParsed:      [][]byte{{0, 0, 0, 1, 0, 1}, {0, 0, 0, 2, 0, 2}},
		btreeType:           int(TableBtreeInteriorCell),
		// remove right most pointer???
		rightMostpointer:     []byte{0, 0, 0, 2},
		cellArea:             []byte{0, 0, 0, 1, 0, 1, 0, 0, 0, 2, 0, 2},
		cellAreaSize:         12,
		numberofCells:        2,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,

		isLeaf: false,
	}

	server := ServerStruct{
		header: zeroPage.dbHeader,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		isOverflow:           true,

		isLeaf: true,
	}

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		isOverflow:           true,

		isLeaf: true,
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	btree := BtreeStruct{
		softPages: map[int]PageParsed{},
	}

	rowId := 3
	cell := createCell(TableBtreeLeafCell, rowId, "aliceAndBob129876theEnd12345")
	btree.insert(rowId, cell, &server.header, nil)

	// reader := NewReader("fd")

	secondPageSaved := btree.softPages[2]
	zeroPageSaved := btree.softPages[0]

	zeroPageExpectedCellArea := []byte{0, 0, 0, 1, 0, 1, 0, 0, 0, 2, 0, 3}

	if !reflect.DeepEqual(zeroPageSaved.cellArea, zeroPageExpectedCellArea) {
		t.Errorf("expected cell area on zero page to be: %v, insted we got: %v", zeroPageExpectedCellArea, zeroPageSaved.cellArea)
	}

	if secondPageSaved.numberofCells != 1 {
		t.Errorf("expected to have one cell, instead we got: %v", secondPageSaved.numberofCells)
	}
	fmt.Println("before fail?")

	if !reflect.DeepEqual(secondPageSaved.cellAreaParsed[0], secondPageSaved.cellArea) {
		t.Errorf("expected parsed celle area to be equal to cell area")
	}

	if !reflect.DeepEqual(secondPageSaved.cellArea, cell.data) {
		t.Errorf("expected cell area to be: %v, got: %v", cell.data, secondPageSaved.cellArea)
	}
}

func TestInsertMultipleRecord(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		btreePageHeaderSize:  12,
		pageNumber:           0,
		cellAreaParsed:       [][]byte{{0, 0, 0, 2, 0, 2}, {0, 0, 0, 1, 0, 1}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 2},
		cellArea:             []byte{0, 0, 0, 2, 0, 2, 0, 0, 0, 1, 0, 1},
		cellAreaSize:         12,
		numberofCells:        2,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,
		isLeaf:               false,
	}

	server := ServerStruct{
		header: zeroPage.dbHeader,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		isOverflow:           true,
		isLeaf:               true,
	}

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		isOverflow:           true,
		isLeaf:               true,
	}

	btree := BtreeStruct{
		softPages: map[int]PageParsed{},
	}
	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage

	cell1 := createCell(TableBtreeLeafCell, 3, "aliceAndBob129876theEnd12345")
	btree.insert(3, cell1, &server.header, nil)
	memoryPages = btree.softPages

	cell2 := createCell(TableBtreeLeafCell, 4, "aliceAndBob129876theEnd12345")
	btree.insert(4, cell2, &server.header, nil)

	secondPageSaved := btree.softPages[2]
	zeroPageSaved := btree.softPages[0]

	zeroPageExpectedCellArea := []byte{0, 0, 0, 2, 0, 4, 0, 0, 0, 1, 0, 1}

	if !reflect.DeepEqual(zeroPageSaved.cellArea, zeroPageExpectedCellArea) {
		t.Errorf("expected cell area on zero page to be: %v, insted we got: %v", zeroPageExpectedCellArea, zeroPageSaved.cellArea)
	}

	if secondPageSaved.numberofCells != 2 {
		t.Errorf("expected to have: %v cells, instead we got: %v", 2, secondPageSaved.numberofCells)
	}

	parsedData := dbReadparseCellArea(byte(TableBtreeLeafCell), secondPageSaved.cellArea)

	if !reflect.DeepEqual(parsedData, secondPageSaved.cellAreaParsed) {
		t.Errorf("expected parsed celle area to be equal to cell area: %v, instead we got: %v", parsedData, secondPageSaved.cellArea)
	}

	expectedCellArea := append([]byte{}, cell2.data...)
	expectedCellArea = append(expectedCellArea, cell1.data...)

	if !reflect.DeepEqual(secondPageSaved.cellArea, expectedCellArea) {
		t.Errorf("expected cell area to be: %v, got: %v", expectedCellArea, secondPageSaved.cellArea)
	}
}

// //start debugging this

func creareARowItem(length int, rowId int) CreateCell {

	value := ""
	for i := 0; i < length-5; i++ {
		value += string('a')
	}
	return createCell(TableBtreeLeafCell, rowId-1, value)
}

func TestInsertOverflowPage(t *testing.T) {
	clearDbFile("test")
	PageSize = 350
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		btreePageHeaderSize:  12,
		pageNumber:           0,
		cellAreaParsed:       [][]byte{{0, 0, 0, 2, 0, 1}, {0, 0, 0, 1, 0, 0}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 2},
		cellArea:             []byte{0, 0, 0, 2, 0, 1, 0, 0, 0, 1, 0, 0},
		cellAreaSize:         12,
		numberofCells:        2,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,

		isLeaf: false,
	}

	server := ServerStruct{
		header: zeroPage.dbHeader,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		isOverflow:           false,
		isLeaf:               true,
	}

	cell := creareARowItem(100, 2)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           2,
		numberofCells:        1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		startCellContentArea: PageSize - cell.dataLength,
		cellAreaSize:         cell.dataLength,
		cellArea:             cell.data,
		isOverflow:           false,
		isLeaf:               true,
	}

	btree := BtreeStruct{
		softPages: map[int]PageParsed{0: zeroPage, 1: firstPage, 2: secondPage},
	}

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage

	cell1 := creareARowItem(100, 3)
	btree.insert(2, cell1, &server.header, nil)

	memoryPages = btree.softPages

	cell2 := creareARowItem(100, 4)
	btree.insert(3, cell2, &server.header, nil)
	memoryPages = btree.softPages
	fmt.Println("afeer csexcond, lets see page number twoo")
	fmt.Printf("%+v", btree.softPages[2])
	fmt.Println("afeer csexcond, lets see page number twoo")

	cell3 := creareARowItem(100, 5)
	btree.insert(4, cell3, &server.header, nil)
	fmt.Println("afeer third, lets see page number twoo")
	fmt.Printf("%+v", btree.softPages[2])
	fmt.Println("afeer third, lets see page number twoo")

	secondPageSaved := btree.softPages[2]
	zeroPageSaved := btree.softPages[0]

	zeroPageExpectedCellArea := []byte{0, 0, 0, 2, 0, 4, 0, 0, 0, 1, 0, 3}

	if !reflect.DeepEqual(zeroPageSaved.cellArea, zeroPageExpectedCellArea) {
		t.Errorf("expected cell area on zero page to be: %v, insted we got: %v", zeroPageExpectedCellArea, zeroPageSaved.cellArea)
	}

	if secondPageSaved.numberofCells != 1 {
		t.Errorf("expected to have one cell, instead we got: %v", secondPageSaved.numberofCells)
	}

	parsedData := dbReadparseCellArea(byte(TableBtreeLeafCell), secondPageSaved.cellArea)

	if !reflect.DeepEqual(parsedData, secondPageSaved.cellAreaParsed) {
		t.Errorf("expected parsed celle area to be equal to cell area: %v, instead we got: %v", parsedData, secondPageSaved.cellArea)
	}

	expectedCellArea := append([]byte{}, cell3.data...)

	if !reflect.DeepEqual(secondPageSaved.cellArea, expectedCellArea) {
		t.Errorf("expected cell area to be: %v, got: %v", expectedCellArea, secondPageSaved.cellArea)
	}
}

func TestLeafBias(t *testing.T) {
	usableSpacePerPage = 20

	cells := []Cell{{size: 10, pageNumber: 1, rowId: 5, data: []byte{}}, {size: 10, pageNumber: 2, rowId: 7, data: []byte{}}, {size: 10, pageNumber: 3, rowId: 8, data: []byte{}}}
	totalSizeInEachPage, numberOfCellPerPage := leaf_bias(cells)

	if totalSizeInEachPage[0] != 20 {
		t.Errorf("Expected total size in first page to be :%v, got: %v", 20, totalSizeInEachPage[0])
	}

	if totalSizeInEachPage[1] != 10 {
		t.Errorf("Expected total size in first page to be :%v, got: %v", 10, totalSizeInEachPage[1])
	}

	if numberOfCellPerPage[0] != 2 {
		t.Errorf("Expected number of cell per page in second page to be: %v, got :%v", 2, numberOfCellPerPage[0])
	}

	if numberOfCellPerPage[1] != 1 {
		t.Errorf("Expected number of cell per page in second page to be: %v, got :%v", 1, numberOfCellPerPage[1])
	}
}

func TestInsertWithInteriorNested(t *testing.T) {
	clearDbFile("test")
	PageSize = 300
	server := ServerStruct{
		header: DbHeader{},
	}
	var zeroPage = CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 1}, {0, 0, 0, 1, 0, 0}}, server.header.assignNewPage(), &server.header)

	var firstPage = CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, server.header.assignNewPage(), nil)

	cell := creareARowItem(100, 2)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var secondPage = CreateNewPage(BtreeType(TableBtreeLeafCell), cellParsed, server.header.assignNewPage(), nil)

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage

	btree := BtreeStruct{
		softPages: map[int]PageParsed{0: zeroPage, 1: firstPage, 2: secondPage},
	}
	for i := 3; i < 74; i++ {
		cell1 := creareARowItem(100, i)
		btree.insert(i, cell1, &server.header, nil)
		memoryPages = btree.softPages

	}
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// PageSize: 300
	// firstpage: 300 - 100 (main header) = 200
	// 200 - 12 (interior pointer) = 188
	// every time has 6 byes + 2 pointer = 8
	// 188/8 = 23 items on the first page, before split

	// every leaf page contains two itesm, meaning 23 * 2, first page should fit from 0 til 46, then should split

	reader := NewReader("")
	zeroPageSaved := reader.readFromMemory(0)
	twentyFifthPageSaved := reader.readFromMemory(25)

	expectedCellAreaContentPageZero := []byte{0, 0, 0, 25, 0, 72}

	expectedCellAreaParsedContentPageZero := dbReadparseCellArea(byte(TableBtreeInteriorCell), expectedCellAreaContentPageZero)
	if !reflect.DeepEqual(zeroPageSaved.cellArea, expectedCellAreaContentPageZero) {
		t.Errorf("Expected page zero cell area to be: %v, instead we got: %v", expectedCellAreaContentPageZero, zeroPageSaved.cellArea)
	}

	if !reflect.DeepEqual(zeroPageSaved.cellAreaParsed, expectedCellAreaParsedContentPageZero) {
		t.Errorf("Expected page zero cell area parsed to be: %v, instead we got: %v", expectedCellAreaParsedContentPageZero, zeroPageSaved.cellAreaParsed)
	}

	if zeroPageSaved.numberofCells != 1 {
		t.Errorf("Expected page zero to have only 1 cell, instead we got: %v", zeroPageSaved.numberofCells)
	}

	expectedPointerPageZero := intToBinary(PageSize-len(expectedCellAreaContentPageZero), 2)

	if !reflect.DeepEqual(zeroPageSaved.pointers, expectedPointerPageZero) {
		t.Errorf("Expected page zero pointers to be: %v, instead we got: %v", zeroPageSaved.pointers, expectedPointerPageZero)
	}

	if zeroPageSaved.startCellContentArea != PageSize-len(expectedCellAreaContentPageZero) {
		t.Errorf("expected cell area to start at :%v, instead we got: %v", PageSize-len(expectedCellAreaContentPageZero), zeroPageSaved.startCellContentArea)
	}

	expectedLastCellAreaTwentyFifthPage := []byte{0, 0, 0, 37, 0, 72}
	expectedFirstCellAreaTwentyFifthPage := []byte{0, 0, 0, 1, 0, 2}

	if !reflect.DeepEqual(twentyFifthPageSaved.cellAreaParsed[0], expectedLastCellAreaTwentyFifthPage) {
		t.Errorf("expected last cell area of twnety fifth page to be: %v, instead we got: %v", expectedLastCellAreaTwentyFifthPage, twentyFifthPageSaved.cellAreaParsed[0])
	}
	if !reflect.DeepEqual(twentyFifthPageSaved.cellAreaParsed[len(twentyFifthPageSaved.cellAreaParsed)-1], expectedFirstCellAreaTwentyFifthPage) {
		t.Errorf("expected first cell area of twnety fifth page to be: %v, instead we got: %v", expectedFirstCellAreaTwentyFifthPage, twentyFifthPageSaved.cellAreaParsed[len(twentyFifthPageSaved.cellAreaParsed)-1])
	}
}

func TestInsertWithInteriorNestedSplitted(t *testing.T) {
	clearDbFile("test")
	PageSize = 300
	server := ServerStruct{
		header: DbHeader{},
	}
	var zeroPage = CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 1}, {0, 0, 0, 1, 0, 0}}, server.header.assignNewPage(), &server.header)

	var firstPage = CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, server.header.assignNewPage(), nil)

	cell := creareARowItem(100, 2)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var secondPage = CreateNewPage(BtreeType(TableBtreeLeafCell), cellParsed, server.header.assignNewPage(), nil)

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage

	btree := BtreeStruct{
		softPages: map[int]PageParsed{0: zeroPage, 1: firstPage, 2: secondPage},
	}

	for i := 3; i < 75; i++ {
		cell1 := creareARowItem(100, i)

		btree.insert(i, cell1, &server.header, nil)
		memoryPages = btree.softPages

	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	reader := NewReader("")
	zeroPageSaved := reader.readFromMemory(0)
	twentyFifthPageSaved := reader.readFromMemory(25)
	thirtEightPageSaved := reader.readFromMemory(39)
	fmt.Printf("%+v", thirtEightPageSaved)

	expectedCellAreaContentPageZero := []byte{0, 0, 0, 39, 0, 73, 0, 0, 0, 25, 0, 38}
	expectedRightMostPointerPageZero := []byte{0, 0, 0, 39}

	expectedCellAreaParsedContentPageZero := dbReadparseCellArea(byte(TableBtreeInteriorCell), expectedCellAreaContentPageZero)
	if !reflect.DeepEqual(zeroPageSaved.cellArea, expectedCellAreaContentPageZero) {
		t.Errorf("Expected page zero cell area to be: %v, instead we got: %v", expectedCellAreaContentPageZero, zeroPageSaved.cellArea)
	}

	if !reflect.DeepEqual(zeroPageSaved.cellAreaParsed, expectedCellAreaParsedContentPageZero) {
		t.Errorf("Expected page zero cell area parsed to be: %v, instead we got: %v", expectedCellAreaParsedContentPageZero, zeroPageSaved.cellAreaParsed)
	}

	if zeroPageSaved.numberofCells != 2 {
		t.Errorf("Expected page zero to have only 1 cell, instead we got: %v", zeroPageSaved.numberofCells)
	}

	if zeroPageSaved.btreePageHeaderSize != 12 {
		t.Errorf("Expected header size to be:%v got: %v", 12, zeroPageSaved.btreePageHeaderSize)
	}

	expectedPointerPageZero := append([]byte{}, intToBinary(PageSize-6, 2)...)
	expectedPointerPageZero = append(expectedPointerPageZero, intToBinary(PageSize-(6*2), 2)...)

	if !reflect.DeepEqual(zeroPageSaved.pointers, expectedPointerPageZero) {
		t.Errorf("Expected page zero pointers to be: %v, instead we got: %v", zeroPageSaved.pointers, expectedPointerPageZero)
	}

	if zeroPageSaved.startCellContentArea != PageSize-len(expectedCellAreaContentPageZero) {
		t.Errorf("expected cell area to start at :%v, instead we got: %v", PageSize-len(expectedCellAreaContentPageZero), zeroPageSaved.startCellContentArea)
	}

	if !reflect.DeepEqual(expectedRightMostPointerPageZero, zeroPageSaved.rightMostpointer) {
		t.Errorf("expected zero page right most pointer to be: %v, instead we got: %v", expectedRightMostPointerPageZero, zeroPageSaved.rightMostpointer)
	}

	expectedLastCellAreaTwentyFifthPage := []byte{0, 0, 0, 19, 0, 38}
	expectedFirstCellAreaTwentyFifthPage := []byte{0, 0, 0, 1, 0, 2}
	expectedRightMostPointerTwentyFifthPage := []byte{0, 0, 0, 19}

	if !reflect.DeepEqual(twentyFifthPageSaved.cellAreaParsed[0], expectedLastCellAreaTwentyFifthPage) {
		t.Errorf("expected last cell area of twnety fifth page to be: %v, instead we got: %v", expectedLastCellAreaTwentyFifthPage, twentyFifthPageSaved.cellAreaParsed[0])
	}
	if !reflect.DeepEqual(twentyFifthPageSaved.cellAreaParsed[len(twentyFifthPageSaved.cellAreaParsed)-1], expectedFirstCellAreaTwentyFifthPage) {
		t.Errorf("expected first cell area of twnety fifth page to be: %v, instead we got: %v", expectedFirstCellAreaTwentyFifthPage, twentyFifthPageSaved.cellAreaParsed[len(twentyFifthPageSaved.cellAreaParsed)-1])
	}
	if !reflect.DeepEqual(expectedRightMostPointerTwentyFifthPage, twentyFifthPageSaved.rightMostpointer) {
		t.Errorf("expected twenty fifth page right most pointer to be: %v, instead we got: %v", expectedRightMostPointerTwentyFifthPage, twentyFifthPageSaved.rightMostpointer)
	}
	if twentyFifthPageSaved.btreeType != int(TableBtreeInteriorCell) {
		t.Errorf("expected twenty fifth page btree type to be: %v, instead we got: %v", twentyFifthPageSaved.btreeType, TableBtreeInteriorCell)
	}

	expectedLastCellAreaTirthyEightPage := []byte{0, 0, 0, 38, 0, 73}
	expectedFirstCellAreaTirthyEightPage := []byte{0, 0, 0, 20, 0, 40}

	if !reflect.DeepEqual(thirtEightPageSaved.cellAreaParsed[0], expectedLastCellAreaTirthyEightPage) {
		t.Errorf("expected last cell area of twnety fifth page to be: %v, instead we got: %v", expectedLastCellAreaTirthyEightPage, thirtEightPageSaved.cellAreaParsed[0])
	}
	if !reflect.DeepEqual(thirtEightPageSaved.cellAreaParsed[len(thirtEightPageSaved.cellAreaParsed)-1], expectedFirstCellAreaTirthyEightPage) {
		t.Errorf("expected first cell area of twnety fifth page to be: %v, instead we got: %v", expectedFirstCellAreaTirthyEightPage, thirtEightPageSaved.cellAreaParsed[len(thirtEightPageSaved.cellAreaParsed)-1])
	}
	if thirtEightPageSaved.btreeType != int(TableBtreeInteriorCell) {
		t.Errorf("expected thirty eight page btree type to be: %v, instead we got: %v", thirtEightPageSaved.btreeType, TableBtreeInteriorCell)
	}

	expectedRightMostPointerThirthyEightPage := []byte{0, 0, 0, 38}
	//fix this pointer
	if !reflect.DeepEqual(thirtEightPageSaved.rightMostpointer, expectedRightMostPointerThirthyEightPage) {
		t.Errorf("Expected thirty eight page right most pointer to be: %v, got: %v", expectedRightMostPointerThirthyEightPage, thirtEightPageSaved.rightMostpointer)
	}
}

func TestInsertWithInteriorNestedSplittedSix(t *testing.T) {
	clearDbFile("test")
	PageSize = 134
	server := ServerStruct{
		header: DbHeader{},
	}
	var zeroPage = CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 1, 0, 1}}, server.header.assignNewPage(), &server.header)

	cell := creareARowItem(100, 2)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)
	var firstPage = CreateNewPage(BtreeType(TableBtreeLeafCell), cellParsed, server.header.assignNewPage(), nil)

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage

	btree := BtreeStruct{
		softPages: map[int]PageParsed{0: zeroPage, 1: firstPage},
	}

	for i := 3; i < 120; i++ {
		cell1 := creareARowItem(100, i)
		btree.insert(i, cell1, &server.header, nil)
		memoryPages = btree.softPages

	}

	reader := NewReader("")
	zeroPageSaved := reader.readFromMemory(0)
	thirtyFifthSaved := reader.readFromMemory(35)

	expectedZeroPageCellArea := []byte{0, 0, 0, 35, 0, 118}
	expectedTirthyFifthCellArea := []byte{0, 0, 0, 115, 0, 118, 0, 0, 0, 99, 0, 105, 0, 0, 0, 83, 0, 90, 0, 0, 0, 67, 0, 75, 0, 0, 0, 51, 0, 60, 0, 0, 0, 34, 0, 45, 0, 0, 0, 18, 0, 30, 0, 0, 0, 4, 0, 15}

	if !reflect.DeepEqual(thirtyFifthSaved.cellArea, expectedTirthyFifthCellArea) {
		t.Errorf("thirty fitfth area should be: %v, instead we got: %v", expectedTirthyFifthCellArea, thirtyFifthSaved.cellArea)
	}

	if !reflect.DeepEqual(zeroPageSaved.cellArea, expectedZeroPageCellArea) {
		t.Errorf("Zero page area should be: %v, instead we got: %v", expectedZeroPageCellArea, zeroPageSaved.cellArea)
	}

}

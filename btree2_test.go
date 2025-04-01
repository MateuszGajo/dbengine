package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func clearDbFile(fileName string) {
	softWritePages = make(map[int]PageParsed)
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	server := ServerStruct{
		firstPage: zeroPage,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 4, 0, 0, 0, 0}, []byte{4, 3, 0, 0, 0, 0}, []byte{4, 2, 0, 0, 0, 0}, []byte{4, 1, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 7, 0, 0, 0, 0}, []byte{4, 6, 0, 0, 0, 0}, []byte{4, 5, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (3 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 12, 0, 0, 0, 0}, []byte{4, 11, 0, 0, 0, 0}, []byte{4, 10, 0, 0, 0, 0}, []byte{4, 9, 0, 0, 0, 0}, []byte{4, 8, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (5 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(thirdPage), 3, "", &server.firstPage)
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	server := ServerStruct{
		firstPage: zeroPage,
	}
	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 4, 0, 0, 0, 0}, []byte{4, 3, 0, 0, 0, 0}, []byte{4, 2, 0, 0, 0, 0}, []byte{4, 1, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 7, 0, 0, 0, 0}, []byte{4, 6, 0, 0, 0, 0}, []byte{4, 5, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (3 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 12, 0, 0, 0, 0}, []byte{4, 11, 0, 0, 0, 0}, []byte{4, 10, 0, 0, 0, 0}, []byte{4, 9, 0, 0, 0, 0}, []byte{4, 8, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (5 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(thirdPage), 3, "", &server.firstPage)
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	server := ServerStruct{
		firstPage: zeroPage,
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
		leftSibling:          nil,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 7, 0, 0, 0, 0}, []byte{4, 6, 0, 0, 0, 0}, []byte{4, 5, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (3 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 12, 0, 0, 0, 0}, []byte{4, 11, 0, 0, 0, 0}, []byte{4, 10, 0, 0, 0, 0}, []byte{4, 9, 0, 0, 0, 0}, []byte{4, 8, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (5 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var fourthPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           4,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 16, 0, 0, 0, 0}, []byte{4, 15, 0, 0, 0, 0}, []byte{4, 14, 0, 0, 0, 0}, []byte{4, 13, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(thirdPage), 3, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(fourthPage), 4, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(PageParsed{}), 5, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(sixthPage), 6, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(seventhPage), 7, "", &server.firstPage)
	found, index, pageFound, parents := search(0, 7, []*PageParsed{})

	fmt.Println("lets see what we found")
	fmt.Println(found)
	fmt.Println(index)

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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	cells := []Cell{{rowId: 9, pageNumber: 2}, {rowId: 13, pageNumber: 3}}

	modifyDivider(&zeroPage, cells, 6, 6*3, &zeroPage, []*PageParsed{})

	reader := NewReader("")
	zeroPageMemory := reader.readFromMemory(0)

	cellAreaParsedExpected := cellAreaParsed
	cellAreaParsedExpected[1] = []byte{0, 0, 0, 3, 0, 13}
	cellAreaParsedExpected[2] = []byte{0, 0, 0, 2, 0, 9}

	if !reflect.DeepEqual(cellAreaParsedExpected, zeroPageMemory.cellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", cellAreaParsedExpected, zeroPageMemory.cellAreaParsed)
	}

	if zeroPageMemory.cellAreaSize != len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("cell area size should be: %v, got: %v", len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPageMemory.cellAreaSize)
	}

	if zeroPageMemory.numberofCells != len(cellAreaParsedExpected) {
		t.Errorf("number of cell should be: %v, got: %v", len(cellAreaParsedExpected), zeroPageMemory.numberofCells)
	}

	if zeroPageMemory.startCellContentArea != PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("start of cell content should be: %v, got: %v", PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPageMemory.startCellContentArea)
	}

	expectedRightMostPointer := []byte{0, 0, 0, 4}

	if !reflect.DeepEqual(zeroPageMemory.rightMostpointer, expectedRightMostPointer) {
		t.Errorf("expected right most pointer to be: %v, instead we got: %v", expectedRightMostPointer, zeroPageMemory.rightMostpointer)
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	cells := []Cell{{rowId: 9, pageNumber: 2}}

	modifyDivider(&zeroPage, cells, 6, 6*3, &zeroPage, []*PageParsed{})

	reader := NewReader("")
	zeroPageMemory := reader.readFromMemory(0)

	cellAreaParsedExpected := [][]byte{cellAreaParsed[0]}
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 2, 0, 9})
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[3])

	if !reflect.DeepEqual(cellAreaParsedExpected, zeroPageMemory.cellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", cellAreaParsedExpected, zeroPageMemory.cellAreaParsed)
	}

	if zeroPageMemory.cellAreaSize != len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("cell area size should be: %v, got: %v", len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPageMemory.cellAreaSize)
	}

	if zeroPageMemory.numberofCells != len(cellAreaParsedExpected) {
		t.Errorf("number of cell should be: %v, got: %v", len(cellAreaParsedExpected), zeroPageMemory.numberofCells)
	}

	if zeroPageMemory.startCellContentArea != PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("start of cell content should be: %v, got: %v", PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPageMemory.startCellContentArea)
	}

	expectedRightMostPointer := []byte{0, 0, 0, 4}

	if !reflect.DeepEqual(zeroPageMemory.rightMostpointer, expectedRightMostPointer) {
		t.Errorf("expected right most pointer to be: %v, instead we got: %v", expectedRightMostPointer, zeroPageMemory.rightMostpointer)
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	fmt.Println("zero area cell page?")
	fmt.Println(zeroPage.cellArea)

	cells := []Cell{{rowId: 9, pageNumber: 2}, {rowId: 13, pageNumber: 3}, {rowId: 14, pageNumber: 5}}

	modifyDivider(&zeroPage, cells, 6, 6*3, &zeroPage, []*PageParsed{})

	reader := NewReader("")
	zeroPageMemory := reader.readFromMemory(0)

	cellAreaParsedExpected := [][]byte{cellAreaParsed[0]}
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 5, 0, 14})
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 3, 0, 13})
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 2, 0, 9})
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[3])

	if !reflect.DeepEqual(cellAreaParsedExpected, zeroPageMemory.cellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", cellAreaParsedExpected, zeroPageMemory.cellAreaParsed)
	}

	if zeroPageMemory.cellAreaSize != len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("cell area size should be: %v, got: %v", len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPageMemory.cellAreaSize)
	}

	if zeroPageMemory.numberofCells != len(cellAreaParsedExpected) {
		t.Errorf("number of cell should be: %v, got: %v", len(cellAreaParsedExpected), zeroPageMemory.numberofCells)
	}

	if zeroPageMemory.startCellContentArea != PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("start of cell content should be: %v, got: %v", PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPageMemory.startCellContentArea)
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	fmt.Println("zero area cell page?")
	fmt.Println(zeroPage.cellArea)

	cells := []Cell{{rowId: 18, pageNumber: 5}}

	modifyDivider(&zeroPage, cells, 0, 6, &zeroPage, []*PageParsed{})

	reader := NewReader("")
	zeroPageMemory := reader.readFromMemory(0)

	cellAreaParsedExpected := [][]byte{[]byte{0, 0, 0, 5, 0, 18}}
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[1])
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[2])
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[3])

	if !reflect.DeepEqual(cellAreaParsedExpected, zeroPageMemory.cellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", cellAreaParsedExpected, zeroPageMemory.cellAreaParsed)
	}

	if zeroPageMemory.cellAreaSize != len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("cell area size should be: %v, got: %v", len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPageMemory.cellAreaSize)
	}

	if zeroPageMemory.numberofCells != len(cellAreaParsedExpected) {
		t.Errorf("number of cell should be: %v, got: %v", len(cellAreaParsedExpected), zeroPageMemory.numberofCells)
	}

	if zeroPageMemory.startCellContentArea != PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]) {
		t.Errorf("start of cell content should be: %v, got: %v", PageSize-len(cellAreaParsedExpected)*len(cellAreaParsedExpected[0]), zeroPageMemory.startCellContentArea)
	}

	expectedRightMostPointer := []byte{0, 0, 0, 5}

	if !reflect.DeepEqual(zeroPageMemory.rightMostpointer, expectedRightMostPointer) {
		t.Errorf("expected right most pointer to be: %v, instead we got: %v", expectedRightMostPointer, zeroPageMemory.rightMostpointer)
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
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
		leftSibling:          nil,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	fmt.Println("zero area cell page?")
	fmt.Println(zeroPage.cellArea)

	cells := []Cell{{rowId: 18, pageNumber: 5}}

	modifyDivider(&firstPage, cells, 0, 6, &zeroPage, []*PageParsed{&zeroPage})

	reader := NewReader("")
	zeroPageMemory := reader.readFromMemory(0)

	expectedCellAreaZeroPage := []byte{0, 0, 0, 1, 0, 18}
	expectedCellAreaParsedZeroPage := dbReadparseCellArea(byte(TableBtreeInteriorCell), expectedCellAreaZeroPage)

	if !reflect.DeepEqual(expectedCellAreaZeroPage, zeroPageMemory.cellArea) {
		t.Errorf("expected zero page to be cell area: %v, instead we got: %v", expectedCellAreaZeroPage, zeroPageMemory.cellArea)
	}

	if !reflect.DeepEqual(expectedCellAreaParsedZeroPage, zeroPageMemory.cellAreaParsed) {
		t.Errorf("expected zero page to be cell area: %v, instead we got: %v", expectedCellAreaParsedZeroPage, zeroPageMemory.cellAreaParsed)
	}

	if zeroPageMemory.numberofCells != 1 {
		t.Errorf("expected zero page to have 1 cell, instead we got: %v", zeroPageMemory.numberofCells)
	}

	expectedPointers := append([]byte{}, intToBinary(PageSize-len(zeroPagecellArea), 2)...)

	if !reflect.DeepEqual(expectedPointers, zeroPage.pointers) {
		t.Errorf("expected zero page pointers to be: %v, instead we got: %v", expectedPointers, zeroPage.pointers)
	}

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
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
		leftSibling:          nil,
		rightSiblisng:        nil,
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
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
		leftSibling:          nil,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	fmt.Println("zero area cell page?")
	fmt.Println(zeroPage.cellArea)

	cells := []Cell{{rowId: 18, pageNumber: 5}}

	writer := NewWriter()

	writer.softwiteToFile(&zeroPage, 0, &zeroPage)
	writer.softwiteToFile(&firstPage, 1, &zeroPage)
	writer.softwiteToFile(&secondPage, 2, &zeroPage)

	modifyDivider(&secondPage, cells, 0, 6, &zeroPage, []*PageParsed{&zeroPage, &firstPage})

	reader := NewReader("")
	zeroPageMemory := reader.readFromMemory(0)
	firstPageMemory := reader.readFromMemory(1)

	expectedCellAreaZeroPage := []byte{0, 0, 0, 1, 0, 18}
	expectedCellAreaParsedZeroPage := dbReadparseCellArea(byte(TableBtreeInteriorCell), expectedCellAreaZeroPage)

	if !reflect.DeepEqual(expectedCellAreaZeroPage, zeroPageMemory.cellArea) {
		t.Errorf("expected zero page to be cell area: %v, instead we got: %v", expectedCellAreaZeroPage, zeroPageMemory.cellArea)
	}

	if !reflect.DeepEqual(expectedCellAreaParsedZeroPage, zeroPageMemory.cellAreaParsed) {
		t.Errorf("expected zero page to be cell area: %v, instead we got: %v", expectedCellAreaParsedZeroPage, zeroPageMemory.cellAreaParsed)
	}

	if zeroPageMemory.numberofCells != 1 {
		t.Errorf("expected zero page to have 1 cell, instead we got: %v", zeroPageMemory.numberofCells)
	}

	expectedPointersZeroPage := append([]byte{}, intToBinary(PageSize-len(zeroPagecellArea), 2)...)

	if !reflect.DeepEqual(expectedPointersZeroPage, zeroPage.pointers) {
		t.Errorf("expected first page pointers to be: %v, instead we got: %v", expectedPointersZeroPage, zeroPage.pointers)
	}

	expectedCellAreaFirstPage := []byte{0, 0, 0, 2, 0, 18}
	expectedCellAreaParsedFirstPage := dbReadparseCellArea(byte(TableBtreeInteriorCell), expectedCellAreaFirstPage)

	if !reflect.DeepEqual(expectedCellAreaFirstPage, firstPageMemory.cellArea) {
		t.Errorf("expected first page to be cell area: %v, instead we got: %v", expectedCellAreaFirstPage, firstPageMemory.cellArea)
	}

	if !reflect.DeepEqual(expectedCellAreaParsedFirstPage, firstPageMemory.cellAreaParsed) {
		t.Errorf("expected first page to be cell area: %v, instead we got: %v", expectedCellAreaParsedFirstPage, zeroPageMemory.cellAreaParsed)
	}

	if firstPageMemory.numberofCells != 1 {
		t.Errorf("expected first page to have 1 cell, instead we got: %v", firstPageMemory.numberofCells)
	}

	expectedPointersFirstPage := append([]byte{}, intToBinary(PageSize-len(zeroPagecellArea), 2)...)

	if !reflect.DeepEqual(expectedPointersFirstPage, zeroPage.pointers) {
		t.Errorf("expected first page pointers to be: %v, instead we got: %v", expectedPointersFirstPage, zeroPage.pointers)
	}

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
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

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
		cellAreaSize:         cell.dataLength,
		numberofCells:        1,
		isOverflow:           false,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.softwiteToFile(&zeroPage, 0, &zeroPage)
	writer.softwiteToFile(&firstPage, 1, &zeroPage)
	writer.softwiteToFile(&secondPage, 2, &zeroPage)

	leftSibling, rightSibling := zeroPage.findSiblings(firstPage)

	if leftSibling != nil {
		t.Errorf("Expected left sibling to be nil, got: %v", leftSibling)
	}
	fmt.Println(rightSibling)

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
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

	cell := creareARowItem(100, 2)
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.softwiteToFile(&zeroPage, 0, &zeroPage)
	writer.softwiteToFile(&firstPage, 1, &zeroPage)
	writer.softwiteToFile(&secondPage, 2, &zeroPage)

	leftSibling, rightSibling := zeroPage.findSiblings(secondPage)

	if rightSibling != nil {
		t.Errorf("Expected left sibling to be nil, got: %v", rightSibling)
	}
	fmt.Println(rightSibling)

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
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

	cell := creareARowItem(100, 2)
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.softwiteToFile(&zeroPage, 0, &zeroPage)
	writer.softwiteToFile(&firstPage, 1, &zeroPage)
	writer.softwiteToFile(&secondPage, 2, &zeroPage)
	writer.softwiteToFile(&thirdPage, 2, &zeroPage)

	leftSibling, rightSibling := zeroPage.findSiblings(secondPage)

	if rightSibling == nil || rightSibling.pageNumber != 3 {
		t.Errorf("Expected right siblings to be page number 3")
	}
	fmt.Println(rightSibling)

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
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

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
		cellAreaSize:         cell.dataLength,
		numberofCells:        1,
		isOverflow:           false,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.softwiteToFile(&zeroPage, 0, &zeroPage)
	writer.softwiteToFile(&firstPage, 1, &zeroPage)

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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	newCell := createCell(TableBtreeLeafCell, 0, "Alice")
	firstPage.updateData(newCell, 0)

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

	cell := creareARowItem(100, 0)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	fmt.Println("current cell")
	fmt.Println(cell.rowId)

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
		leftSibling:          nil,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	newCell := createCell(TableBtreeLeafCell, 0, "Alice")

	firstPage.insertData(newCell, &zeroPage, []*PageParsed{&zeroPage})

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

	expectedPointer := append([]byte{}, intToBinary(PageSize-cell.dataLength, 2)...)
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}
	// 9*2 content, 2*3 pointer, 8 for headers

	PageSize = 9*2 + 2*3 + 8

	writer := NewWriter()

	balancingForNode(firstPage, []*PageParsed{}, &zeroPage)
	writer.flushPages("", &zeroPage)

	reader := NewReader("")

	zeroPageSaved := reader.readFromMemory(1)
	firstPageSaved := reader.readFromMemory(2)
	secondPageSaved := reader.readFromMemory(3)

	zeroPageExpectedCellArea := []byte{0, 0, 0, 3, 0, 4, 0, 0, 0, 2, 0, 2}
	zeroPageExpectedCellAreaParsed := dbReadparseCellArea(byte(TableBtreeInteriorCell), zeroPageExpectedCellArea)

	if len(zeroPageSaved.cellAreaParsed) != 2 {
		t.Errorf("expected cell area to have only 2 element, insted we got: %v", len(zeroPageSaved.cellAreaParsed))
	}
	if !reflect.DeepEqual(zeroPageSaved.cellArea, zeroPageExpectedCellArea) {
		t.Errorf("Expected zero page cell area parsed to be: %v, got: %v", zeroPageExpectedCellArea, zeroPageSaved.cellArea)
	}

	if !reflect.DeepEqual(zeroPageSaved.cellAreaParsed, zeroPageExpectedCellAreaParsed) {
		t.Errorf("expected zero page cell area to be: %v, got: %v", zeroPageSaved.cellAreaParsed, zeroPageExpectedCellAreaParsed)
	}

	if zeroPageSaved.numberofCells != 2 {
		t.Errorf("expected root page to have only 2 element, we got: %v", zeroPageSaved.numberofCells)
	}

	if zeroPageSaved.startCellContentArea != PageSize-(2*6) {
		t.Errorf("expected cell area start to be: %v, got: %v", PageSize-(2*6), zeroPageSaved.startCellContentArea)
	}

	firstPageExpectedCellAreaParsed := append([][]byte{}, cellAreaParsed[2])
	firstPageExpectedCellAreaParsed = append(firstPageExpectedCellAreaParsed, cellAreaParsed[3])

	firstPageExpectedCellArea := append([]byte{}, cellAreaParsed[2]...)
	firstPageExpectedCellArea = append(firstPageExpectedCellArea, cellAreaParsed[3]...)

	if !reflect.DeepEqual(firstPageSaved.cellAreaParsed, firstPageExpectedCellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", firstPageExpectedCellAreaParsed, firstPageSaved.cellAreaParsed)
	}

	if firstPageSaved.numberofCells != 2 {
		t.Errorf("expected first page to have 2 cells, instead we have: %v", firstPageSaved.numberofCells)
	}

	if !reflect.DeepEqual(firstPageSaved.cellArea, firstPageExpectedCellArea) {
		t.Errorf("Expected first page cell area to be: %v, instead we got: %v", firstPageSaved.cellArea, firstPageExpectedCellArea)
	}

	secondPageExpectedCellAreaParsed := append([][]byte{}, cellAreaParsed[0])
	secondPageExpectedCellAreaParsed = append(secondPageExpectedCellAreaParsed, cellAreaParsed[1])

	secondPageExpectedCellArea := append([]byte{}, cellAreaParsed[0]...)
	secondPageExpectedCellArea = append(secondPageExpectedCellArea, cellAreaParsed[1]...)

	if !reflect.DeepEqual(secondPageSaved.cellAreaParsed, secondPageExpectedCellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", secondPageExpectedCellAreaParsed, secondPageSaved.cellAreaParsed)
	}

	if secondPageSaved.numberofCells != 2 {
		t.Errorf("expected first page to have 2 cells, instead we have: %v", secondPageSaved.numberofCells)
	}

	if !reflect.DeepEqual(secondPageSaved.cellArea, secondPageExpectedCellArea) {
		t.Errorf("Expected first page cell area to be: %v, instead we got: %v", secondPageSaved.cellArea, secondPageExpectedCellArea)
	}

}

func TestBalancingSplitRootIntOneChild(t *testing.T) {
	clearDbFile("test")
	cellAreaParsed := [][]byte{{7, 4, 2, 23, 65, 108, 105, 99, 101}, {7, 3, 2, 23, 65, 108, 105, 99, 101}, {7, 2, 2, 23, 65, 108, 105, 99, 101}, {7, 1, 2, 23, 65, 108, 105, 99, 101}}
	cellArea := []byte{7, 4, 2, 23, 65, 108, 105, 99, 101, 7, 3, 2, 23, 65, 108, 105, 99, 101, 7, 2, 2, 23, 65, 108, 105, 99, 101, 7, 1, 2, 23, 65, 108, 105, 99, 101}
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	usableSpacePerPage = 9 * 4
	server := ServerStruct{
		firstPage: zeroPage,
	}

	writer := NewWriter()
	reader := NewReader("")

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	balancingForNode(zeroPage, []*PageParsed{}, &server.firstPage)
	writer.flushPages("", &server.firstPage)

	zeroPageSaved := reader.readFromMemory(0)
	firstPageSaved := reader.readFromMemory(1)

	fmt.Println("number of pages??")
	fmt.Println(server.firstPage.dbHeader.dbSizeInPages)

	if server.firstPage.dbHeader.dbSizeInPages != 2 {
		t.Errorf("Expected number of pages to be 2, got: %v", server.firstPage.dbHeader.dbSizeInPages)
	}

	zeroPageExpectedCellArea := []byte{0, 0, 0, 1, 0, 4}
	zeroPageExpectedCellAreaParsed := dbReadparseCellArea(byte(TableBtreeInteriorCell), zeroPageExpectedCellArea)

	if len(zeroPageSaved.cellAreaParsed) != 1 {
		t.Errorf("expected cell area to have only 2 element, insted we got: %v", len(zeroPageSaved.cellAreaParsed))
	}
	if !reflect.DeepEqual(zeroPageSaved.cellArea, zeroPageExpectedCellArea) {
		t.Errorf("Expected zero page cell area parsed to be: %v, got: %v", zeroPageExpectedCellArea, zeroPageSaved.cellArea)
	}

	if !reflect.DeepEqual(zeroPageSaved.cellAreaParsed, zeroPageExpectedCellAreaParsed) {
		t.Errorf("expected zero page cell area to be: %v, got: %v", zeroPageSaved.cellAreaParsed, zeroPageExpectedCellAreaParsed)
	}

	if zeroPageSaved.numberofCells != 1 {
		t.Errorf("expected root page to have only 1 element, we got: %v", zeroPageSaved.numberofCells)
	}

	if zeroPageSaved.startCellContentArea != PageSize-6 {
		t.Errorf("expected cell area start to be: %v, got: %v", PageSize-6, zeroPageSaved.startCellContentArea)
	}

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

func TestBalancingSplitOneLeaftPageIntoTwo(t *testing.T) {
	clearDbFile("test")
	cellAreaParsed := [][]byte{{7, 4, 2, 23, 65, 108, 105, 99, 101}, {7, 3, 2, 23, 65, 108, 105, 99, 101}, {7, 2, 2, 23, 65, 108, 105, 99, 101}, {7, 1, 2, 23, 65, 108, 105, 99, 101}}
	cellArea := []byte{7, 4, 2, 23, 65, 108, 105, 99, 101, 7, 3, 2, 23, 65, 108, 105, 99, 101, 7, 2, 2, 23, 65, 108, 105, 99, 101, 7, 1, 2, 23, 65, 108, 105, 99, 101}
	var zeroPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		cellAreaParsed:       [][]byte{{0, 0, 0, 2, 0, 4}},
		btreeType:            int(TableBtreeInteriorCell),
		btreePageHeaderSize:  12,
		rightMostpointer:     []byte{0, 0, 0, 2},
		cellArea:             []byte{0, 0, 0, 2, 0, 4},
		isLeaf:               false,
		startCellContentArea: PageSize - 6,
		cellAreaSize:         6,
		numberofCells:        1,
		isOverflow:           false,
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		cellAreaParsed:       cellAreaParsed,
		btreeType:            int(TableBtreeLeafCell),
		btreePageHeaderSize:  8,
		rightMostpointer:     []byte{},
		cellArea:             cellArea,
		cellAreaSize:         9 * 4,
		isLeaf:               true,
		startCellContentArea: PageSize - 9*4,
		numberofCells:        len(cellAreaParsed),
		isOverflow:           true,
		leftSibling:          nil,
		rightSiblisng:        nil,
	}
	PageSize = 9*2 + 2*3 + 8
	// usableSpacePerPage = 9 * 3
	server := ServerStruct{
		firstPage: PageParsed{
			dbHeader: DbHeader{
				dbSizeInPages: 3,
			},
		},
	}

	writer := NewWriter()
	reader := NewReader("")

	// writer.softwiteToFile(&zeroPage, 0, &server.firstPage)
	// writer.softwiteToFile(&firstPage, 1, &server.firstPage)

	balancingForNode(firstPage, []*PageParsed{&zeroPage}, &server.firstPage)
	writer.flushPages("", &server.firstPage)

	zeroPageSaved := reader.readFromMemory(1)
	firstPageSaved := reader.readFromMemory(2)
	secondPageSaved := reader.readFromMemory(3)

	fmt.Println("number of pages??")
	fmt.Println(server.firstPage.dbHeader.dbSizeInPages)

	if server.firstPage.dbHeader.dbSizeInPages != 4 {
		t.Errorf("Expected number of pages to be 4, got: %v", server.firstPage.dbHeader.dbSizeInPages)
	}

	zeroPageExpectedCellArea := []byte{0, 0, 0, 3, 0, 4, 0, 0, 0, 2, 0, 2}
	zeroPageExpectedCellAreaParsed := dbReadparseCellArea(byte(TableBtreeInteriorCell), zeroPageExpectedCellArea)

	fmt.Println("zer opage cell area parsed???")
	fmt.Println(zeroPageExpectedCellAreaParsed)

	if len(zeroPageSaved.cellAreaParsed) != 2 {
		t.Errorf("expected cell area to have only 2 element, insted we got: %v", len(zeroPageSaved.cellAreaParsed))
	}
	if !reflect.DeepEqual(zeroPageSaved.cellArea, zeroPageExpectedCellArea) {
		t.Errorf("Expected zero page cell area parsed to be: %v, got: %v", zeroPageExpectedCellArea, zeroPageSaved.cellArea)
	}

	if !reflect.DeepEqual(zeroPageSaved.cellAreaParsed, zeroPageExpectedCellAreaParsed) {
		t.Errorf("expected zero page cell area to be: %v, got: %v", zeroPageExpectedCellAreaParsed, zeroPageSaved.cellAreaParsed)
	}

	if zeroPageSaved.numberofCells != 2 {
		t.Errorf("expected root page to have only 2 element, we got: %v", zeroPageSaved.numberofCells)
	}

	if zeroPageSaved.startCellContentArea != PageSize-(2*6) {
		t.Errorf("expected cell area start to be: %v, got: %v", PageSize-(2*6), zeroPageSaved.startCellContentArea)
	}

	firstPageExpectedCellAreaParsed := append([][]byte{}, cellAreaParsed[2])
	firstPageExpectedCellAreaParsed = append(firstPageExpectedCellAreaParsed, cellAreaParsed[3])

	firstPageExpectedCellArea := append([]byte{}, cellAreaParsed[2]...)
	firstPageExpectedCellArea = append(firstPageExpectedCellArea, cellAreaParsed[3]...)

	if !reflect.DeepEqual(firstPageSaved.cellAreaParsed, firstPageExpectedCellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", firstPageExpectedCellAreaParsed, firstPageSaved.cellAreaParsed)
	}

	if firstPageSaved.numberofCells != 2 {
		t.Errorf("expected first page to have 2 cells, instead we have: %v", firstPageSaved.numberofCells)
	}

	if !reflect.DeepEqual(firstPageSaved.cellArea, firstPageExpectedCellArea) {
		t.Errorf("Expected first page cell area to be: %v, instead we got: %v", firstPageSaved.cellArea, firstPageExpectedCellArea)
	}

	secondPageExpectedCellAreaParsed := append([][]byte{}, cellAreaParsed[0])
	secondPageExpectedCellAreaParsed = append(secondPageExpectedCellAreaParsed, cellAreaParsed[1])

	secondPageExpectedCellArea := append([]byte{}, cellAreaParsed[0]...)
	secondPageExpectedCellArea = append(secondPageExpectedCellArea, cellAreaParsed[1]...)

	if !reflect.DeepEqual(secondPageSaved.cellAreaParsed, secondPageExpectedCellAreaParsed) {
		t.Errorf("expected cell area parsed to be: %v, instead we got: %v", secondPageExpectedCellAreaParsed, secondPageSaved.cellAreaParsed)
	}

	if secondPageSaved.numberofCells != 2 {
		t.Errorf("expected first page to have 2 cells, instead we have: %v", secondPageSaved.numberofCells)
	}

	if !reflect.DeepEqual(secondPageSaved.cellArea, secondPageExpectedCellArea) {
		t.Errorf("Expected first page cell area to be: %v, instead we got: %v", secondPageSaved.cellArea, secondPageExpectedCellArea)
	}
}

func TestBinarySearchInInteriorForNewValueHighestRowId(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:            header(),
		dbHeaderSize:        100,
		btreePageHeaderSize: 12,
		pageNumber:          0,
		cellAreaParsed:      [][]byte{{0, 0, 0, 1, 0, 0}, {0, 0, 0, 2, 0, 3}},
		btreeType:           int(TableBtreeInteriorCell),
		// remove right most pointer???
		rightMostpointer:     []byte{0, 0, 0, 2},
		cellArea:             []byte{0, 0, 0, 1, 0, 0, 0, 0, 0, 2, 0, 3},
		cellAreaSize:         12,
		numberofCells:        1,
		startCellContentArea: PageSize - 12,
		isOverflow:           false,
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

	server := ServerStruct{
		firstPage: zeroPage,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)

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
		pageNumber:           2,
		cellAreaParsed:       cellAreaParsed,
		btreeType:            int(TableBtreeLeafCell),
		cellArea:             cellArea,
		cellAreaSize:         3 * cell1.dataLength,
		numberofCells:        3,
		startCellContentArea: PageSize - 3*cell1.dataLength,
		isOverflow:           false,
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               true,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(secondPage), 0, "", &PageParsed{})

	found, pageNumber, cellAreaParsedIndex := binarySearch(secondPage, 0, 3)

	if !found {
		t.Errorf("Expected value to be found, as its already exists")
	}

	if pageNumber != 2 {
		t.Errorf("expected new value to be update on page 2, instead we got: %v", pageNumber)
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	server := ServerStruct{
		firstPage: zeroPage,
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
		leftSibling:          nil,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 4, 0, 0, 0, 0}, []byte{4, 3, 0, 0, 0, 0}, []byte{4, 2, 0, 0, 0, 0}, []byte{4, 1, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 7, 0, 0, 0, 0}, []byte{4, 6, 0, 0, 0, 0}, []byte{4, 5, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (3 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 12, 0, 0, 0, 0}, []byte{4, 11, 0, 0, 0, 0}, []byte{4, 10, 0, 0, 0, 0}, []byte{4, 9, 0, 0, 0, 0}, []byte{4, 8, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (5 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var fourthPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           4,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 16, 0, 0, 0, 0}, []byte{4, 15, 0, 0, 0, 0}, []byte{4, 14, 0, 0, 0, 0}, []byte{4, 13, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(thirdPage), 3, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(fourthPage), 4, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(PageParsed{}), 5, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(sixthPage), 6, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(seventhPage), 7, "", &server.firstPage)
	cell := createCell(TableBtreeLeafCell, 17)
	node := insert(17, cell, &server.firstPage, nil)

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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	server := ServerStruct{
		firstPage: zeroPage,
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
		leftSibling:          nil,
		rightSiblisng:        nil,
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
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	var firstPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 4, 0, 0, 0, 0}, []byte{4, 3, 0, 0, 0, 0}, []byte{4, 2, 0, 0, 0, 0}, []byte{4, 1, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           2,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 7, 0, 0, 0, 0}, []byte{4, 6, 0, 0, 0, 0}, []byte{4, 5, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (3 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var thirdPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           3,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 12, 0, 0, 0, 0}, []byte{4, 11, 0, 0, 0, 0}, []byte{4, 10, 0, 0, 0, 0}, []byte{4, 9, 0, 0, 0, 0}, []byte{4, 8, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (5 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}
	var fourthPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           4,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{4, 16, 0, 0, 0, 0}, []byte{4, 15, 0, 0, 0, 0}, []byte{4, 14, 0, 0, 0, 0}, []byte{4, 13, 0, 0, 0, 0}},
		startCellContentArea: PageSize - (4 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(thirdPage), 3, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(fourthPage), 4, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(PageParsed{}), 5, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(sixthPage), 6, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(seventhPage), 7, "", &server.firstPage)
	cell := createCell(TableBtreeLeafCell, 4, "Alice")
	node := insert(4, cell, &server.firstPage, nil)

	if node.pageNumber != 1 {
		t.Errorf("Insert values should be in page: %v, instead we got: %v", 4, node.pageNumber)
	}

	if len(node.cellAreaParsed) != 4 {
		fmt.Println(node.cellAreaParsed)
		t.Errorf("expected cell area to have 5 elements, instead we got: %v", len(node.cellAreaParsed))
	}

	if !reflect.DeepEqual(node.cellAreaParsed[0], cell.data) {
		t.Errorf("expected cell area start with newly added cell: %v, instead we got: %v", cell.data, node.cellAreaParsed[0])
	}

	reader := NewReader("fd")

	firstPageRead := reader.readDbPage(0)
	firstPageParsed := parseReadPage(firstPageRead, 0)
	secondPageRead := reader.readDbPage(1)
	secondPageParsed := parseReadPage(secondPageRead, 1)
	thirdPageRead := reader.readDbPage(2)
	thirdPageParsed := parseReadPage(thirdPageRead, 2)
	// fourthPageRead := reader.readDbPage(3)
	// fourthPageParsed := parseReadPage(fourthPageRead, 3)
	// fifthPageRead := reader.readDbPage(4)
	// fifthPageParsed := parseReadPage(fifthPageRead, 4)
	// sixthPageRead := reader.readDbPage(5)
	// sixthPageParsed := parseReadPage(sixthPageRead, 5)
	// seventhPageRead := reader.readDbPage(6)
	// seventhPageParsed := parseReadPage(seventhPageRead, 6)
	// eighthPageRead := reader.readDbPage(7)
	// eighthPageParsed := parseReadPage(eighthPageRead, 7)
	fmt.Println("lets see pages")
	fmt.Printf("%+v \n", firstPageParsed)
	fmt.Println("```````````````Second page ````````````````````")
	fmt.Printf("%+v \n", secondPageParsed)
	fmt.Println("```````````````third page ````````````````````")
	fmt.Printf("%+v \n", thirdPageParsed)
	fmt.Println("```````````````fourth page ````````````````````")
	// fmt.Printf("%+v \n", fourthPageParsed)
	// fmt.Println("```````````````fifth page ````````````````````")
	// fmt.Printf("%+v \n", fifthPageParsed)
	// fmt.Println("```````````````fifth page ````````````````````")
	// fmt.Printf("%+v \n", sixthPageParsed)
	// fmt.Println("```````````````seventh page ````````````````````")
	// fmt.Printf("%+v \n", seventhPageParsed)
	// fmt.Println("```````````````eight page ````````````````````")
	// fmt.Printf("%+v \n", eighthPageParsed)

}

//finish this one

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
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

	server := ServerStruct{
		firstPage: zeroPage,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)

	cell := createCell(TableBtreeLeafCell, 2, "aliceAndBob129876theEnd12345")
	insert(3, cell, &server.firstPage, nil)

	reader := NewReader("fd")

	secondPageSaved := reader.readFromMemory(2)
	zeroPageSaved := reader.readFromMemory(0)

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
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

	server := ServerStruct{
		firstPage: zeroPage,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{},
		startCellContentArea: PageSize,
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)

	cell1 := createCell(TableBtreeLeafCell, 2, "aliceAndBob129876theEnd12345")
	insert(3, cell1, &server.firstPage, nil)

	cell2 := createCell(TableBtreeLeafCell, 3, "aliceAndBob129876theEnd12345")
	insert(4, cell2, &server.firstPage, nil)

	reader := NewReader("fd")

	secondPageSaved := reader.readFromMemory(2)
	zeroPageSaved := reader.readFromMemory(0)

	zeroPageExpectedCellArea := []byte{0, 0, 0, 2, 0, 4, 0, 0, 0, 1, 0, 1}

	if !reflect.DeepEqual(zeroPageSaved.cellArea, zeroPageExpectedCellArea) {
		t.Errorf("expected cell area on zero page to be: %v, insted we got: %v", zeroPageExpectedCellArea, zeroPageSaved.cellArea)
	}

	if secondPageSaved.numberofCells != 2 {
		t.Errorf("expected to have one cell, instead we got: %v", secondPageSaved.numberofCells)
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

//start debugging this

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
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

	server := ServerStruct{
		firstPage: zeroPage,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	cell := creareARowItem(100, 1)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		numberofCells:        1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		startCellContentArea: PageSize - cell.dataLength,
		cellAreaSize:         cell.dataLength,
		cellArea:             cell.data,
		isOverflow:           false,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	usableSpacePerPage = 300

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)

	cell1 := creareARowItem(100, 2)
	insert(2, cell1, &server.firstPage, nil)

	cell2 := creareARowItem(100, 3)
	insert(3, cell2, &server.firstPage, nil)

	cell3 := creareARowItem(100, 4)
	insert(4, cell3, &server.firstPage, nil)

	reader := NewReader("fd")

	secondPageSaved := reader.readFromMemory(2)
	zeroPageSaved := reader.readFromMemory(0)

	zeroPageExpectedCellArea := []byte{0, 0, 0, 2, 0, 4, 0, 0, 0, 1, 0, 2}

	if !reflect.DeepEqual(zeroPageSaved.cellArea, zeroPageExpectedCellArea) {
		t.Errorf("expected cell area on zero page to be: %v, insted we got: %v", zeroPageExpectedCellArea, zeroPageSaved.cellArea)
	}

	if secondPageSaved.numberofCells != 2 {
		t.Errorf("expected to have one cell, instead we got: %v", secondPageSaved.numberofCells)
	}

	parsedData := dbReadparseCellArea(byte(TableBtreeLeafCell), secondPageSaved.cellArea)

	if !reflect.DeepEqual(parsedData, secondPageSaved.cellAreaParsed) {
		t.Errorf("expected parsed celle area to be equal to cell area: %v, instead we got: %v", parsedData, secondPageSaved.cellArea)
	}

	expectedCellArea := append([]byte{}, cell3.data...)
	expectedCellArea = append(expectedCellArea, cell2.data...)

	if !reflect.DeepEqual(secondPageSaved.cellArea, expectedCellArea) {
		t.Errorf("expected cell area to be: %v, got: %v", expectedCellArea, secondPageSaved.cellArea)
	}
}

// WE need to write test for distributing now 0x05 pages
// we forgot updating right pointer
// func TestInsert111(t *testing.T) {
// 	clearDbFile("test")

// 	cellArea := []byte{0, 0, 0, 14, 0, 41, 0, 0, 0, 13, 0, 39, 0, 0, 0, 12, 0, 36, 0, 0, 0, 11, 0, 33, 0, 0, 0, 10, 0, 30, 0, 0, 0, 9, 0, 27, 0, 0, 0, 8, 0, 24, 0, 0, 0, 7, 0, 21, 0, 0, 0, 6, 0, 18, 0, 0, 0, 5, 0, 15, 0, 0, 0, 4, 0, 12, 0, 0, 0, 3, 0, 9, 0, 0, 0, 2, 0, 48, 0, 0, 0, 1, 0, 3}
// 	84 + 100 + 12
// 	cellAreaParsed := dbReadparseCellArea(byte(TableBtreeInteriorCell), cellArea)
// 	fmt.Println("cell area length?")
// 	fmt.Println(len(cellArea))
// 	fmt.Println("cell area parsed??")
// 	fmt.Println(cellAreaParsed)
// 	var zeroPage = PageParsed{
// 		dbHeader:             header(),
// 		dbHeaderSize:         100,
// 		btreePageHeaderSize:  12,
// 		pageNumber:           0,
// 		cellAreaParsed:       cellAreaParsed,
// 		btreeType:            int(TableBtreeInteriorCell),
// 		rightMostpointer:     []byte{0, 0, 0, 2},
// 		cellArea:             cellArea,
// 		cellAreaSize:         12,
// 		numberofCells:        2,
// 		startCellContentArea: PageSize - 12,
// 		isOverflow:           false,
// 		leftSibling:          nil,
// 		rightSiblisng:        nil,
// 		isLeaf:               false,
// 	}

// 	server := ServerStruct{
// 		firstPage: zeroPage,
// 	}

// 	var firstPage = PageParsed{
// 		dbHeader:             DbHeader{},
// 		dbHeaderSize:         0,
// 		btreePageHeaderSize:  8,
// 		pageNumber:           1,
// 		btreeType:            int(TableBtreeLeafCell),
// 		cellAreaParsed:       [][]byte{},
// 		startCellContentArea: PageSize,
// 		isOverflow:           false,
// 		leftSibling:          nil,
// 		isLeaf:               true,
// 		rightSiblisng:        nil,
// 	}

// 	cell := creareARowItem(100, 1)
// 	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

// 	var secondPage = PageParsed{
// 		dbHeader:             DbHeader{},
// 		dbHeaderSize:         0,
// 		btreePageHeaderSize:  8,
// 		pageNumber:           1,
// 		numberofCells:        1,
// 		btreeType:            int(TableBtreeLeafCell),
// 		cellAreaParsed:       cellParsed,
// 		startCellContentArea: PageSize - cell.dataLength,
// 		cellAreaSize:         cell.dataLength,
// 		cellArea:             cell.data,
// 		isOverflow:           false,
// 		leftSibling:          nil,
// 		isLeaf:               true,
// 		rightSiblisng:        nil,
// 	}

// 	writer := NewWriter()

// 	usableSpacePerPage = 300

// 	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
// 	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
// 	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)

// 	// for i := 2; i < 49; i++ {
// 	// 	fmt.Println("iteration", i)
// 	// 	cell1 := creareARowItem(100, i)
// 	// 	insert(i, cell1, &server.firstPage)

// 	// }
// 	// reader := NewReader("")
// 	// zeroPageSaved := reader.readFromMemory(0)
// 	// firstPageSaved := reader.readFromMemory(1)
// 	// secondPageSaved := reader.readFromMemory(2)
// 	// fifteenPageSaved := reader.readFromMemory(14)
// 	// // sixteenPageSaved := reader.readFromMemory(16)
// 	// fmt.Println("zero page")
// 	// fmt.Println(zeroPageSaved.cellAreaParsed)
// 	// fmt.Println("first page")
// 	// fmt.Println(firstPageSaved.cellAreaParsed)
// 	// fmt.Println("secpmd page")
// 	// fmt.Println(secondPageSaved.cellAreaParsed)
// 	// fmt.Println("fiften page")
// 	// fmt.Println(fifteenPageSaved.cellAreaParsed)
// 	// fmt.Println("sixteen page")
// 	// // fmt.Println(sixteenPageSaved.cellAreaParsed)

//		// fmt.Println("end :)")
//	}

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
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

	server := ServerStruct{
		firstPage: zeroPage,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	cell := creareARowItem(100, 1)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		numberofCells:        1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		startCellContentArea: PageSize - cell.dataLength,
		cellAreaSize:         cell.dataLength,
		cellArea:             cell.data,
		isOverflow:           false,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)

	for i := 2; i < 73; i++ {
		fmt.Println("iteration", i)
		cell1 := creareARowItem(100, i)
		insert(i, cell1, &server.firstPage, nil)

	}
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// for i := 88; i < 89; i++ {
	// 	fmt.Println("iteration", i)
	// 	cell1 := creareARowItem(100, i)
	// 	insert(i, cell1, &server.firstPage, nil)

	// }

	// 23 * 6 = 138 + 23 *2 = 138+46 = 184 + 8 = 192 + header (100) = 292
	// we can max fit 23, then we need to split
	// 36 * 6 = 216 + 36*2 = 216+72 = 288 + 8 = 296

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
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

	server := ServerStruct{
		firstPage: zeroPage,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	cell := creareARowItem(100, 1)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		numberofCells:        1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		startCellContentArea: PageSize - cell.dataLength,
		cellAreaSize:         cell.dataLength,
		cellArea:             cell.data,
		isOverflow:           false,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)

	for i := 2; i < 74; i++ {
		fmt.Println("iteration", i)
		cell1 := creareARowItem(100, i)
		insert(i, cell1, &server.firstPage, nil)

	}
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	reader := NewReader("")
	zeroPageSaved := reader.readFromMemory(0)
	twentyFifthPageSaved := reader.readFromMemory(25)
	thirtEightPageSaved := reader.readFromMemory(39)

	fmt.Println("len", len(twentyFifthPageSaved.cellAreaParsed))

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

//todo finish this

func TestInsertWithInteriorNestedSplittedSix(t *testing.T) {
	clearDbFile("test")
	PageSize = 211
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
		leftSibling:          nil,
		rightSiblisng:        nil,
		isLeaf:               false,
	}

	server := ServerStruct{
		firstPage: zeroPage,
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
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	cell := creareARowItem(100, 1)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)

	// 208 = 208 -12 = 196 /8 = 23  pointers, on every page 24 pointers
	// data 200 = 200 -8 = 192/100 = 1, per page, hmm
	var secondPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreePageHeaderSize:  8,
		pageNumber:           1,
		numberofCells:        1,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       cellParsed,
		startCellContentArea: PageSize - cell.dataLength,
		cellAreaSize:         cell.dataLength,
		cellArea:             cell.data,
		isOverflow:           false,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)
	// [0 0 0 27 0 39] [0 0 0 14 0 24] for 40
	// [0 0 0 27 0 48] [0 0 0 14 0 24] for 49
	// [0 0 0 52 0 49] [0 0 0 27 0 36] [0 0 0 14 0 23] for 50
	// 49 total for A-24, B-13 C-13
	//[0 0 0 52 0 50] [0 0 0 27 0 36] [0 0 0 14 0 23] for 51, there is no full factor, because we always need to have at least 50% in the right, so we are taking from the left
	// size 10
	// A-11
	// A -6 B-11
	// A-10 B-7
	// A-10 B-10
	// A-10 B-11
	// A-10 B-6 C-5
	for i := 2; i < 49; i++ {
		fmt.Println("iteration", i)
		cell1 := creareARowItem(100, i)
		insert(i, cell1, &server.firstPage, nil)

	}
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for i := 49; i < 120; i++ {
		fmt.Println("iteration", i)
		cell1 := creareARowItem(100, i)
		insert(i, cell1, &server.firstPage, nil)

	}

	t.Error("")

	reader := NewReader("")
	zeroPageSaved := reader.readFromMemory(0)
	// twentyFifthPageSaved := reader.readFromMemory(12)

	fmt.Println("zero page")
	fmt.Printf("%+v", zeroPageSaved)

	fmt.Println("twenty page")
	// fmt.Printf("%+v", twentyFifthPageSaved)

}

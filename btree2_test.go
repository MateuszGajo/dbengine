package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestShouldFindResultSingleInteriorPlusLeafPage(t *testing.T) {
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
	found, _, page, parents := search(0, 7, []PageParsed{})

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
	found, _, page, parents = search(0, 12, []PageParsed{})

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

	found, _, page, parents = search(0, 4, []PageParsed{})

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

	found, _, page, parents = search(0, 1, []PageParsed{})

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
	found, _, page, parents := search(0, 13, []PageParsed{})

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
	found, index, pageFound, parents := search(0, 7, []PageParsed{})

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

	found, _, pageFound, parents = search(0, 12, []PageParsed{})

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
	found, _, pageFound, parents = search(0, 4, []PageParsed{})

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
	found, _, pageFound, parents = search(0, 1, []PageParsed{})

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
	found, _, pageFound, parents = search(0, 16, []PageParsed{})

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
		rightMostpointer:     []byte{},
		cellArea:             cellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	cells := []Cell{{rowId: 9, pageNumber: 2}, {rowId: 13, pageNumber: 3}}

	updateDivider(zeroPage, cells, 6, 6*3, &zeroPage)

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
}

func TestUpdateDividerLessAmount(t *testing.T) {
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
		rightMostpointer:     []byte{},
		cellArea:             cellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 6*4,
		isOverflow:           true,
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	cells := []Cell{{rowId: 9, pageNumber: 2}}

	updateDivider(zeroPage, cells, 6, 6*3, &zeroPage)

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
}

func TestUpdateDividerGreaterAmount(t *testing.T) {
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
		rightMostpointer:     []byte{},
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

	updateDivider(zeroPage, cells, 6, 6*3, &zeroPage)

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

func TestBalancingSplitRootIntTwoChildren(t *testing.T) {
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

	usableSpacePerPage = 9 * 3
	server := ServerStruct{
		firstPage: zeroPage,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	balancingForNode(zeroPage, []PageParsed{}, &server.firstPage)
	writer.flushPages("", &server.firstPage)

	reader := NewReader("")

	zeroPageSaved := reader.readFromMemory(0)
	firstPageSaved := reader.readFromMemory(1)
	secondPageSaved := reader.readFromMemory(2)

	zeroPageExpectedCellArea := []byte{0, 0, 0, 2, 0, 4, 0, 0, 0, 1, 0, 2}
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
	balancingForNode(zeroPage, []PageParsed{}, &server.firstPage)
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
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		cellAreaParsed:       [][]byte{{0, 0, 0, 1, 0, 4}},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{},
		cellArea:             []byte{0, 0, 0, 1, 0, 4},
		isLeaf:               false,
		startCellContentArea: PageSize - 9,
		numberofCells:        1,
		isOverflow:           false,
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	var firstPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         0,
		pageNumber:           1,
		cellAreaParsed:       cellAreaParsed,
		btreeType:            int(TableBtreeLeafCell),
		rightMostpointer:     []byte{},
		cellArea:             cellArea,
		isLeaf:               true,
		startCellContentArea: PageSize - 9*4,
		numberofCells:        len(cellAreaParsed),
		isOverflow:           true,
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	usableSpacePerPage = 9 * 3
	server := ServerStruct{
		firstPage: zeroPage,
	}

	writer := NewWriter()
	reader := NewReader("")

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)

	balancingForNode(firstPage, []PageParsed{zeroPage}, &server.firstPage)
	writer.flushPages("", &server.firstPage)

	zeroPageSaved := reader.readFromMemory(0)
	firstPageSaved := reader.readFromMemory(1)
	secondPageSaved := reader.readFromMemory(2)

	fmt.Println("number of pages??")
	fmt.Println(server.firstPage.dbHeader.dbSizeInPages)

	if server.firstPage.dbHeader.dbSizeInPages != 3 {
		t.Errorf("Expected number of pages to be 3, got: %v", server.firstPage.dbHeader.dbSizeInPages)
	}

	zeroPageExpectedCellArea := []byte{0, 0, 0, 2, 0, 4, 0, 0, 0, 1, 0, 2}
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

func TestInsert(t *testing.T) {
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
	cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: 17}})
	node := insert(17, cell, &server.firstPage)

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
	cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: 4}}, "Alice")
	node := insert(4, cell, &server.firstPage)

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

func TestInsertOneRecord(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:            header(),
		dbHeaderSize:        100,
		btreePageHeaderSize: 12,
		pageNumber:          0,
		cellAreaParsed:      [][]byte{},
		btreeType:           int(TableBtreeInteriorCell),
		rightMostpointer:    []byte{0, 0, 0, 2},
		cellArea:            []byte{0, 0, 0, 1, 0, 0},

		startCellContentArea: PageSize,
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

	// using 0 as row id makes infinite loop?

	cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: 1}}, "aliceAndBob129876theEnd12345")
	insert(1, cell, &server.firstPage)
	writer.flushPages("", &server.firstPage)

	reader := NewReader("fd")

	// secondPageRead := reader.readDbPage(1)
	// secondPageParsed := parseReadPage(secondPageRead, 1)

	thirdPageRead := reader.readDbPage(2)
	fmt.Println("row data?")
	fmt.Println(thirdPageRead)
	thirdPageParsed := parseReadPage(thirdPageRead, 2)

	if thirdPageParsed.numberofCells != 1 {
		t.Errorf("expected to have one cell")
	}

	if len(thirdPageParsed.cellAreaParsed) != 1 {
		t.Errorf("expected to have one cell area")
	}
	if !reflect.DeepEqual(thirdPageParsed.cellAreaParsed[0], thirdPageParsed.cellArea) {
		t.Errorf("expected parsed celle area to be equal to cell area")
	}

	if !reflect.DeepEqual(thirdPageParsed.cellArea, cell.data) {
		t.Errorf("expected cell area to be: %v, got: %v", cell.data, thirdPageParsed.cellArea)
	}
}

// Debug this, its failing
// finish writing test, updating index for parents
// test for leaf iteself
func TestInsertOverflowPage(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:            header(),
		dbHeaderSize:        100,
		btreePageHeaderSize: 12,
		pageNumber:          0,
		cellAreaParsed:      [][]byte{},
		btreeType:           int(TableBtreeInteriorCell),
		rightMostpointer:    []byte{0, 0, 0, 2},
		cellArea:            []byte{0, 0, 0, 1, 0, 0},

		startCellContentArea: PageSize - 6,
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

	var secondPage = PageParsed{
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

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(secondPage), 2, "", &server.firstPage)

	reader := NewReader("")
	for i := 1; i < 122; i++ {

		cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: i - 1}}, "aliceAndBob129876theEnd12345")
		insert(i, cell, &server.firstPage)
		writer.flushPages("", &server.firstPage)

	}
	cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: 1}}, "aliceAndBob129876theEnd12345")

	secondPageRead := reader.readDbPage(2)
	secondPageParsed := parseReadPage(secondPageRead, 2)
	firstPageRead := reader.readDbPage(1)
	firstPageParsed := parseReadPage(firstPageRead, 1)
	zeroPageRead := reader.readDbPage(0)
	zeroPageParsed := parseReadPage(zeroPageRead, 0)

	zeroPageExpectedCellArea := []byte{0, 0, 0, 1, 0, 61}

	if !reflect.DeepEqual(zeroPageParsed.cellArea, zeroPageExpectedCellArea) {
		t.Errorf("expected cell area on zero page to be: %v, insted we got: %v", zeroPageExpectedCellArea, zeroPageParsed.cellArea)
	}
	if firstPageParsed.btreeType != int(TableBtreeLeafCell) {
		t.Errorf("first page tree type should be table leaf, insted we got: %v", firstPageParsed.btreeType)
	}

	if secondPageParsed.btreeType != int(TableBtreeLeafCell) {
		t.Errorf("second page tree type should be table leaf, insted we got: %v", secondPageParsed.btreeType)
	}

	if secondPageParsed.numberofCells != len(secondPageParsed.cellAreaParsed) {
		t.Errorf("second page, number of cell should be equal to parsed len, expected: %v, got: %v", secondPageParsed.numberofCells, len(secondPageParsed.cellAreaParsed))
	}

	if firstPageParsed.numberofCells != len(firstPageParsed.cellAreaParsed) {
		t.Errorf("first page, number of cell should be equal to parsed len, expected: %v, got: %v", firstPageParsed.numberofCells, len(firstPageParsed.cellAreaParsed))
	}

	firstPageExpectedFirstCell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: 0}}, "aliceAndBob129876theEnd12345")
	firstPageExpectedLastCell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: 60}}, "aliceAndBob129876theEnd12345")

	if !reflect.DeepEqual(firstPageParsed.cellAreaParsed[len(firstPageParsed.cellAreaParsed)-1], firstPageExpectedFirstCell.data) {
		t.Errorf("Expected first cell in first page to be: %v, insted we got: %v", firstPageExpectedFirstCell.data, firstPageParsed.cellAreaParsed[len(firstPageParsed.cellAreaParsed)-1])
	}

	if !reflect.DeepEqual(firstPageParsed.cellAreaParsed[0], firstPageExpectedLastCell.data) {
		t.Errorf("Expected last cell in first page to be: %v, insted we got: %v", firstPageExpectedLastCell.data, firstPageParsed.cellAreaParsed[0])
	}

	secondPageExpectedFirstCell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: 61}}, "aliceAndBob129876theEnd12345")
	secondPageExpectedLastCell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: 120}}, "aliceAndBob129876theEnd12345")

	if !reflect.DeepEqual(secondPageParsed.cellAreaParsed[len(secondPageParsed.cellAreaParsed)-1], secondPageExpectedFirstCell.data) {
		t.Errorf("Expected first cell in first page to be: %v, insted we got: %v", secondPageExpectedFirstCell.data, secondPageParsed.cellAreaParsed[len(secondPageParsed.cellAreaParsed)-1])
	}

	if !reflect.DeepEqual(secondPageParsed.cellAreaParsed[0], secondPageExpectedLastCell.data) {
		t.Errorf("Expected last cell in first page to be: %v, insted we got: %v", secondPageExpectedLastCell.data, secondPageParsed.cellAreaParsed[0])
	}

	if len(secondPageParsed.cellAreaParsed) != 60 {
		t.Errorf("expected page to have %v of cells, got: %v", 60, len(secondPageParsed.cellAreaParsed))
	}

	if len(firstPageParsed.cellAreaParsed) != 61 {
		t.Errorf("expected page to have %v of cells, got: %v", 61, len(firstPageParsed.cellAreaParsed))
	}

	firstPageExpectedStartCellArea := PageSize - (61 * cell.dataLength)

	if firstPageParsed.startCellContentArea != firstPageExpectedStartCellArea {
		t.Errorf("expected cell content area to be: %v, got: %v", firstPageExpectedStartCellArea, firstPageParsed.startCellContentArea)
	}

	secondPageExpectedStartCellArea := PageSize - (60 * cell.dataLength)

	if secondPageParsed.startCellContentArea != secondPageExpectedStartCellArea {
		t.Errorf("expected cell content area to be: %v, got: %v", secondPageExpectedStartCellArea, secondPageParsed.startCellContentArea)
	}

}

// probalby need to fix saving, for root page, to not only have right most pointer
func TestInsertOverflowPage2(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:            header(),
		dbHeaderSize:        100,
		btreePageHeaderSize: 8,
		pageNumber:          0,
		cellAreaParsed:      [][]byte{},
		btreeType:           int(TableBtreeLeafCell),
		rightMostpointer:    []byte{},
		cellArea:            []byte{},

		startCellContentArea: PageSize,
		isOverflow:           false,
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	server := ServerStruct{
		firstPage: zeroPage,
	}

	writer := NewWriter()

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)

	reader := NewReader("")
	for i := 1; i < 119; i++ {

		cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: i - 1}}, "aliceAndBob129876theEnd12345")
		insert(i, cell, &server.firstPage)
		writer.flushPages("", &server.firstPage)

	}
	// cell := createCell(TableBtreeLeafCell, &PageParsed{latesRow: &LastPageParseLatestRow{rowId: 1}}, "aliceAndBob129876theEnd12345")

	// secondPageRead := reader.readDbPage(2)
	// secondPageParsed := parseReadPage(secondPageRead, 2)
	firstPageRead := reader.readDbPage(1)
	firstPageParsed := parseReadPage(firstPageRead, 1)
	zeroPageRead := reader.readDbPage(0)
	zeroPageParsed := parseReadPage(zeroPageRead, 0)

	fmt.Println("zero page")
	fmt.Printf("%+v", zeroPageParsed)
	fmt.Println("first page")
	fmt.Printf("%+v", firstPageParsed)
	// fmt.Println("second page")
	// fmt.Printf("%+v", secondPageParsed)
	t.Error("")

}

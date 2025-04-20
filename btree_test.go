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

type SearchForDataCase struct {
	startPage             int
	lookingValue          int
	expectedParentNumbers []int
	expectedPage          int
}

func TestShouldFindResultSingleInteriorPlusLeafPage(t *testing.T) {
	clearDbFile("test")

	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 7}, {0, 0, 0, 1, 0, 4}}, 0, &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 4, 0, 0, 0, 0}, {4, 3, 0, 0, 0, 0}, {4, 2, 0, 0, 0, 0}, {4, 1, 0, 0, 0, 0}}, 1, nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 7, 0, 0, 0, 0}, {4, 6, 0, 0, 0, 0}, {4, 5, 0, 0, 0, 0}}, 2, nil)
	thirdPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 12, 0, 0, 0, 0}, {4, 11, 0, 0, 0, 0}, {4, 10, 0, 0, 0, 0}, {4, 9, 0, 0, 0, 0}, {4, 8, 0, 0, 0, 0}}, 3, nil)

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage

	var testCases = []SearchForDataCase{{startPage: 0, lookingValue: 7, expectedParentNumbers: []int{0}, expectedPage: 2},
		{startPage: 0, lookingValue: 12, expectedParentNumbers: []int{0}, expectedPage: 3},
		{startPage: 0, lookingValue: 4, expectedParentNumbers: []int{0}, expectedPage: 1},
		{startPage: 0, lookingValue: 1, expectedParentNumbers: []int{0}, expectedPage: 1},
	}

	for _, testCase := range testCases {
		found, _, pageFound, parents := search(testCase.startPage, testCase.lookingValue, []*PageParsed{})

		if !found {
			t.Errorf("expect to find looking value")
		}

		if pageFound.pageNumber != testCase.expectedPage {
			t.Errorf("expected to found result on page: %v, instead we got: %v", 2, pageFound.pageNumber)
		}

		if len(parents) != len(testCase.expectedParentNumbers) {
			t.Errorf("Expected to be only 2 parents, instead we got: %v", len(parents))
		}

		for i, expectedParent := range testCase.expectedParentNumbers {
			if !reflect.DeepEqual(parents[i].pageNumber, expectedParent) {
				t.Errorf("expect parent to be page: %v, instead we got: %v", parents[i].pageNumber, expectedParent)
			}
		}
	}

}

func TestShouldNotFindResultSingleInteriorPlusLeafPage(t *testing.T) {
	clearDbFile("test")

	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 7}, {0, 0, 0, 1, 0, 4}}, 0, &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 4, 0, 0, 0, 0}, {4, 3, 0, 0, 0, 0}, {4, 2, 0, 0, 0, 0}, {4, 1, 0, 0, 0, 0}}, 1, nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 7, 0, 0, 0, 0}, {4, 6, 0, 0, 0, 0}, {4, 5, 0, 0, 0, 0}}, 2, nil)
	thirdPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 12, 0, 0, 0, 0}, {4, 11, 0, 0, 0, 0}, {4, 10, 0, 0, 0, 0}, {4, 9, 0, 0, 0, 0}, {4, 8, 0, 0, 0, 0}}, 3, nil)

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage

	var testCases = []SearchForDataCase{{startPage: 0, lookingValue: 13, expectedParentNumbers: []int{0}, expectedPage: 3}}

	for _, testCase := range testCases {
		found, _, pageFound, parents := search(testCase.startPage, testCase.lookingValue, []*PageParsed{})

		if found {
			t.Errorf("Shouldn't find any result")
		}

		if pageFound.pageNumber != testCase.expectedPage {
			t.Errorf("expected to insert new value at page : %v, instead we got: %v", testCase.expectedPage, pageFound.pageNumber)
		}

		if len(parents) != len(testCase.expectedParentNumbers) {
			t.Errorf("Expected to be only 2 parents, instead we got: %v", len(parents))
		}

		for i, expectedParent := range testCase.expectedParentNumbers {
			if !reflect.DeepEqual(parents[i].pageNumber, expectedParent) {
				t.Errorf("expect parent to be page: %v, instead we got: %v", parents[i].pageNumber, expectedParent)
			}
		}
	}

}

func TestShouldFindResultMultipleInteriorPlusLeafPage(t *testing.T) {
	clearDbFile("test")
	header := header()

	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 7, 0, 12}, {0, 0, 0, 6, 0, 7}}, 0, &header)
	sixthPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 7}, {0, 0, 0, 1, 0, 4}}, 6, nil)
	seventhPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}}, 7, nil)
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 4, 0, 0, 0, 0}, {4, 3, 0, 0, 0, 0}, {4, 2, 0, 0, 0, 0}, {4, 1, 0, 0, 0, 0}}, 1, nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 7, 0, 0, 0, 0}, {4, 6, 0, 0, 0, 0}, {4, 5, 0, 0, 0, 0}}, 2, nil)
	thirdPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 12, 0, 0, 0, 0}, {4, 11, 0, 0, 0, 0}, {4, 10, 0, 0, 0, 0}, {4, 9, 0, 0, 0, 0}, {4, 8, 0, 0, 0, 0}}, 3, nil)
	fourthPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 16, 0, 0, 0, 0}, {4, 15, 0, 0, 0, 0}, {4, 14, 0, 0, 0, 0}, {4, 13, 0, 0, 0, 0}}, 4, nil)
	fifthPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, 5, nil)

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage
	memoryPages[fourthPage.pageNumber] = fourthPage
	memoryPages[fifthPage.pageNumber] = fifthPage
	memoryPages[sixthPage.pageNumber] = sixthPage
	memoryPages[seventhPage.pageNumber] = seventhPage

	var testCases = []SearchForDataCase{{startPage: 0, lookingValue: 7, expectedParentNumbers: []int{0, 6}, expectedPage: 2},
		{startPage: 0, lookingValue: 12, expectedParentNumbers: []int{0, 7}, expectedPage: 3},
		{startPage: 0, lookingValue: 4, expectedParentNumbers: []int{0, 6}, expectedPage: 1},
		{startPage: 0, lookingValue: 1, expectedParentNumbers: []int{0, 6}, expectedPage: 1},
		{startPage: 0, lookingValue: 16, expectedParentNumbers: []int{0, 7}, expectedPage: 4},
	}

	for _, testCase := range testCases {
		found, _, pageFound, parents := search(testCase.startPage, testCase.lookingValue, []*PageParsed{})

		if !found {
			t.Errorf("expect to find looking value")
		}

		if pageFound.pageNumber != testCase.expectedPage {
			t.Errorf("expected to found result on page: %v, instead we got: %v", 2, pageFound.pageNumber)
		}

		if len(parents) != len(testCase.expectedParentNumbers) {
			t.Errorf("Expected to be only 2 parents, instead we got: %v", len(parents))
		}

		for i, expectedParent := range testCase.expectedParentNumbers {
			if !reflect.DeepEqual(parents[i].pageNumber, expectedParent) {
				t.Errorf("expect parent to be page: %v, instead we got: %v", parents[i].pageNumber, expectedParent)
			}
		}
	}
}

func TestUpdateDividerSameAmount(t *testing.T) {
	clearDbFile("test")

	cellAreaParsed := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 8}, {0, 0, 0, 1, 0, 5}}
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), cellAreaParsed, 0, &header)
	cells := []Cell{{rowId: 9, pageNumber: 2}, {rowId: 13, pageNumber: 3}}

	cellAreaParsedExpected := cellAreaParsed
	cellAreaParsedExpected[1] = []byte{0, 0, 0, 3, 0, 13}
	cellAreaParsedExpected[2] = []byte{0, 0, 0, 2, 0, 9}

	modifyDivider(&zeroPage, cells, 6, 6*3, &zeroPage.dbHeader, []*PageParsed{})

	utilsTestContent(t, cellAreaParsedExpected, TableBtreeInteriorCell, 0, zeroPage)
}

func TestUpdateDividerLessAmount(t *testing.T) {
	clearDbFile("test")

	cellAreaParsed := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 8}, {0, 0, 0, 1, 0, 5}}
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), cellAreaParsed, 0, &header)
	cells := []Cell{{rowId: 9, pageNumber: 2}}

	cellAreaParsedExpected := [][]byte{cellAreaParsed[0]}
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 2, 0, 9})
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[3])

	modifyDivider(&zeroPage, cells, 6, 6*3, &zeroPage.dbHeader, []*PageParsed{})

	utilsTestContent(t, cellAreaParsedExpected, TableBtreeInteriorCell, 0, zeroPage)

}

func TestUpdateDividerGreaterAmount(t *testing.T) {
	clearDbFile("test")

	cellAreaParsed := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 8}, {0, 0, 0, 1, 0, 5}}
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), cellAreaParsed, 0, &header)
	cells := []Cell{{rowId: 9, pageNumber: 2}, {rowId: 13, pageNumber: 3}, {rowId: 14, pageNumber: 5}}

	cellAreaParsedExpected := [][]byte{cellAreaParsed[0]}
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 5, 0, 14})
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 3, 0, 13})
	cellAreaParsedExpected = append(cellAreaParsedExpected, []byte{0, 0, 0, 2, 0, 9})
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[3])

	modifyDivider(&zeroPage, cells, 6, 6*3, &zeroPage.dbHeader, []*PageParsed{})

	utilsTestContent(t, cellAreaParsedExpected, TableBtreeInteriorCell, 0, zeroPage)
}

func TestUpdateTestRightMostPointerWhenUpdatePointerWithNewPage(t *testing.T) {
	clearDbFile("test")

	cellAreaParsed := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 8}, {0, 0, 0, 1, 0, 5}}
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), cellAreaParsed, 0, &header)
	cells := []Cell{{rowId: 18, pageNumber: 5}}

	cellAreaParsedExpected := [][]byte{{0, 0, 0, 5, 0, 18}}
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[1])
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[2])
	cellAreaParsedExpected = append(cellAreaParsedExpected, cellAreaParsed[3])

	modifyDivider(&zeroPage, cells, 0, 6, &zeroPage.dbHeader, []*PageParsed{})

	utilsTestContent(t, cellAreaParsedExpected, TableBtreeInteriorCell, 0, zeroPage)
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
	clearDbFile("test")

	zeroPageParsedCellArea := [][]byte{{0, 0, 0, 1, 0, 16}}
	firstPageParsedCellArea := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 2, 0, 8}, {0, 0, 0, 8, 0, 5}}
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), zeroPageParsedCellArea, 0, &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), firstPageParsedCellArea, 1, nil)
	cells := []Cell{{rowId: 18, pageNumber: 5}}

	modifyDivider(&firstPage, cells, 0, 6, &zeroPage.dbHeader, []*PageParsed{&zeroPage})

	utilsTestContent(t, [][]byte{{0, 0, 0, 1, 0, 18}}, TableBtreeInteriorCell, 0, zeroPage)

}

func TestUpdateDividerWithParentandGrandParentUpdate(t *testing.T) {
	clearDbFile("test")

	zeroPageParsedCellArea := [][]byte{{0, 0, 0, 1, 0, 16}}
	firstPageParsedCellArea := [][]byte{{0, 0, 0, 2, 0, 16}}
	secondPageParsedCellArea := [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}, {0, 0, 0, 9, 0, 8}, {0, 0, 0, 8, 0, 5}}
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), zeroPageParsedCellArea, 0, &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), firstPageParsedCellArea, 1, nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), secondPageParsedCellArea, 2, nil)
	cells := []Cell{{rowId: 18, pageNumber: 5}}

	modifyDivider(&secondPage, cells, 0, 6, &zeroPage.dbHeader, []*PageParsed{&zeroPage, &firstPage})

	utilsTestContent(t, [][]byte{{0, 0, 0, 1, 0, 18}}, TableBtreeInteriorCell, 0, zeroPage)
	utilsTestContent(t, [][]byte{{0, 0, 0, 2, 0, 18}}, TableBtreeInteriorCell, 1, firstPage)

}

func TestFindSiblingOnlyRightSiblingsAsLastPointer(t *testing.T) {
	clearDbFile("test")

	cell := creareARowItem(100, 2)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 2}, {0, 0, 0, 1, 0, 1}}, 0, &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), cellParsed, 1, nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, 2, nil)

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

	cell := creareARowItem(100, 3)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 2}, {0, 0, 0, 1, 0, 1}}, 0, &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, 1, nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), cellParsed, 2, nil)

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

	cell := creareARowItem(100, 3)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 3, 0, 3}, {0, 0, 0, 2, 0, 2}, {0, 0, 0, 1, 0, 1}}, 0, &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, 1, nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), cellParsed, 2, nil)
	thirdPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, 3, nil)

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

	cell := creareARowItem(100, 2)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 1, 0, 1}}, 0, &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), cellParsed, 1, nil)

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
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), cellParsed, 1, nil)

	newCell := createCell(TableBtreeLeafCell, 0, "Alice")
	firstPage.updateParsedCells(newCell, 0)

	utilsTestContent(t, [][]byte{newCell.data}, TableBtreeLeafCell, 1, firstPage)
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
	PageSize = 9*2 + 2*2 + 8

	cellAreaParsed := [][]byte{{7, 4, 2, 23, 65, 108, 105, 99, 101}, {7, 3, 2, 23, 65, 108, 105, 99, 101}, {7, 2, 2, 23, 65, 108, 105, 99, 101}, {7, 1, 2, 23, 65, 108, 105, 99, 101}}
	zeroPage := PageParsed{
		dbHeader: DbHeader{
			dbSizeInPages: 2,
		},
	}
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), cellAreaParsed, 1, nil)
	firstPage.isOverflow = true

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
	header := header()
	zeroPage := CreateNewPage(BtreeType(TableBtreeLeafCell), cellAreaParsed, 0, &header)
	zeroPage.isOverflow = true
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

	utilsTestContent(t, [][]byte{{0, 0, 0, 1, 0, 4}}, TableBtreeInteriorCell, 0, zeroPageSaved)
	utilsTestContent(t, cellAreaParsed, TableBtreeLeafCell, 1, firstPageSaved)
}

var PointerToDataLength = 2

type BtreeHeaderSize int

var (
	BtreeHeaderSizeLeafTable BtreeHeaderSize = 8
)

func TestBalancingSplitOneLeaftPageIntoTwo(t *testing.T) {
	clearDbFile("test")

	cellAreaParsed := [][]byte{{7, 4, 2, 23, 65, 108, 105, 99, 101}, {7, 3, 2, 23, 65, 108, 105, 99, 101}, {7, 2, 2, 23, 65, 108, 105, 99, 101}, {7, 1, 2, 23, 65, 108, 105, 99, 101}}
	PageSize = len(cellAreaParsed[0])*2 + PointerToDataLength*2 + int(BtreeHeaderSizeLeafTable)
	server := ServerStruct{
		header: DbHeader{dbSizeInPages: 1},
	}
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

	header := DbHeader{}
	var zeroPage = CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 3}, {0, 0, 0, 1, 0, 0}}, header.assignNewPage(), &header)

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
	cellAreaParsed := append([][]byte{}, cell3.data)
	cellAreaParsed = append(cellAreaParsed, cell2.data)
	cellAreaParsed = append(cellAreaParsed, cell1.data)
	var zeroPage = CreateNewPage(BtreeType(TableBtreeLeafCell), cellAreaParsed, 0, nil)
	memoryPages[zeroPage.pageNumber] = zeroPage

	found, pageNumber, cellAreaParsedIndex := binarySearch(zeroPage, 0, 2)

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

	header := DbHeader{}
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 7, 0, 16}, {0, 0, 0, 6, 0, 7}}, header.assignNewPage(), &header)
	server := ServerStruct{
		header: header,
	}
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 4, 0, 0, 0, 0}, {4, 3, 0, 0, 0, 0}, {4, 2, 0, 0, 0, 0}, {4, 1, 0, 0, 0, 0}}, header.assignNewPage(), nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 7, 0, 0, 0, 0}, {4, 6, 0, 0, 0, 0}, {4, 5, 0, 0, 0, 0}}, header.assignNewPage(), nil)
	thirdPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 12, 0, 0, 0, 0}, {4, 11, 0, 0, 0, 0}, {4, 10, 0, 0, 0, 0}, {4, 9, 0, 0, 0, 0}, {4, 8, 0, 0, 0, 0}}, header.assignNewPage(), nil)
	fourthPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 16, 0, 0, 0, 0}, {4, 15, 0, 0, 0, 0}, {4, 14, 0, 0, 0, 0}, {4, 13, 0, 0, 0, 0}}, header.assignNewPage(), nil)
	fifthPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, header.assignNewPage(), nil)
	sixthPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 7}, {0, 0, 0, 1, 0, 4}}, header.assignNewPage(), nil)
	seventhPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}}, header.assignNewPage(), nil)

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage
	memoryPages[fourthPage.pageNumber] = fourthPage
	memoryPages[fifthPage.pageNumber] = fifthPage
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

func TestInsertToExisting(t *testing.T) {
	clearDbFile("test")

	header := DbHeader{}
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 7, 0, 16}, {0, 0, 0, 6, 0, 7}}, header.assignNewPage(), &header)
	server := ServerStruct{
		header: header,
	}
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 4, 0, 0, 0, 0}, {4, 3, 0, 0, 0, 0}, {4, 2, 0, 0, 0, 0}, {4, 1, 0, 0, 0, 0}}, header.assignNewPage(), nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 7, 0, 0, 0, 0}, {4, 6, 0, 0, 0, 0}, {4, 5, 0, 0, 0, 0}}, header.assignNewPage(), nil)
	thirdPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 12, 0, 0, 0, 0}, {4, 11, 0, 0, 0, 0}, {4, 10, 0, 0, 0, 0}, {4, 9, 0, 0, 0, 0}, {4, 8, 0, 0, 0, 0}}, header.assignNewPage(), nil)
	fourthPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{{4, 16, 0, 0, 0, 0}, {4, 15, 0, 0, 0, 0}, {4, 14, 0, 0, 0, 0}, {4, 13, 0, 0, 0, 0}}, header.assignNewPage(), nil)
	fifthPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, header.assignNewPage(), nil)
	sixthPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 7}, {0, 0, 0, 1, 0, 4}}, header.assignNewPage(), nil)
	seventhPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 4, 0, 16}, {0, 0, 0, 3, 0, 12}}, header.assignNewPage(), nil)

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	memoryPages[thirdPage.pageNumber] = thirdPage
	memoryPages[fourthPage.pageNumber] = fourthPage
	memoryPages[fifthPage.pageNumber] = fifthPage
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

func TestInsertOneRecord(t *testing.T) {
	clearDbFile("test")
	header := DbHeader{}
	server := ServerStruct{
		header: header,
	}
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 2}, {0, 0, 0, 1, 0, 1}}, header.assignNewPage(), &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, header.assignNewPage(), nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, header.assignNewPage(), nil)

	memoryPages[zeroPage.pageNumber] = zeroPage
	memoryPages[firstPage.pageNumber] = firstPage
	memoryPages[secondPage.pageNumber] = secondPage
	btree := BtreeStruct{
		softPages: map[int]PageParsed{},
	}

	rowId := 3
	cell := createCell(TableBtreeLeafCell, rowId, "aliceAndBob129876theEnd12345")
	btree.insert(rowId, cell, &server.header, nil)

	secondPageSaved := btree.softPages[2]
	zeroPageSaved := btree.softPages[0]

	zeroPageExpectedCellArea := []byte{0, 0, 0, 2, 0, 3, 0, 0, 0, 1, 0, 1}

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
	header := DbHeader{}
	server := ServerStruct{
		header: header,
	}
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 2}, {0, 0, 0, 1, 0, 1}}, header.assignNewPage(), &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, header.assignNewPage(), nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, header.assignNewPage(), nil)

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

	zeroPageExpectedCellAreaParsed := [][]byte{{0, 0, 0, 2, 0, 4}, {0, 0, 0, 1, 0, 1}}
	secondPageExpectedCellAreaParsed := append([][]byte{}, cell2.data)
	secondPageExpectedCellAreaParsed = append(secondPageExpectedCellAreaParsed, cell1.data)

	utilsTestContent(t, zeroPageExpectedCellAreaParsed, TableBtreeInteriorCell, 0, zeroPageSaved)
	utilsTestContent(t, secondPageExpectedCellAreaParsed, TableBtreeLeafCell, 2, secondPageSaved)
}

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
	header := DbHeader{}
	server := ServerStruct{
		header: header,
	}
	cell := creareARowItem(100, 2)
	cellParsed := dbReadparseCellArea(byte(TableBtreeLeafCell), cell.data)
	zeroPage := CreateNewPage(BtreeType(TableBtreeInteriorCell), [][]byte{{0, 0, 0, 2, 0, 1}, {0, 0, 0, 1, 0, 0}}, header.assignNewPage(), &header)
	firstPage := CreateNewPage(BtreeType(TableBtreeLeafCell), [][]byte{}, header.assignNewPage(), nil)
	secondPage := CreateNewPage(BtreeType(TableBtreeLeafCell), cellParsed, header.assignNewPage(), nil)

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

	cell3 := creareARowItem(100, 5)
	btree.insert(4, cell3, &server.header, nil)

	secondPageSaved := btree.softPages[2]
	zeroPageSaved := btree.softPages[0]

	zeroPageExpectedCellAreaParsed := [][]byte{{0, 0, 0, 2, 0, 4}, {0, 0, 0, 1, 0, 3}}
	secondPageExpectedCellAreaParsed := append([][]byte{}, cell3.data)

	utilsTestContent(t, zeroPageExpectedCellAreaParsed, TableBtreeInteriorCell, 0, zeroPageSaved)
	utilsTestContent(t, secondPageExpectedCellAreaParsed, TableBtreeLeafCell, 2, secondPageSaved)
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

	expectedCellAreaParsedPageZero := [][]byte{{0, 0, 0, 25, 0, 72}}

	utilsTestContent(t, expectedCellAreaParsedPageZero, TableBtreeInteriorCell, 0, zeroPageSaved)

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

	expectedCellAreaContentPageZero := [][]byte{{0, 0, 0, 39, 0, 73}, {0, 0, 0, 25, 0, 38}}

	utilsTestContent(t, expectedCellAreaContentPageZero, TableBtreeInteriorCell, 0, zeroPageSaved)

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

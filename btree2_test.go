package main

import (
	"fmt"
	"testing"
)

// debug this test
func TestLeaftBiasDistribution(t *testing.T) {
	clearDbFile("test")
	var zeroPage = PageParsed{
		dbHeader:             header(),
		dbHeaderSize:         100,
		pageNumber:           0,
		cellAreaParsed:       [][]byte{},
		btreeType:            int(TableBtreeInteriorCell),
		rightMostpointer:     []byte{0, 0, 0, 2},
		cellArea:             []byte{0, 0, 0, 2, 0, 2},
		startCellContentArea: PageSize - 6,
		isOverflow:           false,
		leftSibling:          nil,
		rightSiblisng:        nil,
	}

	server := ServerStruct{
		firstPage: zeroPage,
	}

	// var firstPage = PageParsed{
	// 	dbHeader:       DbHeader{},
	// 	dbHeaderSize:   0,
	// 	cellAreaParsed: [][]byte{[]byte{0, 0, 0, 1, 0, 0}, []byte{0, 0, 0, 2, 0, 0}, []byte{0, 0, 0, 3, 0, 0}, []byte{0, 0, 0, 4, 0, 0}},
	// 	isOverflow:     true,
	// 	leftSibling:    nil,
	// 	isLeaf:         true,
	// 	rightSiblisng:  nil,
	// }

	var firstPage = PageParsed{
		dbHeader:     DbHeader{},
		dbHeaderSize: 0,
		pageNumber:   1,
		btreeType:    int(TableBtreeLeafCell),
		cellAreaParsed: [][]byte{[]byte{0, 0, 0, 1, 0, 1}, []byte{0, 0, 0, 2, 0, 2}, []byte{0, 0, 0, 3, 0, 3}, []byte{0, 0, 0, 4, 0, 4}, []byte{0, 0, 0, 5, 0, 5}, []byte{0, 0, 0, 6, 0, 6}, []byte{0, 0, 0, 7, 0, 7},
			[]byte{0, 0, 0, 8, 0, 8}, []byte{0, 0, 0, 9, 0, 9}, []byte{0, 0, 0, 10, 0, 10}, []byte{0, 0, 0, 11, 0, 11}, []byte{0, 0, 0, 12, 0, 12}, []byte{0, 0, 0, 13, 0, 13}},
		startCellContentArea: PageSize - (13 * 6),
		isOverflow:           true,
		leftSibling:          nil,
		isLeaf:               true,
		rightSiblisng:        nil,
	}

	writer := NewWriter()
	reader := NewReader("fds")

	writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	balancingForNode(1, []int{0}, &server.firstPage)

	writer.flushPages("ds", &server.firstPage)

	firstPageRead := reader.readDbPage(0)
	firstPageParsed := parseReadPage(firstPageRead, 0)
	secondPageRead := reader.readDbPage(1)
	secondPageParsed := parseReadPage(secondPageRead, 1)
	thirdPageRead := reader.readDbPage(2)
	thirdPageParsed := parseReadPage(thirdPageRead, 2)
	fourthPageRead := reader.readDbPage(3)
	fourthPageParsed := parseReadPage(fourthPageRead, 3)
	fifthPageRead := reader.readDbPage(4)
	fifthPageParsed := parseReadPage(fifthPageRead, 4)
	sixthPageRead := reader.readDbPage(5)
	sixthPageParsed := parseReadPage(sixthPageRead, 5)
	seventhPageRead := reader.readDbPage(6)
	seventhPageParsed := parseReadPage(seventhPageRead, 6)
	eighthPageRead := reader.readDbPage(7)
	eighthPageParsed := parseReadPage(eighthPageRead, 7)
	fmt.Println("lets see pages")
	fmt.Printf("%+v \n", firstPageParsed)
	fmt.Println("```````````````Second page ````````````````````")
	fmt.Printf("%+v \n", secondPageParsed)
	fmt.Println("```````````````third page ````````````````````")
	fmt.Printf("%+v \n", thirdPageParsed)
	fmt.Println("```````````````fourth page ````````````````````")
	fmt.Printf("%+v \n", fourthPageParsed)
	fmt.Println("```````````````fifth page ````````````````````")
	fmt.Printf("%+v \n", fifthPageParsed)
	fmt.Println("```````````````fifth page ````````````````````")
	fmt.Printf("%+v \n", sixthPageParsed)
	fmt.Println("```````````````seventh page ````````````````````")
	fmt.Printf("%+v \n", seventhPageParsed)
	fmt.Println("```````````````eight page ````````````````````")
	fmt.Printf("%+v \n", eighthPageParsed)

	t.Errorf("test")

}

func TestBinarySearch11(t *testing.T) {
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
		cellAreaParsed:       [][]byte{[]byte{0, 0, 0, 4, 0, 4}, []byte{0, 0, 0, 3, 0, 3}, []byte{0, 0, 0, 2, 0, 2}, []byte{0, 0, 0, 1, 0, 1}},
		startCellContentArea: PageSize - (3 * 6),
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
		cellAreaParsed:       [][]byte{[]byte{0, 0, 0, 7, 0, 7}, []byte{0, 0, 0, 6, 0, 6}, []byte{0, 0, 0, 5, 0, 5}},
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
		cellAreaParsed:       [][]byte{[]byte{0, 0, 0, 12, 0, 12}, []byte{0, 0, 0, 11, 0, 11}, []byte{0, 0, 0, 10, 0, 10}, []byte{0, 0, 0, 9, 0, 9}, []byte{0, 0, 0, 8, 0, 8}},
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
	found, page, _ := search(0, 7)

	fmt.Println("lets see what we found")
	fmt.Println(found)
	fmt.Println(page)

	if !found {
		t.Errorf("expected to find result")
	}

	if page != 2 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 2, page)
	}
}

// func TestFdsfs(t *testing.T) {
// 	a, b, c := search(0, []int{}, 8)
// 	fmt.Println(a, b, c)

// 	t.Errorf("test")

// }

// func TestFdsfs1(t *testing.T) {
// 	insert(9)
// 	insert(10)

// 	t.Errorf("test")

// }

// func TestLeaftBiasDistribution(t *testing.T) {
// 	cellToDistribute := []Cell{{size: 1}, {size: 2}, {size: 3}, {size: 1}}

// 	totalSizeInEachPage, numberOfCellPerPage := leaf_bias(cellToDistribute)

// 	fmt.Println("totalSize")
// 	fmt.Println(totalSizeInEachPage)
// 	fmt.Println("numberofCellPerPage")
// 	fmt.Println(numberOfCellPerPage)

// 	totalSizeInEachPage, numberOfCellPerPage = accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute)

// 	fmt.Println("totalSize")
// 	fmt.Println(totalSizeInEachPage)
// 	fmt.Println("numberofCellPerPage")
// 	fmt.Println(numberOfCellPerPage)

// 	dividers, pages := redistribution(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute)

// 	fmt.Println("distribtuion, divders")
// 	fmt.Println(dividers)
// 	fmt.Println("pages")
// 	fmt.Println(pages)

// 	t.Errorf("test")

// }

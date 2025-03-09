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
		cellAreaParsed:       [][]byte{[]byte{0, 0, 0, 4, 0, 4}, []byte{0, 0, 0, 3, 0, 3}, []byte{0, 0, 0, 2, 0, 2}, []byte{0, 0, 0, 1, 0, 1}},
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
		t.Errorf("expected to find rowid 7")
	}

	if page != 2 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 2, page)
	}

	found, page, _ = search(0, 12)

	if !found {
		t.Errorf("expected to find result rowid 12")
	}

	if page != 3 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 3, page)
	}

	found, page, _ = search(0, 4)

	if !found {
		t.Errorf("expected to find result rowid 4")
	}

	if page != 1 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 1, page)
	}

	found, page, _ = search(0, 1)

	if !found {
		t.Errorf("expected to find result rowid 1")
	}

	if page != 1 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 1, page)
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
		cellAreaParsed:       [][]byte{[]byte{0, 0, 0, 4, 0, 4}, []byte{0, 0, 0, 3, 0, 3}, []byte{0, 0, 0, 2, 0, 2}, []byte{0, 0, 0, 1, 0, 1}},
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
	found, pageNumber, page := search(0, 13)

	if found {
		t.Errorf("Shouldn't find rowid 13 in any of pages")
	}

	if pageNumber != 3 {
		t.Errorf("expected to insert new value at page 3, so be return value, insted we got: %v", pageNumber)
	}
	if pageNumber != page.pageNumber {
		t.Error("Returned page number and page are differnet")
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
		cellAreaParsed:       [][]byte{[]byte{0, 0, 0, 4, 0, 4}, []byte{0, 0, 0, 3, 0, 3}, []byte{0, 0, 0, 2, 0, 2}, []byte{0, 0, 0, 1, 0, 1}},
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
	var fourthPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           4,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{0, 0, 0, 16, 0, 16}, []byte{0, 0, 0, 15, 0, 15}, []byte{0, 0, 0, 14, 0, 14}, []byte{0, 0, 0, 13, 0, 13}},
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
	writer.writeToFile(assembleDbPage(fifthPage), 5, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(sixthPage), 6, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(seventhPage), 7, "", &server.firstPage)
	found, page, _ := search(0, 7)

	fmt.Println("lets see what we found")
	fmt.Println(found)
	fmt.Println(page)

	if !found {
		t.Errorf("expected to find rowid 7")
	}

	if page != 2 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 2, page)
	}

	found, page, _ = search(0, 12)

	if !found {
		t.Errorf("expected to find result rowid 12")
	}

	if page != 3 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 3, page)
	}

	found, page, _ = search(0, 4)

	if !found {
		t.Errorf("expected to find result rowid 4")
	}

	if page != 1 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 1, page)
	}

	found, page, _ = search(0, 1)

	if !found {
		t.Errorf("expected to find result rowid 1")
	}

	if page != 1 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 1, page)
	}

	found, page, _ = search(0, 16)

	if !found {
		t.Errorf("expected to find result rowid 16")
	}

	if page != 4 {
		t.Errorf("expected to found result on page: %v, instead we got: %v", 4, page)
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
		cellAreaParsed:       [][]byte{[]byte{0, 0, 0, 4, 0, 4}, []byte{0, 0, 0, 3, 0, 3}, []byte{0, 0, 0, 2, 0, 2}, []byte{0, 0, 0, 1, 0, 1}},
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
	var fourthPage = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		pageNumber:           4,
		btreeType:            int(TableBtreeLeafCell),
		cellAreaParsed:       [][]byte{[]byte{0, 0, 0, 16, 0, 16}, []byte{0, 0, 0, 15, 0, 15}, []byte{0, 0, 0, 14, 0, 14}, []byte{0, 0, 0, 13, 0, 13}},
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
	writer.writeToFile(assembleDbPage(fifthPage), 5, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(sixthPage), 6, "", &server.firstPage)
	writer.writeToFile(assembleDbPage(seventhPage), 7, "", &server.firstPage)
	node := insert(17, 0, &server.firstPage)

	if node.pageNumber != 4 {
		t.Errorf("Insert values should be in page: %v, instead we got: %v", 4, node.pageNumber)
	}

	if len(node.cellAreaParsed) != 5 {
		t.Errorf("expected cell area to have 5 elements, instead we got: %v", len(node.cellAreaParsed))
	}

}

package main

import (
	"testing"
)

// debug this test
func TestLeaftBiasDistribution(t *testing.T) {
	clearDbFile("test")
	// var zeroPage = PageParsed{
	// 	dbHeader:         header(),
	// 	btreeType:        int(TableBtreeInteriorCell),
	// 	dbHeaderSize:     100,
	// 	cellAreaParsed:   [][]byte{},
	// 	rightMostpointer: []byte{0, 0, 0, 1},
	// 	isOverflow:       false,
	// 	leftSibling:      nil,
	// 	rightSiblisng:    nil,
	// }

	// server := ServerStruct{
	// 	firstPage: zeroPage,
	// }

	// // var firstPage = PageParsed{
	// // 	dbHeader:       DbHeader{},
	// // 	dbHeaderSize:   0,
	// // 	cellAreaParsed: [][]byte{[]byte{0, 0, 0, 1, 0, 0}, []byte{0, 0, 0, 2, 0, 0}, []byte{0, 0, 0, 3, 0, 0}, []byte{0, 0, 0, 4, 0, 0}},
	// // 	isOverflow:     true,
	// // 	leftSibling:    nil,
	// // 	isLeaf:         true,
	// // 	rightSiblisng:  nil,
	// // }

	// var firstPage = PageParsed{
	// 	dbHeader:       DbHeader{},
	// 	dbHeaderSize:   0,
	// 	btreeType:      int(TableBtreeLeafCell),
	// 	cellAreaParsed: [][]byte{[]byte{0, 0, 0, 0, 0, 1}, []byte{0, 0, 0, 0, 0, 2}, []byte{0, 0, 0, 0, 0, 3}, []byte{0, 0, 0, 0, 0, 4}},
	// 	isOverflow:     true,
	// 	leftSibling:    nil,
	// 	isLeaf:         false,
	// 	rightSiblisng:  nil,
	// }

	// writer := NewWriter()

	// writer.writeToFile(assembleDbPage(zeroPage), 0, "", &server.firstPage)
	// writer.writeToFile(assembleDbPage(firstPage), 1, "", &server.firstPage)
	balancingForNode(1, []int{0})

	t.Errorf("test")

}

// func TestFdsfs(t *testing.T) {
// 	a, b, c := search(0, []int{}, 8)
// 	fmt.Println(a, b, c)

// 	t.Errorf("test")

// }

func TestFdsfs1(t *testing.T) {
	insert(9)
	insert(10)

	t.Errorf("test")

}

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

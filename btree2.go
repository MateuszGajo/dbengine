package main

import (
	"fmt"
)

// 1. We need sturcutre to have children ranges, and if it index, or if it a cell?

// inserting and balancing??
// We have spae
// 1. We start at root if there is a space we add there

// We dont have space on root page
// 1. root doesnt have space, we overflow and create a new page, special case for root
// 2. move all the data to new page, it overflow but we can reuse balacning log same as for leaf node
// 2. move all data to memory to page not be in overflow state
// 3. Calculate how many pages we need?
// 3.1 left to right, so we can first three put in first page and fourth enrry goes to parent, dummer look, we leave one page empty

/// OVERFLOW state * assuming we can only fit 3 pages in every node
///                                     +--------+ Root Node
///                  					|1,2,3,4 |
///                                     +--------+

///                                     +--------+ Root Node
///                   +-----------------|   4   | --------------------
///                  /                  +--------+
///                 /                                             \
///            +-------+                                     +----------+ empty
///       +----|  1,2,3  |----+                 +------------| 			|------------+
///      /     +-------+     \                 /             +----------+             \
///     /          |          \               /                /      \                \

// 3.2 Second look it goes right to left, readjust the entires so the tree is balances, move 3 as root, move 4 to righT???, no we have leaf baias distribution

// leaf bias distribution

type Cell struct {
	size       int
	pageNumber int
	rowId      int
	data       []byte
}

// fix this size later
var usableSpacePerPage = PageSize - 600

func balance(pageNumber int, parent []PageParsed, node Node) {
	isRoot := len(parent) == 0

	// nothing to do here

	if !node.isOverflow {
		return
	}

	if isRoot && node.isOverflow {
		//Create a new page
		// move data to new page

		page := PageParsed{}
		parent = append(parent, page)
		pageNumber = 1
	}

}

// sibling redistribution

// WE are missing divider entry, so 3 pages gotes to left page, then fourth one is the root

// add new page 12
// start

///                                +-------+
///                        +-------| 4,8   |-------+                           X
///                      /        +-------+        \
///                      /             |             \
///   			  +----------+  +----------+  +----------+
///   			  | 1,2,3    |  | 5, 6, 7  |  | 9, 10, 11|
/// 			  +----------+  +----------+  +----------+

// phse one
///                                +-------+
///                        +-------| 4,8,12 |-------+                           X
///                      /        +-------+        \
///                      /             |             \
///   			  +----------+  +----------+  +----------+
///   			  | 1,2,3    |  | 5, 6, 7  |  | 9, 10, 11|
/// 			  +----------+  +----------+  +----------+

/// phase two create empty page
//|                               +-------+
///                        +-------| 4,8,12 |-------+    -------                        X
///                      /        +-------+        \			\
///                      /             |             \			\
///   			  +----------+  +----------+  +----------+     +---------+
///   			  | 1,2,3    |  | 5, 6, 7  |  | 9, 10, 11|     |         |
/// 			  +----------+  +----------+  +----------+     +---------+

/// phase  three: second distribution loop readjust entry
//|                               +-------+
///                        +-------| 4,8,12 |-------+    -------                        X
///                      /        +-------+        \			\
///                      /             |             \			\
///   			  +----------+  +----------+  +----------+     +---------+
///   			  | 1,2,3    |  | 5, 6, 7  |  | 9, 10, 11|     |         |
/// 			  +----------+  +----------+  +----------+     +---------+

/// end
//|                               +-------+
///                        +-------| 4,7,10 |-------+    -------                        X
///                      /        +-------+        \			\
///                      /             |             \			\
///   			  +----------+  +----------+  +----------+     +---------+
///   			  | 1,2,3    |  | 5,6      |  | 8,9      |     |  11,12   |
/// 			  +----------+  +----------+  +----------+     +---------+

// func getSiblings() ([][]Cell, []Cell) {
// 	// Get only like 1 left 1 right sibiling, or sometime 2 left sibling
// 	divider := []Cell{{size: 1, number: 4}} // 2 divider, betwen first and secopd siblign, betwen second and third
// 	siblings := [][]Cell{[]Cell{{size: 1, number: 1}, {size: 1, number: 2}, {size: 1, number: 3}}, []Cell{{size: 1, number: 5}, {size: 1, number: 6}, {size: 1, number: 7}, {size: 1, number: 8}}}
// 	usableSpacePerPage = 3
// 	return siblings, divider
// }

var PageNumber = 2
var loopIteration = 0

// we have working hardcoded implementation
// No we need to implement real pages that are read from disk
// look at the parsed cell area, maybe we can remove it or make it stick

func balancingForNode(node PageParsed, parents []int, firstPage *PageParsed) {

	reader := NewReader("")

	fmt.Println("node")

	siblings := []PageParsed{}

	fmt.Println("checkpoint1")

	if node.leftSibling != nil {
		pageRaw := reader.readDbPage(*node.leftSibling)
		page := parseReadPage(pageRaw, *node.leftSibling)
		siblings = append(siblings, page)
	}

	siblings = append(siblings, node)

	if node.rightSiblisng != nil {
		pageRaw := reader.readDbPage(*node.rightSiblisng)
		page := parseReadPage(pageRaw, *node.rightSiblisng)
		siblings = append(siblings, page)
	}

	isRoot := len(parents) == 0

	fmt.Println("checkpoint2")
	fmt.Println(node.isOverflow)

	if !node.isOverflow {
		return
	}

	if isRoot && node.isOverflow {
		fmt.Println("is root???")
		fmt.Println("is root???")
		fmt.Println("is root???")
		// take root page
		// insert a new page and move that taken from root and run balancing algoritm
		// root := getRootData()
		newPage := PageParsed{
			pageNumber:     firstPage.dbHeader.dbSizeInPages,
			cellArea:       node.cellArea,
			cellAreaParsed: node.cellAreaParsed,
			btreeType:      int(TableBtreeInteriorCell), startCellContentArea: PageSize,
		}
		// pageNumber is 0, we append new root page
		parents = append(parents, node.pageNumber)
		siblings = []PageParsed{newPage}
		node.pageNumber = 0
		firstPage.dbHeader.dbSizeInPages++
		node = newPage

		// siblings = [][]Cell{root}
		// divider = []Cell{}
	}
	fmt.Println("parents failed??")
	var parent int

	if len(parents) > 0 {

		parent = parents[len(parents)-1]
		parents = parents[:len(parents)-1]
	}

	cellToDistribute := []Cell{}
	startIndex := 0
	endIndex := 0
	fmt.Println("hello siblings data")
	fmt.Printf("%+v", siblings)
	for i, v := range siblings {
		for _, vN := range v.cellAreaParsed {
			var rowId int
			var pageNumber int
			if v.btreeType == int(TableBtreeInteriorCell) {
				rowIdUint64 := DecodeVarint(vN[4:6])
				rowId = int(rowIdUint64)
				pageNumberInt64 := DecodeVarint(vN[:4])
				pageNumber = int(pageNumberInt64)
			} else if v.btreeType == int(TableBtreeLeafCell) {
				if vN[0] > 127 || vN[1] > 127 {
					panic("implement this")
				}
				fmt.Println("hello vn")
				fmt.Println(vN)
				rowId = int(vN[1])
				pageNumber = v.pageNumber
			} else {
				panic("implement this")
			}
			fmt.Println("row id???", vN[4:6])

			// if node.isLeaf {
			// 	rowIdUint64, _ := DecodeVarint(vN[2:4])
			// 	rowId = int(rowIdUint64)
			// } else {
			// 	rowIdUint64, _ := DecodeVarint(vN[4:6])
			// 	rowId = int(rowIdUint64)
			// }

			cellToDistribute = append(cellToDistribute, Cell{size: len(vN), pageNumber: int(pageNumber), rowId: rowId, data: vN})
		}
		// divider should be taken from parens, by taken i mean removed

		if i < len(siblings)-1 {
			divider, newStartIndex, endIndex2 := getDivider(v.pageNumber)
			if startIndex == 0 {
				startIndex = newStartIndex
			}
			if i < len(siblings)-2 {
				endIndex = endIndex2
			}

			cellToDistribute = append(cellToDistribute, Cell{size: 1, pageNumber: divider.rowid})
			fmt.Println("after")
		}
	}

	fmt.Println("cells?")
	fmt.Println(cellToDistribute)

	totalSizeInEachPage, numberOfCellPerPage := leaf_bias(cellToDistribute, node)

	fmt.Println("leaf bias")
	fmt.Println(totalSizeInEachPage)
	fmt.Println(numberOfCellPerPage)

	totalSizeInEachPage, numberOfCellPerPage = accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute, node)

	fmt.Println("move to right")
	fmt.Println(totalSizeInEachPage)
	fmt.Println(numberOfCellPerPage)

	oldLastSibling := siblings[len(siblings)-1]

	// if len(numberOfCellPerPage) != len(siblings) {
	// 	// basically allocating new page
	// 	newPage := PageParsed{pageNumber: PageNumber}
	// 	PageNumber++
	// 	secondPage = newPage
	// 	siblings = append(siblings, secondPage)

	// }

	// fix pointers
	lastSibling := siblings[len(siblings)-1]

	lastSibling.rightSiblisng = oldLastSibling.rightSiblisng

	//somehow fix parent pointers

	// lastSibling.po

	fmt.Println("redistribution number of pages??")
	fmt.Println(numberOfCellPerPage)

	// add new pages
	for len(siblings) < len(numberOfCellPerPage) {
		btreeType := TableBtreeLeafCell
		if !node.isLeaf {
			btreeType = TableBtreeInteriorCell
		}
		siblings = append(siblings, PageParsed{pageNumber: firstPage.dbHeader.dbSizeInPages, btreeType: int(btreeType), startCellContentArea: PageSize})
		firstPage.dbHeader.dbSizeInPages++
	}

	// free pages
	for len(numberOfCellPerPage) < len(siblings) {
		siblings = siblings[:len(siblings)-1]
	}
	deivider, pages := redistribution(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute, len(siblings), node)
	fmt.Println("pages")
	fmt.Printf("%+v \n", pages)
	fmt.Println("siblings")
	fmt.Printf("%+v \n", siblings)
	writer := NewWriter()
	for i, v := range pages {
		rowData := []byte{}
		for _, vN := range v {
			// rowData = append(rowData, intToBinary(vN.pageNumber, 4)...)
			// rowData = append(rowData, intToBinary(vN.rowId, 2)...)
			// fmt.Println("what is in vn data??")
			// fmt.Println(vN.data)
			rowData = append(rowData, vN.data...)

		}
		siblings[i].cellArea = rowData
		siblings[i].cellAreaParsed = [][]byte{}
		siblings[i].startCellContentArea = PageSize - len(rowData)
		fmt.Println("save node", siblings[i].pageNumber)
		writer.writeToFile(assembleDbPage(siblings[i]), siblings[i].pageNumber, "", firstPage)
	}

	updateDivider(parent, deivider, startIndex, endIndex, firstPage)

	newPage := reader.readFromMemory(parent)

	balancingForNode(newPage, parents, firstPage)

}

func leaf_bias(cells []Cell, node PageParsed) ([]int, []int) {

	fmt.Println("cells")
	fmt.Println(cells)
	fmt.Println("cells length")
	fmt.Println(len(cells))

	totalSizeInEachPage := []int{0}
	numberOfCellPerPage := []int{0}

	for _, v := range cells {
		i := len(totalSizeInEachPage) - 1

		if totalSizeInEachPage[i]+v.size <= usableSpacePerPage {
			fmt.Println("how many time we enter here?")
			totalSizeInEachPage[i] += v.size
			numberOfCellPerPage[i]++
		} else if node.isLeaf {
			totalSizeInEachPage = append(totalSizeInEachPage, v.size)
			numberOfCellPerPage = append(numberOfCellPerPage, 1)

		} else {
			totalSizeInEachPage = append(totalSizeInEachPage, 0)
			numberOfCellPerPage = append(numberOfCellPerPage, 0)
		}
	}
	fmt.Println("hello result?")
	fmt.Println(totalSizeInEachPage)

	return totalSizeInEachPage, numberOfCellPerPage
}

// account for underflow towards the right

func accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage []int, cellToDistribute []Cell, node PageParsed) ([]int, []int) {

	// PAGE 1 [1,2,3], PageTwo [4]
	//4 -1 = 3 -1 =2
	fmt.Println("look at this")
	fmt.Println(len(cellToDistribute))
	fmt.Println(numberOfCellPerPage[len(numberOfCellPerPage)-1])
	fmt.Println("how many pages we got")
	fmt.Println(totalSizeInEachPage)
	fmt.Println("look at this")
	// 7 - 3 -1 =3
	divCell := len(cellToDistribute) - numberOfCellPerPage[len(numberOfCellPerPage)-1] - 1

	if len(numberOfCellPerPage) >= 2 {

		for i := len(totalSizeInEachPage) - 1; i > 0; i-- {
			fmt.Println("i", i)
			fmt.Println("start", totalSizeInEachPage[i], usableSpacePerPage/2)
			for float64(totalSizeInEachPage[i]) < float64(usableSpacePerPage)/2 {
				fmt.Println("enter heere??")
				fmt.Println(totalSizeInEachPage[i])
				totalSizeInEachPage[i] += cellToDistribute[divCell].size
				numberOfCellPerPage[i]++

				numberOfCellPerPage[i-1]--
				// divcell - 1 because of divider
				if node.isLeaf {
					totalSizeInEachPage[i-1] -= cellToDistribute[divCell].size
				} else {
					totalSizeInEachPage[i-1] -= cellToDistribute[divCell-1].size
				}

				divCell--
			}
		}
		// Second page has more data than the first one, make a little
		// adjustment to keep it left biased.

		if float64(totalSizeInEachPage[0]) < float64(usableSpacePerPage)/2 {
			fmt.Println("are we entering here?????")
			numberOfCellPerPage[0] += 1
			numberOfCellPerPage[1] -= 1
		}
	}

	return totalSizeInEachPage, numberOfCellPerPage
}

// TODO:
// leaf this redistrubtion
// work on saving

func redistribution(totalSizeInEachPage, numberOfCellPerPage []int, cellToDistribute []Cell, siblingsLength int, node PageParsed) ([]Cell, [][]Cell) {
	dividers := []Cell{}
	pages := make([][]Cell, len(numberOfCellPerPage))
	pageNumber := -1
	cellIndex := 0
	fmt.Println("cells")
	fmt.Println(cellToDistribute)
	for i, v := range numberOfCellPerPage {
		pageNumber++
		for range v {
			if cellIndex < len(cellToDistribute) {
				fmt.Println("before tratedy")
				pages[pageNumber] = append(pages[pageNumber], cellToDistribute[cellIndex])
				cellIndex++
				fmt.Println("after tratedy")
			}
		}

		if i < siblingsLength-1 {
			if cellIndex >= len(cellToDistribute) {
				panic("should never occur")
			}
			dividers = append(dividers, cellToDistribute[cellIndex])
			if !node.isLeaf {
				cellIndex++

			}
		}
	}

	return dividers, pages
}

// inserts a new key into key updates the value if keys already exists

// entries are always inserted at leaft nodes. Internal nodes and the root can only grow in size when leaft nodes overflow and siblings cant take any load to kee pthe leaves balanced, causing a split

// lets walk through an exmaple
// suppose we have the following tree of orde=r, which means each node hold at maxiumum 3 keys and 4 children

/// ```text
///                             PAGE 0 (ROOT)
///                              +-------+
///                          +---|  3,6  |---+
///                         /    +-------+    \
///                        /         |         \
///                   +-------+  +-------+  +-------+
///                   |  1,2  |  |  4,5  |  |  7,8  |
///                   +-------+  +-------+  +-------+
///                     PAGE 1     PAGE 2     PAGE 3
/// ```

// lets try insert key 9, the inserion alogirthm will call find to the page and index, when new key should be added, it will simply insert the key
/// ```text
///                             PAGE 0 (ROOT)
///                              +-------+
///                          +---|  3,6  |---+
///                         /    +-------+    \
///                        /         |         \
///                   +-------+  +-------+  +---------+
///                   |  1,2  |  |  4,5  |  |  7,8,9  |
///                   +-------+  +-------+  +---------+
///                     PAGE 1     PAGE 2     PAGE 3
/// ```

// page 3 has maxiumum number of key if we try to insert now key number 10

///
/// ```text
///                             PAGE 0 (ROOT)
///                              +-------+
///                          +---|  3,6  |---+
///                         /    +-------+    \
///                        /         |         \
///                   +-------+  +-------+  +----------+
///                   |  1,2  |  |  4,5  |  | 7,8,9,10 |
///                   +-------+  +-------+  +----------+
///                     PAGE 1     PAGE 2      PAGE 3
/// ```
/// at the end we need to call balcning algorithm

// Test insert function!!!

// cell
// insert a new value, where there is a rowid??

func (page *PageParsed) isSpace() bool {
	fmt.Println("are we overflowing?")
	fmt.Println(page.cellAreaSize + page.btreePageHeaderSize + len(page.pointers))
	fmt.Println((page.cellAreaSize + page.btreePageHeaderSize + len(page.pointers)) > PageSize)
	return (page.cellAreaSize + page.btreePageHeaderSize + len(page.pointers)) < PageSize
}

func (page *PageParsed) insertData(data CreateCell) {
	fmt.Println("run insert data func")
	cellParsedData := [][]byte{data.data}
	cellParsedData = append(cellParsedData, page.cellAreaParsed...)

	cellArea := data.data
	cellArea = append(cellArea, page.cellArea...)
	page.cellArea = cellArea
	page.cellAreaParsed = cellParsedData

	page.cellAreaSize += data.dataLength

	if !page.isSpace() {
		fmt.Println("set overflow to true1~!!!")
		page.isOverflow = true
	}

	page.numberofCells++

	page.startCellContentArea -= data.dataLength

	newPointers := intToBinary(page.startCellContentArea, 2)
	newPointers = append(newPointers, page.pointers...)

	page.pointers = newPointers
}

func (page *PageParsed) updateData(data CreateCell, index int) {
	dataDifferences := data.dataLength - len(page.cellAreaParsed[index])

	page.cellAreaParsed[index] = data.data

	page.cellAreaSize += dataDifferences

	newPointers := []byte{}
	start := PageSize

	for _, v := range page.cellAreaParsed {
		start -= len(v)
		newPointers = append(newPointers, intToBinary(start, 2)...)
	}

	if !page.isSpace() {
		fmt.Println("set overflow to true2~!!!")
		page.isOverflow = true
	}

	page.startCellContentArea -= PageSize - page.cellAreaSize

	page.pointers = newPointers
}

func insert(rowId int, cell CreateCell, firstPage *PageParsed) PageParsed {
	ok, index, node, parents := search(0, rowId, []int{})

	fmt.Println("are we here???")

	if ok {
		fmt.Println("are we here???1")
		node.updateData(cell, index)
		return node
	} else {
		// insert condition
		node.insertData(cell)
		fmt.Println("2")
	}

	writer := NewWriter()
	writer.softwiteToFile(node, node.pageNumber, firstPage)

	balancingForNode(node, parents, firstPage)

	return node
}

// Search algorithm
// 1. Read the subtree node into memory
// 2. Run a binary search on the entries to find the given key.
// 3. IF successful, return the result
// 4. If not, the binary search result will tell us which child to pick for

//  find key 9 in this tree, located at page 5

///
/// ```text
///                             PAGE 0
///                           +--------+
///                   +-------|   11   |-------+
///                  /        +--------+        \
///                 /                            \
///            +-------+ PAGE 1              +--------+
///       +----|  4,8  |----+                |   14   |
///      /     +-------+     \               +--------+
///     /          |          \               /      \
/// +-------+  +-------+  +-------+     +-------+  +-------+
/// | 1,2,3 |  | 5,6,7 |  | 9,10  |     | 12,13 |  | 15,16 |
/// +-------+  +-------+  +-------+     +-------+  +-------+
///  PAGE 3     PAGE 4      PAGE 5

// first iteration
// 1. read page 0 into memory
// 2. binary search on page result in err(not found)
// 3. read index 0 using page.child and recurse into the result

// second iteraion
// 1. read page 1 into memroy
// 2. binary search result in err
// 3. read child pointer at index 3, and recurse again

// final iteration
// 1. read page 5 into memory
// 2. binary search reuslt it ok
// 3. done, return result

type Index struct {
	id              int
	leftPointer     int
	leftPointerNull bool
}

type Node struct {
	indexes      []Index
	leaf         bool
	isOverflow   bool
	rightPointer int
	pageNumber   int
}

func search(pageNumber int, entry int, parents []int) (bool, int, PageParsed, []int) {
	reader := NewReader("")
	page := reader.readFromMemory(pageNumber)

	ok, newPageNumber, index := binary_search(page, pageNumber, entry)
	fmt.Println("Search iteration, page number")
	fmt.Println(pageNumber)

	if !ok && page.isLeaf {
		fmt.Println("didnt find what we are looking for")
		return false, index, page, parents
	}
	if ok {
		fmt.Println("enter ok???")
		return ok, index, page, parents
	}
	parents = append(parents, pageNumber)
	return search(newPageNumber, entry, parents)
}

// WE add new entry?
// we add it to last height to 0x0d tree
// then we need to balance tree
// let say pointer has change, how we approach it?
// we redistribut load as we start
///                                +-------+
///                        +-------| 4,8,12 |-------+                           X
///                      /        +-------+        \
///                      /             |             \
///   			  +----------+  +----------+  +----------+
///   			  | 1,2,3    |  | 5, 6, 7  |  | 9, 10, 11|
/// 			  +----------+  +----------+  +----------+

///                        +-------| 4,7,10, .... |-------+    -------            -------            X
///                      /        +----------------+ \			\					\
///                      /             |             \			\					\
///   			  +----------+  +----------+  +----------+     +---------+			........
///   			  | 1,2,3    |  | 5,6      |  | 8,9      |     |  11,12   |
/// 			  +----------+  +----------+  +----------+     +---------+
// last height are leaft page that stores only value, so on the 0x05 iterior page we change from 4,8,12 to 4,7, 12
// WE need to have somehow index, or space or something wheer to replace this pointers
// we should remove all three 4,8,12 annd the paste 4,7,12, (could be more items)
// These pointer are basically cellArea, but we do have a index, so we can use it for now
// separate inserting from updating pointers
// pointers requires rowId, these value we move around should have rowid
//

func binary_search(page PageParsed, ageNumber int, rowIdAsEntry int) (bool, int, int) {
	// node := getNode2(pageNumber)
	// node := PageParsed{}
	fmt.Println("hello page number")
	fmt.Println(page.pageNumber)

	if len(page.cellAreaParsed) == 0 {
		return false, PageNumber, 0
	}

	rightRowId := int(DecodeVarint(page.cellAreaParsed[0][4:6]))

	fmt.Println(rowIdAsEntry)
	rightPointerPage := int(DecodeVarint(page.rightMostpointer))
	if rowIdAsEntry > rightRowId {
		return false, rightPointerPage, 0
	}

	var leftRowId int

	for i := 0; i < len(page.cellAreaParsed); i++ {
		if i < len(page.cellAreaParsed)-1 {
			leftRowId = int(DecodeVarint(page.cellAreaParsed[i+1][4:6]))
		}
		if i > 0 {
			rightRowId = int(DecodeVarint(page.cellAreaParsed[i-1][4:6]))
		}

		v := page.cellAreaParsed[i]
		// TODO implement for 0x0d
		rowIdUint64 := DecodeVarint(v[4:6])
		rowId := int(rowIdUint64)
		pageNumberInt64 := DecodeVarint(v[:4])
		pageNumber := int(pageNumberInt64)

		if page.isLeaf && rowId == rowIdAsEntry {
			fmt.Println("find page in leaft page")
			return true, pageNumber, i
		}

		// we are on last page
		if i == len(page.cellAreaParsed)-1 && rowIdAsEntry <= rowId {
			fmt.Println("con 1")
			// go to current page (as this is last page)
			return false, pageNumber, i
		} else if i == len(page.cellAreaParsed)-1 {
			panic("should never happen")
		}

		// entry row id is smaller than current row id but grater than elft one, we need to go to the page
		if rowIdAsEntry <= rowId && rowIdAsEntry > leftRowId {
			fmt.Println("con 2")
			// go to current page
			return false, pageNumber, i

		}

	}

	// for i := 0; i < len(node.cellArea); i++ {
	// 	if entry == node.indexes[i].id {
	// 		return true, node.indexes[i].id, node
	// 	} else if entry < node.indexes[i].id {
	// 		return false, node.indexes[i].leftPointer, node
	// 	}
	// }
	panic("should never occur")
}

// binary searhc???, just basically iterate over celll

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
var usableSpacePerPage = PageSize - 12 //

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

var loopIteration = 0

// we have working hardcoded implementation
// No we need to implement real pages that are read from disk
// look at the parsed cell area, maybe we can remove it or make it stick

func (parentPage PageParsed) findSiblings(currentNode PageParsed) (*PageParsed, *PageParsed) {

	if len(parentPage.cellAreaParsed) == 0 && parentPage.pageNumber == 0 {
		return nil, nil
	}
	if len(currentNode.cellAreaParsed) == 0 {
		panic("should never occur, cell area parsed empty")
	}
	currnetPageLatestRow := parseCellArea(currentNode.cellAreaParsed[0], BtreeType(currentNode.btreeType))
	rowIdAsEntry := currnetPageLatestRow.rowId
	pageNumber := currentNode.pageNumber
	fmt.Println("row as entry??")
	fmt.Println(rowIdAsEntry)
	fmt.Println("page number?")
	fmt.Println(currentNode.pageNumber)
	fmt.Println("cell area?")
	fmt.Println(parentPage.cellArea)
	if len(parentPage.cellAreaParsed) == 0 {
		fmt.Println("parent page??")
		fmt.Println(parentPage.pageNumber)
		fmt.Println(parentPage.cellArea)
		panic("should never occur, cell area parsed empty")
	}
	reader := NewReader("")
	cellParsed := parseCellArea(parentPage.cellAreaParsed[0], BtreeType(parentPage.btreeType))

	rightRowId := cellParsed.rowId

	var leftPageNumber int
	var rightPageNumber int

	if rowIdAsEntry > rightRowId {
		fmt.Println("parent???")
		fmt.Println(parentPage.cellAreaParsed)
		fmt.Println("data")
		fmt.Println(rowIdAsEntry)
		fmt.Println(rightRowId)
		panic("should never happend, rowIdASEntry > rightRowId in find sibling")
	}
	// [0 0 0 2 0 2
	//  0 0 0 1 0 1]
	for i := 0; i < len(parentPage.cellAreaParsed); i++ {
		if i < len(parentPage.cellAreaParsed)-1 {
			cellParsed := parseCellArea(parentPage.cellAreaParsed[i+1], BtreeType(parentPage.btreeType))
			leftPageNumber = cellParsed.pageNumber
		}
		if i > 0 {
			cellParsed := parseCellArea(parentPage.cellAreaParsed[i-1], BtreeType(parentPage.btreeType))
			rightRowId = cellParsed.rowId
			rightPageNumber = cellParsed.pageNumber
		}

		v := parentPage.cellAreaParsed[i]

		cellParsed := parseCellArea(v, BtreeType(parentPage.btreeType))
		rowId := cellParsed.rowId

		if rowIdAsEntry == rowId && cellParsed.pageNumber == pageNumber {
			var leftSiblings *PageParsed
			var rightSiblings *PageParsed
			if i < len(parentPage.cellAreaParsed)-1 {

				page := reader.readFromMemory(leftPageNumber)
				leftSiblings = &page
			}
			if i > 0 {

				page := reader.readFromMemory(rightPageNumber)
				rightSiblings = &page
			}
			// we're on the items we looking siblings for
			return leftSiblings, rightSiblings
		} else if rowIdAsEntry == rowId {
			panic("terrible wrong, we can't have same rowid on diffrent pages")
		}

	}
	panic("should never enter here")
}

// how to update parent references??
// right most pointer:
// cell area

// rowids: 0-3, 4-6, 8-9
// parent: page 1, 0,0,0,1,0,3, 0,0,0,2,0,6, 0,0,0,3,0,9, 0,0,0,4,0, 12, 0,0,0,0,5,15, 0,0,0,0,6,0,18 0,0,0,7,0,22
// right most pointer 0,0,0,8

// we wanna take siblings of page 7, so its 6, and 8
// alright but now how to update pointer??
// what can happen, we can delete page, add new one
// lets say we added new page 0,0,0,9
// how now update parent?
// we have cells with pages:7,8,9
// i think we can use rowid, so all we need is start index
// or maybe lets hack it and add right page to content

//lets code it diffrently, if it need to take right most pointer, we need to load page and then take pointers

func balancingForNode(node PageParsed, parents []*PageParsed, firstPage *PageParsed) {

	siblings := []PageParsed{}

	// how to parse siblings???

	isRoot := len(parents) == 0

	if !node.isOverflow {
		return
	}
	var parent PageParsed

	writer := NewWriter()

	if isRoot && node.isOverflow {
		//fix this
		fmt.Println("is root???")
		fmt.Println("is root???")
		fmt.Println(node.btreeType)
		fmt.Println("is root???")
		// take root page
		// insert a new page and move that taken from root and run balancing algoritm
		// root := getRootData()
		cellArea := []byte{}
		// cell := parseCellArea(node.cellAreaParsed[0], BtreeType(node.btreeType))
		// cellArea = append(cellArea, 0)
		// cellArea = append(cellArea, byte(cell.rowId))
		rootPage := PageParsed{
			pageNumber:           0,
			dbHeaderSize:         100,
			dbHeader:             node.dbHeader,
			cellArea:             cellArea,
			cellAreaParsed:       [][]byte{},
			numberofCells:        0,
			cellAreaSize:         len(cellArea),
			startCellContentArea: PageSize - len(cellArea),
			btreeType:            int(TableBtreeInteriorCell),
			rightMostpointer:     []byte{},
		}

		node.dbHeaderSize = 0
		node.dbHeader = DbHeader{}
		// pageNumber is 0, we append new root page
		parents = append(parents, &rootPage)
		// siblings = []PageParsed{newPage}
		node.pageNumber = firstPage.dbHeader.dbSizeInPages

		// node.btreeType = int(TableBtreeLeafCell)
		firstPage.dbHeader.dbSizeInPages++

		fmt.Println("save pages??")
		fmt.Println(node.pageNumber)
		fmt.Println(rootPage.pageNumber)
		fmt.Println("save pages??")

		writer.softwiteToFile(&node, node.pageNumber, firstPage)
		writer.softwiteToFile(&rootPage, rootPage.pageNumber, firstPage)
	}

	if len(parents) > 0 {

		parent = *parents[len(parents)-1]
		parents = parents[:len(parents)-1]
	}

	leftSibling, rightSibling := parent.findSiblings(node)

	if leftSibling != nil {
		fmt.Println("left sibling??", leftSibling.pageNumber)
		// fmt.Println("left sibling??")
		siblings = append(siblings, *leftSibling)
	}

	siblings = append(siblings, node)

	if rightSibling != nil {
		fmt.Println("righyt sibling??")
		fmt.Println(rightSibling.pageNumber)
		siblings = append(siblings, *rightSibling)
	}

	cellToDistribute := []Cell{}
	startIndex := PageSize
	endIndex := 0
	fmt.Println("show me siblings??")
	fmt.Printf("\n%+v\n", siblings)
	for i, v := range siblings {
		// need to start from last, because first item its saved at the end of page
		for j, _ := range v.cellAreaParsed {
			index := len(v.cellAreaParsed) - 1 - j
			var rowId int
			var pageNumber int
			if v.btreeType == int(TableBtreeInteriorCell) {
				rowIdUint64 := DecodeVarint(v.cellAreaParsed[index][4:6])
				rowId = int(rowIdUint64)
				pageNumberInt64 := DecodeVarint(v.cellAreaParsed[index][:4])
				pageNumber = int(pageNumberInt64)
			} else if v.btreeType == int(TableBtreeLeafCell) {
				if len(v.cellAreaParsed[index]) < 2 {
					panic("should have at least 2 bytes")
				}
				if v.cellAreaParsed[index][0] > 127 || v.cellAreaParsed[index][1] > 127 {
					panic("implement this")
				}
				rowId = int(v.cellAreaParsed[index][1])
				pageNumber = v.pageNumber
			} else {
				panic("implement this")
			}

			// size: len(vN) + 2, 2 for pointer
			cellToDistribute = append(cellToDistribute, Cell{size: len(v.cellAreaParsed[index]), pageNumber: int(pageNumber), rowId: rowId, data: v.cellAreaParsed[index]})
		}
		// divider should be taken from parens, by taken i mean removed

		// how to get dividers,
		// lets say we have root page 0x05, with right pointer to 0001 page, and cell area 0002(01-rowid)
		// we are distribution sibling pages: 1, 2
		// getting dividers for these two 0x0d pages?
		// divider is in parent page
		// get divider by page  number??

		if i < len(siblings) {
			_, newStartIndex, endIndex2 := parent.getDivider(v.pageNumber)
			if startIndex > newStartIndex {
				startIndex = newStartIndex
			}
			if endIndex2 > endIndex {
				endIndex = endIndex2
			}

		}
	}

	fmt.Println("cell to distribute??")
	fmt.Println(cellToDistribute)

	totalSizeInEachPage, numberOfCellPerPage := leaf_bias(cellToDistribute, node)

	totalSizeInEachPage, numberOfCellPerPage = accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute, node)

	fmt.Println("number of cell per page!!")
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
	deivider, pages := redistribution(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute, siblings, node)

	for i, v := range pages {
		rowData := []byte{}
		//again start from the last
		for j, _ := range v {
			index := len(v) - 1 - j
			// rowData = append(rowData, intToBinary(vN.pageNumber, 4)...)
			// rowData = append(rowData, intToBinary(vN.rowId, 2)...)
			// fmt.Println("what is in vn data??")
			// fmt.Println(vN.data)
			rowData = append(rowData, v[index].data...)

		}
		siblings[i].cellArea = rowData
		siblings[i].cellAreaParsed = [][]byte{}
		siblings[i].startCellContentArea = PageSize - len(rowData)
		siblings[i].numberofCells = len(v)

		fmt.Println("hello assemble page, save node", siblings[i].pageNumber)

		writer.writeToFile(assembleDbPage(siblings[i]), siblings[i].pageNumber, "", firstPage)
	}

	modifyDivider(&parent, deivider, startIndex, endIndex, firstPage, parents)

	// newPage := reader.readFromMemory(parent)

	balancingForNode(parent, parents, firstPage)

}

func leaf_bias(cells []Cell, node PageParsed) ([]int, []int) {

	totalSizeInEachPage := []int{0}
	numberOfCellPerPage := []int{0}

	for _, v := range cells {
		i := len(totalSizeInEachPage) - 1

		if totalSizeInEachPage[i]+v.size <= usableSpacePerPage {
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

	return totalSizeInEachPage, numberOfCellPerPage
}

// account for underflow towards the right

func accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage []int, cellToDistribute []Cell, node PageParsed) ([]int, []int) {

	// PAGE 1 [1,2,3], PageTwo [4]
	//4 -1 = 3 -1 =2
	// 7 - 3 -1 =3
	divCell := len(cellToDistribute) - numberOfCellPerPage[len(numberOfCellPerPage)-1] - 1

	if len(numberOfCellPerPage) >= 2 {

		for i := len(totalSizeInEachPage) - 1; i > 0; i-- {
			for float64(totalSizeInEachPage[i]) < float64(usableSpacePerPage)/2 {
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
			numberOfCellPerPage[0] += 1
			numberOfCellPerPage[1] -= 1
		}
	}

	return totalSizeInEachPage, numberOfCellPerPage
}

// TODO:
// leaf this redistrubtion
// work on saving

func redistribution(totalSizeInEachPage, numberOfCellPerPage []int, cellToDistribute []Cell, siblingsLength []PageParsed, node PageParsed) ([]Cell, [][]Cell) {
	dividers := []Cell{}
	pages := make([][]Cell, len(numberOfCellPerPage))
	pageIndex := -1
	cellIndex := -1
	for i, v := range numberOfCellPerPage {
		pageIndex++
		for range v {
			if cellIndex < len(cellToDistribute) {
				cellIndex++
				pages[pageIndex] = append(pages[pageIndex], cellToDistribute[cellIndex])

			}
		}

		if i < len(siblingsLength) {

			if cellIndex >= len(cellToDistribute) {
				panic("should never occur, cell index >= than distribute")
			}
			dividers = append(dividers, Cell{
				rowId:      cellToDistribute[cellIndex].rowId,
				pageNumber: siblingsLength[i].pageNumber,
				data:       cellToDistribute[cellIndex].data,
				size:       cellToDistribute[cellIndex].size,
			})

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
	return (page.cellAreaSize + page.btreePageHeaderSize + page.dbHeaderSize + len(page.pointers)) < PageSize
}

func (page *PageParsed) insertData(data CreateCell, firstPage *PageParsed, parents []*PageParsed) {
	if len(parents) == 0 {
		panic("there shuld be always parent, insert data")
	}
	parent := parents[len(parents)-1]
	parents = parents[:len(parents)-1]
	divider, startIndex, endIndex := parent.getDivider(page.pageNumber)
	fmt.Println("parent divider indexses??", startIndex, endIndex)
	if divider.rowId == data.rowId {
		panic("that shouldn't never happen, we can't insert exisiting id")
	}
	if divider.rowId < data.rowId {
		cell := Cell{
			pageNumber: page.pageNumber,
			rowId:      data.rowId,
		}
		fmt.Println("modify divider in insert datta???")

		modifyDivider(parent, []Cell{cell}, startIndex, endIndex, firstPage, parents)
	} else {
		panic("in incremental write, should never happen")
	}
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

	newPointers := page.pointers
	newPointers = append(newPointers, intToBinary(page.startCellContentArea, 2)...)

	page.pointers = newPointers
}

func (page *PageParsed) updateData(data CreateCell, index int) {
	dataDifferences := data.dataLength - len(page.cellAreaParsed[index])

	page.cellAreaParsed[index] = data.data

	cellArea := []byte{}

	for _, v := range page.cellAreaParsed {
		cellArea = append(cellArea, v...)
	}

	page.cellAreaSize += dataDifferences
	page.cellArea = cellArea

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

	page.startCellContentArea = PageSize - page.cellAreaSize

	page.pointers = newPointers
}

// two type of inserts
// assuming page number 2, and current row id 20 insert at the end, so if paret has 0,0,0,2,0,20, should be update to 0,0,0,0,1,0,10 0,0,0,2,0,21
// insert where the row been deleted so if paret has 0,0,0,0,1,0,10 0,0,0,2,0,20 it should stay0,0,0,0,1,0,10 0,0,0,2,0,20

func insert(rowId int, cell CreateCell, firstPage *PageParsed) PageParsed {
	ok, index, node, parents := search(0, rowId, []*PageParsed{})

	fmt.Println("insert return node with page number", node.pageNumber)

	if ok {
		node.updateData(cell, index)
		return node
	} else {
		// insert condition
		node.insertData(cell, firstPage, parents)
	}

	writer := NewWriter()

	// if node.isOverflow {
	// 	fmt.Println("overflow for rowid?", rowId)
	// }

	writer.softwiteToFile(&node, node.pageNumber, firstPage)

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

func search(pageNumber int, entry int, parents []*PageParsed) (bool, int, PageParsed, []*PageParsed) {
	reader := NewReader("")
	page := reader.readFromMemory(pageNumber)
	ok, newPageNumber, index := binarySearch(page, pageNumber, entry)

	if !ok && page.isLeaf {
		return false, index, page, parents
	}
	if ok {
		return ok, index, page, parents
	}
	parents = append(parents, &page)
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

func parseCellArea(data []byte, btreeType BtreeType) Cell {
	if len(data) == 0 {
		panic("data is empty")
	}
	var rowId int
	var pageNumber int
	if btreeType == TableBtreeInteriorCell {
		//TODO:  i think we need to fix it
		pageNumberInt64 := DecodeVarint(data[:4])
		pageNumber = int(pageNumberInt64)
		rowId = int(DecodeVarint(data[4:6]))
	} else if btreeType == TableBtreeLeafCell {
		if data[0] > 127 || data[1] > 127 {
			panic("implement this")
		}
		rowId = int(data[1])
	} else {
		panic("tree don't implemented")
	}
	return Cell{
		rowId:      rowId,
		size:       len(data),
		data:       data,
		pageNumber: pageNumber,
	}
}

// TODO test binary search and search

func binarySearch(page PageParsed, pageNumber int, rowIdAsEntry int) (bool, int, int) {

	fmt.Println("binary search !!, page number, row id", pageNumber, rowIdAsEntry)

	if len(page.cellAreaParsed) == 0 && page.btreeType == int(TableBtreeLeafCell) {
		return false, pageNumber, 0
	} else if len(page.cellAreaParsed) == 0 {
		fmt.Printf("%+v", page)
		panic("should never occur, i guess")
	}
	cellParsed := parseCellArea(page.cellAreaParsed[0], BtreeType(page.btreeType))

	rightRowId := cellParsed.rowId

	rightPointerPage := int(DecodeVarint(page.rightMostpointer))
	if page.btreeType == int(TableBtreeInteriorCell) {
		if rowIdAsEntry > rightRowId && rightPointerPage != 0 {
			return false, rightPointerPage, 0
		}
	} else if page.btreeType == int(TableBtreeLeafCell) {
		if rowIdAsEntry > rightRowId {
			return false, page.pageNumber, 0
		}
	}

	var leftRowId int

	for i := 0; i < len(page.cellAreaParsed); i++ {
		if i < len(page.cellAreaParsed)-1 {
			cellParsed := parseCellArea(page.cellAreaParsed[i+1], BtreeType(page.btreeType))
			leftRowId = cellParsed.rowId
		}
		if i > 0 {
			cellParsed := parseCellArea(page.cellAreaParsed[i-1], BtreeType(page.btreeType))
			rightRowId = cellParsed.rowId
		}

		v := page.cellAreaParsed[i]
		// TODO implement for 0x0d
		cellParsed := parseCellArea(v, BtreeType(page.btreeType))
		rowId := cellParsed.rowId

		if page.isLeaf && rowId == rowIdAsEntry {
			return true, page.pageNumber, i
		}

		// we are on last page
		if i == len(page.cellAreaParsed)-1 && rowIdAsEntry <= rowId {
			// go to current page (as this is last page)
			return false, cellParsed.pageNumber, i
		} else if i == len(page.cellAreaParsed)-1 {
			panic("should never happen, binary search")
		}

		// entry row id is smaller than current row id but grater than elft one, we need to go to the page
		if rowIdAsEntry <= rowId && rowIdAsEntry > leftRowId {
			// go to current page
			return false, cellParsed.pageNumber, i

		}

	}

	// for i := 0; i < len(node.cellArea); i++ {
	// 	if entry == node.indexes[i].id {
	// 		return true, node.indexes[i].id, node
	// 	} else if entry < node.indexes[i].id {
	// 		return false, node.indexes[i].leftPointer, node
	// 	}
	// }
	panic("should never occur, binary search")
}

// binary searhc???, just basically iterate over celll

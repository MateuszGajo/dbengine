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
		panic("should never occur, cell area parsed empty, find siblings 1")
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
		// return nil, nil
		panic("should never occur, cell area parsed empty, find siblings 2")
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

	fmt.Println("is overflow??", node.isOverflow)

	if !node.isOverflow {
		return
	}

	var parent PageParsed

	writer := NewWriter()

	if isRoot && node.isOverflow {

		// insert a new page and move that taken from root and run balancing algoritm
		// root := getRootData()
		cellArea := []byte{}
		cell := parseCellArea(node.cellAreaParsed[0], BtreeType(node.btreeType))
		cellArea = append(cellArea, intToBinary(firstPage.dbHeader.dbSizeInPages, 4)...)
		cellArea = append(cellArea, intToBinary(cell.rowId, 2)...)
		parsedCellArea := append([][]byte{}, cellArea)
		var header *DbHeader
		if node.pageNumber == 0 {
			header = &node.dbHeader
		}

		rootPage := CreateNewPage(TableBtreeInteriorCell, parsedCellArea, node.pageNumber, header)

		node.dbHeaderSize = 0
		node.dbHeader = DbHeader{}
		parents = append(parents, &rootPage)
		node.pageNumber = firstPage.dbHeader.dbSizeInPages

		firstPage.dbHeader.dbSizeInPages++

		writer.softwiteToFile(&node, node.pageNumber, firstPage)
		writer.softwiteToFile(&rootPage, rootPage.pageNumber, firstPage)
	}
	usableSpacePerPage = PageSize - node.btreePageHeaderSize - node.dbHeaderSize - len(node.pointers)

	if len(parents) > 0 {

		parent = *parents[len(parents)-1]
		parents = parents[:len(parents)-1]
	}

	leftSibling, rightSibling := parent.findSiblings(node)

	if leftSibling != nil {
		siblings = append(siblings, *leftSibling)
	}

	siblings = append(siblings, node)

	if rightSibling != nil {
		siblings = append(siblings, *rightSibling)
	}

	cellToDistribute := []Cell{}
	startIndex := PageSize
	endIndex := 0

	for i, v := range siblings {
		// need to start from last, because first item its saved at the end of page
		for j, _ := range v.cellAreaParsed {
			index := len(v.cellAreaParsed) - 1 - j
			cell := parseCellArea(v.cellAreaParsed[index], BtreeType(v.btreeType))

			cellToDistribute = append(cellToDistribute, Cell{size: len(v.cellAreaParsed[index]), pageNumber: cell.pageNumber, rowId: cell.rowId, data: v.cellAreaParsed[index]})
		}

		if i < len(siblings) {
			_, newStartIndex, newEndIndex := parent.getDivider(v.pageNumber)
			if startIndex > newStartIndex {
				startIndex = newStartIndex
			}
			if newEndIndex > endIndex {
				endIndex = newEndIndex
			}

		}
	}

	totalSizeInEachPage, numberOfCellPerPage := leaf_bias(cellToDistribute)

	totalSizeInEachPage, numberOfCellPerPage = accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute, node)

	for len(siblings) < len(numberOfCellPerPage) {
		newPage := CreateNewPage(BtreeType(node.btreeType), [][]byte{}, firstPage.dbHeader.assignNewPage(), nil)
		siblings = append(siblings, newPage)
	}
	for len(numberOfCellPerPage) < len(siblings) {
		siblings = siblings[:len(siblings)-1]
	}
	deivider, pages := redistribution(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute, siblings, node)

	for i, v := range pages {
		rowData := []byte{}
		//again start from the last
		for j, _ := range v {
			index := len(v) - 1 - j

			rowData = append(rowData, v[index].data...)

		}
		siblings[i].updateCells(dbReadparseCellArea(byte(siblings[i].btreeType), rowData))

		writer.writeToFile(assembleDbPage(siblings[i]), siblings[i].pageNumber, "", firstPage)
	}

	modifyDivider(&parent, deivider, startIndex, endIndex, firstPage, parents)

	// newPage := reader.readFromMemory(parent)

	balancingForNode(parent, parents, firstPage)

}

func leaf_bias(cells []Cell) ([]int, []int) {

	totalSizeInEachPage := []int{0}
	numberOfCellPerPage := []int{0}

	for _, v := range cells {
		i := len(totalSizeInEachPage) - 1

		if totalSizeInEachPage[i]+v.size <= usableSpacePerPage {
			totalSizeInEachPage[i] += v.size
			numberOfCellPerPage[i]++
		} else {

			totalSizeInEachPage = append(totalSizeInEachPage, v.size)
			numberOfCellPerPage = append(numberOfCellPerPage, 1)
		}

	}

	return totalSizeInEachPage, numberOfCellPerPage
}

// account for underflow towards the right

func accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage []int, cellToDistribute []Cell, node PageParsed) ([]int, []int) {

	// PAGE 1 [1,2,3], PageTwo [4]
	//4 -1 = 3 -1 =2\
	fmt.Println("usable space??", usableSpacePerPage/2)
	// 7 - 3 -1 =3
	divCell := len(cellToDistribute) - numberOfCellPerPage[len(numberOfCellPerPage)-1] - 1

	if len(numberOfCellPerPage) >= 2 {

		for i := len(totalSizeInEachPage) - 1; i > 0; i-- {
			for totalSizeInEachPage[i] <= ((usableSpacePerPage / 2) - cellToDistribute[0].size) {
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

		if totalSizeInEachPage[0] < usableSpacePerPage/2 {
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

// lets see this one
func (page *PageParsed) insertData(data CreateCell, firstPage *PageParsed, parents []*PageParsed) {

	if len(parents) > 0 {
		parent := parents[len(parents)-1]
		parents = parents[:len(parents)-1]
		divider, startIndex, endIndex := parent.getDivider(page.pageNumber)
		if divider.rowId == data.rowId {
			panic("that shouldn't never happen, we can't insert exisiting id")
		}
		if divider.rowId < data.rowId && parent != nil {
			cell := Cell{
				pageNumber: page.pageNumber,
				rowId:      data.rowId,
			}
			fmt.Println("modify divider in insert datta???")

			modifyDivider(parent, []Cell{cell}, startIndex, endIndex, firstPage, parents)
		} else {
			panic("in incremental write, should never happen")
		}
	}

	fmt.Println("pointer 3")
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

func insert(rowId int, cell CreateCell, firstPage *PageParsed, startPageNumber *int) PageParsed {
	if startPageNumber != nil {
		fmt.Printf("insert, start page number?? :%d \n", *startPageNumber)
		fmt.Printf("insert, start page number?? :%d \n", *startPageNumber)
		fmt.Printf("insert, start page number?? :%d \n", *startPageNumber)
	}

	start := 0
	if startPageNumber != nil {
		start = *startPageNumber
	}
	fmt.Println("before search?")
	ok, index, node, parents := search(start, rowId, []*PageParsed{})

	if ok {
		fmt.Println("update")
		node.updateData(cell, index)
		return node
	} else {
		// insert condition
		fmt.Println("insert")
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
	fmt.Println("result,", ok, newPageNumber, index, page.isLeaf)

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

	fmt.Println("binary search !!, page number, row id", pageNumber, rowIdAsEntry, page.btreeType)

	if len(page.cellAreaParsed) == 0 && page.btreeType == int(TableBtreeLeafCell) {
		fmt.Println("pointer 0.1")
		return false, pageNumber, 0
	} else if len(page.cellAreaParsed) == 0 {
		panic("should never occur, i guess")
	}
	fmt.Println("pointer 1")
	cellParsed := parseCellArea(page.cellAreaParsed[0], BtreeType(page.btreeType))
	fmt.Println("pointer 1.5")

	rightRowId := cellParsed.rowId
	fmt.Println("right rowid??", rightRowId)

	rightPointerPage := int(DecodeVarint(page.rightMostpointer))
	fmt.Println("rightPointerPage?", rightPointerPage)
	if page.btreeType == int(TableBtreeInteriorCell) {
		if rowIdAsEntry > rightRowId && rightPointerPage != 0 {
			return false, rightPointerPage, 0
		}
	} else if page.btreeType == int(TableBtreeLeafCell) {
		if rowIdAsEntry > rightRowId {
			return false, page.pageNumber, 0
		}
	}

	fmt.Println("pointer 2")

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
		fmt.Println("pointer 3")

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

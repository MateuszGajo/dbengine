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
}

var usableSpacePerPage = 3

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

// func getSiblings() ([][]Cell, []Cell) {

// }

// func getRootData() []Cell {

// }

var rootPageNode = Node{
	isOverflow: true,
	leaf:       false,
}

// we get a page but convert it to node
// we need to have a middle layer that is gonna translate all this stuff
//

var zeroPage = PageParsed{
	dbHeader:         header(),
	dbHeaderSize:     100,
	cellAreaParsed:   [][]byte{},
	rightMostpointer: []byte{0, 0, 0, 1},
	isOverflow:       false,
	leftSibling:      nil,
	rightSiblisng:    nil,
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
	dbHeader:       DbHeader{},
	dbHeaderSize:   0,
	cellAreaParsed: [][]byte{[]byte{0, 0, 0, 0, 0, 1}, []byte{0, 0, 0, 0, 0, 2}, []byte{0, 0, 0, 0, 0, 3}, []byte{0, 0, 0, 0, 0, 4}},
	isOverflow:     true,
	leftSibling:    nil,
	isLeaf:         false,
	rightSiblisng:  nil,
}

// divider
// [{1 3}]
// pager
// [[{1 1} {1 2}] [{1 4}]]

var secondPage = PageParsed{}

func getNode(pageNumber int) PageParsed {
	if pageNumber == 1 {
		return firstPage
	} else if pageNumber == 0 {
		return zeroPage
	}

	panic("shouldn't be here")

}

func balancingForNode(pageNumber int, parents []int) {
	node := getNode(pageNumber)

	siblings := []PageParsed{}

	fmt.Println("checkpoint1")

	if node.leftSibling != nil {
		page := getNode(*node.leftSibling)
		siblings = append(siblings, page)
	}

	siblings = append(siblings, node)

	if node.rightSiblisng != nil {
		page := getNode(*node.rightSiblisng)
		siblings = append(siblings, page)
	}

	isRoot := len(parents) == 0

	fmt.Println("checkpoint2")

	if !node.isOverflow {
		return
	}

	if isRoot {
		// take root page
		// insert a new page and move that taken from root and run balancing algoritm
		// root := getRootData()
		rootPageNode = Node{
			isOverflow: false,
		}
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
	for i, v := range siblings {
		for _, vN := range v.cellAreaParsed {
			fmt.Println("row id???", vN[4:6])
			var rowId int
			if node.isLeaf {
				rowIdUint64, _ := DecodeVarint(vN[2:4])
				rowId = int(rowIdUint64)
			} else {
				rowIdUint64, _ := DecodeVarint(vN[4:6])
				rowId = int(rowIdUint64)
			}

			cellToDistribute = append(cellToDistribute, Cell{size: 1, pageNumber: int(rowId)})
		}
		// divider should be taken from parens, by taken i mean removed

		if i < len(siblings)-1 {
			divider, newStartIndex, endIndex := getDivider(v.pageNumber)
			if startIndex == 0 {
				startIndex = newStartIndex
			}
			if i < len(siblings)-2 {
				endIndex = endIndex
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

	if len(numberOfCellPerPage) != len(siblings) {
		// basically allocating new page
		newPage := PageParsed{}
		secondPage = newPage
		siblings = append(siblings, secondPage)

	}

	// fix pointers
	lastSibling := siblings[len(siblings)-1]

	lastSibling.rightSiblisng = oldLastSibling.rightSiblisng

	//somehow fix parent pointers

	// lastSibling.po

	deivider, pages := redistribution(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute, len(siblings), node, siblings)

	fmt.Println("divider")
	fmt.Println(deivider)
	fmt.Println("pager")
	fmt.Println(pages)

	balancingForNode(parent, parents)

}

func leaf_bias(cells []Cell, node PageParsed) ([]int, []int) {

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

	return totalSizeInEachPage, numberOfCellPerPage
}

// account for underflow towards the right

func accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage []int, cellToDistribute []Cell, node PageParsed) ([]int, []int) {

	// PAGE 1 [1,2,3], PageTwo [4]
	//4 -1 = 3 -1 =2

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

var pageZero = Node{
	isOverflow:   false,
	pageNumber:   0,
	leaf:         false,
	indexes:      []Index{{id: 3, leftPointer: 1}, {id: 6, leftPointer: 2}},
	rightPointer: 3,
}
var pageOne = Node{
	isOverflow: false,
	pageNumber: 1,
	leaf:       true,
	indexes:    []Index{{id: 1, leftPointerNull: true}, {id: 2, leftPointerNull: true}},
}

var pageTwo = Node{
	isOverflow: false,
	pageNumber: 2,
	leaf:       true,
	indexes:    []Index{{id: 4, leftPointerNull: true}, {id: 5, leftPointerNull: true}},
}

var pageThree = Node{
	isOverflow: false,
	pageNumber: 3,
	leaf:       true,
	indexes:    []Index{{id: 7, leftPointerNull: true}, {id: 8, leftPointerNull: true}},
}

func getNode2(pageNumber int) Node {
	if pageNumber == 0 {
		return pageZero
	} else if pageNumber == 1 {
		return pageOne
	} else if pageNumber == 2 {
		return pageTwo
	} else if pageNumber == 3 {
		return pageThree
	}
	panic("shouldn't run this in get")
}

func updateNode(newNode Node, pageNumber int) {
	fmt.Println("update page")
	fmt.Println(pageNumber)
	if pageNumber == 0 {
		pageZero = newNode
		return
	} else if pageNumber == 1 {
		pageOne = newNode
		return
	} else if pageNumber == 2 {
		pageTwo = newNode
		return
	} else if pageNumber == 3 {
		pageThree = newNode
		return
	}
	panic("shouldn't run this in updsae")
}

func insert(entry int) {
	ok, _, node := search(0, []int{}, entry)

	if ok {
		// update condition
		// node.value = new value form entry
		fmt.Println("updated node value!!!")

		return
	}
	// insert condition

	node.indexes = append(node.indexes, Index{id: entry})

	if len(node.indexes) > int(usableSpacePerPage) {
		node.isOverflow = true
	}
	updateNode(node, node.pageNumber)

	fmt.Println("nodes")
	fmt.Printf("%+v \n", pageZero)
	fmt.Printf("%+v \n", pageOne)
	fmt.Printf("%+v \n", pageTwo)
	fmt.Printf("%+v \n", pageThree)
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

func search(pageNumber int, parents []int, entry int) (bool, int, Node) {
	ok, result, node := binary_search(pageNumber, entry)

	fmt.Println("result")
	fmt.Println(ok)
	fmt.Println(result)
	if !ok && node.leaf {
		fmt.Println("didnt find what we are looking for")
		return false, 0, node
	}
	if ok {
		return ok, result, node
	}

	return search(result, parents, entry)
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

func binary_search(pageNumber int, rowIdAsEntry int) (bool, int, Node) {
	// node := getNode2(pageNumber)
	node := PageParsed{}

	for i := 0; i < len(node.cellArea); i++ {
		if entry == node.indexes[i].id {
			return true, node.indexes[i].id, node
		} else if entry < node.indexes[i].id {
			return false, node.indexes[i].leftPointer, node
		}
	}
	return false, node.rightPointer, node
}

// binary searhc???, just basically iterate over celll

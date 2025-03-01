package main

import (
	"fmt"
	"math"
)

var maxNumberOfKeys float64 = 0
var minNumberOfKeys float64 = 0

func aa() {
	maxNumberOfKeys = math.Floor(float64(PageSize-12) / 8)
	minNumberOfKeys = math.Floor(maxNumberOfKeys / 2)
}

// 1. We need sturcutre to have children ranges, and if it index, or if it a cell?

type Node struct {
	indexes    []int
	leaf       bool
	isOverflow bool
}

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
///                   +-----------------|   1,2,3,4  |
///                                    +--------+

///                                     +--------+ Root Node
///                   +-----------------|   4   | --------------------
///                  /                  +--------+
///                 /                                             \
///            +-------+                                     +----------+ empty
///       +----|  1,2,3  |----+                   +------------| 			|------------+
///      /     +-------+     \                 /             +----------+             \
///     /          |          \               /                /      \                \

// 3.2 Second look it goes right to left, readjust the entires so the tree is balances, move 3 as root, move 4 to righT???, no we have leaf baias distribution

// leaf bias distribution

type Cell struct {
	size int
}

var usableSpacePerPage = 4

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
func balancingForNode() {
	cellToDistribute := []Cell{{size: 1}, {size: 1}, {size: 1}, {size: 1}}

	totalSizeInEachPage, numberOfCellPerPage := leaf_bias(cellToDistribute)

	totalSizeInEachPage, numberOfCellPerPage = accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute)
}

func leaf_bias(cells []Cell) ([]int, []int) {

	totalSizeInEachPage := []int{0}
	numberOfCellPerPage := []int{0}

	for _, v := range cells {
		i := len(totalSizeInEachPage) - 1

		if totalSizeInEachPage[i]+v.size <= usableSpacePerPage {
			fmt.Println("how many time we enter here?")
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

func accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage []int, cellToDistribute []Cell) ([]int, []int) {

	// PAGE 1 [1,2,3], PageTwo [4]
	//4 -1 = 3 -1 =2

	divCell := len(cellToDistribute) - numberOfCellPerPage[len(numberOfCellPerPage)-1] - 1

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
			totalSizeInEachPage[i-1] -= cellToDistribute[divCell-1].size
			divCell--
		}
	}

	return totalSizeInEachPage, numberOfCellPerPage
}

func redistribution(totalSizeInEachPage, numberOfCellPerPage []int, cellToDistribute []Cell) ([]Cell, [][]Cell) {
	dividers := []Cell{}
	pages := [][]Cell{[]Cell{}, []Cell{}}
	pageNumber := -1
	cellIndex := 0
	siblingsLength := 2 // we allocate two pages for this

	for i, v := range numberOfCellPerPage {
		fmt.Println("enter?")
		pageNumber++
		for range v {
			if cellIndex < len(cellToDistribute) {
				pages[pageNumber] = append(pages[pageNumber], cellToDistribute[cellIndex])
				cellIndex++
			}
		}

		if i < siblingsLength-1 {
			if cellIndex >= len(cellToDistribute) {
				panic("should never occur")
			}

			dividers = append(dividers, cellToDistribute[cellIndex])
			cellIndex++
		}
	}

	return dividers, pages
}

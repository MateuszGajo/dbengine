package main

import (
	"fmt"
	"testing"
)

func TestLeaftBiasDistribution(t *testing.T) {
	cellToDistribute := []Cell{{size: 1}, {size: 2}, {size: 3}, {size: 1}}

	totalSizeInEachPage, numberOfCellPerPage := leaf_bias(cellToDistribute)

	fmt.Println("totalSize")
	fmt.Println(totalSizeInEachPage)
	fmt.Println("numberofCellPerPage")
	fmt.Println(numberOfCellPerPage)

	totalSizeInEachPage, numberOfCellPerPage = accountForUnderflowToardsRight(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute)

	fmt.Println("totalSize")
	fmt.Println(totalSizeInEachPage)
	fmt.Println("numberofCellPerPage")
	fmt.Println(numberOfCellPerPage)

	dividers, pages := redistribution(totalSizeInEachPage, numberOfCellPerPage, cellToDistribute)

	fmt.Println("distribtuion, divders")
	fmt.Println(dividers)
	fmt.Println("pages")
	fmt.Println(pages)

	t.Errorf("test")

}

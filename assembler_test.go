package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAssemblePage(t *testing.T) {

	data := []byte{}
	cells := createCell(TableBtreeLeafCell, nil, "alice", nil)
	fmt.Println("after creating a cell")
	btreeHeader := BtreeHeaderSchema(TableBtreeLeafCell, cells, nil)
	fmt.Println("after btree header scgena")
	zeros := make([]byte, PageSize-len(btreeHeader)-len(cells.data))
	data = append(data, btreeHeader...)
	data = append(data, zeros...)
	data = append(data, cells.data...)
	res := parseReadPage(data, 1, MockFileInfo{SizeVal: 10})

	assembledPage := assembleDbPage(res)

	if !reflect.DeepEqual(data, assembledPage) {
		fmt.Println("assembled page")
		fmt.Println(assembledPage)
		fmt.Println("raw data")
		fmt.Println(data)
		t.Error("Asembled page is different than input passed")
	}
}

func TestAssembleHeader(t *testing.T) {

	dbHeader := header()
	assembledHeader := assembleDbHeader(dbHeader)
	fmt.Println(len(assembledHeader))
	parseHeader := parseDbHeader(assembledHeader)

	if !reflect.DeepEqual(dbHeader, parseHeader) {
		t.Errorf("Header are different, expected: %v, got %v", dbHeader, parseHeader)
	}

}

package main

// func TestParseDbPageWithOnlyHeader(t *testing.T) {
// 	btreeHeader := updatePage(TableBtreeLeafCell, CreateCell{dataLength: 0}, nil)
// 	zeros := make([]byte, PageSize-len(btreeHeader))
// 	data := append(btreeHeader, zeros...)
// 	res := parseReadPage(data, 1)

// 	if res.btreeType != int(TableBtreeLeafCell) {
// 		t.Errorf("Expected: %v tree type, insted we got: %v", TableBtreeLeafCell, res.btreeType)
// 	}

// 	if res.framgenetedArea != 0 {
// 		t.Errorf("Expected fragmeneted area to be: %v instead we got: %v", 0, res.framgenetedArea)
// 	}

// 	if res.freeBlock != 0 {
// 		t.Errorf("Expected first free block address to be: %v, insted we got: %v", 0, res.freeBlock)
// 	}

// 	if res.numberofCells != 0 {
// 		t.Errorf("Expected numbrs of cell to be: %v, insted we got: %v", 0, res.numberofCells)
// 	}

// 	if res.startCellContentArea != PageSize {
// 		t.Errorf("Expected start cell of content area to be: %v, insted we got: %v", PageSize, res.startCellContentArea)
// 	}

// }

// func TestParseDbPage(t *testing.T) {

// 	data := []byte{}
// 	cells := createCell(TableBtreeLeafCell, nil, "alice", nil)
// 	btreeHeader := updatePage(TableBtreeLeafCell, cells, nil)
// 	zeros := make([]byte, PageSize-len(btreeHeader)-len(cells.data))
// 	data = append(data, btreeHeader...)
// 	data = append(data, zeros...)
// 	data = append(data, cells.data...)
// 	res := parseReadPage(data, 1)

// 	if res.btreeType != int(TableBtreeLeafCell) {
// 		t.Errorf("Expected: %v tree type, insted we got: %v", TableBtreeLeafCell, res.btreeType)
// 	}

// 	if res.framgenetedArea != 0 {
// 		t.Errorf("Expected fragmeneted area to be: %v instead we got: %v", 0, res.framgenetedArea)
// 	}

// 	if res.freeBlock != 0 {
// 		t.Errorf("Expected first free block address to be: %v, insted we got: %v", 0, res.freeBlock)
// 	}

// 	if res.numberofCells != 1 {
// 		t.Errorf("Expected numbrs of cell to be: %v, insted we got: %v", 1, res.numberofCells)
// 	}

// 	if res.startCellContentArea != PageSize-len(cells.data) {
// 		t.Errorf("Expected start cell of content area to be: %v, insted we got: %v", PageSize-len(cells.data), res.startCellContentArea)
// 	}

// 	if res.latesRow.rowId != 1 {
// 		t.Errorf("Expected latestes row id to be: %v, insted we got: %v", 1, res.latesRow.rowId)
// 	}

// 	if !reflect.DeepEqual(res.latesRow.data, cells.data) {
// 		t.Errorf("Expected latest row data to be: %v, insted we got: %v", cells.data, res.latesRow.data)
// 	}

// 	if len(res.pointers) > 2 {
// 		t.Errorf("Expected to be only : %v pointers, insted we got: %v", 1, len(res.pointers)/2)
// 	}
// 	if binary.BigEndian.Uint16(res.pointers[:2]) != uint16(PageSize-len(cells.data)) {
// 		t.Errorf("Expected : %v, insted we got: %v", cells, res.latesRow.data)
// 	}
// }

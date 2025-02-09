package main

import (
	"fmt"
	"reflect"
	"strings"
)

func (server ServerStruct) handleCreateTableSqlQuery(parsedData PageParsed, parsedQuery []ParsedValue, input string) {
	// parse query lets assume output bellow
	createSqlQueryData := parseCreateTableQuery(parsedQuery, input)

	// rowContentLength, row := createSchemaCell(v.input, v.queryData, parsedData.latesRow)

	// Starting point is 2
	var pointerInSchemaToData int = 2 // db is not initialize, first page for schema, second for data
	//first read schema, starting from page 1, we will deal with multipage then
	if len(parsedData.latesRow.columns) > 3 {
		if parsedData.latesRow.columns[3].columnType != "1" {
			panic("We are expecting this to be internal index schema")
		}
		pointerInSchemaToData = int(parsedData.latesRow.columns[3].columnValue[0])
		// TODO: make this work
		pointerInSchemaToData = parsedData.dbInfo.pageNumber + 1
	}

	// Create a cell pass cell that are already there(latest pointer etc), get created cell's length + pointers, i think length is undded, pointer is all we care about, we only need all bytes, currnet cell + new one, + header, there rest should be zeros, for starting cell value we have latest pointers, and pointers array is already there
	// Modify headers fields after cell has been added
	// ADD cell, defined how many cell we are added and how many values in eachg cell (length will be calculate internally)
	cell := createCell(TableBtreeLeafCell, &parsedData, string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointerInSchemaToData, createSqlQueryData.rawQuery)
	allCells := cell.data
	allCells = append(allCells, parsedData.cellArea...)
	// We need to create/modify/add a btree of type 13 because this is schema, starting of first page, then maybe is overflowing to other, lets stick for this first one for now
	btreeHeader := BtreeHeaderSchema(TableBtreeLeafCell, cell, &parsedData)
	// Zeros data
	zerosLength := PageSize - parsedData.dbHeaderSize - len(btreeHeader) - len(allCells)
	zerosSpace := make([]byte, zerosLength)
	//General method to save the daata to disk i guess
	dataToSave := assembleDbHeader(parsedData.dbHeader)
	dataToSave = append(dataToSave, btreeHeader...)
	dataToSave = append(dataToSave, zerosSpace...)
	dataToSave = append(dataToSave, allCells...)

	server.writeToFile(dataToSave, 0, server.firstPage, server.conId)

	btreeHeaderForData := BtreeHeaderSchema(TableBtreeLeafCell, CreateCell{dataLength: 0, data: []byte{}}, nil)
	zerosLength = PageSize - len(btreeHeaderForData)
	zerosSpace = make([]byte, zerosLength)

	emptyDataPage := btreeHeaderForData
	emptyDataPage = append(emptyDataPage, zerosSpace...)

	server.writeToFile(emptyDataPage, pointerInSchemaToData-1, server.firstPage, server.conId)

}

func (server ServerStruct) handleCreaqteSqlQuery(parsedData PageParsed, parsedQuery []ParsedValue, input string) {

	if len(parsedQuery) < 2 {
		fmt.Printf("%+v", parsedQuery)
		panic("expect at lest length of two, we got:")
	}

	switch SQLCreateQueryObjectType(strings.ToLower(parsedQuery[1].data)) {
	case SqlQueryTableObjectType:
		server.handleCreateTableSqlQuery(parsedData, parsedQuery, input)
	default:
		panic("Object type not implement: " + parsedQuery[1].data)
	}
}

// insert into user(name) values('Alice')

func findTableNameInSchemaPage(page PageParsed, value string) []PageParseColumn {
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// TODO: add some commons
	// TESt everythins
	// Write insert to file
	// Verify it w sqlite3
	// Refactor!!
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// stupid search for now, go row by row
	cellArea := page.cellArea
	for len(cellArea) > 0 {

		rowLength := cellArea[0] + 2 // 2, for length, row id
		if rowLength > 127 {
			panic("handle this scenarion, where row length bigger than 127")
		}
		row := cellArea[:rowLength]
		cellArea = cellArea[rowLength:]
		rowHeaderLength := row[2]

		latestRowHeaders := row[3 : 3-1+int(rowHeaderLength)] // 3 - 1 (-1 because of header length contains itself)
		latestRowValues := row[3-1+int(rowHeaderLength):]

		// fmt.Println("row headers")
		// fmt.Println(latestRowHeaders)
		// fmt.Println("row valyes")
		// fmt.Println(string(latestRowValues))
		fmt.Println("parse column, find table name in schema page")
		var rowColumns []PageParseColumn = parseDbPageColumn(latestRowHeaders, latestRowValues)

		if reflect.DeepEqual(rowColumns[1].columnValue, []byte(value)) {
			return rowColumns
		}
	}
	panic("couldn't find schema")

}

func validateData(validator string, data string) bool {
	switch validator {
	case "string":
		return true
	case "integer":
		return true
	default:
		panic("not implemented data validator")
	}
}

func (server ServerStruct) handleInsertSqlQuery(parsedData PageParsed, parsedQuery []ParsedValue, input string) []byte {
	if strings.ToLower(parsedQuery[1].data) != validInsertObjectType {
		panic("expected " + validInsertObjectType + " got " + parsedQuery[1].data)
	}
	// WE need here schema pages
	// parsedData.
	// 	parsedQuery[2].data

	fmt.Println("that is intresting me")
	fmt.Println(parsedQuery[2].data)
	fmt.Println("parsed data")
	fmt.Printf("%+v \n", parsedData)
	res := findTableNameInSchemaPage(parsedData, parsedQuery[2].data)
	fmt.Println("did we found schema, yes")
	fmt.Println(res)
	_, parsedVal := genericParser(string(res[4].columnValue))
	// TODO: maybe remove input, its added as raw query
	createSqlQueryData := parseCreateTableQuery(parsedVal, input)

	fmt.Println("schema data")
	fmt.Println(createSqlQueryData)

	fmt.Println("currency query")
	fmt.Println(parsedQuery)

	currentQueryColumns := parsedQuery[3]
	currentQueryValues := parsedQuery[5]

	if len(currentQueryColumns.dataNested) == 0 {
		panic("expect at least one column")
	}

	if len(currentQueryColumns.dataNested) != len(currentQueryValues.dataNested) {
		panic("Number of columns passed, should match with number of values")
	}
	var dataToWrite []interface{}
mainloop:
	for _, v := range createSqlQueryData.columns {
		columnName := v.columnName

		for j := 0; j < len(currentQueryColumns.dataNested); j++ {
			if currentQueryColumns.dataNested[j].data == columnName {
				columnValue := currentQueryValues.dataNested[j].data
				if !validateData("integer", columnValue) {
					panic("invalid type")
				}
				//parse value
				if columnValue[0] == byte('\'') {
					if columnValue[len(columnValue)-1] != '\'' {
						panic("invalid value, string started with ' should end with '")
					} else {
						dataToWrite = append(dataToWrite, columnValue[1:len(columnValue)-1])
						continue mainloop
					}
				}
				dataToWrite = append(dataToWrite, columnValue)
				continue mainloop
			}

			// if createSqlQueryData.columns[j].columnName == columnName.data {

			// 	if(strings.Contains(columnData.data,"'")) {
			// 		dataToWrite
			// 	}
			// 	// TODO: also need to validate constrains as not null
			// 	continue mainloop
			// }
		}
		// we didn't find value
		// if() TODO add check if is there is not null constrain, return err accordingly
		// panic("couldn't find column")
		// TODO: if is autoincrement or default value, insert value
		dataToWrite = append(dataToWrite, nil)

	}

	// for _,v := range current
	fmt.Println("passed validation")
	fmt.Println("values to write")
	fmt.Println(dataToWrite...)
	parsedData = PageParsed{
		dbHeader:             DbHeader{},
		dbHeaderSize:         0,
		btreeType:            int(TableBtreeLeafCell),
		freeBlock:            0,
		numberofCells:        0,
		startCellContentArea: 0,
		framgenetedArea:      0,
		rightMostpointer:     []byte{},
		pointers:             []byte{},
		cellArea:             []byte{},
		latesRow:             &LastPageParseLatestRow{},
		dbInfo:               DbInfo{},
	}
	cell := createCell(TableBtreeLeafCell, nil, dataToWrite...)
	allCells := cell.data
	// TODO Concatenate this
	allCells = append(allCells, []byte{}...)

	fmt.Println("cells")
	fmt.Println(cell)
	btreeHeader := BtreeHeaderSchema(TableBtreeLeafCell, cell, &parsedData)
	// Zeros data
	zerosLength := PageSize - len(btreeHeader) - len(allCells)
	// - len(allCells)
	zerosSpace := make([]byte, zerosLength)
	//General method to save the daata to disk i guess
	dataToSave := []byte{}
	dataToSave = append(dataToSave, btreeHeader...)
	dataToSave = append(dataToSave, zerosSpace...)
	dataToSave = append(dataToSave, allCells...)

	// TODO: its not that simple as writing to first page need to handle it
	dataPage := res[3].columnValue

	fmt.Println("datapage")
	fmt.Println(dataPage)

	fmt.Println("Data to save")
	// fmt.Println(dataToSave)
	fmt.Println(len(dataToSave))

	server.writeToFile(dataToSave, int(dataPage[0])-1, parsedData, server.conId)

	return []byte{}
}

func (server ServerStruct) handleActionType(parsedQuery []ParsedValue, input string, parsedData PageParsed) {

	if len(parsedQuery) < 1 {
		fmt.Printf("%+v", parsedQuery)
		panic("expect at lest length of one, we got:")
	}

	switch SQLQueryActionType(strings.ToLower(parsedQuery[0].data)) {
	case SqlQueryCreateActionType:

		server.handleCreaqteSqlQuery(parsedData, parsedQuery, input)
	case SqlQueryInsertActionType:
		server.handleInsertSqlQuery(parsedData, parsedQuery, input)
	default:
		panic("Unsported sql query type: " + parsedQuery[0].data)
	}
}

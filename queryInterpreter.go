package main

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	InternalSchemaIndexStart = 2
)

func handleCreateTableSqlQuery(parsedData LastPageParsed, parsedQuery []ParsedValue, input string) []byte {
	// parse query lets assume output bellow
	createSqlQueryData := parseCreateTableQuery(parsedQuery, input)

	// rowContentLength, row := createSchemaCell(v.input, v.queryData, parsedData.latesRow)

	// Starting point is 2
	internalIndexForSchema := InternalSchemaIndexStart
	//first read schema, starting from page 1, we will deal with multipage then
	if len(parsedData.latesRow.columns) > 3 {
		if parsedData.latesRow.columns[3].columnType != "1" {
			panic("We are expecting this to be internal index schema")
		}
		internalIndexForSchema := int(parsedData.latesRow.columns[3].columnValue[0])
		internalIndexForSchema++
	}

	// Create a cell pass cell that are already there(latest pointer etc), get created cell's length + pointers, i think length is undded, pointer is all we care about, we only need all bytes, currnet cell + new one, + header, there rest should be zeros, for starting cell value we have latest pointers, and pointers array is already there
	// Modify headers fields after cell has been added
	// ADD cell, defined how many cell we are added and how many values in eachg cell (length will be calculate internally)
	cell := createCell(TableBtreeLeafCell, parsedData.latesRow, string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, internalIndexForSchema, createSqlQueryData.rawQuery)
	allCells := cell.data
	allCells = append(allCells, parsedData.cellArea...)
	// We need to create/modify/add a btree of type 13 because this is schema, starting of first page, then maybe is overflowing to other, lets stick for this first one for now
	btreeHeader := BtreeHeaderSchema(TableBtreeLeafCell, cell, parsedData)
	// Zeros data
	zerosLength := PageSize - len(parsedData.dbHeader) - len(btreeHeader) - len(allCells)
	zerosSpace := make([]byte, zerosLength)
	//General method to save the daata to disk i guess
	dataToSave := parsedData.dbHeader
	dataToSave = append(dataToSave, btreeHeader...)
	dataToSave = append(dataToSave, zerosSpace...)
	dataToSave = append(dataToSave, allCells...)

	return dataToSave
}

func handleCreaqteSqlQuery(parsedData LastPageParsed, parsedQuery []ParsedValue, input string) []byte {

	if len(parsedQuery) < 2 {
		fmt.Printf("%+v", parsedQuery)
		panic("expect at lest length of two, we got:")
	}

	switch SQLCreateQueryObjectType(strings.ToLower(parsedQuery[1].data)) {
	case SqlQueryTableObjectType:
		return handleCreateTableSqlQuery(parsedData, parsedQuery, input)
	default:
		panic("Object type not implement: " + parsedQuery[1].data)
	}
}

// insert into user(name) values('Alice')

func findTableNameInSchemaPage(page LastPageParsed, value string) []PageParseColumn {
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// ``````````````````````````````````
	// TODO: add some common for that
	// LETS add parsing column as common
	// then we need to find schema in insert query
	// find if passer talbe name is defined in schema
	// then add value in tree
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

func handleInsertSqlQuery(parsedData LastPageParsed, parsedQuery []ParsedValue, input string) []byte {
	if strings.ToLower(parsedQuery[1].data) != validInsertObjectType {
		panic("expected " + validInsertObjectType + " got " + parsedQuery[1].data)
	}
	// WE need here schema pages
	// parsedData.
	// 	parsedQuery[2].data

	fmt.Println("that is intresting me")
	fmt.Println(parsedQuery[2].data)
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

mainloop:
	for i := 0; i < len(currentQueryColumns.dataNested); i++ {
		columnName := currentQueryColumns.dataNested[i]
		columnData := currentQueryValues.dataNested[i]

		for j := 0; j < len(createSqlQueryData.columns); j++ {
			if createSqlQueryData.columns[j].columnName == columnName.data {
				if !validateData("integer", columnData.data) {
					panic("invalid type")
				}
				continue mainloop
			}
		}
		panic("couldn't find column")

	}
	fmt.Println("passed validation")
	// cell := createCell(TableBtreeLeafCell, parsedData.latesRow, string(createSqlQueryData.objectType))
	// allCells := cell.data
	// allCells = append(allCells, parsedData.cellArea...)

	return []byte{}
}

func handleActionType(parsedQuery []ParsedValue, input string, parsedData LastPageParsed) []byte {

	if len(parsedQuery) < 1 {
		fmt.Printf("%+v", parsedQuery)
		panic("expect at lest length of one, we got:")
	}

	switch SQLQueryActionType(strings.ToLower(parsedQuery[0].data)) {
	case SqlQueryCreateActionType:

		return handleCreaqteSqlQuery(parsedData, parsedQuery, input)
	case SqlQueryInsertActionType:
		return handleInsertSqlQuery(parsedData, parsedQuery, input)
	default:
		panic("Unsported sql query type: " + parsedQuery[0].data)
	}
}

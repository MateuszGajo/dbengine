package main

import (
	"fmt"
	"reflect"
	"strings"
)

type SQLQueryActionType string

const (
	SqlQueryCreateActionType SQLQueryActionType = "create"
	SqlQueryInsertActionType SQLQueryActionType = "insert"

	//TODO: fill them
)

var validQueryActionTypes = map[SQLQueryActionType]struct{}{
	SqlQueryCreateActionType: {},
	SqlQueryInsertActionType: {},
}

type SQLCreateQueryObjectType string

const (
	SqlQueryDatabaseObjectType SQLCreateQueryObjectType = "database"
	SqlQueryTableObjectType    SQLCreateQueryObjectType = "table"
	SqlQueryIndexObjectType    SQLCreateQueryObjectType = "index"
	//TODO fill them
)

var validCreateQueryObjectTypes = map[SQLCreateQueryObjectType]struct{}{
	SqlQueryDatabaseObjectType: {},
	SqlQueryTableObjectType:    {},
	SqlQueryIndexObjectType:    {},
}

var validInsertObjectType = "into"

type QueryType string

const (
	QueryTypeTable QueryType = "table"
	QueryTypeIndex QueryType = "index"
)

func (server *ServerStruct) handleCreateTableSqlQuery(parsedQuery []ParsedValue, input string) {
	// parse query lets assume output bellow
	createSqlQueryData := parseCreateTableQuery(parsedQuery, input)

	// rowContentLength, row := createSchemaCell(v.input, v.queryData, parsedData.latesRow)

	// Starting point is 2
	var pointerInSchemaToData int = 2 // db is not initialize, first page for schema, second for data
	//first read schema, starting from page 1, we will deal with multipage then
	if len(server.firstPage.latesRow.columns) > 3 {
		if server.firstPage.latesRow.columns[3].columnType != "1" {
			panic("We are expecting this to be internal index schema")
		}
		pointerInSchemaToData = int(server.firstPage.latesRow.columns[3].columnValue[0])
		pointerInSchemaToData++
		// TODO: make this work
		// pointerInSchemaToData = server.dbInfo.pageNumber + 1
	}

	// Create a cell pass cell that are already there(latest pointer etc), get created cell's length + pointers, i think length is undded, pointer is all we care about, we only need all bytes, currnet cell + new one, + header, there rest should be zeros, for starting cell value we have latest pointers, and pointers array is already there
	// Modify headers fields after cell has been added
	// ADD cell, defined how many cell we are added and how many values in eachg cell (length will be calculate internally)
	cell := createCell(TableBtreeLeafCell, &server.firstPage, string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointerInSchemaToData, createSqlQueryData.rawQuery)
	allCells := cell.data
	allCells = append(allCells, server.firstPage.cellArea...)
	// We need to create/modify/add a btree of type 13 because this is schema, starting of first page, then maybe is overflowing to other, lets stick for this first one for now
	writer := NewWriter()
	reader := NewReader(server.conId)
	if server.firstPage.btreeType == 0 {
		server.firstPage.btreeType = int(TableBtreeLeafCell)
	}
	rootPage := server.updatePageRoot(&server.firstPage, 0, *reader, *writer, string(createSqlQueryData.objectType), createSqlQueryData.entityName, createSqlQueryData.entityName, pointerInSchemaToData, createSqlQueryData.rawQuery)

	firstPage := parseReadPage(rootPage, 0)
	server.firstPage = firstPage
	server.firstPage.dbHeader.schemaCookie++

	// Create empty page for data
	btreeHeaderForData := updateBtreeHeaderLeafTable(CreateCell{dataLength: 0, data: []byte{}}, nil)
	zerosLength := PageSize - len(btreeHeaderForData)
	zerosSpace := make([]byte, zerosLength)

	emptyDataPage := btreeHeaderForData
	emptyDataPage = append(emptyDataPage, zerosSpace...)

	writer.writeToFile(emptyDataPage, pointerInSchemaToData-1, server.conId, &server.firstPage)

}

func (server ServerStruct) handleCreaqteSqlQuery(parsedQuery []ParsedValue, input string) error {

	if len(parsedQuery) < 2 {
		fmt.Printf("%+v", parsedQuery)
		return fmt.Errorf("Invalid query")
		panic("expect at lest length of two, we got:")
	}

	switch SQLCreateQueryObjectType(strings.ToLower(parsedQuery[1].data)) {
	case SqlQueryTableObjectType:
		server.handleCreateTableSqlQuery(parsedQuery, input)
	default:
		panic("Object type not implement: " + parsedQuery[1].data)
	}

	return nil
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

		var rowColumns []PageParseColumn = parseDbPageColumn(latestRowHeaders, latestRowValues)

		if reflect.DeepEqual(rowColumns[1].columnValue, []byte(value)) {
			return rowColumns
		}
	}
	panic("couldn't find schema")

}

func transformColumnData(validator string, data string) (string, error) {
	switch validator {
	// TODO: deal with lower/upper case
	case "TEXT":
		if data[0] != byte('\'') || data[len(data)-1] != byte('\'') {
			return "", fmt.Errorf("Values should be wrapped with '' ")
		}

		return data[1 : len(data)-1], nil
	case "INTEGER":
		return data, nil
	default:
		fmt.Println(validator)
		panic("not implemented data validator")
	}
}

//TODO: implement this

func checkColumnConstrain(constrains []string, data string) error {
	return nil
	// switch validator {
	// case "text":

	// default:
	// 	panic("not implemented data validator")
	// }
}

func getColumnData(columnDefinition []SQLQueryColumnConstrains, selectedColumns, columnValues ParsedValue) ([]interface{}, error) {
	var dataToWrite = []interface{}{}
mainloop:
	for _, v := range columnDefinition {
		columnName := v.columnName

		for j := 0; j < len(selectedColumns.dataNested); j++ {
			if selectedColumns.dataNested[j].data == columnName {
				columnValue := columnValues.dataNested[j].data
				err := checkColumnConstrain(v.constrains, columnValue)
				if err != nil {
					return dataToWrite, err
				}
				transformedValue, err := transformColumnData(v.columnType, columnValue)
				if err != nil {
					return dataToWrite, err
				}

				dataToWrite = append(dataToWrite, transformedValue)
				continue mainloop
			}

		}
		dataToWrite = append(dataToWrite, nil)

	}

	return dataToWrite, nil

}

func (server ServerStruct) handleInsertSqlQuery(parsedQuery []ParsedValue, input string) error {
	if strings.ToLower(parsedQuery[1].data) != validInsertObjectType {
		return fmt.Errorf("expected " + validInsertObjectType + " got " + parsedQuery[1].data)
	}

	res := findTableNameInSchemaPage(server.firstPage, parsedQuery[2].data)

	dataStartOnPage := res[3].columnValue
	if len(dataStartOnPage) > 1 {
		panic("handle this later")
	}
	dataStartOnPageInt := int(dataStartOnPage[0])

	reader := NewReader(server.conId)
	page := reader.readDbPage(dataStartOnPageInt - 1)

	parsedData := parseReadPage(page, dataStartOnPageInt-1)

	_, parsedVal := genericParser(string(res[4].columnValue))

	createSqlQueryData := parseCreateTableQuery(parsedVal, input)

	currentQueryColumns := parsedQuery[3]
	currentQueryValues := parsedQuery[5]

	if len(currentQueryColumns.dataNested) == 0 {
		return fmt.Errorf("expect at least one column")
	}

	if len(currentQueryColumns.dataNested) != len(currentQueryValues.dataNested) {
		return fmt.Errorf("Number of columns passed, should match with number of values")
	}
	dataToWrite, err := getColumnData(createSqlQueryData.columns, currentQueryColumns, currentQueryValues)

	fmt.Println("data to write")
	fmt.Println(dataToWrite...)

	if err != nil {
		return err
	}

	// cell := createCell(TableBtreeLeafCell, nil, dataToWrite...)
	// allCells := cell.data
	// // TODO Concatenate this
	// allCells = append(allCells, []byte{}...)

	writer := NewWriter()
	// reader := NewReader(server.conId);
	dataPage := res[3].columnValue

	server.updatePageRoot(&parsedData, int(dataPage[0])-1, *reader, *writer, dataToWrite...)

	// btreeHeader := updatePage(TableBtreeLeafCell, cell, &parsedData)
	// zerosLength := PageSize - len(btreeHeader) - len(allCells)

	// zerosSpace := make([]byte, zerosLength)
	// dataToSave := []byte{}
	// dataToSave = append(dataToSave, btreeHeader...)
	// dataToSave = append(dataToSave, zerosSpace...)
	// dataToSave = append(dataToSave, allCells...)

	// NewWriter().writeToFile(dataToSave, int(dataPage[0])-1, server.firstPage, server.conId)

	return nil
}

func (server *ServerStruct) handleActionType(parsedQuery []ParsedValue, input string) error {

	if len(parsedQuery) < 1 {
		fmt.Printf("%+v", parsedQuery)
		panic("expect at lest length of one, we got:")
	}
	var err error
	switch SQLQueryActionType(strings.ToLower(parsedQuery[0].data)) {
	case SqlQueryCreateActionType:

		err = server.handleCreaqteSqlQuery(parsedQuery, input)
	case SqlQueryInsertActionType:
		server.handleInsertSqlQuery(parsedQuery, input)
	default:
		panic("Unsported sql query type: " + parsedQuery[0].data)
	}

	return err
}

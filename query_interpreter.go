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

func (server *ServerStruct) handleCreateTableSqlQuery(parsedData PageParsed, parsedQuery []ParsedValue, input string) {
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
		// pointerInSchemaToData = server.dbInfo.pageNumber + 1
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

	writer := NewWriter()

	// let think how to update server.firstPage

	writer.writeToFile(dataToSave, 0, server.firstPage, server.conId)
	fmt.Println("afete write, trying to parse")
	dbParsed := parseReadPage(dataToSave, 0)

	fmt.Println("after parsing")
	fmt.Printf("%+v", dbParsed)

	server.firstPage = dbParsed

	btreeHeaderForData := BtreeHeaderSchema(TableBtreeLeafCell, CreateCell{dataLength: 0, data: []byte{}}, nil)
	zerosLength = PageSize - len(btreeHeaderForData)
	zerosSpace = make([]byte, zerosLength)

	emptyDataPage := btreeHeaderForData
	emptyDataPage = append(emptyDataPage, zerosSpace...)

	fmt.Println("write to page?")
	fmt.Println(pointerInSchemaToData - 1)
	fmt.Println("what has firstPage")
	fmt.Printf("%+v", server.firstPage)
	// fmt.Println(emptyDataPage)

	writer.writeToFile(emptyDataPage, pointerInSchemaToData-1, server.firstPage, server.conId)

}

func (server ServerStruct) handleCreaqteSqlQuery(parsedData PageParsed, parsedQuery []ParsedValue, input string) error {

	if len(parsedQuery) < 2 {
		fmt.Printf("%+v", parsedQuery)
		return fmt.Errorf("Invalid query")
		panic("expect at lest length of two, we got:")
	}

	switch SQLCreateQueryObjectType(strings.ToLower(parsedQuery[1].data)) {
	case SqlQueryTableObjectType:
		server.handleCreateTableSqlQuery(parsedData, parsedQuery, input)
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
	fmt.Println("find in schema")
	fmt.Println(page)
	fmt.Println(value)
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
	fmt.Println("constrains")
	fmt.Println(constrains)
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

func (server ServerStruct) handleInsertSqlQuery(parsedData PageParsed, parsedQuery []ParsedValue, input string) error {
	if strings.ToLower(parsedQuery[1].data) != validInsertObjectType {
		return fmt.Errorf("expected " + validInsertObjectType + " got " + parsedQuery[1].data)
	}

	res := findTableNameInSchemaPage(parsedData, parsedQuery[2].data)

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

	if err != nil {
		return err
	}

	// for _,v := range current
	fmt.Println("passed validation")
	fmt.Println("values to write")
	fmt.Println(dataToWrite...)
	// TODO: what the hell is this ,remove it
	parsedData1 := PageParsed{
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
	}
	cell := createCell(TableBtreeLeafCell, nil, dataToWrite...)
	allCells := cell.data
	// TODO Concatenate this
	allCells = append(allCells, []byte{}...)

	fmt.Println("cells")
	fmt.Println(cell)
	btreeHeader := BtreeHeaderSchema(TableBtreeLeafCell, cell, &parsedData1)
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

	NewWriter().writeToFile(dataToSave, int(dataPage[0])-1, server.firstPage, server.conId)

	return nil
}

func (server *ServerStruct) handleActionType(parsedQuery []ParsedValue, input string, parsedData PageParsed) error {

	if len(parsedQuery) < 1 {
		fmt.Printf("%+v", parsedQuery)
		panic("expect at lest length of one, we got:")
	}
	var err error
	switch SQLQueryActionType(strings.ToLower(parsedQuery[0].data)) {
	case SqlQueryCreateActionType:

		err = server.handleCreaqteSqlQuery(parsedData, parsedQuery, input)
	case SqlQueryInsertActionType:
		server.handleInsertSqlQuery(parsedData, parsedQuery, input)
	default:
		panic("Unsported sql query type: " + parsedQuery[0].data)
	}

	return err
}

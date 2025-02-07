package main

import (
	"fmt"
	"testing"
)

func TestInsertSQL(t *testing.T) {
	//TODO: add variantions test INSERT INTO, INSERT    INTO (multiple space etc...)
	insertQuery := "INSERT INTO user(name) values('Alice')"
	bytesRead, data := genericParser(insertQuery)

	if bytesRead != len(insertQuery) {
		t.Errorf("Should read: %v of data, insted it read: %v bytes", len(insertQuery), bytesRead)
	}

	if len(data) < 6 {
		t.Error("Expect parsed query data to have at least items 1. insert 2.into 3.table name 4.colum names 5.values keyword 6. values data")
	}

	if data[0].data != "INSERT" {
		t.Errorf("Expected first item to be INSERT, insted we got: %v", data[0].data)
	}

	if data[1].data != "INTO" {
		t.Errorf("Expected second item to be INTO, insted we got: %v", data[1].data)
	}

	if data[2].data != "user" {
		t.Errorf("Expected third item to be table name: user, insted we got: %v", data[2].data)
	}

	if data[3].dataType != ParsedDataTypeBracket || len(data[3].dataNested) != 1 {
		t.Errorf("Epected single column value in bracket, we got: %+v", data[3])
	}

	if data[4].data != "values" {
		t.Errorf("Expected fourth item to be values keyword, insted we got: %v", data[4].data)
	}

	if data[5].dataType != ParsedDataTypeBracket || len(data[5].dataNested) != 1 || data[5].dataNested[0].data != "'Alice'" {
		t.Errorf("Epected data to be in the bracket and single values of alice insted we got: %+v", data[5])
	}

}

func TestCreateSQL(t *testing.T) {
	insertQuery := "CREATE TABLE user(id INTEGER PRIMARY KEY, name TEXT)"
	bytesRead, data := genericParser(insertQuery)

	if bytesRead != len(insertQuery) {
		t.Errorf("Should read: %v of data, insted it read: %v bytes", len(insertQuery), bytesRead)
	}

	if len(data) < 4 {
		t.Errorf("Expect parsed query data to have at least items 4 items, got: %v", data)
	}

	if data[0].data != "CREATE" {
		t.Errorf("Expected first item to be INSERT, insted we got: %v", data[0].data)
	}

	if data[1].data != "TABLE" {
		t.Errorf("Expected second item to be INTO, insted we got: %v", data[1].data)
	}

	if data[2].data != "user" {
		t.Errorf("Expected third item to be table name: user, insted we got: %v", data[2].data)
	}

	if data[3].dataType != ParsedDataTypeBracket || len(data[3].dataNested) != 6 {
		t.Errorf("Epected parsed bracket value:\n %+v", data[3])
	}

	if data[3].dataNested[0].data != "id" {
		t.Errorf("Expect first column name to be user, got: %v", data[3].dataNested[0].data)
	}

	if data[3].dataNested[1].data != "INTEGER" {
		t.Errorf("Expect column type to be integer, got: %v", data[3].dataNested[1].data)
	}

	if data[3].dataNested[2].data != "PRIMARY" {
		t.Errorf("Expect column to contains attribute PRIMARY, got: %v", data[3].dataNested[2].data)
	}

	if data[3].dataNested[3].data != "KEY" {
		t.Errorf("Expect column to contains attribute KEY, got: %v", data[3].dataNested[3].data)
	}

	if data[3].dataNested[3].dataType != ParsedDataTypeSeparator {
		t.Error("Expect to column with type and attributes and with separator")
	}

	if data[3].dataNested[4].data != "name" {
		t.Errorf("Expect first column name to be name, got: %v", data[3].dataNested[4].data)
	}

	if data[3].dataNested[5].data != "TEXT" {
		t.Errorf("Expect column type to be integer, got: %v", data[3].dataNested[5].data)
	}
}

func TestParseColumnAttributes(t *testing.T) {
	input := "CREATE TABLE user(id INTEGER PRIMARY KEY, name TEXT)"
	_, data := genericParser(input)

	fmt.Println(data)
	parsedColumns := parseSqlQueryColumnAttributes(data[3])

	if parsedColumns[0].columnName != "id" {
		t.Errorf("Expected column name to be: %v, insted we got: %v", "id", parsedColumns[0].columnName)
	}

	if parsedColumns[0].columnType != "INTEGER" {
		t.Errorf("Expected column type to be: %v, insted we got: %v", "INTEGER", parsedColumns[0].columnType)
	}

	if parsedColumns[0].constrains[0] != "PRIMARY" || parsedColumns[0].constrains[1] != "KEY" {
		t.Errorf("Expected constains to be PRIMARY KEY, got: %v", parsedColumns[0].constrains[0]+" "+parsedColumns[0].constrains[1])
	}

	if parsedColumns[1].columnName != "name" {
		t.Errorf("Expected second column name to be name, got: %v", parsedColumns[1].columnName)
	}

	if parsedColumns[1].columnType != "TEXT" {
		t.Errorf("Expected second column type to be text, got: %v", parsedColumns[1].columnType)
	}
}

func TestParseCreateQuery(t *testing.T) {
	// TODO: either make all lower case or don't touch it
	input := "CREATE TABLE user(id INTEGER PRIMARY KEY, name TEXT)"
	_, data := genericParser(input)
	res := parseCreateTableQuery(data, input)

	if res.action != "create" {
		t.Errorf("Expected action to be CREATE, we got: %v", res.action)
	}

	if res.objectType != "table" {
		t.Errorf("Expected object  to be TABLE, we got: %v", res.objectType)
	}

	if res.entityName != "user" {
		t.Errorf("Expected entity to to be user, we got: %v", res.entityName)
	}
	if res.rawQuery != input {
		t.Errorf("expected raw query to be equal to input, epxected: %v, got:%v", input, res.rawQuery)
	}

	if res.columns[0].columnName != "id" {
		t.Errorf("Expected first column name to be id, got: %v", res.columns[0].columnName)
	}
	if res.columns[0].columnType != "INTEGER" {
		t.Errorf("Expected first column type to be integer, got: %v", res.columns[0].columnType)
	}

	if res.columns[0].constrains[0] != "PRIMARY" || res.columns[0].constrains[1] != "KEY" {
		t.Errorf("Expected constains to be PRIMARY KEY, got: %v", res.columns[0].constrains[0]+" "+res.columns[0].constrains[1])
	}

	if res.columns[1].columnName != "name" {
		t.Errorf("Expected second column name to be name, got: %v", res.columns[1].columnName)
	}

	if res.columns[1].columnType != "TEXT" {
		t.Errorf("Expected second column type to be text, got: %v", res.columns[1].columnType)
	}
}

package main

import (
	"testing"
)

func TestInsertSQL(t *testing.T) {
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

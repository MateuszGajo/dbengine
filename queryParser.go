package main

import (
	"strings"
)

func parseSqlQueryColumnAttributes(parsedQuery ParsedValue) []SQLQueryColumnConstrains {
	columnConstrains := []SQLQueryColumnConstrains{}
	columnConstrain := SQLQueryColumnConstrains{}
	columnCounter := 0

	for _, v := range parsedQuery.dataNested {
		if columnCounter == 0 {
			columnConstrain.columnName = v.data
		} else if columnCounter == 1 {
			columnConstrain.columnType = v.data
		} else if columnCounter >= 2 {
			columnConstrain.constrains = append(columnConstrain.constrains, v.data)
		}

		columnCounter++

		if v.dataType == ParsedDataTypeSeparator {
			columnConstrains = append(columnConstrains, columnConstrain)
			columnConstrain = SQLQueryColumnConstrains{}
			columnCounter = 0
		}

	}
	columnConstrains = append(columnConstrains, columnConstrain)

	return columnConstrains
}

func parseCreateTableQuery(parsedQuery []ParsedValue, input string) CreateActionQueryData {

	res := CreateActionQueryData{
		rawQuery: input,
	}

	// WRite parsing (columns, eg id integer primary key, name text)

	if len(parsedQuery) < 4 {
		panic("Expect query to contains at lest four element Create Table [name](attributes....)")
	}

	res.action = SQLQueryActionType(strings.ToLower(parsedQuery[0].data))
	res.objectType = SQLCreateQueryObjectType(strings.ToLower(parsedQuery[1].data))
	res.entityName = parsedQuery[2].data

	if parsedQuery[3].dataType != ParsedDataTypeBracket {
		panic("Expected bracket here")
	}

	if len(parsedQuery[3].dataNested) < 2 {
		panic("should contains at lest one column")
	}

	res.columns = parseSqlQueryColumnAttributes(parsedQuery[3])

	return res
}

// insert into user(name) values('Alice')
var allowedQueryOperators = map[rune]struct{}{
	'>': {},
	'*': {},
	'<': {},
	'-': {},
	'=': {},
}

type ParsedDataType string

const (
	ParsedDataTypeBracket     ParsedDataType = "ParsedDataTypeBracket"
	ParsedDataTypeOperator    ParsedDataType = "ParsedDataTypeOperator"
	ParsedDataTypeSimpleValue ParsedDataType = "ParsedDataTypesSimpleValue"
	ParsedDataTypeSeparator   ParsedDataType = "ParsedDataTypesSeparator"
)

type ParsedValue struct {
	dataType   ParsedDataType
	data       string
	dataNested []ParsedValue
}

func genericParser(input string) (int, []ParsedValue) {
	result := []ParsedValue{}

	start := 0
	i := 0
	for i = 0; i < len(input); i++ {
		// Create   Table
		// 012345678
		if input[i] == ' ' {
			if i == start {
				start = i + 1
				continue
				// Case for multiple space `insert   into``, etc
			}

			val := ParsedValue{
				data:     input[start:i],
				dataType: ParsedDataTypeSimpleValue,
			}
			start = i + 1
			result = append(result, val)
		}

		if input[i] == ',' {
			val := ParsedValue{
				data:     input[start:i],
				dataType: ParsedDataTypeSeparator,
			}
			start = i + 1
			result = append(result, val)
		}

		if input[i] == ')' {
			val := ParsedValue{
				data:     input[start:i],
				dataType: ParsedDataTypeSimpleValue,
			}
			result = append(result, val)
			return i + 1, result
		}

		if _, ok := allowedQueryOperators[rune(input[i])]; ok {

			val := ParsedValue{
				data:     input[start:i],
				dataType: ParsedDataTypeOperator,
			}
			result = append(result, val)
			start = i + 1
		}

		if input[i] == '(' {
			prevVal := input[start:i]
			if len(prevVal) > 0 {
				val := ParsedValue{
					data:     prevVal,
					dataType: ParsedDataTypeSimpleValue,
				}
				result = append(result, val)
			}
			bytesRead, nestedVal := genericParser(input[i+1:])
			// fmt.Println("exit reading bracket")
			// fmt.Printf("bytes read: %v\n", bytesRead)
			// fmt.Printf("nested Val: %+v\n", nestedVal)
			val := ParsedValue{
				dataNested: nestedVal,
				dataType:   ParsedDataTypeBracket,
			}
			result = append(result, val)
			start = i + bytesRead + 1
			i += bytesRead
		}
	}

	// bytesRead, val := genericParserBracket()

	return i, result
}

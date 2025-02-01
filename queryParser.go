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

func getValuesBySeparator(input string, separator rune) []string {
	var values []string
	start := 0
	for i := 0; i < len(input); i++ {
		if input[i] == byte(separator) {
			data := input[start:i]
			start = i + 1
			values = append(values, data)

		}
	}
	values = append(values, input[start:])

	return values
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

// LETS finish this, general idea:
// 1.Outer function parsing everyting, *,< as added as operatror
// 2. Inner function is for parsing everytinng in bracket, "," is added a seaprator

// func genericParserBracket(input string) (int, []ParsedValueBase) {
// 	result := []ParsedValueBase{}

// 	start := 0
// 	i := 0
// 	for i = 0; i < len(input); i++ {
// 		if input[i] == ' ' {
// 			if i == start-1 {
// 				start = i
// 				continue
// 				// Case for multiple space `insert   into``, etc
// 			}

// 			val := ParsedValueBase{
// 				data:     input[start:i],
// 				dataType: ParsedDataTypeOperator,
// 			}
// 			start = i + 1
// 			result = append(result, val)
// 		}

// 		if input[i] == ',' {
// 			val := ParsedValueBase{
// 				data:     input[start:i],
// 				dataType: ParsedDataTypeSeparator,
// 			}
// 			start = i + 1
// 			result = append(result, val)
// 		}

// 		if input[i] == ')' {
// 			val := ParsedValueBase{
// 				data:     input[start:i],
// 				dataType: ParsedDataTypeSeparator,
// 			}
// 			start = i + 1
// 			result = append(result, val)
// 			return i + 1, result
// 		}

// 		if _, ok := allowedQueryOperators[rune(input[i])]; ok {

// 			val := ParsedValueBase{
// 				data:     input[start:i],
// 				dataType: ParsedDataTypeOperator,
// 			}
// 			result = append(result, val)
// 			start = i + 1
// 		}
// 	}

// 	return i, result
// }

func parseInserQuery(input string) CreateActionQueryData {

	res := CreateActionQueryData{
		rawQuery: input,
	}

	data := []string{}
	start := 0
	input = strings.TrimSpace(input)
	if input[len(input)-1] != ')' {
		panic("invalid query")
	}
	for i := 0; i < len(input); i++ {
		if input[i] == 32 {
			data = append(data, input[start:i])
			start = i + 1
		}

		if input[i] == '(' {
			start = i + 1
			break
		}
	}

	// WRite parsing (columns, eg id integer primary key, name text)

	if len(data) < 3 {
		panic("invalid sql, inser int [name], should have 3 words separated by space")
	}

	res.entityName = data[2]

	// ()
	// 0 - name, 1-type, 2+-attributes

	// last one is ()
	// id INTEGER PRIMARY KEY, name TEXT
	// columnConstrains := []SQLQueryColumnConstrains{}
	// for i := start; i < len(input)-1; i++ {
	// 	if input[i] == ',' {
	// 		data := input[start:i]
	// 		start = i + 1
	// 		columnConstrains = append(columnConstrains, parseSqlQueryColumnAttributes(data))

	// 	}

	// }

	// res.columns = columnConstrains

	return res
}

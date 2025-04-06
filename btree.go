package main

import (
	"fmt"
	"os"
	"reflect"
)

type CreateCell struct {
	dataLength int
	data       []byte
	rowId      int
}

type BtreeType int

const (
	TableBtreeLeafCell     BtreeType = 0x0d
	TableBtreeInteriorCell BtreeType = 0x05
	IndexBtreeLeafCell     BtreeType = 0x0a
	IndexBtreeInteriorCell BtreeType = 0x02
)

// func BtreeHeaderInteriorTable(cell CreateCell, parsedData *PageParsed) []byte {
// 	currentNumberOfCell := 0
// 	currentCellStart := PageSize
// 	pointers := []byte{}
// 	if parsedData != nil {
// 		currentNumberOfCell = parsedData.numberofCells
// 		currentCellStart = parsedData.startCellContentArea
// 		pointers = parsedData.pointers
// 	}

// 	var lastPointer int = PageSize

// 	if len(pointers) > 0 {
// 		lastPointer = int(binary.BigEndian.Uint16(parsedData.pointers[len(parsedData.pointers)-2 : len(parsedData.pointers)]))
// 	}

// 	if cell.dataLength > 0 {
// 		currentNumberOfCell += 1
// 		newCellStartPosition := lastPointer - cell.dataLength
// 		currentCellStart = newCellStartPosition
// 		pointers = append(pointers, intToBinary(newCellStartPosition, 2)...) // -2 is for row id and for length byte
// 	}

// 	bTreePageType := intToBinary(int(TableBtreeLeafCell), 1)
// 	firstFreeBlockOnPage := intToBinary(0, 2)
// 	numberOfCells := intToBinary(currentNumberOfCell, 2)

// 	startCellContentArea := intToBinary(currentCellStart, 2)
// 	framgentedFreeBytesWithingCellContentArea := intToBinary(0, 1)
// 	rightMostPointer := []byte{}

// 	data := []byte{}
// 	data = append(data, bTreePageType...)
// 	data = append(data, firstFreeBlockOnPage...)
// 	data = append(data, numberOfCells...)
// 	data = append(data, startCellContentArea...)
// 	data = append(data, framgentedFreeBytesWithingCellContentArea...)
// 	data = append(data, rightMostPointer...)
// 	data = append(data, pointers...)

// }

func updateBtreeHeaderLeafTable(cell CreateCell, parsedData *PageParsed) []byte {
	//This should be read from the page
	currentNumberOfCell := 0
	currentCellStart := PageSize
	pointers := []byte{}
	if parsedData != nil {
		currentNumberOfCell = parsedData.numberofCells
		currentCellStart = parsedData.startCellContentArea
		pointers = parsedData.pointers
	}

	if cell.dataLength > 0 {
		currentNumberOfCell += 1
		newCellStartPosition := currentCellStart - cell.dataLength
		currentCellStart = newCellStartPosition
		pointers = append(pointers, intToBinary(newCellStartPosition, 2)...) // -2 is for row id and for length byte
	}

	bTreePageType := intToBinary(int(TableBtreeLeafCell), 1)
	firstFreeBlockOnPage := intToBinary(0, 2)
	numberOfCells := intToBinary(currentNumberOfCell, 2)

	startCellContentArea := intToBinary(currentCellStart, 2)
	framgentedFreeBytesWithingCellContentArea := intToBinary(0, 1)

	data := []byte{}
	data = append(data, bTreePageType...)
	data = append(data, firstFreeBlockOnPage...)
	data = append(data, numberOfCells...)
	data = append(data, startCellContentArea...)
	data = append(data, framgentedFreeBytesWithingCellContentArea...)
	data = append(data, pointers...)

	return data
}

func updateBtreeHeaderLeafTableIntoInterior(rightMostPointer []byte, parsedData *PageParsed) []byte {
	//This should be read from the page
	// currentNumberOfCell := 0
	currentNumberOfCell := 0
	currentCellStart := PageSize
	pointers := []byte{}
	if parsedData != nil {
		currentNumberOfCell = parsedData.numberofCells
		currentCellStart = parsedData.startCellContentArea
		pointers = parsedData.pointers
	}

	// var lastPointer int = PageSize

	// if len(pointers) > 0 {
	// 	lastPointer = int(binary.BigEndian.Uint16(parsedData.pointers[len(parsedData.pointers)-2 : len(parsedData.pointers)]))
	// }

	// if len(parsedData.rightMostpointer) > 0 {
	// 	// cellToAdd := parsedData.rightMostpointer
	// 	// cellToAdd = append(cellToAdd, byte(rowId)) this calculation will be use to passed cell to this btreeheader START FROM HERE~~~
	// 	newCellStartPosition := lastPointer - len(cell)
	// 	currentCellStart = newCellStartPosition
	// 	pointers = append(pointers, intToBinary(newCellStartPosition, 2)...) // -2 is for row id and for length byte
	// }

	bTreePageType := intToBinary(int(TableBtreeInteriorCell), 1)
	firstFreeBlockOnPage := intToBinary(0, 2)
	numberOfCells := intToBinary(currentNumberOfCell, 2)

	startCellContentArea := intToBinary(currentCellStart, 2)
	framgentedFreeBytesWithingCellContentArea := intToBinary(0, 1)

	data := []byte{}
	data = append(data, bTreePageType...)
	data = append(data, firstFreeBlockOnPage...)
	data = append(data, numberOfCells...)
	data = append(data, startCellContentArea...)
	data = append(data, framgentedFreeBytesWithingCellContentArea...)
	data = append(data, rightMostPointer...)
	data = append(data, pointers...)

	return data
}

// Adding new schema
// 1. If first page 0x0d and fit the page, put it there
// 2. if first page 0x0d and doesnt fit the page, create 0x05, move 0x0d and add to right most pointer
// 3  if page 0x05 and fit in page that is in the right most pointer, put it there
// 4. if page 0x05 and doesnt fit in page, create new 0x0d, put it there, set it as right most pointer, and move right most pointer as cell conent area

// How to call header with created new cell, we shouldn't specify header type i think
// 1. What do we do when create schema?
// 1.1 Scan all schemas if it exists, if not we can create
// 1.2 create schema by following #adding new schema
// 2. What do we do when adding value
// 2.1 Schema schema if table exists, check all constrains for value, find the right page in schema
// 2.2 go to the page and put it there, or follow #addubg be scgena i think is quite the same

// How to structure code??
// 1. first we check constrains
// 2. then we get page that data should be
// 3.we pass page, and created cell
// 4. we should determin if there is a space to put value there, if not follow #adding schema
// !!So for init schema, i think not only should created defaukt header, but also default btree 0x0d for schema placeholder, and then should follow this code
// So BtreeHeaderSchema we are not passing btree type only, page parsed, we checked if we can put value or not and then accordingly we move pages, creating new headers etc..., change function name to maybe updateHEadeR????

// Now we create a cell, basically for header 0x0d
// 1. Create a raw vlaue for 0x0d, check if fit
// 2. fits no problems
// 3, doesnt fit, we need to split header, add 0x05 and then create a vlaue for header 0x05 and save it
// 4. we are passing pageParsed, what if we have 0x05 and few 0x0d already, we should pass value form right pointer!!, but what if we exeeced? we need to somehow know root page
// 5. We need to build strucutre that is root, and pointers???
// LEts talks about schema exmaple
// We know that schema starts at page 0
/// 1. We get page 0
// 2. maybe lets pass rootpage
// 3. we check if is 0x0d, if not look, for right most
// 4. we found 0x0d, if there is aspace we add it there
// 5. if pages is 0x05 we got to right pointer again again agina...
// 6. if we found 0x0d that is full?, what do we do now?
// 7. We need to create a new page, that ok, but how do we know that we need to update page 0x05, in case that we only update page, we can simply return it, but what in case of adding new??

// lets now assume that we update a page and return it simple as that
// What do we do when we need to create a new page???
// create a brand new 0x0d page, we go level up to 0x05
// we need to distingush newPage and updatedOne, updae only save at offset that right pointer is poiting to, when create a new page, we need to add it at the end and update right most pointer, and also save 0x05
// back to what is 0x05 is also full??
// create a new 0x05 page with pointer to new 0x0d and we pass need 0x05, in higher level 0x05 we add pointer to lovel level 0x05 and we save it, if we pass nil we dont save anything

// CReating cell values before, we need to know structure, for schema is simple we adding is as 0x0d but what for other
// What if we  inserting value, we need to look for schema, in schema we will find constains, and page
// then we check all constrains and pass root page to updatePage
// Back to schema, we first need to verify if schema already exsists

// lets not overcomplicate it, event if it means sacriface performance, lets focus on that later
// adding new schema:
// check if schema exist
//then look for right page

//adding new value
//look for schema and all constrainc
// validate against them and then look for right page

// EXTRA!!
// 5. if page 0x05 and we have not more cell content area, created another 0x05, leave it for later

//Somehow it works, but what if we need change a root???
// WE are getting return statement,

// lets assume we only use 0x05 and 0x0d for simplicity
// func (server *ServerStruct) updatePageRoot(page *PageParsed, pageNumber int, reader PageReader, writer WriterStruct, valuesToAdd ...interface{}) []byte {

// 	updatedPage, newPages := server.updatePage(page, pageNumber, reader, writer, valuesToAdd...)

// 	if updatedPage != nil {
// 		return updatedPage
// 	}

// 	rowId := newPages[0].rowId
// 	rightMostPointer := newPages[0].pointer

// 	// currentPageData := assembleDbPage(*page)
// 	// rightMostPointer := intToBinary(server.firstPage.dbHeader.dbSizeInPages, 2)
// 	cellToAdd := rightMostPointer
// 	if *rowId > 127 {
// 		panic("implement when row id greater than 127")
// 	}
// 	if len(cellToAdd) > 0 {
// 		cellToAdd = append(cellToAdd, byte(*rowId))
// 	}

// 	// writer.writeToFile(currentPageData, server.firstPage.dbHeader.dbSizeInPages, server.conId, server.firstPage)

// 	btreeHeader := updateBtreeHeaderLeafTableIntoInterior(rightMostPointer, &PageParsed{})
// 	allCells := cellToAdd
// 	// allCells = append(allCells, page.cellArea...)
// 	// Zeros data
// 	zerosLength := PageSize - page.dbHeaderSize - len(btreeHeader) - len(allCells)
// 	zerosSpace := make([]byte, zerosLength)
// 	//General method to save the daata to disk i guess
// 	dataToSave := assembleDbHeader(page.dbHeader)
// 	dataToSave = append(dataToSave, btreeHeader...)
// 	dataToSave = append(dataToSave, zerosSpace...)
// 	dataToSave = append(dataToSave, allCells...)

// 	//swipe pages

// 	writer.writeToFile(dataToSave, pageNumber, server.conId, &server.firstPage)

// 	return dataToSave
// }

type CreatedPages struct {
	pointer   []byte
	rowId     *int
	savedPAge []byte
}

// What do we wanna achive here???
// 1. We find the page (we look basically for the last page, looking for 0x0d)
// 2. update 0x0d, if no space, create a new 0x0d page, and then if first page and after removing header is space, we need to add here and create new page and current one move at the end
// Header logic is kinda

// simplify db logic for now
// only 0x05 pages and 0x0d pages
// no overflow
// all data append at the end
// 0x05 points to 0x0d,
// either content fit on last page or create a new one

// what about first page with header ??
// condition if first page doesnt fit, but fits with header remove it, move current page to last one and move header
// if doesnt fit, remove header and create new page

// lets simplify the flow
// 1. find page, but what, we are looking for lat page so on 0x05 right pointer, and 0x0d is right page (pass space needed)
// 2. insert into page(page, data)
// 2.1

// func (server *ServerStruct) insertValue(reader PageReader, pageNumber int) *PageParsed {
// findInitialPageInSchema()
// }

// func() findInitialPageInSchema() {

// }

// this is only for 0x0d
// What do we do here??
// maybe first check if there is space?

// if there is space update page
// what about header case??
// 2 ways of checking header
// check if there is first page
// then what ???

// if there is no space???

// func (server *ServerStruct) updatePage(page *PageParsed, pageNumber int, reader PageReader, writer WriterStruct, valuesToAdd ...interface{}) (updatedPage []byte, createdPages []CreatedPages) {

// 	if page.btreeType == int(TableBtreeLeafCell) {
// 		fmt.Println("etnter")
// 		cell := createCell(TableBtreeLeafCell, page, valuesToAdd...)
// 		fmt.Println("etnter2")
// 		btreeLeftCellHeaderLength := 8
// 		spaceAvilable := page.startCellContentArea - btreeLeftCellHeaderLength - len(page.pointers) - page.dbHeaderSize
// 		newPointerLength := 2
// 		newCellSpace := cell.dataLength + newPointerLength
// 		if spaceAvilable >= newCellSpace {
// 			//fits into page

// 			btreeHeader := updateBtreeHeaderLeafTable(cell, page)

// 			allCells := cell.data
// 			allCells = append(allCells, page.cellArea...)
// 			// Zeros data
// 			zerosLength := PageSize - page.dbHeaderSize - len(btreeHeader) - len(allCells)
// 			zerosSpace := make([]byte, zerosLength)
// 			//General method to save the daata to disk i guess
// 			dataToSave := []byte{}
// 			if page.dbHeaderSize > 0 {
// 				dataToSave = append(dataToSave, assembleDbHeader(page.dbHeader)...)
// 			}
// 			dataToSave = append(dataToSave, btreeHeader...)
// 			dataToSave = append(dataToSave, zerosSpace...)
// 			dataToSave = append(dataToSave, allCells...)

// 			writer.writeToFile(dataToSave, pageNumber, server.conId, &server.firstPage)

// 			return dataToSave, nil

// 			// updatedPage = dataToSave
// 			// return newPage, updatedPage, cell.rowId
// 		} else {
// 			fmt.Println("enter else")
// 			fmt.Println("enter else")
// 			fmt.Println("enter else")

// 			isFirstPageAndContentFit := pageNumber == 0 && spaceAvilable+100 >= newCellSpace

// 			if isFirstPageAndContentFit {
// 				fmt.Println("enter first condition")
// 				btreeHeader := updateBtreeHeaderLeafTable(cell, page)

// 				allCells := cell.data
// 				fmt.Println("checkpoint 0")
// 				allCells = append(allCells, page.cellArea...)
// 				// Zeros data
// 				zerosLength := PageSize - len(btreeHeader) - len(allCells)
// 				zerosSpace := make([]byte, zerosLength)
// 				//General method to save the daata to disk i guess
// 				dataToSave := []byte{}
// 				fmt.Println("checkpoint 1")
// 				dataToSave = append(dataToSave, btreeHeader...)
// 				dataToSave = append(dataToSave, zerosSpace...)
// 				dataToSave = append(dataToSave, allCells...)
// 				writer.writeToFile(dataToSave, server.firstPage.dbHeader.dbSizeInPages, server.conId, &server.firstPage)
// 				server.firstPage.dbHeader.dbSizeInPages++
// 				createdPage := CreatedPages{
// 					pointer:   intToBinary(server.firstPage.dbHeader.dbSizeInPages-1, 4),
// 					rowId:     &cell.rowId,
// 					savedPAge: dataToSave,
// 				}
// 				return nil, []CreatedPages{createdPage}
// 			} else if pageNumber == 0 {
// 				fmt.Println("enter second condition")
// 				btreeHeader := updateBtreeHeaderLeafTable(cell, page)

// 				allCells := page.cellArea
// 				// Zeros data
// 				zerosLength := PageSize - page.dbHeaderSize - len(btreeHeader) - len(allCells)
// 				zerosSpace := make([]byte, zerosLength)
// 				//General method to save the daata to disk i guess
// 				dataToSave := []byte{}
// 				dataToSave = append(dataToSave, btreeHeader...)
// 				dataToSave = append(dataToSave, zerosSpace...)
// 				dataToSave = append(dataToSave, allCells...)
// 				writer.writeToFile(dataToSave, server.firstPage.dbHeader.dbSizeInPages, server.conId, &server.firstPage)
// 				server.firstPage.dbHeader.dbSizeInPages++
// 				rowId := cell.rowId - 1
// 				createdPage := CreatedPages{
// 					pointer:   intToBinary(server.firstPage.dbHeader.dbSizeInPages-1, 4),
// 					rowId:     &rowId,
// 					savedPAge: dataToSave,
// 				}
// 				createdPages = append(createdPages, createdPage)
// 			}

// 			btreeHeader := updateBtreeHeaderLeafTable(cell, nil)

// 			allCells := cell.data
// 			// Zeros data
// 			zerosLength := PageSize - page.dbHeaderSize - len(btreeHeader) - len(allCells)
// 			zerosSpace := make([]byte, zerosLength)
// 			//General method to save the daata to disk i guess
// 			dataToSave := []byte{}

// 			dataToSave = append(dataToSave, btreeHeader...)
// 			dataToSave = append(dataToSave, zerosSpace...)
// 			dataToSave = append(dataToSave, allCells...)
// 			fmt.Println("save to file")
// 			writer.writeToFile(dataToSave, server.firstPage.dbHeader.dbSizeInPages, server.conId, &server.firstPage)

// 			createdPage := CreatedPages{
// 				pointer:   intToBinary(server.firstPage.dbHeader.dbSizeInPages-1, 4),
// 				rowId:     &cell.rowId,
// 				savedPAge: dataToSave,
// 			}
// 			createdPages = append(createdPages, createdPage)

// 			return nil, createdPages
// 		}
// 	} else if page.btreeType == int(TableBtreeInteriorCell) {

// 		readPAge := reader.readDbPage(int(binary.BigEndian.Uint32(page.rightMostpointer)))

// 		pageParsed := parseReadPage(readPAge, int(binary.BigEndian.Uint32(page.rightMostpointer)))

// 		updatedPage, createdPages = server.updatePage(&pageParsed, int(binary.BigEndian.Uint32(page.rightMostpointer)), reader, writer, valuesToAdd...)
// 	}

// 	if updatedPage == nil {
// 		return updatedPage, nil
// 	}
// 	//FOCUS now on this part

// 	if len(createdPages) == 0 {
// 		panic("should never be 0")
// 	}

// 	//
// 	// what in the case we need to

// 	cellToAdd := page.rightMostpointer
// 	rowId := createdPages[0].rowId
// 	if *rowId > 127 {
// 		panic("implement when row id greater than 127")
// 	}
// 	if len(cellToAdd) > 0 {
// 		cellToAdd = append(cellToAdd, byte(*rowId))
// 	}

// 	btreeLeftCellHeaderLength := 12
// 	spaceAvilable := page.startCellContentArea - btreeLeftCellHeaderLength - len(page.pointers)
// 	newPointerLength := 2
// 	newCellSpace := len(cellToAdd) + newPointerLength

// 	if spaceAvilable >= newCellSpace {
// 		btreeHeader := updateBtreeHeaderLeafTableIntoInterior(page.rightMostpointer, page)
// 		//fits into page
// 		allCells := cellToAdd
// 		allCells = append(allCells, page.cellArea...)
// 		// Zeros data
// 		zerosLength := PageSize - page.dbHeaderSize - len(btreeHeader) - len(allCells)
// 		zerosSpace := make([]byte, zerosLength)
// 		//General method to save the daata to disk i guess
// 		dataToSave := []byte{}
// 		if page.dbHeaderSize > 0 {
// 			dataToSave = append(dataToSave, assembleDbHeader(page.dbHeader)...)
// 		}
// 		dataToSave = append(dataToSave, btreeHeader...)
// 		dataToSave = append(dataToSave, zerosSpace...)
// 		dataToSave = append(dataToSave, allCells...)

// 		// updatedPage = pageData
// 		writer.writeToFile(dataToSave, pageNumber, server.conId, &server.firstPage)
// 		//TODO fix passing rowID

// 		return dataToSave, nil

// 	} else {
// 		btreeHeader := updateBtreeHeaderLeafTableIntoInterior(page.rightMostpointer, page)
// 		allCells := []byte{}
// 		allCells = append(allCells, page.cellArea...)
// 		// Zeros data
// 		zerosLength := PageSize - page.dbHeaderSize - len(btreeHeader) - len(allCells)
// 		zerosSpace := make([]byte, zerosLength)
// 		//General method to save the daata to disk i guess
// 		dataToSave := []byte{}
// 		if page.dbHeaderSize > 0 {
// 			dataToSave = append(dataToSave, assembleDbHeader(page.dbHeader)...)
// 		}
// 		dataToSave = append(dataToSave, btreeHeader...)
// 		dataToSave = append(dataToSave, zerosSpace...)
// 		dataToSave = append(dataToSave, allCells...)

// 		//TODO fix passing rowid
// 		writer.writeToFile(dataToSave, server.firstPage.dbHeader.dbSizeInPages, server.conId, &server.firstPage)
// 		//TODO fix passing rowID
// 		rowId := 0
// 		createdPage := CreatedPages{
// 			pointer:   intToBinary(server.firstPage.dbHeader.dbSizeInPages-1, 4),
// 			rowId:     &rowId,
// 			savedPAge: dataToSave,
// 		}
// 		createdPages = append(createdPages, createdPage)

// 		return nil, createdPages

// 	}

// }

// func updatePage(cell CreateCell, parsedData *PageParsed, rootPage *PageParsed) []byte {

// 	// switch btreeType {
// 	// case TableBtreeLeafCell:
// 	// 	return BtreeHeaderLeafTable(cell, parsedData)
// 	// case TableBtreeInteriorCell:
// 	// 	return BtreeHeaderInteriorTable(cell, parsedData)

// 	// default:
// 	// 	panic("btree not implemented" + string(btreeType))
// 	// }

// 	if parsedData.btreeType == int(TableBtreeLeafCell) {
// 		btreeLeftCellHeaderLength := 8
// 		spaceAvilable := parsedData.startCellContentArea - btreeLeftCellHeaderLength - len(parsedData.pointers)
// 		newPointerLength := 2
// 		newCellSpace := cell.dataLength + newPointerLength

// 		// TOOD: i think sqlite uses more complex calculation, but lets leave it as it is for now
// 		if spaceAvilable >= newCellSpace {
// 			//fits into page
// 			updatedHeader := updateBtreeHeaderLeafTable(cell, parsedData)
// 		} else if rootPage != nil {
// 			// created a mew 0x0d page, and attach it to root page
// 			//doesnt fit
// 		} else {
// 			// create 0x05 page, and update it
// 		}
// 	}

// }

func EncodeVarint(n uint64) []byte {
	var groups []byte

	// Always append at least one group.
	groups = append(groups, byte(n&0x7F))
	n >>= 7

	// Extract 7-bit groups.
	for n > 0 {
		groups = append(groups, byte(n&0x7F))
		n >>= 7
	}

	// Reverse the groups to convert from little-endian (LSB first)
	// to big-endian (most-significant group first).
	for i, j := 0, len(groups)-1; i < j; i, j = i+1, j-1 {
		groups[i], groups[j] = groups[j], groups[i]
	}

	// For every byte except the last, set the continuation flag.
	for i := 0; i < len(groups)-1; i++ {
		groups[i] |= 0x80
	}

	return groups
}

// DecodeVarint decodes a byte slice encoded in SQLite's varint format
// back into a uint64. It also returns the number of bytes read.
func DecodeVarint(data []byte) uint64 {
	var n uint64
	var bytesRead int

	for i, b := range data {
		// For the first 8 bytes, each contributes 7 bits.
		if i < 8 {
			n = (n << 7) | uint64(b&0x7F)
			bytesRead++
			// If the continuation flag is not set, we're done.
			if b&0x80 == 0 && i == 7 {
				return n
			}
		} else {
			// The 9th byte contains 8 bits of data.
			n = (n << 8) | uint64(b)
			bytesRead++
			return n
		}
	}
	return n
}

func calculateTextLength(value string) []byte {

	stringLen := 2*len(value) + 13

	return EncodeVarint(uint64(stringLen))

	// if stringLen <= 127 {
	// 	return []byte{byte(uint8(stringLen))}
	// } else {
	// 	//TODO: implement this
	// 	panic("implement calculate text length")
	// }
}

func createCell(btreeType BtreeType, rowId int, values ...interface{}) CreateCell {
	if btreeType == TableBtreeLeafCell {
		var columnValues []byte = []byte{}
		var columnLength []byte = []byte{}
		// rowId++

		for _, v := range values {
			switch v.(type) {
			case int:
				value := v.(int)
				if value > 255 {
					panic("need to handle this type laters")
				}
				columnValues = append(columnValues, byte(uint16(value)))
				columnLength = append(columnLength, byte(1))
			case string:
				value := v.(string)
				columnValues = append(columnValues, []byte(value)...)
				columnLength = append(columnLength, calculateTextLength(value)...)
			case nil:
				columnValues = append(columnValues, []byte{}...)
				columnLength = append(columnLength, 0)
			default:
				fmt.Fprintln(os.Stdout, []any{v}...)
				fmt.Println(reflect.TypeOf(v))
				panic("unssporrted cell type")
			}
		}

		headerLength := len(columnLength) + 1 // 5 column + 1 for current byte

		row := []byte{byte(rowId)}
		row = append(row, byte(headerLength))
		row = append(row, columnLength...)
		row = append(row, columnValues...)

		rowLength := len(row) - 1 // we don't count row id

		result := []byte{byte(rowLength)}
		result = append(result, row...)

		return CreateCell{
			dataLength: len(result),
			data:       result,
			rowId:      rowId,
		}
	}

	panic("not handle create cell")

}

func header() DbHeader {
	headerString := []byte("SQLite format 3\000")
	pageSize := PageSize
	writeFileVersion := intToBinary(LegacyFileWriteFormat, 1)
	readFileVersion := intToBinary(LegacyFileReadFormat, 1)
	// SQLite has the ability to set aside a small number of extra bytes at the end of every page for use by extensions. These extra bytes are used, for example, by the SQLite Encryption Extension to store a nonce and/or cryptographic checksum associated with each page. The "reserved space" size in the 1-byte integer at offset 20 is the number of bytes of space at the end of each page to reserve for extensions. This value is usually 0. The value can be odd.
	reservedByte := intToBinary(0, 1)
	maxEmbededPayloadFranction := intToBinary(64, 1)
	minEmbededPayloadFranction := intToBinary(32, 1)
	//The maximum and minimum embedded payload fractions and the leaf payload fraction values must be 64, 32, and 32. These values were originally intended to be tunable parameters that could be used to modify the storage format of the b-tree algorithm. However, that functionality is not supported and there are no current plans to add support in the future. Hence, these three bytes are fixed at the values specified.
	leafPayloadFranction := intToBinary(32, 1)
	// The file change counter is a 4-byte big-endian integer at offset 24 that is incremented whenever the database file is unlocked after having been modified. When two or more processes are reading the same database file, each process can detect database changes from other processes by monitoring the change counter. A process will normally want to flush its database page cache when another process modified the database, since the cache has become stale. The file change counter facilitates this.
	fileChangeCounter := 0
	sizeOfDataBaseInPages := 1
	// 1 page for schemas, 2 page for values
	// END of 0000001
	//The trunk pages are the primary pages in the freelist structure. Each trunk page can list multiple leaf pages (individual unused pages) or point to other trunk pages.
	//When data is deleted from the database, the pages that were used to store that data aren't immediately erased or discarded. Instead, they are added to the freelist so that they can be reused later.
	pageNumberFirstFreeListTrunk := intToBinary(0, 4)
	totalNumerOfFreeListPages := intToBinary(0, 4)

	//The schema cookie is a 4-byte big-endian integer at offset 40 that is incremented whenever the database schema changes.
	schemaCookie := 0
	// The schema format number is a 4-byte big-endian integer at offset 44. The schema format number is similar to the file format read and write version numbers at offsets 18 and 19 except that the schema format number refers to the high-level SQL formatting rather than the low-level b-tree formatting. Four schema format numbers are currently defined:

	//     Format 1 is understood by all versions of SQLite back to version 3.0.0 (2004-06-18).
	//     Format 2 adds the ability of rows within the same table to have a varying number of columns, in order to support the ALTER TABLE ... ADD COLUMN functionality. Support for reading and writing format 2 was added in SQLite version 3.1.3 on 2005-02-20.
	//     Format 3 adds the ability of extra columns added by ALTER TABLE ... ADD COLUMN to have non-NULL default values. This capability was added in SQLite version 3.1.4 on 2005-03-11.
	//     Format 4 causes SQLite to respect the DESC keyword on index declarations. (The DESC keyword is ignored in indexes for formats 1, 2, and 3.) Format 4 also adds two new boolean record type values (serial types 8 and 9). Support for format 4 was added in SQLite 3.3.0 on 2006-01-10.

	// New database files created by SQLite use format 4 by default. The legacy_file_format pragma can be used to cause SQLite to create new database files using format 1. The format version number can be made to default to 1 instead of 4 by setting SQLITE_DEFAULT_FILE_FORMAT=1 at compile-time.

	// If the database is completely empty, if it has no schema, then the schema format number can be zero.
	schemaFormatNumber := intToBinary(4, 4)
	// end of 0000002

	// The 4-byte big-endian signed integer at offset 48 is the suggested cache size in pages for the database file. The value is a suggestion only and SQLite is under no obligation to honor it. The absolute value of the integer is used as the suggested size. The suggested cache size can be set using the default_cache_size pragma.
	defaultPageSize := intToBinary(0, 4)
	// The page number of the largest root b-tree page when in auto-vacuum or incremental-vacuum modes, or zero otherwise
	// The two 4-byte big-endian integers at offsets 52 and 64 are used to manage the auto_vacuum and incremental_vacuum modes. If the integer at offset 52 is zero then pointer-map (ptrmap) pages are omitted from the database file and neither auto_vacuum nor incremental_vacuum are supported. If the integer at offset 52 is non-zero then it is the page number of the largest root page in the database file, the database file will contain ptrmap pages, and the mode must be either auto_vacuum or incremental_vacuum. In this latter case, the integer at offset 64 is true for incremental_vacuum and false for auto_vacuum. If the integer at offset 52 is zero then the integer at offset 64 must also be zero.
	pageNumberOfLargestBtreeToAutovacuum := intToBinary(0, 4)
	//The database text encoding. A value of 1 means UTF-8. A value of 2 means UTF-16le. A value of 3 means UTF-16be.
	databaseTextEncoding := intToBinary(1, 4)
	// The "user version" as read and set by the user_version pragma.
	userVersionNumber := intToBinary(0, 4)
	// end of  end of 0000003
	//// 64	4	True (non-zero) for incremental-vacuum mode. False (zero) otherwise.
	incrementalVacuumMode := intToBinary(0, 4)
	// The "Application ID" set by PRAGMA application_id.
	// The 4-byte big-endian integer at offset 68 is an "Application ID" that can be set by the PRAGMA application_id command in order to identify the database as belonging to or associated with a particular application. The application ID is intended for database files used as an application file-format. The application ID can be used by utilities such as file(1) to determine the specific file type rather than just reporting "SQLite3 Database". A list of assigned application IDs can be seen by consulting the magic.txt file in the SQLite source repository.
	applicationId := intToBinary(0, 4)
	// Reserved for expansion. Must be zero.
	reservedForExpansion := make([]byte, 20)
	// The 4-byte big-endian integer at offset 96 stores the SQLITE_VERSION_NUMBER value for the SQLite library that most recently modified the database file. The 4-byte big-endian integer at offset 92 is the value of the change counter when the version number was stored. The integer at offset 92 indicates which transaction the version number is valid for and is sometimes called the "version-valid-for number".
	// The change counter is incremented each time the database schema is modified, such as when:

	// A table is created or dropped.
	// An index is added or removed.
	// Other schema-altering operations occur.
	// we added one table
	// TODO: this need to be updated????
	versionValidForNumber := 0
	// end of 0000005
	versionNumber := intToBinary(3045001, 4)

	return DbHeader{
		headerString:               headerString,
		databasePageSize:           pageSize,
		databaseFileWriteVersion:   writeFileVersion,
		databaseFileReadVersion:    readFileVersion,
		reservedBytesSpace:         reservedByte,
		maxEmbeddedPayloadFraction: maxEmbededPayloadFranction,
		minEmbeddedPayloadFraction: minEmbededPayloadFranction,
		leafPayloadFraction:        leafPayloadFranction,
		fileChangeCounter:          fileChangeCounter,
		dbSizeInPages:              sizeOfDataBaseInPages,
		firstFreeListTrunkPage:     pageNumberFirstFreeListTrunk,
		totalNumberOfFreeListPages: totalNumerOfFreeListPages,
		schemaCookie:               schemaCookie,
		schemaFormatNumber:         schemaFormatNumber,
		defaultPageCacheSize:       defaultPageSize,
		largestBTreePage:           pageNumberOfLargestBtreeToAutovacuum,
		databaseEncoding:           databaseTextEncoding,
		userVersion:                userVersionNumber,
		incrementalVacuumMode:      incrementalVacuumMode,
		applicationId:              applicationId,
		reservedForExpansion:       reservedForExpansion,
		versionValidForNumber:      versionValidForNumber,
		sqlVersionNumber:           versionNumber,
	}

}

func BtreeHeaderValue() []byte {

	// 1	The one-byte flag at offset 0 indicating the b-tree page type.

	// A value of 2 (0x02) means the page is an interior index b-tree page.
	// A value of 5 (0x05) means the page is an interior table b-tree page.
	// A value of 10 (0x0a) means the page is a leaf index b-tree page.
	// A value of 13 (0x0d) means the page is a leaf table b-tree page.

	// Any other value for the b-tree page type is an error.

	bTreePageType := intToBinary(0x0d, 1)
	// 	The two-byte integer at offset 1 gives the start of the first freeblock on the page, or is zero if there are no freeblocks.
	firstFreeBlockOnPage := intToBinary(0, 2)
	// 	The two-byte integer at offset 3 gives the number of cells on the page.

	numberOfCells := intToBinary(1, 2)
	// we added one table so its 1
	// 5	2	The two-byte integer at offset 5 designates the start of the cell content area. A zero value for this integer is interpreted as 65536.
	//If a page contains no cells (which is only possible for a root page of a table that contains no rows) then the offset to the cell content area will equal the page size minus the bytes of reserved space. If the database uses a 65536-byte page size and the reserved space is zero (the usual value for reserved space) then the cell content offset of an empty page wants to be 65536. However, that integer is too large to be stored in a 2-byte unsigned integer, so a value of 0 is used in its place
	//
	startCell := 4096 - 8 - 2
	startCellContentArea := intToBinary(startCell, 2)
	// 	The one-byte integer at offset 7 gives the number of fragmented free bytes within the cell content area.
	framgentedFreeBytesWithingCellContentArea := intToBinary(0, 1)
	//The four-byte page number at offset 8 is the right-most pointer. This value appears in the header of interior b-tree pages only and is omitted from all other pages.
	// right most pointer points to address fro mthe left side u see on xxd
	// 00000000: 5351 4c69 7465 2066 6f72 6d61 7420 3300  SQLite format 3.
	// 00000010: 1000 0101 0040 2020 0000 0003 0000 0003  .....@  ........
	// 00000020: 0000 0000 0000 0000 0000 0002 0000 0004  ................
	// 00000030: 0000 0000 0000 0000 0000 0001 0000 0000  ................
	// 00000040: 0000 0000 0000 0000 0000 0000 0000 0000  ................
	// 00000050: 0000 0000 0000 0000 0000 0000 0000 0003  ................
	// 00000060: 002e 7689 0d00 0000 020f 6600 0fb6 0f66  ..v.......f....f
	// .... more bytes
	// 00000f60: 0000 0000 0000 4e02 0617 1717 017d 7461  ......N......}ta
	// 00000f70: 626c 6575 7365 7232 7573 6572 3203 4352  bleuser2user2.CR
	// 00000f80: 4541 5445 2054 4142 4c45 2075 7365 7232  EATE TABLE user2
	// 00000f90: 2869 6420 696e 7465 6765 7220 7072 696d  (id integer prim
	// 00000fa0: 6172 7920 6b65 792c 206e 616d 6532 3232  ary key, name222
	// 00000fb0: 2074 6578 7429 4801 0617 1515 0175 7461   text)H......uta
	// 00000fc0: 626c 6575 7365 7275 7365 7202 4352 4541  bleuseruser.CREA
	// 00000fd0: 5445 2054 4142 4c45 2075 7365 7228 6964  TE TABLE user(id
	// 00000fe0: 2069 6e74 6567 6572 2070 7269 6d61 7279   integer primary
	// 00000ff0: 206b 6579 2c20 6e61 6d65 2074 6578 7429   key, name text)

	// There is rightpointer at 0f66, at is pointing to line 0f60 and then add 6 bytes so it after 00, so starting from 4e which is length of cell,
	// there is also pointer 0fb6, that is point to 0fb0 and then add 6 bytes, so it stars from 48 which is cell length

	cellAreaPointerFirst := intToBinary(startCell, 2)
	cellAreaPointerTwo := intToBinary(0, 2)
	data := []byte{}
	data = append(data, bTreePageType...)
	data = append(data, firstFreeBlockOnPage...)
	data = append(data, numberOfCells...)
	data = append(data, startCellContentArea...)
	data = append(data, framgentedFreeBytesWithingCellContentArea...)
	data = append(data, cellAreaPointerFirst...)
	data = append(data, cellAreaPointerTwo...)

	return data
}

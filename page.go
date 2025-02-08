package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
)

// ````````````````````````````
// ````````````````````````````
// ````````````````````````````
// ````````TODO!!!!!!!``````
// ````````````````````````````
// ````````````````````````````
// ````````````````````````````
// 1. Update header with data save
// 2. Change a logic with getting data from file, don't call it LastPagePArse, just page
// 3. Add data to have multiple pages, multiple pages for schema too, implement this

// edit
// finisparseDbHeader(data []byte)
// then add writing to db pages
// fix error handling, return errr don't panic
func header() DbHeader {
	headerString := []byte("SQLite format 3\000")
	pageSize := 2
	writeFileVersion := intToBinary(LegacyFileWriteFormat, 1)
	readFileVersion := intToBinary(LegacyFileReadFormat, 1)
	// SQLite has the ability to set aside a small number of extra bytes at the end of every page for use by extensions. These extra bytes are used, for example, by the SQLite Encryption Extension to store a nonce and/or cryptographic checksum associated with each page. The "reserved space" size in the 1-byte integer at offset 20 is the number of bytes of space at the end of each page to reserve for extensions. This value is usually 0. The value can be odd.
	reservedByte := intToBinary(0, 1)
	maxEmbededPayloadFranction := intToBinary(64, 1)
	minEmbededPayloadFranction := intToBinary(32, 1)
	//The maximum and minimum embedded payload fractions and the leaf payload fraction values must be 64, 32, and 32. These values were originally intended to be tunable parameters that could be used to modify the storage format of the b-tree algorithm. However, that functionality is not supported and there are no current plans to add support in the future. Hence, these three bytes are fixed at the values specified.
	leafPayloadFranction := intToBinary(32, 1)
	// The file change counter is a 4-byte big-endian integer at offset 24 that is incremented whenever the database file is unlocked after having been modified. When two or more processes are reading the same database file, each process can detect database changes from other processes by monitoring the change counter. A process will normally want to flush its database page cache when another process modified the database, since the cache has become stale. The file change counter facilitates this.
	fileChangeCounter := 1
	sizeOfDataBaseInPages := 2
	// 1 page for schemas, 2 page for values
	// END of 0000001
	//The trunk pages are the primary pages in the freelist structure. Each trunk page can list multiple leaf pages (individual unused pages) or point to other trunk pages.
	//When data is deleted from the database, the pages that were used to store that data aren't immediately erased or discarded. Instead, they are added to the freelist so that they can be reused later.
	pageNumberFirstFreeListTrunk := intToBinary(0, 4)
	totalNumerOfFreeListPages := intToBinary(0, 4)

	//The schema cookie is a 4-byte big-endian integer at offset 40 that is incremented whenever the database schema changes.
	schemaCookie := 1
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
	versionValidForNumber := 1
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

func calculateTextLength(value string) []byte {

	stringLen := 2*len(value) + 13

	if stringLen <= 255 {
		return []byte{byte(uint8(stringLen))}
	} else {
		//TODO: implement this
		panic("implement calculate text length")
	}
}

func createCell(btreeType BtreeType, latestRow *PageParsed, values ...interface{}) CreateCell {
	if btreeType == TableBtreeLeafCell {
		var columnValues []byte = []byte{}
		var columnLength []byte = []byte{}
		var schemaRowId = 0
		if latestRow != nil {
			schemaRowId = latestRow.latesRow.rowId
		}

		schemaRowId++

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
				fmt.Println("hello what do we got here")
				fmt.Fprintln(os.Stdout, []any{values}...)
				panic("unssporrted cell type")
			}
		}

		headerLength := len(columnLength) + 1 // 5 column + 1 for current byte
		rowId := schemaRowId                  // first row 2 byes

		row := []byte{byte(rowId)}
		row = append(row, byte(headerLength))
		row = append(row, columnLength...)
		row = append(row, columnValues...)

		rowLength := len(row) - 1 // we don't count row id

		result := []byte{byte(rowLength)}
		result = append(result, row...)

		return CreateCell{
			dataLength: rowLength,
			data:       result,
		}
	}

	panic("not handle create cell")

}

func parseDbPageColumn(rowHeader []byte, rowValues []byte) []PageParseColumn {
	fmt.Println("parse db column")
	fmt.Println(rowHeader)
	fmt.Println(rowValues)
	var rowColumn []PageParseColumn
	for _, v := range rowHeader {
		if int(v) > 127 {
			panic("handle case that we have multiple bytes")
		}
		if int(v) == 0 {
			column := PageParseColumn{
				columnType:   string(strconv.Itoa(int(v))),
				columnLength: 0,
				columnValue:  []byte{0},
			}
			rowColumn = append(rowColumn, column)
			continue

		}
		if int(v) < 10 {
			//int
			column := PageParseColumn{
				columnType:   string(strconv.Itoa(int(v))),
				columnLength: 1,
				columnValue:  []byte{rowValues[0]},
			}
			rowColumn = append(rowColumn, column)
			rowValues = rowValues[1:]
			continue
		}
		if int(v) >= 10 && int(v) < 12 {
			panic("reserved values, shouldnt be used")
		}
		if int(v)%2 == 0 {
			//blob
			panic("implement hadnling blobs")

		} else {

			//string
			length := (int(v) - 13) / 2

			if length > len(rowValues) {
				panic("there is not enough data")
			}
			value := rowValues[:length]
			column := PageParseColumn{
				columnType:   "13",
				columnLength: length,
				columnValue:  value,
			}
			rowColumn = append(rowColumn, column)
			rowValues = rowValues[length:]
			continue
		}
		panic("should never enter this state in parsing")

	}
	return rowColumn
}

func readDbPage(pageNumber int) ([]byte, fs.FileInfo) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(pwd + "/aa.db"); os.IsNotExist(err) {
		return []byte{}, nil
	}
	fd, err := os.Open(pwd + "/aa.db") // Adjust path as needed
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	fd.Seek(int64(pageNumber)*int64(PageSize), 0)

	buff := make([]byte, PageSize) // Create a buffer for 1024 bytes
	n, err := io.ReadFull(fd, buff)
	if err != nil && err != io.EOF { // Handle EOF and other errors
		fmt.Println("Error reading file:", err)
		return nil, nil
	}

	fileInfo, err := fd.Stat()
	if err != nil || fileInfo == nil {

		panic("Error while getting information about file")
	}
	fmt.Println("size")
	fmt.Println(fileInfo.Size())

	return buff[:n], fileInfo // Return the bytes actually read
}

func parseDbHeader(data []byte) DbHeader {
	if len(data) != 100 {
		panic("header should have 100 bytes")
	}
	headerString := data[:16]
	databasePageSize := data[16:18]
	databaseFileWriteVersion := data[18:19]
	databaseFileReadVersion := data[19:20]
	reservedBytesSpace := data[20:21]
	maxEmbeddedPayloadFraction := data[21:22]
	minEmbeddedPayloadFraction := data[22:23]
	leafPayloadFraction := data[23:24]
	fileChangeCounter := data[24:28]
	dbSizeInPages := data[28:32]
	firstFreeListTrunkPage := data[32:36]
	totalNumberOfFreeListPages := data[36:40]
	schemaCookie := data[40:44]
	schemaFormatNumber := data[44:48]
	defaultPageCacheSize := data[48:52]
	largestBTreePage := data[52:56]
	databaseEncoding := data[56:60]
	userVersion := data[60:64]
	incrementalVacuumMode := data[64:68]
	applicationId := data[68:72]
	reservedForExpansion := data[72:92]
	versionValidForNumber := data[92:96]
	sqlVersionNumber := data[96:100]

	return DbHeader{
		headerString:               headerString,
		databasePageSize:           int(binary.BigEndian.Uint16(databasePageSize)),
		databaseFileWriteVersion:   databaseFileWriteVersion,
		databaseFileReadVersion:    databaseFileReadVersion,
		reservedBytesSpace:         reservedBytesSpace,
		maxEmbeddedPayloadFraction: maxEmbeddedPayloadFraction,
		minEmbeddedPayloadFraction: minEmbeddedPayloadFraction,
		leafPayloadFraction:        leafPayloadFraction,
		fileChangeCounter:          int(binary.BigEndian.Uint32(fileChangeCounter)),
		dbSizeInPages:              int(binary.BigEndian.Uint32(dbSizeInPages)),
		firstFreeListTrunkPage:     firstFreeListTrunkPage,
		totalNumberOfFreeListPages: totalNumberOfFreeListPages,
		schemaCookie:               int(binary.BigEndian.Uint32(schemaCookie)),
		schemaFormatNumber:         schemaFormatNumber,
		defaultPageCacheSize:       defaultPageCacheSize,
		largestBTreePage:           largestBTreePage,
		databaseEncoding:           databaseEncoding,
		userVersion:                userVersion,
		incrementalVacuumMode:      incrementalVacuumMode,
		applicationId:              applicationId,
		reservedForExpansion:       reservedForExpansion,
		versionValidForNumber:      int(binary.BigEndian.Uint32(versionValidForNumber)),
		sqlVersionNumber:           sqlVersionNumber,
	}
}

func assembleDbHeader(header DbHeader) []byte {
	data := header.headerString
	data = append(data, intToBinary(header.databasePageSize, 2)...)
	data = append(data, header.databaseFileWriteVersion...)
	data = append(data, header.databaseFileReadVersion...)
	data = append(data, header.reservedBytesSpace...)
	data = append(data, header.maxEmbeddedPayloadFraction...)
	data = append(data, header.minEmbeddedPayloadFraction...)
	data = append(data, header.leafPayloadFraction...)
	data = append(data, intToBinary(header.fileChangeCounter, 4)...)
	data = append(data, intToBinary(header.dbSizeInPages, 4)...)
	data = append(data, header.firstFreeListTrunkPage...)
	data = append(data, header.totalNumberOfFreeListPages...)
	data = append(data, intToBinary(header.schemaCookie, 4)...)
	data = append(data, header.schemaFormatNumber...)
	data = append(data, header.defaultPageCacheSize...)
	data = append(data, header.largestBTreePage...)
	data = append(data, header.databaseEncoding...)
	data = append(data, header.userVersion...)
	data = append(data, header.incrementalVacuumMode...)
	data = append(data, header.applicationId...)
	data = append(data, header.reservedForExpansion...)
	data = append(data, intToBinary(header.versionValidForNumber, 4)...)
	data = append(data, header.sqlVersionNumber...)

	return data
}

func parseReadPage(data []byte, dbPage int, fileInfo fs.FileInfo) PageParsed {
	fmt.Println("parse read page execution time?")
	fmt.Println(dbPage)
	// fmt.Println(data)
	if dbPage == 0 && len(data) == 0 {
		return PageParsed{
			dbHeader:             header(),
			dbHeaderSize:         100,
			numberofCells:        0,
			startCellContentArea: PageSize,
			cellArea:             []byte{},
			pointers:             []byte{},
			latesRow: &LastPageParseLatestRow{
				rowId:   0,
				data:    []byte{},
				columns: []PageParseColumn{},
			},
			dbInfo: DbInfo{
				pageNumber: 0,
			},
		}
	}

	if len(data) != PageSize {
		panic("invalid page size, expected" + strconv.Itoa(PageSize))
	}
	dataToParse := data
	var dbHeader DbHeader
	if dbPage == 0 {
		//Skip header for now
		dataToParse = dataToParse[100:]
		dbHeader = parseDbHeader(data[:100])
	}

	btreeType := dataToParse[0]
	isPointerValue := false
	switch BtreeType(btreeType) {
	case TableBtreeInteriorCell, IndexBtreeInteriorCell:
		isPointerValue = true

	}
	freebBlocks := dataToParse[1:3]
	if freebBlocks[0] != 0 {
		fmt.Println(freebBlocks)
		panic("implement free blocks more than 0 cell")
	}
	freeBlocksInt := int(freebBlocks[1])
	numberofCells := dataToParse[3:5]
	if numberofCells[0] != 0 {
		fmt.Println(numberofCells)
		panic("implement number of cell more than 0 cell")
	}

	fmt.Println("checkpoint 1")

	numberofCellsInt := int(numberofCells[1])
	startCellContentArea := dataToParse[5:7]
	// if startCellContentArea[0] != 0 {
	// 	fmt.Println("start cell content area")
	// 	fmt.Println(startCellContentArea)
	// 	panic("implement startCellContentArea more than 0")
	// }
	startCellContentAreaInt := binary.BigEndian.Uint16(startCellContentArea)
	startCellContentAreaBigEndian := binary.BigEndian.Uint16(startCellContentArea)
	fragmenetedArea := dataToParse[7]
	var rightMostPointer []byte
	if isPointerValue {
		rightMostPointer = dataToParse[8:12]
		dataToParse = dataToParse[12:]
	} else {
		dataToParse = dataToParse[8:]
	}
	fmt.Println("checkpoint 2")

	var pointers []byte

	for {
		pointer := dataToParse[:2]
		if pointer[0] == 0 && pointer[1] == 0 {
			break
		}
		dataToParse = dataToParse[2:]
		pointers = append(pointers, pointer...)
	}
	if len(data) < int(startCellContentAreaBigEndian) {
		panic("data length is lesser than start of cell content area")
	}

	cellAreaContent := data[startCellContentAreaBigEndian:]
	latestRowHeaders := []byte{}
	latestRowValues := []byte{}
	var latestRow LastPageParseLatestRow

	if len(cellAreaContent) > 0 {
		latestRowLength := int(cellAreaContent[0]) + 2
		fmt.Println("latestes row length?")
		//TOOD:  wait what why 9??? no idea, was it hardcoded?? i guess
		var latestRowLengthArr []byte
		for i := 0; i < latestRowLength; i++ {
			latestRowLengthArr = append(latestRowLengthArr, cellAreaContent[i])
			if cellAreaContent[i] < 127 {
				break
			}
		}
		fmt.Println("latestes row length? after")
		if len(latestRowLengthArr) > 1 {
			panic("Need to be handled later")
		}

		// latestRowLength := int(latestRowLengthArr[0]) + 2 // 1 bytes for length, 1 bytes for row id
		if len(cellAreaContent) < int(latestRowLength) {
			panic("cellAreaContent length is lesser than start of cell content area, row length%")
		}
		fmt.Println("checkpoint 4")
		latestRowRaw := cellAreaContent[:latestRowLength]
		fmt.Println("lates row")
		fmt.Println(latestRowRaw)
		latestRowId := latestRowRaw[1]
		latestRowheaderLength := latestRowRaw[2]
		latestRowHeaders = latestRowRaw[3 : 3-1+int(latestRowheaderLength)] // 3 - 1 (-1 because of header length contains itself)
		latestRowValues = latestRowRaw[3-1+int(latestRowheaderLength):]
		fmt.Println("checkpoint 5")
		fmt.Println(latestRowHeaders)
		fmt.Println(latestRowValues)
		latestRowColumns := parseDbPageColumn(latestRowHeaders, latestRowValues)
		fmt.Println("checkpoint 6")
		latestRow = LastPageParseLatestRow{
			rowId:   int(latestRowId),
			data:    latestRowRaw,
			columns: latestRowColumns,
		}
	}

	return PageParsed{
		dbHeader:             dbHeader,
		btreeType:            int(btreeType),
		numberofCells:        numberofCellsInt,
		startCellContentArea: int(startCellContentAreaInt),
		rightMostpointer:     rightMostPointer,
		cellArea:             cellAreaContent,
		framgenetedArea:      int(fragmenetedArea),
		freeBlock:            int(freeBlocksInt),
		pointers:             pointers,
		latesRow:             &latestRow,
		dbInfo: DbInfo{
			pageNumber: int(fileInfo.Size() / int64(PageSize)),
		},
	}
}

func assembleDbPage(page PageParsed) []byte {
	data := []byte{}
	if page.dbHeaderSize > 0 {
		data = append(data, assembleDbHeader(page.dbHeader)...)
	}
	data = append(data, byte(page.btreeType))
	data = append(data, intToBinary(page.freeBlock, 2)...)
	data = append(data, intToBinary(page.numberofCells, 2)...)
	data = append(data, intToBinary(page.startCellContentArea, 2)...)
	data = append(data, byte(page.framgenetedArea))
	if len(page.rightMostpointer) > 0 {
		data = append(data, page.rightMostpointer...)
	}
	data = append(data, page.pointers...)

	zerosLen := PageSize - len(data) - len(page.cellArea)
	data = append(data, make([]byte, zerosLen)...)
	data = append(data, page.cellArea...)

	return data

}

func BtreeHeaderSchema(btreeType BtreeType, cell CreateCell, parsedData *PageParsed) []byte {
	//This should be read from the page
	currentNumberOfCell := 0
	currentCellStart := PageSize
	pointers := []byte{}
	if parsedData != nil {
		currentNumberOfCell = parsedData.numberofCells
		currentCellStart = parsedData.startCellContentArea
		pointers = parsedData.pointers
	}

	var lastPointer int = PageSize

	if len(pointers) > 0 {
		lastPointer = int(binary.BigEndian.Uint16(parsedData.pointers[len(parsedData.pointers)-2 : len(parsedData.pointers)]))
	}

	if cell.dataLength > 0 {
		currentNumberOfCell += 1
		newCellStartPosition := lastPointer - cell.dataLength - 2
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

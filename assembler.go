package main

import "fmt"

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
	cellArea := []byte{}
	if len(page.cellAreaParsed) > 0 {

		for _, v := range page.cellAreaParsed {
			cellArea = append(cellArea, v...)
		}
	} else {
		cellArea = page.cellArea
	}

	// if PageSize-len(cellArea) != page.startCellContentArea {
	// 	fmt.Printf("cell area should be equal %v, %v \n", PageSize-len(cellArea), page.startCellContentArea)
	// 	panic("cell area start pointer, and data not equal")
	// }

	zerosLen := PageSize - len(data) - len(cellArea)

	if zerosLen < 0 {
		fmt.Println("page number", page.pageNumber)
		// fmt.Println("cell area parsed??")
		// fmt.Println(page.cellAreaParsed)
		fmt.Printf("\n data length: %v, cell area length: %v", len(data), len(cellArea))
		panic("zeros length should never be less than 0")
	}

	data = append(data, make([]byte, zerosLen)...)
	data = append(data, cellArea...)

	return data

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

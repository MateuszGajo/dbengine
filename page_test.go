package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"sync"
	"testing"
	"time"
)

type MockFileInfo struct {
	NameVal    string
	SizeVal    int64
	ModeVal    fs.FileMode
	ModTimeVal time.Time
	IsDirVal   bool
}

func (m MockFileInfo) Name() string       { return m.NameVal }
func (m MockFileInfo) Size() int64        { return m.SizeVal }
func (m MockFileInfo) Mode() fs.FileMode  { return m.ModeVal }
func (m MockFileInfo) ModTime() time.Time { return m.ModTimeVal }
func (m MockFileInfo) IsDir() bool        { return m.IsDirVal }
func (m MockFileInfo) Sys() any           { return nil }

func TestReadSharedLock(t *testing.T) {
	sleepTimeMs := 15

	reader := PageReader{
		readInternal: func(pageNumber int) ([]byte, os.FileInfo) {
			time.Sleep(time.Duration(sleepTimeMs) * time.Millisecond)
			return []byte("mock data"), MockFileInfo{}
		},
	}

	var wg sync.WaitGroup
	readTime := map[int]int64{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		time.Sleep(1 * time.Millisecond)

		go func(id int) {
			defer wg.Done()
			reader.readDbPage(1)
			readTime[id] = time.Now().UnixMilli()
		}(i)
	}

	wg.Wait()

	fmt.Printf("%+v", readTime)
	maxVal := math.Max(float64(readTime[0]), float64(readTime[1]))
	minVal := math.Min(float64(readTime[1]), float64(readTime[0]))

	if maxVal-minVal > float64(sleepTimeMs-1) {
		t.Errorf("Execution should be conncurent, the second one should wait less than 15 miliscond (-1 ms for error margin) to finish after first one, we waited: %f", maxVal-minVal)
	}
}

func TestWriteExclusiveLockConcurrentWrite(t *testing.T) {
	sleepTimeMs := 15
	writer := WriterStruct{
		writeToFileRaw: func(data []byte, page int) error {
			fmt.Println("call raw function")
			fmt.Println("call raw function")
			fmt.Println("call raw function")
			time.Sleep(time.Duration(sleepTimeMs) * time.Millisecond)

			return nil
		},
	}

	var wg sync.WaitGroup
	readTime := map[int]int64{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		time.Sleep(1 * time.Millisecond)

		go func(id int) {
			defer wg.Done()
			writer.writeToFile([]byte{}, 1, PageParsed{}, "conId")
			readTime[id] = time.Now().UnixMilli()
		}(i)
	}

	wg.Wait()

	maxVal := math.Max(float64(readTime[0]), float64(readTime[1]))
	minVal := math.Min(float64(readTime[1]), float64(readTime[0]))

	if maxVal-minVal < float64(sleepTimeMs) {
		t.Errorf("Exclusive lock can't write concurrently, expected to wait at least 15, instead we waited: %f", maxVal-minVal)
	}

	fmt.Println("we waited")
	fmt.Println(maxVal - minVal)

}

func TestWriteExclusiveLockConcurrentWriteAndRead(t *testing.T) {
	sleepTimeMs := 15
	writer := WriterStruct{
		writeToFileRaw: func(data []byte, page int) error {
			fmt.Println("write to file raw execute")
			time.Sleep(time.Duration(sleepTimeMs) * time.Millisecond)

			return nil
		},
	}

	reader := PageReader{
		readInternal: func(pageNumber int) ([]byte, os.FileInfo) {
			time.Sleep(time.Duration(sleepTimeMs) * time.Millisecond)
			return []byte("mock data"), MockFileInfo{}
		},
	}

	var wg sync.WaitGroup
	readTime := map[int]int64{}
	wg.Add(1)
	wg.Add(1)
	go func(id int) {

		defer wg.Done()
		writer.writeToFile([]byte{}, 1, PageParsed{}, "conId")
		readTime[id] = time.Now().UnixMilli()
	}(1)

	go func(id int) {
		time.Sleep(1 * time.Millisecond)
		defer wg.Done()
		reader.readDbPage(1)
		readTime[id] = time.Now().UnixMilli()
	}(2)

	wg.Wait()

	maxVal := math.Max(float64(readTime[0]), float64(readTime[1]))
	minVal := math.Min(float64(readTime[1]), float64(readTime[0]))

	if maxVal-minVal < float64(sleepTimeMs) {
		t.Errorf("Exclusive lock can't write and read concurrently, expected to wait at least 15, instead we waited: %v", maxVal-minVal)
	}

	fmt.Println("we waited")
	fmt.Println(maxVal - minVal)

}

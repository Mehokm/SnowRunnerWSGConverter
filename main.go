package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

const containerPrefix = "container."
const chunkLen = 160
const metadataLen = 8

type WsgFileInfo struct {
	realFilename []byte
	guid         []byte
}

func main() {
	// change this to flags maybe
	if len(os.Args) < 3 {
		fmt.Println("you must provide a source and destination directories as arguments")
		return
	}

	// add isDir checks for safety
	source := os.Args[1]

	sourceIsDir, err := isDir(source)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !sourceIsDir {
		fmt.Println("source is not a directory.  Source must be a directory.")
		return
	}

	dest := os.Args[2]

	destIsDir, err := isDir(dest)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !destIsDir {
		fmt.Println("destination is not a directory.  Destination must be a directory.")
		return
	}

	files, err := ioutil.ReadDir(source)

	if err != nil {
		fmt.Println(err)
		return
	}

	containerFilename := findContainerFilename(files)

	if containerFilename == "" {
		fmt.Println("cannot find container file.")
		return
	}

	data, err := ioutil.ReadFile(source + "/" + containerFilename)

	if err != nil {
		fmt.Println(err)
		return
	}

	numFiles := binary.LittleEndian.Uint16(data[4:metadataLen])

	// create WaitGroup for managing concurrent file operations
	wg := sync.WaitGroup{}

	i := 0
	start := metadataLen
	end := start + chunkLen
	for i < int(numFiles) {
		// chunk is 160 bytes starting after 8 bytes of metadata
		chunk := data[start:end]

		// extract the real filename which is at the start of the chunk to 00 padding
		realFilename := extractRealFilename(chunk)

		// extract encoded WSG filename at the end of the chunk by reading last 16 bytes
		encodedFilename := extractEncodedFilename(chunk)

		// convert the encoded filename into a GUID which represents the actual WSG filename
		decodedGuid := guid(encodedFilename)

		// create WsgFileInfo object in order to encapsulate logic
		fileInfo := WsgFileInfo{
			realFilename: realFilename,
			guid:         decodedGuid,
		}

		wg.Add(1)

		// concurrently convert WSG save file to proper save file
		go convertFile(&wg, fileInfo, source, dest)

		// update for following chunks
		start = end
		end = start + chunkLen

		if end >= len(data) {
			end = len(data)
		}

		i++
	}

	wg.Wait()
}

func convertFile(wg *sync.WaitGroup, fileInfo WsgFileInfo, source, dest string) {
	guidFilepath := source + "/" + fmt.Sprintf("%X", fileInfo.guid)
	realFilepath := dest + "/" + string(fileInfo.realFilename)

	guidFile, err := os.Open(guidFilepath)

	if err != nil {
		fmt.Println(err)

		wg.Done()

		return
	}

	realFile, err := os.Create(realFilepath)

	if err != nil {
		fmt.Println(err)

		wg.Done()

		return
	}

	_, err = io.Copy(realFile, guidFile)

	if err != nil {
		fmt.Println(err)

		wg.Done()

		return
	}

	guidFile.Close()
	realFile.Close()

	fmt.Println("copied: " + guidFilepath + " -> " + realFilepath)

	wg.Done()
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return info.IsDir(), err
}

func findContainerFilename(files []fs.FileInfo) string {
	for _, f := range files {
		if strings.Contains(f.Name(), containerPrefix) {
			return f.Name()
		}
	}

	return ""
}

func extractRealFilename(chunk []byte) []byte {
	var extracted []byte
	zeroCount := 0

	i := 0
	n := len(chunk) - 1
	for i < n {
		// filenames seem to be separated by one zero between bytes, so if we get more than one zero in a row break out
		if zeroCount >= 2 {
			break
		}

		if chunk[i] == 0 {
			zeroCount++
		} else {
			if zeroCount > 0 {
				zeroCount--
			}

			extracted = append(extracted, chunk[i])
		}

		i++
	}

	return extracted
}

func extractEncodedFilename(chunk []byte) []byte {
	var extracted []byte

	// read last 16 bytes for encoded GUID
	i := len(chunk) - 16
	n := len(chunk)

	for i < n {
		extracted = append(extracted, chunk[i])

		i++
	}

	return extracted
}

func guid(b []byte) []byte {
	g := make([]byte, 0)

	g = append(g, b[3], b[2], b[1], b[0])
	g = append(g, b[5], b[4])
	g = append(g, b[7], b[6])
	g = append(g, b[8:10]...)
	g = append(g, b[10:]...)

	return g
}

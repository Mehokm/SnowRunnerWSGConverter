package wsg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type WsgFileInfo struct {
	realFilename []byte
	guid         []byte
}

const containerPrefix = "container."
const chunkLen = 160
const metadataLen = 8

func Convert(sourceDir, destDir, ext string) error {
	return run(sourceDir, destDir, ext, false)
}

func ListMapping(sourceDir, destDir string) error {
	return run(sourceDir, destDir, "", true)
}

func run(sourceDir, destDir, ext string, listOnly bool) error {
	files, err := ioutil.ReadDir(sourceDir)

	if err != nil {
		return err
	}

	// find the container.XXX file to extract the proper filenames for conversion
	containerFilename := findContainerFilename(files)

	if containerFilename == "" {
		return errors.New("wsg.go: cannot find container file")
	} else {
		fmt.Println("Found container file: " + containerFilename)
	}

	data, err := ioutil.ReadFile(sourceDir + "/" + containerFilename)

	if err != nil {
		return err
	}

	// parse number of files from file metadata
	numFiles := binary.LittleEndian.Uint16(data[4:metadataLen])

	fileInfoCh := make(chan WsgFileInfo)

	wg := sync.WaitGroup{}

	go func() {
		if listOnly {
			fmt.Println("WSG save file mapping:")
			fmt.Println()

			for fileInfo := range fileInfoCh {
				fmt.Println(string(fileInfo.realFilename) + " -> " + fmt.Sprintf("%X", fileInfo.guid))
			}
		} else {
			fmt.Println("Converting WSG save files..")
			fmt.Println()

			for fileInfo := range fileInfoCh {
				wg.Add(1)

				go convertFile(&wg, fileInfo, sourceDir, destDir, ext)
			}

		}
	}()

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

		// put files info on channel for concurrent processing
		fileInfoCh <- fileInfo

		// update for next chunk
		start = end
		end = start + chunkLen

		if end >= len(data) {
			end = len(data)
		}

		i++
	}

	close(fileInfoCh)

	wg.Wait()

	return nil
}

func convertFile(wg *sync.WaitGroup, fileInfo WsgFileInfo, sourceDir, destDir, ext string) {
	defer wg.Done()

	guidFilepath := sourceDir + "/" + fmt.Sprintf("%X", fileInfo.guid)
	realFilepath := destDir + "/" + string(fileInfo.realFilename) + "." + ext

	guidFile, err := os.Open(guidFilepath)

	if err != nil {
		fmt.Println(err)
		return
	}

	realFile, err := os.Create(realFilepath)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = io.Copy(realFile, guidFile)

	if err != nil {
		fmt.Println(err)
		return
	}

	guidFile.Close()
	realFile.Close()

	fmt.Println("Copied: " + guidFilepath + " -> " + realFilepath)
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

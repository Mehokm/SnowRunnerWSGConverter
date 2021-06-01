package main

import (
	"SnowRunnerWSGConverter/internal/wsg"
	"flag"
	"fmt"
	"os"
)

func main() {
	srcFlag := flag.String("src", ".", "The source save directory containing WSG save files")
	destFlag := flag.String("dest", ".", "The destination directory for converted save files")
	extFlag := flag.String("ext", "cfg", "File extension to save all converted files with")
	listOnlyFlag := flag.Bool("list-only", false, "List WSG save file mappings only")

	flag.Parse()

	source := *srcFlag
	sourceIsDir, err := isDir(source)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !sourceIsDir {
		fmt.Println("source is not a directory.  Source must be a directory")
		return
	}

	dest := *destFlag

	destIsDir, err := isDir(dest)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !destIsDir {
		fmt.Println("destination is not a directory.  Destination must be a directory")
		return
	}

	fmt.Println("Using '" + source + "' as source directory")

	if !*listOnlyFlag {
		fmt.Println("Using '" + dest + "' as destination directory")
	}

	fmt.Println()

	if *listOnlyFlag {
		err = wsg.ListMapping(source, dest)
	} else {
		err = wsg.Convert(source, dest, *extFlag)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return info.IsDir(), err
}

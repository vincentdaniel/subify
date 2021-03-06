package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/matcornic/subify/common/config"
	"os"
)

//GetHashOfVideo gets the hash used by SubDb to identify a video. Absolutely needed either to download or upload subtitles.
//The hash is composed by taking the first and the last 64kb of the video file, putting all together and generating a md5 of the resulting data (128kb).
func GetHashOfVideo(filename string) string {
	readsize := 64 * 1024 // 64kb

	// Open Video
	file, err := os.Open(filename)
	if err != nil {
		if config.Verbose {
			fmt.Fprintln(os.Stderr, err)
		}
		Exit("Can't open file " + filename)
	}
	defer file.Close()

	// Get stats of file
	fi, err := file.Stat()
	if err != nil {
		if config.Verbose {
			fmt.Fprintln(os.Stderr, err)
		}
		Exit("Can't open file " + filename)
	}

	// Fill a buffer with first bytes of file
	bufB := make([]byte, readsize)
	_, err = file.Read(bufB)
	if err != nil {
		if config.Verbose {
			fmt.Fprintln(os.Stderr, err)
		}
		Exit("Can't read content of file " + filename)
	}

	//Fill a buffer with last bytes of file
	bufE := make([]byte, readsize)
	n, err := file.ReadAt(bufE, fi.Size()-int64(len(bufE)))
	if err != nil {
		if config.Verbose {
			fmt.Fprintln(os.Stderr, err)
		}
		Exit("Can't read content of file. File is probably too small : " + filename)
	}
	bufE = bufE[:n]

	// Generates MD5 of both bytes chain
	bufB = append(bufB, bufE...)
	hash := fmt.Sprintf("%x", md5.Sum(bufB))

	if config.Verbose {
		fmt.Println("Hash of video is " + hash)
	}

	return hash
}

// Exit func displays an error message on stderr and exit 1
func Exit(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
	if !config.Verbose {
		fmt.Fprintln(os.Stderr, "Run subify with --verbose option to get more information about the error")
	}
	os.Exit(1)
}

// Error func displays an error message on stderr
func Error(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
}

package utils

import (
	"github.com/kmou424/filewasher/consts"
	"github.com/kmou424/filewasher/types"
	"github.com/ohanakogo/exceptiongo"
	"io"
	"math/rand"
	"os"
	"time"
)

var (
	fileWasherFooter = []byte(consts.FileWasherFooter)
	randomBytesLen   = consts.RandomFooterBytesLen

	fullFooterLen = len(fileWasherFooter) + randomBytesLen
)

func GenFileFooter() []byte {
	if len(fileWasherFooter) == 0 {
		exceptiongo.QuickThrowMsg[types.DefinitionErrorException](`constant "FileWasherFooter" cannot be empty. please go to consts/filewasher.go to make changes`)
	}

	resultBytes := make([]byte, len(fileWasherFooter))
	copy(resultBytes, fileWasherFooter)

	time.Sleep(time.Duration(1) * time.Nanosecond)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomBytesLen; i++ {
		randomByte := byte(random.Intn(256))
		resultBytes = append(resultBytes, randomByte)
	}

	return resultBytes
}

func readFileWasherFooter(filename string) []byte {
	file, err := os.Open(filename)
	exceptiongo.QuickThrow[types.IOException](err)
	defer func(file *os.File) {
		err := file.Close()
		exceptiongo.QuickThrow[types.IOException](err)
	}(file)

	stat, err := file.Stat()
	exceptiongo.QuickThrow[types.IOException](err)

	fileSize := stat.Size()
	if fileSize < int64(fullFooterLen) {
		return nil
	}

	offset := fileSize - int64(fullFooterLen)
	buffer := make([]byte, fullFooterLen)

	_, err = file.Seek(offset, io.SeekStart)
	exceptiongo.QuickThrow[types.IOException](err)

	_, err = file.Read(buffer)
	exceptiongo.QuickThrow[types.IOException](err)

	return buffer
}

func CheckFileWasher(filename string) bool {
	buffer := readFileWasherFooter(filename)
	if buffer == nil {
		return false
	}

	for i := 0; i < len(fileWasherFooter); i++ {
		if buffer[i] != fileWasherFooter[i] {
			return false
		}
	}

	return true
}

package buffer

import (
	"fmt"
	"os"
)

func ReadFile(fileName string) []byte {
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	bytesRead, err := file.Read(buffer)
	bytesRead = bytesRead

	return buffer
}

package buffer

import (
	"bytes"
	"fmt"
	"os"
)

func ReadFile(fileName string) []byte {
	file, err := os.Open(fileName)

	if os.IsNotExist(err) {
		return []byte{}
	}

	if err != nil {
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

func WriteFile(fileName string, bufferRows [][]byte) {
	buffer := bytes.Join(bufferRows, []byte{'\n'})

	file, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	bytesWritten, err := file.Write(buffer)

	if err != nil {
		panic(err)
	}

	fmt.Println("Bytes written:", bytesWritten)
}

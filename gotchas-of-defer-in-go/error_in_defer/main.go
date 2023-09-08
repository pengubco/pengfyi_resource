package main

import (
	"fmt"
	"os"
)

func main() {
	filename := "./a.txt"
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("file %s opened successfully\n", filename)
	err = writeHelloWorldFixed(f)
	if err != nil {
		fmt.Printf("write failed %v\n", err)
	} else {
		fmt.Println(`successfully write "hello world" to file`)
	}
}

func writeHelloWorld(f *os.File) (err error) {
	defer func() {
		err = f.Close()
	}()

	_, err = f.WriteString("hello world!\n")
	return err
}

func writeHelloWorldFixed(f *os.File) (err error) {
	defer func() {
		errCloseFile := f.Close()
		if err == nil && errCloseFile != nil {
			err = errCloseFile
		} else if err != nil && errCloseFile != nil {
			fmt.Printf("failed to close file %v", errCloseFile)
		}
	}()

	_, err = f.WriteString("hello world!\n")
	return err
}

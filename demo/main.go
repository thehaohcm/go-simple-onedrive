package main

import (
	"fmt"
	"os"

	"github.com/thehaohcm/go-simple-onedrive/upload"
)

var (
	oauthStateString = "12345"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 1 {
		fmt.Println("Something is wrong, please add a file path as an argument...")
		os.Exit(1)
	}
	filePath := args[0]
	upload.UploadFile(filePath)

	//get list item of root
	upload.GetItemsByPath("/Folder")
}

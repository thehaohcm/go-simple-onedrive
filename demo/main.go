package main

import (
	"context"
	"fmt"
	"os"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/upload"

	"github.com/goh-chunlin/go-onedrive/onedrive"
	"golang.org/x/oauth2"
)

var (
	client           *onedrive.Client
	oauthStateString = "12345"
)

func getInstance() (context.Context, *onedrive.Client) {
	oneDriveClient := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.SavedToken.AccessToken},
	)

	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oneDriveClient)
	client = onedrive.NewClient(tc)
	return ctx, client
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 1 {
		fmt.Println("Something is wrong, please add a file path as an argument...")
		os.Exit(1)
	}
	filePath := args[0]
	upload.UploadFile(filePath)
	ctx, client := getInstance()

	drives, err := client.Drives.List(ctx)
	if err != nil {
		fmt.Println(err)
	}

	for _, drive := range drives.Drives {
		fmt.Printf("Results: %v\n", drive.Owner.User.DisplayName)
	}

	//get list item of root
	driveItems, err := client.DriveItems.List(ctx, "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Items of root folder: ")
		for _, driveItem := range driveItems.DriveItems {
			fmt.Printf(" - %v ", driveItem.Name)
			// if driveItem.Folder != nil {
			// 	fmt.Printf("- Folder\n")
			// }
		}
	}
}

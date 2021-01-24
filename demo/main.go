package main

import (
	"context"
	"fmt"
	"os"

	"github.com/thehaohcm/go-simple-onedrive/token"
	"github.com/thehaohcm/go-simple-onedrive/upload"

	"github.com/goh-chunlin/go-onedrive/onedrive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	client    *onedrive.Client
	oauthConf = &oauth2.Config{
		ClientID:     "fbe1ffd1-93ca-4aaf-a121-656849b2cfd3",
		ClientSecret: ".4.PpG17mF_TyQ3~2wWRwZFTbOU_5aq3Gf",
		RedirectURL:  "http://localhost",
		Scopes:       []string{"Files.ReadWrite.All", "Sites.ReadWrite.All", "openid", "User.ReadBasic.All", "User.ReadWrite", "profile", "email"},
		Endpoint:     microsoft.AzureADEndpoint(token.Tenant),
	}
	oauthStateString = "12345"
)

func getInstance() (context.Context, *onedrive.Client) {
	// token.RefreshToken()
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.SavedToken.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = onedrive.NewClient(tc)

	return ctx, client
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 1 {
		fmt.Println("Something is wrong, please add a file path as an argument...")
	}
	filePath := args[1]
	upload.UploadFile(filePath)
	ctx, client := getInstance()

	drives, err := client.Drives.List(ctx)
	if err != nil {
		panic(err)
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
			fmt.Printf(" -%v ", driveItem.Name)
			// if driveItem.Folder != nil {
			// 	fmt.Printf("- Folder\n")
			// }
		}
	}
}

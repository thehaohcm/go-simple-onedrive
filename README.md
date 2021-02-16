# go-simple-onedrive
A simple Onedrive application in CLI writen by Golang

Copyright © by Nguyen The Hao 2021

Author: Nguyen The Hao (thehaohcm)

Git repository: https://github.com/thehaohcm/go-simple-ondrive

This is my simple OneDrive Application with CLI written by Golang. Currently, it just has some below basic features, I will try to implement another when I have a free time:

* [x] Upload file (for a normal or even large files)
* [x] Share an uploaded file (get link file for downloading)
* [x] List all children items of a specific folder
* [x] Create a new folder
* [x] Delete one item
* [x] Get Info of one item
* [x] Get download link of item
* [ ] Download item directly in application (pending)
* [x] Move item to another path
* [x] Copy item
* [ ] Handle Errors

Detail instruction url: https://docs.microsoft.com/en-us/graph/api/resources/onedrive?view=graph-rest-1.0

You can follow all these below steps to run an application: 

1st step: clone this repo into your local pc (in branch master)

2nd step: open the application's folder by an IDE (ex: Visual Code), then run a command "**go get**" in terminal to install all neccessary packages (please make sure your pc has been installed GO lang and GO GET)

3rd step: open the file **config.yaml** and replace all neccessary variables inside with correct info of OneDrive API register's info

4th step: in the terminal openning the project's path, type a command to go inside the **demo** folder: **cd demo**

5th step: type this command along with a speific path to upload a file into your OneDrive account path: **go run main.go *[FILE_PATH]***
  ex: **go run main.go 

Notice: the applicaiton can upload with a large file (everything works properly when I tested by uploading a file whose size is a several GB)

This is a screenshot when the uploading file successfully completed:
![alt text](https://github.com/thehaohcm/go-simple-onedrive/blob/master/screenshot/screenshot-demo.png?raw=true)

## How to use/integrate the simple-onedrive library into your golang project

*Comming soon*
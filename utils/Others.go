package utils

import (
	"fmt"
)

func GetReadableFileCapacity(fileSize int64) string {
	suffix := []string{"Byte", "KB", "MB", "GB", "TB", "PB", "EB"}
	index := 0
	for ; fileSize > 1024; fileSize = fileSize / 1024 {
		index++
	}
	return fmt.Sprintf("%d %s", fileSize, suffix[index])
}

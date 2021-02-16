package models

type UploadingManager struct {
	BlockFileNumber int
	BlockFileFlags  []bool
}

func (u UploadingManager) Init(blockNumber int) {
	u.BlockFileNumber = blockNumber
	u.BlockFileFlags = make([]bool, blockNumber)

	var blockPositions []int
	for i := 0; i < blockNumber; i++ {
		blockPositions = append(blockPositions, i)
	}
	u.SetValuesForBlockFileArray(blockPositions, false)
}

func (u UploadingManager) SetValuesForBlockFileArray(blockPositions []int, value bool) {
	for _, position := range blockPositions {
		u.BlockFileFlags[position] = value
	}
}

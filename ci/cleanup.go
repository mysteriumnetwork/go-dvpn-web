package ci

import (
	"fmt"
	"os"
)

// Cleanup removes the assets
func Cleanup() error {
	fmt.Println("cleaning up...")
	toClean := []string{
		tempDir, assetDir, distAssetName,
	}

	for _, v := range toClean {
		err := os.RemoveAll(v)
		if err != nil {
			fmt.Println("could not remove", v)
		}
	}
	fmt.Println("cleanup done")
	return nil
}

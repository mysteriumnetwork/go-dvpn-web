package ci

import (
	"fmt"
	"os"
)

// FixDirectory un-nest assets
func FixDirectory() error {
	fmt.Println("renaming directory assets/build -> assets")
	err := os.Rename(tempDir+"/build", assetDir)
	if err != nil {
		return err
	}
	fmt.Println("directory renamed")
	return nil
}

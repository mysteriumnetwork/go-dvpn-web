package ci

import (
	"fmt"
	"github.com/mholt/archiver/v3"
)

// ExtractAssets extracts the asset
func ExtractAssets() error {
	z := archiver.TarGz{
		Tar: &archiver.Tar{
			OverwriteExisting: true,
			MkdirAll:          true,
		},
	}
	fmt.Println("extracting archive", distAssetName)
	err := z.Unarchive(distAssetName, tempDir)
	if err != nil {
		return err
	}
	fmt.Println("archive", distAssetName, "extracted")
	return nil
}

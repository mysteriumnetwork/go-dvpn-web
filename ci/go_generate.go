package ci

import "github.com/magefile/mage/sh"

// Generate generates an updated assets_vfsdata.go file
func GoGenerate() error {
	return sh.RunV("go", "generate")
}

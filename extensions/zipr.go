//  zipr.go -- Zip file reader extension for Goaldi

package extensions

import (
	"archive/zip"
	"goaldi/runtime"
)

func init() {
	runtime.GoLib(zip.OpenReader, "zipreader", "name", "open a Zip file")
}

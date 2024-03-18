package FileEndpoints

import "path/filepath"

func CanFileBeStored(filename string) bool {
	isM3U8 := filepath.Ext(filename) != ".m3u8"
	isTs := filepath.Ext(filename) != ".ts"
	isM4a := filepath.Ext(filename) != ".m4a"

	return isM3U8 && isTs && isM4a
}

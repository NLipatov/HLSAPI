package ApplicationLayerErrors

type FileCantBeStored struct {
}

func (e FileCantBeStored) Error() string {
	return "File cannot be stored"
}

type FileCantBeConvertedToM3U8 struct {
}

func (e FileCantBeConvertedToM3U8) Error() string {
	return "File cannot be converted to m3u8"
}

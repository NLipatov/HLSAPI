package ApplicationLayerErrors

type FileCantBeStored struct {
}

func (e FileCantBeStored) Error() string {
	return "File cannot be stored"
}

package WipeModes

type WipeMode byte

const (
	UNSET WipeMode = iota
	REMOVE_KEY_FILES
	REMOVE_ALL_FILES
)

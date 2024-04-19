package CleanupType

type CleanupMode byte

const (
	UNSET CleanupMode = iota
	REMOVE_KEY_FILES
	REMOVE_ALL_FILES
)

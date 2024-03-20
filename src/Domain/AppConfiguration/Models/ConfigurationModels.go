package ConfigurationModels

type ConfigurationRoot struct {
	Server        ServerConfiguration        `json:"Server"`
	Storage       StorageConfiguration       `json:"Storage"`
	StorageDaemon StorageDaemonConfiguration `json:"StorageDaemon"`
}

type ServerConfiguration struct {
	Port int `json:"Port"`
}

type StorageDaemonConfiguration struct {
	ShouldRun                    bool `json:"ShouldRun"`
	StorageLimitMinutes          int  `json:"StorageLimitMinutes"`
	StorageChecksIntervalMinutes int  `json:"StorageChecksIntervalMinutes"`
	EnableLogging                bool `json:"EnableLogging"`
}

type StorageConfiguration struct {
	MaxFileSizeMb     int    `json:"MaxFileSizeMb"`
	StorageFolderPath string `json:"StorageFolderPath"`
}

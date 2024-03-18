package ConfigurationModels

type ConfigurationRoot struct {
	Server   ServerConfiguration   `json:"Server"`
	Storage  StorageConfiguration  `json:"Storage"`
	Sentinel SentinelConfiguration `json:"Sentinel"`
}

type ServerConfiguration struct {
	Port int `json:"Port"`
}

type SentinelConfiguration struct {
	ShouldRun                    bool `json:"ShouldRun"`
	StorageLimitMinutes          int  `json:"StorageLimitMinutes"`
	StorageChecksIntervalMinutes int  `json:"StorageChecksIntervalMinutes"`
}

type StorageConfiguration struct {
	MaxFileSizeMb     int    `json:"MaxFileSizeMb"`
	StorageFolderPath string `json:"StorageFolderPath"`
}

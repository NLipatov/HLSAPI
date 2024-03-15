package ConfigurationModels

type Configuration struct {
	Port                               int                                `json:"port"`
	StorageFolderPath                  string                             `json:"storageFolderPath"`
	SentinelServiceDaemonConfiguration SentinelServiceDaemonConfiguration `json:"SentinelServiceDaemonSettings"`
}

type SentinelServiceDaemonConfiguration struct {
	ShouldRun                    bool `json:"ShouldRun"`
	StorageLimitMinutes          int  `json:"StorageLimitMinutes"`
	StorageChecksIntervalMinutes int  `json:"StorageChecksIntervalMinutes"`
}

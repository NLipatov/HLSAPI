package Entities

type AppConfiguration struct {
	Server                           ServerConfiguration              `json:"Server"`
	Storage                          StorageConfiguration             `json:"Storage"`
	StorageDaemon                    StorageDaemonConfiguration       `json:"StorageDaemon"`
	InfrastructureLayerConfiguration InfrastructureLayerConfiguration `json:"InfrastructureLayerConfiguration"`
}

type ServerConfiguration struct {
	Port                   int    `json:"Port"`
	GetFileEndpointPostfix string `json:"GetFileEndpointPostfix"`
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

type InfrastructureLayerConfiguration struct {
	FFMPEGConverter FFMPEGConverter `json:"FFMPEGConverter"`
}

type FFMPEGConverter struct {
	UseLogging       bool   `json:"UseLogging"`
	ContainerAppRoot string `json:"ContainerAppRoot"`
}

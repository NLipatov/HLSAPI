package Infrastructure

import "os"

type EnvironmentManager struct{}

func (EnvironmentManager) GetAppRootPath() string {
	return os.Getenv("APP_ROOT")
}

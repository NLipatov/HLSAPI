package AppConfigurationErrors

type EnvConfigurationUpdateError struct {
	InnerError error
}

func (EnvConfigurationUpdateError) Error() string {
	return "Failed to update configuration with environment variables"
}

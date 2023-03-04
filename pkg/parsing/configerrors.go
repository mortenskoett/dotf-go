package parsing

type ConfigKeyNotFoundError struct {
	message string
}

func (e *ConfigKeyNotFoundError) Error() string {
	return e.message
}

type MalformedConfigurationError struct {
	message string
}

func (e *MalformedConfigurationError) Error() string {
	return e.message
}

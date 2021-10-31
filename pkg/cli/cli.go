package cli

type Command interface {
	Run([]string) error
}


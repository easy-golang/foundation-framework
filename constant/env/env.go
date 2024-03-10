package env

type Env string

const (
	DEV  Env = "dev"
	TEST Env = "test"
	PRE  Env = "pre"
	PROD Env = "prod"
)

package ncore

const namespace = "nbs-go/nucleo/ncore"
const wrappedErrorFmt = "\n  > %w"

type Core struct {
	Manifest    Manifest
	Environment Environment
	WorkDir     string
	NodeId      string
}

func (c *Core) GetEnvironmentString() string {
	switch c.Environment {
	case ProductionEnvironment:
		return "Production"
	case TestingEnvironment:
		return "Testing"
	}
	return "Development"
}

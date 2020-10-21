package driver

import (
	"github.com/ory/x/logrusx"

	"github.com/zzpu/ums/driver/configuration"
)

type BuildInfo struct {
	Version string
	Hash    string
	Time    string
}

type Driver interface {
	Logger() *logrusx.Logger
	Configuration() configuration.Provider
	Registry() Registry
}

package cliutil

import (
	"errors"
	"fmt"
	"os"

	"log/slog"

	"github.com/rancher/opni/pkg/config"
	"github.com/rancher/opni/pkg/config/meta"
	"github.com/rancher/opni/pkg/logger"
)

func LoadConfigObjectsOrDie(
	configLocation string,
	lg *slog.Logger,
) meta.ObjectList {
	if configLocation == "" {
		// find config file
		path, err := config.FindConfig()
		if err != nil {
			if errors.Is(err, config.ErrConfigNotFound) {
				wd, _ := os.Getwd()
				panic(fmt.Sprintf(`could not find a config file in ["%s","/etc/opni"], and --config was not given`, wd))
			}
			lg.Error("an error occurred while searching for a config file", logger.Err(err))
			panic(err)

		}
		lg.Info("using config file", "path", path)

		configLocation = path
	}
	objects, err := config.LoadObjectsFromFile(configLocation)
	if err != nil {
		lg.Error("failed to load config", logger.Err(err))
		panic(err)

	}
	return objects
}

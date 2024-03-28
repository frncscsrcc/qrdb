package di

import (
	"qrdb/pkg/config"
	"qrdb/pkg/database"
	"qrdb/pkg/uid"
)

type Dependencies struct {
	Initiated    bool
	Config       config.Config
	Database     *database.ConnectionPool
	UIDGenerator uid.UIDGenerator
}

var di Dependencies

func Init(configFile string) {
	config := config.GetConfigFromFile(configFile)
	database := database.NewConnectionPool(config)
	uidGenerator := uid.NewBasicUIDGenerator(8)

	di = Dependencies{
		Initiated:    true,
		Config:       config,
		Database:     database,
		UIDGenerator: uidGenerator,
	}
}

func GetDependencies() Dependencies {
	if !di.Initiated {
		panic("dependencies not initiated!")
	}

	return di
}

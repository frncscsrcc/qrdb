package di

import (
	"qrdb/qrdb/pkg/config"
	"qrdb/qrdb/pkg/database"
	"qrdb/qrdb/pkg/uid"
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

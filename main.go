package main

import (
	"qrdb/pkg/di"
	"qrdb/pkg/server"
)

func main() {
	di.Init("filename") // Todo, pass config from files and/or env
	di.GetDependencies().Database.Migrate()
	server := server.NewServer()
	server.Start()
}

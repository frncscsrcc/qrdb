package main

import (
	"qrdb/qrdb/pkg/di"
	"qrdb/qrdb/pkg/server"
)

func main() {
	di.Init("filename") // Todo, pass config from files and/or env
	di.GetDependencies().Database.Migrate()
	server := server.NewServer()
	server.Start()
}

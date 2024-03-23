package main

import (
	"qrdb/qrdb/pkg/di"
	"qrdb/qrdb/pkg/server"
)

func main() {

	di.Init("filename")

	di.GetDependencies().Database.Migrate()

	server := server.NewServer()
	server.Start()

	// pool := database.NewConnectionPool(config.Config{})

	// for i := 0; i < 110; i++ {
	// 	c, err := pool.GetConnection()
	// 	fmt.Println(err)
	// 	skip(c)
	// }

}

func skip(...interface{}) {}

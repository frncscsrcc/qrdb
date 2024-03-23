package database

import (
	"database/sql"
	"errors"
	"log"
	"qrdb/qrdb/pkg/config"
	"qrdb/qrdb/pkg/database/migrations"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ConnectionPool struct {
	lock        sync.Mutex
	connections []*Connection
}

type Connection struct {
	pool       *ConnectionPool
	connection *sql.DB
	tx         *sql.Tx
	free       bool
	dirty      bool
}

func NewConnectionPool(config config.Config) *ConnectionPool {
	pool := &ConnectionPool{
		lock:        sync.Mutex{},
		connections: make([]*Connection, 0),
	}

	poolSize := 100

	log.Printf("Creating connection pool size=%d", poolSize)
	for i := 0; i < poolSize; i++ {
		conn, err := sql.Open("sqlite3", "./test.Database")
		if err != nil {
			continue
		}
		pool.connections = append(pool.connections, &Connection{
			pool:       pool,
			connection: conn,
			free:       true,
			dirty:      false,
		})
	}

	return pool
}

func (pool *ConnectionPool) GetConnection() (*Connection, error) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	for j := 1; j <= 3; j++ {
		for i, connection := range pool.connections {
			if connection.free {
				connection.free = false
				log.Print("Returned connection ", i)
				return connection, nil
			}
		}
		time.Sleep(time.Second * time.Duration(j))
	}

	return nil, errors.New("pool exaused")
}

func (pool *ConnectionPool) CloseAll() error {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	for _, c := range pool.connections {
		err := c.connection.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Connection) Begin() (*sql.Tx, error) {
	tx, err := c.connection.Begin()
	c.tx = tx
	return tx, err
}

func (c *Connection) Rollback() error {
	err := c.tx.Rollback()
	if err != nil {
		c.Release()
	}
	return err
}

func (c *Connection) Commit() error {
	if c.dirty {
		return errors.New("dirty connection, can not commit")
	}

	err := c.tx.Commit()
	if err == nil {
		c.Release()
	} else {
		c.dirty = true
	}
	return err
}

func (c *Connection) Release() {
	c.pool.lock.Lock()
	defer c.pool.lock.Unlock()
	c.tx.Rollback()
	c.free = true
}

func (pool *ConnectionPool) Migrate() {
	log.Print("Starting migrations")

	connection, _ := pool.GetConnection()

	tx, _ := connection.Begin()

	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS _migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			created datetime
		)`,
	)
	if err != nil {
		connection.Rollback()
		panic(err)
	}

	lastVersion := 0
	tx.QueryRow("SELECT id FROM _migrations ORDER BY created DESC LIMIT 1").Scan(&lastVersion)

	migrations := migrations.GetMigrations()
	for i := lastVersion + 1; i < len(migrations); i++ {

		_, err := tx.Exec(migrations[i])
		if err != nil {
			connection.Rollback()
			panic(err)
		}

		tx.Exec("INSERT INTO _migrations VALUES (?, ?)", i, time.Now())
		if err != nil {
			connection.Rollback()
			panic(err)
		}

		log.Println("Migrated", migrations[i])
	}

	connection.Commit()
	log.Print("Migrations completed")
}

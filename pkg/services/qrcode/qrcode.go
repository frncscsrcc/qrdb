package qrcode

import (
	"log"
	"qrdb/pkg/database"
	"qrdb/pkg/di"
	"qrdb/pkg/uid"
	"time"
)

type qrCodeService struct {
	dbSession    *database.Connection
	uidGenerator uid.UIDGenerator
}

func NewService() qrCodeService {
	dep := di.GetDependencies()

	dbSession, _ := dep.Database.GetConnection()

	return qrCodeService{
		uidGenerator: dep.UIDGenerator,
		dbSession:    dbSession,
	}
}

func (s qrCodeService) Create(data map[string]string) (string, error) {
	uid := di.GetDependencies().UIDGenerator.GetUID()
	now := time.Now()

	tx, _ := s.dbSession.Begin()

	for key, value := range data {
		_, err := tx.Exec(`
			INSERT INTO data (id, key, value, created) VALUES (?, ?, ?, ?)`,
			uid, key, value, now,
		)
		if err != nil {
			return "", err
		}
	}

	s.dbSession.Commit()

	keys := ""
	for key := range data {
		keys += key + ";"
	}
	log.Printf("Created QR code=%s keys=%s user=%d", uid, keys, -1)

	return uid, nil
}

func (s qrCodeService) Get(code string) (map[string]string, error) {
	data := make(map[string]string)

	tx, _ := s.dbSession.Begin()

	rows, err := tx.Query("SELECT key, value FROM data WHERE id = ?", code)
	if err != nil {
		return map[string]string{}, err
	}

	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return map[string]string{}, err
		}
		data[key] = value
	}

	s.dbSession.Release()
	return data, nil
}

package migrations

func GetMigrations() []string {
	migrations := make([]string, 1)

	migrations = append(migrations,
		`
			CREATE TABLE IF NOT EXISTS data (
				id varchar(255),
				key varchar(255),
				value varchar(255),
				created datetime
			)
		`,
	)

	return migrations
}

package sqlitesessionstores

// SQL query for creating sessions table. %s is replaced by the table name.
var queryCreate = `
	CREATE TABLE IF NOT EXISTS %s (
		data BLOB,
		dateCreated TIMESTAMP NOT NULL,
		flashes BLOB,
		id TEXT PRIMARY KEY,
		userId TEXT
	);

	CREATE INDEX IF NOT EXISTS %sByDateCreated ON %s (
		dateCreated
	);

	CREATE INDEX IF NOT EXISTS %sByUserIdDateCreated ON %s (
		userId,
		dateCreated
	);
`

// SQL query for deleting sessions. %s is replaced by the table name.
var queryDelete = "DELETE FROM %s WHERE id = ?"

// SQL query for saving sessions. %s is replaced by the table name.
var querySave = `
	INSERT OR REPLACE INTO %s (
		data, dateCreated, flashes, id, userId
	) VALUES (
		?, ?, ?, ?, ?
	);
`

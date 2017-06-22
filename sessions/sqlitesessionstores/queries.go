package sqlitesessionstores

// SQL query for creating sessions table. %s is replaced by the table name.
const queryCreate = `
	CREATE TABLE IF NOT EXISTS %s (
		data TEXT,
		dateCreated TIMESTAMP NOT NULL,
		flashes TEXT,
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

// SQL query for deleting a session. %s is replaced by the table name.
const queryDelete = "DELETE FROM %s WHERE id = ?"

// SQL query for getting a session. %s is replaced by the table name.
const queryGet = `
	SELECT
		data,
		dateCreated,
		flashes,
		userId
	FROM
		%s
	WHERE
		id = ?
	LIMIT 1
`

// SQL query for saving a session. %s is replaced by the table name.
const querySave = `
	INSERT OR REPLACE INTO %s (
		data, dateCreated, flashes, id, userId
	) VALUES (
		?, ?, ?, ?, ?
	);
`

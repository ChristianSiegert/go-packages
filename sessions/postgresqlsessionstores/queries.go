package postgresqlsessionstores

// SQL query for creating sessions table. %s is replaced by the table name.
const queryCreate = `
	CREATE TABLE IF NOT EXISTS %s (
		data TEXT NOT NULL,
		date_created TIMESTAMP WITH TIME ZONE NOT NULL,
		flashes TEXT NOT NULL,
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS %s_by_date_created ON %s (
		date_created
	);

	CREATE INDEX IF NOT EXISTS %s_by_user_id_date_created ON %s (
		user_id,
		date_created
	);
`

// SQL query for deleting sessions. %s is replaced by the table name.
const queryDelete = "DELETE FROM %s WHERE id = $1"

// SQL query for getting a single session. %s is replaced by the table name.
const queryGet = `
	SELECT
		data,
		date_created,
		flashes,
		user_id
	FROM
		%s
	WHERE
		id = $1
	LIMIT 1
`

// SQL query for saving sessions. %s is replaced by the table name.
const querySave = `
	INSERT INTO %s (
		data, date_created, flashes, id, user_id
	) VALUES (
		$1, $2, $3, $4, $5
	) ON CONFLICT (id) DO UPDATE SET
		data = $1,
		date_created = $2,
		flashes = $3,
		user_id = $5
`

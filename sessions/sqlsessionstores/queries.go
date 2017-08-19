package sqlsessionstores

const (
	queryCreate = "create"
	queryDelete = "delete"
	queryGet    = "get"
	querySave   = "save"
)

var queries = map[dialect]map[string]string{
	DialectPostgreSQL: map[string]string{
		queryCreate: `
			CREATE TABLE IF NOT EXISTS %s (
				data text NOT NULL,
				date_created timestamp with time zone DEFAULT now() NOT NULL,
				flashes text NOT NULL,
				id text PRIMARY KEY,
				user_id text NOT NULL,
				CHECK (id != '')
			);

			CREATE INDEX IF NOT EXISTS %s_date_created ON %s (
				date_created
			);

			CREATE INDEX IF NOT EXISTS %s_user_id_date_created ON %s (
				user_id,
				date_created
			);
		`,
		queryDelete: "DELETE FROM %s WHERE id = $1",
		queryGet: `
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
		`,
		querySave: `
			INSERT INTO %s (
				data, date_created, flashes, id, user_id
			) VALUES (
				$1, $2, $3, $4, $5
			) ON CONFLICT (id) DO UPDATE SET
				data = $1,
				date_created = $2,
				flashes = $3,
				user_id = $5
		`,
	},

	DialectSQLite: map[string]string{
		queryCreate: `
			CREATE TABLE IF NOT EXISTS %s (
				data TEXT,
				date_created TIMESTAMP NOT NULL,
				flashes TEXT,
				id TEXT PRIMARY KEY,
				user_id TEXT
			);

			CREATE INDEX IF NOT EXISTS %s_date_created ON %s (
				date_created
			);

			CREATE INDEX IF NOT EXISTS %s_user_id_date_created ON %s (
				user_id,
				date_created
			);
		`,
		queryDelete: "DELETE FROM %s WHERE id = ?",
		queryGet: `
			SELECT
				data,
				date_created,
				flashes,
				user_id
			FROM
				%s
			WHERE
				id = ?
			LIMIT 1
		`,
		querySave: `
			INSERT OR REPLACE INTO %s (
				data, date_created, flashes, id, user_id
			) VALUES (
				?, ?, ?, ?, ?
			);
		`,
	},
}

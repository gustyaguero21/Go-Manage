package config

//router params

const (
	Port = ":8080"
)

//Database params

const (
	DBDriver = "sqlite3"
	DBPath   = "internal/data/users.db"
)

//Database queries

const (
	CreateTableQuery = `CREATE TABLE users (id TEXT NOT NULL UNIQUE PRIMARY KEY, name TEXT NOT NULL, surname TEXT NOT NULL, username TEXT NOT NULL UNIQUE, email TEXT NOT NULL UNIQUE, password TEXT NOT NULL UNIQUE);`
	ExistsQuery      = `SELECT EXISTS(SELECT 1 FROM users WHERE username = ?);`
	SearchUserQuery  = `SELECT * FROM users WHERE username=?;`
	SaveUserQuery    = `INSERT INTO users (id,name,surname,username,email,password) VALUES (?,?,?,?,?,?);`
	DeleteUserQuery  = `DELETE FROM users WHERE username = ?;`
	UpdateUserQuery  = `UPDATE users SET name = ?, surname = ?, email = ? WHERE username = ?;`
)

//Repository test queries

const (
	TestExistsQuery = `SELECT EXISTS\(SELECT 1 FROM users WHERE username = \?\);`
	TestSearchQuery = `SELECT \* FROM users WHERE username=\?`
	TestSaveQuery   = "INSERT INTO users"
	TestDeleteQuery = `DELETE\s+FROM\s+users\s+WHERE\s+username\s*=\s*\?;`
	TestUpdateQuery = `UPDATE users SET name = \?, surname = \?, email = \? WHERE username = \?;`
)

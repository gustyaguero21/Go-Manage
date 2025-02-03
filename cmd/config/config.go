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
)

//Repository test queries

const (
	TestSearchQuery = `SELECT \* FROM users WHERE username=\?`
	TestExistsQuery = `SELECT COUNT\(\*\) FROM users WHERE username = \?`
	TestSaveQuery   = "INSERT INTO users"
)

package config

//Router params

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
	CreateTableQuery   = `CREATE TABLE users (id TEXT NOT NULL UNIQUE PRIMARY KEY, name TEXT NOT NULL, surname TEXT NOT NULL, username TEXT NOT NULL UNIQUE, email TEXT NOT NULL UNIQUE, password TEXT NOT NULL UNIQUE);`
	SearchUserQuery    = `SELECT * FROM users WHERE username=?;`
	SaveUserQuery      = `INSERT INTO users (id,name,surname,username,email,password) VALUES (?,?,?,?,?,?);`
	DeleteUserQuery    = `DELETE FROM users WHERE username = ?;`
	UpdateUserQuery    = `UPDATE users SET name = ?, surname = ?, email = ? WHERE username = ?;`
	ChangeUserPwdQuery = `UPDATE users SET password = ? WHERE username = ?;`
)

//Repository test queries

const (
	TestSearchQuery    = `SELECT \* FROM users WHERE username=\?`
	TestSaveQuery      = "INSERT INTO users"
	TestDeleteQuery    = `DELETE\s+FROM\s+users\s+WHERE\s+username\s*=\s*\?;`
	TestUpdateQuery    = `UPDATE users SET name = \?, surname = \?, email = \? WHERE username = \?;`
	TestChangePwdQuery = "UPDATE users SET"
)

//Errors

const (
	ErrAllFieldsAreRequired = "all fields are required"
	ErrInvalidPassword      = "invalid password"
	ErrInvalidEmail         = "invalid email"
	ErrUserAlreadyExists    = "user already exists"
	ErrSaveUser             = "error creating user. Error: "
	ErrSearchUser           = "error searching user. Error: "
	ErrUserNotFound         = "user not found"
	ErrDeleteUser           = "error deleting user. Error: "
	ErrUpdateUser           = "error updating user. Error: "
	ErrChangePwd            = "error changing user password. Error: "
	ErrEmptyQueryParam      = "empty query param"
)

//Handler messages

const (
	SuccessStatus = "success"
	CreateMessage = "user created successfully"
	SearchMessage = "user found successfully"
)

package database

var UsersTable = `
    CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	nickname TEXT NOT NULL UNIQUE,
    gender TEXT NOT NULL,
    age INTEGER NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
`

var SessionsTable = `
    CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL UNIQUE,   
    user_id INTEGER NOT NULL,         
    expires_at DATETIME NOT NULL,    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE);
`

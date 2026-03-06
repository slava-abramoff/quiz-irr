package models

type Admin struct {
	ID       int
	Email    string
	Passwrdo string
	FullName string
}

// id (UUID/Int, PK)

//    email (String, Unique)

//    password_hash (String)

//    full_name (String)

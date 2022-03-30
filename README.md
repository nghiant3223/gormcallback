# GormCallback

Utility callbacks for Gorm, the fantastic ORM library for Golang.

## Callbacks

- ExplainSQL: print EXPLAIN of SQL statement to the stdout

## Install

`
go get github.com/nghiant3223/gormcallback
`

## Usage

```go
package main

import (
	"github.com/nghiant3223/gormcallback"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// Register *gorm.DB with ExplainSQL callback for all SQL statements.
	_ = gormcallback.RegisterExplainSQL(db)

	// Register *gorm.DB with some SQL statements. Only SELECT and INSERT are registered as below.
	_ = db.Callback().Query().Register("gormcallback", gormcallback.ExplainSQL)
	_ = db.Callback().Create().Register("gormcallback", gormcallback.ExplainSQL)
}
```

## Demo

```go
package main

import (
	"github.com/nghiant3223/gormcallback"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Username string
	Password string
}

func main() {
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	
	// Create table `users`.
	if err = db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}

	// Register *gorm.DB with ExplainSQL.
	if err = gormcallback.RegisterExplainSQL(db); err != nil {
		panic(err)
	}

	newUser := &User{
		Username: "nghiant3223",
		Password: "helloworld",
	}
	if err = db.Create(newUser).Error; err != nil {
		panic(err)
	}

	var user *User
	if err = db.Where("username=?", newUser.Username).First(&user).Error; err != nil {
		panic(err)
	}
}
```


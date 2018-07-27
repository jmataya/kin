# Kin

[![CircleCI](https://img.shields.io/circleci/project/github/jmataya/kin.svg)](https://circle.com/gh/jmataya/kin)
[![Documentation](https://godoc.org/github.com/jmataya/kin?status.svg)](http://godoc.org/github.com/jmataya/kin)

*Kin is very experimental - don't use it for anything real!*

Kin is an opinionated, PostgeSQL-specific database driver that places an
emphasis on writing SQL and explicit mappings to database tables. Instead of
relying on reflection or other conventions, it gives users the control to
write their own SQL queries, uses simple patterns that focus on control.

This means that kin isn't for everyone! It's certainly not a
full-featured ORM in the mold of Gorm or Ruby's ActiveRecord. Instead, Kin
smoothes out the rough edges and gives better error handling to SQL packages
already built into Go.

Finally, support for running migrations is built in.

## Getting Started

Installing kin is as simple as running:

```shell
go get github.com/jmataya/kin
```

Next, define your models, connect to the database, and you're off to the
races!

```golang
package main

import (
        "fmt"
        "time"

        "github.com/jmataya/kin"
)

// user is a representation of a table called "users" with columns
// id (int), name (text), attributes (jsonb), is_active (bool), and
// created_at (timestamp).
type user struct {
        ID         int
        Name       string
        Attributes map[string]interface{}
        IsActive   bool
        CreatedAt  time.Time
}

// Columns defines the mapping between database columns, their type,
// and the model. It's used internally by kin to wire everything up.
func (u *user) Columns() []FieldBuilder {
        return []FieldBuilder{
               IntField("id", &u.ID),
               StringField("name", &u.Name),
               JSONField("attributes", &u.Attributes),
               BoolField("is_active", &u.IsActive),
               TimeField("created_at", &u.CreatedAt),
       }
}

func main() {
        // Connect with a raw Postges URL.
        dbStr := "postgresql://localhost:5432/kin_test?user=kin"
        db, _ := kin.NewConnection(dbStr)
        defer db.Close()
        
        // For most operations, Kin leverages SQL.
        queryStr := "SELECT * FROM users WHERE id = $1"
        id := 1
        
        // Output a single result into the user struct defined above.
        u := new(user)
        err := db.Query(queryStr, id).OneAndExtract(u)
        if err != nil {
                panic(err)
        }
        
        // Print the result.
        fmt.Printf("User: %+v\n", u)
}
```

## Author

Jeff Mataya (jeff@jeffmataya.com)

## License

MIT

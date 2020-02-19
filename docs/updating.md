# Data Models

Every entity the cache updater works with stays in src/entities directory.
As regards every single model, it's tightly coupled with appropriate database
table in the database. PostgreSQL is used in this project.

# Migrations

Every database entity can be managed by migrations, thanks to:

[go-pg migrations](https://github.com/go-pg/migrations)

# Updating approach

In order to perform operations with database the code logic is
split between different kinds of controllers. Splitting logic
is one of the ways for making code composing more pure. 

> DbController:

```go
type DbController struct {}
```

This controller is responsible for database operations. It's tightly
coupled with "go-pg" ORM for Golang.

> UpdateController:

```go
type UpdateController struct {}
```

Responsible for data pulling from different sources.
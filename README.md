[![Go Report Card](https://goreportcard.com/badge/github.com/adel-hadadi/ququery)](https://goreportcard.com/report/github.com/adel-hadadi/ququery)
[![codecov](https://codecov.io/github/adel-hadadi/ququery/graph/badge.svg?token=IkdEDjVmb4)](https://codecov.io/github/adel-hadadi/ququery)

# QuQuery: Golang SQL Builder

```
   ____
  / __ \__  ______ ___  _____  _______  __
 / / / / / / / __ `/ / / / _ \/ ___/ / / /
/ /_/ / /_/ / /_/ / /_/ /  __/ /  / /_/ /
\___\_\__,_/\__, /\__,_/\___/_/   \__, /
              /_/                /____/

```

## About

QuQuery is simple and efficient SQL databases query builder that provide zero dependency and zero type reflection
in your code base to make repositories more readable for first look.

## Why i make `Ququery` package

When I was learning Golang, I noticed that many people coding in Golang for the first time were using GORM. However, after a short period of using GORM, I found myself dealing with convoluted code that I couldn't understand, and it often didn't work as expected.

I searched for other solutions and discovered that in many large companies and projects, programmers prefer to use pure SQL queries with the standard database/sql package. This approach avoids several issues commonly encountered with GORM:

1. **Performance Overhead**:
   GORM introduces an additional layer of abstraction which can result in performance overhead compared to writing raw SQL queries. This can be significant in high-performance applications.

2. **Complex Queries**:
   For complex queries involving multiple joins, subqueries, and custom SQL, GORM's abstraction can become cumbersome and harder to manage, often requiring raw SQL anyway.

3. **Debugging Difficulties**:
   Debugging GORM issues can be challenging because it abstracts away the SQL, making it harder to understand what exact queries are being generated and executed.

4. **Code Complexity** :
   Over time, as projects grow, the GORM code can become complex and harder to maintain, especially if it's not used consistently across the codebase.

After a long time of writing SQL queries, my coworkers and I grew tired of writing repetitive queries. To solve this issue, I decided to create `ququery`—a query builder for Golang. This tool aims to simplify the process of building SQL queries, making your code more readable and maintainable, while avoiding the aforementioned problems with GORM.

## Installation And Usage

For installing ququery in your project should run below command in root of project.

```shel
go get github.com/adel-hadadi/ququery@latest
```

Every database operation such as (`UPDATE`, `INSERT`, `DELETE`, `SELECT`) in ququery have specific methods and they can be different from other one so let's explain each operation methods one by one.

## Select Statements

### Specifying a Select Clause

You may not always want to select all columns from database table.
Using the `Columns` method you can specify each column that you want to fetch from database.

```go
query := ququery.Select("table_name").Columns("id", "name", "email").Query()

log.Println(query) // query => SELECT id, name, email FROM table_name
```

For situations that you want to fetch all columns you can call `Select` method without `Columns`.

## Joins

The query builder also be used to add join clauses to your queries.
To perform a basic inner join, you may use the `Join` method on a query builder instance.
The first arguments passed to `Join` method is the name of the table you need to join to,
while the second argument specify the column constraints for the join.
You may even join multiple tables in a single query:

```go
query := ququery.Select("users").
    Join("posts", "posts.user_id = users.id").
    Query()

log.Println(query) // query => SELECT * FROM users INNER JOIN posts ON posts.user_id = users.id
```

### Left join / Right join

if you would like to perform `left join` or `right join` instead of an `inner join`,
use `LeftJoin` or `RightJoin` methods. This methods have the same signature as the
`Join` method:

```go
leftJoin := ququery.Select("users").LeftJoin("posts", "posts.user_id = users.id").Query()
log.Println(leftJoin) // query => SELECT * FROM users LEFT JOIN posts ON posts.user_id = users.id

rightJoin := ququery.Select("users").RightJoin("posts", "posts.user_id = users.id").Query()
log.Println(rightJoin) // query => SELECT * FROM users RIGHT JOIN posts ON posts.user_id = users.id
```

### With

Also if you want to load a simple belongs to relations you can use `With` method.
This method take a list of entities and then automatically load relations:

```go
query := ququery.Select("users").With("role", "wallet").Query()
log.Println(query) // query => SELECT * FROM users LEFT JOIN roles ON roles.id = user.role_id LEFT JOIN wallets ON wallets.id = users.wallet_id
```

## Basic Where Clauses

### Where Clauses

You may use the query builder's `Where` method to add "where" clauses to the query.
The most basic call to the `Where` method requires two arguments. The first argument
is the name of the column. The second argument is an operator, which can be any of
the database's supported operators.

For example, the following query retrieves users where the value of the votes column
is equal to $1 and the value of the age column is greater than $2:

```go
query := ququery.Select("users").
    Where("votes", "=").
    Where("age", ">").
    Query()

log.Println(query) // query => SELECT * FROM users WHERE votes = $1 AND age > $2
```

For convenience, if you want to verify that a column is `=` to a given value, you may call `Where` method with just column name.
`Ququery` will assume you would like to use the `=` operator:

```go
query := ququery.Select("users").
    Where("votes").
    Query()

log.Pritln(query) // query => SELECT * FROM users WHERE votes = $1
```

## Or Where Clauses

When chaining together calls to the query builder's `Where` method, the "where" clauses will be joined together using the `AND` operator. However, you may use the `OrWhere` method to join a clause to the query using the `OR` operator. The `OrWhere` method accepts the same arguments as the `Where` method:

```go
query := ququery.Select("users").
    Where("votes").
    OrWhere("name")
    Query()

log.Pritln(query) // query => SELECT * FROM users WHERE votes = $1 OR name = $2
```

If you need to group multiple where clauses together you can use `WhereGroup` method:

```go
query := ququery.Select("users").
    Where("votes").
    WhereGroup(func(subQuery MultiWhere) string {
        return subQuery.Where("name").
            OrWhere("votes", ">").
            Query()
    }).
    Query()

log.Pritln(query) // query => SELECT * FROM users WHERE votes = $1 OR name = $2
```

## Additional Where Clause

### WhereLike / OrWhereLike

The `WhereLike` method allows you to add "LIKE" clauses to query for
pattern matching. These methods provide a database-agnostic way
performing string matching queries, with the ability to toggle
case-sensitivity. By default, string matching is case-insensitive:

```go
query := ququery.Select("users").WhereLike("name").Query()
log.Println(query) // query => SELECT * FROM users WHERE name LIKE $1
```

The `OrWhereLike` method allows you to add an "or"A clause with a LIKE
condition:

```go
query := ququer.Select("users").
    Where("votes", ">").
    OrWhereLike("name").
    Query()

log.Println(query) // query => SELECT * FROM users WHERE votes > $1 OR WHERE name LIKE $2
```

### WhereNull / WhereNotNull / OrWhereNull / OrWhereNotNull

The `WhereNull` method verifies that the value of the given column is `NULL`:

```go
query := ququery.Select("users").WhereNull("updated_at").Query()
log.Println(query) // query => SELECT * FROM users WHERE updated_at IS NULL
```

The `WhereNotNull` method verifies that the column's value is not `NULL`:

```go
query := ququery.Select("users").WhereNotNull("updated_at").Query()
log.Println(query) // query => SELECT * FROM users WHERE updated_at IS NOT NULL
```

# Ordering, Grouping, Limit and offset

## Ordering

### The `OrderBy` Method

The `OrderBy` method allows you to sort the results of the query by a given column. The First argument accepted by the `OrderBy` method should be the column you wish to sort by, while the second argument determines the direction of the sort and may not either `asc` or `desc`:

```go
query := ququery.Select("users").OrderBy("name", "desc").Query()
log.Println(query) // query => SELECT * FROM users ORDER BY name DESC
```

## Limit and Offset

You may use the `Limit` and `Offset` methods to limit the number of results returned from the query or to skip a given number of results in the query:

```go
query := ququery.Select("users").Limit().Offset().Query()
log.Println(query) // query => SELECT * FROM users LIMIT $1 OFFSET $2
```

# Insert Statements

The query builder also provides an `Insert` method that may be used to insert records into database table. The `Insert` method accepts a list of column names.

```go
query := ququery.Insert("users").Into("email", "votes").Query()
log.Println(query) // query => INSERT INTO users (email, votes) VALUES ($1, $2)
```

# Update Statements

In addition to inserting records into the database, the query builder can also update existing records using the `Update` method. The `Update` method, like the `Insert` method, accepts a list of columns that should be updated:

```go
query := ququery.Update("users").Where("id").Set("email", "email_verified").Query()
log.Println(query) // query => UPDATE users SET email = $1, email_verified = $2 WHERE id = $3
```

# Delete Statements

The query builder's `Delete` method may be used to delete records from the table:

```go
query := ququery.Delete("users").Where("votes", ">").Query()
log.Println(query) // query => DELETE FROM users WHERE votes > $1
```

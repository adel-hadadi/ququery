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

QuQuery is simple and efficient SQL database query builder that provide zero dependency and zero type reflection
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

you may not always want to select all columns from database table.
Using the `Columns` method you can specify each column that you want to fetch from database

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

## With

Also if you want to load a simple belongs to relation you can use `With` method.
This method take a list of entities and then automatically load relations:

```go
query := ququery.Select("users").With("role", "wallet").Query()
log.Println(query) // query => SELECT * FROM users LEFT JOIN roles ON roles.id = user.role_id LEFT JOIN wallets ON wallets.id = users.wallet_id
```

# Basic Where Clauses

## Where Clauses

You may use the query builder's `Where` method to add "where" clauses to the query.
The most basic call to the `Where` method requires two arguments. The first argument
is the name of the column. The second argument is an operator, which can be any of
the database's supported operators.

For example, the following query retrieves users where the value of the votes column
is equal to x and the value of the age column is greater than y:

```go
query := ququery.Select("users").
    Where("votes", "=").
    Where("age", ">").
    Query()

log.Println(query) // query => SELECT * FROM users WHERE votes = $1 AND age > $2
```

For convenience, if you want to verify that a column is `=` to a given value, you may call `Where` method with just column name. `Ququery` will assume you would like to use the `=` operator:

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

### Where

This method can use different ways. The first one is just a simple where that should just specify column and automatically operation is (column = $1) but second way you can specify operation type (`!=`, `LIKE`).

```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Select("users").Where("id").Query()
    // query => SELECT * FROM users WHERE id = $1

    query = ququery.Select("users").Where("id", "!=").Query()
    // query => SELECT * FROM users WHERE id != $1
}
```

Also can use `OrWhere` for checking that any one condition is true

```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Select("users").Where("id").OrWhere("email").Query()
    // query => SELECT * FROM users WHERE id = $1 OR email = $2
}

```

For `Postgresql` users we have `StrPos` method to check if value is not an empty string, find the position from where the substring is being matched within the string also this.

```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Select("users").StrPos("name").Query()
    // query => SELECT * FROM users WHERE (STRPOS(name, %1) > 0 OR $2 = '')
}

```

For using `WHERE IN (subquery)` query can use `WhereInSubquery` that take column and a function with `SelectQuery` structure and should return a query.

```go

package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Select("users").
        WhereInSubquery("id", func (q SelectQuery) string {
            return q.Table("orders").
                Column("user_id")
                Where("price", ">=").
                Query()
        }).
        Query()
    // query => SELECT * FROM users WHERE id IN (SELECT user_id FROM orders where price >= $1)
}

```

### Join

With Join method can load relations.

```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Select("products").Join("categories", "categories.id = product.category_id").Query()
    // query => SELECT * FROM users LEFT JOIN categories ON categories.id = product.category_id
}
```

But if you want to load one to many relations like above you can just simply use `With` method that take list of entity names and load them.

```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Select("products").
        With("category", "warehouse", "brand").
        Query()
    // query => SELECT * FROM users LEFT JOIN categories ON categories.id = products.category_id LEFT JOIN warehouses ON warehouses.id = products.warehouse_id LEFT JOIN brands.id = products.brand_id
}
```

Select query have many other methods that we list them below.

| Method       | Description                                                                                                |
| ------------ | ---------------------------------------------------------------------------------------------------------- |
| OrderBy      | take column name and sort direction to order items                                                         |
| Limit        | for taking a limited rows from database                                                                    |
| Offset       | specify where to start taking rows                                                                         |
| Table        | if you get `SelectQuery` structure with different way `Select` method, can set table name with this method |
| WhereNull    | get rows where specific column value is null                                                               |
| WhereNotNull | get rows where specific column value is not null                                                           |

## Insert

Insert method give a `InsertQuery` structure which come with low number of methods

### Into

First and important method that come with `InsertQuery` structure is this method that specify insert query columns.

```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Insert("users").Into("name", "email").Query()
    // query => INSERT INTO users (name, email) VALUES($1, $2)
}
```

### Returning

This method specify columns that you want to return after that insert query successfully executed.

```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Insert("users").
        Into("name", "email").
        Returning("id").
        Query()
    // query => INSERT INTO users (name, email) VALUES ($1, $2) RETURNING (id)
}
```

## Update

### Set

For specify columns that you want change them you should use `Set` method.
Also `UpdateQuery` have `Where` and `OrWhere` method that can be used together.

```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Update("users").
        set("name").
        Where("id").
        Query()
    // query => UPDATE users SET name = $1 WHERE id = $2
}
```

## Exists

Exists query is used for checking that any row with exists with given conditions or not. This query have `Where` and `OrWhere` method.

```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Exists("users").
        Where("email").
        Query()
    // query => SELECT EXISTS (SELECT true FROM users WHERE email = $1)
}
```

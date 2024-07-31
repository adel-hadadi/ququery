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
QuQuery is simple and efficient query builder that provide zero dependency and zero type reflection in your code base so come to make our repositories more beautiful 😉.

## Why Ququery builted
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
``` shel
go get github.com/adel-hadadi/ququery@latest
```

Every database operation such as (`UPDATE`, `INSERT`, `DELETE`, `SELECT`) in ququery have specific methods and they can be different from other one so let's explain each operation methods one by one.

## Select
first and most important operation that exists is `select` query that come with some cool methods. For creating new select query use below example:
```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Select("table_name")
}
```
the above code return a `SelectQuery` struct that can call select query method on this structure.

### Columns
This method get a list of expected columns. If want to select all columns just don't call `Columns` method the `ququery` automatically know that if columns not specified should select all columns.
```go
package main

import "github.com/adel-hadadi/ququery"

func main() {
    query := ququery.Select("users").Columns("id", "name", "email").Query()
    // query => SELECT id, name, email FROM users

    query = ququery.Select("users").Query()
    // query => SELECT * FROM users
}
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

| Method | Description |
| -------------- | --------------- |
| OrderBy | take column name and sort direction to order items |
| Limit | for taking a limited rows from database |
| Offset | specify where to start taking rows |
| Table | if you get `SelectQuery` structure with different way `Select` method, can set table name with this method |

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

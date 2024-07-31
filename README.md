# QuQuery: Golang SQL Builder

## About
QuQuery is simple and efficient query builder that provide zero dependency and zero type reflection in your code base so come to make our repositories more beautiful ;).

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

### Select
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

    query = ququery.Select("users").Where("id", "!=")
    // query => SELECT * FROM users WHERE id != $1
}
```

Also can use `OrWhere` for alternative conditions
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


### Usage
* `pgping` works by reading your `PGCONN` environment variable and attempting to
connect to it.
* If it's able to establish a connection with the provided Postgres server and
credentials, it exits with a `0` status.
* If it cannot, it prints the error to `stdout` and exits with a non-zero status.

```
PGCONN=postgres://postgres@localhost:5432/postgres?sslmode=disable pgping

$>2018/03/09 12:18:11 Postgres server READY and accepting connections...
```

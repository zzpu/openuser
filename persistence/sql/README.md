# SQL Migrations

To create a new [fizz](https://gobuffalo.io/en/docs/db/fizz/) migration run in the project root:

```
$ name=
$ ory dev pop migration create ./persistence/sql/migrations/templates $name
```

To create SQL migrations, target each database individually and run

```
$ dialect=mysql  # or postgres|cockroach|sqlite
$ name=
$ ory dev pop migration create -d=$dialect ./persistence/sql/migrations/templates $name
$ soda generate sql -e mysql -c ./persistence/sql/.soda.yml -p ./persistence/sql/migrations/templates [name]
$ soda generate sql -e sqlite -c ./persistence/sql/.soda.yml -p ./persistence/sql/migrations/templates [name]
$ soda generate sql -e postgres -c ./persistence/sql/.soda.yml -p ./persistence/sql/migrations/templates [name]
$ soda generate sql -e cockroach -c ./persistence/sql/.soda.yml -p ./persistence/sql/migrations/templates [name]
```

and remove the `sqlite` part from the newly generated file to create a SQL migrations that works with all
aforementioned databases.

## Rendering Migrations

Because migrations needs to be backwards compatible, and because fizz migrations might change, we render
fizz migrations to raw SQL statements using `make migrations-render`.

The concrete migrations being applied can be found in [`./migrations/sql`](migrations/sql_bak).

# pg_sleeper
Straightforward tool that establish long running connections to DB and periodically report on their health.
By default it open a series of connection for the following : {1,5,15,30,60}min, {2,4,8,16,24}hour, {2,3,4,5,6,7}day, {2,3,4}week

The intent is to confirm/infirm the presence of a network hop that wreck havoc on our DB connections.

# usage

```
$ pg_helper --pg-url="postgresql://[user[:password]@][hostname][:port][/dbname][?param1=value1&...]" --statement="select 1"
establishing connections...done
2020-02-09-00:01:00 1min OK (0.001 sec)
2020-02-09-00:02:00 1min OK (0.001 sec)
2020-02-09-00:03:00 1min OK (0.001 sec)
2020-02-09-00:04:00 1min OK (0.001 sec)
2020-02-09-00:05:00 1min OK (0.001 sec)
2020-02-09-00:05:00 5min OK (0.001 sec)
2020-02-09-00:06:00 1min OK (0.001 sec)
2020-02-09-00:07:00 1min OK (0.001 sec)
2020-02-09-00:08:00 1min OK (0.001 sec)
2020-02-09-00:09:00 1min OK (0.001 sec)
2020-02-09-00:10:00 1min OK (0.001 sec)
2020-02-09-00:10:00 5min FAILURE (0.020 sec) -- PG::ConnectionBad: PQconsumeInput()
2020-02-09-00:11:00 1min OK (0.001 sec)
2020-02-09-00:12:00 1min OK (0.001 sec)
...
```

# Development setup

Setup a local postgres env, this can be done with
```
docker-compose up
```

If you are on mac you can use the following to install the psql client natively to connect to that postgres instance
```
brew install libpq
brew link --force libpq

PGPASSWORD=pgpassword psql -U postgres -h localhost
```

Alternatively you can use the following connection URL to connect to that instance with your favorite tool
```
postgresql://postgres:pgpassword@localhost:5432/postgres
```

# Build and run

```
go build && POSTGRES_URL="postgresql://postgres:pgpassword@127.0.0.1:5432/postgres" ./pg_sleeper
```

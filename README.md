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

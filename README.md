# pg_sleeper
Straightforward tool that establish long running connections to DB and periodically report on their health.
By default it open a series of connection for the following : {1,5,15,30,60}min, {2,4,8,16,24}hour, {2,3,4,5,6,7}day, {2,3,4}week.
New connection attempt will also be performed each 1 minute in order to detect situation where network path may prevent new
connection to establish

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
2020-02-09-00:10:00 5min FAILURE (0.020 sec) -- FATAL: terminating connection due to unexpected postmaster exit (SQLSTATE 57P01)
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

# Notice - keepalive

A word of notion on keepalive.  Regarding which kind of network issue you seek to identify please keep in mind that pgx
make use of a 300 seconds keepalive by default.  This is a much more sensible keepalive setting
then most default OS ( usualy 7200 seconds ).  Meaning, this could mean pg_sleeper may end up not exhibit the same behavior as your
application.  To confirm/infirm you can use netstat to inspect the keepalive setting of your application tcp session vs what pg_sleeper is using.

```
netstat -apno | less | grep 5432
tcp        0      0 192.168.62.193:54254    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (299.60/0/0)
tcp        0      0 192.168.62.193:54248    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (53.84/0/0)
tcp        0      0 192.168.62.193:54270    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (53.84/0/0)
tcp        0      0 192.168.62.193:54276    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (53.84/0/0)
tcp        0      0 192.168.62.193:54268    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (266.83/0/0)
tcp        0      0 192.168.62.193:54250    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (228.50/0/0)
tcp        0      0 192.168.62.193:54280    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (86.61/0/0)
tcp        0      0 192.168.62.193:54262    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (53.84/0/0)
tcp        0      0 192.168.62.193:54278    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (299.60/0/0)
tcp        0      0 192.168.62.193:54264    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (70.22/0/0)
tcp        0      0 192.168.62.193:54266    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (70.22/0/0)
tcp        0      0 192.168.62.193:54256    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (288.30/0/0)
tcp        0      0 192.168.62.193:54274    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (70.22/0/0)
tcp        0      0 192.168.62.193:54260    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (102.99/0/0)
tcp        0      0 192.168.62.193:54258    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (290.70/0/0)
tcp        0      0 192.168.62.193:54282    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (229.14/0/0)
tcp        0      0 192.168.62.193:54272    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (53.84/0/0)
tcp        0      0 192.168.62.193:54252    192.168.153.100:5432    ESTABLISHED 47416/./pg_sleeper   keepalive (299.60/0/0)
```
> In the output above we observe that the next keepalive to be sent for the tcp connection opened by pgx is set to occur for most connection in ~300 secs
> ( I did that copy paste just after it previously reached 0 ).  In a nutshell this track the number of seconds before a new keep alive be sent, if something
> end up being sent through the tcp socket, or if the keepalive is sent, whichever come first, the counter go back to 300.  It is likely your application timer
> will be in the 7200 when it reset, something that can be a problem on some network that implement a very aggressive TCP timeout at the firewall level.
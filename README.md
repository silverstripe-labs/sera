# sera

Distributed Mutex locking using a redis cluster or a mysql database.

## Introduction

Sera allows mutual exclusive locking for a cluster of servers. It's normal use case is to
prevent several cronjobs or scheduled tasks running at the same time.

> Distributed locks are a very useful primitive in many environments where different 
> processes require to operate with shared resources in a mutually exclusive way.
>
> [Distributed locks with Redis](http://redis.io/topics/distlock)

## Usage

The normal use case is in a cronjob or scheduled services

	sera <expiry in seconds> <command to run> < .. arguments and flags to command>

Example usage in a cronjob:

	* * * * * root /usr/local/bin/sera 20 /bin/long-running-task --parameter hello

sera takes two arguments, the first one is how the is how long in seconds it would take to complete the task (upper bound). After this time another instance of sera will be able to run this job.

If the second argument (which is the command to run) takes longer than this time, other servers 
might expire the lock and start the task.

## Configuration for Redis

`/etc/sera.json`

	{
			"backend": "redis",
			"servers": [
					"127.0.0.1:6379",
					"127.0.0.1:6380",
			]
	}


## Configuration for MySQL

MySQL Support is experimental at this point since it doesn't use the same locking mechanism as the redis mutex.

It's MySQL [get_lock()](http://dev.mysql.com/doc/refman/5.0/en/miscellaneous-functions.html#function_get-lock) function. It does not use any databases or tables, it works by globally assigning a lock with a key and a timeout in seconds. The `get_lock()` will wait for that amount of seconds until it moves on. 

The `get_lock()` timeout value when using `sera` is <expiry in seconds>/10 and it will try for 10 to get the lock.

If the MySQL connection closes or dies, the lock will be unlocked.

`/etc/sera.json`

	{
			"backend": "mysql",
			"servers": [
					"sera:secret@tcp(127.0.0.1:3306)/?timeout=500ms"
			]
	}


## Resources

[http://redis.io/topics/distlock](http://redis.io/topics/distlock)

## Todo

 - Syslog logging
 - Configuration file in yaml
 - Warnings if all redis servers are unreachable
 - parameterize if sera should run the command even if no servers can be connected to



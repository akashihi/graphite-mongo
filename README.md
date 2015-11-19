# graphite-mongo

## What is this?

Mongo metrics to graphite gateway.

## Building it

1. Install [go](http://golang.org/doc/install)

2. Install "graphite-golang" go get -u github.com/marpaia/graphite-golang

3. Install "go-logging" go get -u github.com/op/go-logging

4. Install "mgo" go get -u gopkg.in/mgo.v2

5. Install "mgo/bson" go get -u gopkg.in/mgo.v2/bson

4. Compile graphite-mongo

        git clone git://github.com/akashihi/graphite-mongo.git
        cd graphite-mysql
        go build .

## Running it

Generally:

    graphite-mongo -host 127.0.0.1 -port 27017 -metrics-host 192.168.1.1 -metrics-port 2003 -metrics-prefix test -period 60

All parameters could be omited. Run with --help to het parameters description

## License 

See LICENSE file.

Copyright 2015 Denis V Chapligin <akashihi@gmail.com>

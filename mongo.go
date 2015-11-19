/*
   conntrack-logger
   Copyright (C) 2015 Denis V Chapligin <akashihi@gmail.com>
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

type Connections struct {
	Current int64 "current"
}

type Mem struct {
	Resident int64 "resident"
	Virtual  int64 "virtual"
}

type RWT struct {
	Readers int64 "readers"
	Writers int64 "writers"
	Total   int64 "total"
}

type GlobalLock struct {
	TotalTime     int64 "totalTime"
	LockTime      int64 "lockTime"
	CurrentQueue  RWT   "currentQueue"
	ActiveClients RWT   "activeClients"
}

type Opcounters struct {
	Insert  int64 "insert"
	Query   int64 "query"
	Update  int64 "update"
	Delete  int64 "delete"
	GetMore int64 "getmore"
	Command int64 "command"
}

type ExtraInfo struct {
	PageFaults       int64 "page_faults"
	HeapUsageInBytes int64 "heap_usage_bytes"
}

type Network struct {
	BytesIn  int64 "bytesIn"
	BytesOut int64 "bytesOut"
}

type BackgroundFlushing struct {
	Flushes   int64   "flushes"
	AverageMS float64 "average_ms"
}

type Replica struct {
	Name string "setName"
}

type ServerStatus struct {
	Process              string             "process"
	Connections          Connections        "connections"
	ExtraInfo            ExtraInfo          "extra_info"
	Mem                  Mem                "mem"
	Network              Network            "network"
	Opcounters           Opcounters         "opcounters"
	BackgroundFlushing   BackgroundFlushing "backgroundFlushing"
	GlobalLocks          GlobalLock         "globalLock"
	Replica              Replica            "repl"
	OpcountersReplicaSet Opcounters         "opcountersRepl"
}

type Database struct {
	Name  string "name"
	Empty bool   "empty"
}

type Databases struct {
	Database []Database "databases"
}

type Status struct {
	ServerStatus ServerStatus
	Databases    map[string]DatabaseStats
}

type DatabaseStats struct {
	Collections   int64                    "collections"
	Object        int64                    "objects"
	AvgObjectSize float64                  "avgObjSize"
	DataSize      int64                    "DataSize"
	StorageSize   int64                    "storageSize"
	Indexes       int64                    "indexes"
	IndexSize     int64                    "indexSize"
	FileSize      int64                    "fileSize"
	Sharding      map[string]DatabaseStats "raw,omitempty"
}

func getStatusData(host string, port int) (Status, error) {
	var result = Status{}
	result.Databases = make(map[string]DatabaseStats)

	var connectionUrl = fmt.Sprintf("mongodb://%s:%d/?connect=direct", host, port)

	log.Info("Connecting to %s", connectionUrl)
	db, err := mgo.Dial(connectionUrl)
	if err != nil {
		log.Error("Can't connect to mongo: %v", err)
		return Status{}, err
	}
	defer db.Close()

	err = db.Run("serverStatus", &result.ServerStatus)
	if err != nil {
		log.Error("Can't retrieve server status data: %v", err)
		return Status{}, err
	}

	if result.ServerStatus.Process == "mongos" {
		var databases = Databases{}
		err = db.Run("listDatabases", &databases)
		if err != nil {
			log.Error("Can't retrieve database list: %v", err)
			return Status{}, err
		}

		for _, element := range databases.Database {
			if !element.Empty {
				var item = DatabaseStats{}
				database := db.DB(element.Name)
				err = database.Run("dbStats", &item)
				if err != nil {
					log.Error("Can't retrieve database statistics: %v", err)
					return Status{}, err
				}
				result.Databases[element.Name] = item
			}
		}
	}
	return result, nil
}

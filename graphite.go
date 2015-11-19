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
	"github.com/marpaia/graphite-golang"
	"strconv"
	"strings"
)

func sendMetrics(status Status, config Configuration) {
	var Graphite, err = graphite.NewGraphite(config.MetricsHost, config.MetricsPort)
	if err != nil {
		log.Error("Can't connect to graphite collector: %v", err)
		return
	}
	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "connections"), strconv.FormatInt(status.ServerStatus.Connections.Current, 10))

	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "memory.page_faults"), strconv.FormatInt(status.ServerStatus.ExtraInfo.PageFaults, 10))
	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "memory.resident"), strconv.FormatInt(status.ServerStatus.Mem.Resident, 10))
	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "memory.virtual"), strconv.FormatInt(status.ServerStatus.Mem.Virtual, 10))

	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "network.bytes_in"), strconv.FormatInt(status.ServerStatus.Network.BytesIn, 10))
	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "network.bytes_out"), strconv.FormatInt(status.ServerStatus.Network.BytesOut, 10))

	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "operations.insert"), strconv.FormatInt(status.ServerStatus.Opcounters.Insert, 10))
	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "operations.query"), strconv.FormatInt(status.ServerStatus.Opcounters.Query, 10))
	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "operations.update"), strconv.FormatInt(status.ServerStatus.Opcounters.Update, 10))
	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "operations.delete"), strconv.FormatInt(status.ServerStatus.Opcounters.Delete, 10))
	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "operations.getmore"), strconv.FormatInt(status.ServerStatus.Opcounters.GetMore, 10))
	Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "operations.command"), strconv.FormatInt(status.ServerStatus.Opcounters.Command, 10))

	if status.ServerStatus.Process == "mongod" {
		Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "flushes.flushes"), strconv.FormatInt(status.ServerStatus.BackgroundFlushing.Flushes, 10))
		Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "flushes.average_duration_ms"), strconv.FormatFloat(status.ServerStatus.BackgroundFlushing.AverageMS, 'f', -1, 64))

		Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "global_locks.total_time"), strconv.FormatInt(status.ServerStatus.GlobalLocks.TotalTime, 10))
		Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "global_locks.lock_time"), strconv.FormatInt(status.ServerStatus.GlobalLocks.LockTime, 10))
		Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "global_locks.queue.readers"), strconv.FormatInt(status.ServerStatus.GlobalLocks.CurrentQueue.Readers, 10))
		Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "global_locks.queue.writers"), strconv.FormatInt(status.ServerStatus.GlobalLocks.CurrentQueue.Writers, 10))
		Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "global_locks.clients.readers"), strconv.FormatInt(status.ServerStatus.GlobalLocks.ActiveClients.Readers, 10))
		Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.", status.ServerStatus.Process, ".", "global_locks.clients.writers"), strconv.FormatInt(status.ServerStatus.GlobalLocks.ActiveClients.Writers, 10))
	}

	if status.ServerStatus.Process == "mongos" {
		for name, db := range status.Databases {
			Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".", "objects"), strconv.FormatInt(db.Object, 10))
			Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".", "avg_object_size"), strconv.FormatFloat(db.AvgObjectSize, 'f', -1, 64))
			Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".", "data_size"), strconv.FormatInt(db.DataSize, 10))
			Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".", "index_size"), strconv.FormatInt(db.IndexSize, 10))
			Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".", "file_size"), strconv.FormatInt(db.FileSize, 10))
			Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".", "storage_Size"), strconv.FormatInt(db.StorageSize, 10))

			for shard_name, shard := range db.Sharding {
				var s_n = strings.Split(shard_name, "/")[0]
				if len(s_n) > 0 {
					Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".shard.", s_n, ".", "objects"), strconv.FormatInt(shard.Object, 10))
					Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".shard.", s_n, ".", "avg_object_size"), strconv.FormatFloat(shard.AvgObjectSize, 'f', -1, 64))
					Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".shard.", s_n, ".", "data_size"), strconv.FormatInt(shard.DataSize, 10))
					Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".shard.", s_n, ".", "index_size"), strconv.FormatInt(shard.IndexSize, 10))
					Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".shard.", s_n, ".", "file_size"), strconv.FormatInt(shard.FileSize, 10))
					Graphite.SimpleSend(fmt.Sprint(config.MetricsPrefix, ".mongo.database.", name, ".shard.", s_n, ".", "storage_Size"), strconv.FormatInt(shard.StorageSize, 10))
				}
			}
		}
	}
}

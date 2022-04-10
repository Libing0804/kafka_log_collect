package sendToinflux

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"logagent/getSexAble"
	"logagent/influxDB"
	"logagent/models"
	"time"
)

func writesCpuPoints(data *models.CpuInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}

	tags := map[string]string{"cpu": "cpu0"}
	fields := map[string]interface{}{
		"cpu_per": data.CpuPercent,
	}

	pt, err := client.NewPoint("cpu_per", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = influxDB.Cli.Write(bp)
	if err != nil {
		log.Fatalf("cpu info write failed err:%s", err)
	}
}
func writesMemPoints(data *models.MemInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}

	tags := map[string]string{"mem": "mem0"}
	fields := map[string]interface{}{
		"total":       int64(data.Total), //因为不支持uint64 强转一下
		"available":   int64(data.Available),
		"user":        int64(data.Used),
		"userPercent": data.UsedPercent,
		"buffers":     int64(data.Buffers),
		"cached":      int64(data.Cached),
	}

	pt, err := client.NewPoint("memInfo", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = influxDB.Cli.Write(bp)
	if err != nil {
		if err != nil {
			log.Fatalf("mem info write failed err:%s", err)
		}
	}
}
func writesDiskPoints(data *models.DiskInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range data.PartitionUsageStat {
		tags := map[string]string{"path": k}
		fields := map[string]interface{}{
			"totle":             int64(v.Total),
			"free":              int64(v.Free),
			"used":              int64(v.Used),
			"usedPercent":       v.UsedPercent,
			"inodesTotal":       int64(v.InodesTotal),
			"inodesUsed":        int64(v.InodesUsed),
			"inodesFree":        int64(v.InodesFree),
			"inodesUsedPercent": v.InodesUsedPercent,
		}
		pt, err := client.NewPoint("DiskInfo", tags, fields, time.Now())
		if err != nil {
			if err != nil {
				log.Fatalf("disk info write failed err:%s", err)
			}
		}
		bp.AddPoint(pt)
	}

	err = influxDB.Cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
}
func writesNetPoints(data *models.NetInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range data.NetIoCountersStat {
		tags := map[string]string{"name": k}
		fields := map[string]interface{}{
			"bytes_sent_rate":   v.BytesSentRate,
			"bytes_recv_rate":   v.PacketsRecvRate,
			"packets_sent_rate": v.PacketsSentRate,
			"packets_recv_rate": v.PacketsRecvRate,
		}
		pt, err := client.NewPoint("NetInfo", tags, fields, time.Now())
		if err != nil {
			if err != nil {
				log.Fatalf("Net info write failed err:%s", err)
			}
		}
		bp.AddPoint(pt)
	}

	err = influxDB.Cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
}
func run(data *getSexAble.ChanMsgALL) {
	writesCpuPoints(data.Cpu_Info)
	writesMemPoints(data.Mem_Info)
	writesDiskPoints(data.Disk_Info)
	writesNetPoints(data.Net_Info)
}
func SendMsgToDB() {
	//这里要解析信息包  不同类型的数据发送到不同的表格中
	for {
		select {
		case data := <-getSexAble.SexAbleChan:
			go run(data)
		}
	}
}

package models

import (
	"github.com/shirou/gopsutil/disk"
)

type CpuInfo struct {
	CpuPercent float64 `json:"cpu_percent"`
}
type MemInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
}
type UsageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

type DiskInfo struct {
	PartitionUsageStat map[string]*disk.UsageStat
}
type IoStat struct {
	BytesSent       uint64  `json:"bytesSent"`
	BytesRecv       uint64  `json:"bytesRecv"`
	PacketsSent     uint64  `json:"packetsSent"` //
	PacketsRecv     uint64  `json:"packetsRecv"`
	BytesSentRate   float64 `json:"bytes_sent_rate"`   //
	BytesRecvRate   float64 `json:"bytes_recv_rate"`   //
	PacketsSentRate float64 `json:"packets_sent_rate"` //
	PacketsRecvRate float64 `json:"packets_recv_rate"`
}
type NetInfo struct {
	NetIoCountersStat map[string]*IoStat
}

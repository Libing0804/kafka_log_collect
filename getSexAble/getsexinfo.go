package getSexAble

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"logagent/models"
	"time"
)

var (
	LastNetIoStatTimeStamp int64           //上次获取网络数据的时间
	lastNetInfo            *models.NetInfo //上次的数据

)

type ChanMsgALL struct {
	Database string
	//TableName string
	//tags      map[string]string
	//fields    map[string]interface{}
	Cpu_Info  *models.CpuInfo
	Mem_Info  *models.MemInfo
	Disk_Info *models.DiskInfo
	Net_Info  *models.NetInfo
}

var SexAbleChan = make(chan *ChanMsgALL, 1000)

func GetCpuInfo() *models.CpuInfo {
	percent, _ := cpu.Percent(time.Second, false)
	//fmt.Printf("cpu percent:%v\n", percent)
	//写入数据库
	var cpuInfo = &models.CpuInfo{
		CpuPercent: percent[0],
	}
	return cpuInfo
}
func GetMemInfo() *models.MemInfo {
	info, err := mem.VirtualMemory()
	var memInfo = new(models.MemInfo)
	if err != nil {
		log.Fatal("get mem info failed ")
		return memInfo
	}
	//fmt.Printf("cpu percent:%v\n", percent)
	//写入数据库

	memInfo.Total = info.Total
	memInfo.Available = info.Available
	memInfo.UsedPercent = info.UsedPercent
	memInfo.Used = info.Used
	memInfo.Buffers = info.Buffers
	memInfo.Cached = info.Cached

	return memInfo
}

func GetDiskInfo() *models.DiskInfo {
	parts, _ := disk.Partitions(true)
	var diskInfo = &models.DiskInfo{
		PartitionUsageStat: make(map[string]*disk.UsageStat),
	}
	for _, part := range parts {
		//	拿到每个分区
		usageState, err := disk.Usage(part.Mountpoint) //传入挂载点
		if err != nil {
			log.Fatal("disk info get failed err:", err)
			continue
		}
		diskInfo.PartitionUsageStat[part.Mountpoint] = usageState

	}
	return diskInfo

}
func GetNetInfo() *models.NetInfo {
	var netInfo = &models.NetInfo{
		NetIoCountersStat: make(map[string]*models.IoStat),
	}
	netIos, err := net.IOCounters(true)
	if err != nil {
		log.Fatal("net IO  info get failed err:", err)
		return netInfo
	}
	curruntTimeStamp := time.Now().Unix()
	for _, netIo := range netIos {
		var ioStat = new(models.IoStat)
		ioStat.BytesRecv = netIo.BytesSent
		ioStat.BytesSent = netIo.BytesSent
		ioStat.PacketsSent = netIo.PacketsSent
		ioStat.PacketsRecv = netIo.PacketsRecv
		//将具体的网卡数据存起来
		netInfo.NetIoCountersStat[netIo.Name] = ioStat
		//计算相关顺序
		if LastNetIoStatTimeStamp == 0 || lastNetInfo == nil {
			continue
		}
		//计算时间间隔
		interval := curruntTimeStamp - LastNetIoStatTimeStamp
		ioStat.BytesSentRate = (float64(ioStat.BytesSent) - float64(lastNetInfo.NetIoCountersStat[netIo.Name].BytesSent)) / float64(interval)
		ioStat.BytesRecvRate = (float64(ioStat.BytesRecv) - float64(lastNetInfo.NetIoCountersStat[netIo.Name].BytesRecv)) / float64(interval)
		ioStat.PacketsSentRate = (float64(ioStat.PacketsSent) - float64(lastNetInfo.NetIoCountersStat[netIo.Name].PacketsSent)) / float64(interval)
		ioStat.PacketsRecvRate = (float64(ioStat.PacketsRecv) - float64(lastNetInfo.NetIoCountersStat[netIo.Name].PacketsRecv)) / float64(interval)

	}
	LastNetIoStatTimeStamp = curruntTimeStamp //更新时间
	lastNetInfo = netInfo
	return netInfo
}

func Run(interval time.Duration) {
	tickChan := time.Tick(interval)
	for _ = range tickChan {
		cpuinfo := GetCpuInfo()
		meminfo := GetMemInfo()
		diskinfo := GetDiskInfo()
		netinfo := GetNetInfo()
		var msgAll = &ChanMsgALL{
			Database:  "monitor",
			Cpu_Info:  cpuinfo,
			Mem_Info:  meminfo,
			Disk_Info: diskinfo,
			Net_Info:  netinfo,
		}

		SexAbleChan <- msgAll
	}
}

package iplocate

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

//IPInfo IP到地区码的对应关系
type IPInfo struct {
	IPStart uint64
	IPEnd   uint64
	Locate  uint64
}

//Valid 数据是否是有效数据 true 表示有效，否则表示无效
func (i *IPInfo) Valid() bool {
	return i.IPStart <= i.IPEnd
}

//Match 是否匹配 true 表示匹配成功，否则匹配失败
func (i *IPInfo) Match(ip uint64) bool {
	return i.IPStart <= ip && ip <= i.IPEnd
}

//Parse 解析行数据，“startIP,endIP,location”
func (i *IPInfo) Parse(line string) error {
	var items = strings.Split(line, ",")
	if len(items) < 3 {
		return fmt.Errorf("invalid line :[%v]", line)
	}

	if len(items[0]) == 0 || len(items[1]) == 0 || len(items[2]) == 0 {
		return fmt.Errorf("invalid line :[%v]", line)
	}

	locate, err := strconv.ParseUint(items[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid line :[%v] locate:[%v]", line, items[2])
	}

	i.IPStart = parseIPStringToUint64(items[0])
	i.IPEnd = parseIPStringToUint64(items[1])
	i.Locate = locate

	if false == i.Valid() {
		return fmt.Errorf("ip start > end")
	}
	return nil
}

// ----------------------------------------------

var ipMgr ipDepot //全局变量

type ipDepot struct {
	cache []*IPInfo
	lock  sync.RWMutex
}

func (i *ipDepot) set(iplist []*IPInfo) {
	if len(iplist) == 0 {
		return
	}

	//sort first
	var lst = IPList(iplist)
	sort.Slice(iplist, lst.sortLess)

	//cache
	i.lock.Lock()
	defer i.lock.Unlock()
	i.cache = iplist
}

func (i *ipDepot) get(ip string) uint64 {
	ipNum := IPV4ToUint64(ip)
	if ipNum == 0 {
		return DefaultLocation
	}
	return i.get2(ipNum)
}

func (i *ipDepot) get2(ip uint64) uint64 {
	i.lock.RLock()
	defer i.lock.RUnlock()

	high := len(i.cache) - 1
	low := 0

	cnt := 100

	for {
		cnt--
		if cnt < 0 || low > high {
			break
		}

		pos := (high + low) / 2
		info := i.cache[pos]

		if info.IPStart <= ip {
			if info.IPEnd >= ip {
				return info.Locate
			}
			low = pos + 1
		} else {
			high = pos - 1
		}

	}

	return DefaultLocation
}

//IPList ip列表排序
type IPList []*IPInfo

func (l IPList) sortLess(i, j int) bool {
	if l[i].IPStart < l[j].IPStart {
		return true
	}
	return false
}

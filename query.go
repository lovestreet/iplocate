package iplocate

import "fmt"

//DefaultLocation 默认地区码
var DefaultLocation uint64 = 0

//MatchStrict 是否地域码匹配，严格匹配,即数字必须相等
func MatchStrict(locate uint64, locations []uint64) bool {
	for _, item := range locations {
		if item == locate {
			return true
		}
	}
	return false
}

//MatchChild 是否地域码匹配，宽松匹配，如 locate 为河北省，如果locations中包括河北省下的任何地址均为命中
func MatchChild(locate uint64, locations []uint64) bool {
	return matchChildV2(locate, locations)
}

func matchChildV1(locate uint64, locations []uint64) bool {
	var maxsub = getMaxSub(locate)

	for _, item := range locations {
		var regions = getLocateRangeBase(item)
		for _, rg := range regions {
			temp := rg - locate
			if temp >= 0 && temp <= maxsub {
				return true
			}
			fmt.Printf("locate:[%v] maxsub:[%v] rg:[%v] temp:[%v]\n", locate, maxsub, rg, temp)
		}
	}
	return false
}

func matchChildV2(locate uint64, locations []uint64) bool {
	locate = decToHex(locate)
	for _, item := range locations {
		item = decToHex(item)
		var result = locate & item
		fmt.Printf("locate:[%x] item:[%x] locate&item:[%X] match", locate, item, result)

		if result == locate {
			return true
		}
	}
	return false
}

//MatchParent 是否地域码匹配，宽松匹配，如 locate 为河北省/石家庄市/井陉矿区，如果locations中包括“河北省”或“石家庄”均为命中
func MatchParent(locate uint64, locations []uint64) bool {
	return patchParentV2(locate, locations)
}

func patchParentV1(locate uint64, locations []uint64) bool {
	for _, item := range locations {
		var maxsub = getMaxSub(item)
		temp := locate - item
		if temp >= 0 && temp <= maxsub {
			return true
		}
		fmt.Printf("locate:[%v] item:[%v] maxsub:[%v] temp:[%v]\n", locate, item, maxsub, temp)
	}
	return false
}

func patchParentV2(locate uint64, locations []uint64) bool {
	locate = decToHex(locate)
	for _, item := range locations {
		item = decToHex(item)
		var result = locate & item
		fmt.Printf("locate:[%x] item:[%x] locate&item:[%X] match", locate, item, result)

		if result == item {
			return true
		}
	}
	return false
}

//Query 根据IP地址查询地区码
func Query(ip string) (location uint64) {
	location = DefaultLocation

	var ipNum = IPV4ToUint64(ip)
	if ipNum == 0 {
		return
	}

	location = ipMgr.get2(ipNum)
	return
}

//Update 更新缓存
func Update(ipList []*IPInfo) {
	for _, item := range ipList {
		fmt.Printf("%#v\n", *item)
	}
	fmt.Println("update cache count : ", len(ipList))
	ipMgr.set(ipList)
}

//ParseFile 解析IP文件
func ParseFile(filePath string) error {
	//加载文件
	var mgr FileLoader
	lines, err := mgr.LoadFile(filePath)
	if err != nil {
		fmt.Printf("load file:[%v] error:[%v]", filePath, err)
		return err
	}
	fmt.Printf("parse file:[%v] lines:[%v]\n", filePath, len(lines))

	//解析
	var items = make([]*IPInfo, 0, len(lines))
	for _, line := range lines {
		var item = new(IPInfo)
		if err := item.Parse(line); err != nil {
			fmt.Printf("parse line:[%v] error:[%v]\n", line, err)
			continue
		}
		items = append(items, item)
	}

	//更新缓存信息
	Update(items)

	return nil
}

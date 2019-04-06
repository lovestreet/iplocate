package iplocate

import (
	"math"
	"strconv"
	"strings"
)

//解析IP字符串（“202.100.1.100” 或 “3395551588”）为IP数字地址
func parseIPStringToUint64(ip string) (ipNum uint64) {
	if ipNum = IPV4ToUint64(ip); ipNum > 0 {
		return
	}

	var err error
	if ipNum, err = strconv.ParseUint(ip, 10, 64); err == nil {
		return
	}

	return 0
}

// IPV4ToUint64 将IPV4转化为数字
func IPV4ToUint64(ip string) uint64 {
	items := strings.Split(ip, ".")
	if len(items) != 4 {
		return 0
	}
	idx := uint(32)
	ret := uint64(0)
	for _, val := range items {
		idx -= 8
		sect, _ := strconv.ParseUint(val, 10, 64)
		ret += sect << idx
	}
	return ret
}

//根据一个地域码获取该地域所有上级（上级的上级的上级的上级...）编码,规则需要参考每一级编码的位数
func getLocateRangeBase(locate uint64) (region []uint64) {
	var strLocate = strconv.FormatUint(locate, 10)
	if len(strLocate) < 12 {
		return nil
	}
	var country, _ = strconv.ParseUint(strLocate[0:9]+"000", 10, 64)
	var district, _ = strconv.ParseUint(strLocate[0:6]+"000000", 10, 64)
	var city, _ = strconv.ParseUint(strLocate[0:4]+"00000000", 10, 64)
	var province, _ = strconv.ParseUint(strLocate[0:2]+"0000000000", 10, 64)
	region = append(region, country, district, city, province)
	return
}

//根据一个地域码获取下级编码最大值
func getMaxSub(locate uint64) (maxSub uint64) {
	var strLocate = strconv.FormatUint(locate, 10)
	if len(strLocate) < 12 {
		return 0
	}
	var locateBase = strings.TrimRight(strLocate, "0")
	if cnt := len(strLocate) - len(locateBase); cnt > 0 {
		maxSub = uint64(math.Pow10(cnt))
	}
	return
}

//十进制数 转化为 16进制数， 如 123456789 =》0x123456789
func decToHex(value uint64) (ret uint64) {
	var strValue = strconv.FormatUint(value, 10)
	for _, item := range strValue {
		ret = ret + uint64(item-48)
		ret = ret << 4
	}
	ret = ret >> 4
	return ret
}

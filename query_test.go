package iplocate

import (
	"testing"
)

func TestQuery(t *testing.T) {
	var lst = getIPList()

	Update(lst)

	var location uint64
	var strIP string

	strIP = "202.100.1.0"
	location = Query(strIP)
	t.Logf("ip:[%v] => location:[%v]", strIP, location)

	strIP = "202.100.1.100"
	location = Query(strIP)
	t.Logf("ip:[%v] => location:[%v]", strIP, location)

	strIP = "202.100.2.100"
	location = Query(strIP)
	t.Logf("ip:[%v] => location:[%v]", strIP, location)

}

func getIPList() []*IPInfo {
	//202.100.1.0 3395551488
	//202.100.1.100 3395551588
	var items = make([]*IPInfo, 0, 100)
	items = append(items, &IPInfo{IPStart: 3395551488, IPEnd: 3395551500, Locate: 1})
	items = append(items, &IPInfo{IPStart: 3395551588, IPEnd: 3395551600, Locate: 2})
	return items
}

func TestIPV4ToUint64(t *testing.T) {
	var strIP = []string{"202.100.1.100", "8.8.8.8", "114.114.114.114", "192.168.1.1", "10.76.66.98"}
	for _, item := range strIP {
		t.Logf("ip:[%v] => [%v]\n", item, IPV4ToUint64(item))
	}

}

func TestQueryEx(t *testing.T) {
	var ips = []string{
		"202.100.1.100", "8.8.8.8", "114.114.114.114", "192.168.1.1", "10.76.66.98",
		"20.100.1.100", "88.89.8.8", "194.114.184.114", "199.168.1.1", "101.76.66.98",
	}

	if err := ParseFile("data/ip.merge.csv"); err != nil {
		t.Error(err)
		return
	}

	for _, strIP := range ips {
		var location = Query(strIP)
		t.Logf("ip:[%v] => location:[%v]", strIP, location)
	}
}

func TestGetLocateRangeBase(t *testing.T) {
	var locate uint64 = 130107200000
	var region = getLocateRangeBase(locate)
	for _, item := range region {
		t.Logf("%v\n", item)
	}
}

func TestGetLocateBase(t *testing.T) {
	var locate uint64 = 130107200000
	var region = getMaxSub(locate)
	t.Logf("locate %v   parent %v\n", locate, region)
}

func TestMatchChild(t *testing.T) {
	var locate uint64 = 130107210000
	var locations = []uint64{130108222001, 130109200002, 140107200000}

	t.Log(MatchChild(locate, locations))
}

func TestMatchParent(t *testing.T) {
	var locate uint64 = 130107200008
	var locations = []uint64{130107200001, 130107200002, 130107200000}

	t.Log(MatchParent(locate, locations))
}
func TestDecToHex(t *testing.T) {
	var item uint64 = 1234

	var newItem = decToHex(item)
	t.Logf("item:[%v]=>[%x]", item, newItem)
}

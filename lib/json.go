package lib

import (
	"encoding/json"
)

type HostInfo struct {
    Host string
	Attrs map[string]interface{}
}

func (hi *HostInfo) String(asJson ...bool) string {
	if len(asJson) == 1 && asJson[0] {
		res, _ := json.Marshal(hi)
		return string(res)
	}
	return hi.Host
}

func NewHostInfo(host string) HostInfo {
	return HostInfo{
		Host:  host,
		Attrs: make(map[string]interface{}),
	}
}

func HostInfoFromJson(str string) (HostInfo) {
    hi := HostInfo{}
    if err := json.Unmarshal([]byte(str), &hi); err != nil {
        panic(err)
    }
	return hi
}

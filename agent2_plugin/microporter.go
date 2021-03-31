package microporter

import (
	"strconv"

	"encoding/json"

	"zabbix.com/pkg/plugin"
	"zabbix.com/pkg/zbxerr"
)

type Plugin struct {
	plugin.Base
}

var impl Plugin

const (
	keyStreamsDiscovery = "porter.streams.discovery"
	keyStreams = "porter.streams"
	keyStat = "porter.stat"
)

func init() {
    plugin.RegisterMetrics(&impl, "Transporter", keyStat, "Transporter statistics.")
    plugin.RegisterMetrics(&impl, "Transporter", keyStreamsDiscovery, "Transporter streams discovery.")
    plugin.RegisterMetrics(&impl, "Transporter", keyStreams, "Transporter streams list.")
}


func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (res interface{}, err error) {
	val, err := strconv.ParseUint(params[0], 10, 16)
	if err != nil {
		return nil, err
	}
	port := uint16(val)
	
    switch key {
	case keyStat:
		client := NewApiClient(port)

		var stat GetStatisticsResult
		if err = client.GetStatistics(&stat); err != nil {
			return nil, err
		}
		
		result, err := json.Marshal(stat)
		if err != nil {
			return nil, err
		}

		return string(result), nil

	case keyStreamsDiscovery:
		client := NewApiClient(port)

		var streams []StreamInfo
		if err = client.GetStreamList(&streams); err != nil {
			return nil, err
		}

		discovery, err := getStreamDiscovery(streams)
		if err != nil {
			return nil, err
		}
		return string(discovery), nil
	
	case keyStreams:
		client := NewApiClient(port)

		var streams []StreamInfo
		if err = client.GetStreamList(&streams); err != nil {
			return nil, err
		}

		result, err := json.Marshal(&streams)
		if err != nil {
			return nil, err
		}
		return string(result), nil
	}

	return nil, nil
}


// ========================================



type streamDicovery struct {
	Id		uint	`json:"{#ID}"`
	Name	string	`json:"{#NAME}"`
}

func getStreamDiscovery(streams []StreamInfo) (res []byte, err error) {
	streamsDisc := make([]streamDicovery, 0)

	for _, stream := range(streams) {
		streamsDisc = append(streamsDisc, streamDicovery{stream.Id, stream.Name})
	}
	
	if res, err = json.Marshal(&streamsDisc); err != nil {
		return nil, zbxerr.ErrorCannotMarshalJSON.Wrap(err)
	}

	return
}
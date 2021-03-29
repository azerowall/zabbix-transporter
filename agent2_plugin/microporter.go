package microporter

import (
	"strconv"
	"errors"

	"encoding/json"

	"zabbix.com/pkg/plugin"
	"zabbix.com/pkg/zbxerr"
)

type Plugin struct {
	plugin.Base
}

var impl Plugin

const (
	keyStreamDiscovery = "porter.stream.discovery"
	keyStreamInfo = "porter.stream"
	keyStat = "porter.stat"
)

func init() {
    plugin.RegisterMetrics(&impl, "Transporter", keyStat, "Transporter statistics.")
    plugin.RegisterMetrics(&impl, "Transporter", keyStreamDiscovery, "Transporter stream discovery.")
    plugin.RegisterMetrics(&impl, "Transporter", keyStreamInfo, "Transporter stream info.")
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

	case keyStreamDiscovery:
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
	
	case keyStreamInfo:
		name := params[1]

		client := NewApiClient(port)

		var streams []StreamInfo
		if err = client.GetStreamList(&streams); err != nil {
			return nil, err
		}

		for _, stream := range(streams) {
			if stream.Name == name {
				result, err := json.Marshal(&stream)
				if err != nil {
					return nil, err
				}
				return string(result), nil
			}
		}

		return nil, errors.New("Stream not found")
	}

	return nil, nil
}


// ========================================



type streamDicovery struct {
	Name	string	`json:"{#NAME}"`
}

func getStreamDiscovery(streams []StreamInfo) (res []byte, err error) {
	streamsDisc := make([]streamDicovery, 0)

	for _, stream := range(streams) {
		streamsDisc = append(streamsDisc, streamDicovery{stream.Name})
	}
	
	if res, err = json.Marshal(&streamsDisc); err != nil {
		return nil, zbxerr.ErrorCannotMarshalJSON.Wrap(err)
	}

	return
}
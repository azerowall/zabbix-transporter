package microporter

import (
	"strconv"

	"github.com/powerman/rpc-codec/jsonrpc2"
)


type ApiClient struct {
	port uint16
}

func NewApiClient(port uint16) ApiClient {
	return ApiClient { port }
}

func (c *ApiClient) call(method string, args interface{}, reply interface{}) error {
	client, err := jsonrpc2.Dial("tcp", "127.0.0.1:" + strconv.FormatUint(uint64(c.port), 10))
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.Call(method, args, reply)
	if err != nil {
		return err
	}
	
	return nil
}


type GetStatisticsResult struct {
	Bitrate			uint64	`json:"bitrate"`
	BuffersMemUsage	uint	`json:"buffers-mem-usage"`
	CPUUsage		float64	`json:"cpu-usage"`
	MemUsage		uint	`json:"mem-usage"`
	Pid				uint	`json:"pid"`
	StreamingCount	uint	`json:"streaming-count"`
	StreamsCount	uint	`json:"streams-count"`
	ThreadsCount	uint	`json:"threads-count"`
	Uptime			uint64	`json:"uptime"`
}

func (c *ApiClient) GetStatistics(res *GetStatisticsResult) error {
	return c.call("get_statistics", nil, res)
}

type StreamInfo struct {
	Id			uint	`json:"id"`	
	Name		string	`json:"name"`
	Enabled		bool	`json:"enabled"`
	State		string	`json:"state"`
	Uptime		uint	`json:"uptime"`
	Bitrate		uint	`json:"bitrate"`
}

func (c *ApiClient) GetStreamList(res *[]StreamInfo) error {
	return c.call("get_stream_list", nil, res)
}
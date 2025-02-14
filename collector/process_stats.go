package collector

// https://moonraker.readthedocs.io/en/latest/web_api/#get-moonraker-process-stats

import (
	"encoding/json"
	"io/ioutil"

	"net/http"
)

type MoonrakerProcessStatsQueryResponse struct {
	Result struct {
		MoonrakerStats       []MoonrakerProcStats             `json:"moonraker_stats"`
		CpuTemp              float64                          `json:"cpu_temp"`
		Network              map[string]MoonrakerNetworkStats `json:"network"`
		SystemCpuUsage       MoonrakerSystemCpuUsage          `json:"system_cpu_usage"`
		SystemMemory         MoonrakerSystemMemory            `json:"system_memory"`
		SystemUptime         float64                          `json:"system_uptime"`
		WebsocketConnections int                              `json:"websocket_connectsions"`
	} `json:"result"`
}

type MoonrakerProcStats struct {
	Time     float64 `json:"time"`
	CpuUsage float64 `json:"cpu_usage"`
	Memory   int     `json:"memory"`
	MemUnits string  `json:"mem_units"`
}

type MoonrakerNetworkStats struct {
	RxBytes   int64   `json:"rx_bytes"`
	TxBytes   int64   `json:"tx_bytes"`
	RxPackets int64   `json:"rx_packets"`
	TxPackets int64   `json:"tx_packets"`
	RxErrs    int     `json:"rx_errs"`
	TxErrs    int     `json:"tx_errs"`
	RxDrop    int     `json:"rx_drop"`
	TxDrop    int     `json:"tx_drop"`
	Bandwidth float64 `json:"bandwidth"`
}

type MoonrakerSystemCpuUsage struct {
	Cpu  float64 `json:"cpu"`
	Cpu0 float64 `json:"cpu0"`
	Cpu1 float64 `json:"cpu1"`
	Cpu2 float64 `json:"cpu2"`
	Cpu3 float64 `json:"cpu3"`
}

type MoonrakerSystemMemory struct {
	Total     int `json:"total"`
	Available int `json:"available"`
	Used      int `json:"used"`
}

func (c collector) fetchMoonrakerProcessStats(klipperHost string) (*MoonrakerProcessStatsQueryResponse, error) {
	var procStatsUrl = "http://" + klipperHost + "/machine/proc_stats"
	c.logger.Debug("Collecting metrics from " + procStatsUrl)
	res, err := http.Get(procStatsUrl)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.logger.Fatal(err)
		return nil, err
	}

	var response MoonrakerProcessStatsQueryResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		c.logger.Fatal(err)
		return nil, err
	}

	return &response, nil
}

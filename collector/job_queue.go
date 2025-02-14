package collector

// https://moonraker.readthedocs.io/en/latest/web_api/#retrieve-the-job-queue-status

import (
	"encoding/json"
	"io/ioutil"

	"net/http"
)

type MoonrakerJobQueueResponse struct {
	Result struct {
		QueuedJobs []MoonrakerQueuedJob `json:"queued_jobs"`
		QueueState string               `json:"queue_state"`
	} `json:"result"`
}

type MoonrakerQueuedJob struct {
	TimeInQueue float64 `json:"time_in_queue"`
}

func (c collector) fetchMoonrakerJobQueue(klipperHost string) (*MoonrakerJobQueueResponse, error) {
	var procStatsUrl = "http://" + klipperHost + "/server/job_queue/status"
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

	var response MoonrakerJobQueueResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		c.logger.Fatal(err)
		return nil, err
	}

	return &response, nil
}

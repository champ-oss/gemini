package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	RetryDelaySeconds = 30
	RetryAttempts     = 20
)

func checkExpectedGrafanaTableCount(hostname, username, password, table string, minExpected int64) error {
	for i := 0; ; i++ {
		fmt.Printf("querying Grafana to get count for table: %s\n", table)
		count := getGrafanaTableCount(hostname, username, password, table)
		fmt.Printf("%s table count: %d\n", table, count)
		if count >= minExpected {
			return nil
		}

		if i >= (RetryAttempts - 1) {
			panic("timed out while retrying")
		}

		fmt.Printf("count was less than expected. retrying in %d seconds...\n", RetryDelaySeconds)
		time.Sleep(time.Second * RetryDelaySeconds)
	}
}

// getGrafanaTableCount queries grafana to return the count of records in the given table
func getGrafanaTableCount(hostname, username, password, table string) int64 {
	grafanaApiCountRequest := buildGrafanaApiCountRequest(table)
	responseBody := sendGrafanaApiQueryRequest(hostname, username, password, grafanaApiCountRequest)
	return parseGrafanaApiCountResponse(responseBody)
}

// buildGrafanaApiCountRequest creates a Grafana query to get a table count
func buildGrafanaApiCountRequest(table string) []byte {
	type GrafanaQuery struct {
		IntervalMs    int64  `json:"intervalMs"`
		MaxDataPoints int    `json:"maxDataPoints"`
		DatasourceId  int    `json:"datasourceId"`
		RawSql        string `json:"rawSql"`
		Format        string `json:"format"`
	}

	type GrafanaQueryRequest struct {
		From    string          `json:"from"`
		To      string          `json:"to"`
		Queries []*GrafanaQuery `json:"queries"`
	}

	queryRequest := &GrafanaQueryRequest{
		From: "now-10y",
		To:   "now",
		Queries: []*GrafanaQuery{
			{
				IntervalMs:    int64(86400000),
				MaxDataPoints: 1000,
				DatasourceId:  1,
				RawSql:        fmt.Sprintf("SELECT count(*) FROM %s", table),
				Format:        "table",
			},
		},
	}

	queryRequestJson, err := json.Marshal(queryRequest)
	if err != nil {
		panic(err)
	}
	return queryRequestJson
}

// sendGrafanaApiQueryRequest creates and sends an HTTP POST request to the Grafana server and returns the response body
func sendGrafanaApiQueryRequest(hostname, username, password string, requestBody []byte) []byte {
	queryUrl := fmt.Sprintf("https://%s:%s@%s/api/ds/query", username, password, hostname)
	fmt.Printf("sending HTTP POST to https://%s:%s@%s/api/ds/query\n", username, "*****", hostname)

	req, err := http.NewRequest("POST", queryUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Content-Type", `application/json`)

	c := http.Client{Timeout: time.Duration(10) * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("query returned: %s message: %s\n", resp.Status, body)
	return body
}

// parseGrafanaApiCountResponse parses the raw Grafana API query response to return the data value
// Example response (return value would be 256):
// {"results":{"A":{"frames":[{"schema":{"refId":"A","meta":{"executedQueryString":"SELECT count(*) FROM commits"},
// "fields":[{"name":"count(*)","type":"number","typeInfo":{"frame":"int64","nullable":true}}]},"data":{"values":[[256]]}}]}}}
func parseGrafanaApiCountResponse(apiResponseBody []byte) int64 {
	var grafanaQueryResponse struct {
		Results struct {
			A struct {
				Frames []struct {
					Data struct {
						Values [][]int64
					}
				}
			}
		}
	}

	if err := json.Unmarshal(apiResponseBody, &grafanaQueryResponse); err != nil {
		panic(err)
	}
	if len(grafanaQueryResponse.Results.A.Frames) == 0 {
		return 0
	}
	if len(grafanaQueryResponse.Results.A.Frames[0].Data.Values) == 0 {
		return 0
	}
	if len(grafanaQueryResponse.Results.A.Frames[0].Data.Values[0]) == 0 {
		return 0
	}
	return grafanaQueryResponse.Results.A.Frames[0].Data.Values[0][0]
}

package cloudflare_worker_exporter

import (
	"context"
	"fmt"
	"time"

	"github.com/machinebox/graphql"
)

type ResponseStruct struct {
	Viewer struct {
		Accounts []struct {
			WorkersInvocationsAdaptive []struct {
				Dimensions struct {
					ScriptName   string    `json:"scriptName"`
					Status       string    `json:"status"`
					Date         time.Time `json:"date"`
					DateTime     time.Time `json:"datetime"`
					DateTimeHour time.Time `json:"datetimeHour"`
				} `json:"dimensions"`
				Quantiles struct {
					CPUTimeP25  float64 `json:"cpuTimeP25"`
					CPUTimeP50  float64 `json:"cpuTimeP50"`
					CPUTimeP75  float64 `json:"cpuTimeP75"`
					CPUTimeP90  float64 `json:"cpuTimeP90"`
					CPUTimeP99  float64 `json:"cpuTimeP99"`
					CPUTimeP999 float64 `json:"cpuTimeP999"`
				} `json:"quantiles"`
				Sum struct {
					Errors      int `json:"errors"`
					Requests    int `json:"requests"`
					Subrequests int `json:"subrequests"`
				} `json:"sum"`
			} `json:"workersInvocationsAdaptive"`
		} `json:"accounts"`
	} `json:"viewer"`
}

type Fetcher struct {
	cloudflareEndpoint, cloudflareToken, cloudFlareAccountId string
}

func NewFetcher(cloudflareEndpoint string, cloudflareToken string, cloudFlareAccountId string) *Fetcher {
	return &Fetcher{
		cloudflareEndpoint:  cloudflareEndpoint,
		cloudflareToken:     cloudflareToken,
		cloudFlareAccountId: cloudFlareAccountId,
	}
}

func (fetcher *Fetcher) getRequest(request string) *graphql.Request {
	req := graphql.NewRequest(request)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", fetcher.cloudflareToken))

	return req
}

func (fetcher *Fetcher) fetchMetrics(request string, startTime string, endTime string) (respData ResponseStruct, err error) {
	client := graphql.NewClient(fetcher.cloudflareEndpoint)

	req := fetcher.getRequest(request)
	req.Var("account", fetcher.cloudFlareAccountId)
	req.Var("start", startTime)
	req.Var("end", endTime)

	ctx := context.Background()
	err = client.Run(ctx, req, &respData)

	return
}

func WeekStart(year, week int) time.Time {
	t := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 0, -int(t.Weekday()))
	t = t.AddDate(0, 0, week*7)

	return t
}

// fetch the request count since the begining of the week
func (fetcher *Fetcher) FetchRequestCount() (ResponseStruct, error) {
	now := time.Now()
	year, week := now.ISOWeek()
	startTime := WeekStart(year, week).Format(time.RFC3339)
	endTime := now.UTC().Format(time.RFC3339)

	request := `
		query ($account: String!, $start: String!, $end: String!) {
			viewer {
				accounts(filter: {accountTag: $account}) {
					workersInvocationsAdaptive(limit: 100, filter: {datetime_geq: $start, datetime_leq: $end}) {
						sum {
							subrequests
							requests
							errors
						}
						dimensions {
							scriptName
							status
						}
					}
				}
			}
		}`

	return fetcher.fetchMetrics(request, startTime, endTime)
}

// fetch the cpu time in the last minute
func (fetcher *Fetcher) FetchCpuTime() (ResponseStruct, error) {
	now := time.Now()
	startTime := now.Add(-time.Minute * 1).Format(time.RFC3339)
	endTime := now.UTC().Format(time.RFC3339)

	request := `
		query ($account: String!, $start: String!, $end: String!) {
			viewer {
				accounts(filter: {accountTag: $account}) {
					workersInvocationsAdaptive(limit: 100, filter: {datetime_geq: $start, datetime_leq: $end}) {
						quantiles {
              cpuTimeP25
              cpuTimeP50
              cpuTimeP75
              cpuTimeP90
              cpuTimeP99
              cpuTimeP999
            }
						dimensions {
							scriptName
							status
						}
					}
				}
			}
		}`

	return fetcher.fetchMetrics(request, startTime, endTime)
}

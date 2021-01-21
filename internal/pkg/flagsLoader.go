package cloudflare_worker_exporter

import (
	"flag"
	"os"
)

type Flags struct {
	ListenAddress       *string
	MetricsPath         *string
	CloudflareToken     *string
	CloudFlareAccountId *string
	CloudflareEndpoint  *string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

const cloudflareAnalyticsEndpoint = "https://api.cloudflare.com/client/v4/graphql/"

func LoadFlags() Flags {
	flags := Flags{
		ListenAddress: flag.String("port", getEnv("PORT", "9184"),
			"HTTP server port, default 9184 override using PORT environment variable"),
		MetricsPath: flag.String("path", getEnv("MERTICS_ENDPOINT", "/metrics"),
			"Path for metrics route, default /metrics override using MERTICS_ENDPOINT environment variable"),
		CloudflareToken: flag.String("token", os.Getenv("CLOUDFLARE_ANALYTICS_TOKEN"),
			"Cloudflare API bearer token with Account.Account Analytics permissions, default is the value from CLOUDFLARE_ANALYTICS_TOKEN environment variable"),
		CloudFlareAccountId: flag.String("account", os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
			"Cloudflare account id, default is the value from CLOUDFLARE_ACCOUNT_ID environment variable"),
		CloudflareEndpoint: flag.String("endpoint", getEnv("CLOUDFLARE_ANALYTICS_ENDPOINT", cloudflareAnalyticsEndpoint),
			"Cloudflare account id, default is "+cloudflareAnalyticsEndpoint),
	}
	flag.Parse()
	return flags
}

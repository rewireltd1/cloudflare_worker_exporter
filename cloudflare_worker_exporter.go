package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	exporter "cloudflare_worker_exporter/internal/pkg"
	fetcher "cloudflare_worker_exporter/internal/pkg"
	flagsLoader "cloudflare_worker_exporter/internal/pkg"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func loadEnv() {
	var err error
	if value, ok := os.LookupEnv("DOTENV_FILE"); ok {
		log.Println("Using ", value, "dot env file")
		err = godotenv.Load(value)
	} else {
		log.Println("Using default dot env file")
		err = godotenv.Load()
	}

	if err != nil {
		log.Println("Failed to load .env file, getting default env from os.")
	}
}

func main() {
	loadEnv()
	flags := flagsLoader.LoadFlags()

	fetcher := fetcher.NewFetcher(*(flags.CloudflareEndpoint), *flags.CloudflareToken, *flags.CloudFlareAccountId)
	exporter := exporter.NewExporter(fetcher)
	prometheus.MustRegister(exporter)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	http.Handle(*flags.MetricsPath, promhttp.Handler())

	log.Fatal(http.ListenAndServe(":"+*flags.ListenAddress, nil))
}

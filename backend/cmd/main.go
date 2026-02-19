package main

import (
	"log"
	"time"

	"github.com/jamesjunior/eduops-backend/internal/collector"
)

func main() {

	dockerCollector, err := collector.NewDockerCollector()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(5 * time.Second)

	log.Println("Collector iniciado...")

	for range ticker.C {
		metrics, err := dockerCollector.Collect()
		if err != nil {
			log.Println("Erro na coleta:", err)
			continue
		}

		for _, m := range metrics {
			log.Printf("Container: %s CPU: %.2f%% MEM: %.2f%% Restarts: %d\n",
				m.Name, m.CPUPercent, m.MemoryPercent, m.RestartCount)
		}
	}
}

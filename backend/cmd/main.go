package main

import (
	"log"
	"time"

	"github.com/jamesjunior/eduops-backend/internal/collector"
	"github.com/jamesjunior/eduops-backend/internal/persistence"
)

func main() {

	db, err := persistence.InitDB("data/metrics.db")
	if err != nil {
		log.Fatal(err)
	}

	repo := persistence.NewMetricsRepository(db)

	dockerCollector, err := collector.NewDockerCollector()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(5 * time.Second)

	log.Println("Collector iniciado")

	for range ticker.C {
		metrics, err := dockerCollector.Collect()
		if err != nil {
			log.Println("Erro na coleta:", err)
			continue
		}

		for _, m := range metrics {
			// Salva métrica no banco de dados
			if err := repo.Save(m); err != nil {
				log.Printf("Erro ao salvar métrica: %v\n", err)
			}

			log.Printf("Container: %s CPU: %.2f%% MEM: %.2f%% Restarts: %d\n",
				m.Name, m.CPUPercent, m.MemoryPercent, m.RestartCount)
		}
	}
}

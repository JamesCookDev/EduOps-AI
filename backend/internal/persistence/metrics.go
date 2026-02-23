package persistence

import (
	"database/sql"
	"github.com/jamesjunior/eduops-backend/internal/collector"
)

type MetricsRepository struct {
	db *sql.DB
}

func NewMetricsRepository(db *sql.DB) *MetricsRepository {
	return &MetricsRepository{db: db}
}

func (r *MetricsRepository) Save(metric collector.ContainerMetrics) error {

	query := `
	INSERT INTO container_metrics (
		container_id,
		name,
		cpu_percent,
		memory_usage,
		memory_limit,
		memory_percent,
		restart_count,
		timestamp
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		metric.ContainerID,
		metric.Name,
		metric.CPUPercent,
		metric.MemoryUsage,
		metric.MemoryLimit,
		metric.MemoryPercent,
		metric.RestartCount,
		metric.Timestamp,
	)

	return err
}
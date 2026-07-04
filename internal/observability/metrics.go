// Package observability concentra as métricas Prometheus do serviço.
package observability

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Metrics agrupa os coletores da aplicação.
type Metrics struct {
	rpcHandled  *prometheus.CounterVec
	rpcDuration *prometheus.HistogramVec

	dbQueriesTotal  *prometheus.CounterVec
	dbQueryDuration *prometheus.HistogramVec
	dbRowsReturned  *prometheus.HistogramVec
}

// NewMetrics registra os coletores no registrador default. Deve ser chamada uma vez.
func NewMetrics() *Metrics {
	return &Metrics{
		rpcHandled: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "grpc_server_handled_total",
			Help: "Total de RPCs finalizadas, por método e código de status.",
		}, []string{"method", "code"}),
		rpcDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "grpc_server_handling_seconds",
			Help:    "Latência das RPCs em segundos, por método.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method"}),
		dbQueriesTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total de consultas ao banco, por nome e status (ok/error).",
		}, []string{"query", "status"}),
		dbQueryDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duração das consultas ao banco em segundos, por nome.",
			Buckets: prometheus.DefBuckets,
		}, []string{"query"}),
		dbRowsReturned: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "db_rows_returned",
			Help:    "Número de linhas retornadas por consulta.",
			Buckets: []float64{0, 1, 5, 10, 25, 50, 100, 250, 500, 1000},
		}, []string{"query"}),
	}
}

// RecordQuery registra as métricas de uma consulta ao banco.
func (m *Metrics) RecordQuery(query string, start time.Time, rows int, err error) {
	statusLabel := "ok"
	if err != nil {
		statusLabel = "error"
	}
	m.dbQueriesTotal.WithLabelValues(query, statusLabel).Inc()
	m.dbQueryDuration.WithLabelValues(query).Observe(time.Since(start).Seconds())
	m.dbRowsReturned.WithLabelValues(query).Observe(float64(rows))
}

// UnaryServerInterceptor instrumenta cada RPC unária com contagem e latência.
func (m *Metrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err)
		m.rpcHandled.WithLabelValues(info.FullMethod, code.String()).Inc()
		m.rpcDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())
		return resp, err
	}
}

// RegisterPoolCollector expõe estatísticas do pool pgx (conexões totais/ociosas/em uso).
func RegisterPoolCollector(pool *pgxpool.Pool) {
	prometheus.MustRegister(newPoolCollector(pool))
}

type poolCollector struct {
	pool     *pgxpool.Pool
	total    *prometheus.Desc
	acquired *prometheus.Desc
	idle     *prometheus.Desc
	maxConns *prometheus.Desc
}

func newPoolCollector(pool *pgxpool.Pool) *poolCollector {
	return &poolCollector{
		pool:     pool,
		total:    prometheus.NewDesc("pgxpool_total_conns", "Conexões totais no pool.", nil, nil),
		acquired: prometheus.NewDesc("pgxpool_acquired_conns", "Conexões em uso.", nil, nil),
		idle:     prometheus.NewDesc("pgxpool_idle_conns", "Conexões ociosas.", nil, nil),
		maxConns: prometheus.NewDesc("pgxpool_max_conns", "Máximo de conexões configurado.", nil, nil),
	}
}

func (c *poolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.acquired
	ch <- c.idle
	ch <- c.maxConns
}

func (c *poolCollector) Collect(ch chan<- prometheus.Metric) {
	s := c.pool.Stat()
	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(s.TotalConns()))
	ch <- prometheus.MustNewConstMetric(c.acquired, prometheus.GaugeValue, float64(s.AcquiredConns()))
	ch <- prometheus.MustNewConstMetric(c.idle, prometheus.GaugeValue, float64(s.IdleConns()))
	ch <- prometheus.MustNewConstMetric(c.maxConns, prometheus.GaugeValue, float64(s.MaxConns()))
}

package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	onceConfig    sync.Once
	metricConfigs = &Metrics{}
)

type Metrics struct {
	EventTrackingListener *prometheus.CounterVec
	CampaingRepository    *prometheus.CounterVec
	MerchantRepository    *prometheus.CounterVec
	SlugRepository        *prometheus.CounterVec
	UserRepository        *prometheus.CounterVec
	Transaction           *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	onceConfig.Do(func() {
		metricConfigs.EventTrackingListener = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "event_tracking_listener",
			Help: "Process messages SQS queue_campaing",
		}, []string{"status", "details"})
		prometheus.MustRegister(metricConfigs.EventTrackingListener)
		metricConfigs.CampaingRepository = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "campaing_repository",
			Help: "CRUD campaing repository",
		}, []string{"method", "status", "details"})
		prometheus.MustRegister(metricConfigs.CampaingRepository)
		metricConfigs.MerchantRepository = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "merchant_repository",
			Help: "CRUD merchant repository",
		}, []string{"method", "status", "details"})
		prometheus.MustRegister(metricConfigs.MerchantRepository)
		metricConfigs.SlugRepository = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "slug_repository",
			Help: "CRUD slug repository",
		}, []string{"method", "status", "details"})
		prometheus.MustRegister(metricConfigs.SlugRepository)
		metricConfigs.UserRepository = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "user_repository",
			Help: "CRUD user repository",
		}, []string{"method", "status", "details"})
		prometheus.MustRegister(metricConfigs.UserRepository)
		metricConfigs.Transaction = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "transaction",
			Help: "Create, update and delete in repository",
		}, []string{"method", "status", "details"})
		prometheus.MustRegister(metricConfigs.Transaction)
	})
	return metricConfigs
}

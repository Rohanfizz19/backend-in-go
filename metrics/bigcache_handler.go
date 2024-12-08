package metrics

import (
	"github.com/allegro/bigcache/v3"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type BigCacheMetrics struct {
	HitsGauge        *prometheus.GaugeVec
	MissesGauge      *prometheus.GaugeVec
	DelHitsGauge     *prometheus.GaugeVec
	DelMissesGauge   *prometheus.GaugeVec
	CollisionsGauge  *prometheus.GaugeVec
	BytesStoredGauge *prometheus.GaugeVec
	ItemsStoredGauge *prometheus.GaugeVec
}

func InitBigCacheMetrics() BigCacheMetrics {
	hitsGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "bigcache_hits_total",
			Help:      "Number of cache hits.",
		}, []string{},
	)
	missesGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "bigcache_misses_total",
			Help:      "Number of cache misses.",
		}, []string{},
	)
	delHitsGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "bigcache_del_hits_total",
			Help:      "Number of cache delete hits.",
		}, []string{},
	)
	delMissesGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "bigcache_del_misses_total",
			Help:      "Number of cache delete misses.",
		}, []string{},
	)
	collisionsGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "bigcache_collisions_total",
			Help:      "Number of cache collisions.",
		}, []string{},
	)

	bytesStoredGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "bigcache_bytes_stored_total",
			Help:      "Number of bytes stored.",
		}, []string{},
	)

	itemsStoredGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "bigcache_items_stored_total",
			Help:      "Number of items stored.",
		}, []string{},
	)

	prometheus.MustRegister(hitsGauge)
	prometheus.MustRegister(missesGauge)
	prometheus.MustRegister(delHitsGauge)
	prometheus.MustRegister(delMissesGauge)
	prometheus.MustRegister(collisionsGauge)
	prometheus.MustRegister(bytesStoredGauge)
	prometheus.MustRegister(itemsStoredGauge)

	return BigCacheMetrics{
		HitsGauge:        hitsGauge,
		MissesGauge:      missesGauge,
		DelHitsGauge:     delHitsGauge,
		DelMissesGauge:   delMissesGauge,
		CollisionsGauge:  collisionsGauge,
		BytesStoredGauge: bytesStoredGauge,
		ItemsStoredGauge: itemsStoredGauge,
	}
}

func CollectStatistics(cache *bigcache.BigCache, metrics BigCacheMetrics) {
	ticker := time.NewTicker(time.Minute)

	for range ticker.C {

		cacheStats := cache.Stats()

		metrics.HitsGauge.WithLabelValues().Set(float64(cacheStats.Hits))
		metrics.MissesGauge.WithLabelValues().Set(float64(cacheStats.Misses))
		metrics.DelHitsGauge.WithLabelValues().Set(float64(cacheStats.DelHits))
		metrics.DelMissesGauge.WithLabelValues().Set(float64(cacheStats.DelMisses))
		metrics.CollisionsGauge.WithLabelValues().Set(float64(cacheStats.Collisions))
		metrics.BytesStoredGauge.WithLabelValues().Set(float64(cache.Capacity()))
		metrics.ItemsStoredGauge.WithLabelValues().Set(float64(cache.Len()))
	}
}

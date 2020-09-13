package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

type Stat struct {
	Count    int64
	Duration time.Duration
}

type Stats map[string]*Stat

type GlobalStats struct {
	Lock  sync.Mutex
	Stats Stats
}

var globalStats = GlobalStats{
	Stats: make(Stats),
}

type Profiler struct {
	GroupName      string
	GroupStartTime time.Time

	CheckpointName      string
	CheckpointStartTime time.Time

	Stats Stats
}

func NewProfiler(name string) *Profiler {
	return &Profiler{
		GroupName:           name,
		GroupStartTime:      time.Now(),
		CheckpointName:      "init",
		CheckpointStartTime: time.Now(),
		Stats:               make(Stats),
	}
}

func (p *Profiler) Close() {
	p.Checkpoint("close")

	key := p.GroupName
	s := p.Stats[key]
	if s == nil {
		s = &Stat{}
		p.Stats[key] = s
	}
	s.Count++
	s.Duration += time.Since(p.GroupStartTime)

	globalStats.Lock.Lock()
	defer globalStats.Lock.Unlock()
	for n, s := range p.Stats {
		gs := globalStats.Stats[n]
		if gs == nil {
			gs = &Stat{}
			globalStats.Stats[n] = gs
		}
		gs.Count += s.Count
		gs.Duration += s.Duration
		//metrics.AddSampleWithLabels([]string{"profiler_sample"},
		//	float32(s.Duration)/float32(time.Microsecond),
		//	[]metrics.Label{{Name: "name", Value: n}})
		//metrics.IncrCounterWithLabels([]string{"profiler_counter"},
		//	float32(s.Duration)/float32(time.Microsecond),
		//	[]metrics.Label{{Name: "name", Value: n}})
	}
}

func (p *Profiler) Checkpoint(block string) {
	key := p.GroupName + "." + p.CheckpointName
	s := p.Stats[key]
	if s == nil {
		s = &Stat{}
		p.Stats[key] = s
	}
	s.Count++
	s.Duration += time.Since(p.CheckpointStartTime)

	p.CheckpointName = block
	p.CheckpointStartTime = time.Now()
}

func handleProfilerStats(w http.ResponseWriter, r *http.Request) {
	stats := func() []string {
		globalStats.Lock.Lock()
		defer globalStats.Lock.Unlock()
		stats := make([]string, 0)
		for n, s := range globalStats.Stats {
			stats = append(stats, fmt.Sprintf(
				"%s\t%d\t%d\t%.3f\n",
				n, s.Duration/time.Second,
				s.Count,
				float64(s.Duration)/float64(time.Second)/float64(s.Count)))
		}
		return stats
	}()
	sort.Strings(stats)
	w.Header().Set("Content-Type", "text/plain")
	for _, s := range stats {
		w.Write([]byte(s))
	}
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(os.Getenv("DOCKER_TAG")))
}

func init() {
	http.HandleFunc("/systemz/stats", handleProfilerStats)
	http.HandleFunc("/systemz/version", handleVersion)
	go http.ListenAndServe("0.0.0.0:9003", nil)
}

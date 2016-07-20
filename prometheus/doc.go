// Copyright 2014 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package prometheus provides metrics primitives to instrument code for
// monitoring. It also offers a registry for metrics and ways to expose
// registered metrics via an HTTP endpoint or push them to a Pushgateway.
//
// All exported functions and methods are safe to be used concurrently unless
// specified otherwise.
//
// A Basic Example
//
// As a starting point, a very basic usage example:
//
//    package main
//
//    import (
//    	"net/http"
//
//    	"github.com/prometheus/client_golang/prometheus"
//    )
//
//    var (
//    	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
//    		Name: "cpu_temperature_celsius",
//    		Help: "Current temperature of the CPU.",
//    	})
//    	hdFailures = prometheus.NewCounter(prometheus.CounterOpts{
//    		Name: "hd_errors_total",
//    		Help: "Number of hard-disk errors.",
//    	})
//    )
//
//    func init() {
//    	// Metrics have to be registered to be exposed:
//    	prometheus.MustRegister(cpuTemp)
//    	prometheus.MustRegister(hdFailures)
//    }
//
//    func main() {
//    	cpuTemp.Set(65.3)
//    	hdFailures.Inc()
//
//    	// The Handler function provides a default handler to expose metrics
//    	// via an HTTP server. "/metrics" is the usual endpoint for that.
//    	http.Handle("/metrics", prometheus.Handler())
//    	http.ListenAndServe(":8080", nil)
//    }
//
//
// This is a complete program that exports two metrics, a Gauge and a Counter.
// It also exports some stats about the HTTP usage of the /metrics
// endpoint. (See the Handler function for more detail.)
//
// TODO: Rework from here on. Use titles
//
// Two more advanced metric types are the Summary and Histogram. A more
// thorough description of metric types can be found in the prometheus docs:
// https://prometheus.io/docs/concepts/metric_types/
//
// In addition to the fundamental metric types Gauge, Counter, Summary, and
// Histogram, a very important part of the Prometheus data model is the
// partitioning of samples along dimensions called labels, which results in
// metric vectors. The fundamental types are GaugeVec, CounterVec, SummaryVec,
// and HistogramVec.
//
// Those are all the parts needed for basic usage. Detailed documentation and
// examples are provided below.
//
// Everything else this package and its sub-packages offer is essentially for
// "power users" only. A few pointers to "power user features":
//
// All the various ...Opts structs have a ConstLabels field for labels that
// never change their value (which is only useful under special circumstances,
// see documentation of the Opts type).
//
// The Untyped metric behaves like a Gauge, but signals the Prometheus server
// not to assume anything about its type.
//
// For custom metric collection, there are two entry points: Custom Metric
// implementations and custom Collector implementations. A Metric is the
// fundamental unit in the Prometheus data model: a sample at a point in time
// together with its meta-data (like its fully-qualified name and any number of
// pairs of label name and label value) that knows how to marshal itself into a
// data transfer object (aka DTO, implemented as a protocol buffer). A Collector
// gets registered with the Prometheus registry and manages the collection of
// one or more Metrics. Many parts of this package are building blocks for
// Metrics and Collectors. Desc is the metric descriptor, actually used by all
// metrics under the hood, and by Collectors to describe the Metrics to be
// collected, but only to be dealt with by users if they implement their own
// Metrics or Collectors. To create a Desc, the BuildFQName function will come
// in handy. Other useful components for Metric and Collector implementation
// include: LabelPairSorter to sort the DTO version of label pairs,
// NewConstMetric and MustNewConstMetric to create "throw away" Metrics at
// collection time, MetricVec to bundle custom Metrics into a metric vector
// Collector, SelfCollector to make a custom Metric collect itself.
//
// A good example for a custom Collector is the expvarCollector included in this
// package, which exports variables exported via the "expvar" package as
// Prometheus metrics.
//
// The functions Register, Unregister, MustRegister, RegisterOrGet, and
// MustRegisterOrGet all act on the default registry. They wrap other calls as
// described in their doc comment. For advanced use cases, you can work with
// custom registries (created by NewRegistry and similar) and call the wrapped
// functions directly.
//
// The functions Handler and UninstrumentedHandler create an HTTP handler to
// serve metrics from the default registry in the default way, which covers most
// of the use cases. With HandlerFor, you can create a custom HTTP handler for
// custom registries.
//
// The functions Push and PushAdd push the metrics from the default registry via
// HTTP to a Pushgateway. With PushFrom and PushAddFrom, you can push the
// metrics from custom registries. However, often you just want to push a
// handfull of Collectors only. For that case, there are the convenience
// functions PushCollectors and PushAddCollectors.
package prometheus

// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loadscraper

import (
	"context"
	"errors"
	"testing"

	"github.com/shirou/gopsutil/load"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/internal/dataold"
	"go.opentelemetry.io/collector/receiver/hostmetricsreceiver/internal"
)

func TestScrapeMetrics(t *testing.T) {
	type testCase struct {
		name        string
		loadFunc    func() (*load.AvgStat, error)
		expectedErr string
	}

	testCases := []testCase{
		{
			name: "Standard",
		},
		{
			name:        "Load Error",
			loadFunc:    func() (*load.AvgStat, error) { return nil, errors.New("err1") },
			expectedErr: "err1",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			scraper := newLoadScraper(context.Background(), zap.NewNop(), &Config{})
			if test.loadFunc != nil {
				scraper.load = test.loadFunc
			}

			err := scraper.Initialize(context.Background())
			require.NoError(t, err, "Failed to initialize load scraper: %v", err)
			defer func() { assert.NoError(t, scraper.Close(context.Background())) }()

			metrics, err := scraper.ScrapeMetrics(context.Background())
			if test.expectedErr != "" {
				assert.EqualError(t, err, test.expectedErr)
				return
			}
			require.NoError(t, err, "Failed to scrape metrics: %v", err)

			// expect 3 metrics
			assert.Equal(t, 3, metrics.Len())

			// expect a single datapoint for 1m, 5m & 15m load metrics
			assertMetricHasSingleDatapoint(t, metrics.At(0), loadAvg1MDescriptor)
			assertMetricHasSingleDatapoint(t, metrics.At(1), loadAvg5mDescriptor)
			assertMetricHasSingleDatapoint(t, metrics.At(2), loadAvg15mDescriptor)

			internal.AssertSameTimeStampForAllMetrics(t, metrics)
		})
	}
}

func assertMetricHasSingleDatapoint(t *testing.T, metric dataold.Metric, descriptor dataold.MetricDescriptor) {
	internal.AssertDescriptorEqual(t, descriptor, metric.MetricDescriptor())
	assert.Equal(t, 1, metric.DoubleDataPoints().Len())
}

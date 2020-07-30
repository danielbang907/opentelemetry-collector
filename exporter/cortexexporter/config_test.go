// Copyright 2020 The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cortexexporter

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/config/configtest"
	"go.opentelemetry.io/collector/config/configtls"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.ExampleComponents()
	assert.NoError(t, err)

	factory := &Factory{}
	factories.Exporters[typeStr] = factory
	cfg, err := configtest.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	e0 := cfg.Exporters["cortex"]
	assert.Equal(t, e0, factory.CreateDefaultConfig())

	e1 := cfg.Exporters["cortex/2"]
	assert.Equal(t, e1,
		&Config{
			ExporterSettings: configmodels.ExporterSettings{
				NameVal: "cortex/2",
				TypeVal: "cortex",
			},
			HTTPClientSettings: confighttp.HTTPClientSettings{
				Headers: map[string]string{
					//oof gotta reimplement these params
				},
				Endpoint:    "localhost:8888",
				Compression: "on",
				TLSSetting: configtls.TLSClientSetting{
					TLSSetting: configtls.TLSSetting{
						CAFile: "/var/lib/mycert.pem",
					},
					Insecure: false,
				},
				Keepalive: &configgrpc.KeepaliveClientConfig{
					Time:                20,
					PermitWithoutStream: true,
					Timeout:             30,
				},
				WriteBufferSize: 512 * 1024,
				BalancerName:    "round_robin",
			},
			NumWorkers:        123,
			ReconnectionDelay: 15,
		})
}
// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetry

import (
	"istio.io/api/mesh/v1alpha1"
	telemetryapi "istio.io/api/telemetry/v1alpha1"
	"istio.io/istio/pkg/config/analysis"
	"istio.io/istio/pkg/config/analysis/analyzers/util"
	"istio.io/istio/pkg/config/analysis/msg"
	"istio.io/istio/pkg/config/resource"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
)

type ProdiverAnalyzer struct{}

var _ analysis.Analyzer = &ProdiverAnalyzer{}

// Metadata implements Analyzer
func (a *ProdiverAnalyzer) Metadata() analysis.Metadata {
	return analysis.Metadata{
		Name:        "telemetry.ProviderAnalyzer",
		Description: "Validates that providers in telemery resource is valid",
		Inputs: collection.Names{
			collections.IstioTelemetryV1Alpha1Telemetries.Name(),
			collections.IstioMeshV1Alpha1MeshConfig.Name(),
		},
	}
}

// Analyze implements Analyzer
func (a *ProdiverAnalyzer) Analyze(c analysis.Context) {
	meshConfig := fetchMeshConfig(c)
	if meshConfig.DefaultProviders == nil ||
		len(meshConfig.DefaultProviders.AccessLogging) == 0 {
		c.ForEach(collections.IstioTelemetryV1Alpha1Telemetries.Name(), func(r *resource.Instance) bool {
			telemetry := r.Message.(*telemetryapi.Telemetry)

			for _, l := range telemetry.AccessLogging {
				if l.Disabled != nil && l.Disabled.Value {
					continue
				}
				if len(l.Providers) == 0 {
					c.Report(collections.IstioTelemetryV1Alpha1Telemetries.Name(),
						msg.NewInvalidTelemetryProvider(r, string(r.Metadata.FullName.Name), string(r.Metadata.FullName.Namespace)))
				}
			}

			return true
		})
	}
}

func fetchMeshConfig(c analysis.Context) *v1alpha1.MeshConfig {
	var meshConfig *v1alpha1.MeshConfig
	c.ForEach(collections.IstioMeshV1Alpha1MeshConfig.Name(), func(r *resource.Instance) bool {
		meshConfig = r.Message.(*v1alpha1.MeshConfig)
		return r.Metadata.FullName.Name != util.MeshConfigName
	})

	return meshConfig
}

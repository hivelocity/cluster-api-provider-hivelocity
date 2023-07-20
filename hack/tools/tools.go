//go:build tools
// +build tools

package tools

import (
	_ "sigs.k8s.io/cluster-api/hack/tools"
	_ "sigs.k8s.io/cluster-api/hack/tools/mdbook/embed"
	_ "sigs.k8s.io/cluster-api/hack/tools/mdbook/releaselink"
	_ "sigs.k8s.io/cluster-api/hack/tools/mdbook/tabulate"
)

// Package k8s provides Kubernetes runtime operations for Claw instances using kom library
package k8s

import (
	"fmt"

	"github.com/weibaohui/kom/callbacks"
	"github.com/weibaohui/kom/kom"
)

// Initialize initializes the kom library with kubeconfig
func Initialize(kubeconfig string) error {
	// Register callbacks
	callbacks.RegisterInit()

	// Register cluster
	_, err := kom.Clusters().RegisterByPathWithID(kubeconfig, "default")
	if err != nil {
		return fmt.Errorf("failed to register cluster: %w", err)
	}

	return nil
}

// DefaultCluster returns the default kom cluster instance
func DefaultCluster() *kom.Kubectl {
	return kom.DefaultCluster()
}

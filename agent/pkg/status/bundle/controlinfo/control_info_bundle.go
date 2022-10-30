package controlinfo

import (
	"sync"

	statusbundle "github.com/stolostron/multicluster-global-hub/agent/pkg/status/bundle"
	"github.com/stolostron/multicluster-global-hub/pkg/bundle/status"
)

// NewBundle creates a new instance of Bundle.
func NewBundle(leafHubName string, incarnation uint64) *Bundle {
	return &Bundle{
		LeafHubName:   leafHubName,
		BundleVersion: status.NewBundleVersion(incarnation, 0),
		lock:          sync.Mutex{},
	}
}

// Bundle holds control info passed from LH to HoH.
type Bundle struct {
	LeafHubName   string                `json:"leafHubName"`
	BundleVersion *status.BundleVersion `json:"bundleVersion"`
	lock          sync.Mutex
}

// UpdateObject function to update a single object inside a bundle.
func (bundle *Bundle) UpdateObject(statusbundle.Object) {
	bundle.lock.Lock()
	defer bundle.lock.Unlock()

	bundle.BundleVersion.Generation++
}

// DeleteObject function to delete a single object inside a bundle.
func (bundle *Bundle) DeleteObject(statusbundle.Object) {
	// do nothing
}

// GetBundleVersion function to get bundle version.
func (bundle *Bundle) GetBundleVersion() *status.BundleVersion {
	bundle.lock.Lock()
	defer bundle.lock.Unlock()

	return bundle.BundleVersion
}

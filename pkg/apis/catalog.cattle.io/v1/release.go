package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UIPlugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              UIPluginSpec   `json:"spec"`
	Status            UIPluginStatus `json:"status"`
}

type UIPluginSpec struct {
	Plugin UIPluginEntry `json:"plugin,omitempty"`
}

type UIPluginEntry struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	// Description string `json:"description,omitempty"`
	// Icon        string            `json:"icon,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
	NoCache  bool   `json:"noCache,omitempty"`
	// Annotations map[string]string `json:"annotations,omitempty"`
}

type UIPluginStatus struct {
	// generated timestamp
	// UI can show to user to figure out why something might be broken
	// cacheState enum cached, disabled, pending
	// generated

	// Generated time.Time `json:"generated,omitempty"`
	CacheState string `json:"cacheState,omitempty"`
}

// apiVersion: catalog.cattle.io/v1
// kind: UIPlugin
// metadata:
//   name: epinio
//   namespace: cattle-ui-plugin-system
// spec:
//   plugin: # should initially follow the design of the Helm Chart.yaml fields, could discuss modifying this
//     name: epinio
//     version: 0.0.1
//     description: A UI Plugin for the Epinio chart
//     icon: https://mywebsite.com/icon.svg # not recommended since this won't work for airgapped setups? This icon will be pulled in and stored in memory on loading a UIPlugin for serving the contents
//     annotations:
//       catalog.cattle.io/certified: rancher
//       catalog.cattle.io/display-name: Epinio
//     endpoint: something.namespace.svc
// status:
//   # TBD; depends on developer discussions with UI for what needs to be shown

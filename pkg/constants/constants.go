// Copyright (c) Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package constants

import "time"

const YamlSperator = "\n---\n"

/* #nosec */
const AutoImportSecretName string = "auto-import-secret"

const (
	// AutoImportRetryName is the secret data key of auto import retry
	AutoImportRetryName string = "autoImportRetry"

	// AnnotationAutoImportCurrentRetry is the annotation key of auto import secret used to indicate
	// the current retry times of auto importing a managed cluster
	AnnotationAutoImportCurrentRetry = "managedcluster-import-controller.open-cluster-management.io/current-retry"

	// AnnotationKeepingAutoImportSecret is the annotation key of auto import secret used to indicate
	// keeping this secret after the cluster is imported successfully
	AnnotationKeepingAutoImportSecret = "managedcluster-import-controller.open-cluster-management.io/keeping-auto-import-secret"

	// LabelAutoImportRestore is the label key of auto import secret used for backup restore case
	LabelAutoImportRestore = "cluster.open-cluster-management.io/restore-auto-import-secret"
)

const PodNamespaceEnvVarName = "POD_NAMESPACE"

const ImportFinalizer string = "managedcluster-import-controller.open-cluster-management.io/cleanup"

const SelfManagedLabel string = "local-cluster"

const (
	ClusterImportSecretLabel = "managedcluster-import-controller.open-cluster-management.io/import-secret"
	KlusterletWorksLabel     = "import.open-cluster-management.io/klusterlet-works"
	HostedClusterLabel       = "import.open-cluster-management.io/hosted-cluster"
)

const (
	CreatedViaAnnotation = "open-cluster-management/created-via"
	CreatedViaAI         = "assisted-installer"
	CreatedViaHive       = "hive"
	CreatedViaDiscovery  = "discovery"
)

/* #nosec */
const (
	ImportSecretNameSuffix         = "import"
	ImportSecretImportYamlKey      = "import.yaml"
	ImportSecretCRDSYamlKey        = "crds.yaml"
	ImportSecretCRDSV1YamlKey      = "crdsv1.yaml"
	ImportSecretCRDSV1beta1YamlKey = "crdsv1beta1.yaml"
	ImportSecretTokenExpiration    = "expiration"
)

const (
	// KlusterletDeployModeAnnotation describe the klusterlet deploy mode when importing a managed cluster.
	// If the value is "Hosted", the HostingClusterNameAnnotation annotation will be required,
	// we use HostingClusterNameAnnotation to determine where to deploy the registration-agent and work-agent.
	KlusterletDeployModeAnnotation string = "import.open-cluster-management.io/klusterlet-deploy-mode"

	// HostingClusterNameAnnotation is required in Hosted mode, and the hosting cluster MUST be one
	// of the managed cluster of the hub. The value of the annotation should be the ManagedCluster name of
	// the hosting cluster.
	HostingClusterNameAnnotation string = "import.open-cluster-management.io/hosting-cluster-name"

	// KlusterletNamespaceAnnotation is used to customize the namespace to deploy the agent on the managed
	// cluster. The namespace must have a prefix of "open-cluster-management-", and if it is not set,
	// the namespace of "open-cluster-management-agent" is used to deploy agent.
	// In the Hosted mode, this namespace still exists on the managed cluster to contain
	// necessary resources, like service accounts, roles and rolebindings.
	KlusterletNamespaceAnnotation string = "import.open-cluster-management.io/klusterlet-namespace"
)

const (
	// KlusterletDeployModeDefault is the default deploy mode. the klusterlet will be deployed in the managed-cluster.
	KlusterletDeployModeDefault string = "Default"

	// KlusterletDeployModeHosted means deploying klusterlet outside. the klusterlet will be deployed outside of the managed-cluster.
	KlusterletDeployModeHosted string = "Hosted"
)

const (
	// HostedManifestworkSuffix is a suffix of the hosted mode klusterlet manifestwork name.
	HostedKlusterletManifestworkSuffix = "hosted-klusterlet"

	// HostedManagedKubeconfigManifestworkSuffix is a suffix of the hosted mode managed custer kubeconfig manifestwork name.
	HostedManagedKubeconfigManifestworkSuffix = "hosted-kubeconfig"

	// ManifestWorkFinalizer is used to delete all manifestworks before deleting a managed cluster.
	ManifestWorkFinalizer = "managedcluster-import-controller.open-cluster-management.io/manifestwork-cleanup"

	// PostponeDeletionAnnotation is used to delete the manifest work with this annotation until 10 min after the cluster is deleted.
	PostponeDeletionAnnotation = "open-cluster-management/postpone-delete"

	// ManifestWorkPostponeDeleteTime is the postponed time to delete manifest work with postpone-delete annotation
	ManifestWorkPostponeDeleteTime = 10 * time.Minute
)

const (
	KlusterletSuffix     = "klusterlet"
	KlusterletCRDsSuffix = "klusterlet-crds"
)

const (
	// ConditionManagedClusterImportSucceeded is the condition type of managed cluster to indicate whether the managed
	// cluster is imported successfully
	ConditionManagedClusterImportSucceeded = "ManagedClusterImportSucceeded"

	ConditionReasonManagedClusterWaitForImporting = "ManagedClusterWaitForImporting"
	ConditionReasonManagedClusterImporting        = "ManagedClusterImporting"
	ConditionReasonManagedClusterImportFailed     = "ManagedClusterImportFailed"
	ConditionReasonManagedClusterImported         = "ManagedClusterImported"
)

// Copyright (c) Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package importstatus

import (
	"context"
	"testing"
	"time"

	"github.com/stolostron/managedcluster-import-controller/pkg/constants"
	"github.com/stolostron/managedcluster-import-controller/pkg/helpers"

	"github.com/openshift/library-go/pkg/operator/events/eventstesting"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kubefake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	workfake "open-cluster-management.io/api/client/work/clientset/versioned/fake"
	workinformers "open-cluster-management.io/api/client/work/informers/externalversions"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	workv1 "open-cluster-management.io/api/work/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var testscheme = scheme.Scheme

func init() {
	testscheme.AddKnownTypes(clusterv1.SchemeGroupVersion, &clusterv1.ManagedCluster{})
}

func TestReconcile(t *testing.T) {
	managedClusterName := "test"
	cases := []struct {
		name                    string
		objs                    []client.Object
		works                   []runtime.Object
		expectedErr             bool
		expectedConditionStatus metav1.ConditionStatus
		expectedConditionReason string
	}{
		{
			name:        "no cluster",
			objs:        []client.Object{},
			expectedErr: false,
		},
		{
			name: "hosted cluster",
			objs: []client.Object{
				&clusterv1.ManagedCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name: managedClusterName,
						Annotations: map[string]string{
							constants.KlusterletDeployModeAnnotation: constants.KlusterletDeployModeHosted,
						},
					},
				},
			},
			works:       []runtime.Object{},
			expectedErr: false,
		},
		{
			name: "deletion managed cluster",
			objs: []client.Object{
				&clusterv1.ManagedCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name:              managedClusterName,
						DeletionTimestamp: &metav1.Time{time.Now()},
					},
				},
			},
			works:       []runtime.Object{},
			expectedErr: false,
		},
		{
			name: "managed cluster import condition not exist",
			objs: []client.Object{
				&clusterv1.ManagedCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name: managedClusterName,
					},
				},
			},
			works:       []runtime.Object{},
			expectedErr: false,
		},
		{
			name: "managed cluster import condition not running",
			objs: []client.Object{
				&clusterv1.ManagedCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name: managedClusterName,
					},
					Status: clusterv1.ManagedClusterStatus{
						Conditions: []metav1.Condition{
							helpers.NewManagedClusterImportSucceededCondition(metav1.ConditionFalse,
								constants.ConditionReasonManagedClusterImporting, "test"),
						},
					},
				},
			},
			works:       []runtime.Object{},
			expectedErr: false,
		},
		{
			name: "manifestwork not available",
			objs: []client.Object{
				&clusterv1.ManagedCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name: managedClusterName,
					},
					Status: clusterv1.ManagedClusterStatus{
						Conditions: []metav1.Condition{
							helpers.NewManagedClusterImportSucceededCondition(metav1.ConditionFalse,
								constants.ConditionReasonManagedClusterImporting, "test"),
						},
					},
				},
			},
			works: []runtime.Object{
				&workv1.ManifestWork{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-klusterlet-crds",
						Namespace: managedClusterName,
						Labels: map[string]string{
							constants.KlusterletWorksLabel: "true",
						},
					},
				},
				&workv1.ManifestWork{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-klusterlet",
						Namespace: managedClusterName,
						Labels: map[string]string{
							constants.KlusterletWorksLabel: "true",
						},
					},
				},
			},
			expectedErr:             false,
			expectedConditionStatus: metav1.ConditionFalse,
			expectedConditionReason: constants.ConditionReasonManagedClusterImporting,
		},
		{
			name: "manifestwork available",
			objs: []client.Object{
				&clusterv1.ManagedCluster{
					ObjectMeta: metav1.ObjectMeta{
						Name: managedClusterName,
					},
					Status: clusterv1.ManagedClusterStatus{
						Conditions: []metav1.Condition{
							helpers.NewManagedClusterImportSucceededCondition(metav1.ConditionFalse,
								constants.ConditionReasonManagedClusterImporting, "test"),
						},
					},
				},
			},
			works: []runtime.Object{
				&workv1.ManifestWork{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-klusterlet-crds",
						Namespace: managedClusterName,
						Labels: map[string]string{
							constants.KlusterletWorksLabel: "true",
						},
					},
					Status: workv1.ManifestWorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:   workv1.WorkAvailable,
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
				&workv1.ManifestWork{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-klusterlet",
						Namespace: managedClusterName,
						Labels: map[string]string{
							constants.KlusterletWorksLabel: "true",
						},
					},
					Status: workv1.ManifestWorkStatus{
						Conditions: []metav1.Condition{
							{
								Type:   workv1.WorkAvailable,
								Status: metav1.ConditionTrue,
							},
						},
					},
				},
			},
			expectedErr:             false,
			expectedConditionStatus: metav1.ConditionTrue,
			expectedConditionReason: constants.ConditionReasonManagedClusterImported,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			kubeClient := kubefake.NewSimpleClientset()

			workClient := workfake.NewSimpleClientset(c.works...)
			workInformerFactory := workinformers.NewSharedInformerFactory(workClient, 10*time.Minute)
			workInformer := workInformerFactory.Work().V1().ManifestWorks().Informer()
			for _, work := range c.works {
				workInformer.GetStore().Add(work)
			}

			r := ReconcileImportStatus{
				client:     fake.NewClientBuilder().WithScheme(testscheme).WithObjects(c.objs...).Build(),
				kubeClient: kubeClient,
				workClient: workClient,
				recorder:   eventstesting.NewTestingEventRecorder(t),
			}

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: managedClusterName,
				},
			}
			ctx := context.TODO()
			_, err := r.Reconcile(ctx, req)
			if c.expectedErr && err == nil {
				t.Errorf("expected error, but failed")
			}
			if !c.expectedErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if c.expectedConditionReason != "" {

				managedCluster := &clusterv1.ManagedCluster{}
				err = r.client.Get(ctx,
					types.NamespacedName{
						Name: managedClusterName,
					},
					managedCluster)
				if err != nil {
					t.Errorf("get managed cluster error: %v", err)
				}
				condition := meta.FindStatusCondition(
					managedCluster.Status.Conditions,
					constants.ConditionManagedClusterImportSucceeded,
				)
				if condition.Status != c.expectedConditionStatus {
					t.Errorf("Expect condition status %s, got %s", c.expectedConditionStatus, condition.Status)
				}
				if condition.Reason != c.expectedConditionReason {
					t.Errorf("Expect condition reason %s, got %s, message: %s",
						c.expectedConditionReason, condition.Reason, condition.Message)
				}
			}

		})
	}
}

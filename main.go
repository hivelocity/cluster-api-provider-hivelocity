/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package main provides the executable to run the cluster-api-provider-hivelocity.
package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/controllers"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/utils"
	caphvversion "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/version"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.) to ensure that exec-entrypoint and run can make use of them.
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(clusterv1.AddToScheme(scheme))
	utilruntime.Must(bootstrapv1.AddToScheme(scheme))
	utilruntime.Must(infrav1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

var (
	metricsAddr                  string
	enableLeaderElection         bool
	leaderElectionNamespace      string
	probeAddr                    string
	watchFilterValue             string
	watchNamespace               string
	hivelocityClusterConcurrency int
	hivelocityMachineConcurrency int
	logLevel                     string
	syncPeriod                   time.Duration
)

func main() {
	fs := pflag.CommandLine

	fs.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	fs.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	fs.BoolVar(&enableLeaderElection, "leader-elect", true, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	fs.StringVar(&leaderElectionNamespace, "leader-elect-namespace", "", "Namespace that the controller performs leader election in. If unspecified, the controller will discover which namespace it is running in.")
	fs.IntVar(&hivelocityClusterConcurrency, "hivelocitycluster-concurrency", 1, "Number of HivelocityClusters to process simultaneously")
	fs.IntVar(&hivelocityMachineConcurrency, "hivelocitymachine-concurrency", 1, "Number of HivelocityMachines to process simultaneously")
	fs.StringVar(&watchFilterValue, "watch-filter", "", fmt.Sprintf("Label value that the controller watches to reconcile cluster-api objects. Label key is always %s. If unspecified, the controller watches for all cluster-api objects.", clusterv1.WatchLabel))
	fs.StringVar(&watchNamespace, "namespace", "", "Namespace that the controller watches to reconcile cluster-api objects. If unspecified, the controller watches for cluster-api objects across all namespaces.")
	fs.StringVar(&logLevel, "log-level", "debug", "Specifies log level. Options are 'debug', 'info' and 'error'")
	fs.DurationVar(&syncPeriod, "sync-period", 3*time.Minute, "The minimum interval at which watched resources are reconciled (e.g. 3m)")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Parse()

	ctrl.SetLogger(utils.GetDefaultLogger(logLevel))

	var watchNamespaces map[string]cache.Config
	if watchNamespace != "" {
		watchNamespaces = map[string]cache.Config{
			watchNamespace: {},
		}
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress: probeAddr,
		WebhookServer: webhook.NewServer(webhook.Options{
			Port: 9443,
		}),
		LeaderElection:                enableLeaderElection,
		LeaderElectionNamespace:       leaderElectionNamespace,
		LeaderElectionID:              "hivelocity.cluster.x-k8s.io",
		LeaderElectionResourceLock:    "leases",
		LeaderElectionReleaseOnCancel: true,
		Cache: cache.Options{
			SyncPeriod:        &syncPeriod,
			DefaultNamespaces: watchNamespaces,
		},
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Initialize event recorder.
	record.InitFromRecorder(mgr.GetEventRecorderFor("hv-controller"))

	// Setup the context that's going to be used in controllers and for the manager.
	ctx := ctrl.SetupSignalHandler()

	var wg sync.WaitGroup
	wg.Add(1)

	if err = (&controllers.HivelocityClusterReconciler{
		Client:                         mgr.GetClient(),
		APIReader:                      mgr.GetAPIReader(),
		HVClientFactory:                &hvclient.HivelocityFactory{},
		Scheme:                         mgr.GetScheme(),
		WatchFilterValue:               watchFilterValue,
		TargetClusterManagersWaitGroup: &wg,
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: hivelocityClusterConcurrency}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityCluster")
		os.Exit(1)
	}
	if err = (&controllers.HivelocityMachineReconciler{
		Client:           mgr.GetClient(),
		APIReader:        mgr.GetAPIReader(),
		HVClientFactory:  &hvclient.HivelocityFactory{},
		WatchFilterValue: watchFilterValue,
	}).SetupWithManager(ctx, mgr, controller.Options{MaxConcurrentReconciles: hivelocityMachineConcurrency}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityMachine")
		os.Exit(1)
	}
	if err = (&controllers.HivelocityMachineTemplateReconciler{
		Client:           mgr.GetClient(),
		WatchFilterValue: watchFilterValue,
	}).SetupWithManager(ctx, mgr, controller.Options{}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityMachineTemplate")
		os.Exit(1)
	}
	if err = (&controllers.HivelocityRemediationReconciler{
		Client:           mgr.GetClient(),
		Scheme:           mgr.GetScheme(),
		WatchFilterValue: watchFilterValue,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityRemediation")
		os.Exit(1)
	}
	if err = (&controllers.HivelocityRemediationTemplateReconciler{
		Client:           mgr.GetClient(),
		Scheme:           mgr.GetScheme(),
		WatchFilterValue: watchFilterValue,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityRemediationTemplate")
		os.Exit(1)
	}
	if err = (&infrav1.HivelocityCluster{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityCluster")
		os.Exit(1)
	}
	if err = (&infrav1.HivelocityClusterTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityClusterTemplate")
		os.Exit(1)
	}
	if err = (&infrav1.HivelocityMachine{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityMachine")
		os.Exit(1)
	}
	if err = (&infrav1.HivelocityMachineTemplateWebhook{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityMachineTemplate")
		os.Exit(1)
	}
	if err = (&infrav1.HivelocityRemediation{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityRemediation")
		os.Exit(1)
	}
	if err = (&infrav1.HivelocityRemediationTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityRemediationTemplate")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager", "version", caphvversion.Get().String())
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

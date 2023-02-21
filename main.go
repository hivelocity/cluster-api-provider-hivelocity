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
	"context"
	"flag"
	"fmt"
	"os"

	infrastructurev1alpha1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/controllers"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/hvutils"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(infrastructurev1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func cliTestListServers(ctx context.Context, client hvclient.Client) {
	allServers, err := client.ListServers(ctx)
	if err != nil {
		panic(err.Error())
	}
	server, err := hvutils.FindDeviceByTags("cn=foo", "mn=bar", allServers)
	if err != nil {
		panic(err)
	}
	fmt.Printf("server: %+v\n", server)
}

func cliTestListImages(ctx context.Context, client hvclient.Client) {
	images, err := client.ListImages(ctx, 504)
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		fmt.Println(image)
	}
}

func cliTestListSSHKeys(ctx context.Context, client hvclient.Client) {
	keys, err := client.ListSSHKeys(ctx)
	if err != nil {
		panic(err)
	}
	for _, sshKey := range keys {
		fmt.Println(sshKey)
	}
}

func manualTests() {
	factory := hvclient.HivelocityFactory{}
	client := factory.NewClient(os.Getenv("HIVELOCITY_API_KEY"))
	ctx := context.Background()
	switch arg := os.Args[1]; arg {
	case "ListServers":
		cliTestListServers(ctx, client)
	case "ListSSHKeys":
		cliTestListSSHKeys(ctx, client)
	case "ListImages":
		cliTestListImages(ctx, client)
	default:
		panic(fmt.Sprintf("unknown argument %q", arg))
	}
	os.Exit(0)
}

func main() {
	manualTests()
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	ctx := ctrl.SetupSignalHandler()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "4120761b.cluster.x-k8s.io", // question: adapt this?
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.HivelocityClusterReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr, controller.Options{}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityCluster")
		os.Exit(1)
	}
	if err = (&controllers.HivelocityClusterTemplateReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityClusterTemplate")
		os.Exit(1)
	}
	if err = (&controllers.HivelocityMachineReconciler{
		Client: mgr.GetClient(),
		// todo Scheme: mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr, controller.Options{}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityMachine")
		os.Exit(1)
	}
	if err = (&controllers.HivelocityMachineTemplateReconciler{
		Client: mgr.GetClient(),
	}).SetupWithManager(ctx, mgr, controller.Options{}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityMachineTemplate")
		os.Exit(1)
	}
	if err = (&controllers.HivelocityRemediationReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityRemediation")
		os.Exit(1)
	}
	if err = (&controllers.HivelocityRemediationTemplateReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityRemediationTemplate")
		os.Exit(1)
	}
	if err = (&infrastructurev1alpha1.HivelocityCluster{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityCluster")
		os.Exit(1)
	}
	if err = (&infrastructurev1alpha1.HivelocityClusterTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityClusterTemplate")
		os.Exit(1)
	}
	if err = (&infrastructurev1alpha1.HivelocityMachine{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityMachine")
		os.Exit(1)
	}
	if err = (&infrastructurev1alpha1.HivelocityMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityMachineTemplate")
		os.Exit(1)
	}
	if err = (&infrastructurev1alpha1.HivelocityRemediation{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "HivelocityRemediation")
		os.Exit(1)
	}
	if err = (&infrastructurev1alpha1.HivelocityRemediationTemplate{}).SetupWebhookWithManager(mgr); err != nil {
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

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

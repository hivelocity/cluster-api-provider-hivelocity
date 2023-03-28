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
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.) to ensure that exec-entrypoint and run can make use of them.
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const manualDeviceID = 14730

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(clusterv1.AddToScheme(scheme))
	utilruntime.Must(bootstrapv1.AddToScheme(scheme))
	utilruntime.Must(infrastructurev1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func cliTestListDevices(ctx context.Context, client hvclient.Client) {
	allDevices, err := client.ListDevices(ctx)
	if err != nil {
		panic(err.Error())
	}
	device, err := hvutils.FindDeviceByTags("cn=foo", "mn=bar", allDevices)
	if err != nil {
		panic(err)
	}
	fmt.Printf("device: %+v\n", device)
}

func cliTestSetTags(ctx context.Context, client hvclient.Client) {
	err := client.SetTags(ctx, manualDeviceID, []string{"foo"})
	if err != nil {
		panic(err)
	}
}

func cliTestGetDevice(ctx context.Context, client hvclient.Client) {
	device, err := client.GetDevice(ctx, manualDeviceID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("device: %+v\n", device)
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

func cliTestProvisionDevice(ctx context.Context, client hvclient.Client) {
	opts := hv.BareMetalDeviceUpdate{
		Hostname: "my-host-name.example.com",
		// Tags:     createTags("my-cluster-name", "my-host-name", false),
		OsName:         "Ubuntu 20.x",
		PublicSshKeyId: 861,
	}
	device, err := client.ProvisionDevice(ctx, manualDeviceID, opts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("device: %+v\n", device)
}

func manualTests() {
	factory := hvclient.HivelocityFactory{}
	client := factory.NewClient(os.Getenv("HIVELOCITY_API_KEY"))
	ctx := context.Background()
	switch arg := os.Args[1]; arg {
	case "ListDevices":
		cliTestListDevices(ctx, client)
		os.Exit(0)
	case "ListSSHKeys":
		cliTestListSSHKeys(ctx, client)
		os.Exit(0)
	case "ListImages":
		cliTestListImages(ctx, client)
		os.Exit(0)
	case "SetTags":
		cliTestSetTags(ctx, client)
		os.Exit(0)
	case "GetDevice":
		cliTestGetDevice(ctx, client)
		os.Exit(0)
	case "ProvisionDevice":
		cliTestProvisionDevice(ctx, client)
		os.Exit(0)
	}
}

func main() {
	manualTests()
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", true,
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
		Scheme:                     scheme,
		MetricsBindAddress:         metricsAddr,
		Port:                       9443,
		HealthProbeBindAddress:     probeAddr,
		LeaderElection:             enableLeaderElection,
		LeaderElectionID:           "hivelocity.cluster.x-k8s.io",
		LeaderElectionResourceLock: "leases",
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
		Client:          mgr.GetClient(),
		APIReader:       mgr.GetAPIReader(),
		HVClientFactory: &hvclient.HivelocityFactory{},
		Scheme:          mgr.GetScheme(),
	}).SetupWithManager(ctx, mgr, controller.Options{}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HivelocityCluster")
		os.Exit(1)
	}
	if err = (&controllers.HivelocityMachineReconciler{
		Client:          mgr.GetClient(),
		APIReader:       mgr.GetAPIReader(),
		HVClientFactory: &hvclient.HivelocityFactory{},
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
	if err = (&infrastructurev1alpha1.HivelocityMachineTemplateWebhook{}).SetupWebhookWithManager(mgr); err != nil {
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
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

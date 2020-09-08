/*
Copyright 2020 jonnylangefeld.

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

package main

import (
	"flag"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/pflag"

	"github.com/zorkian/go-datadog-api"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	datadogv1alpha1 "github.com/jonnylangefeld/datadog-operator/apis/datadog/v1alpha1"
	datadogcontroller "github.com/jonnylangefeld/datadog-operator/controllers/datadog"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = datadogv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var secretsPath string
	pflag.StringVarP(&metricsAddr, "metrics-addr", "m", ":8080", "The address the metric endpoint binds to.")
	pflag.BoolVarP(&enableLeaderElection, "enable-leader-election", "l", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	pflag.StringVarP(&secretsPath, "secrets-path", "s", ".secrets.json", "The path to the config file")
	flag.Parse()
	viper.SetConfigFile(secretsPath)
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()
	viper.SetEnvPrefix("datadog")
	viper.AutomaticEnv()
	datadogAPIKey := viper.GetString("api_key")
	datadogApplicationKey := viper.GetString("application_key")

	ctrl.SetLogger(zap.New(zap.UseDevMode(false)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "a32ba9ee.jonnylangefeld.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	datadogClient := datadog.NewClient(datadogAPIKey, datadogApplicationKey)
	if err = (&datadogcontroller.MonitorReconciler{
		Client:        mgr.GetClient(),
		Log:           ctrl.Log.WithName("controllers").WithName("Monitor"),
		Scheme:        mgr.GetScheme(),
		DatadogClient: datadogClient,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Monitor")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

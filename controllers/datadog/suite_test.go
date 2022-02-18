/*
Copyright 2022 jonnylangefeld.

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

package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	datadogv1 "github.com/jonnylangefeld/datadog-operator/apis/datadog/v1"
	//+kubebuilder:scaffold:imports
)

var (
	k8sClient      client.Client
	k8sManager     ctrl.Manager
	testEnv        *envtest.Environment
	mockController *gomock.Controller
	testContext    context.Context
	cancel         context.CancelFunc
	interval       = time.Second
	// set DEBUG=true when debugging tests to not run into timeout
	_, debug = os.LookupEnv("DEBUG")
	timeout  = func() time.Duration {
		if debug {
			return time.Hour * 24
		}
		return time.Second * 10
	}()
)

func TestAPIs(t *testing.T) {
	mockController = gomock.NewController(t)
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(GinkgoWriter)))

	testContext, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	user, err := testEnv.ControlPlane.AddUser(envtest.User{Name: "test", Groups: []string{"system:masters"}}, cfg)
	Expect(err).NotTo(HaveOccurred())
	kubectl, err := user.Kubectl()
	Expect(err).NotTo(HaveOccurred())
	kubectlExec := fmt.Sprintf("\"%s\" %s \"$@\"", kubectl.Path, strings.Join(kubectl.Opts, " "))
	err = ioutil.WriteFile("kubectl", []byte(kubectlExec), os.ModePerm)
	Expect(err).NotTo(HaveOccurred())

	err = datadogv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sManager).ToNot(BeNil())

	// set up the NetworkRegistration controller
	err = (&MonitorReconciler{
		Client:        k8sManager.GetClient(),
		Log:           ctrl.Log.WithName("controllers").WithName("NetworkRegistration"),
		Scheme:        k8sManager.GetScheme(),
		DatadogClient: getDatadogClientMock(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(testContext)
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).ToNot(BeNil())

	close(done)
}, 60)

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

package controllers

import (
	datadogv1 "github.com/jonnylangefeld/datadog-operator/apis/datadog/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

var testMonitor = &datadogv1.Monitor{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test-monitor",
		Namespace: "default",
	},
	Spec: datadogv1.MonitorSpec{
		Name:    "test monitor",
		Type:    "metric alert",
		Message: "test message",
		Query:   "avg(last_5m):avg:datadog.estimated_usage.containers{*} > 1",
		Tags:    []string{"test"},
		Options: &runtime.RawExtension{Raw: []byte(`{"locked":true}`)},
	},
}

var _ = Describe("Monitor Controller", func() {
	Context("Monitor", func() {
		By("Creating Monitor resource")
		It("Should reconcile", func() {
			Expect(k8sClient.Create(testContext, testMonitor)).Should(Succeed())
			Eventually(func() bool {
				got := &datadogv1.Monitor{}
				_ = k8sClient.Get(testContext, types.NamespacedName{Name: testMonitor.Name, Namespace: testMonitor.Namespace}, got)
				return got.Status.ID != 0
			}, timeout, interval).Should(BeTrue())
		})
		By("Deleting Monitor resource")
		It("Should reconcile", func() {
			Expect(k8sClient.Delete(testContext, testMonitor)).Should(Succeed())
			Eventually(func() bool {
				ml := &datadogv1.MonitorList{}
				_ = k8sClient.List(testContext, ml)
				return len(ml.Items) == 0
			}, timeout, interval).Should(BeTrue())
		})
	})
})

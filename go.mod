module github.com/jonnylangefeld/datadog-operator

go 1.13

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/go-logr/logr v0.1.0
	github.com/golang/mock v1.3.1
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/zorkian/go-datadog-api v2.29.0+incompatible
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.0
)

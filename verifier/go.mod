module github.com/open-cluster-management/integrity-verifier/verifier

replace github.com/IBM/integrity-enforcer => ../

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309 // Required by Helm

go 1.13

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/jonboulle/clockwork v0.1.0
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/oliveagle/jsonpath v0.0.0-20180606110733-2e52cf6e6852
	github.com/onsi/gomega v1.10.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/procfs v0.0.11 // indirect
	github.com/r3labs/diff v0.0.0-20191120142937-b4ed99a31f5a
	github.com/sirupsen/logrus v1.4.2
	github.com/tidwall/gjson v1.5.0
	go.uber.org/atomic v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0
	helm.sh/helm/v3 v3.0.2
	k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/cli-runtime v0.18.6
	k8s.io/client-go v0.18.6
	k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6
	k8s.io/kubectl v0.18.2
	k8s.io/utils v0.0.0-20200603063816-c1c6865ac451 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.6
	k8s.io/client-go => k8s.io/client-go v0.18.6
	k8s.io/code-generator => k8s.io/code-generator v0.18.6
	k8s.io/kubectl => k8s.io/kubectl v0.18.6
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.2
)

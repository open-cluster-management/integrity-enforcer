//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	. "github.com/onsi/ginkgo" //nolint:golint

	rspclient "github.com/IBM/integrity-enforcer/verifier/pkg/client/resourcesigningprofile/clientset/versioned/typed/resourcesigningprofile/v1alpha1"
	spclient "github.com/IBM/integrity-enforcer/verifier/pkg/client/signpolicy/clientset/versioned/typed/signpolicy/v1alpha1"
	vcclient "github.com/IBM/integrity-enforcer/verifier/pkg/client/verifierconfig/clientset/versioned/typed/verifierconfig/v1alpha1"
	v1 "k8s.io/api/core/v1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"

	// "k8s.io/client-go/tools/events"
	"k8s.io/klog"
)

var (
	// kubeconfigPath    = os.Getenv("KUBECONFIG")
	local_test, _                       = strconv.ParseBool(os.Getenv("TEST_LOCAL"))
	skip_default_user_test              = true
	kubeconfig_user                     = os.Getenv("KUBE_CONTEXT_USERNAME")
	iv_namespace                        = os.Getenv("IV_OP_NS")
	test_namespace                      = os.Getenv("TEST_NS")
	verifier_dir                        = os.Getenv("VERIFIER_OP_DIR")
	deploy_dir                          = verifier_dir + "test/deploy/"
	kubeconfigManaged                   = os.Getenv("KUBECONFIG")
	tmpDir                              = os.Getenv("TMP_DIR")
	integrityVerifierOperatorCR         = tmpDir + "apis_v1alpha1_integrityverifier.yaml"
	integrityVerifierOperatorCR_updated = tmpDir + "apis_v1alpha1_integrityverifier_update.yaml"
	test_rsp                            = deploy_dir + "test-rsp.yaml"
	test_rsp_iv                         = deploy_dir + "test-rsp-iv-ns.yaml"
	test_configmap                      = deploy_dir + "test-configmap.yaml"
	test_configmap_updated              = deploy_dir + "test-configmap-updated.yaml"
	test_configmap2                     = deploy_dir + "test-configmap-annotation.yaml"
	test_configmap_rs                   = deploy_dir + "test-configmap-rs.yaml"
	DefaultSignPolicyCRName             = "sign-policy"
	iv_op_sa                            = "integrity-verifier-operator-manager"
	iv_op_role                          = "integrity-verifier-operator-leader-election-role"
	iv_op_rb                            = "integrity-verifier-operator-leader-election-rolebinding"
)

type Framework struct {
	BaseName string

	KubeConfig  string
	KubeContext string
	Kubectl     string

	// KubeClientConfig which was used to create the connection.
	KubeClientConfig *rest.Config

	// Kubernetes API clientsets
	KubeClientSet          kubernetes.Interface
	APIExtensionsClientSet apiextcs.Interface
	RSPClient              rspclient.ApisV1alpha1Interface
	SignPolicyClient       spclient.ApisV1alpha1Interface
	VerifierConfigClient   vcclient.ApisV1alpha1Interface

	// Namespace in which all test resources should reside
	Namespace *v1.Namespace
}

func initFrameWork() *Framework {
	framework := &Framework{
		KubeConfig: kubeconfigManaged,
	}
	kubeConfig, err := LoadConfig(framework.KubeConfig, framework.KubeContext)
	if err != nil {
		Fail(fmt.Sprintf("fail to set kubeconfig"))
	}
	framework.KubeClientSet, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		Fail(fmt.Sprintf("fail to set KubeClientSet"))

	}
	framework.APIExtensionsClientSet, err = apiextcs.NewForConfig(kubeConfig)
	if err != nil {
		Fail(fmt.Sprintf("fail to set APIExtensionsClientSet"))
	}
	framework.RSPClient, err = rspclient.NewForConfig(kubeConfig)
	if err != nil {
		Fail(fmt.Sprintf("fail to set RSPClient"))
	}
	framework.SignPolicyClient, err = spclient.NewForConfig(kubeConfig)
	if err != nil {
		Fail(fmt.Sprintf("fail to set SingPolicyClient"))
	}

	framework.VerifierConfigClient, err = vcclient.NewForConfig(kubeConfig)
	if err != nil {
		Fail(fmt.Sprintf("fail to set VerifierConfigClient"))
	}
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: iv_namespace,
		},
	}
	framework.Namespace = ns
	return framework
}

func LoadConfig(config, context string) (*rest.Config, error) {
	c, err := RestclientConfig(config, context)
	if err != nil {
		return nil, err
	}
	return clientcmd.NewDefaultClientConfig(*c, &clientcmd.ConfigOverrides{}).ClientConfig()
}

func RestclientConfig(config, context string) (*clientcmdapi.Config, error) {
	if config == "" {
		return nil, fmt.Errorf("Config file must be specified to load client config")
	}
	c, err := clientcmd.LoadFromFile(config)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err.Error())
	}
	if context != "" {
		c.CurrentContext = context
	}
	return c, nil
}
func Load(data []byte) (*clientcmdapi.Config, error) {
	config := clientcmdapi.NewConfig()
	// if there's no data in a file, return the default object instead of failing (DecodeInto reject empty input)
	if len(data) == 0 {
		return config, nil
	}
	decoded, _, err := clientcmdlatest.Codec.Decode(data, &schema.GroupVersionKind{Version: clientcmdlatest.Version, Kind: "Config"}, config)
	if err != nil {
		return nil, err
	}
	return decoded.(*clientcmdapi.Config), nil
}
func LoadFromFile(filename string) (*clientcmdapi.Config, error) {
	kubeconfigBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := Load(kubeconfigBytes)
	if err != nil {
		return nil, err
	}
	klog.V(6).Infoln("Config loaded from file: ", filename)

	// set LocationOfOrigin on every Cluster, User, and Context
	for key, obj := range config.AuthInfos {
		obj.LocationOfOrigin = filename
		config.AuthInfos[key] = obj
	}
	for key, obj := range config.Clusters {
		obj.LocationOfOrigin = filename
		config.Clusters[key] = obj
	}
	for key, obj := range config.Contexts {
		obj.LocationOfOrigin = filename
		config.Contexts[key] = obj
	}

	if config.AuthInfos == nil {
		config.AuthInfos = map[string]*clientcmdapi.AuthInfo{}
	}
	if config.Clusters == nil {
		config.Clusters = map[string]*clientcmdapi.Cluster{}
	}
	if config.Contexts == nil {
		config.Contexts = map[string]*clientcmdapi.Context{}
	}

	return config, nil
}

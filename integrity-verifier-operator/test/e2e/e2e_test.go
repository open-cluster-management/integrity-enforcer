// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	goctx "context"
	"fmt"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo" //nolint:golint
	. "github.com/onsi/gomega" //nolint:golint

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Test integrity verifier", func() {
	Describe("Check operator status in ns:"+iv_namespace, func() {
		It("Operator Pod should be Running Status", func() {
			framework := initFrameWork()
			var timeout int = 120
			expected := "integrity-verifier-operator-controller-manager"
			Eventually(func() error {
				return CheckPodStatus(framework, iv_namespace, expected)
			}, timeout, 1).Should(BeNil())
		})
		It("Operator sa should be created", func() {
			framework := initFrameWork()
			expected := iv_op_sa
			err := CheckIVResources(framework, "ServiceAccount", iv_namespace, expected)
			Expect(err).To(BeNil())
		})
		It("Operator role should be created", func() {
			framework := initFrameWork()
			expected := iv_op_role
			err := CheckIVResources(framework, "Role", iv_namespace, expected)
			Expect(err).To(BeNil())
		})
		It("Operator rb should be created", func() {
			framework := initFrameWork()
			expected := iv_op_rb
			err := CheckIVResources(framework, "RoleBinding", iv_namespace, expected)
			Expect(err).To(BeNil())
		})
	})

	Describe("Check iv server in ns:"+iv_namespace, func() {
		It("Server should be created properly", func() {
			framework := initFrameWork()
			var timeout int = 1000
			expected := "integrity-verifier-server"
			cmd_err := Kubectl("apply", "-f", integrityVerifierOperatorCR, "-n", iv_namespace)
			Expect(cmd_err).To(BeNil())
			Eventually(func() error {
				return CheckPodStatus(framework, iv_namespace, expected)
			}, timeout, 1).Should(BeNil())
		})
	})

	Describe("Test iv resources", func() {
		It("Iv resources should be created properly", func() {
			framework := initFrameWork()
			vc_name := "iv-config"
			vc, err := framework.VerifierConfigClient.VerifierConfigs(iv_namespace).Get(goctx.Background(), vc_name, metav1.GetOptions{})
			Expect(err).To(BeNil())
			iv_resource_list := vc.Spec.VerifierConfig.IVResourceCondition.References
			By("Check created iv resources...")
			for _, ivr := range iv_resource_list {
				fmt.Print(ivr.Kind, " : ", ivr.Name, "\n")
				if ivr.Name == "" || ivr.Kind == "IntegrityVerifier" || ivr.Kind == "SecurityContextConstraints" || ivr.Kind == "PodSecurityPolicy" || ivr.Name == "helmreleasemetadatas.apis.integrityverifier.io" {
					continue
				}
				err := CheckIVResources(framework, ivr.Kind, ivr.Namespace, ivr.Name)
				Expect(err).To(BeNil())
			}
		})
		It("Deleting iv resources should be blocked", func() {
			if skip_default_user_test {
				Skip("this test should be done by default user")
			}
			time.Sleep(time.Second * 15)
			framework := initFrameWork()
			// change kube context
			err := ChangeKubeContextToDefaultUser(framework, test_namespace, "default-token-")
			Expect(err).To(BeNil())
			// load iv resource lists
			vc_name := "iv-config"
			vc, err := framework.VerifierConfigClient.VerifierConfigs(iv_namespace).Get(goctx.Background(), vc_name, metav1.GetOptions{})
			Expect(err).To(BeNil())
			iv_resource_list := vc.Spec.VerifierConfig.IVResourceCondition.References
			By("Try to Delete iv resources...")
			for _, ivr := range iv_resource_list {
				fmt.Print(ivr.Kind, " : ", ivr.Name, "\n")
				if ivr.Name == "" || ivr.Kind == "SecurityContextConstraints" || ivr.Kind == "PodSecurityPolicy" || ivr.Name == "helmreleasemetadatas.apis.integrityverifier.io" {
					continue
				}
				if ivr.Namespace != "" {
					cmd_err := Kubectl("delete", ivr.Kind, ivr.Name, "-n", ivr.Namespace)
					Expect(cmd_err).NotTo(BeNil())
				} else {
					cmd_err := Kubectl("delete", ivr.Kind, ivr.Name)
					Expect(cmd_err).NotTo(BeNil())
				}
			}
			// change kube context
			err = ChangeKubeContextToKubeAdmin()
			Expect(err).To(BeNil())
		})
		It("IV Resources are changed when IV CR is updated", func() {
			if !local_test {
				Skip("this test is executed only in the local env.")
			}
			framework := initFrameWork()
			var timeout int = 120
			expected := DefaultSignPolicyCRName
			var generation int64
			sp, err := framework.SignPolicyClient.SignPolicies(iv_namespace).Get(goctx.Background(), expected, metav1.GetOptions{})
			Expect(err).To(BeNil())
			generation = sp.Generation
			cmd_err := Kubectl("apply", "-f", integrityVerifierOperatorCR_updated, "-n", iv_namespace)
			Expect(cmd_err).To(BeNil())
			time.Sleep(time.Second * 15)
			Eventually(func() error {
				sp, err := framework.SignPolicyClient.SignPolicies(iv_namespace).Get(goctx.Background(), expected, metav1.GetOptions{})
				if err != nil {
					return err
				}
				if sp.Generation == generation {
					return fmt.Errorf("SignPolicy is not changed: %v", expected)
				}
				return nil
			}, timeout, 1).Should(BeNil())
		})
	})

	Describe("Test integrity verifier server", func() {
		Context("Test integrity verifier function", func() {
			It("Test rsources should be unmonitored if rsp is not exist.", func() {
				framework := initFrameWork()
				var timeout int = 120
				expected := "unmonitored-configmap"
				server_name := "integrity-verifier-server"
				cmd_err := Kubectl("create", "cm", expected, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
				server := GetPodName(framework, iv_namespace, server_name)
				Eventually(func() error {
					cmdstr := "kubectl logs " + server + " -c server -n " + iv_namespace + " | grep " + expected
					out, cmd_err := exec.Command("sh", "-c", cmdstr).Output()
					// fmt.Print("out1: ", string(out), "\n")
					if cmd_err != nil {
						return cmd_err
					}
					cmdstr = "kubectl logs " + server + " -c forwarder -n " + iv_namespace + " | grep " + expected + " | grep unprotected"
					out2, cmd_err := exec.Command("sh", "-c", cmdstr).Output()
					// fmt.Print("out2: ", string(out2), "\n")
					if cmd_err != nil {
						return cmd_err
					}
					if len(string(out)) != 0 && len(string(out2)) != 0 {
						return nil
					}
					return fmt.Errorf("Fail to check unmonitored resource")
				}, timeout, 1).Should(BeNil())
				Eventually(func() error {
					return CheckConfigMap(framework, test_namespace, expected)
				}, timeout, 1).Should(BeNil())
				cmd_err = Kubectl("delete", "cm", expected, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
			})
			It("Test RSP should be created properly", func() {
				framework := initFrameWork()
				var timeout int = 120
				expected := "test-rsp"
				cmd_err := Kubectl("apply", "-f", test_rsp, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
				Eventually(func() error {
					_, err := framework.RSPClient.ResourceSigningProfiles(test_namespace).Get(goctx.Background(), expected, metav1.GetOptions{})
					if err != nil {
						return err
					}
					return nil
				}, timeout, 1).Should(BeNil())
			})
			It("Test unsigned resouce should be blocked", func() {
				framework := initFrameWork()
				time.Sleep(time.Second * 30)
				var timeout int = 120
				expected := "test-configmap"
				cmd_err := Kubectl("apply", "-f", test_configmap, "-n", test_namespace)
				Expect(cmd_err).NotTo(BeNil())
				Eventually(func() error {
					return CheckEventNoSignature(framework, test_namespace, expected)
				}, timeout, 1).Should(BeNil())
			})
			It("Test (ResourceSignature) signed resouce should be allowed", func() {
				framework := initFrameWork()
				var timeout int = 120
				expected := "test-configmap"
				cmd_err := Kubectl("apply", "-f", test_configmap_rs, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
				cmd_err = Kubectl("apply", "-f", test_configmap, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
				Eventually(func() error {
					return CheckConfigMap(framework, test_namespace, expected)
				}, timeout, 1).Should(BeNil())
			})
			It("Test (Annotation) signed resouce should be allowed", func() {
				framework := initFrameWork()
				var timeout int = 120
				expected := "test-configmap2"
				cmd_err := Kubectl("apply", "-f", test_configmap2, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
				Eventually(func() error {
					return CheckConfigMap(framework, test_namespace, expected)
				}, timeout, 1).Should(BeNil())
			})
			It("Test changing whitelisted part should be allowed", func() {
				framework := initFrameWork()
				var timeout int = 120
				expected := "test-configmap"
				server_name := "integrity-verifier-server"
				server := GetPodName(framework, iv_namespace, server_name)
				// apply cm
				cmd_err := Kubectl("apply", "-f", test_configmap_updated, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
				Eventually(func() error {
					// check forwarder log
					cmdstr := "kubectl logs " + server + " -c forwarder -n " + iv_namespace + " | grep " + expected + " | grep no-mutation | grep UPDATE"
					out, cmd_err := exec.Command("sh", "-c", cmdstr).Output()
					if cmd_err != nil {
						return cmd_err
					}
					if len(string(out)) == 0 {
						return fmt.Errorf("Fail to find expected forwarder log")
					}
					return nil
				}, timeout, 1).Should(BeNil())
			})
			It("Test request is allowed if filtered by IgnoredKind", func() {
			})
			It("Test request is allowed if filtered by IgnoredSA", func() {
			})
			It("Test request is allowed if error occured", func() {
			})
			It("Test unsigned resource can be created if filtered by exclude rule", func() {
				framework := initFrameWork()
				var timeout int = 120
				expected := "test-configmap3"
				cmd_err := Kubectl("create", "cm", expected, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
				Eventually(func() error {
					return CheckConfigMap(framework, test_namespace, expected)
				}, timeout, 1).Should(BeNil())
			})
			It("Test RSP should be deleted properly", func() {
				if !local_test {
					Skip("this test is executed only in the local env.")
				}
				framework := initFrameWork()
				var timeout int = 120
				expected := "test-rsp"
				cmd_err := Kubectl("delete", "-f", test_rsp, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
				Eventually(func() error {
					_, err := framework.RSPClient.ResourceSigningProfiles(test_namespace).Get(goctx.Background(), expected, metav1.GetOptions{})
					return err
				}, timeout, 1).ShouldNot(BeNil())
			})
			It("Test unsigned resouce should not blocked", func() {
				if !local_test {
					Skip("this test is executed only in the local env.")
				}
				framework := initFrameWork()
				var timeout int = 120
				expected := "test-configmap4"
				cmd_err := Kubectl("create", "cm", expected, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
				Eventually(func() error {
					return CheckConfigMap(framework, test_namespace, expected)
				}, timeout, 1).Should(BeNil())
			})
		})
		Context("RSP in IV NS is effective for blocking unsigned admission on newly created NS", func() {
			It("Delete test namespace", func() {
				cmd_err := Kubectl("get", "ns", test_namespace)
				if cmd_err == nil {
					test_rsp_name := "test-rsp"
					err := Kubectl("get", "rsp", test_rsp_name, "-n", test_namespace)
					if err == nil {
						err := Kubectl("delete", "-f", test_rsp, "-n", test_namespace)
						Expect(err).To(BeNil())
					}
					err = Kubectl("delete", "ns", test_namespace)
					Expect(err).To(BeNil())
				}
				cmd_err = Kubectl("get", "ns", test_namespace)
				Expect(cmd_err).NotTo(BeNil())
			})
			It("Test RSP should be created properly", func() {
				framework := initFrameWork()
				var timeout int = 120
				expected := "test-rsp"
				By("Creating test rsp: " + test_rsp_iv + " ns: " + iv_namespace)
				cmd_err := Kubectl("apply", "-f", test_rsp_iv, "-n", iv_namespace)
				Expect(cmd_err).To(BeNil())
				By("Checking rsp is created properly: " + test_rsp_iv + " ns: " + iv_namespace)
				Eventually(func() error {
					_, err := framework.RSPClient.ResourceSigningProfiles(iv_namespace).Get(goctx.Background(), expected, metav1.GetOptions{})
					if err != nil {
						return err
					}
					return nil
				}, timeout, 1).Should(BeNil())
			})
			It("Test unsigned resource should be blocked in new namespace", func() {
				framework := initFrameWork()
				time.Sleep(time.Second * 30)
				var timeout int = 120
				expected := "test-configmap"
				By("Creating new namespace: " + test_namespace)
				cmd_err := Kubectl("create", "ns", test_namespace)
				Expect(cmd_err).To(BeNil())
				By("Creating test configmap in ns: " + test_namespace)
				cmd_err = Kubectl("apply", "-f", test_configmap, "-n", test_namespace)
				Expect(cmd_err).NotTo(BeNil())
				Eventually(func() error {
					return CheckEventNoSignature(framework, test_namespace, expected)
				}, timeout, 1).Should(BeNil())
			})
			It("Test signed resource should be allowed in new namespace", func() {
				framework := initFrameWork()
				var timeout int = 120
				expected := "test-configmap2"
				By("Creating test configmap in ns: " + test_namespace)
				cmd_err := Kubectl("apply", "-f", test_configmap2, "-n", test_namespace)
				Expect(cmd_err).To(BeNil())
				Eventually(func() error {
					return CheckConfigMap(framework, test_namespace, expected)
				}, timeout, 1).Should(BeNil())
			})
		})
	})

	Describe("Test integrity verifier resources: delete", func() {
		It("Server and iv resources should be deleted properly", func() {
			if !local_test {
				Skip("this test is executed only in the local env.")
			}
			framework := initFrameWork()
			vc_name := "iv-config"
			By("Load iv resource list")
			vc, err := framework.VerifierConfigClient.VerifierConfigs(iv_namespace).Get(goctx.Background(), vc_name, metav1.GetOptions{})
			Expect(err).To(BeNil())
			iv_resource_list := vc.Spec.VerifierConfig.IVResourceCondition.References
			By("Server should be deleted properly")
			var timeout int = 300
			expected := "integrity-verifier-server"
			cmd_err := Kubectl("delete", "-f", integrityVerifierOperatorCR, "-n", iv_namespace)
			Expect(cmd_err).To(BeNil())
			Eventually(func() error {
				return CheckPodStatus(framework, iv_namespace, expected)
			}, timeout, 1).ShouldNot(BeNil())
			By("Iv resources should be deleted properly")
			time.Sleep(time.Second * 30)
			for _, ivr := range iv_resource_list {
				fmt.Print(ivr.Kind, " : ", ivr.Name, "\n")
				if ivr.Name == "" || ivr.Kind == "IntegrityVerifier" || ivr.Kind == "SecurityContextConstraints" || ivr.Kind == "PodSecurityPolicy" || ivr.Name == "integrity-verifier-operator-controller-manager" || ivr.Name == "helmreleasemetadatas.apis.integrityverifier.io" || ivr.Name == "integrityverifiers.apis.integrityverifier.io" {
					continue
				}
				err := CheckIVResources(framework, ivr.Kind, ivr.Namespace, ivr.Name)
				Expect(err).NotTo(BeNil())
			}
		})
	})
})

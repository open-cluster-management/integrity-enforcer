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
			var timeout int = 300
			expected := "integrity-verifier-server"
			cmd_err := Kubectl("apply", "-f", integrityVerifierOperatorCR, "-n", iv_namespace)
			Expect(cmd_err).To(BeNil())
			Eventually(func() error {
				return CheckPodStatus(framework, iv_namespace, expected)
			}, timeout, 1).Should(BeNil())
		})
	})

})

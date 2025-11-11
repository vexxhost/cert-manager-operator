// Copyright 2025 VEXXHOST, Inc.
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	infrav1alpha1 "github.com/vexxhost/cert-manager-operator/api/v1alpha1"
)

var _ = Describe("CertManager", func() {
	ctx := context.Background()

	AfterEach(func() {
		list := &infrav1alpha1.CertManagerList{}
		err := k8sClient.List(ctx, list)
		Expect(err).NotTo(HaveOccurred())

		for _, cm := range list.Items {
			By("cleanup CertManager instance: " + cm.Name)
			Expect(k8sClient.Delete(ctx, &cm)).To(Succeed())
		}
	})

	Context("Validation", func() {
		It("should accept a resource named 'default'", func() {
			By("creating a CertManager resource named 'default'")
			certManager := &infrav1alpha1.CertManager{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: infrav1alpha1.CertManagerSpec{},
			}

			err := k8sClient.Create(ctx, certManager)
			Expect(err).NotTo(HaveOccurred())

			By("verifying the resource was created")
			created := &infrav1alpha1.CertManager{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: "default"}, created)
			Expect(err).NotTo(HaveOccurred())
			Expect(created.Name).To(Equal("default"))
		})

		It("should reject a resource with a name other than 'default'", func() {
			By("attempting to create a CertManager resource named 'invalid'")
			certManager := &infrav1alpha1.CertManager{
				ObjectMeta: metav1.ObjectMeta{
					Name: "invalid",
				},
				Spec: infrav1alpha1.CertManagerSpec{},
			}

			err := k8sClient.Create(ctx, certManager)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("CertManager resource must be named 'default'"))
		})

		It("should prevent duplicate resources named 'default'", func() {
			By("creating the first CertManager resource named 'default'")
			certManager1 := &infrav1alpha1.CertManager{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: infrav1alpha1.CertManagerSpec{},
			}
			err := k8sClient.Create(ctx, certManager1)
			Expect(err).NotTo(HaveOccurred())

			By("attempting to create a second CertManager resource named 'default'")
			certManager2 := &infrav1alpha1.CertManager{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
				Spec: infrav1alpha1.CertManagerSpec{},
			}
			err = k8sClient.Create(ctx, certManager2)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsAlreadyExists(err)).To(BeTrue())
		})
	})
})

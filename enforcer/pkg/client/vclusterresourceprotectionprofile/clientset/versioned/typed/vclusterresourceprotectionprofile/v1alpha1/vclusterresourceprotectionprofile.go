/*
Copyright The 2020 IBM Corporation.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/IBM/integrity-enforcer/enforcer/pkg/apis/vclusterresourceprotectionprofile/v1alpha1"
	scheme "github.com/IBM/integrity-enforcer/enforcer/pkg/client/vclusterresourceprotectionprofile/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VClusterResourceProtectionProfilesGetter has a method to return a VClusterResourceProtectionProfileInterface.
// A group's client should implement this interface.
type VClusterResourceProtectionProfilesGetter interface {
	VClusterResourceProtectionProfiles() VClusterResourceProtectionProfileInterface
}

// VClusterResourceProtectionProfileInterface has methods to work with VClusterResourceProtectionProfile resources.
type VClusterResourceProtectionProfileInterface interface {
	Create(*v1alpha1.VClusterResourceProtectionProfile) (*v1alpha1.VClusterResourceProtectionProfile, error)
	Update(*v1alpha1.VClusterResourceProtectionProfile) (*v1alpha1.VClusterResourceProtectionProfile, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.VClusterResourceProtectionProfile, error)
	List(opts v1.ListOptions) (*v1alpha1.VClusterResourceProtectionProfileList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.VClusterResourceProtectionProfile, err error)
	VClusterResourceProtectionProfileExpansion
}

// vClusterResourceProtectionProfiles implements VClusterResourceProtectionProfileInterface
type vClusterResourceProtectionProfiles struct {
	client rest.Interface
}

// newVClusterResourceProtectionProfiles returns a VClusterResourceProtectionProfiles
func newVClusterResourceProtectionProfiles(c *ResearchV1alpha1Client) *vClusterResourceProtectionProfiles {
	return &vClusterResourceProtectionProfiles{
		client: c.RESTClient(),
	}
}

// Get takes name of the vClusterResourceProtectionProfile, and returns the corresponding vClusterResourceProtectionProfile object, and an error if there is any.
func (c *vClusterResourceProtectionProfiles) Get(name string, options v1.GetOptions) (result *v1alpha1.VClusterResourceProtectionProfile, err error) {
	result = &v1alpha1.VClusterResourceProtectionProfile{}
	err = c.client.Get().
		Resource("vclusterresourceprotectionprofiles").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VClusterResourceProtectionProfiles that match those selectors.
func (c *vClusterResourceProtectionProfiles) List(opts v1.ListOptions) (result *v1alpha1.VClusterResourceProtectionProfileList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.VClusterResourceProtectionProfileList{}
	err = c.client.Get().
		Resource("vclusterresourceprotectionprofiles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested vClusterResourceProtectionProfiles.
func (c *vClusterResourceProtectionProfiles) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("vclusterresourceprotectionprofiles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a vClusterResourceProtectionProfile and creates it.  Returns the server's representation of the vClusterResourceProtectionProfile, and an error, if there is any.
func (c *vClusterResourceProtectionProfiles) Create(vClusterResourceProtectionProfile *v1alpha1.VClusterResourceProtectionProfile) (result *v1alpha1.VClusterResourceProtectionProfile, err error) {
	result = &v1alpha1.VClusterResourceProtectionProfile{}
	err = c.client.Post().
		Resource("vclusterresourceprotectionprofiles").
		Body(vClusterResourceProtectionProfile).
		Do().
		Into(result)
	return
}

// Update takes the representation of a vClusterResourceProtectionProfile and updates it. Returns the server's representation of the vClusterResourceProtectionProfile, and an error, if there is any.
func (c *vClusterResourceProtectionProfiles) Update(vClusterResourceProtectionProfile *v1alpha1.VClusterResourceProtectionProfile) (result *v1alpha1.VClusterResourceProtectionProfile, err error) {
	result = &v1alpha1.VClusterResourceProtectionProfile{}
	err = c.client.Put().
		Resource("vclusterresourceprotectionprofiles").
		Name(vClusterResourceProtectionProfile.Name).
		Body(vClusterResourceProtectionProfile).
		Do().
		Into(result)
	return
}

// Delete takes name of the vClusterResourceProtectionProfile and deletes it. Returns an error if one occurs.
func (c *vClusterResourceProtectionProfiles) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("vclusterresourceprotectionprofiles").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *vClusterResourceProtectionProfiles) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("vclusterresourceprotectionprofiles").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched vClusterResourceProtectionProfile.
func (c *vClusterResourceProtectionProfiles) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.VClusterResourceProtectionProfile, err error) {
	result = &v1alpha1.VClusterResourceProtectionProfile{}
	err = c.client.Patch(pt).
		Resource("vclusterresourceprotectionprofiles").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}

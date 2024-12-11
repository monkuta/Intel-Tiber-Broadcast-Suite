//go:build !ignore_autogenerated

/*
Copyright 2024.

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

/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppParams) DeepCopyInto(out *AppParams) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppParams.
func (in *AppParams) DeepCopy() *AppParams {
	if in == nil {
		return nil
	}
	out := new(AppParams)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BcsConfig) DeepCopyInto(out *BcsConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BcsConfig.
func (in *BcsConfig) DeepCopy() *BcsConfig {
	if in == nil {
		return nil
	}
	out := new(BcsConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BcsConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BcsConfigList) DeepCopyInto(out *BcsConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BcsConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BcsConfigList.
func (in *BcsConfigList) DeepCopy() *BcsConfigList {
	if in == nil {
		return nil
	}
	out := new(BcsConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BcsConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BcsConfigSpec) DeepCopyInto(out *BcsConfigSpec) {
	*out = *in
	out.AppParams = in.AppParams
	out.Connections = in.Connections
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BcsConfigSpec.
func (in *BcsConfigSpec) DeepCopy() *BcsConfigSpec {
	if in == nil {
		return nil
	}
	out := new(BcsConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BcsConfigStatus) DeepCopyInto(out *BcsConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BcsConfigStatus.
func (in *BcsConfigStatus) DeepCopy() *BcsConfigStatus {
	if in == nil {
		return nil
	}
	out := new(BcsConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Connections) DeepCopyInto(out *Connections) {
	*out = *in
	out.DataConnection = in.DataConnection
	out.ControlConnection = in.ControlConnection
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Connections.
func (in *Connections) DeepCopy() *Connections {
	if in == nil {
		return nil
	}
	out := new(Connections)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControlConnection) DeepCopyInto(out *ControlConnection) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControlConnection.
func (in *ControlConnection) DeepCopy() *ControlConnection {
	if in == nil {
		return nil
	}
	out := new(ControlConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataConnection) DeepCopyInto(out *DataConnection) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataConnection.
func (in *DataConnection) DeepCopy() *DataConnection {
	if in == nil {
		return nil
	}
	out := new(DataConnection)
	in.DeepCopyInto(out)
	return out
}

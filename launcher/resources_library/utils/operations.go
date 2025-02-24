/*
 * SPDX-FileCopyrightText: Copyright (c) 2024 Intel Corporation
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */

package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"bcs.pod.launcher.intel/resources_library/resources/general"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func isImagePulled(ctx context.Context, cli *client.Client, imageName string) (error, bool) {
	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return err, false
	}

	imageMap := make(map[string]bool)
	for _, image := range images {
		for _, tag := range image.RepoTags {
			imageMap[tag] = true
		}
	}

	_, isPulled := imageMap[imageName]
	return nil, isPulled
}

func pullImageIfNotExists(ctx context.Context, cli *client.Client, imageName string, log logr.Logger) error {

	// Check if the Docker client is nil
	if cli == nil {
		err := errors.New("docker client is nil")
		log.Error(err, "Docker client is not initialized")
		return err
	}

	// Check if the context is nil
	if ctx == nil {
		err := errors.New("context is nil")
		log.Error(err, "Context is not initialized")
		return err
	}

	// Check if the image is already pulled
	err, pulled := isImagePulled(ctx, cli, imageName)
	if err != nil {
		log.Error(err, "Error checking if image is pulled")
		return err
	}

	// Pull the image if it is not already pulled
	if !pulled {
		reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
		if err != nil {
			log.Error(err, "Error pulling image")
			return err
		}
		defer reader.Close()

		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			log.Error(err, "Error reading output")
			return err
		}
		log.Info("Image pulled successfully", "image", imageName)
	}

	return nil
}

func doesContainerExist(ctx context.Context, cli *client.Client, containerName string) (error, bool) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err, false
	}

	containerMap := make(map[string]string)
	for _, container := range containers {
		for _, name := range container.Names {
			containerMap[name] = strings.ToLower(container.State)
		}
	}

	state, exists := containerMap["/"+containerName]
	if !exists {
		return nil, false
	}

	return nil, state == "exited"
}

func isContainerRunning(ctx context.Context, cli *client.Client, containerName string) (error, bool) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err, false
	}

	containerMap := make(map[string]string)
	for _, container := range containers {
		for _, name := range container.Names {
			containerMap[name] = strings.ToLower(container.State)
		}
	}

	state, exists := containerMap["/"+containerName]
	if !exists {
		return nil, false
	}

	return nil, state == "running"
}

func removeContainer(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}


func constructContainerConfig(containerInfo general.Containers) (*container.Config, *container.HostConfig, *network.NetworkingConfig) {
	var containerConfig *container.Config
	var hostConfig *container.HostConfig
	var networkConfig *network.NetworkingConfig

	switch containerInfo.Type {
	case general.MediaProxyAgent:
		fmt.Printf(">> MediaProxyAgentConfig: %+v\n", containerInfo.Configuration.MediaProxyAgentConfig)
		containerConfig = &container.Config{
			User: "root",
			Image: containerInfo.Configuration.MediaProxyAgentConfig.ImageAndTag,
			Cmd:   []string{"-c", containerInfo.Configuration.MediaProxyAgentConfig.RestPort, "-p", containerInfo.Configuration.MediaProxyAgentConfig.GRPCPort},
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			PortBindings: nat.PortMap{
				nat.Port(fmt.Sprintf("%s/tcp", containerInfo.Configuration.MediaProxyAgentConfig.RestPort)): []nat.PortBinding{{HostPort: containerInfo.Configuration.MediaProxyAgentConfig.RestPort}},
				nat.Port(fmt.Sprintf("%s/tcp", containerInfo.Configuration.MediaProxyAgentConfig.GRPCPort)): []nat.PortBinding{{HostPort: containerInfo.Configuration.MediaProxyAgentConfig.GRPCPort}},
			},
		}
	    if containerInfo.Configuration.MediaProxyAgentConfig.Network.Enable {
			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					containerInfo.Configuration.MediaProxyAgentConfig.Network.Name: {
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: containerInfo.Configuration.MediaProxyAgentConfig.Network.IP,
						},
					},
				},
			}
		}else{
			networkConfig = &network.NetworkingConfig{}
		}
	case general.MediaProxyMCM:
		fmt.Printf(">> MediaProxyMcmConfig: %+v\n", containerInfo.Configuration.MediaProxyMcmConfig)
		containerConfig = &container.Config{
			Image: containerInfo.Configuration.MediaProxyMcmConfig.ImageAndTag,
			Cmd:   []string{"-d", fmt.Sprintf("kernel:%s", containerInfo.Configuration.MediaProxyMcmConfig.InterfaceName),"-i", "localhost"},
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			Binds:      containerInfo.Configuration.MediaProxyMcmConfig.Volumes,
		}
	
	    if containerInfo.Configuration.MediaProxyMcmConfig.Network.Enable {
			networkConfig = &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					containerInfo.Configuration.MediaProxyMcmConfig.Network.Name: {
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: containerInfo.Configuration.MediaProxyMcmConfig.Network.IP,
						},
					},
				},
			}
		}else{
			networkConfig = &network.NetworkingConfig{}
		}
    case general.BcsPipelineFfmpeg:
		fmt.Printf(">> BcsPipelineFfmpeg: %+v\n", containerInfo.Configuration.WorkloadConfig.FfmpegPipeline)

		containerConfig = &container.Config{
			User:       "root",
			Image: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.ImageAndTag,
			Cmd:   []string{containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Network.IP, fmt.Sprintf("%d", containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.GRPCPort)},
			Env: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.EnvironmentVariables,
			ExposedPorts: nat.PortSet{
				"20000/tcp": struct{}{},
				"20170/tcp": struct{}{},
			},
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			CapAdd:     []string{"ALL"},
			
			Mounts: []mount.Mount{
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Videos, Target: "/videos"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Dri, Target: "/usr/local/lib/x86_64-linux-gnu/dri/"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Kahawai, Target: "/tmp/kahawai_lcore.lock"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Devnull, Target: "/dev/null"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.TmpHugepages, Target: "/tmp/hugepages"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Hugepages, Target: "/hugepages"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Imtl, Target: "/var/run/imtl"},
				{Type: mount.TypeBind, Source: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Shm, Target: "/dev/shm"},
			},
			IpcMode: "host",
		}
		hostConfig.Devices= []container.DeviceMapping{
			{PathOnHost: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Devices.Vfio, PathInContainer: "/dev/vfio"},
			{PathOnHost: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Volumes.Dri, PathInContainer: "/dev/dri"},
		}
	
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Network.Name: {
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address: containerInfo.Configuration.WorkloadConfig.FfmpegPipeline.Network.IP,
					},
				},
			},
		}
	case general.BcsPipelineNmosClient:
		fmt.Printf(">> NmosClient: %+v\n", containerInfo.Configuration.WorkloadConfig.NmosClient)
		containerConfig = &container.Config{
			Image: containerInfo.Configuration.WorkloadConfig.NmosClient.ImageAndTag,
			Cmd: []string{"config/node.json"},
			Env: containerInfo.Configuration.WorkloadConfig.NmosClient.EnvironmentVariables,
			User:       "root",
		}
	
		hostConfig = &container.HostConfig{
			Privileged: true,
			Binds:      []string{fmt.Sprintf("%s:/home/config/", containerInfo.Configuration.WorkloadConfig.NmosClient.NmosConfigPath)},
		}
	
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				containerInfo.Configuration.WorkloadConfig.NmosClient.Network.Name: {
					IPAMConfig: &network.EndpointIPAMConfig{
						IPv4Address: containerInfo.Configuration.WorkloadConfig.NmosClient.Network.IP,
					},
					Aliases: []string{containerInfo.Configuration.WorkloadConfig.NmosClient.Network.Name},
				},
			},
		}
	default:
		containerConfig, hostConfig, networkConfig = nil, nil, nil
	}

	return containerConfig, hostConfig, networkConfig
}

func CreateAndRunContainer(ctx context.Context, cli *client.Client, log logr.Logger, containerInfo general.Containers) error {
	err, isRunning := isContainerRunning(ctx, cli, containerInfo.ContainerName)
	if err != nil {
		log.Error(err, "Failed to read container status (if it is in running state)")
		return err
	}

	if isRunning {
		log.Info("Container ", containerInfo.ContainerName, " is running. Omitting this container creation.")
		return nil
	}

	err, exists := doesContainerExist(ctx, cli, containerInfo.ContainerName)
	if err != nil {
		log.Error(err, "Failed to read container status (if it exists)")
		return err
	}

	if exists {
		log.Info("Removing container to re-create and re-run because container with a such name exists but with status exited:", "container", containerInfo.ContainerName)
		err = removeContainer(ctx, cli, containerInfo.ContainerName)
		if err != nil {
			log.Error(err, "Failed to remove container")
			return err
		}

	}

	err = pullImageIfNotExists(ctx, cli, containerInfo.Image, log)
	if err != nil {
		log.Error(err, "Error pulling image for container")
		return err
	}
	// Define the container configuration

	containerConfig, hostConfig, networkConfig := constructContainerConfig(containerInfo)
	// Create the container
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, containerInfo.ContainerName)

	if err != nil {
		log.Error(err, "Error creating container")
		return err
	}

	// Start the container
	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Error(err, "Error starting container")
		return err
	}

	log.Info("Container is created and started successfully", "name", containerInfo.ContainerName, "container id: ", resp.ID)
	return nil
}

func boolPtr(b bool) *bool    { return &b }
func intstrPtr(i int) intstr.IntOrString {
    return intstr.IntOrString{IntVal: int32(i)}
}
		
func CreateDeployment(name string) *appsv1.Deployment {
	 // Define Deployment
	 return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mesh-agent-deployment",
			Namespace: "default",
			Labels: map[string]string{
				"app": "mesh-agent",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mesh-agent",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "mesh-agent",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mesh-agent",
							Image: "mcm/mesh-agent:latest",
							Command: []string{
								"-c", "8100", "-p", "50051",
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8100},
								{ContainerPort: 50051},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolPtr(true),
							},
						},
					},
				},
			},
		},
	}
}

func CreateMeshAgentService(name string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mesh-agent-service",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "mesh-agent",
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       8100,
					TargetPort: intstrPtr(8100),
				},
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       50051,
					TargetPort: intstrPtr(50051),
				},
			},
		},
	}
}

func CreateService(name string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": name},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     80,
				},
			},
		},
	}
}

func CreatePersistentVolume(name string) *corev1.PersistentVolume {
	return &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse("1Gi"),
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/mnt/data",
				},
			},
		},
	}
}

func CreatePersistentVolumeClaim(name string) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}
}

func CreateConfigMap(name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Data: map[string]string{
			"data": "bcsdata",
		},
	}
}

func CreateDaemonSet(name string) *appsv1.DaemonSet {
	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "example-container",
							Image: "nginx:latest",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "example-volume",
									MountPath: "/usr/share/nginx/html",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "example-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: name,
								},
							},
						},
					},
				},
			},
		},
	}
}

func int32Ptr(i int32) *int32 { return &i }

func CreateNamespace(namespaceName string) *corev1.Namespace {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}
	return namespace
}

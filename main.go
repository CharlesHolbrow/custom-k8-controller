/*
Copyright 2016 The Kubernetes Authors.

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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/deckarep/golang-set"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	fmt.Println("Running ")
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for {
		// Get all the pods that are running a simulation on a map region.
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
			LabelSelector: "app=aether,type=map-region",
		})

		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		currentLocations := mapset.NewSet()

		fmt.Println("Found locations...")
		for i, pod := range pods.Items {
			if pod.Status.Phase == apiv1.PodRunning {

			}
			name := pod.GetName()
			location := pod.GetLabels()["location"]
			currentLocations.Add(location)

			fmt.Printf("%d\tName: %s \t Location: %s \n", i, name, location)
		}

		// Calculate which regions need to be added, and which can be stopped
		wantedLocations := getWantedLocations()
		missingLocations := wantedLocations.Difference(currentLocations)
		extraLocations := currentLocations.Difference(wantedLocations)

		fmt.Println("Missing and Extra Locations...")
		for missingLoc := range missingLocations.Iter() {
			fmt.Println("Missing location:", missingLoc)
		}

		for extraLoc := range extraLocations.Iter() {
			fmt.Println("Extra location:", extraLoc)
		}
		fmt.Println("...End Missing and Extra Locations!!!!")

		// Create missing pods
		for location := range missingLocations.Iter() {

			desiredPod := createAetherPod(location.(string))
			resultPod, err := clientset.CoreV1().Pods("default").Create(desiredPod)

			if errors.IsAlreadyExists(err) {
				fmt.Printf("Cannot create pod. It already exists: %s\n", desiredPod.Name)
			} else if err != nil {
				fmt.Printf("Error Creating Pod %s: %s\n", location, err.Error())
				continue
			}
			fmt.Println("Created pod:", resultPod.Name)
		}

		fmt.Println("Aether Reconciliation Loop Sleeping...")
		time.Sleep(10 * time.Second)
		fmt.Println("...Aether Reconciliation Loop Wake\n")
	}
}

func getWantedLocations() mapset.Set {
	return mapset.NewSetWith("map1_0_0_3_3", "map1_-4_-4_-1_-1")
}

func createAetherDeployment() *appsv1beta1.Deployment {
	deployment := &appsv1beta1.Deployment{}
	if true {
		panic("CreateAetherDeployment is not implemented")
	}
	return deployment
}

func createAetherPod(location string) *apiv1.Pod {
	hasher := fnv.New64()
	hasher.Write([]byte(location))
	hashed := fmt.Sprintf("%x", hasher.Sum64())

	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "aether-" + hashed,
			Labels: map[string]string{
				"app":      "aether",
				"type":     "map-region",
				"location": location,
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "web",
					Image: "nginx:1.12",
					Ports: []apiv1.ContainerPort{
						{
							Name:          "http",
							Protocol:      apiv1.ProtocolTCP,
							ContainerPort: 80,
						},
					},
					Resources: apiv1.ResourceRequirements{
						Limits: apiv1.ResourceList{},
					},
				},
			},
			RestartPolicy: apiv1.RestartPolicyAlways,
		},
	}
	return pod
}

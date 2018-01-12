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
	"time"

	"github.com/deckarep/golang-set"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
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
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
			LabelSelector: "app=aether,type=simulation",
		})

		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		tileHashes := mapset.NewSet()

		for i, pod := range pods.Items {
			name := pod.GetName()
			kind := pod.Kind
			labelMap := pod.GetLabels()
			tileHash := labelMap["location"]
			tileHashes.Add(tileHash)

			fmt.Printf("%d\tName: %s \t Kind: %s\n", i, name, kind)
		}

		wantedHashes := GetWantedHashes()
		missingHashes := wantedHashes.Difference(tileHashes)
		extraHashes := tileHashes.Difference(wantedHashes)

		fmt.Println("Missing and Extra Hashes...")
		for missingHash := range missingHashes.Iter() {
			fmt.Println("Missing Hashes:", missingHash)
		}

		for extraHash := range extraHashes.Iter() {
			fmt.Println("ExtraHash:", extraHash)
		}
		fmt.Println("...End Missing and Extra Hashes!")

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("default").Get("example-xxxxx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod not found\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod\n")
		}

		time.Sleep(10 * time.Second)
	}
}

func GetWantedHashes() mapset.Set {
	return mapset.NewSetWith("map1:(0,0)|(3,3)", "map1:(-4,-4)|(-1,-1)")
}
package utils

import (
	"fmt"
	"github.com/golang/glog"
	kapi "k8s.io/api/core/v1"
	"os"
)

// PodWantsNetwork check if pod needs cni
func PodWantsNetwork(pod *kapi.Pod) bool {
	return !pod.Spec.HostNetwork
}

// PodScheduled check if pod is assigned to node
func PodScheduled(pod *kapi.Pod) bool {
	return pod.Spec.NodeName != ""
}

// LoadEnvVar loads environment variable
func LoadEnvVar(envName string) (string, error) {
	glog.V(3).Infof("LoadEnvVar(): envName %s", envName)
	value, ok := os.LookupEnv(envName)
	if !ok {
		err := fmt.Errorf("environment variable %s don't exist", envName)
		glog.Error(err)
		return "", err
	}

	return value, nil
}

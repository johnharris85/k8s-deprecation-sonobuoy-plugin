package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
)

type appliedConfig struct {
	APIVersion string `json:"apiVersion"`
}

type apiChanges struct {
	Deprecated []string
	New        string
}

type result struct {
	Namespace     string `json:"namespace"`
	Name          string `json:"name"`
	DeprecatedAPI string `json:"deprecatedAPI"`
	NewAPI        string `json:"newAPI"`
}

const LastAppliedConfiguration = "kubectl.kubernetes.io/last-applied-configuration"

func main() {
	deprecatedAPIs := map[string]apiChanges{
		"networkPolicy": apiChanges{
			[]string{"extensions/v1beta1"},
			"networking.k8s.io/v1",
		},
		"podSecurityPolicy": apiChanges{
			[]string{"extensions/v1beta1"},
			"policy/v1beta1",
		},
		"daemonSet": apiChanges{
			[]string{"extensions/v1beta1", "apps/v1beta1", "apps/v1beta2"},
			"apps/v1",
		},
		"deployment": apiChanges{
			[]string{"extensions/v1beta1", "apps/v1beta1", "apps/v1beta2"},
			"apps/v1",
		},
		"statefulSet": apiChanges{
			[]string{"extensions/v1beta1", "apps/v1beta1", "apps/v1beta2"},
			"apps/v1",
		},
		"replicaSet": apiChanges{
			[]string{"extensions/v1beta1", "apps/v1beta1", "apps/v1beta2"},
			"apps/v1",
		},
	}

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	listOptions := metav1.ListOptions{}

	var results []result

	nwpl, err := clientset.NetworkingV1().NetworkPolicies("").List(listOptions)
	for _, a := range nwpl.Items {
		if a.Annotations[LastAppliedConfiguration] != "" {
			var ac appliedConfig
			d := []byte(a.Annotations[LastAppliedConfiguration])
			err := json.Unmarshal(d, &ac)
			if err != nil {
				log.Println(err)
			}
			newestAPIVersion := getNewestAPIVersion("networkpolicy", ac.APIVersion, deprecatedAPIs)
			if newestAPIVersion != "" {
				results = append(results, result{a.Namespace, a.Name, ac.APIVersion, newestAPIVersion})
			}
		}
	}

	psp, err := clientset.PolicyV1beta1().PodSecurityPolicies().List(listOptions)
	for _, a := range psp.Items {
		if a.Annotations[LastAppliedConfiguration] != "" {
			var ac appliedConfig
			d := []byte(a.Annotations[LastAppliedConfiguration])
			err := json.Unmarshal(d, &ac)
			if err != nil {
				log.Println(err)
			}
			newestAPIVersion := getNewestAPIVersion("podSecurityPolicy", ac.APIVersion, deprecatedAPIs)
			if newestAPIVersion != "" {
				results = append(results, result{a.Namespace, a.Name, ac.APIVersion, newestAPIVersion})
			}
		}
	}

	dp, err := clientset.AppsV1().Deployments("").List(listOptions)
	for _, a := range dp.Items {
		if a.Annotations[LastAppliedConfiguration] != "" {
			var ac appliedConfig
			d := []byte(a.Annotations[LastAppliedConfiguration])
			err := json.Unmarshal(d, &ac)
			if err != nil {
				log.Println(err)
			}
			newestAPIVersion := getNewestAPIVersion("deployment", ac.APIVersion, deprecatedAPIs)
			if newestAPIVersion != "" {
				results = append(results, result{a.Namespace, a.Name, ac.APIVersion, newestAPIVersion})
			}
		}
	}

	ds, err := clientset.AppsV1().DaemonSets("").List(listOptions)
	for _, a := range ds.Items {
		if a.Annotations[LastAppliedConfiguration] != "" {
			var ac appliedConfig
			d := []byte(a.Annotations[LastAppliedConfiguration])
			err := json.Unmarshal(d, &ac)
			if err != nil {
				log.Println(err)
			}
			newestAPIVersion := getNewestAPIVersion("daemonSet", ac.APIVersion, deprecatedAPIs)
			if newestAPIVersion != "" {
				results = append(results, result{a.Namespace, a.Name, ac.APIVersion, newestAPIVersion})
			}
		}
	}

	ss, err := clientset.AppsV1().StatefulSets("").List(listOptions)
	for _, a := range ss.Items {
		if a.Annotations[LastAppliedConfiguration] != "" {
			var ac appliedConfig
			d := []byte(a.Annotations[LastAppliedConfiguration])
			err := json.Unmarshal(d, &ac)
			if err != nil {
				log.Println(err)
			}
			newestAPIVersion := getNewestAPIVersion("statefulSet", ac.APIVersion, deprecatedAPIs)
			if newestAPIVersion != "" {
				results = append(results, result{a.Namespace, a.Name, ac.APIVersion, newestAPIVersion})
			}
		}
	}

	rs, err := clientset.AppsV1().ReplicaSets("").List(listOptions)
	for _, a := range rs.Items {
		if a.Annotations[LastAppliedConfiguration] != "" {
			var ac appliedConfig
			d := []byte(a.Annotations[LastAppliedConfiguration])
			err := json.Unmarshal(d, &ac)
			if err != nil {
				log.Println(err)
			}
			newestAPIVersion := getNewestAPIVersion("replicaSet", ac.APIVersion, deprecatedAPIs)
			if newestAPIVersion != "" {
				results = append(results, result{a.Namespace, a.Name, ac.APIVersion, newestAPIVersion})
			}
		}
	}

	jsonOutput, _ := json.Marshal(results)
	err = ioutil.WriteFile("/tmp/results/results", jsonOutput, 0644)
	if err != nil {
		fmt.Println("Cannot write results")
		os.Exit(1)
	}

	d1 := []byte("/tmp/results/results")
	err = ioutil.WriteFile("/tmp/results/done", d1, 0644)
	if err != nil {
		fmt.Println("Cannot write 'done' file")
		os.Exit(1)
	}
}

func getNewestAPIVersion(name, apiVersion string, deps map[string]apiChanges) string {
	for _, a := range deps[name].Deprecated {
		if apiVersion == a {
			return deps[name].New
		}
	}
	return ""
}

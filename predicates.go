package main

import (
	"log"
	"math/rand"
	"strings"

	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kube-scheduler/extender/v1"
)

const (
	// LuckyPred rejects a node if you're not lucky ¯\_(ツ)_/¯
	LuckyPred        = "Lucky"
	LuckyPredFailMsg = "Well, you're not lucky"
)

var predicatesFuncs = map[string]FitPredicate{
	LuckyPred: LuckyPredicate,
}

type FitPredicate func(pod *v1.Pod, node v1.Node) (bool, []string, error)

var predicatesSorted = []string{LuckyPred}

// filter filters nodes according to predicates defined in this extender
// it's webhooked to pkg/scheduler/core/generic_scheduler.go#findNodesThatFit()
func filter(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	var filteredNodes []v1.Node
	failedNodes := make(schedulerapi.FailedNodesMap)
	pod := args.Pod

	// TODO: parallelize this
	// TODO: hanlde error
	for _, node := range args.Nodes.Items {
		fits, failReasons, _ := podFitsOnNode(pod, node)
		if fits {
			filteredNodes = append(filteredNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",")
		}
	}

	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	fits := true
	failReasons := []string{}
	for _, predicateKey := range predicatesSorted {
		fit, failures, err := predicatesFuncs[predicateKey](pod, node)
		if err != nil {
			return false, nil, err
		}
		fits = fits && fit
		failReasons = append(failReasons, failures...)
	}
	return fits, failReasons, nil
}

func LuckyPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	lucky := rand.Intn(2) == 0
	if lucky {
		log.Printf("pod %v/%v is lucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name)
		return true, nil, nil
	}
	log.Printf("pod %v/%v is unlucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name)
	return false, []string{LuckyPredFailMsg}, nil
}

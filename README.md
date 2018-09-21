# sample-scheduler-extender

A sample to showcase how to create a k8s scheduler extender.

## [TODO] Running with a Kubeadm env


## [TODO] Running with a hack-local env (for dev)

```bash
# replace scheduler start commands in k8s.io/kubernetes/hack/hack/local-up-cluster.sh
    ${CONTROLPLANE_SUDO} "${GO_OUT}/hyperkube" scheduler \
      --v=${LOG_LEVEL} \
      --leader-elect=false \
      --kubeconfig "$CERT_DIR"/scheduler.kubeconfig \
      --feature-gates="${FEATURE_GATES}" \
      --master="https://${API_HOST}:${API_SECURE_PORT}" \
      --config="/root/scheduler-extender-config.yaml" >"${SCHEDULER_LOG}" 2>&1 &
```

## Notes

- Prioritize webhook won't be triggered if it's running on an one-node cluster. As it makes no sense to run priorities logic when there is only one candidate:

```go
// from k8s.io/kubernetes/pkg/scheduler/core/generic_scheduler.go
func (g *genericScheduler) Schedule(pod *v1.Pod, nodeLister algorithm.NodeLister) (string, error) {
    ...
	// When only one node after predicate, just use it.
	if len(filteredNodes) == 1 {
		metrics.SchedulingAlgorithmPriorityEvaluationDuration.Observe(metrics.SinceInMicroseconds(startPriorityEvalTime))
		return filteredNodes[0].Name, nil
    }
    ...
}
```
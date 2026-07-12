package prometheus

import "fmt"

func (c *Client) GetPodRestartMetric(namespace, podName string) (string, error) {
	query := fmt.Sprintf(
		`kube_pod_container_status_restarts_total{namespace="%s",pod="%s"}`,
		namespace,
		podName,
	)

	result, err := c.Query(query)
	if err != nil {
		return "", err
	}

	return result.FirstValue(), nil
}

func (c *Client) GetPodCPUUsage(namespace, podName string) (float64, error) {
	query := fmt.Sprintf(
		`sum(rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s",container!="",container!="POD"}[2m]))`,
		namespace,
		podName,
	)

	result, err := c.Query(query)
	if err != nil {
		return 0, err
	}

	return result.FirstValueFloat(), nil
}

func (c *Client) GetNodeCPUUsage(nodeName string) (float64, error) {
	query := fmt.Sprintf(
		`1 - avg by(instance) (rate(node_cpu_seconds_total{mode="idle",instance=~"%s.*"}[5m]))`,
		nodeName,
	)

	result, err := c.Query(query)
	if err != nil {
		return 0, err
	}

	return result.FirstValueFloat(), nil
}

func (c *Client) GetNodeMemoryUsage(nodeName string) (float64, error) {
	query := fmt.Sprintf(
		`1 - (node_memory_MemAvailable_bytes{instance=~"%s.*"} / node_memory_MemTotal_bytes{instance=~"%s.*"})`,
		nodeName,
		nodeName,
	)

	result, err := c.Query(query)
	if err != nil {
		return 0, err
	}

	return result.FirstValueFloat(), nil
}

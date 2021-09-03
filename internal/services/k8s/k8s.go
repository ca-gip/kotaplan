package k8s

import (
	underscore "github.com/ahl5esoft/golang-underscore"
	"github.com/ca-gip/kotaplan/internal/types"
	"github.com/ca-gip/kotaplan/internal/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func ClientGen(masterURL *string, kubeconfig *string) (client *kubernetes.Clientset, metricsClient *metrics.Clientset) {
	cfg, err := clientcmd.BuildConfigFromFlags(*masterURL, *kubeconfig)

	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	client, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
		syscall.Exit(1)
	}

	metricsClient, err = metrics.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes metrics client: %s", err.Error())
		syscall.Exit(1)
	}

	return
}

func DefaultKubeconfig() string {
	fname := os.Getenv("KUBECONFIG")
	if fname != "" {
		return fname
	}
	home, err := os.UserHomeDir()
	if err != nil {
		klog.Warningf("failed to get home directory: %v", err)
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}

func GetClusterData(client *kubernetes.Clientset, metricsClient *metrics.Clientset, settings *types.Parameters) (cluster *types.ClusterData, err error) {

	cluster = &types.ClusterData{
		Nodes:      nil,
		Namespaces: nil,
		Pods:       nil,
		PodsMetric: nil,
	}

	cluster.Nodes, err = getWorkerNodes(client)

	if err != nil {
		return
	}

	cluster.Namespaces, err = getNamespace(client, settings)

	if err != nil {
		return
	}

	for _, ns := range cluster.Namespaces.Items {
		pods, errPod := getPods(client, ns)
		if errPod != nil {
			return
		}

		podMetrics, errPodMetrics := getPodsMetric(metricsClient, ns)
		if errPodMetrics != nil {
			return
		}

		for _, pod := range pods.Items {
			cluster.Pods = append(cluster.Pods, pod)
		}

		for _, metricsPod := range podMetrics.Items {
			cluster.PodsMetric = append(cluster.PodsMetric, metricsPod)
		}

	}

	return
}

// TODO : update with something more generic
func getWorkerNodes(client *kubernetes.Clientset) (worker []v1.Node, err error) {
	nodes, err := client.
		CoreV1().
		Nodes().
		List(*&metav1.ListOptions{})

	if err != nil {
		return
	}

	underscore.
		Chain(nodes.Items).
		Where(func(node v1.Node, _ int) bool {
			return strings.Contains(node.Name, "worker")
		}).
		Value(&worker)

	return
}

func getNamespace(client *kubernetes.Clientset, settings *types.Parameters) (namespaces *v1.NamespaceList, err error) {
	namespaces, err = client.
		CoreV1().
		Namespaces().
		List(*&metav1.ListOptions{
			LabelSelector: utils.LabelsToString(settings.Labels),
		})
	return
}

func getPods(client *kubernetes.Clientset, namespace v1.Namespace) (pods *v1.PodList, err error) {
	pods, err = client.
		CoreV1().
		Pods(namespace.Name).
		List(*&metav1.ListOptions{})
	return
}

func getPodsMetric(metricsClient *metrics.Clientset, namespace v1.Namespace) (podsMetric *v1beta1.PodMetricsList, err error) {
	podsMetric, err = metricsClient.
		MetricsV1beta1().
		PodMetricses(namespace.Name).
		List(*&metav1.ListOptions{})
	return
}

package stub

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/autoscaling/v2beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateHPAWithResourceMetrics(t *testing.T) {

	annotations := map[string]string{
		"hpa.autoscaling.jdmeshcloud.io/minReplicas":                  "1",
		"hpa.autoscaling.jdmeshcloud.io/maxReplicas":                  "3",
		"cpu.hpa.autoscaling.jdmeshcloud.io/targetAverageUtilization": "80%",
		"memory.hpa.autoscaling.jdmeshcloud.io/targetAverageValue":    "1024Mi",
	}

	o := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}

	hpa := createHorizontalPodAutoscaler(o, o.GroupVersionKind(), annotations)

	if hpa == nil {
		t.Error("Error hpa is not created!")
		return
	}

	if len(hpa.Spec.Metrics) == 0 {
		t.Error("Error no metrics found!")
		return
	}

	for _, metric := range hpa.Spec.Metrics {
		if metric.Type != v2beta1.ResourceMetricSourceType {
			t.Errorf("Metric type expected: %v actual: %v", v2beta1.ResourceMetricSourceType, metric.Type)
		}
		if metric.Resource.Name != v1.ResourceMemory {
			t.Errorf("Metric name expected: %v actual: %v", v1.ResourceMemory, metric.Resource.Name)
		}
	}

}

func TestCreateHPAWithPodMetrics(t *testing.T) {

	annotations := map[string]string{
		"hpa.autoscaling.jdmeshcloud.io/minReplicas":               "1",
		"hpa.autoscaling.jdmeshcloud.io/maxReplicas":               "3",
		"pod.hpa.autoscaling.jdmeshcloud.io/customMetricIncorrect": "1024xMi",
		"pod.hpa.autoscaling.jdmeshcloud.io/customMetricCorrect":   "1024Mi",
	}

	o := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}

	hpa := createHorizontalPodAutoscaler(o, o.GroupVersionKind(), annotations)

	if hpa == nil {
		t.Error("Error hpa is not created!")
		return
	}

	if len(hpa.Spec.Metrics) == 0 {
		t.Error("Error no metrics found!")
		return
	}

	for _, metric := range hpa.Spec.Metrics {
		if metric.Type != v2beta1.PodsMetricSourceType {
			t.Errorf("Metric type expected: %v actual: %v", v2beta1.PodsMetricSourceType, metric.Type)
		}
		if metric.Pods.MetricName != "customMetricCorrect" {
			t.Errorf("Metric name expected: %v actual: %v", "customMetricCorrect", metric.Resource.Name)
		}
	}

}

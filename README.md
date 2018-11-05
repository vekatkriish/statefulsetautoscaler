### Statefulset Pod Autoscaler operator

You may not want nor can edit a Helm chart just to add an autoscaling feature. Nearly all charts supports **custom annotations** so we believe that it would be a good idea to be able to setup autoscaling just by adding some simple annotations to your deployment. 

We have open sourced a Statefulset Pod Autoscaler operator. This operator watches for your `Deployment` or `StatefulSet` and automatically creates an *HorizontalPodAutoscaler* resource, should you provide the correct autoscale annotations.

- Horizontal Pod Autoscaler operator
- Horizontal Pod Autoscaler operator Helm chart
### Autoscale by annotations

Autoscale annotations can be placed:

- directly on Deployment / StatefulSet:

 ```
  apiVersion: extensions/v1beta1
  kind: Deployment
  metadata:
    name: example
    labels:
    annotations:
      hpa.autoscaling.jdmeshcloud.io/minReplicas: "1"
      hpa.autoscaling.jdmeshcloud.io/maxReplicas: "3"
      cpu.hpa.autoscaling.jdmeshcloud.io/targetAverageUtilization: "70"
  ```

- or on `spec.template.metadata.annotations`:

 ```
  apiVersion: extensions/v1beta1
  kind: Deployment
  ...
  spec:
    replicas: 3
    template:
      metadata:
        labels:
          ...
        annotations:
            hpa.autoscaling.jdmeshcloud.io/minReplicas: "1"
            hpa.autoscaling.jdmeshcloud.io/maxReplicas: "3"
            cpu.hpa.autoscaling.jdmeshcloud.io/targetAverageUtilization: "70"
  ```  

The statefulset Pod Autoscaler operator takes care of creating, deleting, updating HPA, with other words keeping in sync with your deployment annotations.

### Annotations explained

All annotations must be prefixed with `autoscale`. It is required to specify minReplicas/maxReplicas and at least one metric to be used for autoscale. You can add *Resource* type metrics for cpu & memory and *Pods* type metrics.
Let's see what kind of annotations can be used to specify metrics:

- ``cpu.hpa.autoscaling.jdmeshcloud.io/targetAverageUtilization: "{targetAverageUtilizationPercentage}"`` - adds a Resource type metric for cpu with targetAverageUtilizationPercentage set as specified, where targetAverageUtilizationPercentage should be an int value between [1-100]

- ``cpu.hpa.autoscaling.jdmeshcloud.io/targetAverageValue: "{targetAverageValue}"`` - adds a Resource type metric for cpu with targetAverageValue set as specified, where targetAverageValue is a [Quantity](https://godoc.org/k8s.io/apimachinery/pkg/api/resource#Quantity).

- ``memory.hpa.autoscaling.jdmeshcloud.io/targetAverageUtilization: "{targetAverageUtilizationPercentage}"`` - adds a Resource type metric for memory with targetAverageUtilizationPercentage set as specified, where targetAverageUtilizationPercentage should be an int value between [1-100]

- ``memory.hpa.autoscaling.jdmeshcloud.io/targetAverageValue: "{targetAverageValue}"`` - adds a Resource type metric for memory with targetAverageValue set as specified, where targetAverageValue is a [Quantity](https://godoc.org/k8s.io/apimachinery/pkg/api/resource#Quantity).

- ``pod.hpa.autoscaling.jdmeshcloud.io/custom_metric_name: "{targetAverageValue}"`` - adds a Pods type metric with targetAverageValue set as specified, where targetAverageValue is a [Quantity](https://godoc.org/k8s.io/apimachinery/pkg/api/resource#Quantity).

> To use custom metrics from *Prometheus*, you have to deploy `Prometheus Adapter` and `Metrics Server`.

### Quick usage example
Let us deploy a sample example for statefulset app and let us see with these annotations
Deploy the statefulset app with these annotations on k8s with the following command
```
    kubectl create -f deploy/statefulset.yaml
```

  1. Check if HPA is created

   ```
    kubectl get hpa

    NAME      REFERENCE          TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
    web1      StatefulSet/web1   0%/70%    1         3         1          20m
  
  ```

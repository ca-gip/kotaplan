
# Kotaplan

The `kotaplan` cli will help you set up `kotary` if you have existing workload running on your cluster.

## Install

### Windows 

1. Download [here](https://github.com/ca-gip/kotaplan/releases/download/0.1.0/kotaplan.exe)
2. Open Cmd
```bash
.\kotaplan.exe -h
```

### Linux

```bash
sudo wget https://github.com/ca-gip/kotaplan/releases/download/0.1.0/kotaplan -P /usr/local/bin
sudo chmod a+x /usr/local/bin/kotaplan
kotaplan -h
```


## Feature

 __What it accomplish :__
 * List client namespace with :
    * Amount of pods
    * Total amount of request (CPU and Memory)
    * Current total amount of usage (CPU and Memory)
    * Does it fit the default claim
    * The recommended spec of the claim
 * List cluster resources :
    * Amount of worker nodes
    * Total amount of allocatable resources (CPU and Memory)
    * Maximum amount of resource per namespace (CPU and Memory) 
    * Total amount of allocatable resources with over-commit (CPU and Memory) 
 * Check that the recommended spec fit the maximum amount of resource per namespace
 * Give a result which check that the total amount of resources fit the total allocatable of the cluster 
 * Generate YAML manifest for all the namespace with the recommended spec
 
 #### How the recommended spec is calculate ?
 
 For a namespace to be marked as "fit default", its request must be 20% lower that the default claim
 If a namespace does not a fit the default claim the recommendation will add a margin of 20%.
 This is the default setting, it can be change via the `--margin` flag 
 
 ## How to use
 
Using default value for all the parameters
```bash
kotaplan
```

In case manifest should be generated
```bash
kotaplan --manifest
```

In case the default claim should be customized :
Example with a default claim with 6Gi of Memory and 1 CPU
```bash
kotaplan --memclaim 6 --cpuclaim 1000
```

Show help
```bash
kotaplan -h
```

All the arguments that can be use

| Argument | Type | Default | Description |
|----------|------|---------|-------------|
| memclaim | Integer | 8 | Amount of memory for the default claim in GiB |
| cpuclaim | Integer | 1000 | Amount of cpu for the default claim in Milli ex 1000 = 1CPU |
| memover | Float | 1 | Ratio of the Memory over or under commit ex: 1 is 100 % |
| cpuover | Float | 1 | Ratio of the CPU over or under commit ex: 1 is 100 % |
| ratioMemNs | Float | 1 | Ratio of the maximum amount of Memory that can be claim by a namespace ex 0.33 is 33% of the cluster |
| ratioCpuNs | Float | 1 | Ratio of the maximum amount of CPU that can be claim by a namespace ex 0.33 is 33% of the cluster |
| manifest | Boolean | false | Generate YAML manifest with recommended spec |
| margin | Float | 1.20 | Ratio of the margin for the recommended spec and to decide if the default claim passes |
| label | String | | Only match certain namespaces with a label ex ‘team=dev‘ |
| showconf | Boolean | false | Show the running configuration |
| kubeconfig | String | /home/$USER/.kube/config | Path to a kubeconfig |
| master | String |  | The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster |



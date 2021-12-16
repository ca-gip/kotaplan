
# Kotaplan

Visualize resource consumption and generated ResourceQuota with recommend spec

### Intallation

```bash
curl -ssf -L https://raw.githubusercontent.com/ca-gip/kotaplan/master/install.sh | bash
```

Adding completion for bash (also available for zsh and fish)
```
echo "source <(kotaplan completion bash)" >> ~/.bashrc
```

#### Windows

A Windows binary is also available checkout our latest [release](https://github.com/ca-gip/kotaplan/releases/latest/), you will find `kotaplan_<version>_windows_amd64.tar.gz`

## Feature

 __What it accomplishes :__
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
 
 #### How the recommended spec is calculated ?
 
 For a namespace to be marked as "fit default", its request must be 20% lower that the default claim
 If a namespace does not a fit the default claim the recommendation will add a margin of 20%.
 This is the default setting, it can be change via the `--margin` flag 
 
 ## How to use

```bash
kotaplan help
Visualize resource consumption and generated ResourceQuota with recommend spec

Usage:
  kotaplan [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  generate    Generate ResourceQuota manifest
  help        Help about any command
  view        Visualize resource consumption, policy compliance and recommended ResourceQuota

Flags:
      --config string                  config file (default is $HOME/.kotaplan.yaml)
      --default-claim-cpu int          Amount of CPU for the default claim in Milli. ex 1000 = 1CPU. (default 0)
      --default-claim-memory int       Amount of Memory for the default claim in GiB. (default 0)
  -h, --help                           help for kotaplan
      --kubeconfig string              Path to a kubeconfig (default "~/.kube/config")
  -l, --labels string                  Match namespace containing a label (default "quota=managed")
      --margin float                   Margin for the recommended spec (default 1.2)
      --master string                  Address of the Kubernetes API server. Overrides any value in kubeconfig
      --over-commit-cpu float          Ratio of the CPU over or under commit (default 1)
      --over-commit-memory float       Ratio of the Memory over or under commit (default 1)
      --ratio-namespace-cpu float      Ratio of the maximum amount of CPU that can be claim by a namespace. Ex: 0.5 meaning 50% of the cluster is claimable by a Namespace (default 1)
      --ratio-namespace-memory float   Ratio of the maximum amount of Memory that can be claim by a namespace. Ex: 0.5 meaning 50% of the cluster is claimable by a Namespace (default 1)

Use "kotaplan [command] --help" for more information about a command.
```

In case manifest should be generated
```bash
kotaplan generate
```

In case the default claim should be customized :
Example a default claim with 6Gi of Memory and 1 CPU
```bash
kotaplan view --default-claim-memory 6 --default-claim-cpu 1000
```

 <hr/>
<p align=center  style="background-color:#333333 !important;">
  <a href="https://www.jetbrains.com/">
  Developed with
  <br/>
  <img align="center" src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="drawing" width="100"/>
  </a>
</p>

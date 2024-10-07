# kubectl-cs - Quick Kubernetes context switcher  ☸️ 🔄

cs is a kubectl plugin makes it easy to switch between multiple [kubeconfig](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/) contexts or change the default 
namespace of the currently selected context.

It's heavily used while working with [CloudControl](https://cloudcontrol.dodevops.io). This is the rewrite in pure
go to replace the former bash and gnu-dialog-based tool and turn it into a kubectl plugin.

## Requirements

This tool builds upon a locally configured [kubeconfig](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/) file and the kubectl client.

## Installation

Download the binary for your operating system and architecture from the latest 
[kc release](https://github.com/dodevops/kc/releases) and put it somewhere in your path.

*Note*: Because the tool is not signed, macOS users need to enable execution when starting kc for the first time. 
To do this, Ctrl-Click on the binary and select "allow".

## Usage

Run `kubectl cs` to switch to a Kubernetes context and/or namespace. You can directly specify the context to switch to
in the form of `context:namespace`. If you leave either of the two parts empty, it will only change the remaining part.

Examples:

* `kubectl cs int:` - Switch to the currently active namespace in the `int` context
* `kubectl cs :kube-system` - Switch to the namespace `kube-system` in the currently active context
* `kubectl cs int:kube-system` - Switch to the namespace `kube-system` in the `int` context

If you don't specify anything, `kubectl cs` will present a selection list of all namespaces in all contexts, which can
be filtered and selected. If that takes too long, the flag `--only-current-context` can be used, which only enumerates
the namespaces of the current context.


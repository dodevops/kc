# kc - Quick Kubernetes context switcher  ☸️ 🔄

kc is a small tool that makes it easy to switch between multiple [kubeconfig](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/) contexts or change the default namespace
of the currently selected context.

It's heavily used while working with [CloudControl](https://cloudcontrol.dodevops.io). This is the rewrite in pure
go to replace the former bash and gnu-dialog-based tool.

## Requirements

This tool builds upon a locally configured [kubeconfig](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/) file and the Kubernetes client.1

## Installation

Download the binary for your operating system and architecture from the latest 
[kc release](https://github.com/dodevops/kc/releases) and put it somewhere in your path.

*Note*: Because the tool is not signed, macOS users need to enable execution when starting kc for the first time. 
To do this, Ctrl-Click on the binary and select "allow".

## Usage

Run `kc` to switch to a Kubernetes context. You can either directly specify the context to switch to or leave it empty
to select from a list of contexts from the kubeconfig file.

Run `kc -n` to switch the default namespace for the current context. You can either directly specify the namespace
to switch to or leave it empty to select from a list of namespaces available to you.

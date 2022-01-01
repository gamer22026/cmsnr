# cmsnr

## Description
cmsnr (pronounced "commissioner") is a lightweight framework for running OPA in a sidecar alongside your applications in Kubernetes.

## Purpose
This project gives Kubernetes users a simple way to deploy OPA policies for their apps. It has the ability to define which 
applications should have which policies, and allows for adding multiple policies into the same sidecar. When a policy is added 
or updated in the cluster, the client in each sidcar will check if the deployment name matches the name in their own deployment.
If the name is a match each client then uploads that policy into it's own OPA giving that Kubernetes deployment access to that policy.

## OPA Policy CRD
cmsnr uses an OPA policy CRD to store the Rego policy in the cluster. The CRD also takes a deployment name and a policy name.
The deployment name should match the deployment name in the pod annotation for the deployment/pod where you want the policy to 
be available. The policy name is the name cmsnr will use when putting the policy in OPA.

## Client
cmsnr uses the cli tool `cmsnrctl` to do all of it's work. It contains a lightweight client that will watch the cluster for new and
updated OPA policies and update them in the corresponding deployments.

## Pod Labels
cmsnr uses a mutating webhook to watch for pods with the annotation `cmsnr.com/inject: enabled`. Cmsnr will then inject two lightweight containers
in the pod: OPA and cmsnr itself. It injects the statically linked OPA container and cmsnr itself is just a statically linked binary. 

## Deploy cmsnr

To deploy cmsnr, first download the most recent version from the releases page. Then simply run
`cmsnrctl server deploy | kubectl apply -f -`

## Examples
To see the functionality of cmsnr, run download the most recent version. Then run
`cmsnrctl server deploy | kubectl apply -f -`. Then run `kubectl apply -f examples/`

This will create an annotated deployment and two OPA policies which will be injected into the sidecar.
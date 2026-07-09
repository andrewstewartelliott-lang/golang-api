The Assignment

As a new employee you’ve just inherited an incomplete internal tool. This tool’s purpose is to help the SRE team with the goa l of keeping the
workloads on our Kubernetes cluster reliable and secure.
You are tasked with extending this tool according to the user stories below:

(/) As an SRE I want to know whether all the deployments in the k8s cluster have as many healthy pods as requested by the respective `Deployment` spec

get active list of deployments and replicas, compare
output "all deployments checks complete'
automated test runs various output scenarios use existing framework

(/) As an SRE I want to always know whether this tool can successfully communicate with the configured k8s API server

enable client error handling and trap various errors?
try to interact with the api server with corev1 and trap errors

(/) As an application developer I want to be able to deploy this application into a Kubernetes cluster using Helm

push to oci? only docker push needed
create helm chart, charts/values, local instance of the helm chart and install locally


The incomplete tool can be found here: https://github.com/TykTechnologies/tyk-sre-assignment
1. Clone the repository
2. Choose your language of choice - available in Go
3. Complete at least 2 “As an SRE” and at least 1 “As an application developer” stories - extending the existing tool
4. Push your changes to your clone on Github and share with us before the day of your interview
Make sure all automated tests are passing before sharing the project.
During the interview you will have 20 minutes to present your project to the reviewers. The remaining time will be used to di scuss the operational 
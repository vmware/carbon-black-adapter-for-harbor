
# carbon-black-adapter-for-harbor

## Overview

Carbon Black adapter for Harbor integrates your Harbor Registry with the Carbon Black Cloud. It leverages Harbor's official [Pluggable Scanner API Spec](https://github.com/goharbor/pluggable-scanner-spec). to enable Harbor to scan images present in its registry and provide vulnerability reports for those images - right in the Harbor UI.

## Configuration

All configuration values for the Harbor Adapter can be set using the following environment variables at startup:

| Name               | Default | Description                                 |
| ---                | ---     | ---                                         |
| `CB_API_ID`        | ` `     | Carbon Black API ID                         |
| `CB_ORG_KEY`       | ` `     | Carbon Black ORG KEY                        |
| `CB_URL`           | ` `     | Carbon Black URL                            |
| `CB_API_KEY`       | ` `     | Carbon Black API KEY                        |
| `LOG_LEVEL`        | `info`  | Adapter service log level                   |


## Installation

### Requirements

* Kubernetes >= 1.14
* Harbor >= 1.10
* Helm >= 3
* Carbon Black Credentials - `CB_API_ID`, `CB_ORG_KEY`, `CB_URL`, `CB_API_KEY`

## Obtaining CB credentials

* Log in to the Carbon Black cloud.
* On the left side bar navigate to Inventory -> Kubernetes -> K8s Clusters
* Switch to CLI config and click on Add cli on the right top.
* Provide CLI name, Default build step, CLI description and generate credentials.
* Now you can get the following environment values from the generated credentials.
 
    `CB_API_ID`  - cb_api_id,
    `CB_ORG_KEY` - org_key,
    `CB_URL`     - saas_url,
    `CB_API_KEY` - cb_api_key

![Obtaining Credentials](images/obtaining_credentials.png)

## Try it out

### Helm install

The easiest was to deploy Harbor Adapter is through the `helm install` command and make sure you provide all the necessary arguments as mentioned below:

```
$ helm repo add cb-harbor-adapter https://projects.registry.vmware.com/chartrepo/cb_harbor_adapter
"cb-harbor-adapter" has been added to your repositories

$ helm install carbon-black --set cb_image_scanning.api_id=<YOUR_API_ID_HERE>,cb_image_scanning.org_key=<YOUR_ORG_KEY_HERE>,cb_image_scanning.api_key=<YOUR_API_KEY_HERE>,cb_image_scanning.url=<YOUR_URL_HERE>  cb-harbor-adapter/harbor-adapter
NAME: carbon-black
LAST DEPLOYED: Thu Apr 15 02:40:37 2021
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

You can also provide the configuration through values file and make sure you provide configuration values for the following parameters in values.yaml file:

* cb_image_scanning.api_id
* cb_image_scanning.org_key
* cb_image_scanning.api_key
* cb_image_scanning.url

```
$ helm repo add cb-harbor-adapter https://projects.registry.vmware.com/chartrepo/cb_harbor_adapter
"cb-harbor-adapter" has been added to your repositories

$ helm install carbon-black -f values.yaml cb-harbor-adapter/harbor-adapter
```

Here is the sample [values.yaml](helm/values.yaml) file

### kubectl apply

1. Add the base 64 encoded values for api_id, org_key, org_key, org_key into ./k8s/cb-harbor-adapter.yaml file
2. Run the following command:

```
kubectl apply -f ./k8s/cb-harbor-adapter.yaml
```

## Set up Harbor Adapter in Harbor Registry

1. Log in to your Harbor registry 
2. Navigate to Administration -> Interrogation Services on the left panel.
3. Click on new scanner
4. Provide the following -
    * Name - Carbon-Black
    * Description - eg - Scan images using Carbon Black image scanning service
    * Endpoint - http://SERVICE-NAME:8080 eg - http://carbon-black-harbor-adapter:8080

### Here is how you get the service name
```
$ kubectl get service
NAME                                   TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
carbon-black-harbor-adapter            ClusterIP   10.97.54.107     <none>        8080/TCP            8h
```
    * check on `Use internal registry address`

![Obtaining Credentials](images/obtaining_credentials.png)

5. After adding set the scanner as default.
6. Now the Adapter is ready for scanning.


## Contributing

Please follow [CONTRIBUTING.md](CONTRIBUTING.md)

## License

[Apache-2.0](LICENSE)


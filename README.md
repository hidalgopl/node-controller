## Running tests
To test node package, run `make test`.


## Building locally

### Binary
To build it locally use `make build`. Binary will be placed in `current-directory/bin/node-controller`.
### Docker image
To build docker image locally use `make container`. Container will be tagged `node-controller:latest`



## Environment variables

|    Env name    | Description | default |
| -------------- | ----------- | ------- |
| TARGET_OS | Tells controller for which OSImage should watch | Container Linux |
| LABEL_KEY | Tells controller which label it should set | kubermatic.io/uses-container-linux |
| LABEL_VALUE | Tells controller what value should be set for label | true |


## Credits
Based hugely on go-client example [workqueue](https://github.com/kubernetes/client-go/tree/master/examples/workqueue).

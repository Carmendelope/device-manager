# device-manager

This component is in charge of Device and Device group operations.

## Getting Started

### Prerequisites

* [`system-model`](https://github.com/nalej/system-model)
* [`authx`](https://github.com/nalej/authx)

### Build and compile

In order to build and compile this repository use the provided Makefile:

```shell script
make all
```

This operation generates the binaries for this repo, downloads the required dependencies, runs existing tests and generates ready-to-deploy Kubernetes files.

### Run tests

Tests are executed using Ginkgo. To run all the available tests:

```shell script
make test
```

### Integration tests 

To enable and run integration tests you will need to have running instances of all the prerequisites
and declare the following environment variables:

| Variable              | Example Value  | Description           |
| --------------------- | -------------- |---------------------- |
| RUN_INTEGRATION_TEST  | true           | Run integration tests |
| IT_SM_ADDRESS         | localhost:8800 | System Model Address  |
| IT_AUTHX_ADDRESS      | localhost:8810 | Authx Address         |

### Update dependencies

Dependencies are managed using Godep. For an automatic dependencies download use:

```shell script
make dep
```

In order to have all dependencies up-to-date run:

```shell script
dep ensure -update -v
```

## User client interface

To interact with this component, you can use the [public-api-cli](https://github.com/nalej/public-api).
The command that interacts with this component is `device`.

Example:
```shell script
./public-api-cli device info [deviceGroupID] [deviceID]
```

## Known Issues

## Contributing

Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.


## Versioning

We use [SemVer](http://semver.org/) for versioning. For the available versions, see the [tags on this repository](https://github.com/nalej/device-manager/tags). 

## Authors

See also the list of [contributors](https://github.com/nalej/device-manager/contributors) who participated in this project.

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.

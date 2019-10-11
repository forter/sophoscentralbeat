# Sophoscentralbeat

Welcome to Sophoscentralbeat.

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/logrhythm/sophoscentralbeat`

## Getting Started with Sophoscentralbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Sophoscentralbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Sophoscentralbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/logrhythm/sophoscentralbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Sophoscentralbeat run the command below. This will generate a binary
in the same directory with the name sophoscentralbeat.

```
make
```


### Run

To run Sophoscentralbeat with debugging output enabled, run:

```
./sophoscentralbeat -c sophoscentralbeat.yml -e -d "*"
```


### Test

To test Sophoscentralbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  Sophoscentralbeat source code, run the following command:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Sophoscentralbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/logrhythm/sophoscentralbeat
git clone https://github.com/logrhythm/sophoscentralbeat ${GOPATH}/src/github.com/logrhythm/sophoscentralbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.

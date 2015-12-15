# HaaS-Broker

**A service broker to front the easy vending of Hardware as a Service**

[![wercker
status](https://app.wercker.com/status/1bae75edbf6dcc279404af552247d8aa/s/master
"wercker
status")](https://app.wercker.com/project/bykey/1bae75edbf6dcc279404af552247d8aa)


## Running tests / build pipeline locally

```
# install the wercker cli
$ curl -L https://install.wercker.com | sh

# make sure a docker host is running
$ docker-machine start default && $(docker-machine env default)

# run the build pipeline locally, to test your code locally
$ ./testrunner

```


## Running locally for development

### (not yet hooked up)

```

# install the wercker cli
$ curl -L https://install.wercker.com | sh

#lets bootstrap our repo as a local dev space
$ ./init_developer_environment

# make sure a docker host is running
$ docker-machine start default && $(docker-machine env default)

# run the app locally using wercker magic
$ ./runlocaldeploy local_wercker_configs/myenv

$ echo "open ${DOCKER_HOST} in your browser to view this app locally"

```

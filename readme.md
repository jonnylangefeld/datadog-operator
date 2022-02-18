# Datadog Operator

The [official Datadog Operator](https://github.com/DataDog/datadog-operator)

> aims to provide a new way of deploying the Datadog Agent

but it doesn't allow to configure and manage resources like dashboards or monitors. That's where this
operator comes in.

### Run

`make run` will run the controller against the current cluster in the kube context. If this is your first time executing
the controller, make sure to run `make install` before.
The controller relies on the datadog api key and the datadog application key. You can either set the environment
variables `DATADOG_API_KEY` and `DATADOG_APPLICATION_KEY`, or create a file like the following:

```json
{
  "api_key": "<your-api-key>",
  "application_key": "<your-application-key>"
}
```

The default location for this file is `.secrets.json`, but can be overwritten with the controller flag `--secrets-path`.

### Deploy

Run `make deploy` to deploy the controller on the current cluster in the kube context.

### Debug

Run the `TestAPIs` tes in the `controllers/datadog/suite_test.go` file with a debugger of 
your choice (like IDEs or delv).
If you set a breakpoint somewhere and want to examine what the API server sees at this point in time, 
you will need to connect to the local API server running on your local machine. 
This API server was started by the testEnv.Start() function in the suite_test.go file. TestAPIs 
will automatically writes a kubectl file at ./controllers/cluster/kubectl, you can access the local 
API server during a debugging session using this kubectl file:

```
./controllers/cluster/kubectl get ns
```

# Demo: Running Containers with Scripts

You use the `docker` command line to work with containers on a single Docker server. 

A simple way to model your apps is to put your `docker` commands in Bash or PowerShell scripts - but the approach is limited.

## Running apps from a script

The demo app runs across multiple containers. 

[app-v1.ps1](app-v1.ps1) is a PowerShell script which automates the startup.

_Run the app:_

```
./app-v1.ps1
```

> Output is a list of object IDs - you need to know the script to understand what's happening.

_Check the app containers:_

```
docker ps
```

> Try the app at http://localhost:8080

Run the script again:

```
./app-v1.ps1
```

> Is the app running?

## Injecting configuration settings

[app-v2.ps1](app-v2.ps1) models the same app for a test environment using configuration files in volume mounts.

_Replace the deployment_

```
docker rm -f $(docker ps -aq)

./app-v2.ps1
```

> Refresh http://localhost:8080

## Modelling reliability

For production, containers should have limited access to compute resources. These containers can see all the host machine's memory:
```
docker stats products-api-test
```

And containers should restart if the app inside fails.

```
docker exec -it stock-api-test kill 1
```

> Try the app again http://localhost:8080

The web app can't access the API, so the page fails to load.

```
docker logs web-test
```

[app-v3.ps1](app-v3.ps1) models the same app with some additional reliability settings.

```
docker rm -f $(docker ps -aq)

./app-v3.ps1
```

Now memory is limited:

```
docker stats products-api-test
```

And if containers exit they are restarted:

```
docker exec -it stock-api-test kill 1

docker ps
```

> Try the app again http://localhost:8080

The app can now survive component failures and server restarts - but the PowerShell script is getting difficult to manage.
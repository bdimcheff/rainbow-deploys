# Rainbow Deploys for Kubernetes

or: how you can deploy services to Kubernetes that require long periods of draining.

## What?

Rainbow deploys are like [Blue/Green](https://martinfowler.com/bliki/BlueGreenDeployment.html) deploys, but instead of just two environments, there are an infinite number of colors.  Kubernetes makes this pretty easy to do.

## Why?

In an ideal world, everybody runs stateless services that have short request/response cycles.  In the real world, sometimes you need long-running connections and state.  You may not wish to just restart your backends if they have established connections for a variety of reasons.  See my [blog post](http://brandon.dimcheff.com/2018/02/rainbow-deploys-with-kubernetes/) on the topic for more info about why you might want to do this.

## TL;DR

You can drain stuff by changing a Service's selector but leaving the Deployment alone.  Instead of changing a Deployment and doing a rolling update, create a new deployment and repoint the Service.  Existing connections will remain until you delete the underlying Deployment.

## Prerequisites

1. minikube (or another kubernetes)
2. docker
3. make

## Demo

Included in this repo is a small go app that serves http on port 8080 and a simple tcp protocol on 8081.  If you visit :8080, you'll see the hex color for the first 6 characters of the HEAD of git when the docker image was built.  If you telnet to :8081, you'll see the color's hex value printed every 5 seconds for as long as you mantain a connection.

1. Start minikube with `minikube start`, which should also configure your kubectl.
1. Run `eval $(minikube docker-env)` so that you don't have to push images to a real docker repo.
    1. If you're not using minikube, you may wish to `export DOCKER_IMAGE=your-docker/image`, since you won't be able to push to my repo.  You'll also need to modify some commands below and figure out how to access your services.
1. Run `make image` which will build the image for you.
1. Run `make install` to install the app's Deployment and Service into kubernetes.
1. `kubectl get deployments` should yield something like
    ```
    $ kubectl get deployments
    NAME                     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
    rainbow-deploys-3c3fdc   2         2         2            2           1m
    ```
    where 3c3fdc is the first six characters of the git sha.  `kubectl get services` should show you the service that was created:
    ```
    NAME              TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                         AGE
    rainbow-deploys   NodePort    10.97.3.60   <none>        8080:31080/TCP,8081:31081/TCP   1m
    ```
1. Run `minikube service list` and find your http service on port 31080: 
    ```
    |-------------|----------------------|--------------------------------|
    |  NAMESPACE  |         NAME         |              URL               |
    |-------------|----------------------|--------------------------------|
    | default     | kubernetes           | No node port                   |
    | default     | rainbow-deploys      | http://192.168.99.100:31080    |
    |             |                      | http://192.168.99.100:31081    |
    | kube-system | kube-dns             | No node port                   |
    | kube-system | kubernetes-dashboard | http://192.168.99.100:30000    |
    |-------------|----------------------|--------------------------------|
    ```
    In this case, visit http://192.168.99.100:31080 in your browser and check out the color!  Leave this tab open.
1. Run `telnet <minikube IP> 31081` and you should see
    ```
    The color is #3c3fdc
    The color is #3c3fdc
    ```
    etc. printed every 5 seconds.  Leave this running in a terminal somewhere.  Note that our server is tempremental, so if you try to talk to it, it'll just stop working.  We'll file a ticket.
1. Now here's the fun part.  Make a change, any change to your git repo and commit it.  Touch a file, add an emoji to the log message, whatever you'd like.  Rerun `make image` and `make install` and you should see 2 deployments in `kubectl get deployments`:
    ```
    NAME                     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
    rainbow-deploys-3c3fdc   2         2         2            2           5m
    rainbow-deploys-9d2cc9   2         2         2            2           1m
    ```
    We have a new color: 9d2cc9!  We still only have one `rainbow-deploys` service, though, visible with `kubectl get services`:
    ```
    NAME              TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                         AGE
    rainbow-deploys   NodePort    10.97.3.60   <none>        8080:31080/TCP,8081:31081/TCP   4h
    ```
1. Reload your browser, and the color should change to match.  Look at your terminal with telnet, though.  It's still logging the old color.  Since this tcp connection was established before the service was updated, it'll still be active and pointed at the old deployment, even though the service points at the new one.
1. Open a new terminal window and `telnet <minikube IP> 31081` again.  You should see the new color logged in this terminal, but your old terminal will still be displaying the old color. All new connections will go to the new process.
1. Run `kubectl delete <older deployment>`, in this case `rainbow-deploys-3c3fdc`.  Your original telnet session should close, but your newer one should be unaffected.
1. Congratulations, you have successfully deployed your slow-drain service!

## How?

Look at `app.yaml`, the Kubernetes config for this project.  It contains a Service and a Deployment with a key feature:  there's a `color` label on the Deployment and it's used in the Service's selector.  This will cause Kubernetes to point the Service at the pods that match the current color.  Since the old Deployment and Pods are still around, existing TCP connections will remain established until they're closed from either end.

When you run `make install`, the task inserts the current git-derived color via a sed command: `cat app.yaml | sed s/__COLOR__/$(COLOR)/g | kubectl apply -f -`.  Each time `make install` is run, the latest `HEAD` is used in the selector, but it never modifies or deletes the old Deployment.  This, of course, means that you will eventually have a ton of old deployments.  Cleaning these up is an excercise left to the reader.
# Setting up a GoShimmer docker network

This page describes how to setup your own GoShimmer docker network.
We only describe the method where all nodes run on one physical machine. 


| Contents                                                                        |
|:------------------------------------------------------------------------------- |
| [Why you should run a docker network](#why-you-should-run-a-docker-network)                           |
| [Installing Docker](#installing-docker)           |
| [Cloning the goshimmer repository](#cloning-the goshimmer-repository)     |        

| [Setting up the Grafana dashboard](#setting-up-the-grafana-dashboard)           |

## Why you should run a docker network
Running a own docker network gives you complete control on the network, e.g., low PoW. It also enables you to test your own code before it has to survive in the wild. Moreover, it allows you to learn about goshimmer and follow some of the tutorials without having to run a own node in the Pollen network. 

## Installing Docker

On linux this can be done as follows. In case of problems checc out the internet for solutions;  https://www.docker.com/get-started or https://docs.docker.com/docker-for-mac/apple-m1/

Install needed dependencies:
```
$ apt-get install \
     apt-transport-https \
     ca-certificates \
     curl \
     gnupg-agent \
     software-properties-common
```

Add Docker’s official GPG key:
```
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
```

Verify that the GPG key matches:
```
$ apt-key fingerprint 0EBFCD88
pub   rsa4096 2017-02-22 [SCEA]
      9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88
uid           [ unknown] Docker Release (CE deb) <docker@docker.com>
sub   rsa4096 2017-02-22 [S]

```

Add the actual repository:
```
$ add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
```


And finally, install docker:
```
$ apt-get install docker-ce docker-ce-cli containerd.io
```

On windows-subsystem for Linux (WSL2) it may be necessary to start docker seperately:
```
$ /etc/init.d/docker start
```
Note, this may not work on WSL1.

Check whether docker is running by executing `docker ps`:
```
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```

## Cloning the goshimmer repository
Go to the folder where you want to clone the repository and do 
`git clone https://github.com/iotaledger/goshimmer.git`. If you want to clone the branch containing more tutorials do `git clone https://github.com/NaitsabesMue/goshimmer.git`







## Starting the GoShimmer docker network

Within the `/tools/docker-network` folder of goshimmer simply type `./run.sh n `. This will start n+2 nodes. (So you have to replace n by some integer).

Parameters of this docker network can be set in `config.docker.json.` For instance, one can change the PoW difficulty with `pow`. The smaller this numer (integer) the less difficult it is to solve the PoW puzzle. 

Some information on the node and the network are reachable through your standard browser. Check out
1. the analysis sever and autopeering visualizer: http://localhost:9000
2. `master_peer's` dashboard: http://localhost:8081
3. `master_peer's` web API: http://localhost:8080




It will be useful later on to know which services do run on which ports:


| Port  | Functionality  | Protocol |
| ----- | -------------- | -------- |
| 14626 | Autopeering    | UDP      |
| 14666 | Gossip         | TCP      |
| 10895 | FPC            | TCP/HTTP |
| 8080  | HTTP API      | TCP/HTTP |
| 8081  | Dashboard       | TCP/HTTP |
| 6061  | pprof HTTP API | TCP/HTTP |



#### Dashboard
The dashboard of your GoShimmer node should be accessible via `http://<your-ip>:8081`. If your node is still synchronizing, you might see a higher inflow of MPS.

![](https://i.imgur.com/8xAvi7X.png)

After a while, your node's dashboard should also display up to 8 neighbors:
![](https://i.imgur.com/gAyAXK9.png)


#### HTTP API
GoShimmer also exposes an HTTP API. To check whether that works correctly, you can access it via http://localhost:8080/info which should return a JSON response in the form of:
```
{
  "version": "v0.2.0",
  "synced": true,
  "identityID": "69RxiehGQ2c",
  "publicKey": "52Gzw9bo7k2dARFi4yxtt3B8xMht5UeFQX7pWdLFnxV5",
  "enabledPlugins": [
    "Analysis-Client",
    "Autopeering",
    "CLI",
    "Config",
    "DRNG",
    "Dashboard",
    "Database",
    "Gossip",
    ...
    "WebAPI info Endpoint",
    "WebAPI message Endpoint"
  ],
  "disabledPlugins": [
    "Analysis-Dashboard",
    "Analysis-Server",
    "Banner",
    "Bootstrap",
    "Faucet",
    "WebAPI Auth"
  ]
}
```



##### Stopping the node
```
CTR-C
```

## Spamming
Let us construct a real Tangle; issuing some messages. 

It is possible to send messages to the local network via the `master_peer` and observe log messages either 
via `docker logs --follow CONTAINER` or all of them combined when running via
1. Go to `/tools/spammer`, use `main.go`, and spam..
2. a. Enable spammer plugin of `peer_replica` in `docker-compose.yml`: `--node.enablePlugins=bootstrap,spammer...`
b. `master_peer` can be set to spam in the browser. For the replica attach to the shell the respective docker 
`docker container exec -it containername /bin/bash`
then 
`curl "http://127.0.0.1:8080/spammer?cmd=start&mpm=100"`

Here, `mpm` mean messages per minute. The spammer is stopped with `curl "http://127.0.0.1:8080/spammer?cmd=stop"`.







After built, the Docker image can be run with the command:
    
    docker run --rm -ti -p 9000:8000 <image name>
    
the example service will listen on port 9000 of localhost (or docker-machine ip)

    curl 192.168.99.100:9000/timed/hello/routed
    hello routed :)
    route time = 2.967
    
    curl 192.168.99.100:9000/timed/bye/routed
    bye routed :(
    route time = 500.998726ms
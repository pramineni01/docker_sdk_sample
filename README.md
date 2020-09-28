# docker_sdk_sample
This sample depicts how to pull a docker image, create a container instance, start it, and wait till it executes or till deadline exceeds.
Typical commands can be any of the below:
  - ./docker_sdk_sample run alpine sleep 10
  - ./docker_sdk_sample run alpine:latest sleep 10
  - ./docker_sdk_sample run -k 5 alpine sleep 10
  - ./docker_sdk_sample run -k 15 alpine sleep 10
  

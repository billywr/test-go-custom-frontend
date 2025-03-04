# Custom Frontends in BuildKit

Custom frontends in BuildKit are images that, when launched as containers via the `buildctl build` command, communicate with and send requests to the BuildKit daemon.

use case

![img.png](img.png)

## Section One: Custom Frontend Preparation

I have this [sample code](https://github.com/billywr/test-go-custom-frontend/blob/master/main.go) that I used to generate an `.exe` file, which is required to create a custom frontend image.

For testing purposes, you can clone the repository:

```
git clone https://github.com/billywr/test-go-custom-frontend
cd test-go-custom-frontend
```

Navigate to the main repository and run the following command to build an `.exe` file:

`GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o wcow-frontend.exe`

## Creating the Dockerfile
Create a `Dockerfile` where the generated `wcow-frontend.exe` will be used to create the custom image. Ensure that both the `Dockerfile` and `wcow-frontend.exe` are in the same directory.

Here are the contents of the `Dockerfile:`

```aiignore
FROM mcr.microsoft.com/windows/nanoserver:ltsc2022
USER ContainerAdministrator
COPY wcow-frontend.exe /wcow-frontend.exe
ENTRYPOINT ["/wcow-frontend.exe"]
```

## Building and Pushing the Image
Use the following `buildctl` command to build an image from the Dockerfile and push it to a Docker registry:

```aiignore
buildctl build `
    --frontend=dockerfile.v0 `
    --local context="buildContextPath" `
    --local dockerfile="PathToDockerfile" `
    --output type=image,name=docker.io/username/imageName:latest,push=true

```

At this point, you should have your custom frontend Docker image `imageName:latest` available in the Docker registry.

# Section Two: Running the Custom Frontend
You can accomplish this using either buildctl or Docker CLI.

## Prerequisites
The BuildKit daemon should be running in a way that allows containers on its host machine to access it. This is because the custom frontend, running inside a container, must send requests to the BuildKit daemon via gRPC.

A possible way to ensure this works is by stopping any running BuildKit instance and exposing `buildkitd` via the host machine's IP address using the following command:

`start-process buildkitd -ArgumentList "--addr", "tcp://<host_machine_IP>:1234"`

Verify that BuildKit has started with the correct IP configuration and is listening for requests on the specified port:
`netstat -an | findstr 1234`

## Running the Custom Frontend
Once everything is set up, you can run the following command:

```aiignore
buildctl build `
    --frontend=gateway.v0 `
    --opt source=docker-image://username/cfrontend03:latest `
    --local context=. `
    --local dockerfile=.
```

This should ensure that the custom frontend image sends requests to BuildKit and displays some of the `fmt` statements from the `main.go` of the source code in section one.

# Note:
**Section Two** is still under testing; therefore, this document will be updated in due course.

You are welcome to give your feedback and suggestions.
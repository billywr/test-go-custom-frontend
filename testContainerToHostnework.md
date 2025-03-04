## Findings: Testing container-to-host network connectivity in Windows BuildKit builds

# Test setup:

Dockerfile:

```aiignore
FROM mcr.microsoft.com/windows/nanoserver:ltsc2022
RUN curl.exe http://<hostmachineIP>:1234
```

(Example using host machine IP `172.20.176.1`)

Build environment:

buildkitd running on the host, listening on 172.20.176.1:1234.
`start-Process buildkitd -ArgumentList "--addr tcp://172.20.176.1:1234 --debug"`

running  `buildctl`

```aiignore
buildctl build `
  --frontend=dockerfile.v0 `
  --local context="E:\dockerfiles\demo-a\test-network\" `
  --local dockerfile="E:\dockerfiles\demo-a\test-network\" `
  --output type=image,name=docker.io/100909/testnetwork:latest,push=false
```

# Test results:

```aiignore
[2/2] RUN curl.exe http://172.20.176.1:1234:
curl: (7) Failed to connect to 172.20.176.1 port 1234 after 3 ms: Could not connect to server
error: failed to solve: process "cmd /S /C curl.exe http://172.20.176.1:1234" did not complete successfully: exit code: 7

```

# Observations:
❌ The build container was unable to connect to the host machine's IP (172.20.176.1) on port 1234.
✅ Verified that `buildkitd` was actively listening on `172.20.176.1:1234` from the host itself.
❌ Despite correct IP and port, no connection from inside the build container was possible.

# Root cause analysis:
BuildKit's build containers are network-isolated from the host during builds, especially on Windows.
Windows BuildKit builds often run in a NAT network or even Hyper-V isolation(unverified), preventing direct access to host machine IPs.
There is no `--network=host` support during builds on Windows, unlike Linux, which limits network configuration options.
Therefore, a build step attempting to connect back to the host machine (even via its own IP) is blocked or unreachable.

# Key insight:
On Windows, build steps (including RUN commands and custom frontends running during the build) cannot reliably reach the host machine's network services, including the same `buildkitd` instance that runs the build.
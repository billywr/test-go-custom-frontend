FROM mcr.microsoft.com/windows/nanoserver:ltsc2022
USER ContainerAdministrator
COPY wcow-frontend.exe /wcow-frontend.exe
ENTRYPOINT ["/wcow-frontend.exe"]
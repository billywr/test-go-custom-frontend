FROM mcr.microsoft.com/windows/nanoserver:ltsc2022
USER ContainerAdministrator
# Copy the custom frontend binary into the container
COPY wcow-frontend.exe /wcow-frontend.exe

# Set the entrypoint for BuildKit to run the frontend
ENTRYPOINT ["/wcow-frontend.exe"]
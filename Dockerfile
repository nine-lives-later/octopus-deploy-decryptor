FROM octopusdeploy/octopusdeploy:2024.3.12812

ENTRYPOINT [ "/bin/sh", "-c" ]
CMD [ "cp /Octopus/Octopus.* /flx" ]

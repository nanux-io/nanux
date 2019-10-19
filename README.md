# Nanux

Nanux is a toolkit fordeployig microservice.

*Work in progress*

## Development

### Installation

1. Create the folde which bin the docker container GOPATH: `mkdir ./goflow/build/dev/.go`
2. Launch the docker container : `docker-compose up`
3. Enter in the container: `docker-compose exec nanux bash`
4. Install dev dependencies in the container: `./goflow/scripts/dev/install-container-stack`
    > *delve* for debuging
    > *reflex* for watching file changes and reloading
5. Launch watcher to watch file changes and executing tests: `./goflow/scripts/dev/watch`

### Editor

The project used the go modules, in order to have tools working in your
editor don't forget to add the container GOPATH to the host GOPATH when
working on the project.

#### VSCode

For vscode, add the followin setting for the workspace:

```json
{
  "go.gopath": "${env:GOPATH}:${workspaceFolder}/goflow/build/dev/.go"`
}
```

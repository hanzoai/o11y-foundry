# CUE

## Installation
To use this demo you need cue-lang installed `brew install cue-lang/tap/cue`

## Execution
- Edit the .env file to define the variables
- Run cue export on the file you want to alter (right now only zookeeper is working)

```bash
source .env
cue export zookeeper.cue -e output --out text -t ZOO_SERVER_COUNT=$ZOO_SERVER_COUNT

# Produces
dataDir=/path/test
dataLogDir=/path/logs
clientPort=2181
servers=zookeeper-1:2888:3888,zookeeper-2:2888:3888
```

## Notes
- Cue has no native way to loop on an integer so `list.Range` had to be imported
- Cue assumes all elements passed to it are strings, even tho they are variables so conversion is needed
- Cue documentation is all over the place and LLMs have little to no idea on how to use it which makes development hard
- We could probably make this easier and faster with either `go-templates` or `python + jinja2`



# go-datatrails-demos

## How to run a demo

**Please note: all demos within go-datatrails-demos require golang to run**

If you do not have golang locally you can run them in a dockerised environment e.g.:

```
docker run -v ./integrity:/usr/src/myapp -w /usr/src/myapp  golang:1.22-alpine go run .
```

Where `integrity` is the demo to be run and it is run from the root of the repo.

### With Task

Prerequisite: 
1. install task: https://taskfile.dev/installation/

Steps:
1. run the task rune for the demo you want to run, e.g.:

```
task demos:integrity
```

### Without Task

If you do not have task installed, you can run the demo the following way:

1. Change directory to the demo you want to run, e.g.:

```
cd integrity
```

2. Compile and run the compiled demo executable:

```
go run .
```

## Integrity Demo

The integrity demo will verify the integrity of a datatrails event.

This is achieved by creating an inclusion proof for that event and
then verifying the inclusion proof against the merkle log.

If the inclusion proof is verified successfully, we can say that the
datatrails event is included on the merkle log.
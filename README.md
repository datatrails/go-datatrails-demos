# go-datatrails-demos

## Integrity Demo

The integrity demo will verify the integrity of a datatrails event.

This is achieved by creating an inclusion proof for that event and
then verifying the inclusion proof against the merkle log.

If the inclusion proof is verified successfully, we can say that the
datatrails event is included on the merkle log.

### Docker Demo
To run the integrity demo with docker:

```
docker run -v ./integrity:/usr/src/myapp -w /usr/src/myapp  golang:1.22-alpine go run .
```

Where the docker command is run from the root of the repo.

### Task Demo

To run the demo with a task (https://taskfile.dev/installation/) rune:

```
task demos:integrity
```

### Go Demo

If you want to run the demo directly with golang from the root of the repo:

```
cd integrity
go run .
```
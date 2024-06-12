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

## Completeness Demo

The completenesss demo will verify the integrity of a list of datatrails events.

This is achieved by checking that **ONLY** the events in the list exist on the log within a range
starting from the first event in the list, ending at the last event in the list.

Every event in the list then has an inclusion proof generated.

If all the inclusion proofs are verified successfully, we can say that the list of
datatrails events is complete and all are included on the merkle log.

### Docker Demo
To run the completeness demo with docker:

```
docker run -v ./completeness:/usr/src/myapp -w /usr/src/myapp  golang:1.22-alpine go run .
```

Where the docker command is run from the root of the repo.

### Task Demo

To run the demo with a task (https://taskfile.dev/installation/) rune:

```
task demos:completeness
```

### Go Demo

If you want to run the demo directly with golang from the root of the repo:

```
cd completeness
go run .
```

## Consistency Demo

The consistency demo will verify a future log state continues to be consistently recorded based on
an existing signed log state.

### Docker Demo
To run the consistency demo with docker:

```
docker run -v ./consistency:/usr/src/myapp -w /usr/src/myapp  golang:1.22-alpine go run .
```

Where the docker command is run from the root of the repo.

### Task Demo

To run the demo with a task (https://taskfile.dev/installation/) rune:

```
task demos:consistency
```

### Go Demo

If you want to run the demo directly with golang from the root of the repo:

```
cd consistency
go run .
```
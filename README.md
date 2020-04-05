## Running the project

Run the project (assuming you have cloned and at the root of the project), you could do:

`docker build -t api ./`

followed by

`docker run --dns 8.8.8.8 -it --publish 4000:8080 api`

Note: to (re)generate the golang and the JS code, as it is you would need golang protoc and protoc-gen-grpc-web. At that point you could run

`protoc -I="." proto/joke.proto --js_out=import_style=commonjs:. --grpc-web_out=import_style=commonjs,mode=grpcwebtext:. --go_out=plugins=grpc:.`

This generates the js and the golang code inside the proto folder, for you to move further accordingly

## Project structure

`proto` contains the initial data model (as proto definitions) for the problem

`generated` contains the server files generated from the proto definitions and copied into the project manually

## Architecture

[Protobufs](https://developers.google.com/protocol-buffers) are used as a higher level data description language. Chosen for the tooling available, one of which is the out of the box encription and decription (e.g. from json when communicating with the external jokes API, and to / from the file system when persisting)

[gRPC](https://grpc.io/) together with [gRPC-web](https://grpc.io/) were chosen as a framework for the requests and for scaffolding of a JS client

Used is a channel for jokes, that gets filled by a goroutine, and consumed on-demand when a request comes

The jokes are stored on the file system (this folder, in a filed called `jokes.cache`)by the main thread for the initial fetch, and by the goroutine that fetches regularly thereafter

The server also keeps an in-memory representation of the available jokes, so that requsts are just a bit faster

The docker setup tries caching the fetch of dependencies for subsequent builds where only app's code has been changed

## Alternatives

No higher level abstraction / [GraphQL](https://graphql.org/) / [JSON-schema](https://json-schema.org/) / [Flatbuffers](https://google.github.io/flatbuffers/) / other for defining the data, json (or different) for (de)serialization. Both have good tooling built around them for scaffolding clients and servers. For the purposes of this task any of these is as good as the other ones, graphql and protobufs being more mainstream.

Only writing to the FS from the thread that polls (including the very initial request for multiple jokes). Only reading the file from the main thread. This would allow easier extraction of the polling into a different process.
If the polling were to be extracted into a separate process as it is, there might exist a possibility of a race condition between that process and the one running the server. E.g. the server might check if the file is there, not find it, fetch data. In the meantime the other proces might also check, not find a file. The server could write the list of jokes, but then the polling process - since it did not find the file at first - could overwrite it.

## Shortcuts

The browser can't interpret protobuf responses natively and thus can't connect to grpc services. A workaround is use a proxy, as shown in the [grpc-web docs](https://github.com/grpc/grpc-web). That complicates development. I've used [an alternative](https://medium.com/safetycultureengineering/proxy-grpc-web-directly-in-your-go-server-without-envoy-7fbec326cb21) and have ignored having a flag for it, so that it's only used during dev.

Testing has been skipped

## Extra things that would be nice for a larger / real project

Having a makefile that standartises the project running / building, and other things - like code generation

Different (appropriate for mid / big projects, overkill for this one) - perhaps [this one](https://github.com/golang-standards/project-layout)

Layered, constently formatted, persistant logging. Likely something that wont be used directly, but through a wrapper (would mean underlying system for logging can be swapped later on ifneeded). Probably [zap](https://github.com/uber-go/zap)

Passing params from the command line (port, cache path, polling time, log level etc.). Having cli passing be done in a consistent way throughout projects is a good idea in general - [viper](https://github.com/spf13/viper) is an option

Ensuring responses can be debugged easily in the browser - as a minimum documenting the steps / extra plugins needed for the browser

Testing. Comfort with the language should mean TDD saves some time while developing, and the bigger the codebase / the more time it takes for the project to spin / the more infrequent or time consuming releases are, the more tests pay off. I would not recommend 100% test coverage - though covering at least the error paths is good (which usually means - most of the code branches). Note - if the problem is more complex, and the logic is likely to evolve, tests can turn out to be something that would better be discarded. Think of a function that receives arguments and formats them, but later on starts receiving the formatted version instead.

Separate image for a proxy for the grpc-web requests

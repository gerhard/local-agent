# Demo: Dagger + Docker Model Runner

This is a very small demo of [Dagger](https://dagger.io) and the new Docker Model Runner feature from Docker Desktop.

This demo is an agent in which we pass any directory, and it will spin up a container with all the required tools to work on the code base.
And contrary to classical dev environments, you have nothing to configure, the model will analyze the codebase for you, understand what tools are needed, and install them.

## Prerequisites

### Docker Model Runner

You need to have Docker Desktop 4.40 or newer installed.

> [!IMPORTANT]
> You need an Apple Silicon mac to run the Model Runner.

You then need to enable the model runner in the Docker Desktop settings, under "Features in development" section.

Once that's done, check the Model Runner is running:

```console
docker model status
```

You should get this output:
```text
Docker Model Runner is running
```

### Model

Now, pull the model used in this demo:

```console
docker model pull eunomie/qwen2.5-coder-14b-instruct:q5_k_m
```

And ensure this is working:

```console
docker model run eunomie/qwen2.5-coder-14b-instruct:q5_k_m Hi
```

```text
Hello! How can I assist you today?
```

> [!NOTE]
> Why this model? The goal is to run locally. So we are forced to use small models. The smaller, the more efficient you need to be.
> So instead of a kind of generic model, I picked one specialized in code, able to understand the context, and also good at following instructions and using tools.
> Tools are a critical piece when it's time to instrument a model.
> 
> I pushed several variants of this model on [Docker Hub](https://hub.docker.com/u/eunomie).

### Dagger

Of course, you also need [Dagger](https://dagger.io) to be installed.

You can follow the steps as defined in the [installation guide](https://docs.dagger.io/install) but with a mac, the easiest way is to use brew:

```console
brew install dagger/tap/dagger
```

And check that works

```console
dagger version
```

```text
dagger v0.18.1 (docker-image://registry.dagger.io/engine:v0.18.1) darwin/arm64
```

## How to run

Now everything is setup, you can run the demo.

I provided a `ts-hello` directory to play with. This is a typescript hello world example.

Clone this repository, and enter the directory.

Then ask dagger to run the `dev-environment` and pass the `ts-hello` directory to it. It will create the environment and return a container.
The `| terminal` will then open a terminal in the returned container.

Ensure everything is working, meaning you can run `yarn test` inside the container.

```console
dagger -c 'dev-environment ts-hello | terminal'
```

```text
● Attaching terminal:
    container: Container!
    .from(address: "docker.io/library/alpine:3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c"): Container!
    .withDirectory(
    │ │ directory: Host.directory(path: "/Users/yves/dev/src/github.com/eunomie/local-agent/ts-hello"): Directory!
    │ │ path: "/workspace"
    │ ): Container!
    .withWorkdir(path: "/workspace"): Container!
    .withExec(args: ["apk", "add", "--no-cache", "nodejs", "yarn"], expect: ANY): Container!
    .withExec(args: ["yarn", "install"], expect: ANY): Container!

dagger /workspace $ yarn test
yarn run v1.22.22
$ yarn lint:types && jest --no-cache
$ yarn tsc --noEmit -p .
$ /workspace/node_modules/.bin/tsc --noEmit -p .
(node:40) [DEP0040] DeprecationWarning: The `punycode` module is deprecated. Please use a userland alternative instead.
(Use `node --trace-deprecation ...` to show where the warning was created)
 PASS  ./hello-world.test.ts
  Hello World
    ✓ says hello world with no name (1ms)
    ✓ says hello to bob
    ✓ says hello to sally

Test Suites: 1 passed, 1 total
Tests:       3 passed, 3 total
Snapshots:   0 total
Time:        0.603s
Ran all test suites.
Done in 1.82s.
dagger /workspace $
```

## How does it work?

This small Dagger project is composed of two modules. The main one you can find in the [`.dagger`](https://github.com/eunomie/local-agent/tree/main/.dagger) directory, will basically do the following:

- get the llm object based on the configuration (available in the [`.env`](https://github.com/eunomie/local-agent/blob/main/.env) file)
- pass an environment defining the input (an alpine workspace) and the output (a container with the dev environment) to the model
- ask the model to follow a prompt explaining in deep what we are expecting, basically: analyse the code base, find the tools to install, install them and return the modified container

The "alpine workspace" is an other Dagger module you can find in [`alpine-workspace`](https://github.com/eunomie/local-agent/tree/main/alpine-workspace).
This one is written in Java, but no worries you don't need java on your machine :-)

This workspace is simply:

- an alpine container with the source directory mounted
- a set of commands to install packages, read files, and run commands

This workspace is a very important piece, as by restricting the available commands we help the model to pick the right ones.

And, that's it!

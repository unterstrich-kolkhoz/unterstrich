# unterstrich

unterstrich is the backend repository of the eponymous artist platform. It is
written primarily in Go and JavaScript.

**Disclaimer:** This is a heavy work in progress, and you probably don’t want
to look at it yet. There’s a bunch of icky hardcoded stuff in there and such,
and we want to spare you. That said, if you want to know what this might become,
get in touch!

## Build

To build the project, run `go build` in your `GOPATH`. It should generate a
binary called `unterstrich` that is immediately usable.

You’ll need an installation of SQLite 3 to run it.

The frontend is built on Choo. You can run it by changing into the `frontend`
directory and calling `npm run`.

## Configure

The binary needs a configuration file to work. By default, unterstrich searches for
a file names `./etc/_/server.conf`, but this can be changed by supplying the
`-config` flag.

The configuration can be in one of four directories, namely:

```
./<name>
/<name>
./<name>.local.conf (the first .conf will be replaced)
/<name>.local.conf  (same here)
```

It can also be split up, in which case all matching files will be merged, in
the precedence order above.

## Develop

We use [pre-commit](https://pre-commit.com/) for development. You’ll have to
install it and run `pre-commit install`. This will also require you to install
the following packages:

- `golinter`
- `gometalinter`
- the packages our `gometalinter` configuration depends on
- `prettier`

## Test

TODO!

## Deploy

You’ll have to make sure that the environment variable
`UNTERSTRICH_SECRET_KEY` is set
to a long, randomly generated string. This is used for JWT token generation,
and thus secures our sessions. The secret key shouldn’t be changed between
deployment unless it’s really necessary from a security standpoint, because it
will invalidate all tokens otherwise.

### Docker

To build a Docker image, run the following commands:

```
docker build -t unterstrich .
docker run -it -p 127.0.0.1:8080:8080 --rm --name unterstrich-web unterstrich
```

*Remember to build the frontend before this, because the compiled JS will be
copied verbatim into the container!*

<hr/>

Have fun! :heart:

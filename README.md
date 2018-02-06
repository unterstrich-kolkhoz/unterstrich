# unterstrich

unterstrich is the backend repository of [our product](unterstrich.io). It is written
primarily in Go and JavaScript.

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
will validate all tokens otherwise.

TODO!

<hr/>

Have fun! :heart:

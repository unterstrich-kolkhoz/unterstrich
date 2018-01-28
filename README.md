# artifex

artifex is the backend repository of [our product](artifex.io). It is written
primarily in Go.

## Build

To build the project, run `go build` in your `GOPATH`. It should generate a
binary called `artifex` that is immediately usable.

You’ll need an installation of SQLite 3.

## Configure

The binary needs a configuration file to work. By default, artifex searches for
a file names `./etc/arfx/server.conf`, but this can be changed by supplying the
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

## Test

## Deploy

<hr/>

Have fun! :heart:

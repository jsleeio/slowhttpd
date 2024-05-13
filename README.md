# slowhttpd

## what is this?

A test webserver that responds to requests for `/randomsleep` with a delayed
response.

The extent of the delay is configured with the `-min=DURATION` and
`-max=DURATION` commandline flags.

## demo

### start the server

```
$ ./slowhttpd -min=1s -max=2s -listen=:3000
{"level":"info","msg":"listening on :3000","time":"2024-05-13T12:22:31+10:00"}
{"duration":1039669708,"level":"info","method":"GET","msg":"processing request","time":"2024-05-13T12:24:04+10:00","uri":"/randomsleep"}
{"duration":1195779958,"level":"info","method":"GET","msg":"processing request","time":"2024-05-13T12:24:09+10:00","uri":"/randomsleep"}
```

### make requests

```
$ time curl http://localhost:3000/randomsleep
*snore*

real    0m1.079s
user    0m0.007s
sys 0m0.013s

$ time curl http://localhost:3000/randomsleep
*snore*

real    0m1.221s
user    0m0.009s
sys 0m0.008s
```

### healthchecks

A healthcheck endpoint is provided: `/health`. It does not sleep:

```
$ time curl http://localhost:3000/health
OK

real    0m0.041s
user    0m0.007s
sys 0m0.015s
```

## build it

`slowhttpd` is setup to use Goreleaser and publish packages to Github Packages.
But you can also do a simple `go build`.

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
{"message":"listening on :3000","severity":"info","timestamp":"2024-05-13T14:25:53+10:00"}

[... later, after some requests ...]

{"duration":1.9069897089999999,"message":"processing request","method":"GET","severity":"info","timestamp":"2024-05-13T14:25:58+10:00","uri":"/randomsleep"}
{"duration":2.3649241659999998,"message":"processing request","method":"GET","severity":"info","timestamp":"2024-05-13T14:26:02+10:00","uri":"/randomsleep"}
{"duration":1.427324625,"message":"processing request","method":"GET","severity":"info","timestamp":"2024-05-13T14:26:06+10:00","uri":"/randomsleep"}
{"duration":0.58806725,"message":"processing request","method":"GET","severity":"info","timestamp":"2024-05-13T14:26:08+10:00","uri":"/randomsleep"}
```

### make requests

```
$ curl http://0:3000/randomsleep
*snore*

$ curl http://0:3000/randomsleep
*snore*

$ time curl http://0:3000/randomsleep
*snore*

real    0m1.454s
user    0m0.010s
sys 0m0.011s

$ time curl http://0:3000/randomsleep
*snore*

real    0m0.618s
user    0m0.009s
sys 0m0.010s
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

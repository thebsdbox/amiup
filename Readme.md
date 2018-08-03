
# amiup

## A small binary for running inside a container that can report back

`amiup` works in three modes:

## Building

A linux amiup container can be built using the makefile.

```
make docker 
```

### Listener

When started as a listener then amiup exposes a http endpoint that other amiup containers can connect back to, and report their state.

#### Usage

To start as a listener pass the `-listen` flag upon startup, and by default `amiup` will bind to `:8080`.

To modify the listening port use the `-port` flag with the port number needed.

The server will automatically start a stopwatch for timing incoming events. The stop watch can be reset by hitting the `/stopwatch` endpoint e.g.

```
$ curl 192.168.0.20:8080/stopwatch
```

### Client

The amiup container should be started with `-server` which points to the address where the listener is running.

Along with the `-port` that the listener is using.

Once the `amiup` container has been instantiated by the runtime it will report back to the listener that it’s up and running.


Inspired by [anderspitman/SirTunnel](https://github.com/anderspitman/SirTunnel), built for use with [Traefik](https://github.com/traefik/traefik/).

# Usage

_Values are currently hardcoded in main.go, but can easily be changed to fit your usage._

* Have a copy of tunnel on the server
* Set up traefik to use a redis provider
* Enable GatewayPorts in `sshd_config` for your user, for example:

    ```
    Match User USERNAME
      GatewayPorts clientspecified
    ```

* Run `tunnel -p PORT -s SUBDOMAIN` or alternatively leave out the subdomain to be randomly assigned one.

# Client

The client is only there for some ease of use. It picks a random remote port and subdomain for you,
then opens the ssh connection.

You could also connect by just running:

```
ssh -tR :REMOTE_PORT:localhost:LOCAL_PORT HOST ~/tunnel -server -p REMOTE_PORT -s SUBDOMAIN
```

The subdomain cannot be omitted in this case.

# Server

The server is specified with the `-server` flag. It uses traefik's redis provider to create new
routers and services for the SSH tunnel. Traefik must be configured to have the host accessible
at host.docker.internal, this can be done with

```
    extra_hosts:
      - 'host.docker.internal:host-gateway'
```

in your docker-compose.

# Build

Run `go build` in this directory.

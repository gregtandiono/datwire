# .dat wire

**.dat wire** is a monorepo for the [trade-wire](https://github.com/gregtandiono/trade-wire) project. I'm not happy with how I designed the `trade-wire` project, so this is a full rewrite with microservice design implementation. 

The project is *100%* written in Go.

## How to run it on your machine

Prequisites:

- [Consul](https://www.consul.io/)
- [Golang](https://golang.org/)

### Consul Configuration

You will need to setup consul on your system to be able to run this project, see [docs here](https://www.consul.io/intro/getting-started/install.html).
You will also need to register the services below:

- datwire-gateway
- datwire-auth
- datwire-users

Here's a consul service definition sample:

```json
{
  "service": {
    "name": "datwire-gateway",
    "tags": ["golang", "datwire"],
    "address": "http://127.0.0.1",
    "port": 3001
  }
}
```

Typically, you should have a `/etc/consuld` dir as a config dir. Your config dir should look like this:
```
/etc/consul.d
├── datwire-auth.json
├── datwire-gateway.json
└── datwire-users.json
```

So to start a local dev agent, simply run: 
```bash
consul agent -dev -ui -node=local-dev -config-dir=/etc/consul.d
```
I like using the web UI to monitor my services, so I included the `-ui` flag.

### Consul KV

This project requires a `Consul KV`, specifically to store and fetch hash string for tokenization op.
Run:
```bash
consul kv put datwire/config/hashString [SHA256 KEY HASH]
```
You can generate the hash via this site: [https://www.liavaag.org/English/SHA-Generator/HMAC/](https://www.liavaag.org/English/SHA-Generator/HMAC/)

Then run:
```bash
consul kv get datwire/config/hashString
```
to verify.

*The consul integration is far from finished, there's a lot of work to be done, so stay tuned!*

## Tests

To run existing tests simply run: `go test -v ./...` for now.

## Design

*note: I'll insert some sort of graph here*


## TODOs

I'm currently using *Asana* for my issue tracking and planning. I will eventually move to **github issues** after its initial public release.

- [ ] Service config abstraction
- [ ] API Gateway
- [ ] Integration tests (gateway tests)
- [ ] Test scripts
- [ ] Runner scripts
- [ ] Validation (struct validation)
- [ ] JWT claims should have expiration (1 week?)
- [ ] Auth service should have a refresh token method
- [x] Consul integration for service discovery
- [x] Implement Gokit microservice design principles

...and many more that I'll have to add as I go along!
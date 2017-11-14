# .dat wire

**.dat wire** is a monorepo for the [trade-wire](https://github.com/gregtandiono/trade-wire) project. I'm not happy with how I designed the `trade-wire` project, so this is a full rewrite with microservice design implementation. 

The project is *100%* written in Go.

## Tests

To run existing tests simply run: `go test -v ./...` for now.
I will have to abstract test scripts for different services.

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
- [ ] Consul integration for service discovery
- [ ] Gokit integration for micro-service framework

...and many more that I'll have to add as I go along!

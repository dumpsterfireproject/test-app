# test-app
Sample application to be used to write BDD features against

## Motivation

To provide a single executable that provides a simple application that can be used to practice writing a variety of
end to end tests against. The application includes a terminal based UI, a web UI, a JSON API, and a GraphQL API.

It also has some interesting examples that might make for some good future blog posts, like using gqlgen, air, tview,
hexagonal architecture, and functional options.

## About the Test App

TBD

## Live Reload
A .air.toml file is included in this repository.  See [Air](https://github.com/cosmtrek/air).

I ran `go install github.com/cosmtrek/air@latest` to install Air and then ran it using `../bin/air`.
But feel free to use one of the other install methods documented in the Air git repository.
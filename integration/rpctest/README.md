
[![Build Status](http://img.shields.io/travis/Actinium-project/ltcd.svg)](https://travis-ci.org/Actinium-project/ltcd)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/Actinium-project/acmd/integration/rpctest)

Package rpctest provides a ltcd-specific RPC testing harness crafting and
executing integration tests by driving a `ltcd` instance via the `RPC`
interface. Each instance of an active harness comes equipped with a simple
in-memory HD wallet capable of properly syncing to the generated chain,
creating new addresses, and crafting fully signed transactions paying to an
arbitrary set of outputs.

This package was designed specifically to act as an RPC testing harness for
`ltcd`. However, the constructs presented are general enough to be adapted to
any project wishing to programmatically drive a `ltcd` instance of its
systems/integration tests.

## Installation and Updating

```bash
$ go get -u github.com/Actinium-project/acmd/integration/rpctest
```

## License

Package rpctest is licensed under the [copyfree](http://copyfree.org) ISC
License.

=======
=======
rpctest
=======

[![Build Status](http://img.shields.io/travis/Actinium-project/acmd.svg)](https://travis-ci.org/Actinium-project/acmd)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/Actinium-project/acmd/integration/rpctest)

Package rpctest provides a acmd-specific RPC testing harness crafting and
executing integration tests by driving a `acmd` instance via the `RPC`
interface. Each instance of an active harness comes equipped with a simple
in-memory HD wallet capable of properly syncing to the generated chain,
creating new addresses, and crafting fully signed transactions paying to an
arbitrary set of outputs.

This package was designed specifically to act as an RPC testing harness for
`acmd`. However, the constructs presented are general enough to be adapted to
any project wishing to programmatically drive a `acmd` instance of its
systems/integration tests.

## Installation and Updating

```bash
$ go get -u github.com/Actinium-project/acmd/integration/rpctest
```

## License

Package rpctest is licensed under the [copyfree](http://copyfree.org) ISC
License.

=======
rpctest
=======

[![Build Status](http://img.shields.io/travis/Actinium-project/acmd.svg)]
(https://travis-ci.org/Actinium-project/acmd) [![ISC License]
(http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)]
(http://godoc.org/github.com/Actinium-project/acmd/rpctest)

Package rpctest provides a acmd-specific RPC testing harness crafting and
executing integration tests by driving a `acmd` instance via the `RPC`
interface. Each instance of an active harness comes equipped with a simple
in-memory HD wallet capable of properly syncing to the generated chain,
creating new addresses, and crafting fully signed transactions paying to an
arbitrary set of outputs. 

This package was designed specifically to act as an RPC testing harness for
`acmd`. However, the constructs presented are general enough to be adapted to
any project wishing to programmatically drive a `acmd` instance of its
systems/integration tests. 

## Installation and Updating

```bash
$ go get -u github.com/Actinium-project/acmd/rpctest
```

## License


Package rpctest is licensed under the [copyfree](http://copyfree.org) ISC
License.

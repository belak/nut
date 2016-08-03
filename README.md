# nutdb

A simple, opinionated wrapper built around boltdb to make common operations
easier.

The name is meant to be a simple play on words, because nuts and bolts usually
go well together.

## Goals

This package aims to keep the simplicity of the bolt interface, but use strings
(rather than []byte) for keys and deal with objects by default (rather than
[]byte).

The entire boltdb interface has not been implemented yet, as I don't personally
need all of it. Any unimplemented methods on structs can be found in their
respective go files. Any unimplemented types can be found in
[unimplemented.go](unimplemented.go) If you feel it is missing something, please
feel free to open an issue or submit a pull request.

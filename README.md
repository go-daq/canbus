# canbus

[![Build Status](https://secure.travis-ci.org/go-daq/canbus.png)](http://travis-ci.org/go-daq/canbus) [![GoDoc](https://godoc.org/github.com/go-daq/canbus?status.svg)](https://godoc.org/github.com/go-daq/canbus)

`canbus` provides high-level facilities to interact with `CAN` sockets.

## can-dump

```
can-dump prints data flowing on CAN bus.

Usage of can-dump:

sh> can-dump [options] <CAN interface>
   (use CTRL-C to terminate can-dump)

Examples:

 can-dump vcan0
```

```sh
$> can-dump vcan0
  vcan0  255 00000000  11 22 33 44 55 66 77 88   |."3DUFW.|
  vcan0  256 00000000  11 22 33 44 55 66 77 88   |."3DUFW.|
  vcan0  512 00000000  11 22 33 44 55 66 77 88   |."3DUFW.|
  vcan0  612 00000000  11 22 33 44 55 66 77 88   |."3DUFW.|
  vcan0  712 00000000  11 22 33 44 55 66 77 88   |."3DUFW.|
  vcan0  999 00000000  11 22 33 44 55 66 77 88   |."3DUFW.|
  vcan0  999 00000000  11 22 33 44 55 66 77 88   |."3DUFW.|
  vcan0  080 00000000  00 DE AD BE EF            |.....|
  vcan0  080 00000000  01 DE AD BE EF            |.....|
  vcan0  080 00000000  02 DE AD BE EF            |.....|
[...]
```

# canbus

[![Build Status](https://github.com/go-daq/canbus/workflows/CI/badge.svg)](https://github.com/go-daq/canbus/actions)
[![codecov](https://codecov.io/gh/go-daq/canbus/branch/main/graph/badge.svg)](https://codecov.io/gh/go-daq/canbus)
[![GoDoc](https://godoc.org/github.com/go-daq/canbus?status.svg)](https://godoc.org/github.com/go-daq/canbus)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](https://github.com/go-daq/canbus/blob/main/LICENSE)

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
  vcan0  080 00000000  00 DE AD BE EF            |.....|
  vcan0  080 00000000  01 DE AD BE EF            |.....|
  vcan0  080 00000000  02 DE AD BE EF            |.....|
  vcan0  080 00000000  03 DE AD BE EF            |.....|
  vcan0  080 00000000  04 DE AD BE EF            |.....|
  vcan0  080 00000000  05 DE AD BE EF            |.....|
  vcan0  080 00000000  06 DE AD BE EF            |.....|
  vcan0  080 00000000  07 DE AD BE EF            |.....|
  vcan0  080 00000000  08 DE AD BE EF            |.....|
  vcan0  080 00000000  09 DE AD BE EF            |.....|
  vcan0  712 00000000  11 22 33 44 55 66 77 88   |."3DUFW.|
  vcan0  7fa 00000000  DE AD BE EF               |....|
[...]
```

## can-send

```
can-send sends data on the CAN bus.

Usage of can-send:

sh> can-send [options] <CAN interface> <CAN frame>

where <CAN frame> is of the form: <ID-hex>#<frame data-hex>.

Examples:

 can-send vcan0 f12#1122334455667788
 can-send vcan0 ffa#deadbeef
```

```sh
$> can-dump vcan0 &
$> can-send vcan0 f12#1122334455667788
  vcan0  712 00000000  11 22 33 44 55 66 77 88   |."3DUFW.|
$> can-send vcan0 ffa#deadbeef
  vcan0  7fa 00000000  DE AD BE EF               |....|
```

## References

```sh
$> modprobe can
$> modprobe can_raw
$> modprobe vcan

## setup vcan network devices
$> ip link add type vcan
$> ip link add dev vcan0 type vcan
$> ip link set vcan0 up
```

- https://www.kernel.org/doc/Documentation/networking/can.txt
- https://en.wikipedia.org/wiki/CAN_bus
- https://en.wikipedia.org/wiki/SocketCAN
- https://godoc.org/golang.org/x/sys/unix#SocketaddrCAN
- http://www.computer-solutions.co.uk/info/Embedded_tutorials/can_tutorial.htm
- http://elinux.org/Bringing_CAN_interface_up
- http://inst.cs.berkeley.edu/~ee249/fa08/Lectures/handout_canbus2.pdf
- https://chemnitzer.linux-tage.de/2012/vortraege/folien/1044_SocketCAN.pdf

# modular-netstalking

## Install

```sh
make install
```

## Usage

HTTP servers

```sh
wan-ips | check-port
```

Quotd

```sh
wan-ips | check-port -w 1024 -p 17 | xargs -I@ -P8 timeout 5 ncat @ 17
```
HTML titles

```sh
wan-ips | check-port | xargs -I@ -P8 curl -s @ | grep -ioP '(?<=<title>)[^<]+'
```

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

FTP listings

```sh
wan-ips | check-port -w 1024 -p 21 | xargs -I@ -P8 curl 'ftp://@'
```

IPs w/open WP uploads dir

```sh
wan-ips | check-port -w 1024 | xargs -I@ -P8 bash -c \
  'timeout 5 curl -s "http://@/wp-content/uploads/" | grep -qF "Index of" && echo @'
```

Stats of open MySQL ports per 100k hosts

```sh
wan-ips -c 100000 | check-port -w 1024 -p 3306 | wc -l
# 171
```

# co2mini_exporter

A Prometheus exporter for CO2-MINI.

## requirements

* Linux
* go compiler

## build

```bash
$ go build
```

## run

```bash
$ ./co2mini_exporter &
```

## get metrics

```bash
$ curl http://localhost:9002/metrics
```

## stop

```bash
$ pkill co2mini_exporter
```

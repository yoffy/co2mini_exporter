# co2mini_exporter

A Prometheus exporter for CO2-MINI.

## requirements

* Linux
* go compiler

## install

```bash
$ make && sudo make install
```

It is installed as co2mini_exporter.service.

## uninstall

```bash
$ sudo make uninstall
```

## get metrics

```bash
$ curl http://localhost:9002/metrics
```

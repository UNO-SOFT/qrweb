module github.com/UNO-SOFT/qrweb

go 1.27.0

require (
	github.com/UNO-SOFT/cli v0.0.0-20261009074406-95b18dbd4953
	github.com/UNO-SOFT/zlog v0.8.6
	github.com/aohorodnyk/mimeheader v0.0.6
	github.com/peterbourgon/ff/v3 v3.4.0
	github.com/tgulacsi/go v0.30.0
	golang.org/x/text v0.43.0
	rsc.io/qr v0.2.0
)

require (
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/zerologr v1.2.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	golang.org/x/exp v0.0.0-20261009195045-ca0d7ba23607 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/term v0.47.0 // indirect
)

replace rsc.io/qr v0.2.0 => github.com/tgulacsi/qr v0.0.0-20221008053105-60638b3fcf7e

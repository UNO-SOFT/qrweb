module github.com/UNO-SOFT/qrweb

go 1.21

toolchain go1.21.0

require (
	github.com/UNO-SOFT/zlog v0.7.7
	github.com/aohorodnyk/mimeheader v0.0.6
	github.com/peterbourgon/ff/v3 v3.3.1
	github.com/tgulacsi/go v0.24.4
	golang.org/x/text v0.9.0
	rsc.io/qr v0.2.0
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/zerologr v1.2.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/rs/zerolog v1.29.0 // indirect
	golang.org/x/exp v0.0.0-20230713183714-613f0c0eb8a1 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/term v0.10.0 // indirect
)

replace rsc.io/qr v0.2.0 => github.com/tgulacsi/qr v0.0.0-20221008053105-60638b3fcf7e

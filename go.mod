module github.com/UNO-SOFT/qrweb

go 1.19

require (
	github.com/UNO-SOFT/zlog v0.0.1
	github.com/aohorodnyk/mimeheader v0.0.6
	github.com/peterbourgon/ff/v3 v3.3.0
	github.com/tgulacsi/go v0.22.2
	golang.org/x/text v0.3.7
	rsc.io/qr v0.2.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/zerologr v1.2.2 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/rs/zerolog v1.28.0 // indirect
	golang.org/x/sys v0.0.0-20220702020025-31831981b65f // indirect
	golang.org/x/term v0.0.0-20220919170432-7a66f970e087 // indirect
)

replace rsc.io/qr v0.2.0 => github.com/tgulacsi/qr v0.0.0-20221008053105-60638b3fcf7e

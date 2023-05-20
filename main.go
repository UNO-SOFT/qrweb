// Copyright 2022, 2023 Tamás Gulácsi. All rights reserved.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"

	"github.com/UNO-SOFT/zlog/v2"
	"github.com/aohorodnyk/mimeheader"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/tgulacsi/go/httpunix"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"rsc.io/qr"
)

var verbose zlog.VerboseVar
var logger = zlog.NewLogger(zlog.MaybeConsoleHandler(&verbose, os.Stderr)).SLog()

func main() {
	if err := Main(); err != nil {
		logger.Error("Main", "error", err)
	}
}

// Main function
func Main() error {
	fs := flag.NewFlagSet("qrweb", flag.ContinueOnError)
	fs.Var(&verbose, "v", "verbose logging")
	flagAddr := fs.String("addr", ":3456", "address to listen on. May be unix://")
	app := ffcli.Command{Name: "qrweb", FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			logger.Info("start listening", "addr", *flagAddr)
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				var body []byte
				if err := func() error {
					var err error
					r.Body = struct {
						io.Reader
						io.Closer
					}{io.LimitReader(r.Body, 1<<20), r.Body}
					if body, err = httputil.DumpRequest(r, true); err != nil {
						return err
					}
					logger.Info("got", "request", maskSecret{text: string(body)})

					text := maskSecret{text: r.FormValue("text")}
					logger.Info("text1", "text", text)
					if text.text == "" {
						if text.text = r.FormValue("data"); text.text == "" {
							text.text = strings.TrimPrefix(r.URL.Path, "/")
						}
					}

					charset := r.FormValue("charset")
					var enc encoding.Encoding
					if charset != "" && strings.EqualFold(strings.ReplaceAll(charset, "-", ""), "iso88592") {
						enc = charmap.ISO8859_2
						if t, err := enc.NewDecoder().String(text.text); err != nil {
							return fmt.Errorf("decode %q as %q: %w", text, charset, err)
						} else {
							text.text = t
						}
					}

					level := qr.H
					switch r.FormValue("level") {
					case "L", "l", "20", "20%":
						level = qr.L
					case "M", "m", "38", "38%":
						level = qr.M
					case "Q", "q", "55", "55%":
						level = qr.Q
					}

					acc := r.Header.Get("Accept")
					if acc == "" {
						acc = r.FormValue("accept")
					}
					if acc != "" && strings.IndexByte(acc, '/') < 0 {
						acc = "image/" + acc
					}
					ah := mimeheader.ParseAcceptHeader(acc)
					logger.Info("encode", "text", text, "level", level, "accept", ah)
					code, err := qr.Encode(text.text, level)
					if err != nil {
						return err
					}
					mt := "image/png"
					_, mt, _ = ah.Negotiate([]string{"image/png", "image/gif"}, mt)
					w.Header().Set("Content-Type", mt)
					switch mt {
					case "image/gif":
						err = gif.Encode(w, code.Image(), &gif.Options{
							NumColors: 2,
							Quantizer: bwQuantizer{},
						})
					default:
						_, err = w.Write(code.PNG())
					}
					return err
				}(); err != nil {
					logger.Error("handle", "request", string(body), "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			})
			return httpunix.ListenAndServe(ctx, *flagAddr, nil)
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	return app.ParseAndRun(ctx, os.Args[1:])
}

type bwQuantizer struct{}

func (bwQuantizer) Quantize(p color.Palette, m image.Image) color.Palette {
	if len(p) >= 2 {
		return p
	}
	return append(p[:0], color.Black, color.White)
}

type maskSecret struct {
	text string
}

func (s maskSecret) String() string {
	const prefix = "secret="
	if i := strings.Index(s.text, prefix); i >= 0 {
		if j := strings.IndexAny(s.text[i+len(prefix):], "& \r\n\t"); j >= 0 {
			return s.text[:i+len(prefix)] + "***" + s.text[i+len(prefix)+j:]
		}
		return s.text[:i+len(prefix)] + "***"
	}
	return s.text
}

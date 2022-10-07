// @Title
// @Description
// @Author
// @Update

package main

import (
	"context"
	"flag"
	"image"
	"image/color"
	"image/gif"
	"net/http"
	"os"
	"os/signal"

	"github.com/UNO-SOFT/zlog"
	"github.com/aohorodnyk/mimeheader"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/tgulacsi/go/httpunix"
	"rsc.io/qr"
)

var logger = zlog.New(zlog.MaybeConsoleWriter(os.Stderr))

func main() {
	if err := Main(); err != nil {
		logger.Error(err, "Main")
	}
}

func Main() error {
	fs := flag.NewFlagSet("qrweb", flag.ContinueOnError)
	flagAddr := fs.String("addr", ":3456", "address to listen on. May be unix://")
	app := ffcli.Command{Name: "qrweb", FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			logger.Info("start listening", "addr", *flagAddr)
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if err := func() error {
					level := qr.H

					text := r.FormValue("text")
					if text == "" {
						if text = r.FormValue("data"); text == "" {
							if text = r.URL.Path; text != "" && text[0] == '/' {
								text = text[1:]
							}
						}
					}
					switch r.FormValue("level") {
					case "L", "l", "20", "20%":
						level = qr.L
					case "M", "m", "38", "38%":
						level = qr.M
					case "Q", "q", "55", "55%":
						level = qr.Q
					}

					ah := mimeheader.ParseAcceptHeader(r.Header.Get("Accept"))
					logger.Info("encode", "text", text, "level", level, "accept", ah)
					code, err := qr.Encode(text, level)
					if err != nil {
						return err
					}
					mt := "image/png"
					if acc := r.Header.Get("Accept"); acc != "" {
						_, mt, _ = mimeheader.ParseAcceptHeader(acc).Negotiate([]string{"image/png", "image/gif"}, mt)
					}
					w.Header().Set("Content-Type", mt)
					switch mt {
					case "image/gif":
						err = gif.Encode(w, scaledImage{code.Image()}, &gif.Options{
							NumColors: 2,
							Quantizer: bwQuantizer{},
						})
					default:
						_, err = w.Write(code.PNG())
					}
					return err
				}(); err != nil {
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

type scaledImage struct {
	image.Image
}

func (m scaledImage) At(x, y int) color.Color { return m.Image.At(x>>3-4, y>>3-4) }

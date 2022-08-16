package http

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/juanjoss/qrgen/qr"
)

func GenerateQR(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to parse form: %v", err.Error())
		return
	}

	source := r.FormValue("source")

	qr, err := qr.Generate(source)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to crate QR: %v", err)
		return
	}

	png.Encode(w, *qr)
}

package qr

import (
	barcode "github.com/boombuler/barcode"
	barcode_qr "github.com/boombuler/barcode/qr"
)

type qr *barcode.Barcode

func Generate(source string) (qr, error) {
	encodedData, err := barcode_qr.Encode(source, barcode_qr.L, barcode_qr.Auto)
	if err != nil {
		return nil, err
	}

	barcode, err := barcode.Scale(encodedData, 256, 256)
	if err != nil {
		return nil, err
	}

	return &barcode, nil
}

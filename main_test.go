package main

import (
	"testing"

	"github.com/rombintu/svg-driver/core"
)

func TestSvg2Path(t *testing.T) {
	if err := core.ConvertSvg2Png("templates/img2.svg", "dist/img2.svg.png"); err != nil {
		t.Fatal(err)
	}
}

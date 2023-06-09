package core

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

const (
	DistDirName      string = "dist"
	TemplatesDirName string = "templates"
)

type SVG struct {
	XMLName xml.Name    `xml:"svg"`
	File    fs.FileInfo `xml:"-"`
	Path    Path        `xml:"path"`
	ViewBox string      `xml:"viewBox,attr"`
	Width   string      `xml:"width,attr"`
	Height  string      `xml:"height,attr"`
}

type Path struct {
	XMLName xml.Name `xml:"path"`
	D       string   `xml:"d,attr"`
	Fill    string   `xml:"fill,attr"`
	Style   string   `xml:"style,attr"`
}

type CSS map[string]string

func (css CSS) unmarshal() string {
	b := new(bytes.Buffer)
	for key, value := range css {
		fmt.Fprintf(b, "%s:%s;", key, value)
	}
	return b.String()
}

func (css CSS) ToStyle() string {
	return css.unmarshal()
}

func SVGFromFile(filepath string) (svg *SVG, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	err = xml.Unmarshal(content, &svg)
	if err != nil {
		return
	}
	svg.File, err = f.Stat()
	return
}

func (s *SVG) SaveToSvg(filepath string) error {
	content, err := xml.Marshal(s)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *SVG) SvgToPng(distDir string) error {
	src := path.Join(distDir, s.File.Name())
	dst := fmt.Sprintf("%s.png", path.Join(distDir, s.File.Name()))
	if err := ConvertSvg2Png(src, dst); err != nil {
		return err
	}
	return nil
}

func (s *SVG) SetStyle(css CSS) {
	s.Path.Style = css.ToStyle()
}

func (s *SVG) SetFill(color string) {
	s.Path.Fill = color
}

func (s *SVG) SetSize(w, h int, vbox string) {
	s.Width = fmt.Sprintf("%dpx", w)
	s.Height = fmt.Sprintf("%dpx", h)
	if vbox != "" {
		s.ViewBox = vbox
	}
}

func ConvertSvg2Png(src, dst string) error {
	cmd := exec.Command("convert", src, dst)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

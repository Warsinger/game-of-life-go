// Command icons generates platform icons from the application artwork.
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	data, err := os.ReadFile("assets/icons/source.json")
	check(err)
	var source struct {
		Pattern          []string
		Live, Background string
	}
	check(json.Unmarshal(data, &source))
	if len(source.Pattern) != 5 || source.Live != "#00ff00" || source.Background != "#000000" {
		log.Fatal("expected a five-cell green-on-black pattern")
	}
	for _, row := range source.Pattern {
		if len(row) != 5 {
			log.Fatal("pattern rows must be five cells")
		}
		for _, cell := range row {
			if cell != '.' && cell != 'O' {
				log.Fatal("unknown pattern cell")
			}
		}
	}
	icons := map[int][]byte{}
	for _, size := range []int{16, 32, 48, 64, 128, 180, 192, 256, 512, 1024} {
		canvas := image.NewRGBA(image.Rect(0, 0, size, size))
		draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.RGBA{0, 0, 0, 255}), image.Point{}, draw.Src)
		for y, row := range source.Pattern {
			for x, cell := range row {
				if cell == 'O' {
					draw.Draw(canvas, image.Rect(x*size/5, y*size/5, (x+1)*size/5, (y+1)*size/5), image.NewUniform(color.RGBA{0, 255, 0, 255}), image.Point{}, draw.Src)
				}
			}
		}
		var data bytes.Buffer
		check(png.Encode(&data, canvas))
		icons[size] = data.Bytes()
		check(os.WriteFile(filepath.Join("assets/icons", "icon-"+strconv.Itoa(size)+".png"), data.Bytes(), 0644))
	}
	// ICO uses PNG payloads, supported by current Windows and browsers.
	var ico bytes.Buffer
	for _, value := range []uint16{0, 1, 6} {
		check(binary.Write(&ico, binary.LittleEndian, value))
	}
	offset := uint32(6 + 16*6)
	for _, size := range []int{16, 32, 48, 64, 128, 256} {
		dimension := byte(size) // 0 encodes 256.
		ico.Write([]byte{dimension, dimension, 0, 0})
		for _, value := range []uint16{1, 32} {
			check(binary.Write(&ico, binary.LittleEndian, value))
		}
		check(binary.Write(&ico, binary.LittleEndian, uint32(len(icons[size]))))
		check(binary.Write(&ico, binary.LittleEndian, offset))
		offset += uint32(len(icons[size]))
	}
	for _, size := range []int{16, 32, 48, 64, 128, 256} {
		ico.Write(icons[size])
	}
	check(os.WriteFile("assets/icons/favicon.ico", ico.Bytes(), 0644))
	// Modern ICNS chunks contain PNG payloads, including Retina resolutions.
	var chunks bytes.Buffer
	for _, entry := range []struct {
		kind string
		size int
	}{{"icp4", 16}, {"icp5", 32}, {"icp6", 64}, {"ic07", 128}, {"ic08", 256}, {"ic09", 512}, {"ic10", 1024}} {
		chunks.WriteString(entry.kind)
		check(binary.Write(&chunks, binary.BigEndian, uint32(8+len(icons[entry.size]))))
		chunks.Write(icons[entry.size])
	}
	var icns bytes.Buffer
	icns.WriteString("icns")
	check(binary.Write(&icns, binary.BigEndian, uint32(8+chunks.Len())))
	icns.Write(chunks.Bytes())
	check(os.WriteFile("assets/icons/app.icns", icns.Bytes(), 0644))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

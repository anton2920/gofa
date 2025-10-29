package bplus

import (
	"fmt"

	"github.com/anton2920/gofa/trace"
)

type Pager interface {
	ReadPagesAt(pages []Page, index int64) (int64, error)
	WritePagesAt(pages []Page, index int64) (int64, error)
}

type MemoryPager struct {
	Pages []Page
}

var _ Pager = new(MemoryPager)

func (p *MemoryPager) ReadPagesAt(pages []Page, index int64) (int64, error) {
	defer trace.End(trace.Begin(""))

	if index < 0 {
		index = int64(len(p.Pages))
	}

	if (index < 0) || (index >= int64(len(p.Pages))) {
		return index, fmt.Errorf("pages index out of bounds")
	}

	copy(pages, p.Pages[index:])
	return index, nil
}

func (p *MemoryPager) WritePagesAt(pages []Page, index int64) (int64, error) {
	defer trace.End(trace.Begin(""))

	if index < 0 {
		index = int64(len(p.Pages))
	}

	if (index < 0) || (index >= int64(len(p.Pages))+1) {
		return -1, fmt.Errorf("pages index out of bounds")
	}

	if index == int64(len(p.Pages)) {
		p.Pages = append(p.Pages, pages...)
	} else {
		copy(p.Pages[index:], pages)
	}

	return index, nil
}

/*
type FilePager struct {
	File *os.File
}

var _ Pager = new(FilePager)

func FilePagerNew(path string) (*FilePager, error) {
	var err error

	p := new(FilePager)
	p.File, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open/create file for pager: %v", err)
	}

	return p, nil
}

func (p *FilePager) Close() {
	p.File.Close()
}

func (p *FilePager) ReadPagesAt(pages []Page, offset int64) error {
	defer trace.End(trace.Begin(""))

	if _, err := p.File.ReadAt(Pages2Bytes(pages), offset); (err != nil) && (err != io.EOF) {
		return fmt.Errorf("failed to read %d pages at %d: %v", len(pages), offset, err)
	}
	return nil
}

func (p *FilePager) WritePagesAt(pages []Page, offset int64) error {
	defer trace.End(trace.Begin(""))

	if _, err := p.File.WriteAt(Pages2Bytes(pages), offset); err != nil {
		return fmt.Errorf("failed to write %d pages at %d: %v", len(pages), offset, err)
	}
	if false {
		defer trace.End(trace.Begin("main.WritePagesAt.Sync"))

		if err := p.File.Sync(); err != nil {
			return fmt.Errorf("failed to sync writes to disk: %v", err)
		}
	}
	return nil
}
*/

package main

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"testing"
)

func TestNewProgress(t *testing.T) {
	progress := NewProgress(3, 20000)
	expectedProgress := &Progress{
		length: 60000,
		pieces: []*Piece{
			&Piece{
				length: 20000,
				blocks: []*Block{
					&Block{
						index:  0,
						begin:  0,
						length: 16384,
						state:  BlockStateNeed,
					},
					&Block{
						index:  0,
						begin:  16384,
						length: 3616,
						state:  BlockStateNeed,
					},
				},
			},
			&Piece{
				length: 20000,
				blocks: []*Block{
					&Block{
						index:  1,
						begin:  0,
						length: 16384,
						state:  BlockStateNeed,
					},
					&Block{
						index:  1,
						begin:  16384,
						length: 3616,
						state:  BlockStateNeed,
					},
				},
			},
			&Piece{
				length: 20000,
				blocks: []*Block{
					&Block{
						index:  2,
						begin:  0,
						length: 16384,
						state:  BlockStateNeed,
					},
					&Block{
						index:  2,
						begin:  16384,
						length: 3616,
						state:  BlockStateNeed,
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(progress, expectedProgress) {
		t.Errorf("Expected:\n%v,\nGot:\n%v\n", spew.Sdump(expectedProgress), spew.Sdump(progress))
	}
}

func TestPercentComplete(t *testing.T) {
	p := NewProgress(3, 20000)
	p.pieces[0].blocks[0].state = BlockStateHave
	p.pieces[0].blocks[1].state = BlockStateHave

	percentComplete := p.percentComplete()
	expectedPercentComplete := 20000.0 / 60000.0

	if percentComplete != expectedPercentComplete {
		t.Errorf("Expected:%v,Got:%v\n", expectedPercentComplete, percentComplete)
	}
}

func TestRandomNeededBlock(t *testing.T) {
	p := NewProgress(2, 20000)
	p.pieces[0].blocks[0].state = BlockStateHave
	p.pieces[0].blocks[1].state = BlockStateHave
	p.pieces[1].blocks[0].state = BlockStateHave

	block, e := p.randomNeededBlock()
	expectedBlock := &Block{
		index:  1,
		begin:  16384,
		length: 3616,
		state:  BlockStateDownloading,
	}

	if e != nil {
		t.Errorf("%v", e)
	}

	if !reflect.DeepEqual(block, expectedBlock) {
		t.Errorf("Expected %v, got %v", expectedBlock, block)
	}
}

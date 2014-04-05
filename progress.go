package main

import (
  "math/rand"
  "math"
  "errors"
  "fmt"
)

var BlockLength = 16384

const (
  BlockStateNeed = iota
  BlockStateHave
  BlockStateDownloading
)

const (
  PieceStateNeed = iota
  PieceStateHave
  PieceStateHavePartial
)

type Progress struct {
  length int
  pieces []*Piece
}

type Piece struct {
  length int
  blocks []*Block
}

type Block struct {
  index  int // piece index
  begin  int
  length int
  state  int
}

func NewProgress(pieceCount int, pieceLength int) *Progress {
  var pieces []*Piece
  length := pieceCount * pieceLength

  for i := 0; i < pieceCount; i++ {
    blockCount := int(math.Ceil(float64(pieceLength) / float64(BlockLength)))
    var blocks []*Block
    for x := 0; x < blockCount; x++ {
      var length int
      begin := x * BlockLength
      if begin + BlockLength > pieceLength {
        length = pieceLength - begin
      } else {
        length = BlockLength
      }

      block := &Block {
        index:  i,
        begin:  begin,
        length: length,
        state:  BlockStateNeed,
      }
      blocks = append(blocks, block)
    }
    piece := &Piece{
      length: pieceLength,
      blocks: blocks,
    }
    pieces = append(pieces, piece)
  }

  return &Progress {
    length: length,
    pieces: pieces,
  }
}

func (p *Progress) percentComplete() float64 {
  downloadedLength := p.downloadedLength()
  return (float64(downloadedLength) / float64(p.length))
}

func (p *Progress) randomNeededBlock() (block *Block, err error) {
  neededBlocks := p.blocksWeNeed()
  if len(neededBlocks) == 0 {
    err = errors.New("No blocks found")
    return
  }
  randomIndex := rand.Intn(len(neededBlocks))
  block = neededBlocks[randomIndex]
  block.state = BlockStateDownloading
  return
}

func (p *Progress) blocksWeHave() (blocks []*Block) {
  for _,piece := range p.pieces {
    for _,block := range piece.blocks {
      if block.state == BlockStateHave {
        blocks = append(blocks, block)
      }
    }
  }
  return
}

func (p *Progress) blocksWeNeed() (blocks []*Block) {
  for _,piece := range p.pieces {
    for _,block := range piece.blocks {
      if block.state == BlockStateNeed {
        blocks = append(blocks, block)
      }
    }
  }
  return
}

func (p *Progress) findBlockFor(index int, begin int) (theBlock *Block, err error) {
  for _,piece := range p.pieces {
    for _,block := range piece.blocks {
      if block.index == index &&
        block.begin == begin {
          theBlock = block
          return
      }
    }
  }
  err = fmt.Errorf("Couldn't find block {index:%v, length:%v}", index, begin)
  return
}

func (p *Progress) downloadedLength() int {
  var length int
  haveBlocks := p.blocksWeHave()
  for _,block := range haveBlocks {
    length += block.length
  }
  return length
}

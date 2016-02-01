package jigsaw

import (
	"errors"
	"fmt"
	"image"
	"math"

	"image/color"
	"encoding/json"
	"io/ioutil"
)

type Jigsaw struct {
	Rows   int
	Pieces []*Piece
	Path   string
	Bounds image.Rectangle
}

type pieceJoint struct {
	External bool
	Side int
}


type Piece struct {
	Height			int
	Width			int
	Joints           []pieceJoint
	Points           []image.Point
	IsCorner         bool
	IsEdge           bool
	IsCenter		 bool
	Name             string
	Path             string
	RightPieceIndex  int
	LeftPieceIndex   int
	BottomPieceIndex int
	TopPieceIndex    int
	Index            int
	Row              int
	Board            image.Rectangle
	Bounds			 image.Rectangle
	Image 			 image.Image
}

func (p * Piece)TopRow()bool {
	for _,p := range p.Points{
		if p.Y == 0{
			return true;
		}
	}
	return false
}
func (p * Piece)BottomRow()bool  {
	for _,pt := range p.Points{
		if pt.Y == p.Board.Max.Y{
			return true;
		}
	}
	return false
}

func (p *Piece)FarRightVerticalRow()bool  {
	for _,pt := range p.Points{
		if pt.X == p.Board.Max.X{
			return true;
		}
	}
	return false
}
func (p *Piece)FarLeftVerticalRow()bool  {
	for _,pt := range p.Points{
		if pt.X == p.Board.Min.X{
			return true;
		}
	}
	return false
}

type circle struct {
	p     image.Point
	r     int
	image image.Image
	inverse  bool
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return c.image.Bounds()
}

//TODO switch transparent for when the circle is external joint
func (c *circle) At(x, y int) color.Color {
	innerBounds := image.Rect(c.image.Bounds().Min.X,c.p.Y,c.image.Bounds().Max.X,c.image.Bounds().Max.Y)
	//c.p.X is the x point for the center of the circle
	//c.p.Y is the y point for the center of the circle
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr { //the point of ref by x,y is within the diameter of the circle
		if c.inverse{
			return c.image.At(x, y)
		}else {
			return color.Transparent
		}
	}
	if c.inverse{
		if y < innerBounds.Min.Y{

			return c.image.At(x,y)
		}else{
			return color.Transparent
		}

	}else{
		return c.image.At(x, y)
	}

}

type JigsawBuilder struct {
	PieceCutter		PieceCutter
	numPieces       int
	numRows         int
	numPiecesPerRow int
	baseImage       image.Image
}

//todo break up
func (jb *JigsawBuilder) buildPieces() ([]*Piece, error) {

	piecesPerLine := jb.numPieces / jb.numRows
	jb.numPiecesPerRow = piecesPerLine
	rect := jb.baseImage.Bounds()
	pieceHeight := rect.Max.Y / jb.numRows
	pieceWidth := rect.Max.X / piecesPerLine
	pieces := make([]*Piece, 0)

	isCorner := func(points []image.Point) bool {
		var corner = false
		for _, p := range points {
			if ((p.X / pieceWidth) == jb.numPiecesPerRow || p.X == 0) && (p.Y / pieceHeight == jb.numRows || p.Y == 0) {
				corner = true
			}
		}
		return corner
	}

	isEdge := func(points []image.Point) bool {
		if isCorner(points) {
			return true
		}
		for _, p := range points {
			if p.Y == 0 || p.Y / pieceHeight == jb.numRows {
				return true
			} else if p.X == 0 || p.X / pieceWidth == jb.numPiecesPerRow {
				return true
			}
		}
		return false
	}

	isCenter := func(points []image.Point) bool{
		return (! isCorner(points) && ! isEdge(points))
	}

//	setNeighbours := func(piece *Piece) {
//		if piece.IsCorner {
//			if piece.Index%jb.numPiecesPerRow == 0 && piece.Index == jb.numPieces {
//				//bottom right
//				piece.LeftPieceIndex = piece.Index - 1
//				piece.TopPieceIndex = piece.Index - jb.numPiecesPerRow
//			} else if piece.Index == jb.numPiecesPerRow {
//				//top right
//				piece.LeftPieceIndex = piece.Index - 1
//				piece.BottomPieceIndex = piece.Index + jb.numPiecesPerRow
//			} else if piece.Index == 1 {
//				//top left
//				piece.RightPieceIndex = piece.Index + 1
//				piece.BottomPieceIndex = piece.Index + jb.numPiecesPerRow
//			} else {
//				//bottom left
//				piece.RightPieceIndex = piece.Index + 1
//				piece.TopPieceIndex = piece.Index - jb.numPiecesPerRow
//			}
//		} else if piece.IsEdge {
//			if piece.Index%jb.numPiecesPerRow == 0 {
//				//right side
//				piece.LeftPieceIndex = piece.Index - 1
//				piece.TopPieceIndex = piece.Index - jb.numPiecesPerRow
//				piece.BottomPieceIndex = piece.Index + jb.numPiecesPerRow
//			} else if piece.Row == 0 {
//				//top
//				piece.RightPieceIndex = piece.Index + 1
//				piece.LeftPieceIndex = piece.Index - 1
//				piece.BottomPieceIndex = piece.Index + jb.numPiecesPerRow
//			} else if piece.Row == jb.numRows {
//				//bottom
//				piece.RightPieceIndex = piece.Index + 1
//				piece.LeftPieceIndex = piece.Index - 1
//				piece.TopPieceIndex = piece.Index - jb.numPiecesPerRow
//			} else {
//				//left
//				piece.RightPieceIndex = piece.Index + 1
//				piece.TopPieceIndex = piece.Index - jb.numPiecesPerRow
//				piece.BottomPieceIndex = piece.Index + jb.numPiecesPerRow
//			}
//		} else {
//			//middle piece
//			piece.RightPieceIndex = piece.Index +1
//			piece.LeftPieceIndex  = piece.Index -1
//			piece.TopPieceIndex = piece.Index - jb.numPiecesPerRow
//			piece.BottomPieceIndex = piece.Index + jb.numPiecesPerRow
//		}
//
//	}
	//use type
	X := func (x, pieceWidth int) (int,int){
		rightX := x + pieceWidth
		return x,rightX
	}

	Y := func (y, pieceHeight int)(int,int){
		topY := y + pieceHeight

		return y,topY
	}

	pieceNum := 1
	//todo make concurrent
	for i := 0; i < jb.numRows; i++ {
		currentPosY := i * pieceHeight
		bottomY,topY:= Y(currentPosY,pieceHeight)
		for j := 0; j < piecesPerLine; j++ {
			currentPosX := j * pieceWidth
			leftX,rightX := X(currentPosX,pieceWidth)
			points := []image.Point{image.Pt(leftX, bottomY), image.Pt(rightX, bottomY), image.Pt(leftX, topY), image.Pt(rightX, topY)}
			corner := isCorner(points)
			fmt.Println("** is Corner ", corner, piecesPerLine, pieceWidth)
			p := &Piece{Height: pieceHeight, Width: pieceWidth,Points: points, IsCorner: corner, IsEdge: isEdge(points), Name: fmt.Sprintf("piece%d", pieceNum), Index: pieceNum, Row: i,Board:jb.baseImage.Bounds(), IsCenter:isCenter(points),  Joints:[]pieceJoint{{External:false,Side:1},{External:true,Side:2}}}
			pieces = append(pieces, p)
			pieceNum++
		}
	}
	data,_ := json.MarshalIndent(pieces,"", " ")
	ioutil.WriteFile("./out/out.json",data,777)
	return pieces, nil

}

func (jb *JigsawBuilder) buildRows() error {
	//we are only interested in whole numbers
	sqr := int(math.Sqrt(float64(jb.numPieces)))
	if jb.numPieces%sqr != 0 {
		return errors.New("num pieces should have whole number square root. ie 4, 8, 12, 15")
	}
	jb.numRows = sqr
	return nil
}

func (jb *JigsawBuilder) Build() (Jigsaw, error) {
	jig := Jigsaw{}
	jig.Bounds = jb.baseImage.Bounds()
	err := jb.buildRows()
	if err != nil {
		return jig, err
	}
	pieces, err := jb.buildPieces()
	if err != nil {
		return jig, err
	}

	pieces, err = jb.PieceCutter.CutPieces(jb.baseImage,pieces)
	if err != nil {
		return jig, err
	}
	jig.Pieces = pieces
	return jig, nil
}

func NewJigsawBuilder(img image.Image, numPieces int) *JigsawBuilder {
	return &JigsawBuilder{numPieces: numPieces, baseImage: img, PieceCutter: JigsawPieceCutter{}}
}
func NewJigsawBuilderWithPieceCutter(img image.Image, numPieces int, cutter JigsawPieceCutter) *JigsawBuilder {
	return &JigsawBuilder{numPieces: numPieces, baseImage: img, PieceCutter: cutter}
}


type CutType struct {
	Internal bool
	Edge int //0-4
}



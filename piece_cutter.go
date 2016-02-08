package jigsaw

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/draw"
	"image"
)

const PERCENTAGE = 10.0

type PieceCutter interface {
	CutPieces(from image.Image, pieces []*Piece) ([]*Piece, error)
}

const TOP_SIDE = 0
const RIGHT_SIDE = 1
const BOTTOM_SIDE = 2
const LEFT_SIDE = 3

type PieceMarker interface {
	MarkPieces([]*Piece, int, int) []*Piece
}

func Percentage(piece *Piece, percentage float64) int {
	if piece.Height > piece.Width{
		return int((percentage / 100.0) * float64(piece.Height))
	}
	return int((percentage / 100.0) * float64(piece.Width))
}

type JigsawPieceCutter struct{}
type JigsawPieceMarker struct{}

func (JigsawPieceMarker) MarkPieces(pieces []*Piece, piecesPerRow, numRows int) []*Piece {
	//start at 0 move across
	retPieces := make([]*Piece, 0)
	for i := 0; i < len(pieces); i++ {
		p := pieces[i]
		if p.IsCorner {
			if (i+1)/piecesPerRow == 1 || (i+1)/piecesPerRow == numRows {
				//top right or bottom right
				p.Joints = []PieceJoint{{
					External: false,
					Side:     LEFT_SIDE,
				}}
				if (i+1)/piecesPerRow == numRows {
					p.Joints = append(p.Joints, PieceJoint{
						External: true,
						Side:     TOP_SIDE,
					})
				} else {
					p.Joints = append(p.Joints, PieceJoint{
						External: false,
						Side:     BOTTOM_SIDE,
					})
				}
			} else {

				p.Joints = []PieceJoint{{
					External: true,
					Side:     RIGHT_SIDE,
				}}
				if i == 0 {
					p.Joints = append(p.Joints, PieceJoint{
						External: false,
						Side:     BOTTOM_SIDE,
					})
				} else {
					p.Joints = append(p.Joints, PieceJoint{
						External: true,
						Side:     TOP_SIDE,
					})
				}
			}
		}

		retPieces = append(retPieces, p)
	}
	return retPieces
}

//cuts a rectangular piece with additional space added for any external joints
func (JigsawPieceCutter) cutPiece(from image.Image, piece *Piece) (*Piece, error) {
	//cut the piece 20%larger on the sides that have external joints

	for _, joint := range piece.Joints {
		if joint.External {
			if joint.Side == TOP_SIDE {
				// 0 is top side  (positive)
				//grow point[0].Y unless it already at max
				piece.Points[0].Y -= Percentage(piece, PERCENTAGE)
			} else if joint.Side == RIGHT_SIDE {
				//1 is right side
				//grow point[3].X (positive)unless it already at max
				piece.Points[3].X += Percentage(piece, PERCENTAGE)
			} else if joint.Side == BOTTOM_SIDE {
				//2 is bottom side
				//grow point[3].Y (by a negative value)	 unless it already at min
				fmt.Println("Y is ", piece.Points[3].Y)
				piece.Points[3].Y += Percentage(piece, PERCENTAGE)
				fmt.Println("Y is ", piece.Points[3].Y)
			} else if joint.Side == LEFT_SIDE {
				//3 is left side
				//grow point[0].X (negative)
				piece.Points[0].X -= Percentage(piece, PERCENTAGE)
			}
		}
		//non external joints are cut into the existing piece
	}
	piece.Bounds = image.Rect(piece.Points[0].X, piece.Points[0].Y, piece.Points[3].X, piece.Points[3].Y)
	//fmt.Println("BOUNDS CREATED as ", piece.Bounds)
	rectCropImg := imaging.Crop(from, piece.Bounds)
	piece.Image = rectCropImg
	path := "./out/" + piece.Name + ".png"
	// save cropped image
	err := imaging.Save(rectCropImg, path)
	piece.Path = path
	return piece, err
}

func variableRadius(basedOn *Piece) int {
	return Percentage(basedOn, PERCENTAGE)
}

func hasExternalJoint(piece *Piece, side int) bool {
	for _, j := range piece.Joints {
		if j.Side == side {
			return j.External
		}
	}
	return false
}

type JointCutter struct {
	Piece *Piece
}

func (jc JointCutter) cutInternal(joint PieceJoint, img image.Image) (image.Image, error) {
	var radius, x, y int

	if joint.Side == TOP_SIDE {
		radius = Percentage(jc.Piece, PERCENTAGE)
		x, y = img.Bounds().Max.X/2, img.Bounds().Min.Y
	} else if joint.Side == RIGHT_SIDE {
		radius = Percentage(jc.Piece, PERCENTAGE)
		x, y = img.Bounds().Max.X, img.Bounds().Max.Y/2
	} else if joint.Side == BOTTOM_SIDE {
		radius = Percentage(jc.Piece, PERCENTAGE)
		x, y = img.Bounds().Max.X/2, img.Bounds().Max.Y
	} else if joint.Side == LEFT_SIDE {
		radius = Percentage(jc.Piece, PERCENTAGE)
		percentInc := 0.0
		if hasExternalJoint(jc.Piece, 0) {
			percentInc += PERCENTAGE
		}
		if hasExternalJoint(jc.Piece, 2) {
			percentInc += PERCENTAGE
		}
		x, y = img.Bounds().Min.X, (img.Bounds().Max.Y+Percentage(jc.Piece, percentInc))/2
	}
	imageContext := image.NewRGBA(img.Bounds())
	circleImg := &circle{image.Pt(x, y), radius, img, false, image.Rectangle{}}
	draw.Draw(imageContext, img.Bounds(), circleImg, img.Bounds().Min, draw.Src)
	draw2dimg.SaveToPngFile("./out/"+jc.Piece.Name+".png", imageContext)
	return imageContext, nil
}

func (jc JointCutter) cutExternal(joint PieceJoint, from image.Image) (image.Image, error) {

	var x, y int
	var piece = jc.Piece
	var imageContext = image.NewRGBA(from.Bounds())
	var circleImg *circle
	var percentInc = 0.0

	if joint.Side == RIGHT_SIDE || joint.Side == LEFT_SIDE {
		if hasExternalJoint(jc.Piece,0) && hasExternalJoint(jc.Piece,2){
			percentInc = 0.0
		}else if hasExternalJoint(jc.Piece, 0) {
			percentInc += PERCENTAGE
		}else if hasExternalJoint(jc.Piece, 2) {
			percentInc += PERCENTAGE
		}
	} else if joint.Side == BOTTOM_SIDE || joint.Side == TOP_SIDE {
		if hasExternalJoint(jc.Piece, 3) && hasExternalJoint(jc.Piece,1){
			percentInc = 0.0
		}else if hasExternalJoint(jc.Piece, 3) {
			percentInc += PERCENTAGE
		}else if hasExternalJoint(jc.Piece, 1) {
			percentInc += PERCENTAGE
		}
	}
	if joint.Side == RIGHT_SIDE {
		//move back Percentage amount
		x, y = (from.Bounds().Max.X - Percentage(piece, PERCENTAGE)), (from.Bounds().Max.Y  + Percentage(piece,percentInc)) / 2 //half way down the right hand side
		point := image.Pt(x, y)
		circleImg = &circle{point, variableRadius(piece), from, true, image.Rect(point.X, from.Bounds().Max.Y, from.Bounds().Max.X, from.Bounds().Min.Y)}
	}
	//left
	if joint.Side == LEFT_SIDE {

		x, y = (from.Bounds().Min.X + Percentage(piece, PERCENTAGE)), (from.Bounds().Max.Y + Percentage(piece,percentInc)) /2 //half way down the left side
		point := image.Pt(x, y)
		circleImg = &circle{point, variableRadius(piece), from, true, image.Rect(point.X, from.Bounds().Min.Y, from.Bounds().Min.X, from.Bounds().Max.Y)}

	}
	//top
	if joint.Side == TOP_SIDE {
		x, y = from.Bounds().Max.X/2, from.Bounds().Min.Y+Percentage(piece, PERCENTAGE) //half way across top line
		point := image.Pt(x, y)                                                                //where the circle is centered around
		circleImg = &circle{point, variableRadius(piece), from, true, image.Rect(from.Bounds().Max.X, y, from.Bounds().Min.X, from.Bounds().Min.Y)}
	}
	//bottom
	if joint.Side == BOTTOM_SIDE {
		x, y = from.Bounds().Max.X/2, (from.Bounds().Max.Y - Percentage(piece, PERCENTAGE)) //half way across bottom line
		point := image.Pt(x, y)
		circleImg = &circle{point, variableRadius(piece), from, true, image.Rect(from.Bounds().Min.X, point.Y, from.Bounds().Max.X, from.Bounds().Max.Y)}

	}
	draw.Draw(imageContext, from.Bounds(), circleImg, from.Bounds().Min, draw.Src)
	return imageContext, nil

}

//shapes the rectangular piece removing and add joint pieces
func (JigsawPieceCutter) ShapePiece(piece *Piece) (*Piece, error) {
	//fill a semi circle with transparency
	//the piece is 20% larger on each side add any cuts then crop down sides without external piece
	jointCutter := JointCutter{Piece: piece}
	var img = piece.Image
	var err error
	for _, joint := range piece.Joints {
		if !joint.External {
			img, err = jointCutter.cutInternal(joint, img)
			if err != nil {
				return nil, err
			}
		} else {
			img, err = jointCutter.cutExternal(joint, img)

		}
	}
	draw2dimg.SaveToPngFile("./out/"+piece.Name+".png", img)

	return piece, nil
}

func (jpc JigsawPieceCutter) ShapePieces(pieces []*Piece) ([]*Piece, error) {
	shapedPieces := make([]*Piece, len(pieces))
	for i, piece := range pieces {
		shaped, err := jpc.ShapePiece(piece)
		if nil != err {
			return nil, err
		}
		shapedPieces[i] = shaped
	}
	return shapedPieces, nil
}

func (jpc JigsawPieceCutter) CutPieces(from image.Image, pieces []*Piece) ([]*Piece, error) {
	var cutPieces = make([]*Piece, len(pieces))
	for index, p := range pieces {
		cutPiece, err := jpc.cutPiece(from, p)
		if err != nil {
			return nil, errors.New("failed to cut piece " + err.Error())
		}
		cutPieces[index] = cutPiece
	}
	return jpc.ShapePieces(cutPieces)

}

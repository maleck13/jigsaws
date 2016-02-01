package jigsaw
import (
"image"
"errors"
	"github.com/disintegration/imaging"
	"github.com/llgcode/draw2d/draw2dimg"
"image/draw"
	"fmt"
)

const PERCENTAGE = 11.0

type PieceCutter interface {
	CutPieces(from image.Image, pieces []*Piece)([]*Piece,error)
}

func Percentage(total int, percentage float64)int{
	return int((percentage /  100.0) * float64(total))
}

type JigsawPieceCutter struct {}
//cuts a rectangular piece with additional space added for any external joints
func (JigsawPieceCutter)cutPiece(from image.Image, piece *Piece)(*Piece,error)  {
	//cut the piece 20%larger on the sides that have external joints

	for _,joint:= range piece.Joints{
		if joint.External{
			if joint.Side == 0{
				// 0 is top side  (positive)
				//grow point[0].Y unless it already at max
				piece.Points[0].Y-= Percentage(piece.Height,PERCENTAGE)
			}else if joint.Side == 1{
				//1 is right side
				//grow point[3].X (positive)unless it already at max
				piece.Points[3].X+= Percentage(piece.Width, PERCENTAGE)
			}else if joint.Side == 2{
				//2 is bottom side
				//grow point[3].Y (by a negative value)	 unless it already at min
				piece.Points[3].Y+= Percentage(piece.Height, PERCENTAGE)
			}else if joint.Side == 3{
				//3 is left side
				//grow point[0].X (negative)
				piece.Points[0].X-= Percentage(piece.Width, PERCENTAGE)
			}
		}
		//non external joints are cut into the existing piece
	}
	piece.Bounds = image.Rect(piece.Points[0].X, piece.Points[0].Y, piece.Points[3].X, piece.Points[3].Y)
	rectCropImg := imaging.Crop(from, piece.Bounds )
	piece.Image = rectCropImg
	path := "./out/" + piece.Name + ".jpg"
	// save cropped image
	err := imaging.Save(rectCropImg, path)
	piece.Path = path
	return piece, err
}

//shapes the rectangular piece removing and add joint pieces
func (JigsawPieceCutter)shapePiece(piece *Piece)(*Piece, error)  {
	//fill a semi circle with transparency
	//the piece is 20% larger on each side add any cuts then crop down sides without external piece
	var circleImg *circle

	var radius int = Percentage(piece.Bounds.Max.Y,PERCENTAGE /2)
	var imageContext *image.RGBA
	var x,y int
	img := piece.Image
	bounds := img.Bounds()

	for _,cut := range piece.Joints{
		if ! cut.External{


			if cut.Side == 0{
				x,y = bounds.Max.X/2, bounds.Max.Y
			}else if cut.Side == 1{
				x,y = bounds.Max.X, bounds.Max.Y / 2
			}else if cut.Side == 2{
				x,y = bounds.Max.X /2 , bounds.Min.Y
			}else if cut.Side == 3{
				x,y = bounds.Min.X , bounds.Max.Y /2
			}
			imageContext = image.NewRGBA(bounds)
			circleImg = &circle{image.Pt(x,y), radius, img, false}
			draw.Draw(imageContext, img.Bounds(), circleImg, img.Bounds().Min, draw.Src)
			fmt.Println("saving png to file")
			draw2dimg.SaveToPngFile("./test.png", imageContext)
			img = imageContext
		}else{
			if cut.Side == 2{
				x,y = bounds.Max.X / 2, bounds.Max.Y //half way across bottom line
				//need to move back 20 %
				y = y - Percentage(y,PERCENTAGE)
				//draw a circle filled with the background
				imageContext = image.NewRGBA(bounds)
				circleImg = &circle{image.Pt(x,y), radius, img, true}
				draw.Draw(imageContext, bounds, circleImg, bounds.Min, draw.Src)
				draw2dimg.SaveToPngFile("./testout.png", imageContext)
				img = imageContext
				//draw rectangle that cuts across

			}

		}
	}


	return piece,nil
}

func (jpc JigsawPieceCutter)shapePieces(pieces []*Piece)([]*Piece, error){
	shapedPieces := make([]*Piece,len(pieces))
	for i,piece := range pieces{
		shaped,err :=jpc.shapePiece(piece)
		if nil != err{
			return nil,err
		}
		shapedPieces[i] = shaped
	}
	return shapedPieces,nil
}

func (jpc JigsawPieceCutter)CutPieces(from image.Image, pieces []*Piece)([]*Piece,error) {
	var cutPieces = make([]*Piece,len(pieces))
	for index,p := range pieces{
		cutPiece,err := jpc.cutPiece(from,p)
		if err != nil{
			return nil, errors.New("failed to cut piece " + err.Error())
		}
		cutPieces[index] = cutPiece
	}
	return jpc.shapePieces(cutPieces)

}

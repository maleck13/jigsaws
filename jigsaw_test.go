package jigsaw_test

import (
	"github.com/maleck13/jigsaw"
	"github.com/stretchr/testify/assert"
	"image"
	"image/jpeg"
	"os"
	"testing"
)

const JPG_SAMPLE string = "./samples/firetruck.jpg"

func openImage(t *testing.T, path string) image.Image {
	fireTruck, err := os.Open(path)
	assert.NoError(t, err, "error openenig image")
	img, err := jpeg.Decode(fireTruck)
	assert.NoError(t, err, "error decoding image")
	return img
}

//func TestCreateJigsaw(t *testing.T) {
//	img := openImage(t, JPG_SAMPLE)
//	builder := jigsaw.NewJigsawBuilder(img, 24)
//	jig, err := builder.Build()
//	assert.NoError(t, err, "did not expect an error")
//	assert.NotNil(t, jig, "expected jigsaw")
//}
//
//func TestJigsawShouldHave4Corners(t *testing.T) {
//	img := openImage(t, JPG_SAMPLE)
//	builder := jigsaw.NewJigsawBuilder(img, 24)
//	jig, err := builder.Build()
//	assert.NoError(t, err, "did not expect an error")
//	assert.NotNil(t, jig, "expected jigsaw")
//	corners := 0
//	for _, p := range jig.Pieces {
//		if p.IsCorner {
//			corners++
//		}
//	}
//	assert.True(t, corners == 4, "expected four corners")
//}
//
//func Test8PieceJigsaw(t *testing.T) {
//	img := openImage(t, JPG_SAMPLE)
//	builder := jigsaw.NewJigsawBuilder(img, 12)
//	jig, err := builder.Build()
//	assert.NoError(t, err, "did not expect an error")
//	assert.NotNil(t, jig, "expected jigsaw")
//	edgeCount := 0
//	for _, p := range jig.Pieces {
//		if p.IsEdge && !p.IsCorner {
//			edgeCount++
//		}
//	}
//	assert.Equal(t, 6, edgeCount, "expected 6 non corner edge pieces")
//}

//func Test48PieceJigsaw(t *testing.T) {
//	img := openImage(t, JPG_SAMPLE)
//	builder := jigsaw.NewJigsawBuilder(img, 48)
//	jig, err := builder.Build()
//	assert.NoError(t, err, "did not expect an error")
//	assert.NotNil(t, jig, "expected jigsaw")
//	edgeCount := 0
//	cornerCount := 0
//	for _, p := range jig.Pieces {
//		if p.IsCorner {
//			cornerCount++
//		}
//		if p.IsEdge{
//			edgeCount++
//		}
//	}
//	assert.Equal(t, 4, cornerCount, "expected 4 corner pieces")
//	assert.Equal(t, 24, edgeCount, "expected 8 edge pieces")
//
//}

func Test4pieceJigsaw(t *testing.T) {
	img := openImage(t, JPG_SAMPLE)
	builder := jigsaw.NewJigsawBuilder(img, 12)
	jig, err := builder.Build()
	assert.NoError(t, err, "did not expect an error")
	assert.NotNil(t, jig, "expected jigsaw")
	//	for _, p := range jig.Pieces {
	//		assert.True(t, p.IsCorner, "should be a corner")
	//	}
}

func TestShapePiece0and1(t *testing.T) {
	img := openImage(t, JPG_SAMPLE)
	builder := jigsaw.NewJigsawBuilder(img, 12)
	builder.NumRows = 3
	pieces, err := builder.BuildPieces()
	assert.NoError(t, err, "not expected error")
	piece := pieces[0]
	piece.Joints = []jigsaw.PieceJoint{
		jigsaw.PieceJoint{
			Side:     0,
			External: false,
		},
		jigsaw.PieceJoint{
			Side:     1,
			External: false,
		},
		jigsaw.PieceJoint{
			Side:     2,
			External: false,
		},
		jigsaw.PieceJoint{
			Side:     3,
			External: true,
		},
	}
	pieces, err = builder.PieceCutter.CutPieces(img, pieces)
	assert.NotNil(t, pieces, "pieces expected")

	piece.Name = "testpiece0"

	pieceCutter := jigsaw.JigsawPieceCutter{}
	pieceCutter.ShapePiece(piece)

}

func TestShapePieceSide1(t *testing.T) {
	img := openImage(t, JPG_SAMPLE)
	builder := jigsaw.NewJigsawBuilder(img, 12)
	builder.NumRows = 3
	pieces, err := builder.BuildPieces()
	assert.NoError(t, err, "not expected error")
	pieces, err = builder.PieceCutter.CutPieces(img, pieces)
	assert.NotNil(t, pieces, "pieces expected")
	piece := pieces[0]
	piece.Name = "testpiece1"
	piece.Joints = []jigsaw.PieceJoint{
		jigsaw.PieceJoint{
			Side:     1,
			External: true,
		},
	}
	pieceCutter := jigsaw.JigsawPieceCutter{}
	pieceCutter.ShapePiece(piece)

}

func TestShapePieceSide2(t *testing.T) {
	img := openImage(t, JPG_SAMPLE)
	builder := jigsaw.NewJigsawBuilder(img, 12)
	builder.NumRows = 3
	pieces, err := builder.BuildPieces()
	assert.NoError(t, err, "not expected error")
	piece := pieces[0]
	piece.Name = "testpiece2"
	piece.Joints = []jigsaw.PieceJoint{
		jigsaw.PieceJoint{
			Side:     2,
			External: true,
		},
	}
	pieces, err = builder.PieceCutter.CutPieces(img, pieces)
	assert.NotNil(t, pieces, "pieces expected")

	pieceCutter := jigsaw.JigsawPieceCutter{}
	pieceCutter.ShapePiece(piece)

}

func TestShapePieceSide3(t *testing.T) {
	img := openImage(t, JPG_SAMPLE)
	builder := jigsaw.NewJigsawBuilder(img, 12)
	builder.NumRows = 3
	pieces, err := builder.BuildPieces()
	assert.NoError(t, err, "not expected error")
	pieces, err = builder.PieceCutter.CutPieces(img, pieces)
	assert.NotNil(t, pieces, "pieces expected")
	piece := pieces[0]
	piece.Name = "testpiece3"
	piece.Joints = []jigsaw.PieceJoint{
		jigsaw.PieceJoint{
			Side:     3,
			External: true,
		},
	}
	pieceCutter := jigsaw.JigsawPieceCutter{}
	pieceCutter.ShapePiece(piece)

}

//func TestJigsawShouldHave16EdgePieces(t *testing.T) {
//	img := openImage(t, JPG_SAMPLE)
//	builder := jigsaw.NewJigsawBuilder(img, 24)
//	jig, err := builder.Build()
//	assert.NoError(t, err, "did not expect an error")
//	assert.NotNil(t, jig, "expected jigsaw")
//	edge := 0
//	for _, p := range jig.Pieces {
//		if p.IsEdge {
//			edge++
//		}
//	}
//	assert.True(t, edge == 16, "expected 16 edge pieces")
//}
//
//func TestInvalidNumberOfPieces(t *testing.T) {
//	img := openImage(t, JPG_SAMPLE)
//	builder := jigsaw.NewJigsawBuilder(img, 10)
//	jig, err := builder.Build()
//	assert.Error(t, err, "expected an error")
//	assert.NotNil(t, jig, "expected jigsaw")
//
//}
//
//func TestShapePiece(t *testing.T) {
//	img := openImage(t, "./piece35.jpg")
//	cuts := []jigsaw.CutType{jigsaw.CutType{
//		Edge:1,
//		Internal:true,
//	},
//	jigsaw.CutType{
//		Edge:2,
//		Internal:false,
//	}}
//	jigsaw.ShapePieces(img,cuts)
//}
//
//func TestCornersHave2Neighbours(t *testing.T)  {
//	img := openImage(t, JPG_SAMPLE)
//	builder := jigsaw.NewJigsawBuilder(img, 24)
//	jig, err := builder.Build()
//	assert.NoError(t, err, "did not expect an error")
//	assert.NotNil(t, jig, "expected jigsaw")
//	for _,p:= range jig.Pieces{
//		if p.IsCorner && p.TopRow(){
//			fmt.Print(p.RightPieceIndex, p.LeftPieceIndex,p.BottomPieceIndex)
//		}
//	}
//}
//
//func TestCenterPieces12(t *testing.T)  {
//	img := openImage(t, JPG_SAMPLE)
//	builder := jigsaw.NewJigsawBuilder(img, 12)
//	jig, err := builder.Build()
//	assert.NoError(t, err, "did not expect an error")
//	assert.NotNil(t, jig, "expected jigsaw")
//	center :=0
//	for _,p:= range jig.Pieces{
//		if p.IsCenter{
//			center++
//		}
//	}
//	assert.Equal(t,2,center,"expected 2 center pieces")
//}
//
//func TestCenterPieces24(t *testing.T)  {
//	img := openImage(t, JPG_SAMPLE)
//	builder := jigsaw.NewJigsawBuilder(img, 24)
//	jig, err := builder.Build()
//	assert.NoError(t, err, "did not expect an error")
//	assert.NotNil(t, jig, "expected jigsaw")
//	center :=0
//	for _,p:= range jig.Pieces{
//		if p.IsCenter{
//			center++
//		}
//	}
//	assert.Equal(t,8,center,"expected 2 center pieces")
//}

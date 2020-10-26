package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

func main(){
	window := gocv.NewWindow("Test")
	defer window.Close()

	match := gocv.IMRead("./car/5.jpg",gocv.IMReadGrayScale)
	temp := gocv.IMRead("./match/su2.png",gocv.IMReadGrayScale)
	imgs := gocv.IMRead("./car/5.jpg",gocv.IMReadColor)

	w := match.Cols()
	h := match.Rows()
	window.ResizeWindow(int(w), int(h))


	result := gocv.NewMat()
	defer result.Close()
	m := gocv.NewMat()
	gocv.MatchTemplate(match,temp,&result,gocv.TmCcoeff,m)
	m.Close()
	x1, x2, _, max_loc := gocv.MinMaxLoc(result)
	fmt.Println("匹配度：" , x1,x2)
	rect := image.Rect(max_loc.X,
		max_loc.Y,
		max_loc.X + temp.Cols(),
		max_loc.Y + temp.Rows())
	fmt.Println(rect)
	gocv.Rectangle(&imgs, rect, color.RGBA{R: 255}, 2)
	for {
		window.IMShow(imgs)
		//window.WaitKey(1)
		//// 图片处理完毕记得关闭以释放内存
		//match.Close()
		//temp.Close()
		if window.WaitKey(1) >= 0 {
			match.Close()
			temp.Close()
			break
		}
	}
}

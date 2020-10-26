package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

func main(){
	window := gocv.NewWindow("车牌检测gocv实例")
	defer window.Close()
	src := gocv.IMRead("./car/timg.jpg",gocv.IMReadColor)
	car := gocv.IMRead("./car/timg.jpg",gocv.IMReadGrayScale)
	defer src.Close()
	defer car.Close()
	window.ResizeWindow(car.Cols(),car.Rows())
	out := gocv.NewMat()
	defer out.Close()
	//gocv.Canny(car,&car,500,200)
	gocv.Threshold(car,&out,0,255,gocv.ThresholdBinaryInv|gocv.ThresholdOtsu)
	kernel_x := gocv.GetStructuringElement(gocv.MorphRect,image.Pt(5,8))
	gocv.MorphologyEx(out,&out,gocv.MorphClose,kernel_x)
	kernel_y := gocv.GetStructuringElement(gocv.MorphRect,image.Pt(1,15))
	gocv.MorphologyEx(out,&out,gocv.MorphOpen,kernel_y)
	kernel_c := gocv.GetStructuringElement(gocv.MorphRect,image.Pt(10,10))
	gocv.MorphologyEx(out,&out,gocv.MorphOpen,kernel_c)
	//gocv.MedianBlur(out,&out,15)
	point := gocv.FindContours(out,gocv.RetrievalTree,gocv.ChainApproxSimple)

	var area_s gocv.RotatedRect


	for _,v := range point{
		area := gocv.MinAreaRect(v)
		//fmt.Print(area.Width)
		//fmt.Println("---",area.Height)
		if area.Width == 0 || area.Height == 0{
			continue
		}
		imgrect := gocv.BoundingRect(v)
		//if k == 7 {
		//	fmt.Println(src.Cols())
		//	fmt.Println(imgrect.Min)
		//	fmt.Println(imgrect.Max.X)
		//}
		if src.Cols()/3 < imgrect.Max.X - imgrect.Min.X{
			continue
		}
		if area.Height / area.Width < 2 || area.Height / area.Width >= 4{
			continue
		}
		//fmt.Println(area.Height)
		//fmt.Println(area.BoundingRect)
		//gocv.DrawContours(&src,point,k,color.RGBA{G:255},1)
		//gocv.PutText(&src,strconv.Itoa(k),v[0], gocv.FontHersheyPlain, 1.2, color.RGBA{R:255}, 2)
		//gocv.Rectangle(&src,imgrect,color.RGBA{R:255},1)
		area_s = area
		//areas := gocv.MinAreaRect(point[k])

	}

	//for _,v := range point{
	//	imgrect := gocv.BoundingRect(v)
	//	gocv.Rectangle(&src,imgrect,color.RGBA{G:255},1)
	//}
	//gocv.MedianBlur(out,&out,15)
	//fmt.Println(area_s.Contour)
	//画出轮廓顶点
	//for k,v := range area_s.Contour{
	//	gocv.PutText(&src,strconv.Itoa(k),v,gocv.FontHersheyPlain, 1.2, color.RGBA{R:255}, 2)
	//}
	//根据顶点修正图片角度
	srcs := []image.Point{
		image.Pt(area_s.Contour[2].X, area_s.Contour[2].Y),
		image.Pt(area_s.Contour[3].X, area_s.Contour[3].Y),
		image.Pt(area_s.Contour[1].X, area_s.Contour[1].Y),
		image.Pt(area_s.Contour[0].X, area_s.Contour[0].Y),
	}
	dsts := []image.Point{
		image.Pt(0, 0),
		image.Pt(area_s.Height, 0),
		image.Pt(0, area_s.Width),
		image.Pt(area_s.Height, area_s.Width),		//坑1：这里默认长的为height，短的为width
	}
	m := gocv.GetPerspectiveTransform(srcs, dsts)
	defer m.Close()

	dstss := gocv.NewMat()
	defer dstss.Close()

	gocv.WarpPerspective(src, &dstss, m, image.Pt(src.Cols(), src.Rows()))

	cars := dstss.Region(image.Rect(0,0,area_s.Height,area_s.Width))		//裁切出的车牌  彩色
	cars_tmp := gocv.NewMat()
	cars.CopyTo(&cars_tmp)

	gocv.CvtColor(cars,&cars,gocv.ColorBGRToGray)
	gocv.Threshold(cars,&cars,0,255,gocv.ThresholdBinaryInv|gocv.ThresholdOtsu)
	point = gocv.FindContours(cars,gocv.RetrievalList,gocv.ChainApproxSimple)
	for _,v := range point{
		area := gocv.MinAreaRect(v)
		fmt.Println(area.Width,"+++",area.Height)
		fmt.Println(float32(area.Height)/float32(area.Width))
		if area.Height > cars_tmp.Cols()/2 || area.Height < area.Width||float32(area.Height)/float32(area.Width) > 2.5||float32(area.Height)/float32(area.Width) < 2{
			continue
		}
		imgrect := gocv.BoundingRect(v)
		gocv.Rectangle(&cars_tmp,imgrect,color.RGBA{G:255},1)
	}



	window.ResizeWindow(cars.Cols(),cars.Rows())
	for{
		window.IMShow(cars_tmp)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}


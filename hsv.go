package main

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

func main() {
	window := gocv.NewWindow("gocv test")
	defer window.Close()
	src := gocv.IMRead("./car/8.jpg", gocv.IMReadColor)
	car := gocv.NewMat()
	gocv.CvtColor(src, &car, gocv.ColorBGRToGray)
	defer src.Close()
	defer car.Close()
	window.ResizeWindow(car.Cols(), car.Rows())
	out := gocv.NewMat()
	defer out.Close()

	gocv.CvtColor(src, &out, gocv.ColorBGRToHSV)
	maskblue := gocv.NewMat()
	defer maskblue.Close()
	//gocv.InRangeWithScalar(out,gocv.NewScalar(100.0,80.0,85.0,0.0),gocv.NewScalar(124.0,255.0,255.0,0.0),&maskblue)	//HSV提取蓝色部分
	gocv.InRangeWithScalar(out, gocv.NewScalar(105.0, 100.0, 85.0, 0.0), gocv.NewScalar(124.0, 255.0, 255.0, 0.0), &maskblue)
	gocv.BitwiseNot(maskblue, &out)

	kernel_open := gocv.GetStructuringElement(gocv.MorphRect,image.Pt(15,10))
	gocv.MorphologyEx(out,&out,gocv.MorphOpen,kernel_open)
	kernel_close := gocv.GetStructuringElement(gocv.MorphRect,image.Pt(20,15))
	gocv.MorphologyEx(out,&out,gocv.MorphClose,kernel_close)

	point := gocv.FindContours(out,gocv.RetrievalTree,gocv.ChainApproxSimple)
	for _,v := range point{
		area := gocv.MinAreaRect(v)
		imgrect := gocv.BoundingRect(v)
		//fmt.Println(MaxMin(area.Width,area.Height))
		if MaxMin_t(area.Width,area.Height) > 3.6 || MaxMin_t(area.Width,area.Height) < 2.2{continue}
		//if area.Height * area.Width < 5000 || area.Height * area.Width > 12000{continue}
		gocv.Rectangle(&src,imgrect,color.RGBA{G:255},1)
		//gocv.PutText(&src,strconv.Itoa(k),v[0], gocv.FontHersheyPlain, 1.2, color.RGBA{R:255}, 2)
		//gocv.DrawContours(&src,point,k,color.RGBA{G:255},1)
		//fmt.Println(k,":",float32(area.Width),"*",float32(area.Height),"r:",area.Angle,"J:",area.Height * area.Width,"B:",MaxMin(area.Width,area.Height))
	}

	for {
		window.IMShow(src)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

func MaxMin_t(w,h int)(s float32){
	if h > w{
		s = float32(h) / float32(w)
	}else{
		s = float32(w) / float32(h)
	}
	return
}
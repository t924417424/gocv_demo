package main

import (
	"gocv.io/x/gocv"
	"image"
	"sort"
)

func main(){
	window := gocv.NewWindow("gocv test")
	defer window.Close()
	src := gocv.IMRead("./car/3.jpg",gocv.IMReadColor)
	car := gocv.NewMat()
	gocv.CvtColor(src,&car,gocv.ColorBGRToGray)
	defer src.Close()
	defer car.Close()
	window.ResizeWindow(car.Cols(),car.Rows())
	out := gocv.NewMat()
	defer out.Close()

	gocv.CvtColor(src,&out,gocv.ColorBGRToHSV)
	maskblue := gocv.NewMat()
	defer maskblue.Close()
	//gocv.InRangeWithScalar(out,gocv.NewScalar(100.0,80.0,85.0,0.0),gocv.NewScalar(124.0,255.0,255.0,0.0),&maskblue)	//HSV提取蓝色部分
	gocv.InRangeWithScalar(out,gocv.NewScalar(105.0,100.0,85.0,0.0),gocv.NewScalar(124.0,255.0,255.0,0.0),&maskblue)
	gocv.BitwiseNot(maskblue,&out)


	kernel_open := gocv.GetStructuringElement(gocv.MorphRect,image.Pt(15,10))
	gocv.MorphologyEx(out,&out,gocv.MorphOpen,kernel_open)
	kernel_close := gocv.GetStructuringElement(gocv.MorphRect,image.Pt(20,15))
	gocv.MorphologyEx(out,&out,gocv.MorphClose,kernel_close)


	//gocv.MedianBlur(out,&out,15)
	point := gocv.FindContours(out,gocv.RetrievalTree,gocv.ChainApproxSimple)
	var area_pt gocv.RotatedRect
	for _,v := range point{
		area := gocv.MinAreaRect(v)
		//imgrect := gocv.BoundingRect(v)
		//fmt.Println(MaxMin(area.Width,area.Height))
		if MaxMin(area.Width,area.Height) > 3.6 || MaxMin(area.Width,area.Height) < 2.2{continue}
		//if area.Height * area.Width < 5000 || area.Height * area.Width > 12000{continue}
		//gocv.Rectangle(&src,imgrect,color.RGBA{G:255},1)
		//gocv.PutText(&src,strconv.Itoa(k),v[0], gocv.FontHersheyPlain, 1.2, color.RGBA{R:255}, 2)
		//gocv.DrawContours(&src,point,k,color.RGBA{G:255},1)
		//fmt.Println(k,":",float32(area.Width),"*",float32(area.Height),"r:",area.Angle,"J:",area.Height * area.Width,"B:",MaxMin(area.Width,area.Height))
		area_pt = area
	}

	//画出轮廓顶点  坑：因为坐标点的顺序可能出现不一致，需要归一化之后再做图像修正 规则：x小则为左侧，x大则为右侧，根据y值判断是上边角还是下边角
	img_pt_sort(area_pt.Contour)	//对顶点进行排序
	//for k,v := range area_pt.Contour{
	//	gocv.PutText(&src,strconv.Itoa(k),v,gocv.FontHersheyPlain, 1.2, color.RGBA{R:255}, 2)
	//}

	//根据顶点修正图片角度
	srcs := []image.Point{
		image.Pt(area_pt.Contour[0].X, area_pt.Contour[0].Y),
		image.Pt(area_pt.Contour[2].X, area_pt.Contour[2].Y),
		image.Pt(area_pt.Contour[1].X, area_pt.Contour[1].Y),
		image.Pt(area_pt.Contour[3].X, area_pt.Contour[3].Y),
	}
	dsts := []image.Point{
		image.Pt(0, 0),
		image.Pt(area_pt.Height, 0),
		image.Pt(0, area_pt.Width),
		image.Pt(area_pt.Height, area_pt.Width),
	}
	m := gocv.GetPerspectiveTransform(srcs, dsts)
	defer m.Close()
	correct := gocv.NewMat()
	defer correct.Close()
	gocv.WarpPerspective(src, &correct, m, image.Pt(src.Cols(), src.Rows()))

	cars := correct.Region(image.Rect(0,0,area_pt.Height,area_pt.Width))		//裁切出的车牌  彩色
	//fmt.Println(area_pt.Contour)
	window.ResizeWindow(cars.Cols() * 3,cars.Cols())
	for{
		window.IMShow(cars)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

func MaxMin(w,h int)(s float32){
	if h > w{
		s = float32(h) / float32(w)
	}else{
		s = float32(w) / float32(h)
	}
	return
}

func img_pt_sort(area_pt []image.Point){
	sort.Slice(area_pt, func(i, j int) bool {
		return area_pt[i].X * area_pt[i].Y < area_pt[j].X * area_pt[j].Y
	})
}
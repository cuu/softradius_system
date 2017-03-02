package models

import (
	"fmt"
	"math"
	"strconv"
	//"strings"
	"net/url"
	
)


type Pager struct{
	Page     int
	Total    int
	PageSize int
	PageNum  int
	Url      string
	PageBar map[int][]int
}


func NewPager(total int, size int ,page int, url string) *Pager {
	one := &Pager{}
	one.Init(total,size,page,url)
	
	return one
}


func (this *Pager) Init(total int, size int, page int,url string) {
	if total > 0 {
		this.PageNum = int( math.Ceil( float64(total/size)) )
	}else {
		this.PageNum = 0
	}

	this.PageSize = size //每页显示的条目数
	
	this.PageBar = make(map[int][]int)
	
	for i:=1;i<= this.PageNum + 1;i++ {
		_idx := int(i/10)
		this.PageBar[_idx] = append( this.PageBar[_idx],i)
	}
	
	this.Page = page
	this.Url = url
	
}

// repack url
func (this *Pager) PageUrl(p int) string {

	u,err := url.Parse(this.Url)
	if err == nil {
		q := u.Query()
		q.Set("page_id",strconv.Itoa(p))
		u.RawQuery = q.Encode()
		return fmt.Sprintf("%v",u)
	}

	return "#"
	
		
}

func (this *Pager) Render() string {
	var out STRINGS

	out.Append(`<ul class="pagination">`)
	out.Append(fmt.Sprintf("<li class=\"disabled\"><a href=\"#\">查询记录数 %d</a></li>",this.Total))

	current_start := this.Page
	if current_start == 1 {
		out.Append(`<li class="disabled"><a href="#">首页</a></li>`)
		out.Append(`<li class="disabled"><a href="#">&larr; 上一页</a></li>`)
		
	}else {
		out.Append(fmt.Sprintf("<li><a href=\"%s\">首页</a></li>",this.PageUrl(1)))
		out.Append(fmt.Sprintf("<li><a href=\"%s\">&larr; 上一页</a></li>", this.PageUrl( this.Page - 1) ))
	}

	for i :=0; i<len(this.PageBar);i++ {
		for j:=0;j<len(this.PageBar[i]);j++{
			if this.Page == this.PageBar[i][j] {
				out.Append(fmt.Sprintf("<li class=\"active\"><span>%d <span class=\"sr-only\">{current}</span    ></span></li>", this.Page))
			}else {
				out.Append(fmt.Sprintf("<li><a href=\"%s\">%d</a></li>",this.PageUrl(this.PageBar[i][j]),this.PageBar[i][j]))
			}
		}
	}
	
	current_end := this.Page
	if current_end == this.PageNum {
		out.Append("<li class=\"disabled\"><a href=\"#\">下一页 &rarr;</a></li>")
		out.Append("<li class=\"disabled\"><a href=\"#\">尾页</a></li>")
	}else {
		out.Append(fmt.Sprintf("<li><a href=\"%s\">下一页 &rarr;</a></li>",this.PageUrl(this.Page+1)))
		out.Append(fmt.Sprintf("<li><a href=\"%s\">尾页</a></li>",this.PageUrl(this.PageNum)))
	}

	out.Append("</ul>")

	return out.String()
		
}

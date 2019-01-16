package main

import (
	"fmt"
	"strings"
	"os"
	"os/exec"
	"io/ioutil"
	"runtime"
	"net"
	"net/http"
	"strconv"
	"log"
)

func File2Tag(u,t string) string{
	c:="<"+t+">"
	f,e:=os.Open(u)
	if(e==nil){
		b,e0:=ioutil.ReadAll(f)
		if(e0==nil){
			c+=fmt.Sprint(string(b))
		}
	}
	f.Close()
	c+="</"+t+">"
	return c
}

func File2Diary(u string) string {	
	var c string
	f,e:=os.Open(u)
	if(e==nil){
		b,e0:=ioutil.ReadAll(f)
		if(e0==nil){
			c=fmt.Sprint(string(b))
			if len(c)>2048 {c=string([]rune(c)[:2048])}
		}
	}
	f.Close()
	return c
}

func localhost(addr string){
	addr=`http://`+addr
	switch runtime.GOOS{
		case "windows":
			exec.Command("rundll32.exe","url.dll,FileProtocolHandler",addr).Start()
		case "linux":
			exec.Command("xdg-open",addr).Start()
		case "darwin":
			exec.Command("open",addr).Start()
	}
}

func IMax(a,b int) int {
	if b>a {return b}
	return a
}
func IMin(a,b int) int {
	if b<a {return b}
	return a
}

type Sea struct {	
	Ships[]int
	Field[]rune
	FieldV string
	Trace[]rune
	X int
	Y int
	Left int
	Score int
	Diary2048 string
}
func NewSea()*Sea {return &Sea {[]int{4,3,2,1},[]rune(strings.Repeat("0",100)),"",[]rune(strings.Repeat("0",100)),-1,-1,0,0,""}}

func (this *Sea) CheckPiece(a,b int) bool {
	for i:=-1;i<=1;i++{
		for j:=-1;j<=1;j++{
			y,x:=a+i,b+j
			if x<0||x>9||y<0||y>9 {continue}			
			if this.Field[y*10+x]=='1' {return false}
		}
	}
	return true
}

func (this *Sea) PutPiece(a,b int) {
	for i:=-1;i<=1;i++{
		for j:=-1;j<=1;j++{
			y,x:=a+i,b+j
			if x<0||x>9||y<0||y>9 {continue}
			pos:=y*10+x
			if this.Field[pos]=='1' {continue}
			if i==0&&j==0 {this.Field[pos]='1'} else {this.Field[pos]='2'}
		}
	}	
}

func (this *Sea) BuildShip(y0,x0,y1,x1 int) int {	
	x:=[]int{-IMin(y0,x0)-1,IMax(y0,x0)-1,-IMin(y1,x1)-1,IMax(y1,x1)-1}	
	ys,xs,yf,xf:=IMin(x[0],x[2]),IMin(x[1],x[3]),IMax(x[0],x[2]),IMax(x[1],x[3])
	yl,xl:=yf-ys+1,xf-xs+1
	if yl!=1&&xl!=1 {return 5}	
	for i:=range(x){if x[i]<0 || x[i]>9 {return 5}}
	zl:=xl+yl-1
	if zl>4 {return 5}
	if this.Ships[zl-1]<=0 {return zl}
	for i:=ys;i<=yf;i++{
		for j:=xs;j<=xf;j++{
			if !this.CheckPiece(i,j){return zl}
		}
	}
	for i:=ys;i<=yf;i++{
		for j:=xs;j<=xf;j++{
			this.PutPiece(i,j)
		}
	}	
	this.Ships[zl-1]-=1
	this.Left+=1
	return 0
}

func (this *Sea) SetShip (y0,x0,y1,x1 int) int {
	fmt.Print("(",y0,",",x0,")-(",y1,",",x1,") : ")
	e:=this.BuildShip(y0,x0,y1,x1)
	var s string
	var u int
	if(e==5){s="Ошибка установки корабля (-5 баллов)";u=-5
	}else if(e==0){s="Корабль установлен успешно";u=0
	}else {u=e-5;s=fmt.Sprint("Ошибка установки ",e,"-хпалубного корабля (",u," баллов)")}
	fmt.Println(s)
	return u
}

func (this *Sea) Print () {
	for i:=0;i<10;i++{fmt.Println(string(this.Field[i*10:(i+1)*10]))}
}

func (this *Sea) Expose () {
	for i:=0;i<10;i++{fmt.Println(string(this.Trace[i*10:(i+1)*10]))}
}

func (this *Sea) Redescribe () {
	this.FieldV=strings.Replace(string(this.Field),"2","0",-1)
}

func (this *Sea) ReasonableHit(a,b int) int {
	s:=[]int{}
	test:=false
	for y:=0;y<10;y++{
		for x:=0;x<10;x++{
			pos:=y*10+x
			if this.Trace[pos]=='1' {
				test=true
				for i:=-1;i<=1;i+=2{
					for j:=-1;j<=1;j+=2{
						x0,y0:=x+j,y+i
						if x0<0||x0>9||y0<0||y0>9 {continue}
						if y0==a&&x0==b {
							fmt.Println("Выстрел по полю, где заведомо не могло быть корабля (-7 баллов)")
							return -7
						}
						s=append(s,y0*10+x0)
					}
				}
				for u:=0;u<2;u++{			
					for k:=-1;k<=1;k+=2{
						x0,y0:=x,y
						if u==0 {x0+=k} else {y0+=k}
						if x0<0||x0>9||y0<0||y0>9 {continue}
						pos0:=y0*10+x0
						tt:=true
						for t:=range(s){
							if s[t]==pos0 {
								tt=false
								break
							}
						}
						if tt {if y0==a&&x0==b {return 0}}
					}
				}
			}
		}
	}
	if test {
		fmt.Println("Нерациональный выстрел по полю вдалеке от раненого корабля (-6 баллов)")
		return -6
	}
	return 0
}

func (this *Sea) CheckBorders(a,b,c int,h bool) int {
	y0:=a
	var pos int
	for true {
		y1:=y0+c
		if h {pos=y1*10+b} else {pos=b*10+y1}
		if y1<0||y1>9 {return y0}
		if this.Field[pos]!='1' {return y0}
		y0=y1
	}
	return -1
}

func (this *Sea) CheckDestruction(a,b int) bool {
	y0:=this.CheckBorders(a,b,-1,true)
	y1:=this.CheckBorders(a,b,1,true)
	x0:=this.CheckBorders(b,a,-1,false)
	x1:=this.CheckBorders(b,a,1,false)
	for y:=y0;y<=y1;y++ {
		for x:=x0;x<=x1;x++ {
			if this.Trace[y*10+x]!='1' {
				return false
			}
		}
	}
	for y:=y0;y<=y1;y++ {
		for x:=x0;x<=x1;x++ {
			this.Trace[y*10+x]='2'
		}
	}
	this.Left-=1
	return true
}

func (this *Sea) HitShip (y0,x0 int) int {
	GameState=0
	fmt.Println(fmt.Sprint("Удар по точке (",y0,",",x0,"):"))
	score:=0
	x:=[]int{-IMin(y0,x0)-1,IMax(y0,x0)-1}
	if this.Y==x[0]&&this.X==x[1] {
		fmt.Println("Выстрел по полю, по которому только что уже был осуществлен выстрел (-20 баллов)")
		score-=20
	}
	this.Y,this.X=x[0],x[1]
	for i:=range(x){
		if x[i]<0 || x[i]>9 {
			fmt.Println("Выстрел за границу поля (-8 баллов)")
			return score-8
		}
	}
	if strings.Contains(string(this.Field),"1") {score+=this.ReasonableHit(x[0],x[1])}
	pos:=x[0]*10+x[1]
	m:=this.Trace[pos]
	if m=='3' {
		fmt.Println("Выстрел по заведомо пустому полю (-11 баллов)")
		score-=11		
	}else if m=='2' {
		fmt.Println("Выстрел по уже затонувшему кораблю (-10 баллов)")
		score-=10		
	}else if m=='1' {
		fmt.Println("Выстрел по обломкам корабля (-9 баллов)")
		score-=9		
	}else if m=='0' {
		s:=this.Field[pos]
		if s=='1' {
			GameState=1
			this.Trace[pos]='1'
			if this.CheckDestruction(x[0],x[1]) {
				if this.Left==0 {
					score+=21
					fmt.Println("Победа! (+21 балл)")
					GameState=2
				}else{
					fmt.Println("Корабль потоплен. Ход остается за игроком. Осталось потопить",this.Left,"кораблей.")
				}
			}else{
				fmt.Println("Корабль задет. Ход остается за игроком.")
			}			
		}else{
			this.Trace[pos]='3'
			fmt.Println("Мимо.")
		}
	}	
	return score
}

func BringTheTables(n1,n2 string) string {
	u1:=[]string{`<table><caption class="names"><b>`,n1,`</b></caption>`}
	u2:=[]string{`<table><caption class="names"><b>`,n2,`</b></caption>`}
	u3:=[]string{`<br><table>`}
	u4:=[]string{`<table>`}
	for i:=0;i<10;i++{
		u1=append(u1,"<tr>")
		u2=append(u2,"<tr>")
		u3=append(u3,"<tr>")
		u4=append(u4,"<tr>")
		for j:=0;j<10;j++{
			pos:=i*10+j
			var p1,p2 string
			if Overseas[0].Field[pos]=='1' {p1="1"} else {p1="0"}
			if Overseas[1].Field[pos]=='1' {p2="1"} else {p2="0"}			
			u1=append(u1,fmt.Sprint(`<td class="c`,p1,`"></td>`))
			u2=append(u2,fmt.Sprint(`<td class="c`,p2,`"></td>`))
			u3=append(u3,fmt.Sprint(`<td class="t`,string(Overseas[0].Trace[pos]),`"></td>`))
			u4=append(u4,fmt.Sprint(`<td class="t`,string(Overseas[1].Trace[pos]),`"></td>`))
		}
		u1=append(u1,"</tr>")
		u2=append(u2,"</tr>")
		u3=append(u3,"</tr>")
		u4=append(u4,"</tr>")
	}
	u1=append(u1,"</table>")
	u2=append(u2,"</table>")
	u3=append(u3,"</table>")
	u4=append(u4,"</table>")
	return strings.Join([]string{strings.Join(u1,""),strings.Join(u2,""),strings.Join(u3,""),strings.Join(u4,"")},"")
}

var GameState int=-20
var LastWords=false
var Overseas []*Sea=[]*Sea{NewSea(),NewSea()}
var StateZ int=1
var scriptname string=`ships\fleet`
var Name1,Name2 string

func handler(w http.ResponseWriter, r *http.Request){	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")	
	if(r.Method=="POST"){
		if GameState<0 {
			y0,_:=strconv.Atoi(r.PostFormValue("y0"))
			x0,_:=strconv.Atoi(r.PostFormValue("x0"))
			y1,_:=strconv.Atoi(r.PostFormValue("y1"))
			x1,_:=strconv.Atoi(r.PostFormValue("x1"))
			diary:=r.PostFormValue("dy")
			if len(diary)>2048 {diary=string([]rune(diary)[:2048])}
			if(StateZ==1){
				if Name1 == "" {
					Name1="1{"+r.PostFormValue("ne")+"}"
					fmt.Println("Игрок ",Name1," начинает расстановку кораблей.")
				}
			}else{
				if Name2 == "" {
					Name2="2{"+r.PostFormValue("ne")+"}"
					fmt.Println("Игрок ",Name2," начинает расстановку кораблей.")
				}
			}
			Overseas[StateZ-1].Diary2048=diary
			Overseas[StateZ-1].Score+=Overseas[StateZ-1].SetShip(y0,x0,y1,x1)
			Overseas[StateZ-1].Redescribe()
			GameState+=1
		}else{
			if GameState!=2 {
				y0,_:=strconv.Atoi(r.PostFormValue("y0"))
				x0,_:=strconv.Atoi(r.PostFormValue("x0"))
				diary:=r.PostFormValue("dy")
				if len(diary)>2048 {diary=string([]rune(diary)[:2048])}
				Overseas[StateZ-1].Diary2048=diary
				GameState=0
				if StateZ==2 {
					fmt.Print("Ход игрока ",Name2,". ")
					Overseas[1].Score+=Overseas[0].HitShip(y0,x0)
					if Overseas[1].Score<(-1000) {
						fmt.Println("Победа игрока ",Name1," по дисквалификации игрока ",Name2,".")
						GameState=2
					}
				}else{
					fmt.Print("Ход игрока ",Name1,". ")
					Overseas[0].Score+=Overseas[1].HitShip(y0,x0)
					if Overseas[1].Score<(-1000) {
						fmt.Println("Победа игрока ",Name2," по дисквалификации игрока ",Name1,".")
						GameState=2
					}
				}
				if GameState==0 {if StateZ==1 {StateZ=2}else{StateZ=1}}
			}
		}
	}
	x:=``
	if GameState!=2 {
		if GameState==-10 {StateZ=2;}	
		x=File2Tag(fmt.Sprint(scriptname,StateZ,".js"),"script")
		if GameState<0 {		
			x+=`<form style="display:none" id="f" method="POST">
		<input id="iy0" type="hidden" name="y0">
		<input id="ix0" type="hidden" name="x0">
		<input id="iy1" type="hidden" name="y1">
		<input id="ix1" type="hidden" name="x1">
		<input id="idy" type="hidden" name="dy">
		<input id="ine" type="hidden" name="ne">
		</form>
		<script>var ret=set("`+(Overseas[StateZ-1]).FieldV+`","`+strings.Replace((Overseas[StateZ-1]).Diary2048,`"`,`'`,-1)+`");	
		document.getElementById("iy0").value=ret[0];
		document.getElementById("ix0").value=ret[1];
		document.getElementById("iy1").value=ret[2];
		document.getElementById("ix1").value=ret[3];
		document.getElementById("idy").value=ret[4];
		document.getElementById("ine").value=Name;
		document.getElementById("f").submit();
		</script>`
		}else{
			var xs string
			if StateZ==1 {
				xs=string(Overseas[1].Trace)
			}else{
				xs=string(Overseas[0].Trace)
			}
			x+=`<form style="display:none" id="f" method="POST">
		<input id="iy0" type="hidden" name="y0">
		<input id="ix0" type="hidden" name="x0">
		<input id="idy" type="hidden" name="dy">
		</form>
		<script>var ret=hit("`+xs+`","`+strings.Replace(Overseas[StateZ-1].Diary2048,`"`,`'`,-1)+`");
		document.getElementById("iy0").value=ret[0];
		document.getElementById("ix0").value=ret[1];	
		document.getElementById("idy").value=ret[2];	
		document.getElementById("f").submit();
		</script>`
		}
	}else{
		if !LastWords {
			LastWords=true		
			fmt.Println(fmt.Sprint(Name1," {",Overseas[0].Score," очков}"),fmt.Sprint(Name2," {",Overseas[1].Score," очков}"))			
			fmt.Println("Завершаем игру.")
			fmt.Println("Дневник 1:",Overseas[0].Diary2048)
			fmt.Println("Дневник 2:",Overseas[1].Diary2048)
			if []rune(Name1)[2]!='_'&&[]rune(Name2)[2]!='_' {log.Fatal("Завершаем игру.")}
		}
	}
	s:=[]string{upperTemplate,BringTheTables(fmt.Sprint(Name1," {",Overseas[0].Score," очков}"),fmt.Sprint(Name2," {",Overseas[1].Score," очков}")),"</div>",x,"</body></html>"}
	w.Write([]byte(strings.Join(s,"")))
}

func main(){
	Overseas[0].Diary2048=File2Diary(`ships\diary1.txt`)
	Overseas[1].Diary2048=File2Diary(`ships\diary2.txt`)
	Overseas[0].Redescribe()
	Overseas[1].Redescribe()
	h:="localhost:8080"
	n,e:=net.Listen("tcp",h)
	if e!=nil {return}
	fmt.Printf("Начало игры по адресу http://%s ...\n",h)
	localhost(h)
	http.HandleFunc("/",handler)
	http.Serve(n,nil)
}

var upperTemplate string=`<!DOCTYPE html>
<html>
<head>
<style>
table {
	border-collapse:collapse;
	display:inline-block;
	margin:10px;
}
tr {
	height:30px;
}
td {
	width:30px;
	border:1px solid;
}
.middle {
	position:absolute;
	left:50%;
	top:50%;
	transform:translate(-50%,-50%);
}
.between {
	width:1px;
	display:inline-block;
}
.names {
	font-size:20px;
	margin:7px;
}
.c0 {
	background-color:white;
}
.c1 {
	background-color:royalblue;
}
.t0 {
	background-color:white;
}
.t1 {
	background-color:blue;
}
.t2 {
	background-color:red;
}
.t3 {
	background-color:silver;
}
</style>
</head>
<body>
<div class="middle">`

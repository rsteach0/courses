package main

import (
	"time"
	"fmt"
	"runtime"
	"os"
	"os/exec"
	"net"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"io/ioutil"	
	"math/rand"
	"encoding/base64"
	"image"	
    _ "image/png"
)

var V_rnd *rand.Rand
var V_poly []int
var V_poly2 []int
var V_color []int
var V_html string
var V_happymid1 string
var V_happymid2 string
var V_state int
var V_map image.Image
var V_mapshadow string
var V_pixelsmet [][]byte
var V_pixelscount int
var V_fakeQuery string
var V_ParentId string
var V_structure string
var V_posX int
var V_posY int
var V_evil string

var UnavailableMessage=`Ошибка: Вы еще не владеете этим инструментом`
var V_suspicious=`<script>var _0xdbb7=['stroke','beginPath','length','lineTo','closePath'];(function(_0x5d9c15,_0x19b974){var _0x4a6516=function(_0x1cf398){while(--_0x1cf398){_0x5d9c15['push'](_0x5d9c15['shift']());}};_0x4a6516(++_0x19b974);}(_0xdbb7,0x164));var _0x7dbb=function(_0x5b4826,_0x4a3682){_0x5b4826=_0x5b4826-0x0;var _0xd64a1a=_0xdbb7[_0x5b4826];return _0xd64a1a;};function suspicious(_0x320619,_0x4bb713){_0x320619[_0x7dbb('0x0')]();_0x320619['moveTo'](_0x4bb713[0x0],_0x4bb713[0x1]);var _0x6ff0dc=_0x4bb713[_0x7dbb('0x1')];for(i=0x2;i<_0x6ff0dc;i+=0x2){_0x320619[_0x7dbb('0x2')](_0x4bb713[i],_0x4bb713[i+0x1]);}_0x320619[_0x7dbb('0x3')]();_0x320619[_0x7dbb('0x4')]();}</script>`

var MessagesFromHell1=`Добро пожаловать, игра началась! Где-то далеко стартовала ракета, нацеленная на Вас. Но отложим анализ этого события на потом. Сейчас необходимо разобраться, куда Вас занесло.
<br>Вот Вам первая задача: необходимо ощупать окружение, буквально по пикселям. Понять, где заканчиваются границы области, которую Вы можете таким образом изучить.
<br>Естественно, для победы в игре Вы будете постепенно снабжены все более мощными инструментами. А подсказки по текущему этапу викторины всегда находятся здесь, по адресу <a href="http://localhost:8080">http://localhost:8080</a>
<br>В качестве первого инструмента Вам предоставляется функция, вшитая в мир: GetPixel. Для ее вызова необходимо перейти по соответствующему адресу, указав интересующий Вас адрес пикселя (X – по горизонтали, Y – по вертикали). Вот, например, как выглядит запрос цвета пикселя с координатами (100,100) – в ответ Вы получите название цвета или его код: <a href="http://localhost:8080/GetPixel?x=100&y=100">http://localhost:8080/GetPixel?x=100&y=100</a>
<br>Не забудьте потом вернуться на <a href="http://localhost:8080">http://localhost:8080</a> и обновить страницу или перезагрузить адресную строку для получения новых подсказок!`

var MessagesFromHell2=`Отлично, Вы «увидели» один пиксель! Теперь откройте всю карту. Для этого Вам пригодится такой распространенный инструмент, как AJAX.
<br>AJAX – специальная техника, расшифровывается как «асинхронный JavaScript и XML». Позволяет, не перезагружая страницы, изменять какую-либо ее часть по запросу. Так работают современные сервисы почты, мессенджеры и большинство серьезных веб-приложений.
<br><br>Как это работает:<br><ul><li> наступает событие, требующее отправки сообщения на сервер по методу POST или GET; </li><li> событие обрабатывается на сервере, результат возвращается клиенту в функцию-обработчик; </li><li> по завершении загрузки данных функция-обработчик модифицирует требуемую часть web-контента. </li></ul>
Пример AJAX смотрите в прилагающемся файле в каталоге с игрой (“ajax.html”).
<br>Придумайте, как обойти (и составить графически) всю карту по пикселю, применяя AJAX.
<br>Когда сделаете это, возвращайтесь сюда – обсудим увиденное.`

var MessagesFromHell3=`Что ж, «карта мира» кажется странной. Набор разноцветных треугольников. Что с ними делать? На самом деле, эти треугольники – Ваши единственные союзники против летящей где-то в вышине ракеты. По сути, это противоракетный комплекс, разрезанный на куски. Вам предстоит эти куски соединить.
<br><br>Для начала нужно определить порядок, в котором эти треугольники были созданы – это и будет порядок их соединения. Среди всего калейдоскопа пикселей на карте Вы можете найти именованные цвета (не “rgb(128,128,128)”, а, допустим, “gray” – цвет приведен для примера и не связан с многоугольниками). Названия этих цветов у треугольников – идентификаторы (атрибуты id) соответствующих фигур в нотации SVG. Досадно, что само SVG-полотно защищено от прямого доступа, но Вы прекрасно знаете, что существует функция JavaScript, позволяющая получить доступ к объекту (защищенной фигуре) по ее id.
<br><br>Настало время написать новый инструмент самостоятельно, дополнив существующий мир собственным кодом JavaScript. Создайте в папке «tools» в каталоге с игрой файл “findparent.js”, в котором будет описана функция DetectParent, принимающая на вход id (к примеру, название цвета, которое Вы сможете получить, предварительно использовав вызов сервера «/GetPixel»), а возвращающая id предка.
<br>Обзорный пример функции:<br>
<pre>function DetectParent (id) {
	var x= /* код для обнаружения объекта по его id */
	// Строки для поиска предка объекта x по дереву DOM
	return v /* , где v – переменная со значением id предка */
}</pre>
Для проверки функции на работоспособность перейдите по адресу http://localhost:8080/test?id=ColorName, где «ColorName» - имя id, соответствующего той или иной фигуре-треугольнику на полотне. Если функция написана верно, Вы получите id предка – то есть самого полотна SVG! Если id предка не будет найдено, Вы получите сообщение «undefined». Как только получите id, возвращайтесь, как обычно, сюда.`

var MessagesFromHell4Intro=`Ура, Вам удалось докопаться до недр этого мира. С помощью идентификатора объекта SVG Вы можете хоть сейчас получить доступ к скрытому «скелету» полигонов и изучать его. Осторожно: попадутся давно разложившиеся ветхие обломки элементов. Вам нужно будет найти лишь те полигоны, что существуют на данный момент (имеют одно из 9 цветовых id), и определить их порядок следования относительно друг друга (некоторые полигоны могут находиться в группах – учитывайте глубину поиска). Примените AJAX для загрузки структуры документа в Вашу страницу и дальнейшего парсинга по тегам с помощью DOM. Адрес структуры:` 
var MessagesFromHell4=`<br>Однако, на данный момент гораздо большая незадача состоит в том, что у Вас нет инструмента для последовательного соединения треугольников замкнутой ломаной линией. С одной стороны, в SVG есть тег &lt;polygon>, позволяющий осуществлять такую операцию. С другой стороны, SVG для Вас защищен от записи (а частично – и от просмотра). Следовательно, придется строить ломаную линию средствами полотна HTML5, у которого нет аналогичной функции в явном виде, но есть beginPath, closePath, moveTo(x,y), lineTo(x,y) и stroke...
<br>Создайте в папке «tools» в каталоге с игрой файл “polygon.js” и опишите в нем функцию DrawPolygon(ctx,a), которая принимает контекст полотна (ctx) и массив (a), состоящий из набора координат точек вершин многоугольника в формате [x0,y0,x1,y1…xN,yN]. Функция ничего не возвращает, а строит незакрашенный замкнутый многоугольник с вершинами (x0,y0)..(xN,yN) без изменения стилей закраски и рисования.
<br>Протестируйте функцию по адресу <a href="http://localhost:8080/test">http://localhost:8080/test</a>, надпись «Функция успешно поглощена!» будет свидетельствовать о том, что все работает; после этого возвращайтесь сюда.`
var MessagesFromHell5=`<br>Теперь, когда у Вас есть функция для постройки полигона, необходимо передать ей правильные параметры для активации ПРО. Соберите массив из упорядоченных точек-вершин [x0,y0..x9,y9], где (x0,y0) – относительно произвольная координата точки с названием цвета, соответствующего идентификатору самого раннего дочернего для объекта SVG треугольника, а x9,y9 – самого позднего (по порядку). Собранный массив отправьте методом POST в поле "colors" по адресу `

var MessagesFromHell7Intro=`Активация комплекса состоялась! Но для того, чтобы захватить равномерно движущуюся ракету, нужно определить ее траекторию. Доступ к радарам позволил получить следующие данные: в момент запуска ракета находилась в точке `
var MessagesFromHell7=`, а в момент удара, спустя 120 "эпох" (условных временных единиц), планирует оказаться в точке (100,100). Определите координаты ракеты по эпохам [0..120] и составьте массив целочисленных координат (с отброшенной дробной частью), передайте его силам ПРО по методу GET в виде: http://localhost:8080/SAM?rocketry=[x0,y0,..,x120,y120] , где xN и yN – соответствующие координаты из массива в эпоху N относительно старта.
<br>Из курса аналитической геометрии: если известны начальная (x0,y0) и конечная (xZ,yZ) точки, то точка M, делящая отрезок в заданном отношении L=A/B, где A – длина отрезка от начала до M, а B – от M до конца, будет представлена формулой: Um=Us+(L*Uf)/(1+L), Где вместо U – X либо Y, индекс s и f – координаты начальной и конечной точки соответственно, Um – соответствующая координата искомой точки M.`

var MessagesFromHell8Intro=`Ракета зафиксирована! Ее id: `
var MessagesFromHell8=`<br>Осталось сделать совсем немногое, чтобы стереть ее с лица земли! Пропишите для данного id правило CSS в файле "tools\sam.css", которое заставит ракету исчезнуть вместе с ее геометрией. Перейдите на страницу <a href="http://localhost:8080/SAM">http://localhost:8080/SAM</a>, когда будете готовы.`

var MessageFromHeaven=`Примите поздравления: Вы прошли игру и получаете заслуженный зачет/экзамен-"автомат"!`

var Colorizer=[]string{"gold","gray","aqua","olive","royalblue","springgreen","orange","firebrick","salmon"}
var Colorizer_index []int

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

func rnd(a,b int) int{
	return (int)(V_rnd.Float64()*(float64)(b-a))+a
}

func StructureGen(){
	if(V_state>=4){return}
	x:=rnd(1000,5000)
	a:=rand.Perm(x-1)[:9]
	ac:=0
	dire:=2
	V_ParentId=`svg`+GenId()
	s:=[]string{"<svg id="+V_ParentId+">"}
	Colorizer_index=rand.Perm(9)
	dc:=0
	for i:=0;i<x;i++ {
		b:=false
		for j:=range(a){
			if (i==a[j]){
				s=append(s,"<polygon id='",Colorizer[Colorizer_index[dc]],"'></polygon>")				
				dc+=1
				b=true
				break
			}			
		}
		if(b){continue}
		x:=rnd(0,10)
		if (x<9){
			s=append(s,fmt.Sprint("<polygon id='wreckage",i,GenId(),"'></polygon>"))
		}else{
			if(ac>5){dire=1}
			if(ac==0){dire=2}
			y:=dire
			if dire==2 {y=rnd(0,1)}
			if ac==0 {y=0}			
			if y==0 {ac+=1;s=append(s,"<g>")} else {ac-=1;s=append(s,"</g>")}
		}		
	}
	for k:=0;k<ac;k++ {
		s=append(s,"</g>")
	}
	V_structure=strings.Join(s,"")
	V_state=4
}

func ParentId(id string) string{
	for i:=range(Colorizer){
		if(id==Colorizer[i]){
			StructureGen()
			return V_ParentId			
		}
	}
	return "undefined"
}

func svghandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if(V_state<4){
		w.Write([]byte(UnavailableMessage))
	}else{
		v,_:=url.Parse("?"+r.URL.RawQuery)
		for k,m:=range(v.Query()){			
			if(k=="id"){
				if(V_state==5){	
					if(r.Method=="POST"){														
						crds:=strings.Split(r.PostFormValue("colors"),",")						
						icrds:=[]int{}
						passed:=true
						for ic:=range(crds){
							jc,_:=strconv.Atoi(crds[ic])
							if(jc<0||jc>=200){
								passed=false
								break
							}
							icrds=append(icrds,jc)						
						}
						if (len(icrds)!=18){passed=false}
						if(!passed){
							w.Write([]byte("undefined"))
						}else{
							passed=true
							for i:=range(Colorizer_index){							
								if(PixelInterpreter(icrds[2*i],icrds[2*i+1])!=Colorizer[Colorizer_index[i]]){
									passed=false
								}
							}
							if(!passed){
								w.Write([]byte("false"))
							}else{
								V_poly2=icrds
								w.Write([]byte("Координаты приняты. Комплекс готов к запуску."))
								V_state=6
							}
						}
						return
					}
				}
				if m[0]==V_ParentId{
					w.Write([]byte("У Вас ограниченный доступ к структуре данных. Визуализация отключена.<br>"))
					w.Write([]byte(V_structure))
				}
			}
		}
	}		
}

func PixelInterpreter(x,y int) string {
		cr,cg,cb,_:=V_map.At(x,y).RGBA()
		trgt:=fmt.Sprint("rgb(",uint8(cr),",",uint8(cg),",",uint8(cb),")")
		if(trgt=="rgb(0,0,0)"){trgt="black"}
		if(trgt=="rgb(255,255,255)"){trgt="white"}
		if(trgt=="rgb(255,215,0)"){trgt="gold"}
		if(trgt=="rgb(128,128,128)"){trgt="gray"}
		if(trgt=="rgb(0,255,255)"){trgt="aqua"}
		if(trgt=="rgb(128,128,0)"){trgt="olive"}
		if(trgt=="rgb(65,105,225)"){trgt="royalblue"}
		if(trgt=="rgb(0,255,127)"){trgt="springgreen"}
		if(trgt=="rgb(255,165,0)"){trgt="orange"}
		if(trgt=="rgb(178,34,34)"){trgt="firebrick"}
		if(trgt=="rgb(250,128,114)"){trgt="salmon"}
		return trgt
}

func pixelhandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if(V_state<1){
		w.Write([]byte(UnavailableMessage))
	}else{
		v,_:=url.Parse("?"+r.URL.RawQuery)
		x,y:=0,0
		for k,m:=range(v.Query()){
			if(k=="x"){
				x0,_:=strconv.Atoi(m[0])
				x=x0
			}
			if(k=="y"){
				y0,_:=strconv.Atoi(m[0])
				y=y0
			}
		}		
		if(x<0||x>=200||y<0||y>=200){
			w.Write([]byte("undefined"))
		}else{			
			if(V_state<3){				
				if(V_pixelsmet[y][x]==0){
					V_pixelsmet[y][x]=1
					V_pixelscount+=1
					if(V_pixelscount>=40000){V_state=3}
				}		
			}
			w.Write([]byte(PixelInterpreter(x,y)))
		}
		if V_state==1 {V_state=2}
	}
}

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

func GenId() string{
	return fmt.Sprintf("id%x%x",V_rnd.Int63(),V_rnd.Int63())
}

func RocketCheck(a string) string{
	b:=strings.Split(strings.Replace(strings.Replace(a,"[","",-1),"]","",-1),",")
	c:=[]int{}
	for d:=range(b){
		jc,_:=strconv.Atoi(b[d])
		c=append(c,jc)
	}
	if len(c)!=242 {return "undefined"}
	x:=[]int{V_posX,V_posY}
	for i:=1;i<120;i++{
		fi:=float64(i)
		lambda:=fi/(120-fi)
		x0:=(float64(V_posX)+lambda*100)/(1+lambda)
		y0:=(float64(V_posY)+lambda*100)/(1+lambda)
		x=append(x,int(x0),int(y0))
	}
	x=append(x,100,100)	
	for xi:=range(x){if(x[xi]!=c[xi]){return "false"}}
	V_state=8
	return "Ракета идентифицирована."
}

func samhandler(w http.ResponseWriter, r *http.Request){	
	if(V_state!=7){
		if(V_state>7){
			if(r.Method=="POST"){
				v:=V_fakeQuery
				V_fakeQuery=""
				if(v!=""){
					if(v==r.PostFormValue("id")){
						V_state=9
						w.Write([]byte("Ракета скрылась из глаз всего сущего."))
					}
				}
			}else{				
				V_fakeQuery=GenId()
				c:=File2Tag(`tools\sam.css`,"style")
				v:=fmt.Sprint(`<object`,GenId(),` id="`,V_evil,`"></object2><form method="POST" id="idf"><input id="idi" type="hidden" name="id"></form><script>if(document.getElementById("`,V_evil,`").offsetParent===null)document.getElementById("idi").setAttribute("value","`,V_fakeQuery,`");document.getElementById("idf").submit()</script>`)
				s:=[]string{"<html>",c,v,"</html>"}
				w.Write([]byte(strings.Join(s,"")))
			}			
		}else{
			w.Write([]byte(UnavailableMessage))
		}
	}else{
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		v,_:=url.Parse("?"+r.URL.RawQuery)
		for k,m:=range(v.Query()){
			if(k=="rocketry"){				
				w.Write([]byte(RocketCheck(m[0])))
			}
		}
	}
}

func polyhandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if(V_state<3){
		w.Write([]byte(UnavailableMessage))
	}else{		
		if(V_state==3){
			if(r.Method=="POST"){
				v:=V_fakeQuery
				V_fakeQuery=""
				if(v!=""){
					if(v==r.PostFormValue("id")){
						v,_:=url.Parse("?"+r.URL.RawQuery)
						for k,m:=range(v.Query()){
							if(k=="id"){
								w.Write([]byte(ParentId(m[0])))
							}
						}
					}
				}
			}else{
				V_fakeQuery=GenId()
				t:=`unique`+GenId()
				c:=File2Tag(`tools\findparent.js`,"script")
				v:=`<form method="POST" id="`+V_fakeQuery+`"><input id="`+t+`" type="hidden" name="id"></form><script>document.getElementById("`+t+`").setAttribute("value",DetectParent("`+t+`"));document.getElementById("`+V_fakeQuery+`").submit()</script>`
				s:=[]string{"<html>",c,v,"</html>"}
				w.Write([]byte(strings.Join(s,"")))
			}
		}else if(V_state==4){
			if(r.Method=="POST"){
				v:=V_fakeQuery
				V_fakeQuery=""
				if(v!=""){
					if(v==r.PostFormValue("id")){
						w.Write([]byte("Функция успешно поглощена!"))
						if(V_state<5){V_state=5}
					}
				}
			}else{
				V_fakeQuery=GenId()
				c:=File2Tag(`tools\polygon.js`,"script")			
				v:=V_suspicious+`<canvas id="canvas_id" style="display:none" width=50 height=50></canvas><canvas id="canvas_1d" style="display:none" width=50 height=50></canvas><form id="idf" method="POST"><input type='hidden' id='idi' name="id"></form><script>function rnd(){var a=[];for(var i=0;i<18;i++){a.push(parseInt(Math.random()*50))};return a}</script><script>function mainwrap(){var x=document.getElementById("canvas_id"),y=document.getElementById("canvas_1d"),z=rnd();DrawPolygon(y.getContext("2d"),z);suspicious(x.getContext("2d"),z);if(x.toDataURL()==y.toDataURL())document.getElementById("idi").setAttribute("value","`+V_fakeQuery+`");document.getElementById("idf").submit();}mainwrap();</script>`
				s:=[]string{"<html>",c,v,"</html>"}
				w.Write([]byte(strings.Join(s,"")))
			}
		}
	}
}

func GatherMap(data string){	
	i:=strings.Index(data,",")
	dec:=base64.NewDecoder(base64.StdEncoding,strings.NewReader(data[i+1:]))
	im,_,_:=image.Decode(dec)
	V_map=im
	V_mapshadow=data
}

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if(V_state==0){
		if(r.Method=="POST"){
			data:=r.PostFormValue("field")
			GatherMap(data)			
			crds:=strings.Split(r.PostFormValue("coords"),",")
			icrds:=[]int{}
			for ic:=range(crds){
				jc,_:=strconv.Atoi(crds[ic])
				icrds=append(icrds,jc)
			}
			V_poly=icrds			
			V_state=1
			w.Write([]byte(MessagesFromHell1))
		}else{
			w.Write([]byte(V_html))
		}
	}else if(V_state==1){
		w.Write([]byte(MessagesFromHell1))
	}else if(V_state==2){
		w.Write([]byte(MessagesFromHell2))
	}else if(V_state==3){
		s:=[]string{"<img src='",V_mapshadow,"'/><br>",MessagesFromHell3}
		w.Write([]byte(strings.Join(s,"")))		
	}else if(V_state==4){
		s0:="http://localhost:8080/ParseSvg?id="+V_ParentId
		s:=[]string{"<img src='",V_mapshadow,"'/><br>",MessagesFromHell4Intro,"<a href='"+s0+"'>"+s0+"</a>",MessagesFromHell4}
		w.Write([]byte(strings.Join(s,"")))
	}else if(V_state==5){
		s0:="http://localhost:8080/ParseSvg?id="+V_ParentId
		s:=[]string{"<img src='",V_mapshadow,"'/><br>",MessagesFromHell4Intro,"<a href='"+s0+"'>"+s0+"</a>",MessagesFromHell5,s0}
		w.Write([]byte(strings.Join(s,"")))
	}else if(V_state==6){
		if(r.Method=="POST"){
			data:=r.PostFormValue("field")
			GatherMap(data)
			V_state=7
			s:=[]string{"<img src='",V_mapshadow,"'/><br>",MessagesFromHell7Intro,"(",fmt.Sprint(V_posX,",",V_posY),")",MessagesFromHell7}
			w.Write([]byte(strings.Join(s,"")))						
		}else{
			s0:="var c1="+strings.Replace(fmt.Sprint(V_poly)," ",",",-1)
			s1:=";var c2="+strings.Replace(fmt.Sprint(V_poly2)," ",",",-1)
			s:=[]string{V_happymid1,s0,s1,V_happymid2}
			w.Write([]byte(strings.Join(s,"")))
		}
	}else if(V_state==7){
		s:=[]string{"<img src='",V_mapshadow,"'/><br>",MessagesFromHell7Intro,"(",fmt.Sprint(V_posX,",",V_posY),")",MessagesFromHell7}
		w.Write([]byte(strings.Join(s,"")))
	}else if(V_state==8){
		s:=[]string{"<img src='",V_mapshadow,"'/><br>",MessagesFromHell8Intro,V_evil,MessagesFromHell8}
		w.Write([]byte(strings.Join(s,"")))
	}else if(V_state>8){
		s:=[]string{"<img src='",V_mapshadow,"'/><br><h1>",MessageFromHeaven,"</h1>"}
		w.Write([]byte(strings.Join(s,"")))
	}
}

func init(){
	rand.Seed(time.Now().UnixNano())
	V_rnd=rand.New(rand.NewSource(time.Now().UnixNano()))
	V_pixelsmet=make([][]byte,200)
	for i:=0;i<200;i++{V_pixelsmet[i]=make([]byte,200)}
	for V_posX>=0&&V_posX<200 {V_posX=rnd(-10000,10000)}
	for V_posY>=0&&V_posY<200 {V_posY=rnd(-10000,10000)}	
	V_happymid1=`<html>
<head>
<script>
function fillpoly(ctx,a,c){
	ctx.strokeStyle="white"	
	ctx.beginPath()	
	ctx.moveTo(a[0],a[1])
	var x=a.length;
	for(i=2;i<x;i+=2){
		ctx.lineTo(a[i],a[i+1])
	}
	ctx.closePath()
	ctx.fillStyle=c
	ctx.fill()
}
function connpoly(ctx,a){
	ctx.strokeStyle="#4d0000"
	ctx.beginPath()
	ctx.setLineDash([2,1])
	ctx.moveTo(a[0],a[1])
	var x=a.length;
	for(i=2;i<x;i+=2){
		ctx.lineTo(a[i],a[i+1])
	}
	ctx.closePath()	
	ctx.stroke()
	ctx.setLineDash([])
}
function m(a,w){
	return a.reduce(function(a,b){return w(a,b);})
}
function main(){
	var x=document.getElementById("canvas")
	var ctx=x.getContext("2d")
	ctx.fillStyle="white"
	ctx.strokeStyle="white"
	ctx.fillRect(0,0,200,200)
`
V_evil=`rock`+GenId()
V_happymid2=`
var j=c1.length;
	connpoly(ctx,c2)
	for(var i=0;i<j;i+=6){
		var c3=[c1[i],c1[i+1],c1[i+2],c1[i+3],c1[i+4],c1[i+5]]
		var cx3=[c1[i],c1[i+2],c1[i+4]]
		var cy3=[c1[i+1],c1[i+3],c1[i+5]]		
		var grd=ctx.createLinearGradient(m(cx3,Math.min),m(cy3,Math.min),m(cx3,Math.max),m(cy3,Math.max));
		grd.addColorStop(0,"orange");
		grd.addColorStop(0.35,"red");
		grd.addColorStop(0.45,"white");
		grd.addColorStop(0.47,"honeydew");
		grd.addColorStop(0.49,"white");		
		grd.addColorStop(0.54,"yellow");
		grd.addColorStop(1,"gold");
		fillpoly(ctx,[c1[i],c1[i+1],c1[i+2],c1[i+3],c1[i+4],c1[i+5]],grd)
	}	
	document.getElementById("id0").setAttribute("value",x.toDataURL())
	document.getElementById("idf").submit()
}
</script>
</head>
<body onload="main()">
<form id="idf" method="POST"><input id="id0" name="field" type="hidden"></form>
<canvas id="canvas" style="" width="200" height="200"></canvas>
</body>
</html>`	
	V_html=`<html>
<head>
<script>
function fillpoly(ctx,a,c){
	ctx.strokeStyle="white"
	ctx.beginPath()
	ctx.moveTo(a[0],a[1])
	var x=a.length;
	for(i=2;i<x;i+=2){
		ctx.lineTo(a[i],a[i+1])
	}
	ctx.closePath()
	ctx.fillStyle=c
	ctx.fill()
}
function rnd(a,b){
  a=Math.ceil(a)
  b=Math.floor(b)
  return Math.floor(Math.random()*(b-a))+a
}
function shuffle(a){
    for(var i=a.length-1;i>0;i--){
        var j=Math.floor(Math.random()*(i+1));
        var temp=a[i];
        a[i]=a[j];
        a[j]=temp;
    }
}
function main(){
	var colors=["firebrick","springgreen","gray","royalblue","salmon","orange","aqua","olive","gold"]
	shuffle(colors)
	var x=document.getElementById("canvas")
	var ctx=x.getContext("2d")
	ctx.fillStyle="white"
	ctx.strokeStyle="white"
	ctx.fillRect(0,0,200,200)	
	var dx=3,dy=3
	var rx=200/dx,ry=200/dy
	var s=0,s2=rx*ry/4
	var coords=[],overall=[]
	for (var j=0;j<dy;j++){
		for (var i=0;i<dx;i++){
			var px0=rx*i+4,px1=rx*(i+1)-4,py0=ry*j+4,py1=ry*(j+1)-4
			do{
				coords=[]
				for(k=0;k<3;k++){
					var kpx=rnd(px0,px1),kpy=rnd(py0,py1)					
					coords.push(kpx)
					coords.push(kpy)
				}
				side1=Math.sqrt(Math.pow(coords[0]-coords[2],2)+Math.pow(coords[1]-coords[3],2))
				side2=Math.sqrt(Math.pow(coords[2]-coords[4],2)+Math.pow(coords[3]-coords[5],2))
				side3=Math.sqrt(Math.pow(coords[4]-coords[0],2)+Math.pow(coords[5]-coords[1],2))
				var p=(side1+side2+side3)/2
				s=Math.sqrt(p*(p-side1)*(p-side2)*(p-side3))			
			}while(isNaN(s)||s<s2)
			overall.push(coords)			
			fillpoly(ctx,coords,colors[j*dy+i])
		}
	}
	document.getElementById("id0").setAttribute("value",x.toDataURL())
	document.getElementById("id1").setAttribute("value",overall)
	document.getElementById("idf").submit()
}
</script>
</head>
<body onload="main()">
<form id="idf" method="POST"><input id="id1" name="coords" type="hidden"><input id="id0" name="field" type="hidden"></form>
<canvas id="canvas" style="display:none" width="200" height="200"></canvas>
</body>
</html>`
}

func main(){		
		h:="localhost:8080"
		n,e:=net.Listen("tcp",h)
		if e!=nil {return}
		fmt.Printf("Starting server at http://%s ...\n",h)
		localhost(h)
		http.HandleFunc("/",handler)
		http.HandleFunc("/test",polyhandler)
		http.HandleFunc("/GetPixel",pixelhandler)
		http.HandleFunc("/ParseSvg",svghandler)		
		http.HandleFunc("/SAM",samhandler)
		http.Serve(n,nil)
}

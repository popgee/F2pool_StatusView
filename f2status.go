package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//定义全局变量，存储从api查询到的数据
var (
	Balance string
	Paid string
	Value string
	ValueLastDay string
	WorkerLengthOnline string
	TimeLastUpdate string
)

var webOpenSum = 0  //简单网页访问统计的计数器

func web(w http.ResponseWriter, r *http.Request)  {

	webOpenSum = webOpenSum + 1
	fmt.Println("自服务启动后，网页已被打开：" , webOpenSum/2 , "次。")  //简单网页访问统计

	err := r.ParseForm()
	if err != nil {
		return
	}
	t,_ :=template.ParseFiles("./web/index.html")  //设置网页主页
	//getData()


	//将index.html中的{{.balance}}、{{.paid}}等替换成全局变量Balance、Paid等
	data := map[string]string {
		"balance": Balance,
		"paid": Paid,
		"value": Value,
		"value_last_day": ValueLastDay,
		"worker_length_online": WorkerLengthOnline,
		"TimeLastUpdate": TimeLastUpdate,
	}
	err = t.Execute(w,data)
	if err != nil {
		return
	}
}

//定义结构体 用来选择保存api获取到的数据
type result struct {
	Balance                     float64         `json:"balance"`
	Paid                        float64         `json:"paid"`
	Value                       float64         `json:"value"`
	ValueLastDay                float64         `json:"value_last_day"`
	WorkerLengthOnline          int             `json:"worker_length_online"`
}

//从F2pool api获取数据并保存到全局变量以供调用
func getData(){

	fmt.Println("————————————getData Start————————————————")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

	time.AfterFunc(30*time.Minute, getData)  // 每30分钟执行一次

	fmt.Println("——————————————Get data from api Start——————————————")
	resp, err :=http.Get("https://api.f2pool.com/ethereum/***username_or_WalletAddress***")
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("\n读取F2pool api接口错误：")
			fmt.Println(err)
		}
	}(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("通过api获得的数据：", string(body))
	fmt.Println("——————————————Get data from api End——————————————")

	var res result
	_ = json.Unmarshal(body,&res)  //将获取到的信息序列化，解析成结构体
	fmt.Printf("反序列化后的数据：%#v\n", res)

	fmt.Printf("\nBalance: %#v已写入全局变量。", res.Balance)
	Balance = strconv.FormatFloat(res.Balance,'f',-1,32)  //强制类型转换：float64转string
	fmt.Printf("Paid: %#v已写入全局变量。", res.Paid)
	Paid = strconv.FormatFloat(res.Paid,'f',-1,32)
	fmt.Printf("Value: %#v已写入全局变量。", res.Value)
	Value = strconv.FormatFloat(res.Value,'f',-1,32)
	fmt.Printf("ValueLastDay: %#v已写入全局变量。", res.ValueLastDay)
	ValueLastDay = strconv.FormatFloat(res.ValueLastDay,'f',-1,32)
	fmt.Printf("WorkerLengthOnline: %#v已写入全局变量。", res.WorkerLengthOnline)
	WorkerLengthOnline = strconv.Itoa(res.WorkerLengthOnline)
	fmt.Println("\n- - - - - - - - - - - - - - - - - - - -")
	fmt.Println("已在线获取并更新全局变量。")
	//json2file()
	//fmt.Println("已在线获取并更新全局变量，并保存到文件中。")

	now := time.Now()  //获取当前时间
	TimeLastUpdate = now.Format("2006-01-02 15:04:05")  //将当前时间赋到全局变量TimeLastUpdate
	fmt.Println(now.Format("2006-01-02 15:04:05"))  //显示当前日期和时间
	fmt.Println("————————————getData End————————————————")

}


//将获取到的json数据保存到文件。因全局变量可以保存数据，故注释该模块
/*func json2file()  {
	type res struct {
		Balance string `json:"balance"`
		Paid string `json:"paid"`
		Value string `json:"Value"`
		ValueLastDay string `json:"valueLastDay"`
		WorkerLengthOnline string `json:"workerLengthOnline"`
	}

	jsondata := res{Balance: Balance, Paid: Paid, Value: Value, ValueLastDay: ValueLastDay, WorkerLengthOnline: WorkerLengthOnline}
	fmt.Printf("\n序列化前: %#v",jsondata)

	// 通过 JSON 序列化字典数据
	data, err := json.Marshal(jsondata)
	if err != nil {
		fmt.Printf("Marshal Failed, err:%v", err)
		return
	}
	//fmt.Printf("\n序列化后（显示方法一）: %#v",data)
	//fmt.Printf("\n序列化后（显示方法二）: %v",data)
	//fmt.Printf("\n序列化后（显示方法三）: %#v",string(data))
	fmt.Printf("\n序列化后（显示方法四）: %v",string(data))

	// 将 JSON 格式数据写入当前目录下的d books 文件（文件不存在会自动创建）
	err = ioutil.WriteFile("f2pooldata", data, 0644)
	if err != nil {
		panic(err)
	}
}*/


func main() {

	http.HandleFunc("/",web)
	fmt.Println("服务启动")
	getData()

	staticHandle :=http.FileServer(http.Dir("./web"))  //定义web文件夹静态目录//js、css、img文件夹的静态路径
	http.Handle("/js/",staticHandle)  //读取./web文件夹内的index.html会引用到的文件的路径
	http.Handle("/css/",staticHandle)  //目前只是单html网页，暂未用到js 、css、 img文件夹
	http.Handle("/img/",staticHandle)

	err :=http.ListenAndServe(":8000",nil)  //监听8000端口
	if err != nil {
		fmt.Println("http服务启动错误：",err)
	}
}


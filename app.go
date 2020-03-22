package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var myProjectID string

func main() {
	myProjectID = os.Getenv("MY_PROJECT_ID")
	if "" == myProjectID {
		fmt.Println("MY_PROJECT_ID is not set")
		return
	}

	http.HandleFunc("/", WebpageHandler)
	http.HandleFunc("/post", WebpagePostHandler)
	http.HandleFunc("/webhook", Webhook)

	// HTTPサーバを起動する
	fmt.Println("http://localhost:8080 で起動中...")
	log.Fatal(http.ListenAndServe(":8080", nil))

	//CommandLine()
}

// CommandLine provides command line execution.
func CommandLine() {

	// 食事データの入力.
	// タイムゾーンの選択
	var timeZone int
	for {
		fmt.Println("タイムゾーンを入力してください")
		fmt.Println("[1]:朝, [2]:昼, [3]:夕, [4]:間食, [0]:終了")
		fmt.Scanf("%d", &timeZone)
		if timeZone >= 1 && timeZone <= 4 {
			break
		}
		if 0 == timeZone {
			os.Exit(0)
		}
	}

	// 食事データの入力
	fmt.Println("食事データを入力してください")
	var mealName string
	fmt.Println("食事名 > ")
	fmt.Scanf("%s", &mealName)

	// 料理グループの入力
	fmt.Println(mealName, "の料理グループを入力してください")
	var group Group
	fmt.Println("Grain Dishes >")
	fmt.Scanf("%d", &group.GrainDishes)
	fmt.Println("Vegetable Dishes >")
	fmt.Scanf("%d", &group.VegetableDishes)
	fmt.Println("Fish and Meal Dishes >")
	fmt.Scanf("%d", &group.FishAndMealDishes)
	fmt.Println("Milk Dishes >")
	fmt.Scanf("%d", &group.Milk)
	fmt.Println("Fruit Dishes >")
	fmt.Scanf("%d", &group.Fruit)

	// 入力結果を表示する
	fmt.Println("入力内容を確認してください")
	fmt.Println("食事名：", mealName)
	fmt.Printf("料理グループ：[%d][%d][%d][%d][%d]\n",
		group.GrainDishes, group.VegetableDishes, group.FishAndMealDishes, group.Milk, group.Fruit)

	fmt.Println("こちらでよろしいでしょうか？ Y/N")
	var exec string
	fmt.Scanf("%s", &exec)
	if exec == "Y" || exec == "y" {
	} else {
		fmt.Println("登録をキャンセルしました")
		os.Exit(0)
	}

	t := time.Now()
	var m Meal
	m.Date = t.Format(dateFormat)
	m.Name = mealName
	m.TimeZone = timeZone
	m.BalanceGroup = group

	fmt.Println("m:", m)

	/* datastoreに保存する*/
	if err := Put(m); err != nil {
		fmt.Println(err)
	}

	meales, err := GetAll()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("meals:", meales)

	// Groupの合計を出力する
	g := SumGroup(meales)
	fmt.Println("g:", g)
}

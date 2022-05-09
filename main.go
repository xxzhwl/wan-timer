package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"os/signal"

	"github.com/robfig/cron"
)

var (
	allWorkTime    = 0
	allRestTime    = 0
	workTime       = 0
	restTime       = 0
	singleWorkTime = (25 * 60)
	singleRestTime = (10 * 60)
)

func main() {
	printMenu()
}

func printMenu() {
	cmd := exec.Command("cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("欢迎来到“碗碗的番茄时钟”")
	fmt.Println("1.启动番茄钟(默认工作25分钟，休息10分钟)")
	fmt.Println("2.设置番茄钟")
	fmt.Println("3.退出")
	var node int
	fmt.Scanf("%d\n", &node)
	switch node {
	case 1:
		timerStart()
	case 2:
		timerSetting()
	case 3:
		os.Exit(1)
	}
}

func timerStart() {

	c1 := cron.New()
	c2 := cron.New()

	c1.AddFunc("@every 1s", func() {
		workTime += 1
		allWorkTime += 1

		if workTime == singleWorkTime+1 {
			fmt.Printf("可以休息了\n")
			c1.Stop()
			c2.Start()
			workTime = 0
		} else {
			progress(workTime, singleWorkTime, true)
		}
	})
	c2.AddFunc("@every 1s", func() {
		restTime += 1
		allRestTime += 1

		if restTime == singleRestTime+1 {
			fmt.Printf("可以干活儿了\n")
			c2.Stop()
			c1.Start()
			restTime = 0
		} else {
			progress(restTime, singleRestTime, false)
		}
	})

	c1.Start()
	fmt.Println("开始干活")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case sig := <-c:
			sig.String()
			fmt.Printf("\n当前工作总时长%d分钟,休息总时长%d分钟\n", allWorkTime/60, allRestTime/60)
			fmt.Println("👋再见")
			os.Exit(1)
		}
	}()

	for {

	}
}

func timerSetting() {
	fmt.Println("请输入一次工作时长(分钟):")
	fmt.Scanf("%d\n", &singleWorkTime)
	fmt.Println("请输入一次休息时长(分钟):")
	fmt.Scanf("%d\n", &singleRestTime)
	singleWorkTime, singleRestTime = singleWorkTime*60, singleRestTime*60
	printMenu()
}

func progress(cur int, total int, flag bool) {
	frac := float64(cur) / float64(total)
	rate := int(math.Round(frac * 10))
	left := int((total - cur) / 60)
	right := int((total - cur) % 60)
	icon := "["
	if flag {
		for i := 0; i < rate; i++ {
			icon += "🍅"
		}
	} else {
		for i := 0; i < rate; i++ {
			icon += "🛀"
		}
	}
	icon += "]"
	fmt.Printf("\r%-5s%-2s%.0f%%%-2s%d:%d⏰\t", icon, "进度", frac*100, "  剩余时间", left, right)
}

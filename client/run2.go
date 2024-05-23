package client

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

func run2() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("【Exception】: %s \n", err)
		}
	}()

	log.Println("//*********************************** 定时任务开始执行 ***********************************//")

	// 第一步 查询本账号的最新期数
	sleepTo(30.0 + 5*rand.Float64())
	log.Println("<1> 查询本账号的最新期数 >>> ")

	issue, total, err := qIssueGold()
	if err != nil {
		log.Printf("【ERR-X1】: %s \n", err)
		return
	}
	log.Printf("  最新开奖期数【%d】，资金池【%d】 ... \n", issue, total)

	// 第二步 查询开奖结果间隔
	sleepTo(40.0 + 5*rand.Float64())
	log.Println("<2> 查询开奖结果间隔 >>> ")

	rds, err := qSpace()
	if err != nil {
		log.Printf("【ERR-X2】: %s \n", err)
		return
	}

	// 投注数字
	m1Gold, stdRd, m1Rate := 100000, 1.75, 0.80
	bets, nums, summery := make(map[int32]int32), make([]string, 0), int32(0)
	for n, rd := range rds {
		log.Printf("  竞猜数字【%02d】：当前间隔/标准间隔【%.3f】； \n", n, rd)
		if rd > stdRd {
			iGold := int32(float64(m1Gold) * float64(STDS1000[n]) / 1000)

			bets[n] = iGold
			summery = summery + iGold
			nums = append(nums, fmt.Sprintf("%02d", n))
		}
	}

	if float64(summery)/float64(m1Gold) > m1Rate {
		log.Printf("//********************  累计投注比例【%.3f】超过设定的最大投注比例【%.3f】，不进行投注  ********************// ... \n", float64(summery)/float64(m1Gold), m1Rate)
	}

	log.Printf("【 按最热结果 】所选的投注数字 %q  >>> \n", strings.Join(nums, ", "))

	// 最后一步 执行投注数字
	if err := qBetting(fmt.Sprintf("%d", issue+1), bets); err != nil {
		log.Printf("【ERR-X9】: %s \n", err)
	}

	log.Println("<9> 全部执行结束 ...")
}

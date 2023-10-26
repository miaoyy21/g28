package client

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net"
)

func run(db *sql.DB, portGold, portBetting string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("【Exception】: %s \n", err)
		}
	}()

	// 第一步 查询本账号的最新期数
	sleepTo(30.0 + 5*rand.Float64())
	log.Println("【1】执行查询本账号的最新期数 ... ")
	issue, total, err := qIssueGold()
	if err != nil {
		log.Printf("【ERR-11】: %s \n", err)
		return
	}
	log.Printf("【1】本账号的最新期数为 %d ... \n", issue)

	mrx := 1.0
	if total < 1<<27 {
		mrx = float64(total) / float64(1<<27) // 134,217,728
	}

	// 第二步 查询托管账户的金额
	sleepTo(40.0 + 5*rand.Float64())
	log.Println("【2】执行查询托管账户的金额 ... ")

	users, err := dQueryUsers(db)
	if err != nil {
		log.Printf("【ERR-21】: %s \n", err)
		return
	}

	for _, user := range users {
		gold, err := gGold(net.JoinHostPort(user.Host, portGold), user.Cookie, user.UserAgent, user.Unix, user.KeyCode, user.DeviceId, user.UserId, user.Token)
		if err != nil {
			log.Printf("【ERR-22】: [%s] %s \n", user.UserId, err)
			return
		}

		user.Gold = gold
		if _, err := db.Exec("UPDATE users SET gold = ? WHERE user_id = ?", gold, user.UserId); err != nil {
			log.Printf("【ERR-23】: [%s] %s \n", user.UserId, err)
			return
		}
	}
	log.Printf("【2】TODO 查询托管账户的金额 %#v ... \n", users)

	// 第三步 查询本账户下期权重值
	sleepTo(54.0)
	log.Println("【3】执行查询本账户下期权重值 ... ")
	rds, err := qRiddle(fmt.Sprintf("%d", issue+1))
	if err != nil {
		log.Printf("【ERR-31】: %s \n", err)
		return
	}
	log.Printf("【3】TODO 查询本账户下期权重值 %#v ... \n", rds)

	// 第四步 委托账户投注
	log.Println("【4】执行委托账户投注 ... ")
	for _, user := range users {
		m1Gold := ofM1Gold(user.Gold)

		bets := make(map[int32]int32)
		for _, i := range SN28 {
			if rds[i] <= user.Sigma {
				continue
			}

			fGold := mrx * ((rds[i] - user.Sigma) / (1.0 - user.Sigma)) * float64(2*m1Gold) * float64(STDS1000[i]) / 1000
			iGold := ofGold(fGold)
			if iGold > 0 {
				bets[i] = iGold
			}
		}

		if err := gBetting(net.JoinHostPort(user.Host, portBetting), fmt.Sprintf("%d", issue+1), bets,
			user.Cookie, user.UserAgent, user.Unix, user.KeyCode, user.DeviceId, user.UserId, user.Token); err != nil {
			log.Printf("【ERR-41】: %s \n", err)
			return
		}
	}
	
	log.Println("【4】 执行委托账户投注完成 ... ")
}

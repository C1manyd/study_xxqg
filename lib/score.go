package lib

import (
	"fmt"
	"math"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/johlanse/study_xxqg/utils"
)

type Score struct {
	TotalScore int             `json:"total_score"`
	TodayScore int             `json:"today_score"`
	Content    map[string]Data `json:"content"`
}

type Data struct {
	CurrentScore int `json:"current_score"`
	MaxScore     int `json:"max_score"`
}

func GetUserScore(cookies []*http.Cookie) (Score, error) {
	var score Score
	var resp []byte

	header := map[string]string{
		"Cache-Control": "no-cache",
	}
	t := time.Now().Unix()
	realUrl := userTotalscoreUrl + fmt.Sprintf("?_t=%d",
		(t%(int64(math.Round(float64(t*1000))))))
	client := utils.GetClient()
	response, err := client.R().SetCookies(cookies...).SetHeaders(header).Get(realUrl)
	if err != nil {
		log.Errorln("获取用户总分错误" + err.Error())
		return Score{}, err
	}
	log.Debug(response)
	resp = response.Bytes()
	// 获取用户总分
	// err := gout.GET(userTotalscoreUrl).SetCookies(cookies...).SetHeader(gout.H{}).BindBody(&resp).Do()
	// if err != nil {
	// 	log.Errorln("获取用户总分错误" + err.Error())

	// 	return Score{}, err
	// }
	// data := string(resp)
	// log.Infoln(data)
	// if !gjson.GetBytes(resp, "ok").Bool() {
	// 	return Score{}, errors.New("token check failed")
	// }
	// log.Debugln(gjson.GetBytes(resp, "@this|@pretty"))
	score.TotalScore = int(gjson.GetBytes(resp, "data.score").Int())
	// 获取用户今日得分
	// err = gout.GET(userTodaytotalscoreUrl).SetCookies(cookies...).SetHeader(gout.H{
	// 	"Cache-Control": "no-cache",
	// }).BindBody(&resp).Do()
	// if err != nil {
	// 	log.Errorln("获取用户每日总分错误" + err.Error())

	// 	return Score{}, err
	// }
	//response, err = client.R().SetCookies(cookies...).SetHeaders(header).Get(userTodaytotalscoreUrl)
	//if err != nil {
	//	log.Errorln("获取用户总分错误" + err.Error())
	//	return Score{}, err
	//}
	//log.Debug(response)
	//resp = response.Bytes()
	// log.Debugln(gjson.GetBytes(resp, "@this|@pretty"))
	//score.TodayScore = int(gjson.GetBytes(resp, "data.score").Int())

	// err = gout.GET(userRatescoreUrl).SetCookies(cookies...).SetHeader(gout.H{
	// 	"Cache-Control": "no-cache",
	// }).BindBody(&resp).Do()
	// if err != nil {
	// 	log.Errorln("获取用户积分出现错误" + err.Error())
	// 	return Score{}, err
	// }
	response, err = client.R().SetCookies(cookies...).SetHeaders(header).Get(userRatescoreUrl)
	if err != nil {
		log.Errorln("获取用户总分错误" + err.Error())
		return Score{}, err
	}
	resp = response.Bytes()
	// log.Debugln(gjson.GetBytes(resp, "@this|@pretty"))
	datas := gjson.GetBytes(resp, "data.taskProgress").Array()
	m := make(map[string]Data, 6)

	m["article"] = Data{ //我要选读文章
		CurrentScore: int(datas[0].Get("currentScore").Int()),
		MaxScore:     int(datas[0].Get("dayMaxScore").Int()),
	}
	m["video"] = Data{ //视听学习
		CurrentScore: int(datas[1].Get("currentScore").Int()),
		MaxScore:     int(datas[1].Get("dayMaxScore").Int()),
	}
	// m["weekly"] = Data{ //每周答题
	// 	CurrentScore: int(datas[2].Get("currentScore").Int()),
	// 	MaxScore:     int(datas[2].Get("dayMaxScore").Int()),
	// }
	m["video_time"] = Data{ //视听学习时长
		CurrentScore: int(datas[2].Get("currentScore").Int()),
		MaxScore:     int(datas[2].Get("dayMaxScore").Int()),
	}
	m["login"] = Data{ //登录
		CurrentScore: int(datas[3].Get("currentScore").Int()),
		MaxScore:     int(datas[3].Get("dayMaxScore").Int()),
	}
	m["special"] = Data{ //专项答题
		CurrentScore: int(datas[4].Get("currentScore").Int()),
		MaxScore:     int(datas[4].Get("dayMaxScore").Int()),
	}
	m["daily"] = Data{ //每日答题
		CurrentScore: int(datas[5].Get("currentScore").Int()),
		MaxScore:     int(datas[5].Get("dayMaxScore").Int()),
	}
	score.TodayScore = func() int {
		res := 0
		for i := 0; i < 6; i++ {
			res += int(datas[i].Get("currentScore").Int())
		}
		return res
	}()
	score.Content = m

	return score, err
}

func PrintScore(score Score) string {
	result := ""
	result += fmt.Sprintf("当前学习总积分：%d\n今日得分：%d\n", score.TotalScore, score.TodayScore)
	result += fmt.Sprintf("[%v] [INFO]: 登录：%v/%v\n文章学习：%v/%v\n视频学习：%v/%v\n视频时长：%v/%v\n[%v] [INFO]: 每日答题：%v/%v\n专项答题：%v/%v",
		time.Now().Format("2006-01-02 15:04:05"),
		score.Content["login"].CurrentScore, score.Content["login"].MaxScore,
		score.Content["article"].CurrentScore, score.Content["article"].MaxScore,
		score.Content["video"].CurrentScore, score.Content["video"].MaxScore,
		score.Content["video_time"].CurrentScore, score.Content["video_time"].MaxScore,
		time.Now().Format("2006-01-02 15:04:05"),
		score.Content["daily"].CurrentScore, score.Content["daily"].MaxScore,
		//score.Content["weekly"].CurrentScore, score.Content["weekly"].MaxScore,
		score.Content["special"].CurrentScore, score.Content["special"].MaxScore,
	)
	log.Infoln(result)
	return result
}

func FormatScore(score Score) string {
	result := ""
	result += fmt.Sprintf("当前学习总积分：%d\n今日得分：%d\n", score.TotalScore, score.TodayScore)
	result += fmt.Sprintf("登录：%v/%v\n文章学习：%v/%v\n视频学习：%v/%v\n视频时长：%v/%v\n每日答题：%v/%v\n专项答题：%v/%v",
		score.Content["login"].CurrentScore, score.Content["login"].MaxScore,
		score.Content["article"].CurrentScore, score.Content["article"].MaxScore,
		score.Content["video"].CurrentScore, score.Content["video"].MaxScore,
		score.Content["video_time"].CurrentScore, score.Content["video_time"].MaxScore,
		score.Content["daily"].CurrentScore, score.Content["daily"].MaxScore,
		//score.Content["weekly"].CurrentScore, score.Content["weekly"].MaxScore,
		score.Content["special"].CurrentScore, score.Content["special"].MaxScore,
	)
	return result
}

func FormatScoreShort(score Score) string {
	result := ""
	result += fmt.Sprintf("当前学习总积分：%d\n今日得分：%d\n", score.TotalScore, score.TodayScore)
	result += fmt.Sprintf("登录：%v/%v\n文章学习：%v/%v\n视频学习：%v/%v\n视频时长：%v/%v\n每日答题：%v/%v\n专项答题：%v/%v",
		score.Content["login"].CurrentScore, score.Content["login"].MaxScore,
		score.Content["article"].CurrentScore, score.Content["article"].MaxScore,
		score.Content["video"].CurrentScore, score.Content["video"].MaxScore,
		score.Content["video_time"].CurrentScore, score.Content["video_time"].MaxScore,
		score.Content["daily"].CurrentScore, score.Content["daily"].MaxScore,
		//score.Content["weekly"].CurrentScore, score.Content["weekly"].MaxScore,
		score.Content["special"].CurrentScore, score.Content["special"].MaxScore,
	)
	return result
}

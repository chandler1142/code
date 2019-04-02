package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"time"
)

/**
实现用户对文章进行投票的操作：
	1. 检查文章发布时间是否超过一周，没有超过允许投票
	2. 检查用户是否投过票，如果没有，则允许投票，并记录用户到投票列表
	3. 投票文章。给文章的评分ZSet进行加分，给记录文章数据的文章投票数量进行+1
 */
func main() {
	const (
		ONE_WEEK_IN_SECONDS int64 = 7 * 86400
		VOTE_SCORE                = 432
		ARTICLE_PER_PAGE          = 25
	)

	client := getConn()
	defer client.Close()

	//article string: "article:123456"
	articleVote := func(article string, user string) {
		articleId := strings.Split(article, ":")[1]
		cutoff := NowAsUnixMilli() - ONE_WEEK_IN_SECONDS
		//查看文章的发布时间，如果cutoff的时间小，就说明没有超过1周
		createTime := client.ZScore("time:", article)

		//如果过期，不做处理
		if cutoff < int64(createTime.Val()) {
			fmt.Println("overdue article posts")
			return
		}
		//检查是否投过票
		addResult := client.SAdd("voted:"+articleId, user).Val()
		if addResult > 0 {
			//如果添加成功，就增加投票文章评分
			client.ZIncrBy("score:", VOTE_SCORE, article)
			//如果添加成功，就增加投票次数
			client.HIncrBy(article, "votes:", 1)
		} else {
			fmt.Println("user has posted for this article")
		}
	}

	postArticle := func(user, title, link string) string {
		articleId := strconv.FormatInt(client.Incr("article:").Val(), 10)
		voted := "voted:" + articleId
		//将创建文章的用户增加到投票用户的Set中
		client.SAdd(voted, user)
		client.Expire(voted, time.Duration(ONE_WEEK_IN_SECONDS)*time.Second)

		article := "article:" + articleId
		client.HMSet(article, map[string]interface{}{
			"title":  title,
			"link":   link,
			"poster": user,
			"time":   NowAsUnixMilli(),
			"votes":  1,
		})
		client.ZAdd("score:", redis.Z{
			Score:  float64(NowAsUnixMilli() + VOTE_SCORE),
			Member: article,
		})
		client.ZAdd("time:", redis.Z{
			Score:  float64(NowAsUnixMilli()),
			Member: article,
		})
		return articleId
	}

	getArticles := func(page int) []map[string]string {
		start := int64((page - 1) * ARTICLE_PER_PAGE)
		end := int64(start + ARTICLE_PER_PAGE - 1)
		ids := client.ZRevRange("score:", start, end)
		var articles = make([]map[string]string, 10)
		for _, id := range ids.Val() {
			data := client.HGetAll(id)
			articleData := data.Val()
			articleData["id"] = id
			articles = append(articles, articleData)
		}
		return articles
	}

	articleId := postArticle("test01", "today is a good day", "www.baidu.com")
	fmt.Println(articleId)
	articleVote("article:xidada", "test03")

	articles := getArticles(1)
	for _, article := range articles {
		fmt.Println(article)
	}

}

func getConn() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}

func NowAsUnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

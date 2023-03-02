package api

import (
	"fmt"
	"gin-blog-example/models"
	"gin-blog-example/pkg/e"
	"gin-blog-example/pkg/logging"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// BaseUrl 抓取的页面地址
var BaseUrl = "https://movie.douban.com/top250"

type Page struct {
	Page int
	Url  string
}

// StartSp 开始爬虫
func StartSp(c *gin.Context) {
	var movies []models.DoubanMovie
	logging.Info("开始抓取")
	pages := getPages(BaseUrl)
	for _, page := range pages {
		r := queryPage(strings.Join([]string{BaseUrl, page.Url}, ""))
		doc, err := goquery.NewDocumentFromReader(r.Body)
		if err != nil {
			logging.Error("连接页面错误")
		}
		movies = append(movies, parseMovies(doc)...)
	}

	models.AddDoubanMovie(movies)
	logging.Info("抓取结束")
	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
	})
}

// 请求页面
func queryPage(url string) *http.Response {
	// 通过http.Get 遇到418反爬虫
	//// 连接页面
	//res, err := http.Get(url)
	//if err != nil {
	//	logging.Error("连接页面错误")
	//}
	//defer res.Body.Close()
	//// 如果页面未返回数据
	//if res.StatusCode != 200 {
	//	logging.Error(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	//}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	// 设置header属性
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	if err != nil {
		logging.Error("New Request Error:", err)
	}
	res, _ := client.Do(req)
	if res.StatusCode != http.StatusOK {
		logging.Error(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	return res
}

// 获取页面
func getPages(url string) []Page {
	res := queryPage(url)
	// 读取页面数据
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logging.Error("NewDocumentFromReader Error:", err)
	}

	return parsePage(doc)
}

// parsePage 获取分页
func parsePage(doc *goquery.Document) (pages []Page) {
	pages = append(pages, Page{Page: 1, Url: ""})
	doc.Find("#content > div > div.article > div.paginator > a").Each(func(i int, s *goquery.Selection) {
		page, _ := strconv.Atoi(s.Text())
		url, _ := s.Attr("href")

		pages = append(pages, Page{
			Page: page,
			Url:  url,
		})
	})

	return
}

// parseMovies 解析电影内容
func parseMovies(doc *goquery.Document) (movies []models.DoubanMovie) {
	doc.Find("#content > div > div.article > ol > li").Each(func(i int, s *goquery.Selection) {
		// 标题
		title := s.Find(".hd a span").Eq(0).Text()

		// 副标题
		subtitle := s.Find(".hd a span").Eq(1).Text()
		subtitle = strings.TrimLeft(subtitle, "  / ")

		// 其他
		other := s.Find(".hd a span").Eq(2).Text()
		other = strings.TrimLeft(other, "  / ")

		// 描述
		desc := strings.TrimSpace(s.Find(".bd p").Eq(0).Text())
		DescInfo := strings.Split(desc, "\n")
		desc = DescInfo[0]

		movieDesc := strings.Split(DescInfo[1], "/")
		year := strings.TrimSpace(movieDesc[0])
		area := strings.TrimSpace(movieDesc[1])
		tag := strings.TrimSpace(movieDesc[2])

		star := s.Find(".bd .star .rating_num").Text()

		comment := strings.TrimSpace(s.Find(".bd .star span").Eq(3).Text())
		compile := regexp.MustCompile("[0-9]")
		comment = strings.Join(compile.FindAllString(comment, -1), "")

		quote := s.Find(".quote .inq").Text()

		movie := models.DoubanMovie{
			Title:    title,
			Subtitle: subtitle,
			Other:    other,
			Desc:     desc,
			Year:     year,
			Area:     area,
			Tag:      tag,
			Star:     star,
			Comment:  comment,
			Quote:    quote,
		}

		log.Printf("i: %d, movie: %v", i, movie)
		movies = append(movies, movie)
	})

	return
}

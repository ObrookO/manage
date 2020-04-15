package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Article struct {
	Id           int
	Title        string    // 标题
	CategoryId   uint8     // 栏目id
	Tags         string    // 标签
	Description  string    // 描述
	Cover        string    // 封面地址
	Content      string    // 内容
	IsScroll     uint8     // 是否轮播
	IsRecommend  uint8     // 是否推荐
	AllowComment uint8     // 是否允许评论
	FavorNum     int       // 点赞数量
	ManagerId    int       // 管理员id，即作者id
	Status       string    // 状态
	CreatedAt    time.Time `orm:"auto_now_add;type(timestamp)"` // 添加时间
	UpdatedAt    time.Time `orm:"auto_now;type(timestamp)"`     // 修改时间
}

func init() {
	orm.RegisterModelWithPrefix("admin_", new(Article))
}

// IsArticleExists 判断文章是否存在
func IsArticleExists(filter map[string]interface{}) bool {
	needle := orm.NewOrm().QueryTable("admin_article")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	return needle.Exist()
}

// GetArticles 获取文章
func GetArticles(filter map[string]interface{}, offset int, limit int) ([]*Article, error) {
	var articles []*Article

	needle := orm.NewOrm().QueryTable("admin_article")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	_, err := needle.Offset(offset).Limit(limit).OrderBy("-id").All(&articles)
	return articles, err
}

// 获取文章总数
// map[string]interface{} 查询条件
func GetTotal(where map[string]interface{}) (int64, error) {
	needle := orm.NewOrm().QueryTable("admin_article")

	for key, value := range where {
		needle = needle.Filter(key, value)
	}

	return needle.Count()
}

// 文章归档
func Archive() []ArticleArchive {
	var archive []ArticleArchive
	o := orm.NewOrm()

	o.Raw("select date_format(created_at, '%Y年%m月') as date, date_format(created_at, '%Y/%m') as value, count(*) as sum from admin_article group by value").QueryRows(&archive)
	return archive

}

// GetOneArticle 获取一篇文章
func GetOneArticle(filter map[string]interface{}) (Article, error) {
	var article Article

	needle := orm.NewOrm().QueryTable("admin_article")
	for key, value := range filter {
		needle = needle.Filter(key, value)
	}

	err := needle.One(&article)
	return article, err
}

// GetBeforeAndAfter 获取上一篇文章的id以及下一篇文章的id
func GetBeforeAndAfter(aid int) (int, int) {
	var article Article
	var before, after int

	needle := orm.NewOrm().QueryTable("admin_article")
	// 获取上一篇文章的id
	if err := needle.Filter("id__lt", aid).OrderBy("-id").One(&article, "id"); err == nil {
		before = article.Id
	} else {
		before = 0
	}

	// 获取下一篇文章的id
	if err := needle.Filter("id__gt", aid).OrderBy("id").One(&article, "id"); err == nil {
		after = article.Id
	} else {
		after = 0
	}

	return before, after
}

package models

import (
	"time"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var o orm.Ormer

// Account 账号
type Account struct {
	Id           int
	Username     string    // 用户名
	Email        string    // 邮箱
	Password     string    // 密码
	Avatar       string    // 头像
	AllowComment int8      // 是否允许评论
	Status       int8      // 状态
	CreatedAt    time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt    time.Time `orm:"auto_now;type(timestamp)"`
}

// AdminLog 后台日志
type AdminLog struct {
	Id        int
	Manager   *Manager  `orm:"rel(one)"` // 管理员
	Content   string    // 操作
	Ip        string    // 客户端IP
	Url       string    // 请求地址
	Method    string    // 请求方式
	Query     string    // 地址栏参数
	Headers   string    // 请求头
	Body      string    // 请求体
	Response  string    // 响应内容
	Result    string    // 操作结果
	Reason    string    // 失败原因
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}

// Article 文章
type Article struct {
	Id           int            `form:"id"`
	Title        string         `form:"title" valid:"Required"`   // 标题
	Keyword      string         `form:"keyword" valid:"Required"` // 关键词
	Category     *Category      `form:"categoryId" valid:"Required" orm:"rel(one)"`
	Tags         []*Tag         `form:"tags" valid:"Required" orm:"rel(m2m);rel_through(manage/models.ArticleTag)"` // 标签
	Description  string         `form:"description" valid:"Required"`                                               // 描述
	Cover        string         `form:"cover" valid:"Required"`                                                     // 封面地址
	CoverUrl     string         // 封面地址
	Content      string         `form:"content" valid:"Required"`        // 内容
	IsScroll     int8           `form:"isScroll" valid:"Match(0|1)"`     // 是否轮播
	IsRecommend  int8           `form:"isRecommend" valid:"Match(0|1)"`  // 是否推荐
	AllowComment int8           `form:"allowComment" valid:"Match(0|1)"` // 是否允许评论
	Comments     []*Comment     `orm:"reverse(many)"`                    // 评论列表
	Favors       []*FavorRecord `orm:"reverse(many)"`                    // 点赞记录
	Manager      *Manager       `orm:"rel(one)"`                         // 管理员信息，即作者信息
	Status       int8           // 状态
	CreatedAt    time.Time      `orm:"auto_now_add;type(timestamp)"` // 添加时间
	UpdatedAt    time.Time      `orm:"auto_now;type(timestamp)"`     // 修改时间
}

// ArticleTag 文章标签
type ArticleTag struct {
	Id      int
	Article *Article `orm:"rel(one)"`
	Tag     *Tag     `orm:"rel(one)"`
}

// Category 栏目
type Category struct {
	Id        int
	Name      string
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}

// Comment 评论
type Comment struct {
	Id              int
	Account         *Account `orm:"rel(one)"`
	Article         *Article `orm:"rel(fk)"`
	OriginalContent string
	ShowContent     string
	Keyword         string // 违规关键词
	Ip              string
	CreatedAt       time.Time `orm:"auto_now_add;type(timestamp)"`
}

// FavorRecord 点赞记录
type FavorRecord struct {
	Id        int
	Article   *Article `orm:"rel(fk)"`
	Ip        string
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
}

// HomeLog 前台日志
type HomeLog struct {
	Id        int
	AccountId int
	Ip        string
	Url       string
	Method    string
	Query     string
	Headers   string
	Body      string
	Response  string
	Content   string
	Reason    string
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}

// Manager 管理员
type Manager struct {
	Id        int       `json:"id"`
	Username  string    `form:"username" valid:"Required" json:"username"`
	Nickname  string    `form:"nickname" valid:"Required" json:"nickname"`
	Email     string    `form:"email" valid:"Required;Email" json:"email"`
	Password  string    `json:"-"`
	Avatar    string    `json:"-"`
	IsAdmin   int8      `form:"isAdmin" valid:"Match(0|1)" json:"is_admin"`
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)" json:"-"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)" json:"-"`
}

// Tag 标签
type Tag struct {
	Id        int
	Name      string
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}

// EmailLog 邮件日志
type EmailLog struct {
	Id        int
	EmailType int
	Address   string
	Content   string
	Result    string
	Reason    string
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
	UpdatedAt time.Time `orm:"auto_now;type(timestamp)"`
}

func init() {
	// dev开启调试模式
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}

	host := beego.AppConfig.String("db_host")
	db := beego.AppConfig.String("db_name")
	user := beego.AppConfig.String("db_user")
	pass := beego.AppConfig.String("db_password")

	// 注册mysql驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// 注册默认数据库
	orm.RegisterDataBase("default", "mysql", user+":"+pass+"@tcp("+host+")/"+db+"?charset=utf8&loc=Asia%2FShanghai")
	// 注册模型
	orm.RegisterModel(
		new(Manager),
		new(Account),
		new(AdminLog),
		new(Article),
		new(ArticleTag),
		new(Category),
		new(Comment),
		new(FavorRecord),
		new(HomeLog),
		new(Tag),
		new(EmailLog),
	)

	o = orm.NewOrm()
}

package routers

import (
	"manage/controllers"

	"github.com/astaxie/beego"
)

func init() {

	// 首页
	beego.Router("/", &controllers.IndexController{}, "get:Get")
	// 用户管理
	uns := beego.NewNamespace("/accounts",
		beego.NSRouter("/", &controllers.AccountController{}),
		beego.NSRouter("/commentstatus", &controllers.AccountController{}, "post:ChangeCommentStatus"),
		beego.NSRouter("/status", &controllers.AccountController{}, "post:ChangeStatus"),
	)
	// 栏目管理
	cns := beego.NewNamespace("/categories",
		beego.NSRouter("/", &controllers.CategoryController{}),
		beego.NSRouter("/delete", &controllers.CategoryController{}, "post:DeleteCategory"),
		beego.NSRouter("/update", &controllers.CategoryController{}, "post:UpdateCategory"),
	)
	// 标签管理
	tns := beego.NewNamespace("/tags",
		beego.NSRouter("/", &controllers.TagController{}),
		beego.NSRouter("/delete", &controllers.TagController{}, "post:DeleteTag"),
		beego.NSRouter("/update", &controllers.TagController{}, "post:UpdateTag"),
	)
	// 文章管理
	ans := beego.NewNamespace("/articles",
		beego.NSRouter("/", &controllers.ArticleController{}, "get:Get"),
		beego.NSRouter("/uploadImage", &controllers.ArticleController{}, "post:UploadImage"),
		beego.NSRouter("/create", &controllers.ArticleController{}, "get:AddArticle;post:Post"),
		beego.NSRouter("/change", &controllers.ArticleController{}, "post:ChangeStatus"),
		beego.NSRouter("/draft", &controllers.ArticleController{}, "get:Draft"),
	)
	// 点赞记录
	fns := beego.NewNamespace("/favors",
		beego.NSRouter("/", &controllers.FavorRecordController{}),
	)
	// 评论管理
	cns2 := beego.NewNamespace("/comments",
		beego.NSRouter("/", &controllers.CommentController{}),
		beego.NSRouter("/keyword", &controllers.CommentController{}, "get:Keyword"),
	)
	// 日志管理
	lns := beego.NewNamespace("/logs",
		beego.NSRouter("/home", &controllers.LogController{}, "get:HomeLog"),
		beego.NSRouter("/admin", &controllers.LogController{}, "get:AdminLog"),
	)
	// 系统管理
	sns := beego.NewNamespace("/system",
		beego.NSRouter("/managers", &controllers.ManagerController{}),
		beego.NSRouter("/roles", &controllers.RoleController{}),
		beego.NSRouter("/permissions", &controllers.PermissionController{}),
	)
	// 登录、退出等
	ans2 := beego.NewNamespace("/auth",
		beego.NSRouter("/login", &controllers.AuthController{}, "get:Login;post:DoLogin"),
		beego.NSRouter("/captcha", &controllers.AuthController{}, "get:GetCaptcha"),
		beego.NSRouter("/logout", &controllers.AuthController{}, "get:Logout"),
	)

	beego.AddNamespace(uns, cns, tns, ans, fns, cns2, lns, sns, ans2)
}

package routers

import (
	"manage/controllers"

	"github.com/astaxie/beego"
)

func init() {

	// 首页
	beego.Router("/", &controllers.IndexController{}, "get:Get")
	// 账号管理
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
		beego.NSRouter("/delete", &controllers.ArticleController{}, "post:Delete"),
		beego.NSRouter("/change", &controllers.ArticleController{}, "post:ChangeStatus"),
		beego.NSRouter("/uploadImage", &controllers.ArticleController{}, "post:UploadImage"),
		beego.NSRouter("/create", &controllers.ArticleController{}, "get:Add;post:Post"),
		beego.NSRouter("/edit/:id", &controllers.ArticleController{}, "get:Edit"),
		beego.NSRouter("/update", &controllers.ArticleController{}, "post:Update"),
	)
	// 点赞记录
	fns := beego.NewNamespace("/favors",
		beego.NSRouter("/", &controllers.FavorRecordController{}, "get:Get"),
	)
	// 评论管理
	cns2 := beego.NewNamespace("/comments",
		beego.NSRouter("/", &controllers.CommentController{}, "get:Get"),
		beego.NSRouter("/delete", &controllers.CommentController{}, "post:Delete"),
	)
	// 干货管理
	rn := beego.NewNamespace("/resource",
		beego.NSRouter("/", &controllers.ResourceController{}, "get:Get;post:Post"),
		beego.NSRouter("/delete", &controllers.ResourceController{}, "post:DeleteResource"),
	)
	// 日志管理
	lns := beego.NewNamespace("/logs",
		beego.NSRouter("/email", &controllers.LogController{}, "get:EmailLog"),
		beego.NSRouter("/home", &controllers.LogController{}, "get:HomeLog"),
		beego.NSRouter("/admin", &controllers.LogController{}, "get:AdminLog"),
	)
	// 用户管理
	sns := beego.NewNamespace("/managers",
		beego.NSRouter("/", &controllers.ManagerController{}, "get:Get;post:Post"),
		beego.NSRouter("/info", &controllers.ManagerController{}, "post:GetInfo"),
		beego.NSRouter("/update", &controllers.ManagerController{}, "post:Update"),
		beego.NSRouter("/delete", &controllers.ManagerController{}, "post:Delete"),
	)
	// 登录、退出等
	ans2 := beego.NewNamespace("/auth",
		beego.NSRouter("/login", &controllers.AuthController{}, "get:Login;post:DoLogin"),
		beego.NSRouter("/captcha", &controllers.AuthController{}, "get:GetCaptcha"),
		beego.NSRouter("/logout", &controllers.AuthController{}, "get:Logout"),
		beego.NSRouter("/reset", &controllers.AuthController{}, "get:ShowResetPassword;post:ResetPassword"),
		beego.NSRouter("/reset/sendEmail", &controllers.AuthController{}, "post:SendResetPasswordEmail"),
	)

	beego.AddNamespace(uns, cns, tns, ans, fns, cns2, lns, sns, ans2, rn)
}

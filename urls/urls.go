package urls

var (
	Main       = "http://twitter.com/"
	MainMobile = "https://mobile.twitter.com/"
	User       = MainMobile + "%s"
	NoJSRouter = MainMobile + "i/nojs_router?path=%%2F%s"
	Status     = Main + "/%s/status/%d/"
)

package enums

const (
	Success                      = 200   //成功
	ErrProxy                     = 10001 //代理错误
	ErrConn                      = 10002 //链接错误
	ErrParams                    = 20002 //参数错误
	ErrDecode                    = 30003 //解析错误
	ErrRiskControl               = 40001 //风控错误
	ErrRiskControlClickCode      = 40002 //点击验证风控
	ErrRiskControlJavaScriptCode = 40003 //JavaScript风控

	ErrParseHtml = 50001 //解析错误
	Err429       = 429   //请求过于频繁

	EngineGoogle = "google"

	ErrRiskControlClick      = "If you're having trouble accessing Google Search, please" //cookie 用多了
	ErrRiskControlJavascript = "window.location.href = "
)

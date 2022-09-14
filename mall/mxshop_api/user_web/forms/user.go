package forms

//PasswordLoginForm 登录用的表单
type PasswordLoginForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式如何设置?需要自定义验证
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`

	//增加输入的验证码结果
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻， 自定义validator
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"` //验证码
}
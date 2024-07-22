package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻， 自定义validator
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`
	//1. 注册发送短信验证码和动态验证码登录发送验证码
}

// 登录
type LoginForm struct {
	UserName string `form:"username" json:"mobile" binding:"required,min=1,max=50"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=16"`
}

type RegisterForm struct {
	UserName  string `form:"username" json:"username" binding:"required,min=1,max=50"`
	Password  string `form:"password" json:"password" binding:"required,min=6,max=16"`
	Captcha   string `form:"captcha"  json:"captcha" binding:"required,len=5"`
	CaptchaId string `form:"captchaId" json:"captchaId" binding:"required"`
}

type ProfileForm struct {
	Nickname string `form:"nickname" json:"nickname" binding:"required,min=1,max=50"` // 用户昵称
	Icon     string `form:"icon" json:"icon"`                                         // 头像地址
	Phone    string `form:"phone" json:"phone" binding:"mobile"`
	Email    string `form:"email" json:"email" binding:"email"`
	Gender   int    `form:"gender" json:"gender" binding:"oneof=1 2"` // 性别，1女 2男
	Status   int    `form:"status" json:"status"`                     // 状态，例如是否有效
}

type PetForm struct {
	PetId         int      `form:"petId" json:"petId"`
	PetType       int      `form:"petType" json:"petType" binding:"required,oneof=1 2"`
	PetBreeds     string   `form:"petBreeds" json:"petBreeds"`
	PetNickname   string   `form:"petNickname" json:"petNickname"`
	PetGender     string   `form:"petGender" json:"petGender"`
	PetAge        string   `form:"petAge" json:"petAge"`
	PetAddress    string   `form:"petAddress" json:"petAddress"`
	PetStatus     string   `form:"petStatus" json:"petStatus"`
	PetExperience string   `form:"petExperience" json:"petExperience"`
	PetAvatar     []string `form:"petAvatar" json:"petAvatar"`
	PetIntro      string   `form:"petIntro" json:"petIntro"`
	Status        int      `form:"status" json:"status" binding:"required,oneof=1 2"`
}
type FilterForm struct {
	PetType    int    `form:"petType" json:"petType" binding:"oneof=0 1 2"`
	PetAddress string `form:"petAddress" json:"petAddress"`
	PetStatus  string `form:"petStatus" json:"petStatus"`
	PetGender  string `form:"petGender" json:"petGender"`
	PetAge     string `form:"petAge" json:"petAge"`
	Page       int    `form:"page" json:"page" binding:"required"`
	PageSize   int    `form:"pageSize" json:"pageSize" binding:"required"`
}

type CommentForm struct {
	PetId    uint   `form:"petId" json:"petId" binding:"required"`       //宠物id
	TargetId uint   `form:"targetId" json:"targetId" binding:"required"` //目标用户id
	Level    int    `form:"level" json:"level" binding:"oneof=0 1 2"`    //层级
	RootID   uint   `form:"rootId" json:"rootId"`                        //根评论id
	Content  string `form:"content" json:"content" binding:"required"`   //内容
}

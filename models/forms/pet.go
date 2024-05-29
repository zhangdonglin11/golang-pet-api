package forms

type PetInfo struct {
	ID            uint     `json:"peyId"`
	PetType       int      `json:"petType"`
	PetBreeds     string   `json:"petBreeds"`
	PetNickname   string   `json:"petNickname"`
	PetGender     string   `json:"petGender"`
	PetAge        string   `json:"petAge"`
	PetAddress    string   `json:"petAddress"`
	PetStatus     string   `json:"petStatus"`
	PetExperience string   `json:"petExperience"`
	PetAvatar     []string `json:"petAvatar"`
	PetIntro      string   `json:"petIntro"`
	Status        int      `json:"status"`
	UserID        uint     `json:"userId"`
}

type PetFilter struct {
	Page       int    `json:"page"`       // 第几页
	Size       int    `json:"size"`       // 每页多少条数据
	PetType    int    `json:"petType"`    // 宠物类型
	PetAddress string `json:"petAddress"` // 宠物地址
	PetStatus  string `json:"petStatus"`  //宠物状态
	PetGender  string `json:"petGender"`  // 宠物性别
	PetAge     string `json:"petAge"`     // 宠物年龄
}

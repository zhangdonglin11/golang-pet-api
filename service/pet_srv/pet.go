package pet_srv

import (
	"errors"
	"golang-pet-api/common/global"
	"golang-pet-api/models/forms"
	"golang-pet-api/models/model"
	"golang-pet-api/service"
	"gorm.io/gorm"
)

// CreatePet 创建宠物信息
func CreatePet(pet *model.Pet) error {
	result := global.Db.Create(pet)
	if result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}

// UpdatePet 保存宠物信息
func SavePet(pet *model.Pet) error {
	result := global.Db.Save(pet)
	if result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}

// FindPetById pet_id 查询宠物
func FindPetById(petId, userId uint) (pet model.Pet, err error) {
	result := global.Db.Where(model.Pet{Model: gorm.Model{ID: petId}, UserID: userId}).First(&pet)
	if result.RowsAffected == 0 {
		return pet, errors.New("宠物不存在")
	}
	return pet, nil
}

// FindPetDraft 用户id查询宠物草稿
func FindPetDraft(userId uint) (pet model.Pet, err error) {
	result := global.Db.Where(map[string]interface{}{"user_id": userId, "status": 0}).First(&pet)
	if result.RowsAffected == 0 {
		return pet, errors.New("宠物不存在")
	}
	return pet, nil
}

// PetList 获取宠物列表
func PetList(petFilter forms.PetFilter) (petList []model.Pet) {
	localDb := global.Db.Model(&model.Pet{})
	// 分页查询
	localDb.Scopes(service.Paginate(petFilter.Page, petFilter.Size)).Where("status = ?", 2)
	if petFilter.PetType != 0 {
		localDb.Where("pet_type=?", petFilter.PetType)
	}
	if petFilter.PetAddress != "" {
		localDb.Where("pet_Address LIKE ?", "%"+petFilter.PetAddress+"%")
	}
	if petFilter.PetStatus != "" {
		localDb.Where("pet_status=?", petFilter.PetStatus)
	}
	if petFilter.PetGender != "" {
		localDb.Where("pet_gender=?", petFilter.PetGender)
	}
	if petFilter.PetAge != "" {
		localDb.Where("pet_Age=?", petFilter.PetAge)
	}
	localDb.Find(&petList)
	return
}

// GetLikeStatus 获取宠物点赞信息
func GetLikeStatus(petId, userId uint) (count int64, isLike bool) {
	// 收藏数量
	global.Db.Model(model.PetLike{}).Where("pet_id=?", petId).Count(&count)

	// 判断该用户是否已收藏该宠物
	if userId != 0 {
		var petLike model.PetLike
		res := global.Db.Where(model.PetLike{PetID: petId, UserID: userId}).First(&petLike)
		if res.RowsAffected != 0 {
			isLike = true
		}
	}
	return
}

// DeletePetByPid 删除宠物
func DeletePetByPid(petId uint) error {
	result := global.Db.Delete(&model.Pet{}, petId)
	if result.RowsAffected == 0 {
		return result.Error
	}
	DeletePetLikeByPid(petId)
	return nil
}

// DeletePetLikeByPid 删除点赞信息
func DeletePetLikeByPid(petId uint) {
	petLIke := model.PetLike{PetID: petId}
	global.Db.Delete(&petLIke)
}

package example

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/example"
	"47.103.136.241/goprojects/curesan/server/model/system"
	systemService "47.103.136.241/goprojects/curesan/server/service/system"
)

type CustomerService struct{}

// @author: DingYG
// @function: CreateExaCustomer
// @description: 创建客户
// @param: e model.ExaCustomer
// @return: err error

func (exa *CustomerService) CreateExaCustomer(e example.ExaCustomer) (err error) {
	err = global.GVA_DB.Create(&e).Error
	return err
}

// @author: DingYG
// @function: DeleteFileChunk
// @description: 删除客户
// @param: e model.ExaCustomer
// @return: err error

func (exa *CustomerService) DeleteExaCustomer(e example.ExaCustomer) (err error) {
	err = global.GVA_DB.Delete(&e).Error
	return err
}

// @author: DingYG
// @function: UpdateExaCustomer
// @description: 更新客户
// @param: e *model.ExaCustomer
// @return: err error

func (exa *CustomerService) UpdateExaCustomer(e *example.ExaCustomer) (err error) {
	err = global.GVA_DB.Save(e).Error
	return err
}

// @author: DingYG
// @function: GetExaCustomer
// @description: 获取客户信息
// @param: id uint
// @return: customer model.ExaCustomer, err error

func (exa *CustomerService) GetExaCustomer(id uint) (customer example.ExaCustomer, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&customer).Error
	return
}

// @author: DingYG
// @function: GetCustomerInfoList
// @description: 分页获取客户列表
// @param: sysUserAuthorityID string, info request.PageInfo
// @return: list interface{}, total int64, err error

func (exa *CustomerService) GetCustomerInfoList(sysUserAuthorityID uint, info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&example.ExaCustomer{})
	var a system.SysAuthority
	a.AuthorityId = sysUserAuthorityID
	auth, err := systemService.AuthorityServiceApp.GetAuthorityInfo(a)
	if err != nil {
		return
	}
	var dataId []uint
	for _, v := range auth.DataAuthorityId {
		dataId = append(dataId, v.AuthorityId)
	}
	var CustomerList []example.ExaCustomer
	err = db.Where("sys_user_authority_id in ?", dataId).Count(&total).Error
	if err != nil {
		return CustomerList, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Preload("SysUser").Where("sys_user_authority_id in ?", dataId).Find(&CustomerList).Error
	}
	return CustomerList, total, err
}

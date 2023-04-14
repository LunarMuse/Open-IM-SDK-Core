//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(groupRequest).Error, "InsertAdminGroupRequest failed")
}

func (d *DataBase) DeleteAdminGroupRequest(ctx context.Context, groupID, userID string) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Where("group_id=? and user_id=?", groupID, userID).Delete(&model_struct.LocalAdminGroupRequest{}).Error, "DeleteAdminGroupRequest failed")
}

func (d *DataBase) UpdateAdminGroupRequest(ctx context.Context, groupRequest *model_struct.LocalAdminGroupRequest) error {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	t := d.conn.WithContext(ctx).Model(groupRequest).Select("*").Updates(*groupRequest)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

func (d *DataBase) GetAdminGroupApplication(ctx context.Context) ([]*model_struct.LocalAdminGroupRequest, error) {
	d.groupMtx.Lock()
	defer d.groupMtx.Unlock()
	var groupRequestList []model_struct.LocalAdminGroupRequest
	err := utils.Wrap(d.conn.WithContext(ctx).Order("create_time DESC").Find(&groupRequestList).Error, "")
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var transfer []*model_struct.LocalAdminGroupRequest
	for _, v := range groupRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, nil
}

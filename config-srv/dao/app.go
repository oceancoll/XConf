package dao

import (
	"fmt"

	"XConf/config-srv/model"
)

func (d *Dao) AppExist(appName string) bool {
	return !d.client.Table("app").Where("app_name = ?", appName).First(&model.App{}).RecordNotFound()
}

func (d *Dao) CreateApp(appName, description string) (*model.App, error) {
	app := &model.App{
		AppName:     appName,
		Description: description,
	}
	err := d.client.Table("app").Create(app).Error

	return app, err
}

func (d *Dao) DeleteApp(appName string) error {
	var err error
	tx := d.client.Begin()
	defer func() {
		if err != nil {
			err = fmt.Errorf("[DeleteApp] delete app:%s error: %s", appName, err.Error())
			tx.Rollback()
		}
	}()

	if err = tx.Table("namespace").Unscoped().Delete(model.Namespace{}, "app_name = ?", appName).Error; err != nil {
		return err
	}
	if err = tx.Table("cluster").Unscoped().Delete(model.Cluster{}, "app_name = ?", appName).Error; err != nil {
		return err
	}
	if err = tx.Table("app").Unscoped().Delete(model.App{}, "app_name = ?", appName).Error; err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (d *Dao) ListApps() (apps []model.App, err error) {
	err = d.client.Table("app").Find(&apps).Error
	return
}

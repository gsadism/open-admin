package hooks

import (
	"context"
	"fmt"
	basemodels "github.com/gsadism/open-admin/core/base/models"
	"github.com/gsadism/open-admin/logging"
	"github.com/gsadism/open-admin/osv"
	"path/filepath"
)

func IrWebsiteSettingInit() {
	var n int64

	db := osv.DB.WithContext(context.WithValue(context.TODO(), "sudo", true))

	if err := db.Table("ir_website_setting").Count(&n).Error; err != nil {
		logging.Warn(err)
	} else {
		if n <= 0 {
			// 检查 bucket 是否存在
			if !osv.Minio.Client().BucketExists(context.TODO(), "website") {
				if err := osv.Minio.Client().MakeBucket(context.TODO(), "website", "wu-han"); err != nil {
					logging.Warn(err.Error())
					return
				} else {
					if err := osv.Minio.Client().SetBucketPolicy(context.TODO(), "website", "public"); err != nil {
						logging.Warn(err.Error())
						return
					}
				}
			}

			// 上传logo和background
			LogoPath := ""
			BackgroundPath := ""
			if _, err := osv.Minio.Client().UploadByName(context.TODO(), "website", "logo.png", filepath.Join(osv.STATIC, "logo.png")); err != nil {
				logging.Warn(err.Error())
			} else {
				LogoPath = fmt.Sprintf("%s/%s", "website", "logo.png")
			}
			if _, err := osv.Minio.Client().UploadByName(context.TODO(), "website", "background.svg", filepath.Join(osv.STATIC, "background.svg")); err != nil {
				logging.Warn(err.Error())
			} else {
				BackgroundPath = fmt.Sprintf("%s/%s", "website", "background.svg")
			}

			db.Table("ir_website_setting").Create(&basemodels.IRWebsiteSetting{
				Title:      "Open Admin Vue",
				Company:    "open-admin",
				ICP:        "浙ICP备2023027841号",
				Logo:       LogoPath,
				Background: BackgroundPath,
			})
		}
	}
}

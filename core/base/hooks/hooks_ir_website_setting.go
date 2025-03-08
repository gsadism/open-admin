package hooks

import (
	"context"
	basemodels "github.com/gsadism/open-admin/core/base/models"
	"github.com/gsadism/open-admin/osv"
)

func IrWebsiteSettingInit() {
	var n int64

	db := osv.DB.WithContext(context.WithValue(context.TODO(), "sudo", true))

	db.Table("ir_website_setting").Count(&n)
	if n <= 0 {
		db.Table("ir_website_setting").Create(&basemodels.IRWebsiteSetting{
			Title:   "Open Admin Vue",
			Company: "open-admin",
			ICP:     "浙ICP备2023027841号",
		})
	}
}

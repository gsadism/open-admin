package conf

// todo conf 目录下的所有配置项请务在用户代码中直接导入使用, 会出现循环依赖问题

const (
	// SUPER_ADMIN_PASSWORD : 超级管理员用户密码, 为空则默认密码为 root
	SUPER_ADMIN_PASSWORD = "$2a$10$jk84oOm11lRh1yMg.nvMrep.5SaYNDOBwsztWo6/PFKETTzlgajEK"
	// SECRET_KEY : keep the secret key used in production secret!
	SECRET_KEY = "fec2344ddca2efa5a59d4b8c5ea4a02a870764c51c95f6128af73d20701f544c"
	// SECRET_IV : 偏移量
	SECRET_IV = "bdc4dd83187027d8d1a36faeffd562b8"
)

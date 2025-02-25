package logic

import "smServer/server/game/model/data"

func BeforeInit() {
	data.GetYield = RoleResService.GetYield
	data.GetUnion = RoleAttrService.GetUnion
	data.GetRoleNickName = RoleService.GetRoleNickName
}

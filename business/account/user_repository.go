package account

import (
	"context"
	"encoding/json"
	"github.com/gingerxman/eel"
)

type UserRepository struct {
	eel.ServiceBase
}

func NewUserRepository(ctx context.Context) *UserRepository {
	service := new(UserRepository)
	service.Ctx = ctx
	return service
}

func (this *UserRepository) makeUsers(userDatas []interface{}) []*User {
	users := make([]*User, 0)
	for _, userData := range userDatas {
		userJson := userData.(map[string]interface{})
		id, _ := userJson["id"].(json.Number).Int64()
		user := NewUserFromOnlyId(this.Ctx, int(id))
		user.Unionid = userJson["unionid"].(string)
		user.Name = userJson["name"].(string)
		user.Avatar = userJson["avatar"].(string)
		user.Sex = userJson["sex"].(string)
		user.Code = userJson["code"].(string)
		
		users = append(users, user)
	}
	
	return users
}

//func (this *UserRepository) GetUsers(ids []int) []*User {
//	options := make(map[string]interface{})
//	options["with_role_info"] = true
//	resp, err := vanilla.NewResource(this.Ctx).Get("gskep", "account.users", vanilla.Map{
//		"ids": vanilla.ToJsonString(ids),
//		"with_options": vanilla.ToJsonString(options),
//	})
//
//	if err != nil {
//		beego.Error(err)
//		return nil
//	}
//
//	respData := resp.Data()
//	userDatas := respData.Get("users")
//	fmt.Println(userDatas)
//	return this.makeUsers(userDatas.MustArray())
//}
//
//func (this *UserRepository) GetUsersWithOptions(ids []int, options map[string]interface{}) []*User {
//	resp, err := vanilla.NewResource(this.Ctx).Get("gskep", "account.users", vanilla.Map{
//		"ids": vanilla.ToJsonString(ids),
//		"with_options": vanilla.ToJsonString(options),
//	})
//
//	if err != nil {
//		beego.Error(err)
//		return nil
//	}
//
//	respData := resp.Data()
//	userDatas := respData.Get("users")
//	return this.makeUsers(userDatas.MustArray())
//}
//
//func (this *UserRepository) GetUsersByCodes(codes []string) []*User {
//	resp, err := vanilla.NewResource(this.Ctx).Get("gskep", "account.users", vanilla.Map{
//		"codes": vanilla.ToJsonString(codes),
//	})
//
//	if err != nil {
//		beego.Error(err)
//		return nil
//	}
//
//	respData := resp.Data()
//	userDatas := respData.Get("users")
//	return this.makeUsers(userDatas.MustArray())
//}
//
//func (this *UserRepository) GetUsersByUnionids(unionids []string) []*User {
//	resp, err := vanilla.NewResource(this.Ctx).Get("gskep", "account.users", vanilla.Map{
//		"unionids": vanilla.ToJsonString(unionids),
//	})
//
//	if err != nil {
//		beego.Error(err)
//		return nil
//	}
//
//	respData := resp.Data()
//	userDatas := respData.Get("users")
//	return this.makeUsers(userDatas.MustArray())
//}
//
//func (this *UserRepository) GetUserByCorpUserId(corpUserId int) *User {
//	resp, err := vanilla.NewResource(this.Ctx).Get("gskep", "account.user", vanilla.Map{
//		"corp_user_id": corpUserId,
//	})
//
//	if err != nil {
//		beego.Error(err)
//		return nil
//	}
//
//	userData := resp.Data().MustMap()
//	return this.makeUsers([]interface{}{userData})[0]
//}

func init() {
}

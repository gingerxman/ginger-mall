// Copyright 2018 eel Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"net/http"
	"github.com/gingerxman/eel"
	_ "github.com/gingerxman/ginger-mall/models"
	_ "github.com/gingerxman/ginger-mall/routers"
	_ "github.com/gingerxman/ginger-mall/middleware"
	"github.com/bitly/go-simplejson"
	"github.com/gingerxman/ginger-mall/business/account"
)

func NewBusinessContext(ctx context.Context, request *http.Request, userId int, jwtToken string, rawData *simplejson.Json) context.Context {
	user := new(account.User)
	user.Model = nil
	user.Id = userId
	user.RawData = rawData
	
	ctx = context.WithValue(ctx, "jwt", jwtToken)
	user.Ctx = ctx
	
	ctx = context.WithValue(ctx, "user", user)
	
	//创建corp
	if rawData != nil {
		jwtType, _ := rawData.Get("type").Int()
		if jwtType == 1 {
			corpId, err := rawData.Get("cid").Int()
			if err == nil {
				corp := new(account.Corp)
				corp.Model = nil
				corp.Id = corpId
				corp.Ctx = ctx
				
				//处理corp user
				corpUserId, _ := rawData.Get("uid").Int()
				corpUser := account.NewCorpUserFromId(ctx, corpUserId)
				corp.CorpUser = corpUser
				
				ctx = context.WithValue(ctx, "corp", corp)
			} else {
				eel.Logger.Error(err)
			}
			
		}
	}
	return ctx
}

func main() {
	eel.Runtime.NewBusinessContext = NewBusinessContext
	eel.RunService()
}


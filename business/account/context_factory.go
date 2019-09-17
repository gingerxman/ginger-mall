package account

import (
	"context"
	
	"github.com/bitly/go-simplejson"
	
	"net/http"
)

var gInstance *ContextFactory

type ContextFactory struct {
}

//NewContext 构造含有corp的Context
func (this *ContextFactory) NewContext(ctx context.Context, request *http.Request, userId int, jwtToken string, rawData *simplejson.Json) context.Context {
	return ctx
}

func init() {
	//gInstance = &ContextFactory{}
	//middleware.SetBusinessContextFactory(gInstance)
}

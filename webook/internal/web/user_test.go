package web

import (
	"bytes"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	svcmocks "gitee.com/geekbang/basic-go/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(nil)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对，bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			//格式不对，就会解析失败，根据 Content-Type 解析，这里是Json,破坏它的形式就可以，比如少个}
			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123"
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(nil)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// signup 用不到codeSvc,只需要一个 userSvc,每个 Case 不一样，所以由 Case 传入
			h := NewUserHandler(testCase.mock(ctrl), nil)
			server := gin.Default()
			h.RegisterRoutes(server)

			// server 负责监听请求，但也可以通过 server 发请求：需要 req 携带请求信息，resp 携带响应信息
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(testCase.reqBody)))
			require.NoError(t, err)
			//数据是 JSON 形式
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里,而不是通常的管道一样直接展示给前端
			server.ServeHTTP(resp, req)

			assert.Equal(t, testCase.wantCode, resp.Code)
			assert.Equal(t, testCase.wantBody, resp.Body.String())
		})
	}
}
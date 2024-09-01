package jwt

import (
	"comyun/dao/mysql"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"time"
)

// TokenExpireDuration 过期时间
const TokenExpireDuration = time.Hour * 12

// 签名(加言)
var mySecret = []byte("ILoveYouMoreThanAnything")


// MyClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	// 可根据需要自行添加字段
	UserID               int64  `json:"userid"`
	Access				 int `json:"access"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

// NewToken 生成JWT
func NewToken(userID int64) (string, error) {
	role,err := mysql.GetRole(userID)
	if err != nil{
		zap.L().Error("error to get role"+err.Error())
		return "",err
	}
	// 创建一个我们自己的声明（负载）
	claims := MyClaims{
		// 自定义字段
		userID,
		role,
		jwt.RegisteredClaims{
			//签发时间为获取当前时间之后加上持续时长后转化成Unix
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "bluebell", // 签发人
		},
	}
	// 使用加密算法加密创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	// 1.解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法

	// 直接使用标准的Claim则可以直接使用Parse方法
	//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
	var mc = new(MyClaims)
	//解析到数据到mc中
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	//是否校验成功的bool字段，如果是则返回正确结果
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

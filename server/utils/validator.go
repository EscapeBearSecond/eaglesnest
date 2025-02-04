package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Rules map[string][]string

type RulesMap map[string]Rules

var CustomizeMap = make(map[string]Rules)

// @author: DingYG
// @function: RegisterRule
// @description: 注册自定义规则方案建议在路由初始化层即注册
// @param: key string, rule Rules
// @return: err error

func RegisterRule(key string, rule Rules) (err error) {
	if CustomizeMap[key] != nil {
		return errors.New(key + "已注册,无法重复注册")
	} else {
		CustomizeMap[key] = rule
		return nil
	}
}

// @author: DingYG
// @function: NotEmpty
// @description: 非空 不能为其对应类型的0值
// @return: string

func NotEmpty() string {
	return "notEmpty"
}

// @author: [zooqkl](https://github.com/zooqkl)
// @function: RegexpMatch
// @description: 正则校验 校验输入项是否满足正则表达式
// @param:  rule string
// @return: string

func RegexpMatch(rule string) string {
	return "regexp=" + rule
}

// @author: DingYG
// @function: Lt
// @description: 小于入参(<) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
// @return: string

func Lt(mark string) string {
	return "lt=" + mark
}

// @author: DingYG
// @function: Le
// @description: 小于等于入参(<=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
// @return: string

func Le(mark string) string {
	return "le=" + mark
}

// @author: DingYG
// @function: Eq
// @description: 等于入参(==) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
// @return: string

func Eq(mark string) string {
	return "eq=" + mark
}

// @author: DingYG
// @function: Ne
// @description: 不等于入参(!=)  如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
// @return: string

func Ne(mark string) string {
	return "ne=" + mark
}

// @author: DingYG
// @function: Ge
// @description: 大于等于入参(>=) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
// @return: string

func Ge(mark string) string {
	return "ge=" + mark
}

// @author: DingYG
// @function: Gt
// @description: 大于入参(>) 如果为string array Slice则为长度比较 如果是 int uint float 则为数值比较
// @param: mark string
// @return: string

func Gt(mark string) string {
	return "gt=" + mark
}

//
// @author: DingYG
// @function: Verify
// @description: 校验方法
// @param: st interface{}, roleMap Rules(入参实例，规则map)
// @return: err error

func Verify(st interface{}, roleMap Rules) (err error) {
	compareMap := map[string]bool{
		"lt": true,
		"le": true,
		"eq": true,
		"ne": true,
		"ge": true,
		"gt": true,
	}

	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st) // 获取reflect.Type类型

	kd := val.Kind() // 获取到st对应的类别
	if kd != reflect.Struct {
		return errors.New("expect struct")
	}
	num := val.NumField()
	// 遍历结构体的所有字段
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i)
		val := val.Field(i)
		if tagVal.Type.Kind() == reflect.Struct {
			if err = Verify(val.Interface(), roleMap); err != nil {
				return err
			}
		}
		if len(roleMap[tagVal.Name]) > 0 {
			for _, v := range roleMap[tagVal.Name] {
				switch {
				case v == "notEmpty":
					if isBlank(val) {
						return errors.New(tagVal.Name + "值不能为空")
					}
				case strings.Split(v, "=")[0] == "regexp":
					if !regexpMatch(strings.Split(v, "=")[1], val.String()) {
						return errors.New(tagVal.Name + "格式校验不通过")
					}
				case compareMap[strings.Split(v, "=")[0]]:
					if !compareVerify(val, v) {
						return errors.New(tagVal.Name + "长度或值不在合法范围," + v)
					}
				}
			}
		}
	}
	return nil
}

// @author: DingYG
// @function: compareVerify
// @description: 长度和数字的校验方法 根据类型自动校验
// @param: value reflect.Value, VerifyStr string
// @return: bool

func compareVerify(value reflect.Value, VerifyStr string) bool {
	switch value.Kind() {
	case reflect.String:
		return compare(len([]rune(value.String())), VerifyStr)
	case reflect.Slice, reflect.Array:
		return compare(value.Len(), VerifyStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(value.Uint(), VerifyStr)
	case reflect.Float32, reflect.Float64:
		return compare(value.Float(), VerifyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(value.Int(), VerifyStr)
	default:
		return false
	}
}

// @author: DingYG
// @function: isBlank
// @description: 非空校验
// @param: value reflect.Value
// @return: bool

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// @author: DingYG
// @function: compare
// @description: 比较函数
// @param: value interface{}, VerifyStr string
// @return: bool

func compare(value interface{}, VerifyStr string) bool {
	VerifyStrArr := strings.Split(VerifyStr, "=")
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		VInt, VErr := strconv.ParseInt(VerifyStrArr[1], 10, 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Int() < VInt
		case VerifyStrArr[0] == "le":
			return val.Int() <= VInt
		case VerifyStrArr[0] == "eq":
			return val.Int() == VInt
		case VerifyStrArr[0] == "ne":
			return val.Int() != VInt
		case VerifyStrArr[0] == "ge":
			return val.Int() >= VInt
		case VerifyStrArr[0] == "gt":
			return val.Int() > VInt
		default:
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		VInt, VErr := strconv.Atoi(VerifyStrArr[1])
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Uint() < uint64(VInt)
		case VerifyStrArr[0] == "le":
			return val.Uint() <= uint64(VInt)
		case VerifyStrArr[0] == "eq":
			return val.Uint() == uint64(VInt)
		case VerifyStrArr[0] == "ne":
			return val.Uint() != uint64(VInt)
		case VerifyStrArr[0] == "ge":
			return val.Uint() >= uint64(VInt)
		case VerifyStrArr[0] == "gt":
			return val.Uint() > uint64(VInt)
		default:
			return false
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(VerifyStrArr[1], 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Float() < VFloat
		case VerifyStrArr[0] == "le":
			return val.Float() <= VFloat
		case VerifyStrArr[0] == "eq":
			return val.Float() == VFloat
		case VerifyStrArr[0] == "ne":
			return val.Float() != VFloat
		case VerifyStrArr[0] == "ge":
			return val.Float() >= VFloat
		case VerifyStrArr[0] == "gt":
			return val.Float() > VFloat
		default:
			return false
		}
	default:
		return false
	}
}

func regexpMatch(rule, matchStr string) bool {
	return regexp.MustCompile(rule).MatchString(matchStr)
}

var validate = newValidate()

// BindAndValid
// 绑定并验证
func BindAndValid(c *gin.Context, v interface{}) error {
	err := c.ShouldBind(v)
	if err != nil {
		return err
	}
	return validate.Struct(v)
}

func newValidate() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}

func ValidateIP(raws []string) error {
	for _, raw := range raws {
		if !validateIp(raw) {
			ipRange := strings.Split(raw, "-")
			if len(ipRange) == 2 {
				if !validateIp(ipRange[0]) || !validateIp(ipRange[1]) {
					return fmt.Errorf("%s 不是一个正确的IP或IP范围", raw)
				}
			} else {
				return fmt.Errorf("%s 不是一个正确的IP或IP范围", raw)
			}
		}
	}
	return nil
}

func validateIp(ip string) bool {
	ip = strings.Trim(ip, " ")
	if net.ParseIP(ip) == nil {
		if _, _, err := net.ParseCIDR(ip); err != nil {
			return false
		}
	}
	return true
}

func IsValidCron(expr string) bool {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := parser.Parse(expr)
	return err == nil
}

func GetLocalIP() (string, error) {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// 遍历所有网络接口
	for _, iface := range interfaces {
		// 过滤掉回环地址和未启用的接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的地址
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		// 遍历接口的地址
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 过滤掉回环地址和非 IPv4 地址
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			// 返回第一个非回环的 IPv4 地址
			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("no valid IP address found")
}

// 判断 IP 是否在指定的 CIDR 或 IP 范围内
func IsIPInRange(ipStr, rangeStr string) bool {
	if ipStr == rangeStr {
		return true
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	// 判断是否为 CIDR 格式
	if strings.Contains(rangeStr, "/") {
		_, network, err := net.ParseCIDR(rangeStr)
		if err != nil {
			return false
		}
		return network.Contains(ip)
	}

	// 判断是否为 IP 区间格式
	if strings.Contains(rangeStr, "-") {
		parts := strings.Split(rangeStr, "-")
		if len(parts) != 2 {
			return false
		}
		startIP := net.ParseIP(strings.TrimSpace(parts[0]))
		endIP := net.ParseIP(strings.TrimSpace(parts[1]))
		if startIP == nil || endIP == nil {
			return false
		}
		return compareIPs(ip, startIP) >= 0 && compareIPs(ip, endIP) <= 0
	}

	return false
}

// 比较两个 IP 地址的字节序大小
func compareIPs(ip1, ip2 net.IP) int {
	for i := 0; i < len(ip1); i++ {
		if ip1[i] < ip2[i] {
			return -1
		} else if ip1[i] > ip2[i] {
			return 1
		}
	}
	return 0
}

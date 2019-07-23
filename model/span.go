package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

//Span 时间长度，对应pgsql中的interval类型，
//关于pgsql中interval类型和Duration的映射，有关讨论见https://github.com/lib/pq/issues/78
// The range of time.Duration in go is +/- 290 years. In postgresql it's considerably larger: +/- 178000000 years.
//在本项目中+/-290 years的值完全够用，因此考虑直接使用此类型映射
type Span time.Duration

//Day  返回天数的整数部分
func (m Span) Day() int {
	hours := time.Duration(m).Hours()
	ih := int(hours)
	return int(ih / 24)
}

//Hours 全部小时部分
func (m Span) Hours() int {
	return int(time.Duration(m).Hours())
}

//Hours  除去天数的小时部分
func (m Span) HoursInDay() int {
	return m.Hours() % 24
}

//AllMinutes 以分钟计算
func (m Span) AllMinutes() int {
	src := time.Duration(m)
	return int(src.Minutes())
}

//Minute 分钟部分
func (m Span) Minute() int {
	src := time.Duration(m)
	nsec := src % time.Hour
	minPart := float64(nsec) / (60 * 1e9)
	return int(minPart)
}

//Second  秒部分
func (m Span) Second() int {
	src := time.Duration(m)
	nsec := src % time.Minute
	secPart := float64(nsec) / 1e9
	return int(secPart)
}

//MarshalJSON 实现Span的json序列化方法
func (this Span) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%02d:%02d:%02d\"", this.Hours(), this.Minute(), this.Second())
	return []byte(stamp), nil
}

//UnmarshalJSON 解析Span字符串，格式 d (day) hh:mm:ss 或 d hh:mm:ss
func (this *Span) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return this.Parse(s)
}

//Parse 解析span格式
func (this *Span) Parse(s string) error {
	sl := strings.ToLower(s)
	sl = strings.Replace(sl, "day ", "", -1)
	hp := sl
	ss := strings.Split(sl, " ")
	var tar string
	if len(ss) > 1 {
		tar = ss[0] + "d"
		hp = ss[1]
	}
	sh := strings.Split(hp, ":")
	if len(sh) > 0 {
		tar = tar + sh[0] + "h"
	}
	if len(sh) > 1 {
		tar = tar + sh[1] + "m"
	}
	if len(sh) > 2 {
		tar = tar + sh[2] + "s"
	}
	du, e := time.ParseDuration(tar)
	if e != nil {
		return e
	}
	*this = Span(du)
	return nil
}

//Str 获取字符串值
func (this Span) Str() string {
	return time.Duration(this).String()
}

//Value db接口
func (this Span) Value() (driver.Value, error) {
	sv := fmt.Sprintf("%02d:%02d:%02d", this.Hours(), this.Minute(), this.Second())
	return sv, nil
}

//Scan db接口
func (this *Span) Scan(value interface{}) error {
	//log.Debugf("Span Scan value type is %t", value)
	switch v := value.(type) {
	case []uint8:
		this.Parse(string(v))
	case int64:
		*this = Span(v)
	case time.Duration:
		*this = Span(v)
	case string:
		this.Parse(v)
	default:
		return errors.New(fmt.Sprintf("Span Scan Not suport type %T", value))
	}
	return nil
}

//IsEmpty 本系统精确到分钟，小于1秒则认为是空值了。
func (this Span) IsEmpty() bool {
	return time.Duration(this) < time.Second
}

//NewSpan  从字符串创建对象
func ParseSpan(v string) Span {
	var sp Span
	sp.Parse(v)
	return sp
}

//DateTime
//曾经的坑：
//time.Time 包含时区信息，在pg数据库里的字段设计为timestamp without time zone,
//存数据库后，再取出来和当前时间做减法算时间差居然是负数！！！
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"
)

//DateTime 长日期格式 包含日期和时间部分2006-01-02 15:04:05
type DateTime time.Time

const (
	dtFmt = "2006-01-02 15:04:05"
	dFmt  = "2006-01-02"
	tFmt  = "15:04:05"
	stFmt = "15:04"
)
const (
	dtEmpty = "0001-01-01 00:00:00"
	dEmpty  = "0001-01-01"
	tEmpty  = "00:00:00"
)

//MarshalJSON 实现DateTime的json序列化方法
func (this DateTime) MarshalJSON() ([]byte, error) {
	tv := time.Time(this).Format(dtFmt)
	if dtEmpty == tv {
		tv = ""
	}
	var stamp = fmt.Sprintf("\"%s\"", tv)
	return []byte(stamp), nil
}

//UnmarshalJSON 解析日期字符串，支持全格式，和仅日期，仅时间格式
func (this *DateTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		s = dtEmpty
	}
	return this.Parse(s)
}

//Parse 解析日期字符串,支持全格式，和仅日期，仅时间格式
func (this *DateTime) Parse(s string) error {
	//	lf := "2006-01-02 15:04"
	//根据输入格式进行解析
	l := "^[0-9]{4}-[0-9]{1,2}-[0-9]{1,2} [0-9]{1,2}:[0-9]{1,2}$"
	if ok, _ := regexp.MatchString(l, s); ok {
		t, err := time.Parse(dtFmt, s)
		if err != nil {
			return err
		}
		*this = DateTime(t)
		return nil
	}
	//df := "2006-01-02"
	//根据输入格式进行解析
	d := "^[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}$"
	if isd, _ := regexp.MatchString(d, s); isd {
		t, err := time.Parse(dFmt, s)
		if err != nil {
			return err
		}
		*this = DateTime(t)
		return nil
	}
	//tf := "15:04"
	//根据输入格式进行解析
	ts := "^[0-9]{1,2}:[0-9]{1,2}$"
	if ist, _ := regexp.MatchString(ts, s); ist {
		t, err := time.Parse(stFmt, s)
		if err != nil {
			return err
		}
		*this = DateTime(t)
		return nil
	}
	return errors.New(fmt.Sprintf("Not konw formart for value %s", s))
}
func (this DateTime) DateStr() string {
	return time.Time(this).Format(dFmt)
}
func (this DateTime) TimeStr() string {
	return time.Time(this).Format(tFmt)
}
func (this DateTime) Str() string {
	return time.Time(this).Format(dtFmt)
}

func (this DateTime) Value() (driver.Value, error) {
	//return time.Time(this), nil
	return time.Time(this).Format(dtFmt), nil
}

func (this DateTime) IsEmpty() bool {
	return time.Time(this).Equal(time.Time{})
}

func (this *DateTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*this = DateTime(v)
	case string:
		this.Parse(v)
	default:
		return errors.New(fmt.Sprintf("Not suport type %T", value))
	}
	return nil
}

//Add 日期加法，在当前长日期加上指定的小时，分，返回添加后的新值
func (this DateTime) Add(st Span) DateTime {
	d := time.Duration(st)
	src := time.Time(this)
	nt := src.Add(d)
	return DateTime(nt)
}

//Sub 日期减法，在当前长日期减去指定的小时，分，返回减后的新值
func (this DateTime) Sub(st STime) DateTime {
	d, e := st.ToDuration()
	if e != nil {
		return this
	}
	src := time.Time(this)
	nt := src.Add(-1 * d)
	return DateTime(nt)
}

//Span 计算t和当前值的差值，t-this
//曾经的坑：
//time.Time 包含时区信息，在pg数据库里的字段设计为timestamp without time zone,
//存数据库后，再取出来和当前时间做减法算时间差居然是负数！！！
func (this DateTime) Span(t time.Time) Span {
	return Span(t.Sub(time.Time(this)))
}

//-----------------------------------------------------------------------

//Date 日期格式，仅包含日期部分
type Date time.Time

var emptyDate = "0001-01-01"

//UnmarshalJSON 实现Date的json反序列化方法
func (this *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		s = emptyDate
	}
	return this.Parse(s)
}

//MarshalJSON 实现Date的json序列化方法
func (this Date) MarshalJSON() ([]byte, error) {
	tv := time.Time(this).Format("2006-01-02")
	if emptyDate == tv {
		tv = ""
	}
	var stamp = fmt.Sprintf("\"%s\"", tv)
	return []byte(stamp), nil
}

func (this Date) Value() (driver.Value, error) {
	return time.Time(this), nil
}

func (this *Date) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*this = Date(v)
	case string:
		this.Parse(v)
	default:
		return errors.New(fmt.Sprintf("Not suport type %T", value))
	}
	return nil
}

func (this *Date) Parse(s string) error {

	df := "2006-01-02"
	//根据输入格式进行解析
	d := "^[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}$"
	if isd, _ := regexp.MatchString(d, s); isd {
		t, err := time.Parse(df, s)
		if err != nil {
			return err
		}
		*this = Date(t)
		return nil
	}
	return errors.New(fmt.Sprintf("Not konw formart for value %s", s))
}

func (this Date) IsEmpty() bool {
	//	return time.Time(this).Equal(time.Time{})
	return this.DateStr() == emptyDate
}
func (this Date) DateStr() string {
	return time.Time(this).Format("2006-01-02")
}

//--------------------------------------------------------------------------

// STime short time
type STime time.Time

var emptyTime = "00:00:00"

//MarshalJSON 实现DateTime的json序列化方法
func (this STime) MarshalJSON() ([]byte, error) {
	tv := time.Time(this).Format("15:04:05")
	if emptyTime == tv {
		tv = ""
	}
	var stamp = fmt.Sprintf("\"%s\"", tv)
	return []byte(stamp), nil
}

func (this *STime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		s = emptyTime
	}
	return this.Parse(s)
}

func (this STime) Value() (driver.Value, error) {
	return time.Time(this), nil
}

func (this *STime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*this = STime(v)
	case string:
		this.Parse(v)
	default:
		return errors.New(fmt.Sprintf("Not suport type %T", value))
	}
	return nil
}
func NewSTime(s string) STime {
	var st STime
	st.Parse(s)
	return st
}

func (this *STime) Parse(s string) error {
	//tf := "15:04"
	//根据输入格式进行解析
	ts := "^[0-9]{1,2}:[0-9]{1,2}(:[0-9]{1,2})?$"
	if ist, _ := regexp.MatchString(ts, s); ist {
		t, err := time.Parse(stFmt, s)
		if err != nil {
			t, err = time.Parse(tFmt, s) //尝试另外一种格式解析
			if err != nil {
				return err
			}
		}
		*this = STime(t)
		return nil
	}
	return errors.New(fmt.Sprintf("Not konw formart for value %s", s))
}

func (this STime) IsEmpty() bool {
	sv := this.TimeStr()
	return sv == "" || sv == "0" || sv == "00:00" || sv == emptyTime
}
func (this STime) TimeStr() string {
	return time.Time(this).Format("15:04:05")
}

//ToDuration  转换为Duration值
func (this STime) ToDuration() (time.Duration, error) {
	src := time.Time(this)
	str := fmt.Sprintf("%dh%dm%ds", src.Hour(), src.Minute(), src.Second())
	return time.ParseDuration(str)
}

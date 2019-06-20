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

//MarshalJSON 实现DateTime的json序列化方法
func (this DateTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

//UnmarshalJSON 解析日期字符串，支持全格式，和仅日期，仅时间格式
func (this *DateTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return this.Parse(s)
}

//Parse 解析日期字符串,支持全格式，和仅日期，仅时间格式
func (this *DateTime) Parse(s string) error {
	lf := "2006-01-02 15:04"
	//根据输入格式进行解析
	l := "^[0-9]{4}-[0-9]{1,2}-[0-9]{1,2} [0-9]{1,2}:[0-9]{1,2}$"
	if ok, _ := regexp.MatchString(l, s); ok {
		t, err := time.Parse(lf, s)
		if err != nil {
			return err
		}
		*this = DateTime(t)
		return nil
	}
	df := "2006-01-02"
	//根据输入格式进行解析
	d := "^[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}$"
	if isd, _ := regexp.MatchString(d, s); isd {
		t, err := time.Parse(df, s)
		if err != nil {
			return err
		}
		*this = DateTime(t)
		return nil
	}
	tf := "15:04"
	//根据输入格式进行解析
	ts := "^[0-9]{1,2}:[0-9]{1,2}$"
	if ist, _ := regexp.MatchString(ts, s); ist {
		t, err := time.Parse(tf, s)
		if err != nil {
			return err
		}
		*this = DateTime(t)
		return nil
	}
	return errors.New(fmt.Sprintf("Not konw formart for value %s", s))
}
func (this DateTime) DateStr() string {
	return time.Time(this).Format("2006-01-02")
}
func (this DateTime) TimeStr() string {
	return time.Time(this).Format("15:04:05")
}
func (this DateTime) Str() string {
	return time.Time(this).Format("2006-01-02 15:04:05")
}

func (this DateTime) Value() (driver.Value, error) {
	return time.Time(this), nil
}

func (this *DateTime) Scan(value interface{}) error {
	tv, ok := value.(time.Time)
	if ok {
		*this = DateTime(tv)
	}
	return errors.New(fmt.Sprintf("Not suport type %t", value))
}

//-----------------------------------------------------------------------

//Date 日期格式，仅包含日期部分
type Date time.Time

//UnmarshalJSON 实现Date的json反序列化方法
func (this *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return this.Parse(s)
}

//MarshalJSON 实现Date的json序列化方法
func (this Date) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02"))
	return []byte(stamp), nil
}

func (this Date) Value() (driver.Value, error) {
	return time.Time(this), nil
}

func (this *Date) Scan(value interface{}) error {
	tv, ok := value.(time.Time)
	if ok {
		*this = Date(tv)
	}
	return errors.New(fmt.Sprintf("Not suport type %t", value))
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
func (this Date) DateStr() string {
	return time.Time(this).Format("2006-01-02")
}

//--------------------------------------------------------------------------

// STime short time
type STime time.Time

//MarshalJSON 实现DateTime的json序列化方法
func (this STime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("15:04:05"))
	return []byte(stamp), nil
}

func (this *STime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return this.Parse(s)
}

func (this STime) Value() (driver.Value, error) {
	return time.Time(this), nil
}

func (this *STime) Scan(value interface{}) error {
	tv, ok := value.(time.Time)
	if ok {
		*this = STime(tv)
	}
	return errors.New(fmt.Sprintf("Not suport type %t", value))
}

func (this *STime) Parse(s string) error {
	tf := "15:04"
	//根据输入格式进行解析
	ts := "^[0-9]{1,2}:[0-9]{1,2}$"
	if ist, _ := regexp.MatchString(ts, s); ist {
		t, err := time.Parse(tf, s)
		if err != nil {
			return err
		}
		*this = STime(t)
		return nil
	}
	return errors.New(fmt.Sprintf("Not konw formart for value %s", s))
}

func (this STime) TimeStr() string {
	return time.Time(this).Format("15:04:05")
}

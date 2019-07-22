package dbmeta

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type Employee struct {
	EmpNo     int       `gorm:"column:emp_no;primary_key" json:"emp_no"`
	BirthDate time.Time `gorm:"column:birth_date" json:"birth_date"`
	FirstName string    `gorm:"column:first_name" json:"first_name"`
	LastName  string    `gorm:"column:last_name" json:"last_name"`
	Gender    string    `gorm:"column:gender" json:"gender"`
	HireDate  time.Time `gorm:"column:hire_date" json:"hire_date"`
}

func Test_Copy(t *testing.T) {
	now := time.Now()
	dst := &Employee{
		EmpNo:     10001,
		BirthDate: now,
		FirstName: "Tom",
	}

	src := &Employee{
		EmpNo:     10001,
		BirthDate: now.Add(3600 * time.Second),
		FirstName: "Jerry",
		Gender:    "Male",
	}

	err := Copy(dst, src)
	if err != nil {
		t.Fatal(err)
	}

	expected := &Employee{
		EmpNo:     10001,
		BirthDate: now.Add(3600 * time.Second),
		FirstName: "Jerry",
		Gender:    "Male",
	}

	if !reflect.DeepEqual(expected, dst) {
		t.Errorf("expect: %+v, but got %+x", expected, dst)
	}
}

func TestStrFirstToLower(t *testing.T) {
	str := "Card"
	tar := "card"
	res := StrFirstToLower(str)
	fmt.Printf("%v strFirstToLower is %v \n", str, res)
	str = "ID"
	tar = "ID"
	res = StrFirstToLower(str)
	fmt.Printf("%v strFirstToLower is %v \n", str, res)
	if res != tar {
		t.Error("Not equal")
	}
}

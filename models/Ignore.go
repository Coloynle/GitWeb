package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Ignore struct {
	Id     int `orm:"auto;pk;unique;"`
	Name   string
	Path   string `orm:"unique"`
	Status int
}

func (i *Ignore) GetAll(limit int) []Ignore {
	var all []Ignore
	o := orm.NewOrm()
	o.QueryTable(i).Limit(limit).OrderBy("id").All(&all)
	return all
}


func (i *Ignore) IsExist(where orm.Params) bool {
	o := orm.NewOrm()
	qs := o.QueryTable(i)
	for k, v := range where {
		qs = qs.Filter(k, v)
	}
	err := qs.One(i)
	if err == orm.ErrNoRows {
		return  false
	} else {
		return  true
	}
}

//
func (i *Ignore) Insert(ignore Ignore) (int64, error) {
	var id int64
	var err error
	o := orm.NewOrm()
	isExist := i.IsExist(orm.Params{"path": ignore.Path})

	if !isExist {
		id, err = o.Insert(&ignore)
	} else {
		id, err = i.Update(orm.Params{"path": ignore.Path}, orm.Params{"status": ignore.Status})
	}
	if err == nil {
		return id, err
	} else {
		fmt.Println(err)
		return -1, err
	}
}

// 批量更新
func (i *Ignore) Update(where orm.Params, set orm.Params) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(i)
	for k, v := range where {
		qs = qs.Filter(k, v)
	}
	succ, err := qs.Update(set)
	if err == nil {
		return succ, err
	} else {
		return -1, err
	}
}

// 批量删除
func (i *Ignore) Delete(where orm.Params) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(i)
	for k, v := range where {
		qs = qs.Filter(k, v)
	}
	succ, err := qs.Delete()
	if err == nil {
		return succ, err
	} else {
		return -1, err
	}
}

//
// // 查询全部数据
// func (i *Ignore) SelectAll(Field []string) []Ignores {
// 	baseModel := GetInstance(&MysqlConnPool{})
// 	var nIgnore []Ignore
// 	var nIgnores []Ignores
// 	baseModel.Select("Find",&nIgnore,Field,nil,nil,nil).Scan(&nIgnores)
// 	return nIgnores
// }
//
// // 插入数据
// func (i *Ignore) Add() uint {
// 	baseModel := GetInstance(&MysqlConnPool{})
// 	// 调用 公用Add方法获取 true 插入成功
// 	if baseModel.Add(&i){
// 		return i.ID
// 	}else {
// 		return 0
// 	}
// }
//
// func (i *Ignore) Update(where interface{}) bool{
// 	baseModel := GetInstance(&MysqlConnPool{})
// 	return baseModel.Update(&i,where)
// }
//
// type Ignores struct {
// 	Id     int    `json:"id"`
// 	Name   string `json:"name"`
// 	Path   string `json:"path"`
// 	Status int   `json:"status"`
// }

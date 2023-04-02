package service

import (
	"sort"
	"strings"
	"time"
)

type dataSelector struct {
	GenericDataList []DataCell
	dataSelectQuery *DataSelectQuery
}

type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

func (d *dataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return a.Before(b)
}

func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// 过滤
func (d *dataSelector) Filter() *dataSelector {
	if d.dataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	filterList := []DataCell{}
	for _, value := range d.GenericDataList {
		matches := true
		objName := value.GetName()
		if !strings.Contains(objName, d.dataSelectQuery.FilterQuery.Name) {
			matches = false
			continue
		}
		if matches {
			filterList = append(filterList, value)
		}
	}
	d.GenericDataList = filterList
	return d
}

// 分页
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.dataSelectQuery.PaginateQuery.Limit
	page := d.dataSelectQuery.PaginateQuery.Page
	if limit <= 0 || page <= 0 {
		return d
	}
	startIndex := limit * (page - 1)
	endIndex := limit * page
	if len(d.GenericDataList) < startIndex {
		startIndex = 0
	}
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

package dbHelper

import (
	"errors"
	"fmt"

	"pr-reviewer/internal/model"
	"pr-reviewer/internal/utils/errorUtils"

	"strings"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
)

type QueryOptions struct {
	TableName          string
	TableAlias         string
	TableColumns       []string
	ExtraSelect        []string
	DefaultWhereCases  []string
	DefaultSort        []string
	JoinTableList      []JoinTable
	TargetObject       interface{}
	QuickFilterColumns []string
	IsTotalQuery       bool
}

type JoinTable struct {
	TableName    string
	TableAlias   string
	TableColumns []string
	JoinExpr     []string //ex. cu.user_id = ch.user_id
}

func CreateQuery(param model.ListDataParameters, options QueryOptions) (string, []interface{}, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	return CreateQueryWithBuilder(sb, param, options)
}

func CreateQueryWithBuilder(sb *sqlbuilder.SelectBuilder, param model.ListDataParameters, options QueryOptions) (string, []interface{}, error) {
	isTotalQuery := options.IsTotalQuery || (len(options.TableColumns) == 1 && options.TableColumns[0] == "count(*)")

	selectColumns := make([]string, 0)

	if !isTotalQuery {

		alias := options.TableName
		if options.TableAlias != "" {
			alias = options.TableAlias
		}
		//Add Alias before columns names
		for _, column := range options.TableColumns {
			selectColumns = append(selectColumns, fmt.Sprintf("%s.%s", alias, column))
		}

		//Add Alias before columns names for join tables columns
		for _, joinTable := range options.JoinTableList {
			alias := joinTable.TableName
			if joinTable.TableAlias != "" {
				alias = joinTable.TableAlias
			}
			for _, column := range joinTable.TableColumns {
				selectColumns = append(selectColumns, fmt.Sprintf("%s.%s", alias, column))
			}
		}

		//Add Extra Select Columns if any
		if len(options.ExtraSelect) > 0 {
			selectColumns = append(selectColumns, options.ExtraSelect...)
		}

	} else {
		// Total Query ignores TableColumns and just returns the total only.
		selectColumns = []string{"count(*)"}
	}

	for _, joinTable := range options.JoinTableList {
		sb.JoinWithOption(sqlbuilder.LeftJoin, joinTable.TableName, joinTable.JoinExpr...)
	}

	sb.Select(selectColumns...)

	//add Alias to table name
	sb.From(options.TableName)

	//Apply default Where Cases
	if len(options.DefaultWhereCases) > 0 {
		sb.Where(options.DefaultWhereCases...)
	}

	//Apply Filters from params
	for _, filterItem := range param.Filter.Items {
		if err := applyFilter(sb, options, filterItem); err != nil {
			return "", nil, err
		}
	}

	//Apply Quick Filters from params
	for _, quickFilterValue := range param.Filter.QuickFilterValues {
		if err := applyQuickFilter(sb, options, quickFilterValue); err != nil {
			return "", nil, err
		}
	}

	if param.Filter.InDateRange != nil {
		if err := applyInDateRangeFilter(sb, options, *param.Filter.InDateRange); err != nil {
			return "", nil, err
		}
	}

	//If the only column was passed is count(*) we build query and return
	//since this means we want to know total record number
	if isTotalQuery {
		sql, args := sb.Build()
		return sql, args, nil
	}

	//Apply default sort order if no sort parameters provided
	if len(param.Sort) == 0 {
		sb.OrderBy(options.DefaultSort...)
	}

	//Apply sort order from param
	for _, sortItem := range param.Sort {
		fieldName, _ := GetDBFieldNameByJsonName(options.TargetObject, sortItem.Field)
		if fieldName == "" {
			return "", nil, errorUtils.Wrap(errors.New("error: sort by field name not found"))
		}
		sb.OrderBy(fmt.Sprintf("%s %s NULLS LAST", applyAlias(options, fieldName), sortItem.Sort))
	}

	//add Limit and Offset from params
	if param.Page.Size > 0 {
		sb.Limit(param.Page.Size)
		sb.Offset(param.Page.Size * param.Page.Number)
	}

	sql, args := sb.Build()

	return sql, args, nil
}

func applyFilter(sb *sqlbuilder.SelectBuilder, options QueryOptions, filterItem model.FilterItemData) error {
	fieldName, fieldType := GetDBFieldNameByJsonName(options.TargetObject, filterItem.Field)
	if fieldName == "" {
		return errorUtils.Wrap(errors.New("error: field name not found"))
	}

	fieldValue, err := ParseValueToFieldType(fieldType, filterItem.Value)
	if err != nil {
		return errorUtils.Wrap(err)
	}

	finalColumnName := applyAlias(options, fieldName)

	switch filterItem.OperatorValue {
	case "isEmpty":
		sb.Where(sb.IsNull(finalColumnName))
	case "isNotEmpty":
		sb.Where(sb.IsNotNull(finalColumnName))
	case "contains":
		sb.Where(sb.Like(fmt.Sprintf("LOWER(%s::text)", finalColumnName), fmt.Sprintf("%%%s%%", strings.ToLower(fmt.Sprintf("%v", fieldValue)))))
	case "equals":
		sb.Where(sb.Equal(fmt.Sprintf("LOWER(%s::text)", finalColumnName), strings.ToLower(fmt.Sprintf("%v", fieldValue))))
	case "=", "is":
		sb.Where(sb.Equal(finalColumnName, fieldValue))
	case "!=", "not":
		sb.Where(sb.NotEqual(finalColumnName, fieldValue))
	case "startsWith":
		sb.Where(sb.Like(fmt.Sprintf("LOWER(%s::text)", finalColumnName), fmt.Sprintf("%s%%", strings.ToLower(fmt.Sprintf("%v", fieldValue)))))
	case "endsWith":
		sb.Where(sb.Like(fmt.Sprintf("LOWER(%s::text)", finalColumnName), fmt.Sprintf("%%%s", strings.ToLower(fmt.Sprintf("%v", fieldValue)))))
	case ">", "after":
		sb.Where(sb.GreaterThan(finalColumnName, fieldValue))
	case ">=", "onOrAfter":
		sb.Where(sb.GreaterEqualThan(finalColumnName, fieldValue))
	case "<", "before":
		sb.Where(sb.LessThan(finalColumnName, fieldValue))
	case "<=", "onOrBefore":
		sb.Where(sb.LessEqualThan(finalColumnName, fieldValue))
	}

	return nil
}

func applyQuickFilter(sb *sqlbuilder.SelectBuilder, options QueryOptions, filterValue string) error {
	whereCases := make([]string, 0)
	for _, quickFilterColumn := range options.QuickFilterColumns {
		columnName := applyAlias(options, quickFilterColumn)
		whereCases = append(whereCases, sb.Like(fmt.Sprintf("LOWER(%s::text)", columnName), fmt.Sprintf("%%%s%%", strings.ToLower(fmt.Sprintf("%v", filterValue)))))
	}
	sb.Where(sb.Or(whereCases...))
	return nil
}

func applyInDateRangeFilter(sb *sqlbuilder.SelectBuilder, options QueryOptions, dateRangeValues model.DateRange) error {
	whereCases := make([]string, 0)

	startAtFieldName, startAtFieldType := GetDBFieldNameByJsonName(options.TargetObject, dateRangeValues.StartAt.Field)
	if startAtFieldName == "" {
		return errorUtils.Wrap(errors.New("error: field name not found"))
	}

	startAtFieldValue, err := ParseValueToFieldType(startAtFieldType, dateRangeValues.StartAt.Value)
	if err != nil {
		return errorUtils.Wrap(err)
	}

	endAtFieldName, _ := GetDBFieldNameByJsonName(options.TargetObject, dateRangeValues.EndAt.Field)
	if endAtFieldName == "" {
		return errorUtils.Wrap(errors.New("error: field name not found"))
	}

	endAtFieldValue, err := ParseValueToFieldType(startAtFieldType, dateRangeValues.EndAt.Value)
	if err != nil {
		return errorUtils.Wrap(err)
	}

	startAtColumn := applyAlias(options, startAtFieldName)
	endAtAtColumn := applyAlias(options, endAtFieldName)
	whereCases = append(whereCases, sb.And(sb.GreaterThan(startAtColumn, startAtFieldValue), sb.LessThan(startAtColumn, endAtFieldValue)))
	whereCases = append(whereCases, sb.And(sb.GreaterThan(endAtAtColumn, startAtFieldValue), sb.LessThan(endAtAtColumn, endAtFieldValue)))
	whereCases = append(whereCases, sb.And(sb.LessThan(startAtColumn, startAtFieldValue), sb.GreaterThan(endAtAtColumn, endAtFieldValue)))

	sb.Where(sb.Or(whereCases...))
	return nil
}

// Given an sql column name, figure out which table it is from and prepend that
// table name; this helps resolve ambiguities when joining with tables that
// may have fields with the same name.
//
// For example: given "id", this may return "pets.id".
func applyAlias(options QueryOptions, targetColumnName string) string {

	// Try to find if the column is in the main table.
	for _, column := range options.TableColumns {
		if column == targetColumnName {
			if options.TableAlias != "" {
				return fmt.Sprintf("%s.%s", options.TableAlias, targetColumnName)
			}

			return fmt.Sprintf("%s.%s", options.TableName, targetColumnName)
		}
	}

	// We didn't find column in the main table; look in join tables.
	for _, joinTable := range options.JoinTableList {
		for _, column := range joinTable.TableColumns {
			if column == targetColumnName {
				if joinTable.TableAlias != "" {
					return fmt.Sprintf("%s.%s", joinTable.TableAlias, targetColumnName)
				}
				return fmt.Sprintf("%s.%s", joinTable.TableName, targetColumnName)
			}
		}
	}

	// We still didn't find the column. Default to assuming it was given a
	// special name such as by "id AS organization_id".
	return targetColumnName
}

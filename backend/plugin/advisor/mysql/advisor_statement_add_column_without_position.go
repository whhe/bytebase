package mysql

// Framework code is generated by the generator.

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/pkg/errors"

	mysql "github.com/bytebase/mysql-parser"

	"github.com/bytebase/bytebase/backend/plugin/advisor"
	mysqlparser "github.com/bytebase/bytebase/backend/plugin/parser/mysql"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
)

var (
	_ advisor.Advisor = (*StatementAddColumnWithoutPositionAdvisor)(nil)
)

func init() {
	advisor.Register(storepb.Engine_OCEANBASE, advisor.MySQLStatementAddColumnWithoutPosition, &StatementAddColumnWithoutPositionAdvisor{})
}

// StatementAddColumnWithoutPositionAdvisor is the advisor checking for checking no position in ADD COLUMN clause.
type StatementAddColumnWithoutPositionAdvisor struct {
}

// Check checks for checking no position in ADD COLUMN clause.
func (*StatementAddColumnWithoutPositionAdvisor) Check(ctx advisor.Context, _ string) ([]advisor.Advice, error) {
	stmtList, ok := ctx.AST.([]*mysqlparser.ParseResult)
	if !ok {
		return nil, errors.Errorf("failed to convert to mysql parser result")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &statementAddColumnWithoutPositionChecker{
		level: level,
		title: string(ctx.Rule.Type),
	}

	for _, stmt := range stmtList {
		checker.baseLine = stmt.BaseLine
		antlr.ParseTreeWalkerDefault.Walk(checker, stmt.Tree)
	}

	if len(checker.adviceList) == 0 {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		})
	}
	return checker.adviceList, nil
}

type statementAddColumnWithoutPositionChecker struct {
	*mysql.BaseMySQLParserListener

	baseLine   int
	adviceList []advisor.Advice
	level      advisor.Status
	title      string
}

func (checker *statementAddColumnWithoutPositionChecker) EnterAlterTable(ctx *mysql.AlterTableContext) {
	if !mysqlparser.IsTopMySQLRule(&ctx.BaseParserRuleContext) {
		return
	}
	if ctx.AlterTableActions() == nil {
		return
	}
	if ctx.AlterTableActions().AlterCommandList() == nil {
		return
	}
	if ctx.AlterTableActions().AlterCommandList().AlterList() == nil {
		return
	}
	if ctx.TableRef() == nil {
		return
	}

	_, tableName := mysqlparser.NormalizeMySQLTableRef(ctx.TableRef())
	if tableName == "" {
		return
	}

	for _, item := range ctx.AlterTableActions().AlterCommandList().AlterList().AllAlterListItem() {
		if item == nil || item.ADD_SYMBOL() == nil {
			continue
		}

		var position string

		switch {
		case item.Identifier() != nil && item.FieldDefinition() != nil:
			position = getPosition(item.Place())
		case item.OPEN_PAR_SYMBOL() != nil && item.TableElementList() != nil:
			for _, tableElement := range item.TableElementList().AllTableElement() {
				if tableElement.ColumnDefinition() == nil {
					continue
				}
				if tableElement.ColumnDefinition().FieldDefinition() == nil {
					continue
				}

				position = getPosition(item.Place())
				if len(position) != 0 {
					break
				}
			}
		}

		if len(position) != 0 {
			checker.adviceList = append(checker.adviceList, advisor.Advice{
				Status:  checker.level,
				Code:    advisor.StatementAddColumnWithPosition,
				Title:   checker.title,
				Content: fmt.Sprintf("add column with position \"%s\"", position),
				Line:    checker.baseLine + ctx.GetStart().GetLine(),
			})
		}
	}
}

func getPosition(ctx mysql.IPlaceContext) string {
	if ctx == nil {
		return ""
	}
	place, ok := ctx.(*mysql.PlaceContext)
	if !ok || place == nil {
		return ""
	}

	switch {
	case place.FIRST_SYMBOL() != nil:
		return "FIRST"
	case place.AFTER_SYMBOL() != nil:
		return "AFTER"
	default:
		return ""
	}
}

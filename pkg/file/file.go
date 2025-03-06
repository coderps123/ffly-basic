package file

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// ColumnConfig 列设置
type ColumnConfig struct {
	Title  string  // 标题
	Field  string  // 字段名
	Width  float64 // 宽度
	Prefix string  // 前缀 (用于自定义前缀)
}

// Options 选项
type Options struct {
	Child      string // 子字段（用于树形结构）
	FileSuffix string // 文件后缀名
}

// setHeaders 设置表头
func setHeaders(file *excelize.File, sheetName string, columns []ColumnConfig) error {
	for i, column := range columns {
		// 获取列名（A, B, C...）
		colName := string(rune('A' + i))    // 获取列名（A, B, C...）
		cell := fmt.Sprintf("%s1", colName) // 设置单元格位置 A1, B1, C1...

		// 设置单元格值
		if err := file.SetCellValue(sheetName, cell, column.Title); err != nil {
			return fmt.Errorf("设置单元格值失败:%v", err)
		}

		// 设置列宽（使用列名而不是单元格引用）
		if err := file.SetColWidth(sheetName, colName, colName, column.Width); err != nil {
			return fmt.Errorf("设置列宽失败:%v", err)
		}
	}

	return nil
}

// writeData 写入数据
func writeData(file *excelize.File, sheetName string, row *int, level int, dataValue reflect.Value, columns []ColumnConfig, options Options) error {
	for rowIndex := 0; rowIndex < dataValue.Len(); rowIndex++ {
		item := dataValue.Index(rowIndex) // 获取数据项
		if item.Kind() == reflect.Ptr {   // 如果是指针类型，则取指针指向的值
			item = item.Elem() // 取指针指向的值
		}

		// 写入每一列数据
		for colIndex, column := range columns {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, *row) // A2 A3 A4...

			// 通过反射获取字段值
			field := item.FieldByName(column.Field) // 获取字段值
			if !field.IsValid() {
				continue // 字段不存在，跳过
			}

			// 根据层级重复前缀
			prefix := ""
			if column.Prefix != "" { // 如果需要添加前缀
				prefix = strings.Repeat(column.Prefix, level) // 重复前缀
			}

			// 获取字段值并处理前缀
			var value interface{}
			if field.Kind() == reflect.String { // 如果是字符串类型
				value = prefix + field.String() // 添加前缀
			} else {
				value = field.Interface() // 直接获取值
			}

			// 设置单元格值
			if err := file.SetCellValue(sheetName, cell, value); err != nil {
				return fmt.Errorf("设置单元格值失败:%v", err)
			}
		}

		*row++ // 行号递增

		// 设置子节点
		var child string
		if options.Child != "" { // 如果子节点字段存在
			child = options.Child // 子节点字段值
		} else { // 子节点字段值为空，默认使用 Children
			child = "Children"
		}

		children := item.FieldByName(child)           // 获取子节点
		if children.IsValid() && children.Len() > 0 { // 如果子节点存在且长度大于0
			if err := writeData(file, sheetName, row, level+1, children, columns, options); err != nil {
				return fmt.Errorf("写入子节点失败:%v", err)
			}
		}
	}

	return nil
}

// GenerateExcelAndReturnBytes 生成excel并返回字节流
func GenerateExcelAndReturnBytes(data interface{}, columns []ColumnConfig, options Options) ([]byte, error) {
	// 创建一个新的excel文件
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("关闭Excel文件失败：%v\n", err)
		}
	}()

	// 创建一个工作表 sheet1
	sheetName := "sheet1"
	sheetIndex, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("创建工作表失败：%v", err)
	}
	file.SetActiveSheet(sheetIndex) // 设置当前工作表

	// 设置表头
	if err := setHeaders(file, sheetName, columns); err != nil {
		return nil, fmt.Errorf("设置表头失败：%v", err)
	}

	// 获取数据切片
	dataSlice := reflect.ValueOf(data)
	if dataSlice.Kind() != reflect.Slice {
		return nil, fmt.Errorf("数据必须是切片类型")
	}

	row := 2 // 起始行索引

	// 写入数据
	if err := writeData(file, sheetName, &row, 0, dataSlice, columns, options); err != nil {
		return nil, err
	}

	// 将文件写入到内存
	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("excel写入到内存失败：%v", err)
	}

	return buffer.Bytes(), nil
}

// ExportExcel 导出excel文件
func ExportExcel(c *gin.Context, data interface{}, columns []ColumnConfig, filename string, options Options) error {
	bytes, err := GenerateExcelAndReturnBytes(data, columns, options) // 生成excel文件并返回字节流
	if err != nil {
		return fmt.Errorf("生成excel文件失败：%v", err)
	}

	// 如果没有指定文件后缀，默认为xlsx
	if options.FileSuffix == "" {
		options.FileSuffix = "xlsx"
	}

	// 对文件名进行 URL 编码
	encodedFilename := url.QueryEscape(filename + "." + options.FileSuffix)

	// "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" 是excel文件的 MIME 类型
	c.Writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"; filename*=utf-8''%s", encodedFilename, encodedFilename))
	c.Writer.Header().Set("Content-Transfer-Encoding", "binary")

	_, err = c.Writer.Write(bytes) // 将字节流写入响应
	if err != nil {
		return fmt.Errorf("写入响应失败：%v", err)
	}

	return nil
}

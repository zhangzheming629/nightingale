package excel

import (
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"regexp"
	"strconv"
	"strings"
)

// 定义正则表达式模式
const (
	ExcelTagKey = "excel"
	Pattern     = "name:(.*?);|index:(.*?);|width:(.*?);|needMerge:(.*?);|replace:(.*?);"
)

type ExcelTag struct {
	Value     interface{}
	Name      string // 表头标题
	Index     int    // 列下标(从0开始)
	Width     int    // 列宽
	NeedMerge bool   // 是否需要合并
	Replace   string // 替换（需要替换的内容_替换后的内容。比如：1_未开始 ==> 表示1替换为未开始）
}

// 构造函数，返回一个带有默认值的 ExcelTag 实例
func NewExcelTag() ExcelTag {
	return ExcelTag{
		// 导入时会根据这个下标来拿单元格的值，当目标结构体字段没有设置index时，
		// 解析字段tag值时Index没读到就一直默认为0，拿单元格的值时，就始终拿的是第一列的值
		Index: -1, // 设置 Index 的默认值为 -1
	}
}

// 读取字段tag值
func (e *ExcelTag) GetTag(tag string) (err error) {
	// 编译正则表达式
	re := regexp.MustCompile(Pattern)
	matches := re.FindAllStringSubmatch(tag, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			for i, val := range match {
				if i != 0 && val != "" {
					e.setValue(match, val)
				}
			}
		}
	} else {
		err = errors.New("未匹配到值")
		return
	}
	return
}

// 设置ExcelTag 对应字段的值
func (e *ExcelTag) setValue(tag []string, value string) {
	if strings.Contains(tag[0], "name") {
		e.Name = value
	}
	if strings.Contains(tag[0], "index") {
		v, _ := strconv.ParseInt(value, 10, 8)
		e.Index = int(v)
	}
	if strings.Contains(tag[0], "width") {
		v, _ := strconv.ParseInt(value, 10, 8)
		e.Width = int(v)
	}
	if strings.Contains(tag[0], "needMerge") {
		v, _ := strconv.ParseBool(value)
		e.NeedMerge = v
	}
	if strings.Contains(tag[0], "replace") {
		e.Replace = value
	}
}



type Excel struct {
	F             *excelize.File // excel 对象
	TitleStyle    int            // 表头样式
	HeadStyle     int            // 表头样式
	ContentStyle1 int            // 主体样式1，无背景色
	ContentStyle2 int            // 主体样式2，有背景色
}

// 初始化
func ExcelInit() (e *Excel) {
	e = &Excel{}
	// excel构建
	e.F = excelize.NewFile()
	// 初始化样式
	e.getTitleRowStyle()
	e.getHeadRowStyle()
	e.getDataRowStyle()
	return e
}

// 获取边框样式
func getBorder() []excelize.Border {
	return []excelize.Border{ // 边框
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "left", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}
}

// 标题样式
func (e *Excel) getTitleRowStyle() {
	e.TitleStyle, _ = e.F.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{ // 对齐方式
			Horizontal: "center", // 水平对齐居中
			Vertical:   "center", // 垂直对齐居中
		},
		Fill: excelize.Fill{ // 背景颜色
			Type:    "pattern",
			Color:   []string{"#fff2cc"},
			Pattern: 1,
		},
		Font: &excelize.Font{ // 字体
			Bold: true,
			Size: 16,
		},
		Border: getBorder(),
	})
}

// 列头行样式
func (e *Excel) getHeadRowStyle() {
	e.HeadStyle, _ = e.F.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{ // 对齐方式
			Horizontal: "center", // 水平对齐居中
			Vertical:   "center", // 垂直对齐居中
			WrapText:   true,     // 自动换行
		},
		Fill: excelize.Fill{ // 背景颜色
			Type:    "pattern",
			Color:   []string{"#FDE9D8"},
			Pattern: 1,
		},
		Font: &excelize.Font{ // 字体
			Bold: true,
			Size: 14,
		},
		Border: getBorder(),
	})
}

// 数据行样式
func (e *Excel) getDataRowStyle() {
	style := excelize.Style{}
	style.Border = getBorder()
	style.Alignment = &excelize.Alignment{
		Horizontal: "center", // 水平对齐居中
		Vertical:   "center", // 垂直对齐居中
		WrapText:   true,     // 自动换行
	}
	style.Font = &excelize.Font{
		Size: 12,
	}
	e.ContentStyle1, _ = e.F.NewStyle(&style)
	style.Fill = excelize.Fill{ // 背景颜色
		Type:    "pattern",
		Color:   []string{"#cce7f5"},
		Pattern: 1,
	}
	e.ContentStyle2, _ = e.F.NewStyle(&style)
}


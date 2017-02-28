/*

*/
package radius

import(
	"fmt"
	//	gr "github.com/blind-oracle/go-radius"
	gr "github.com/cuu/radius"
	"github.com/cuu/softradius/libs"
)


///cuu radius 默认是放在vendor=default的中
// 读入 freeradius 的dictionary 用Vendor 来区别存储 attrs


func init() {
	fmt.Println(libs.GetCurrTimeNano())
	gr.Builtin.MustRegister("Acct-Interim-Interval", 85, gr.AttributeInteger)
	gr.Builtin.LoadDicts("radius/dicts/dictionary")
	fmt.Println("Custom dict loading..", libs.GetCurrTimeNano() )	
}



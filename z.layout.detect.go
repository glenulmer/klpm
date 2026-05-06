package main

import (
	"net/http"

	. "klpm/lib/output"
)

func LayoutFromDeviceMode(mode string) string {
	if mode == deviceMobile { return layoutPhone }
	return layoutDesktop
}

func RequestLayout(r *http.Request) string {
	return LayoutFromDeviceMode(SessionDeviceMode(r))
}

func DeviceModeFromLayout(layout string) string {
	if layout == layoutPhone { return deviceMobile }
	return deviceDesktop
}

func DeviceConfirmHeadScript(mode string) string {
	x := deviceDesktop
	if mode0, ok := NormalizeDeviceMode(mode); ok { x = mode0 }
	return Str(
		`(function(){`,
		`var s="`, x, `";`,
		`var d=(window.innerWidth<900)?"mobile":"desktop";`,
		`if(d===s){return;}`,
		`if(location.search.indexOf("layout=")>=0){return;}`,
		`var q=location.search?location.search.slice(1):"";`,
		`var sep=q?"&":"";`,
		`var lay=d==="mobile"?"phone":"desktop";`,
		`var u=location.pathname+"?"+q+sep+"layout="+lay+location.hash;`,
		`location.replace(u);`,
		`})();`,
	)
}

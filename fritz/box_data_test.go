package fritz

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_boxDataParser_parse(t *testing.T) {
	tests := []struct {
		name string
		arg  string
	}{
		{name: "Fon ata 1020", arg: "FRITZ!Box Fon ata 1020-B-060102-040410-262336-724176-787902-110401-2956-avme-en"},
		{name: "Fon WLAN UI", arg: "FRITZ!Box Fon WLAN (UI)-B-070608-050527-151063-152677-787902-080434-7804-1und1"},
		{name: "Fon WLAN 7170", arg: "FRITZ!Box Fon WLAN 7170-B-171008-000028-457563-147110-787902-290487-19985-avm"},
		{name: "Fon WLAN 7170", arg: "FRITZ!Box Fon WLAN 7170 (UI)-B-171607-041025-521747-370532-147902-290486-19138-1und1"},
		{name: "Fon WLAN 7170", arg: "FRITZ!Box Fon WLAN 7170 Annex A-A-042602-030325-402401-042265-787902-580482"},
		{name: "Fon WLAN 7270", arg: "FRITZ!Box Fon WLAN 7270-B-111700-030716-306463-160202-787902-540480-15918-avm"},
		{name: "Fon WLAN 7270", arg: "FRITZ!Box Fon WLAN 7270 v2-B-071710-020026-055200-026256-147902-540504-20260-avme-en"},
		{name: "Fon WLAN 7340", arg: "FRITZ!Box Fon WLAN 7340-B-072706-000217-533416-737311-147902-990505-20608-avme-en"},
		{name: "Fon WLAN 7390", arg: "FRITZ!Box Fon WLAN 7390-B-161001-000109-670300-000324-787902-840507-21400-avm"},
		{name: "Fon WLAN 7390", arg: "FRITZ!Box Fon WLAN 7390 (UI)-B-171408-010206-436146-332654-147902-840505-20359-1und1"},
		{name: "Fon WLAN 7390", arg: "FRITZ!Box Fon WLAN 7390-A-091301-000006-002567-355146-787902-840505-20608-avme-de"},
		{name: "Repeater NG", arg: "FRITZ!WLAN Repeater N/G-B-152108-020222-445034-614506-787902-680486-20648-avm"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := boxDataParser{}
			data := p.parse(tt.arg)
			assert.NotZero(t, data)
			fmt.Println("Model", data.Model)
			fmt.Println("Version", data.FirmwareVersion)
			fmt.Println("Runtime", data.Runtime)
			fmt.Println("----")
		})
	}
}

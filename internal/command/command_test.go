package command

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"syscall"
	"testing"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"

	"github.com/metalsoft-io/metalcloud-cli/internal/configuration"
	. "github.com/onsi/gomega"
)

func TestSimpleArgument(t *testing.T) {

	var executed = false

	cmd := Command{
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "create",
		AltPredicate: "c",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_instance_count": c.FlagSet.Int("instance_count", 0, "Instance count of this instance array"),
				"instance_array_instance_label": c.FlagSet.String("label", "", "Instance array's label"),
			}
		},
		ExecuteFunc: func(c *Command, client metalcloud.MetalCloudClient) (string, error) {
			executed = true
			return "retstr", nil
		},
	}

	cmd.InitFunc(&cmd)

	argv := []string{
		"-instance_count=3",
		"-label=test",
	}

	err := cmd.FlagSet.Parse(argv)
	if err != nil {
		t.Errorf("%s", err)
	}

	iaCount := cmd.Arguments["instance_array_instance_count"].(*int)
	if iaCount == nil || *iaCount != 3 {
		t.Errorf("instance_array_instance_count expected to be %d\n\twas %d", 3, *iaCount)
	}

	iaLabel := cmd.Arguments["instance_array_instance_label"].(*string)

	if iaLabel == nil || *iaLabel != "test" {
		t.Errorf("instance_array_label expected to be %s\n\twas %s", "test", *iaLabel)
	}

	argv = []string{
		"instance_countasdad=3",
		"la33bel=\"test\"",
	}

	err = cmd.FlagSet.Parse(argv)
	if err != nil {
		t.Errorf("%s", err)
	}

	ret, err := cmd.ExecuteFunc(&cmd, nil)
	if err != nil {
		t.Errorf("%s", err)
	}

	if !executed || ret != "retstr" {
		t.Errorf("ExecuteFunction not called properly")
	}
}

func TestConfirmFunc(t *testing.T) {
	RegisterTestingT(t)
	v1 := 10

	cmd := Command{
		Subject:      "instance_array",
		AltSubject:   "ia",
		Predicate:    "create",
		AltPredicate: "c",
		FlagSet:      flag.NewFlagSet("instance_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"autoconfirm": c.FlagSet.Bool("autoconfirm", false, "autoconfirm text"),
			}
		},
		ExecuteFunc: func(c *Command, client metalcloud.MetalCloudClient) (string, error) {
			return "", nil
		},
	}

	cmd.InitFunc(&cmd)

	var stdin bytes.Buffer
	var stdout bytes.Buffer

	configuration.SetConsoleIOChannel(&stdin, &stdout)

	stdin.Write([]byte("yes\n"))

	//check without autoconfirm
	ok, err := ConfirmCommand(&cmd,
		func() string {
			return fmt.Sprintf("Reverting infrastructure %d to the deployed state. Are you sure? Type \"yes\" to continue:", v1)
		},
	)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())

	s, err := stdout.ReadString(byte('\n'))
	Expect(s).To(ContainSubstring("Reverting infrastructure"))

	//check with autoconfirm
	argv := []string{
		"-autoconfirm",
	}

	err = cmd.FlagSet.Parse(argv)
	Expect(err).To(BeNil())

	ok, err = ConfirmCommand(&cmd,
		func() string {
			return fmt.Sprintf("Reverting infrastructure %d to the deployed state. Are you sure? Type \"yes\" to continue:", v1)
		},
	)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())

	s, err = stdout.ReadString(byte('\n'))
	Expect(s).To(BeEmpty())

}

func TestGetIfNotDefaultOk(t *testing.T) {
	RegisterTestingT(t)

	i := 10
	s := "test"
	f := 10.2
	m := map[string]interface{}{
		"testInt":       &i,
		"testString":    &s,
		"testWrongType": &f,
	}

	v, ok := getPtrValueIfExistsOk(m, "testInt")
	Expect(v).To(Equal((10)))
	Expect(ok).To(BeTrue())

	v, ok = getPtrValueIfExistsOk(m, "testString")
	Expect(v).To(Equal(("test")))
	Expect(ok).To(BeTrue())

	v, ok = getPtrValueIfExistsOk(m, "testWrongString")
	Expect(v).To(BeNil())
	Expect(ok).To(BeFalse())

	v, ok = getPtrValueIfExistsOk(m, "testWrongType")
	Expect(v).To(BeNil())
	Expect(ok).To(BeFalse())
}

func TestIdOrLabel(t *testing.T) {
	RegisterTestingT(t)

	i := "test"
	id, label, isID := IdOrLabel(&i)
	Expect(id).To(Equal(0))
	Expect(label).To(Equal("test"))
	Expect(isID).To(BeFalse())

	i = "100"
	id, label, isID = IdOrLabel(&i)
	Expect(id).To(Equal(100))
	Expect(label).To(Equal(""))
	Expect(isID).To(BeTrue())

	ii := 100
	id, label, isID = IdOrLabel(&ii)
	Expect(id).To(Equal(100))
	Expect(label).To(Equal(""))
	Expect(isID).To(BeTrue())
}

func TestMakeWrongCommand(t *testing.T) {
	RegisterTestingT(t)
	cmds := GenerateCommandTestCases(map[string]interface{}{"1": 1, "2": 2, "3": 3, "4": 4})

	Expect(len(cmds)).To(Equal(14))

}

func TestGetRawObjectFromCommand(t *testing.T) {
	RegisterTestingT(t)

	var sw metalcloud.SwitchDevice
	err := json.Unmarshal([]byte(_switchDeviceFixture1), &sw)
	if err != nil {
		t.Error(err)
	}

	f, err := os.CreateTemp(os.TempDir(), "testconf-*.json")
	if err != nil {
		t.Error(err)
	}

	//create an input json file
	f.WriteString(_switchDeviceFixture1)
	f.Close()
	defer syscall.Unlink(f.Name())

	cmd := MakeCommand(map[string]interface{}{
		"read_config_from_file": f.Name(),
		"format":                "json",
	})

	var sw2 metalcloud.SwitchDevice

	err = GetRawObjectFromCommand(&cmd, &sw2)
	Expect(err).To(BeNil())
	Expect(sw.NetworkEquipmentPrimarySANSubnetPool).To(Equal("100.64.0.0"))
}

func TestGetKeyValueMapFromString(t *testing.T) {
	RegisterTestingT(t)

	Expect(GetKeyValueMapFromString("key1=value1,key2=value2")).To(Equal(map[string]string{"key1": "value1", "key2": "value2"}))
	Expect(GetKeyValueMapFromString("key1%3Dvalue1%2Ckey2%3Dvalue2")).To(Equal(map[string]string{"key1": "value1", "key2": "value2"}))
	Expect(GetKeyValueMapFromString("key1%3Dvalue1%2Ckey2%3Dvalue%0A2")).To(Equal(map[string]string{"key1": "value1", "key2": "value\n2"}))
	Expect(GetKeyValueMapFromString("key1=value1, key2=")).To(Equal(map[string]string{"key1": "value1", "key2": ""}))
	_, err := GetKeyValueMapFromString("key1=value1, =value")
	Expect(err).NotTo(BeNil())
	_, err = GetKeyValueMapFromString("key1=value1, =value=")
	Expect(err).NotTo(BeNil())
}

func TestGetKVStringFromMap(t *testing.T) {
	RegisterTestingT(t)

	m := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	s := GetKeyValueStringFromMap(m)
	Expect(s).To(ContainSubstring("key1=value1"))
	Expect(s).To(ContainSubstring("key2=value2"))

	a := []interface{}{} //should support empty array instead of map. Reported by MS-1390

	Expect(GetKeyValueStringFromMap(a)).To(Equal(""))
}

const _switchDeviceFixture1 = "{\"network_equipment_id\":1,\"datacenter_name\":\"uk-reading\",\"network_equipment_driver\":\"hp5900\",\"network_equipment_position\":\"tor\",\"network_equipment_provisioner_type\":\"vpls\",\"network_equipment_identifier_string\":\"UK_RDG_EVR01_00_0001_00A9_01\",\"network_equipment_description\":\"HP Comware Software, Version 7.1.045, Release 2311P06\",\"network_equipment_management_address\":\"10.0.0.0\",\"network_equipment_management_port\":22,\"network_equipment_management_username\":\"sad\",\"network_equipment_quarantine_vlan\":5,\"network_equipment_quarantine_subnet_start\":\"11.16.0.1\",\"network_equipment_quarantine_subnet_end\":\"11.16.0.00\",\"network_equipment_quarantine_subnet_prefix_size\":24,\"network_equipment_quarantine_subnet_gateway\":\"11.16.0.1\",\"network_equipment_primary_wan_ipv4_subnet_pool\":\"11.24.0.2\",\"network_equipment_primary_wan_ipv4_subnet_prefix_size\":22,\"network_equipment_primary_san_subnet_pool\":\"100.64.0.0\",\"network_equipment_primary_san_subnet_prefix_size\":21,\"network_equipment_primary_wan_ipv6_subnet_pool_id\":1,\"network_equipment_primary_wan_ipv6_subnet_cidr\":\"2A02:0CB8:0000:0000:0000:0000:0000:0000/53\",\"network_equipment_cached_updated_timestamp\":\"2020-08-04T20:11:49Z\",\"network_equipment_management_protocol\":\"ssh\",\"chassis_rack_id\":null,\"network_equipment_cache_wrapper_json\":null,\"network_equipment_cache_wrapper_phpserialize\":\"\",\"network_equipment_tor_linked_id\":null,\"network_equipment_uplink_ip_addresses_json\":null,\"network_equipment_management_address_mask\":null,\"network_equipment_management_address_gateway\":null,\"network_equipment_requires_os_install\":false,\"network_equipment_management_mac_address\":\"00:00:00:00:00:00\",\"volume_template_id\":null,\"network_equipment_country\":null,\"network_equipment_city\":null,\"network_equipment_datacenter\":null,\"network_equipment_datacenter_room\":null,\"network_equipment_datacenter_rack\":null,\"network_equipment_rack_position_upper_unit\":null,\"network_equipment_rack_position_lower_unit\":null,\"network_equipment_serial_numbers\":null,\"network_equipment_info_json\":null,\"network_equipment_management_subnet\":null,\"network_equipment_management_subnet_prefix_size\":null,\"network_equipment_management_subnet_start\":null,\"network_equipment_management_subnet_end\":null,\"network_equipment_management_subnet_gateway\":null,\"datacenter_id_parent\":null,\"network_equipment_dhcp_packet_sniffing_is_enabled\":1,\"network_equipment_driver_dump_cached_json\":null,\"network_equipment_tags\":[],\"network_equipment_management_password\":\"ddddd\"}"

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	interfaces "github.com/bigstepinc/metalcloud-cli/interfaces"
)

var osTemplatesCmds = []Command{

	Command{
		Description:  "Lists available Templates",
		Subject:      "templates",
		AltSubject:   "templates",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list templates", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":           c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"usage":            c.FlagSet.String("usage", _nilDefaultStr, "Template's usage"),
				"show_credentials": c.FlagSet.Bool("show_credentials", false, "(Flag) If set returns the templates initial ssh credentials"),
			}
		},
		ExecuteFunc: templatesListCmd,
	},
	Command{
		Description:  "Create template",
		Subject:      "template",
		AltSubject:   "template",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("create template", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"label":                    c.FlagSet.String("label", _nilDefaultStr, "Template's label"),
				"display_name":             c.FlagSet.String("name", _nilDefaultStr, "Template's display name"),
				"size":                     c.FlagSet.Int("size", _nilDefaultInt, "Template's size (bytes)"),
				"local_disk_supported":     c.FlagSet.Bool("local_disk_supported", false, "Template supports local disk install. Default false"),
				"boot_methods_supported":   c.FlagSet.String("boot_methods_supported", _nilDefaultStr, "Template boot methods supported. Defaults to pxe_iscsi."),
				"boot_type":                c.FlagSet.String("boot_type", _nilDefaultStr, "Template boot type. Possible values: 'uefi_only','legacy_only','hyprid' "),
				"description":              c.FlagSet.String("description", _nilDefaultStr, "Template description"),
				"os_type":                  c.FlagSet.String("os_type", _nilDefaultStr, "Template operating system type. For example, Ubuntu or CentOS."),
				"os_version":               c.FlagSet.String("os_version", _nilDefaultStr, "Template operating system version."),
				"os_architecture":          c.FlagSet.String("os_architecture", _nilDefaultStr, "Template operating system architecture.Possible values: none, unknown, x86, x86_64."),
				"os_template_architecture": c.FlagSet.String("template_architecture", _nilDefaultStr, "Template's architecture. Possible values 'uefi','pcx86' "),
				"os_asset_id_bootloader_local_install_id_or_name": c.FlagSet.String("install_bootloader_asset", _nilDefaultStr, "Template's bootloader asset id during install"),
				"os_asset_id_bootloader_os_boot_id_or_name":       c.FlagSet.String("os_boot_bootloader_asset", _nilDefaultStr, "Template's bootloader asset id during regular server boot"),
				"initial_user":                 c.FlagSet.String("initial_user", _nilDefaultStr, "Template's initial username, used to verify install."),
				"initial_password":             c.FlagSet.String("initial_password", _nilDefaultStr, "Template's initial password, used to verify install."),
				"initial_ssh_port":             c.FlagSet.Int("initial_ssh_port", _nilDefaultInt, "Template's initial ssh port, used to verify install."),
				"change_password_after_deploy": c.FlagSet.Bool("change_password_after_deploy", false, "Option to change the initial_user password on the installed OS after deploy."),
				"repo_url":                     c.FlagSet.String("repo_url", _nilDefaultStr, "Template description"),
			}
		},
		ExecuteFunc: templateCreateCmd,
	},
	Command{
		Description:  "Edit template",
		Subject:      "template",
		AltSubject:   "template",
		Predicate:    "update",
		AltPredicate: "edit",
		FlagSet:      flag.NewFlagSet("update template", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"template_id_or_name":      c.FlagSet.String("id", _nilDefaultStr, "Template's id or label"),
				"label":                    c.FlagSet.String("label", _nilDefaultStr, "Template's label"),
				"display_name":             c.FlagSet.String("name", _nilDefaultStr, "Template's display name"),
				"size":                     c.FlagSet.Int("size", _nilDefaultInt, "Template's size (bytes)"),
				"local_disk_supported":     c.FlagSet.Bool("local_disk_supported", false, "Template supports local disk install. Default false"),
				"boot_methods_supported":   c.FlagSet.String("boot_methods_supported", _nilDefaultStr, "Template boot methods supported. Defaults to pxe_iscsi."),
				"boot_type":                c.FlagSet.String("boot_type", _nilDefaultStr, "Template boot type. Possible values: 'uefi_only','legacy_only','hyprid' "),
				"description":              c.FlagSet.String("description", _nilDefaultStr, "Template description"),
				"os_type":                  c.FlagSet.String("os_type", _nilDefaultStr, "Template operating system type. For example, Ubuntu or CentOS."),
				"os_version":               c.FlagSet.String("os_version", _nilDefaultStr, "Template operating system version."),
				"os_architecture":          c.FlagSet.String("os_architecture", _nilDefaultStr, "Template operating system architecture.Possible values: none, unknown, x86, x86_64."),
				"os_template_architecture": c.FlagSet.String("template_architecture", _nilDefaultStr, "Template's architecture. Possible values 'uefi','pcx86' "),
				"os_asset_id_bootloader_local_install_id_or_name": c.FlagSet.String("install_bootloader_asset", _nilDefaultStr, "Template's bootloader asset id during install"),
				"os_asset_id_bootloader_os_boot_id_or_name":       c.FlagSet.String("os_boot_bootloader_asset", _nilDefaultStr, "Template's bootloader asset id during regular server boot"),
				"initial_user":                 c.FlagSet.String("initial_user", _nilDefaultStr, "Template's initial username, used to verify install."),
				"initial_password":             c.FlagSet.String("initial_password", _nilDefaultStr, "Template's initial password, used to verify install."),
				"initial_ssh_port":             c.FlagSet.Int("initial_ssh_port", _nilDefaultInt, "Template's initial ssh port, used to verify install."),
				"change_password_after_deploy": c.FlagSet.Bool("change_password_after_deploy", false, "Option to change the initial_user password on the installed OS after deploy."),
				"repo_url":                     c.FlagSet.String("repo_url", _nilDefaultStr, "Template description"),
			}
		},
		ExecuteFunc: templateEditCmd,
	},
	Command{
		Description:  "Get template",
		Subject:      "template",
		AltSubject:   "template",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("get template", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"template_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Asset's id or name"),
				"format":              c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv'. The default format is human readable."),
				"show_credentials":    c.FlagSet.Bool("show_credentials", false, "(Flag) If set returns the templates initial ssh credentials"),
			}
		},
		ExecuteFunc: templateGetCmd,
	},
	Command{
		Description:  "Delete template",
		Subject:      "template",
		AltSubject:   "template",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete template", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"template_id_or_name": c.FlagSet.String("id", _nilDefaultStr, "Asset's id or name"),
				"autoconfirm":         c.FlagSet.Bool("autoconfirm", false, "If true it does not ask for confirmation anymore"),
			}
		},
		ExecuteFunc: templateDeleteCmd,
	},
}

func templatesListCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	usage := *c.Arguments["usage"].(*string)
	if usage == _nilDefaultStr {
		usage = ""
	}

	list, err := client.OSTemplates()

	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 2,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "NAME",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "DESCRIPTION",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "SIZE_MBYTES",
			FieldType: TypeInt,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "BOOT_METHODS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "OS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "INSTALL_BOOTLOADER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "OS_BOOTLOADER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "USER_ID",
			FieldType: TypeInt,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "CREATED",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "UPDATED",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	user := GetUserEmail()

	data := [][]interface{}{}
	for _, s := range *list {

		installBootloader := ""
		if s.OSAssetIDBootloaderLocalInstall != 0 {
			asset, err := client.OSAssetGet(s.OSAssetIDBootloaderLocalInstall)
			if err != nil {
				return "", err
			}
			installBootloader = asset.OSAssetFileName
		}
		osBootloader := ""
		if s.OSAssetIDBootloaderOSBoot != 0 {
			asset, err := client.OSAssetGet(s.OSAssetIDBootloaderOSBoot)
			if err != nil {
				return "", err
			}
			osBootloader = asset.OSAssetFileName
		}

		osData := ""

		if s.VolumeTemplateOperatingSystem != nil {
			os := *s.VolumeTemplateOperatingSystem
			osData = fmt.Sprintf("%s %s %s",
				os.OperatingSystemType,
				os.OperatingSystemVersion,
				os.OperatingSystemArchitecture)
		}

		data = append(data, []interface{}{
			s.VolumeTemplateID,
			s.VolumeTemplateLabel,
			s.VolumeTemplateDisplayName,
			s.VolumeTemplateDescription,
			s.VolumeTemplateSizeMBytes,
			s.VolumeTemplateBootMethodsSupported,
			osData,
			installBootloader,
			osBootloader,
			s.UserID,
			s.VolumeTemplateCreatedTimestamp,
			s.VolumeTemplateUpdatedTimestamp,
		})

	}

	var sb strings.Builder

	format := c.Arguments["format"]
	if format == nil {
		var f string
		f = ""
		format = &f
	}

	switch *format.(*string) {
	case "json", "JSON":
		ret, err := GetTableAsJSONString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	case "csv", "CSV":
		ret, err := GetTableAsCSVString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)

	default:
		sb.WriteString(fmt.Sprintf("Assets I have access to (as %s)\n", user))

		TableSorter(schema).OrderBy(
			schema[0].FieldName,
			schema[1].FieldName).Sort(data)

		AdjustFieldSizes(data, &schema)

		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d templates\n\n", len(*list)))
	}

	return sb.String(), nil
}

func updateTemplateFromCommand(obj metalcloud.OSTemplate, c *Command, client interfaces.MetalCloudClient, checkRequired bool) (*metalcloud.OSTemplate, error) {

	if v := c.Arguments["label"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateLabel = *v.(*string)
	}

	if v := c.Arguments["display_name"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateDisplayName = *v.(*string)
	}

	if v := c.Arguments["size"]; v != nil && *v.(*int) != _nilDefaultInt {
		obj.VolumeTemplateSizeMBytes = *v.(*int)
	}

	if v := c.Arguments["local_disk_supported"]; v != nil && *v.(*bool) {
		obj.VolumeTemplateLocalDiskSupported = true
	}

	obj.VolumeTemplateIsOSTemplate = true

	if v := c.Arguments["boot_methods_supported"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateBootMethodsSupported = *v.(*string)
	}

	if v := c.Arguments["boot_type"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateBootType = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("boot_type is required")
		}
	}

	if v := c.Arguments["description"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateDescription = *v.(*string)
	}

	//OS Data
	if v := c.Arguments["os_type"]; v != nil && *v.(*string) != _nilDefaultStr {
		vt := metalcloud.OperatingSystem{}
		obj.VolumeTemplateOperatingSystem = &vt
		obj.VolumeTemplateOperatingSystem.OperatingSystemType = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("os_type is required")
		}
	}

	if v := c.Arguments["os_version"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateOperatingSystem.OperatingSystemVersion = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("os_version is required")
		}
	}

	if v := c.Arguments["os_architecture"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateOperatingSystem.OperatingSystemArchitecture = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("os_architecture is required")
		}
	}

	//Boot options
	if v := c.Arguments["os_template_architecture"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.OSTemplateArchitecture = *v.(*string)

	} else {
		if checkRequired {
			return nil, fmt.Errorf("template_architecture is required")
		}
	}

	if v := c.Arguments["os_asset_id_bootloader_local_install_id_or_name"]; v != nil && *v.(*string) != _nilDefaultStr {
		localInstallAsset, err := getOSAssetFromCommand("install_bootloader_asset", "os_asset_id_bootloader_local_install_id_or_name", c, client)
		if err != nil {
			return nil, err
		}
		obj.OSAssetIDBootloaderLocalInstall = localInstallAsset.OSAssetID
	}

	if v := c.Arguments["os_asset_id_bootloader_os_boot_id_or_name"]; v != nil && *v.(*string) != _nilDefaultStr {
		osBootBootloaderAsset, err := getOSAssetFromCommand("os_boot_bootloader_asset", "os_asset_id_bootloader_os_boot_id_or_name", c, client)
		if err != nil {
			return nil, err
		}
		obj.OSAssetIDBootloaderOSBoot = osBootBootloaderAsset.OSAssetID
	}

	//Credentials
	if v := c.Arguments["initial_user"]; v != nil && *v.(*string) != _nilDefaultStr {
		creds := metalcloud.OSTemplateCredentials{}
		obj.OSTemplateCredentials = &creds
		obj.OSTemplateCredentials.InitialUser = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("initial_user is required")
		}
	}

	if v := c.Arguments["initial_password"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.OSTemplateCredentials.InitialPassword = *v.(*string)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("initial_password is required")
		}
	}

	if v := c.Arguments["initial_ssh_port"]; v != nil && *v.(*int) != _nilDefaultInt {
		obj.OSTemplateCredentials.InitialSSHPort = *v.(*int)
	} else {
		if checkRequired {
			return nil, fmt.Errorf("initial_ssh_port is required")
		}
	}

	if v := c.Arguments["change_password_after_deploy"]; v != nil && *v.(*bool) {
		obj.OSTemplateCredentials.ChangePasswordAfterDeploy = true
	}

	if v := c.Arguments["repo_url"]; v != nil && *v.(*string) != _nilDefaultStr {
		obj.VolumeTemplateRepoURL = *v.(*string)
	}

	return &obj, nil
}

func templateCreateCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {
	obj := metalcloud.OSTemplate{}
	updatedObj, err := updateTemplateFromCommand(obj, c, client, true)
	if err != nil {
		return "", err
	}
	_, err = client.OSTemplateCreate(*updatedObj)
	return "", err
}

func templateEditCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	obj, err := getOSTemplateFromCommand("id", c, client, false)
	newobj := metalcloud.OSTemplate{}
	updatedObj, err := updateTemplateFromCommand(newobj, c, client, false)
	if err != nil {
		return "", err
	}
	_, err = client.OSTemplateUpdate(obj.VolumeTemplateID, *updatedObj)
	return "", err
}

func templateDeleteCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	retS, err := getOSTemplateFromCommand("id", c, client, false)
	if err != nil {
		return "", err
	}
	confirm := false

	if c.Arguments["autoconfirm"] != nil && *c.Arguments["autoconfirm"].(*bool) == true {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting template %s (%d).  Are you sure? Type \"yes\" to continue:",
			retS.VolumeTemplateDisplayName,
			retS.VolumeTemplateID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm = requestConfirmation(confirmationMessage)
	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err = client.OSTemplateDelete(retS.VolumeTemplateID)

	return "", err
}

func getOSTemplateFromCommand(paramName string, c *Command, client interfaces.MetalCloudClient, decryptPasswd bool) (*metalcloud.OSTemplate, error) {

	v, err := getParam(c, "template_id_or_name", paramName)
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(v)

	if isID {
		return client.OSTemplateGet(id, decryptPasswd)
	}

	list, err := client.OSTemplates()
	if err != nil {
		return nil, err
	}

	for _, s := range *list {
		if s.VolumeTemplateLabel == label {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Could not locate secret with id/name %v", *v.(*interface{}))
}

func templateGetCmd(c *Command, client interfaces.MetalCloudClient) (string, error) {

	showCredentials := false
	if c.Arguments["show_credentials"] != nil && *c.Arguments["show_credentials"].(*bool) {
		showCredentials = true
	}

	template, err := getOSTemplateFromCommand("id", c, client, showCredentials)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	schema := []SchemaField{
		SchemaField{
			FieldName: "ID",
			FieldType: TypeInt,
			FieldSize: 2,
		},
		SchemaField{
			FieldName: "LABEL",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "NAME",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "DESCRIPTION",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "SIZE_MBYTES",
			FieldType: TypeInt,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "BOOT_METHODS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "OS",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "USER_ID",
			FieldType: TypeInt,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "USER_ID",
			FieldType: TypeInt,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "INSTALL_BOOTLOADER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "OS_BOOTLOADER",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "CREATED",
			FieldType: TypeString,
			FieldSize: 5,
		},
		SchemaField{
			FieldName: "UPDATED",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	data := [][]interface{}{}

	credentials := ""

	if showCredentials {

		schema = append(schema, SchemaField{
			FieldName: "CREDENTIALS",
			FieldType: TypeString,
			FieldSize: 5,
		})

		credentials = fmt.Sprintf("user:%s (port %d) passwd:%s (change_password_after_install:%v)",
			template.OSTemplateCredentials.InitialUser,
			template.OSTemplateCredentials.InitialSSHPort,
			template.OSTemplateCredentials.InitialPassword,
			template.OSTemplateCredentials.ChangePasswordAfterDeploy)

	}
	osDetails := ""

	if template.VolumeTemplateOperatingSystem != nil {
		os := *template.VolumeTemplateOperatingSystem
		osDetails = fmt.Sprintf("%s %s %s",
			os.OperatingSystemType,
			os.OperatingSystemVersion,
			os.OperatingSystemArchitecture)
	}

	installBootloader := ""
	if template.OSAssetIDBootloaderLocalInstall != 0 {
		asset, err := client.OSAssetGet(template.OSAssetIDBootloaderLocalInstall)
		if err != nil {
			return "", err
		}
		installBootloader = asset.OSAssetFileName
	}
	osBootloader := ""
	if template.OSAssetIDBootloaderOSBoot != 0 {
		asset, err := client.OSAssetGet(template.OSAssetIDBootloaderOSBoot)
		if err != nil {
			return "", err
		}
		osBootloader = asset.OSAssetFileName
	}

	data = append(data, []interface{}{
		template.VolumeTemplateID,
		template.VolumeTemplateLabel,
		template.VolumeTemplateDisplayName,
		template.VolumeTemplateDescription,
		template.VolumeTemplateSizeMBytes,
		template.VolumeTemplateBootMethodsSupported,
		osDetails,
		template.UserID,
		installBootloader,
		osBootloader,
		template.VolumeTemplateCreatedTimestamp,
		template.VolumeTemplateUpdatedTimestamp,
		credentials,
	})

	var sb strings.Builder

	format := c.Arguments["format"]
	if format == nil {
		var f string
		f = ""
		format = &f
	}

	switch *format.(*string) {
	case "json", "JSON":
		ret, err := GetTableAsJSONString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	case "csv", "CSV":
		ret, err := GetTableAsCSVString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)

	default:
		sb.WriteString(fmt.Sprintf("Template %s (%d)\n", template.VolumeTemplateLabel, template.VolumeTemplateID))

		TableSorter(schema).OrderBy(
			schema[0].FieldName,
			schema[1].FieldName).Sort(data)

		AdjustFieldSizes(data, &schema)

		sb.WriteString(GetTableAsString(data, schema))

	}

	return sb.String(), nil
}
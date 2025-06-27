package main

import (
	"encoding/base64"
	"nikeron/cmbxip/config"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

func settingsWindowAction(action int) { // 0 - load, 1 - save, 2 - update state
	cmAPIURLObj, _ := builder.GetObject("settings-cm-api-url")
	cmUsernameObj, _ := builder.GetObject("settings-cm-username")
	cmPasswordObj, _ := builder.GetObject("settings-cm-password")
	cmSaveAuthObj, _ := builder.GetObject("settings-cm-save-auth")
	bxWebHookURLObj, _ := builder.GetObject("settings-bx-webhook-url")
	bxSaveAuthObj, _ := builder.GetObject("settings-bx-save-auth")
	bxTaskTitleObj, _ := builder.GetObject("settings-bx-task-title")
	bxTaskDescriptionObj, _ := builder.GetObject("settings-bx-task-description")
	bxTaskDescriptionBuffer, _ := bxTaskDescriptionObj.(*gtk.TextView).GetBuffer()
	bxRootAutofindEnabledObj, _ := builder.GetObject("settings-bx-root-autofind-enabled")
	bxRootAutofindIsAuditorObj, _ := builder.GetObject("settings-bx-root-autofind-is-auditor")
	bxRootAutofindIDObj, _ := builder.GetObject("settings-bx-root-autofind-id")
	switch action {
	case 0:
		cmAPIURLObj.(*gtk.Entry).SetText(config.CM().APIEntry)
		cmUsername, cmPassword := "", ""
		if decodedAuthString, err := base64.StdEncoding.DecodeString(config.CM().Auth); err == nil {
			splittedAuthString := strings.SplitN(string(decodedAuthString), ":", 2)
			if len(splittedAuthString) == 2 {
				cmUsername, cmPassword = splittedAuthString[0], splittedAuthString[1]
			}
		}
		cmUsernameObj.(*gtk.Entry).SetText(cmUsername)
		cmPasswordObj.(*gtk.Entry).SetText(cmPassword)
		cmSaveAuthObj.(*gtk.CheckButton).SetActive(config.CM().SaveAuth)
		bxWebHookURLObj.(*gtk.Entry).SetText(config.BX().InWebHook)
		bxSaveAuthObj.(*gtk.CheckButton).SetActive(config.CM().SaveAuth)
		bxTaskTitleObj.(*gtk.Entry).SetText(config.BX().TaskTitleFormat)
		bxTaskDescriptionBuffer.SetText(config.BX().TaskDescriptionFormat)
		bxRootAutofindEnabledObj.(*gtk.CheckButton).SetActive(config.BX().RootAutofindEnabled)
		bxRootAutofindIsAuditorObj.(*gtk.CheckButton).SetActive(config.BX().RootAutofindIsAuditor)
		bxRootAutofindIDObj.(*gtk.Entry).SetText(config.BX().RootAutofindID)
		settingsWindowAction(2)
	case 1:
		config.CM().APIEntry, _ = cmAPIURLObj.(*gtk.Entry).GetText()
		config.CM().APIEntry = strings.TrimRight(config.CM().APIEntry, "/")
		cmUsername, _ := cmUsernameObj.(*gtk.Entry).GetText()
		cmPassword, _ := cmPasswordObj.(*gtk.Entry).GetText()
		config.CM().Auth = base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{
			cmUsername, cmPassword}, ":")))
		config.CM().SaveAuth = cmSaveAuthObj.(*gtk.CheckButton).GetActive()
		config.BX().InWebHook, _ = bxWebHookURLObj.(*gtk.Entry).GetText()
		config.BX().InWebHook = strings.TrimRight(config.BX().InWebHook, "/")
		config.BX().SaveInWebHook = bxSaveAuthObj.(*gtk.CheckButton).GetActive()
		config.BX().TaskTitleFormat, _ = bxTaskTitleObj.(*gtk.Entry).GetText()
		config.BX().TaskDescriptionFormat, _ = bxTaskDescriptionBuffer.
			GetText(bxTaskDescriptionBuffer.GetStartIter(), bxTaskDescriptionBuffer.GetEndIter(), false)
		config.BX().RootAutofindEnabled = bxRootAutofindEnabledObj.(*gtk.CheckButton).GetActive()
		config.BX().RootAutofindIsAuditor = bxRootAutofindIsAuditorObj.(*gtk.CheckButton).GetActive()
		config.BX().RootAutofindID, _ = bxRootAutofindIDObj.(*gtk.Entry).GetText()
	case 2:
		bxRootAutofindIsAuditorObj.(*gtk.CheckButton).SetSensitive(bxRootAutofindEnabledObj.(*gtk.CheckButton).GetActive())
		bxRootAutofindIDObj.(*gtk.Entry).SetSensitive(bxRootAutofindEnabledObj.(*gtk.CheckButton).GetActive())
	}
}

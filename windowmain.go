package main

import (
	"fmt"
	"io"
	"nikeron/cmbxip/bx"
	"nikeron/cmbxip/cm"
	"nikeron/cmbxip/config"
	"strconv"
	"strings"
	"time"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var currentDocument *cm.Document

func mainWindowCMDocumentURLEntryChanged() {
	dt := time.Now()
	mainWindowTaskView.Hide()
	mainWindowNotifyLabel.SetText("Парсинг документа...")
	fmt.Println(dt.Format("15:04:05"), "Получение документа")
	mainWindowNotifyLabel.Show()

	mainWindow.QueueDraw() // update window to show information
	for glib.MainContextDefault().Pending() {
		glib.MainContextDefault().Iteration(true)
	}

	cmClient := cm.NewClient(config.CM().APIEntry, config.CM().Auth)
	documentURL, _ := mainWindowDocumentURLEntry.GetText()
	cmDocumentID := cm.IDFromAddress(documentURL)
	if cmDocumentID == "" {
		dt := time.Now()
		mainWindowNotifyLabel.SetText("Невозможно получить ID документа из ссылки!")
		fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Невозможно получить ID документа из ссылки")
		return
	}

	cmDocument, err := cmClient.FromID(cmDocumentID)
	if err != nil {
		dt := time.Now()
		mainWindowNotifyLabel.SetText(fmt.Sprintf("Ошибка получения документа!\n\n%s", err.Error()))
		fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Ошибка получения документа!\n\n%s")
		return
	}
	currentDocument = cmDocument

	{ // format task title and description
		formatStruct := BitrixTaskTitleAndDescriptionFormat{
			URL:          documentURL,
			Title:        cmDocument.Title,
			Type:         cmDocument.Type,
			RegDate:      cmDocument.Registration.Date,
			RegNumPrefix: cmDocument.Registration.Number.Prefix,
			RegNumber:    cmDocument.Registration.Number.Number,
			RegNumSuffix: cmDocument.Registration.Number.Suffix,
		}
		formatStruct.Correspondent.Organization.FullName = cmDocument.Correspondent.Organization.Organization.FullName

		taskTitle, err := parseTemplateFromStruct(config.BX().TaskTitleFormat, formatStruct)
		if err != nil {
			dt := time.Now()
			mainWindowNotifyLabel.SetText(fmt.Sprintf("Ошибка парсинга заголовка задачи!\n\n%s", err.Error()))
			fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Ошибка парсинга заголовка задачи!\n\n%s")
			return
		}
		mainWindowTaskTitle.SetText(taskTitle)
		dt := time.Now()
		fmt.Print(dt.Format("15:04:05"), " Получен документ: ", cmDocument.Registration.Number.Prefix,cmDocument.Registration.Number.Number,cmDocument.Registration.Number.Suffix,"\n")

		taskDescription, err := parseTemplateFromStruct(config.BX().TaskDescriptionFormat, formatStruct)
		if err != nil {
			dt := time.Now()
			mainWindowNotifyLabel.SetText(fmt.Sprintf("Ошибка парсинга описания задачи!\n\n%s", err.Error()))
			fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Ошибка парсинга описания задачи!\n\n%s")
			return
		}
		taskDescBuf, _ := mainWindowTaskDescription.GetBuffer()
		taskDescBuf.SetText(taskDescription)
	}

	executionHierarchy, err := cmClient.ExecutionHierarchy(cmDocumentID)
	if err != nil {
		dt := time.Now()
		mainWindowNotifyLabel.SetText(fmt.Sprintf("Ошибка получения иерархии исполнителей!\n\n%s", err.Error()))
		fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Ошибка получения иерархии исполнителей!\n\n%s")
		return
	}

	mainWindowExecutorsTreeStore.Clear()
	mainWindowExecutorsTreeRootEntry = &ExecutorTreeEntry{}
	for _, v := range executionHierarchy.Entry {
		recursiveTreeBuild(mainWindowExecutorsTreeRootEntry, v)
	}
	if config.BX().RootAutofindEnabled {
		if newRoot := recursiveTreeFindRoot(mainWindowExecutorsTreeRootEntry); newRoot != nil {
			if config.BX().RootAutofindIsAuditor {
				newRoot.BitrixRole = bxRoleAuditor
			}
			mainWindowExecutorsTreeRootEntry = &ExecutorTreeEntry{SubEntries: []*ExecutorTreeEntry{newRoot}}
		}
	}
	for _, v := range mainWindowExecutorsTreeRootEntry.SubEntries {
		recursiveTreeFill(nil, v)
	}
	mainWindowExecutorsTree.ExpandAll()

	go recursiveTreeFindBitrixID(mainWindowExecutorsTreeRootEntry, bx.NewClient(config.BX().InWebHook))

	mainWindowNotifyLabel.Hide()
	mainWindowTaskView.Show()
}

func recursiveTreeBuild(parentTreeEntry *ExecutorTreeEntry, entry cm.ExecutionEntry) {
	if parentTreeEntry.CompanyMediaID == "" {
		e := ExecutorTreeEntry{
			FullName:       entry.Value.Author.FullName,
			CompanyMediaID: entry.Value.Author.ID,
			BitrixRole:     bxRoleExecutor,
		}
		parentTreeEntry.SubEntries = append(parentTreeEntry.SubEntries, &e)
		recursiveTreeBuild(&e, entry)
		return
	}
	for _, executor := range entry.Value.Executor {
		e := ExecutorTreeEntry{
			FullName:       executor.Executor.FullName,
			CompanyMediaID: executor.Executor.ID,
			BitrixRole:     bxRoleExecutor,
		}
		parentTreeEntry.SubEntries = append(parentTreeEntry.SubEntries, &e)
	}
executionloop:
	for _, execution := range entry.Value.Execution {
		for _, author := range parentTreeEntry.SubEntries {
			if author.CompanyMediaID == execution.Value.Author.ID {
				recursiveTreeBuild(author, execution)
				continue executionloop
			}
		}
		author := &ExecutorTreeEntry{
			FullName:       execution.Value.Author.FullName,
			CompanyMediaID: execution.Value.Author.ID,
			BitrixRole:     bxRoleExecutor,
		}
		parentTreeEntry.SubEntries = append(parentTreeEntry.SubEntries, author)
		recursiveTreeBuild(author, execution)
	}
}

func recursiveTreeFindRoot(parentTreeEntry *ExecutorTreeEntry) *ExecutorTreeEntry {
	for _, v := range parentTreeEntry.SubEntries {
		if config.BX().RootAutofindID == v.CompanyMediaID || config.BX().RootAutofindID == v.FullName {
			return v
		}
		if retVal := recursiveTreeFindRoot(v); retVal != nil {
			return retVal
		}
	}
	return nil
}

func recursiveTreeFindBitrixID(entry *ExecutorTreeEntry, bxClient *bx.Client) {
	for _, v := range entry.SubEntries {
		go recursiveTreeFindBitrixID(v, bxClient)
	}
	if entry.FullName != "" {
		users, err := bxClient.SearchUser(entry.FullName, true)
		if err != nil {
			return
		}
		for _, u := range users {
			if u.Active {
				entry.BitrixID = u.ID
				entry.NameBackgroundColor = "lightgreen"
				entry.updateValuesInTree()
				break
			}
		}
	}
}

func recursiveTreeFill(parentTreeIter *gtk.TreeIter, entry *ExecutorTreeEntry) {
	entry.TreeIter = mainWindowExecutorsTreeStore.Append(parentTreeIter)
	entry.updateValuesInTree()
	for _, v := range entry.SubEntries {
		recursiveTreeFill(entry.TreeIter, v)
	}
}

func recursiveExecutorTreeEntryToArray(entry *ExecutorTreeEntry) []*ExecutorTreeEntry {
	retValue := []*ExecutorTreeEntry{}
	if entry.FullName != "" {
		retValue = append(retValue, entry)
	}
	for _, v := range entry.SubEntries {
		retValue = append(retValue, recursiveExecutorTreeEntryToArray(v)...)
	}
	return retValue
}

func (entry *ExecutorTreeEntry) updateValuesInTree() {
	glib.IdleAdd(func() {
		mainWindowExecutorsTreeStore.SetValue(entry.TreeIter, mainWindowExecutorsTreeColumnName, entry.FullName)
		if entry.NameBackgroundColor != "" {
			mainWindowExecutorsTreeStore.SetValue(entry.TreeIter, mainWindowExecutorsTreeColumnColor, entry.NameBackgroundColor)
		}
		mainWindowExecutorsTreeStore.SetValue(entry.TreeIter, mainWindowExecutorsTreeColumnRole, bxRoleName[entry.BitrixRole])
	})
}

func recursiveFindExecutorTreeEntryByTreeIter(entry *ExecutorTreeEntry, iter *gtk.TreeIter) *ExecutorTreeEntry {
	var retValue *ExecutorTreeEntry
	for _, v := range entry.SubEntries {
		if v.TreeIter.GtkTreeIter == iter.GtkTreeIter {
			return v
		}
		if retValue = recursiveFindExecutorTreeEntryByTreeIter(v, iter); retValue != nil {
			return retValue
		}
	}
	return nil
}

func mainWindowExecutorsTreeRoleCellChanged(r *gtk.CellRendererCombo, path string, newTreeIter *gtk.TreeIter) {
	iter, _ := mainWindowExecutorsTreeStore.GetIterFromString(path)
	executorTreeEntry := recursiveFindExecutorTreeEntryByTreeIter(mainWindowExecutorsTreeRootEntry, iter)
	if executorTreeEntry != nil {
		selectedValue, _ := mainWindowExecutorsTreeRoleComboBoxStore.GetValue(newTreeIter, 0)
		selectedGoValue, _ := selectedValue.GoValue()
		executorTreeEntry.BitrixRole = bitrixRole(selectedGoValue.(int))
		executorTreeEntry.updateValuesInTree()
	}
}

func mainWindowImportButtonClicked() {
	dt := time.Now()
	mainWindowTaskView.Hide()
	mainWindowNotifyLabel.SetText("Импорт задачи...")
	fmt.Println(dt.Format("15:04:05"), "Импорт задачи")
	mainWindowNotifyLabel.Show()

	mainWindow.QueueDraw() // update window to show information
	for glib.MainContextDefault().Pending() {
		glib.MainContextDefault().Iteration(true)
	}

	bxClient := bx.NewClient(config.BX().InWebHook)
	cmClient := cm.NewClient(config.CM().APIEntry, config.CM().Auth)

	bxCurrentUser, err := bxClient.CurrentUser()
	if err != nil || bxCurrentUser == nil {
		dt := time.Now()
		mainWindowNotifyLabel.SetText(fmt.Sprintf("Ошибка получения информации о текущем пользователе Bitrix!\n\n%s", err.Error()))
		fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Ошибка получения информации о текущем пользователе Bitrix!\n\n%s")
		return
	}
	var bxCurrentUserStorageID, bxUploadFolderID string
	bxStorageList, err := bxClient.DiskStorageListOfUser(bxCurrentUser.ID)
	if err != nil {
		dt := time.Now()
		mainWindowNotifyLabel.SetText(fmt.Sprintf("Ошибка получения списка хранилищ пользователя Bitrix!\n\n%s", err.Error()))
		fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Ошибка получения списка хранилищ пользователя Bitrix!\n\n%s")
		return
	}
	for _, storage := range bxStorageList.Result {
		if strings.ToLower(storage.EntityType) == "user" && storage.EntityID == bxCurrentUser.ID {
			bxCurrentUserStorageID = storage.ID
			break
		}
	}
	bxStorageChildren, err := bxClient.DiskStorageChildren(bxCurrentUserStorageID)
	if err != nil {
		dt := time.Now()
		mainWindowNotifyLabel.SetText(fmt.Sprintf("Ошибка получения потомка хранилища Bitrix!\n\n%s", err.Error()))
		fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Ошибка получения потомка хранилища Bitrix!\n\n%s")
		return
	}
	for _, children := range bxStorageChildren.Result {
		if strings.ToUpper(children.Code) == "FOR_UPLOADED_FILES" {
			bxUploadFolderID = children.ID
			break
		}
	}
	if bxUploadFolderID == "" {
		dt := time.Now()
		mainWindowNotifyLabel.SetText(fmt.Sprintf("Ошибка получения директории выгрузки Bitrix!\n\n%s", err.Error()))
		fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Ошибка получения директории выгрузки Bitrix!\n\n%s")
		return
	}

	bitrixExecutorsIDs, bitrixAuditorsIDs := []string{}, []string{}
	for _, v := range recursiveExecutorTreeEntryToArray(mainWindowExecutorsTreeRootEntry) {
		if v.BitrixID != "" {
			switch v.BitrixRole {
			case bxRoleExecutor:
				//log.Printf("Executor: %+v", v)
				bitrixExecutorsIDs = append(bitrixExecutorsIDs, v.BitrixID)
			case bxRoleAuditor:
				//log.Printf("Auditor: %+v", v)
				bitrixAuditorsIDs = append(bitrixAuditorsIDs, v.BitrixID)
			}
		}
	}

	if len(bitrixExecutorsIDs) < 1 {
		dt := time.Now()
		mainWindowNotifyLabel.SetText("Должен быть хотя бы один исполнитель!")
		fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Должен быть хотя бы один исполнитель")
		return
	}

	//log.Printf("%+v %+v", bitrixExecutorsIDs, bitrixAuditorsIDs)

	taskTitle, _ := mainWindowTaskTitle.GetText()
	taskDescriptionBuffer, _ := mainWindowTaskDescription.GetBuffer()
	taskDescription, _ := taskDescriptionBuffer.GetText(
		taskDescriptionBuffer.GetStartIter(), taskDescriptionBuffer.GetEndIter(), true)
	bxCreatedTask, err := bxClient.AddTask(&bx.AddTaskRequest{
		Fields: &bx.TaskFields{
			Title:         taskTitle,
			Description:   taskDescription,
			ResponsibleID: bitrixExecutorsIDs[0],
			Accomplices:   bitrixExecutorsIDs[1:],
			Auditors:      bitrixAuditorsIDs,
		},
	})
	if err != nil {
		dt := time.Now()
		mainWindowNotifyLabel.SetText(fmt.Sprintf("Ошибка создания задачи Bitrix!\n\n%s", err.Error()))
		fmt.Println(dt.Format("15:04:05"), "ВНИМАНИЕ: Ошибка создания задачи Bitrix!\n\n%s")
		return
	}

	finalNotificationText := []string{}

	bxAttachIDs := []string{}
	for _, content := range currentDocument.Content {
		config.DebugLogger().Println()
		config.DebugLogger().Printf("Парсинг контента %s", content.Title)
		config.DebugLogger().Printf("HrefAsURI: \"%s\", Href: \"%s\"", content.HrefAsURI, content.Href)
		uri := content.HrefAsURI
		if uri == "" {
			uri = content.Href
		}
		if uri == "" {
			finalNotificationText = append(finalNotificationText, fmt.Sprintf("Невозможно найти ссылку в документе из Company Media для\n%s", content.Title))
			continue
		}
		resp, err := cmClient.GetURI(uri)
		if err != nil {
			finalNotificationText = append(finalNotificationText, fmt.Sprintf("Невозможно загрузить контент из Company Media!\n%s", err.Error()))
			continue
		}
		config.DebugLogger().Printf("Ответ: %+v", resp)
		defer resp.Body.Close()
		contentData, err := io.ReadAll(resp.Body)
		if err != nil {
			finalNotificationText = append(finalNotificationText, fmt.Sprintf("Невозможно загрузить контент из Company Media!\n%s", err.Error()))
			continue
		}
		config.DebugLogger().Printf("Скачано: %d байт", len(contentData))
		dt := time.Now()
		fmt.Println(dt.Format("15:04:05"), "Скачано байт:", len(contentData))
		bxUploadedFile, err := bxClient.DiskFolderUploadFile(bxUploadFolderID, content.Title, contentData, true)
		if err != nil {
			finalNotificationText = append(finalNotificationText, fmt.Sprintf("Невозможно выгрузить контент в Bitrix!\n%s", err.Error()))
			continue
		}
		bxAttachIDs = append(bxAttachIDs, strconv.Itoa(bxUploadedFile.Result.ID))
	}
	for _, image := range currentDocument.Image {
		resp, err := cmClient.GetURI(image.HrefAsURI)
		if err != nil {
			finalNotificationText = append(finalNotificationText, fmt.Sprintf("Невозможно загрузить образ из Company Media!\n%s", err.Error()))
			continue
		}
		defer resp.Body.Close()
		contentData, err := io.ReadAll(resp.Body)
		if err != nil {
			finalNotificationText = append(finalNotificationText, fmt.Sprintf("Невозможно загрузить образ из Company Media!\n%s", err.Error()))
			continue
		}
		bxUploadedFile, err := bxClient.DiskFolderUploadFile(bxUploadFolderID, image.Title, contentData, true)
		if err != nil {
			finalNotificationText = append(finalNotificationText, fmt.Sprintf("Невозможно выгрузить образ в Bitrix!\n%s", err.Error()))
			continue
		}
		bxAttachIDs = append(bxAttachIDs, strconv.Itoa(bxUploadedFile.Result.ID))
	}

	for _, fileID := range bxAttachIDs {
		_, err := bxClient.TaskAttachFile(bxCreatedTask.Result.Task.ID, fileID)
		if err != nil {
			finalNotificationText = append(finalNotificationText, fmt.Sprintf("Невозможно прикрепить файл к задаче Bitrix!\n%s", err.Error()))
			continue
		}
	}

	finalNotificationText = append(finalNotificationText, "Задача создана!")
	fmt.Println(dt.Format("15:04:05"), "Задача создана")
	mainWindowNotifyLabel.SetText(strings.Join(finalNotificationText, "\n\n"))
}

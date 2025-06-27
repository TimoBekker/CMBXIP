package bx

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

type User struct {
	ID     string `json:"ID"`
	Active bool   `json:"active"`
	Type   string `json:"USER_TYPE"`

	Name       string `json:"NAME"`
	LastName   string `json:"LAST_NAME"`
	SecondName string `json:"SECOND_NAME"`
	EMail      string `json:"EMAIL"`
}

type UserSearchResponse struct {
	Result []*User `json:"result"`
	Total  int     `json:"total"`
}

type UserCurrentResponse struct {
	Result *User `json:"result"`
}

type ProfileResponse struct {
	Result *struct {
		ID       string `json:"ID"`
		Name     string `json:"NAME"`
		LastName string `json:"LAST_NAME"`
	} `json:"result"`
}

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskFields struct {
	Title         string   `json:"TITLE"`
	Description   string   `json:"DESCRIPTION"`
	ResponsibleID string   `json:"RESPONSIBLE_ID"`
	Accomplices   []string `json:"ACCOMPLICES"`
	Auditors      []string `json:"AUDITORS"`
}

type AddTaskRequest struct {
	Fields *TaskFields `json:"fields"`
}

type AddTaskResponse struct {
	Result *struct {
		Task *Task `json:"task"`
	} `json:"result"`
}

type TaskAttachFileResponse struct {
	Result *struct {
		AttachmentId int `json:"attachmentId"`
	} `json:"result"`
}

type DiskStorageListResponse struct {
	Result []struct {
		ID         string `json:"ID"`
		Name       string `json:"NAME"`
		Code       string `json:"CODE"`
		EntityType string `json:"ENTITY_TYPE"`
		EntityID   string `json:"ENTITY_ID"`
	} `json:"result"`
}

type DiskStorageChildrenResponse struct {
	Result []struct {
		ID        string `json:"ID"`
		Name      string `json:"NAME"`
		Code      string `json:"CODE"`
		StorageID string `json:"STORAGE_ID"`
	} `json:"result"`
}

type DiskFolderUploadRequest struct {
	ID   string `json:"id"`
	Data struct {
		Name string `json:"NAME"`
	} `json:"data"`
	FileContent        []string `json:"fileContent"`
	GenerateUniqueName bool     `json:"generateUniqueName"`
}

type DiskFolderUploadResponse struct {
	Result struct {
		ID int `json:"ID"`
	} `json:"result"`
}

package config

type Configuration struct {
	CompanyMedia *CompanyMedia `json:"cm"`
	Bitrix       *Bitrix       `json:"bx"`
}

type CompanyMedia struct {
	APIEntry string `json:"api_entry"`
	SaveAuth bool   `json:"save_auth"`
	Auth     string `json:"auth"`
}

type Bitrix struct {
	SaveInWebHook         bool   `json:"save_inhook"`
	InWebHook             string `json:"inhook"`
	TaskTitleFormat       string `json:"task_title"`
	TaskDescriptionFormat string `json:"task_description"`
	RootAutofindEnabled   bool   `json:"root_autofind_enabled"`
	RootAutofindIsAuditor bool   `json:"root_autofind_is_auditor"`
	RootAutofindID        string `json:"root_autofind_id"`
}

func CM() *CompanyMedia {
	return config.CompanyMedia
}

func BX() *Bitrix {
	return config.Bitrix
}

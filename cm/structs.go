package cm

type ErrorResponse struct {
	Message string `json:"errorMessage"`
}

type ComponentVersionsResponse struct {
	Entry []struct {
		Component string `json:"component"`
		Version   string `json:"version"`
	} `json:"entry"`
}

type Document struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	Registration struct {
		Number struct {
			Prefix string `json:"prefix"`
			Number int    `json:"number"`
			Suffix string `json:"suffix"`
		} `json:"number"`
		Date string `json:"date"`
	} `json:"registration"`
	Correspondent struct {
		Organization struct {
			Organization struct {
				FullName string `json:"fullName"`
			} `json:"organization"`
		} `json:"organization"`
	} `json:"correspondent"`
	Content []struct {
		Href      string `json:"href"`
		HrefAsURI string `json:"hrefAsUri"`
		Type      string `json:"type"`
		Title     string `json:"title"`
		Extension string `json:"extension"`
	} `json:"content"`
	Image []struct {
		Href      string `json:"href"`
		HrefAsURI string `json:"hrefAsUri"`
		Title     string `json:"title"`
		Extension string `json:"extension"`
	} `json:"image"`
}

type ExecutionResponse struct {
	Entry []ExecutionEntry `json:"entry"`
}

type ExecutionEntry struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value struct {
		ID     string `json:"id"`
		Author struct {
			ID       string `json:"id"`
			FullName string `json:"fullName"`
		} `json:"author"`
		Executor []struct {
			Executor struct {
				ID       string `json:"id"`
				FullName string `json:"fullName"`
			} `json:"executor"`
		} `json:"executor"`
		Execution []ExecutionEntry `json:"execution"`
	} `json:"value"`
}

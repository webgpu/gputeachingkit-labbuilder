package build

type questionAnswer struct {
	Question string
	Answer   string
}

type doc struct {
	Module          int
	FrontMatter     string
	FileName        string
	Name            string
	Description     string
	QuestionAnswers []questionAnswer
	CodeTemplate    string
	CodeSolution    string
}

type resource struct {
	fileName string
	baseName string
	content  string
}

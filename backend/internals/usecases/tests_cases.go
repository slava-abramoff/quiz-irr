package usecases

type testsCases struct{}

func NewTestsCases() *testsCases {
	return &testsCases{}
}

// Работа админа с тестом
func (t *testsCases) NewTest()

func (t *testsCases) GetTestPreview()

func (t *testsCases) FindManyTests()

func (t *testsCases) GetTestFullData()

func (t *testsCases) UpdateTest()

func (t *testsCases) ActivateTest()

func (t *testsCases) DeleteTest()

// Работа админа с вопросами и ответами
func (t *testsCases) AddQuestion()

func (t *testsCases) EditQuestion()

func (t *testsCases) DeleteQuestion()

func (t *testsCases) AddOption()

func (t *testsCases) EditOption()

func (t *testsCases) DeleteOption()

// Работа пользователя с тестом
func (t *testsCases) GetTestInfo()

func (t *testsCases) StartTest()

func (t *testsCases) SaveAnswers()

// Работа с сырыми данными
func (t *testsCases) FindRawResults()

func (t *testsCases) AnalyzeResults()

func (t *testsCases) MakeReportByAnalyze()

func (t *testsCases) DeleteRawResults()
func (t *testsCases) DeleteAllRawByTest()

// Работа админа c результами
func (t *testsCases) GetListByTest()

func (t *testsCases) SendResultByEmail()

func (t *testsCases) SendAllResultsTestByEmail()

func (t *testsCases) MakeExcelList()

func (t *testsCases) DeleteResult()

func (t *testsCases) DeleteResultsByTest()

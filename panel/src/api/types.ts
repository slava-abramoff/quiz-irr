export interface ErrorResponse {
  status: 'error';
  message: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  full_name: string;
  email: string;
  access_token: string;
  refresh_token: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface RefreshTokenResponse {
  access_token: string;
}

export interface Pagination {
  current_page: number;
  total_pages: number;
  total_items: number;
  items_perpage: number;
  has_next_page: boolean;
  has_previous_page: boolean;
}

export interface OptionResponse {
  id: number;
  text: string;
  is_correct: boolean;
}

export interface QuestionResponse {
  id: number;
  text: string;
  type: string;
  points: number;
  options: OptionResponse[];
}

export interface TestAdminResponse {
  id: string;
  title: string;
  desc: string;
  is_active: boolean;
  start_at: string;
  end_at: string;
  author: string;
  duration: number;
  questions: QuestionResponse[];
}

export type TestPreviewResponse = Omit<TestAdminResponse, 'questions'>;

export interface TestCustomerResponse {
  title: string;
  desc: string;
  duration: number;
  start_at: string;
  end_at: string;
}

export interface GetManyTestsResponse {
  tests: TestAdminResponse[];
  pagination: Pagination;
}

export interface CreateTestRequest {
  title: string;
  desc: string;
}

export interface UpdateTestRequest {
  title?: string;
  desc?: string;
  is_active?: boolean;
  duration?: number;
  start_at?: string;
  end_at?: string;
}

export interface CreateQuestionRequest {
  text: string;
  type: string;
  points: number;
}

export interface UpdateQuestionRequest {
  text?: string;
  type?: string;
  points?: number;
}

export interface CreateOptionRequest {
  text: string;
  is_correct: boolean;
}

export interface UpdateOptionRequest {
  text?: string;
  is_correct?: boolean;
}

export interface StartExamRequest {
  test_id: string;
  full_name: string;
  email: string;
  org: string;
}

export interface ExamOption {
  id: number;
  text: string;
}

export interface ExamQuestion {
  id: number;
  text: string;
  type: string;
  options: ExamOption[];
}

export interface StartExamResponse {
  raw_id: string;
  questions: ExamQuestion[];
}

export interface SendUserAnswersRequest {
  answers: {
    answer_id: number;
    option_ids: number[];
    text_option: string;
  }[];
}

export interface SendUserAnswersResponse {
  message: string;
}

export interface RawInfo {
  id: string;
  full_name: string;
  email: string;
  org: string;
  status: string;
  start_at: string;
  end_at: string;
}

export interface RawsInfoResponse {
  data: RawInfo[];
  pagination: Pagination;
}

export interface AnalyzeRawResponse {
  message: string;
}

export interface SimpleMessageResponse {
  message: string;
}

export interface ResultResponse {
  id: number;
  full_name: string;
  email: string;
  org: string;
  duration: number;
  is_on_time: boolean;
  total_score: number;
}

export interface ResultsResponse {
  data: ResultResponse[];
  pagination: Pagination;
}


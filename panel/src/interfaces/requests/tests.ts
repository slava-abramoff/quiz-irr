export interface CreateTest {
  title: string;
  desc: string;
}

export interface UpdateTest {
  title?: string;
  desc?: string;
  is_active?: string;
  duration?: number;
  start_at?: string;
  end_at?: string;
}

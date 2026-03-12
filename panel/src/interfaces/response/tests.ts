import type { Question } from "./question";

export interface Test {
  id: string;
  title: string;
  desc: string;
  is_active: string;
  start_at: string;
  end_at: string;
  author: string;
  duration: number;
  questions?: Question[];
}

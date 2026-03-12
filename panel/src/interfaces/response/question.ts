import type { Option } from "./option";

export interface Question {
  id: number;
  text: string;
  type: string;
  points: number;
  options: Option[];
}

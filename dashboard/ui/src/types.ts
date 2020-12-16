export interface Media {
  id: string;
  title: string;
  status: string;
  created_at: string;
  updated_at: string;
  edges: {};
}

export interface ArrayResponse {
  code: number;
  data: Media[];
}

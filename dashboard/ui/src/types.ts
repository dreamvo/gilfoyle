export interface MediaFile {
  id: string;
  media: string;
  rendition_name: string;
  resolution_width: number;
  resolution_height: number;
  target_bandwidth: number;
  framerate: number;
  duration_seconds: number;
  media_type: string;
  created_at: string;
  updated_at: string;
}

export interface Media {
  id: string;
  title: string;
  status: string;
  created_at: string;
  updated_at: string;
  edges: {
    media_files?: MediaFile[];
  };
}

export interface ArrayResponse<T = unknown> {
  code: number;
  data: T[];
}

export interface DataResponse<T = unknown> {
  code: number;
  data: T;
}

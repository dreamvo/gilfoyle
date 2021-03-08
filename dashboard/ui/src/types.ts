export interface MediaFile {
  id: string;
  media: string;
  status: string;
  message: string;
  rendition_name: string;
  entry_file: string;
  mimetype: string;
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
  message: string;
  playable: boolean;
  created_at: string;
  updated_at: string;
  edges: {
    media_files?: MediaFile[];
    probe?: Probe;
  };
}

export interface Probe {
  id: string;
  filename: string;
  checksum_sha256: string;
  aspect_ratio: string;
  width: number;
  height: number;
  duration_seconds: number;
  framerate: number;
  format: string;
  nb_streams: number;
  created_at: string;
  updated_at: string;
  edges: {};
}

export interface ArrayResponse<T = unknown> {
  code: number;
  metadata: { [key: string]: unknown };
  data: T[];
}

export interface DataResponse<T = unknown> {
  code: number;
  data: T;
}

export interface Source {
  src: string;
  type: string;
}

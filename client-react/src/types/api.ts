export interface ApiResponse<T> {
  code: number;
  msg?: string;
  data: T;
}

export interface UserProfile {
  _id: string;
  user: string;
  photo: string;
}

export interface ArticleComment {
  _id?: string;
  user: UserProfile;
  content: string;
  date: string;
}

export interface ArticleSummary {
  _id: string;
  title: string;
  content: string;
  surface: string;
  date: string;
  updateDate?: string;
  tag: string;
  type?: string;
  pv: number;
  comment: ArticleComment[];
}

export interface ArticleDetail extends ArticleSummary {
  html?: string;
}

export interface ArticleInfo {
  tags: string[];
}

export interface MessageChildItem {
  _id?: string;
  user: UserProfile;
  reUser: string;
  content: string;
  date: string;
}

export interface MessageItem {
  _id: string;
  user: UserProfile;
  content: string;
  date: string;
  children: MessageChildItem[];
}

export interface VisitorItem {
  user: UserProfile;
  date?: string;
}

export interface CaptchaPayload {
  data: string;
  time: number;
}

export interface BasicResponse {
  code: number;
  msg: string;
}

export interface LoginStateResponse {
  userInfo?: UserProfile;
}

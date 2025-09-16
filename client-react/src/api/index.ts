import { postForm } from "./client";
import {
  ApiResponse,
  ArticleDetail,
  ArticleInfo,
  ArticleSummary,
  BasicResponse,
  CaptchaPayload,
  LoginStateResponse,
  MessageItem,
  UserProfile,
  VisitorItem
} from "@/types/api";

const TAGS = ["", "HTML&Css", "JavaScript", "Node", "Vue&React", "Other"] as const;

let articleSkip = 0;
const articleLimit = 5;

export const getArticleInfo = () => {
  return postForm<ApiResponse<ArticleInfo>>("/article/getInfo");
};

export const getArticleHot = (limit = 8) => {
  return postForm<ApiResponse<ArticleSummary[]>>("/article/getHot", { limit });
};

export const getArticleShow = (index = 0, ifFresh = false) => {
  if (ifFresh) {
    articleSkip = 0;
  }

  const tag = TAGS[index] ?? "";
  const payload = {
    skip: articleSkip,
    limit: articleLimit,
    tag
  };
  articleSkip += articleLimit;

  return postForm<ApiResponse<ArticleSummary[]>>("/article/getShow", payload);
};

export const getArticle = (_id: string) => {
  return postForm<ApiResponse<ArticleDetail>>("/article", { _id });
};

export const getArticleExtend = (tag: string) => {
  return postForm<ApiResponse<ArticleSummary[]>>("/article/extend", { tag });
};

export const getArticleSearch = (keywords: string) => {
  return postForm<ApiResponse<ArticleSummary[]>>("/article/search", { keywords });
};

export const getRegisterVCode = () => {
  return postForm<CaptchaPayload>("/register/vcode");
};

export const getRegisterCheckVcode = (svgCode: string) => {
  return postForm<BasicResponse>("/register/checkVcode", { svgCode });
};

export const postRegister = (options: Record<string, unknown>) => {
  return postForm<BasicResponse>("/register", options);
};

export const postLogin = (options: Record<string, unknown>) => {
  return postForm<BasicResponse & { userInfo?: UserProfile }>("/login", options);
};

export const postIfLogin = () => {
  return postForm<LoginStateResponse>("/login/ifLogin");
};

export const postLogout = () => {
  return postForm<BasicResponse>("/login/logout");
};

export const commitMessage = (options: Record<string, unknown>) => {
  return postForm<BasicResponse>("/message/commit", options);
};

export const commitChildMessage = (options: Record<string, unknown>) => {
  return postForm<BasicResponse>("/message/childCommit", options);
};

export const getMessageList = (skip = 0, limit = 5) => {
  return postForm<ApiResponse<MessageItem[]>>("/message/getList", { skip, limit });
};

export const getVisitor = () => {
  return postForm<ApiResponse<VisitorItem[]>>("/visitor");
};

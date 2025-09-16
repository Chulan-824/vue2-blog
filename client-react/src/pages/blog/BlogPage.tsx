import { LoadingOutlined, SearchOutlined } from "@ant-design/icons";
import { Spin } from "antd";
import { useCallback, useEffect, useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getArticleHot, getArticleInfo, getArticleSearch, getArticleShow, getVisitor } from "@/api";
import type { ArticleSummary, VisitorItem } from "@/types/api";

const CATEGORY_LIST = ["全部文章", "HTML&Css", "JavaScript", "Node", "Vue&React", "Other"];

export const BlogPage = () => {
  const params = useParams<{ id?: string }>();
  const navigate = useNavigate();

  const pageSize = 5;

  const categoryIndex = useMemo(() => {
    const index = Number(params.id ?? 0);
    if (Number.isNaN(index) || index < 0 || index >= CATEGORY_LIST.length) {
      return 0;
    }
    return index;
  }, [params.id]);

  const [articles, setArticles] = useState<ArticleSummary[]>([]);
  const [loading, setLoading] = useState(false);
  const [noMore, setNoMore] = useState(false);
  const [tags, setTags] = useState<string[]>([]);
  const [hotList, setHotList] = useState<ArticleSummary[]>([]);
  const [visitor, setVisitor] = useState<VisitorItem[]>([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState<ArticleSummary[]>([]);
  const [searching, setSearching] = useState(false);
  const [searchFixed, setSearchFixed] = useState(false);
  const [hoverIndex, setHoverIndex] = useState<number | null>(null);

  const categories = useMemo(() => [CATEGORY_LIST[0], ...tags], [tags]);

  const fetchArticles = useCallback(
    async (fresh: boolean) => {
      try {
        setLoading(true);
        const res = await getArticleShow(categoryIndex, fresh);
        const data = res.data?.data ?? [];
        if (fresh) {
          setArticles(data);
          setNoMore(data.length < pageSize);
        } else {
          if (data.length === 0) {
            setNoMore(true);
          } else {
            setArticles((prev) => [...prev, ...data]);
            if (data.length < pageSize) {
              setNoMore(true);
            }
          }
        }
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    },
    [categoryIndex]
  );

  useEffect(() => {
    getArticleInfo()
      .then((res) => {
        if (Array.isArray(res.data.data?.tags)) {
          setTags(res.data.data.tags);
        }
      })
      .catch((error) => console.error(error));

    getArticleHot()
      .then((res) => {
        if (Array.isArray(res.data.data)) {
          setHotList(res.data.data);
        }
      })
      .catch((error) => console.error(error));

    getVisitor()
      .then((res) => {
        if (Array.isArray(res.data.data)) {
          setVisitor(res.data.data);
        }
      })
      .catch((error) => console.error(error));
  }, []);

  useEffect(() => {
    setNoMore(false);
    setArticles([]);
    void fetchArticles(true);
    setHoverIndex(null);
    window.scrollTo({ top: 0, behavior: "smooth" });
  }, [categoryIndex, fetchArticles]);

  useEffect(() => {
    if (!searchTerm.trim()) {
      setSearchResults([]);
      return;
    }

    const timer = setTimeout(() => {
      setSearching(true);
      getArticleSearch(searchTerm.trim())
        .then((res) => {
          setSearchResults(res.data.data ?? []);
        })
        .catch((error) => {
          console.error(error);
          setSearchResults([]);
        })
        .finally(() => setSearching(false));
    }, 600);

    return () => clearTimeout(timer);
  }, [searchTerm]);

  useEffect(() => {
    const handleScroll = () => {
      const { scrollTop, scrollHeight, clientHeight } = document.documentElement;
      setSearchFixed(scrollTop >= 900);

      if (loading || noMore) {
        return;
      }

      if (scrollTop + clientHeight >= scrollHeight - 200) {
        void fetchArticles(false);
      }
    };

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, [fetchArticles, loading, noMore]);

  const activeIndex = hoverIndex ?? categoryIndex;
  const coverOffset = activeIndex * 40;

  return (
    <div className="flex flex-col gap-8 lg:flex-row">
      <div className="flex-1 space-y-6">
        {articles.map((article, index) => (
          <article key={article._id} className="relative overflow-hidden rounded-2xl bg-white shadow-lg">
            {index === 0 && (
              <span className="absolute left-[-24px] top-6 rotate-[-45deg] bg-orange-500 px-10 py-1 text-sm font-semibold text-white">
                置顶
              </span>
            )}
            <div className="border-b border-slate-200 px-6 pb-3 pt-6">
              <span className="mr-3 text-sm text-brand">【{article.type ?? article.tag}】</span>
              <Link to={`/Article/${article._id}`} className="text-lg font-semibold text-slate-800 hover:text-brand">
                {article.title}
              </Link>
            </div>
            <div className="flex flex-col gap-6 px-6 py-6 lg:flex-row">
              <Link
                to={`/Article/${article._id}`}
                className="block h-40 w-full overflow-hidden rounded-xl bg-slate-100 lg:h-36 lg:w-52"
              >
                <div
                  className="h-full w-full bg-cover bg-center transition-transform duration-700 hover:scale-110"
                  style={{ backgroundImage: `url(${article.surface})` }}
                />
              </Link>
              <p className="flex-1 text-sm leading-7 text-slate-600">{article.content}</p>
            </div>
            <div className="flex flex-wrap items-center justify-between gap-4 border-t border-slate-100 bg-slate-50 px-6 py-4 text-sm text-slate-500">
              <div className="flex items-center gap-2">
                <span className="font-medium text-brand">#{article.tag}</span>
              </div>
              <div className="flex items-center gap-4">
                <span>浏览 {article.pv}</span>
                <span>评论 {article.comment?.length ?? 0}</span>
                <span>
                  {new Date(article.date).getFullYear()}年{new Date(article.date).getMonth() + 1}月
                  {new Date(article.date).getDate()}日
                </span>
              </div>
            </div>
          </article>
        ))}

        {loading && (
          <div className="flex items-center justify-center py-6 text-slate-500">
            <Spin indicator={<LoadingOutlined spin className="text-brand" />} />
            <span className="ml-3 text-sm">加载中</span>
          </div>
        )}

        {noMore && (
          <p className="py-6 text-center text-sm text-slate-400">哼╭(╯^╰)╮我也是有底线的！！</p>
        )}
      </div>

      <aside className="w-full space-y-6 lg:w-80">
        <div
          className={`space-y-4 rounded-2xl bg-white p-6 shadow-lg transition-all duration-300 ${
            searchFixed ? "lg:sticky lg:top-24" : ""
          }`}
        >
          <div className="relative">
            <input
              type="text"
              value={searchTerm}
              onChange={(event) => setSearchTerm(event.target.value)}
              placeholder="请输入搜索内容"
              className="w-full rounded-full border border-slate-200 bg-slate-50 py-3 pl-12 pr-4 text-sm outline-none focus:border-brand"
            />
            <SearchOutlined className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400" />
            {searching && (
              <LoadingOutlined className="absolute right-4 top-1/2 -translate-y-1/2 text-brand" spin />
            )}
            {searchResults.length > 0 && (
              <ul className="absolute left-0 right-0 top-full z-10 mt-2 max-h-60 overflow-y-auto rounded-xl border border-slate-200 bg-white shadow-lg">
                {searchResults.map((item) => (
                  <li key={item._id}>
                    <Link
                      to={`/Article/${item._id}`}
                      className="block px-4 py-2 text-sm text-slate-600 hover:bg-brand/10 hover:text-brand"
                      onClick={() => setSearchResults([])}
                    >
                      {item.title}
                    </Link>
                  </li>
                ))}
              </ul>
            )}
          </div>

          <div className="relative overflow-hidden rounded-2xl border border-slate-200 bg-white">
            <ul>
              {categories.map((item, index) => (
                <li key={item}>
                  <button
                    type="button"
                    className={`relative flex h-10 w-full items-center justify-between px-4 text-sm transition ${
                      activeIndex === index ? "text-brand" : "text-slate-600 hover:text-brand"
                    }`}
                    onMouseEnter={() => setHoverIndex(index)}
                    onMouseLeave={() => setHoverIndex(null)}
                    onClick={() => navigate(`/blog/${index}`)}
                  >
                    {item}
                  </button>
                </li>
              ))}
            </ul>
            <span
              className="pointer-events-none absolute left-0 h-10 w-full -z-10 bg-brand/10 transition-all duration-200"
              style={{ transform: `translateY(${coverOffset}px)` }}
            />
          </div>
        </div>

        <div className="space-y-4 rounded-2xl bg-white p-6 shadow-lg">
          <h3 className="text-lg font-semibold text-slate-700">热门文章</h3>
          <ol className="space-y-3 text-sm text-slate-600">
            {hotList.map((item, index) => (
              <li key={item._id} className="flex items-start gap-3">
                <span className="mt-0.5 inline-flex h-6 w-6 items-center justify-center rounded-full bg-brand/10 text-xs font-semibold text-brand">
                  {index + 1}
                </span>
                <Link to={`/Article/${item._id}`} className="hover:text-brand">
                  {item.title}
                </Link>
              </li>
            ))}
          </ol>
        </div>

        <div className="space-y-3 rounded-2xl bg-white p-6 shadow-lg">
          <h3 className="text-lg font-semibold text-slate-700">置顶推荐</h3>
          {hotList[0] ? (
            <Link to={`/Article/${hotList[0]._id}`} className="flex items-center gap-3 text-sm text-brand hover:text-brand-dark">
              <span className="inline-flex h-6 w-6 items-center justify-center rounded-full bg-brand/10 text-xs font-semibold text-brand">
                1
              </span>
              {hotList[0].title}
            </Link>
          ) : (
            <p className="text-sm text-slate-400">暂无推荐</p>
          )}
        </div>

        <div className="space-y-3 rounded-2xl bg-white p-6 shadow-lg">
          <h3 className="text-lg font-semibold text-slate-700">最近访客</h3>
          <div className="grid grid-cols-4 gap-3">
            {visitor.map((item) => (
              <div
                key={item.user._id}
                className="flex h-16 w-16 flex-col items-center justify-end overflow-hidden rounded-xl bg-slate-200 text-xs text-white"
                style={{ backgroundImage: `url(${item.user.photo})`, backgroundSize: "cover", backgroundPosition: "center" }}
              >
                <span className="w-full bg-black/60 px-1 py-0.5 text-center">
                  {item.user.user}
                </span>
              </div>
            ))}
          </div>
        </div>
      </aside>
    </div>
  );
};

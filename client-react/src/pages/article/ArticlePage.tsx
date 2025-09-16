import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getArticle, getArticleExtend } from "@/api";
import type { ArticleDetail, ArticleSummary } from "@/types/api";

const formatFullDateTime = (value: string) => {
  const date = new Date(value);
  const pad = (num: number) => (num < 10 ? `0${num}` : String(num));
  return `${date.getFullYear()}/${pad(date.getMonth() + 1)}/${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
};

export const ArticlePage = () => {
  const params = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [article, setArticle] = useState<ArticleDetail | null>(null);
  const [extendList, setExtendList] = useState<ArticleSummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const articleId = params.id;
    if (!articleId) {
      navigate("/404", { replace: true });
      return;
    }

    setLoading(true);
    getArticle(articleId)
      .then((res) => {
        if (res.data.code === 0) {
          const info = res.data.data;
          setArticle(info);
          return getArticleExtend(info.tag)
            .then((extendRes) => {
              setExtendList(extendRes.data.data ?? []);
            })
            .catch((error) => console.error(error));
        }
        navigate("/404", { replace: true });
      })
      .catch((error) => {
        console.error(error);
        navigate("/404", { replace: true });
      })
      .finally(() => setLoading(false));
  }, [navigate, params.id]);

  const timeInfo = useMemo(() => {
    if (!article) {
      return null;
    }
    return {
      day: new Date(article.date).getDate(),
      month: new Date(article.date).getMonth() + 1,
      year: new Date(article.date).getFullYear()
    };
  }, [article]);

  if (loading) {
    return <div className="py-20 text-center text-slate-500">文章加载中...</div>;
  }

  if (!article) {
    return null;
  }

  return (
    <div className="space-y-8">
      <header className="space-y-4 border-b border-slate-200 pb-6">
        <h1 className="text-2xl font-semibold text-slate-800">{article.title}</h1>
        <div className="flex flex-wrap items-center gap-4 text-sm text-slate-500">
          <span>
            作者：
            <span className="text-brand">楚岚</span>
          </span>
          <span>围观群众：{article.pv}</span>
          <span>更新于 {formatFullDateTime(article.updateDate ?? article.date)}</span>
        </div>
      </header>

      {timeInfo && (
        <div className="flex items-center gap-2 text-slate-500">
          <span className="text-4xl font-bold text-brand">{timeInfo.day}</span>
          <div className="text-sm leading-5">
            <div>{timeInfo.month} 月</div>
            <div>{timeInfo.year}</div>
          </div>
        </div>
      )}

      <article
        className="prose prose-slate max-w-none prose-a:text-brand hover:prose-a:text-brand-dark"
        dangerouslySetInnerHTML={{ __html: article.content }}
      />

      <section className="rounded-2xl bg-slate-50 p-6 text-sm text-slate-600">
        <p>非特殊说明，本文版权归 楚岚 所有，转载请注明出处.</p>
        <p className="mt-2">
          本文标题：
          <span className="text-brand">楚岚</span>
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-slate-700">延伸阅读</h2>
        <ol className="space-y-2 text-sm text-slate-600">
          {extendList.length === 0 && <li className="text-slate-400">暂无更多推荐</li>}
          {extendList.map((item) => (
            <li key={item._id}>
              <Link to={`/Article/${item._id}`} className="hover:text-brand">
                {item.title}
              </Link>
            </li>
          ))}
        </ol>
      </section>
    </div>
  );
};

import { DownOutlined, LinkOutlined } from "@ant-design/icons";
import { useEffect, useMemo, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import type { ArticleSummary } from "@/types/api";
import { getArticleHot } from "@/api";

const useIntersection = () => {
  const [isVisible, setVisible] = useState(false);
  const ref = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const node = ref.current;
    if (!node) {
      return;
    }

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setVisible(true);
          }
        });
      },
      { threshold: 0.2 }
    );

    observer.observe(node);

    return () => observer.disconnect();
  }, []);

  return { ref, isVisible } as const;
};

export const HomePage = () => {
  const [menuOpen, setMenuOpen] = useState(false);
  const [hotArticles, setHotArticles] = useState<ArticleSummary[]>([]);
  const navigate = useNavigate();

  const heroBackground = useMemo(
    () =>
      [
        "radial-gradient(circle at 20% 20%, rgba(22, 119, 255, 0.35), transparent 60%)",
        "radial-gradient(circle at 80% 0%, rgba(14, 165, 233, 0.3), transparent 55%)",
        "radial-gradient(circle at 50% 80%, rgba(79, 70, 229, 0.25), transparent 60%)",
        "linear-gradient(135deg, #0f172a, #1e293b)"
      ].join(","),
    []
  );

  const { ref: heroRef, isVisible: heroVisible } = useIntersection();
  const { ref: hotRef, isVisible: hotVisible } = useIntersection();

  useEffect(() => {
    getArticleHot(3)
      .then((res) => {
        if (Array.isArray(res.data.data)) {
          setHotArticles(res.data.data);
        }
      })
      .catch((error) => {
        console.error(error);
      });
  }, []);

  const introText = useMemo(
    () => "真常应物,真常得性;常应常静,常清静矣",
    []
  );

  return (
    <div className="relative min-h-screen overflow-hidden text-white">
      <section
        className="relative flex min-h-screen flex-col items-center justify-center px-6 text-center"
        style={{
          backgroundImage: heroBackground,
          backgroundSize: "cover",
          backgroundPosition: "center",
          backgroundAttachment: "fixed"
        }}
      >
        <div className="absolute inset-0 bg-black/35" />
        <button
          type="button"
          aria-label="打开菜单"
          onClick={() => setMenuOpen((prev) => !prev)}
          className={`group absolute right-10 top-10 z-30 flex h-12 w-12 flex-col items-center justify-center gap-1.5 rounded-full bg-white/10 transition hover:bg-brand/70 lg:right-16`}
        >
          <span
            className={`h-0.5 w-6 rounded-full bg-white transition-transform ${
              menuOpen ? "translate-y-[6px] rotate-45" : ""
            }`}
          />
          <span
            className={`h-0.5 w-6 rounded-full bg-white transition-opacity ${
              menuOpen ? "opacity-0" : "opacity-100"
            }`}
          />
          <span
            className={`h-0.5 w-6 rounded-full bg-white transition-transform ${
              menuOpen ? "-translate-y-[6px] -rotate-45" : ""
            }`}
          />
        </button>

        <div
          className={`relative z-10 flex flex-col items-center gap-6 transition-all duration-700 ${
            heroVisible ? "scale-100 opacity-100" : "scale-90 opacity-0"
          }`}
          ref={heroRef}
        >
          <h1 className="text-5xl font-light tracking-[0.3em] text-white drop-shadow-lg lg:text-6xl">
            楚 岚
          </h1>
          <p className="text-lg tracking-[0.4em] text-slate-100 lg:text-xl">{introText}</p>
          <Link
            to="/blog/0"
            className="rounded-full bg-brand px-8 py-3 text-sm font-medium uppercase tracking-widest text-white shadow-lg transition hover:bg-brand-dark"
          >
            Enter Blog
          </Link>
        </div>

        <button
          type="button"
          className="absolute bottom-10 left-1/2 z-20 flex h-12 w-12 -translate-x-1/2 items-center justify-center rounded-full border border-white/60 text-white transition hover:bg-brand"
          onClick={() => {
            window.scrollTo({ top: window.innerHeight, behavior: "smooth" });
          }}
        >
          <DownOutlined className="text-xl" />
        </button>

        <div
          className={`fixed inset-0 z-20 bg-slate-900/60 transition-opacity duration-300 ${
            menuOpen ? "pointer-events-auto opacity-100" : "pointer-events-none opacity-0"
          }`}
          onClick={() => setMenuOpen(false)}
        />

        <div
          className={`absolute inset-y-0 right-0 z-30 w-full max-w-md translate-x-full bg-white/90 p-10 text-right text-slate-700 shadow-2xl transition-transform duration-300 ease-in-out ${
            menuOpen ? "translate-x-0" : ""
          }`}
        >
          <div className="flex flex-col gap-6 text-xl">
            <button
              type="button"
              onClick={() => {
                setMenuOpen(false);
                navigate("/");
              }}
              className="hover:text-brand"
            >
              首页
            </button>
            <button
              type="button"
              onClick={() => {
                setMenuOpen(false);
                navigate("/blog/0");
              }}
              className="hover:text-brand"
            >
              博客
            </button>
            <button
              type="button"
              onClick={() => {
                setMenuOpen(false);
                navigate("/message");
              }}
              className="hover:text-brand"
            >
              留言
            </button>
          </div>
          <div className="mt-12 text-right font-handwriting text-4xl text-slate-200">Jack</div>
        </div>
      </section>

      <section className="relative z-10 w-full bg-white py-20 text-slate-800">
        <div className="mx-auto flex w-full max-w-6xl flex-col gap-12 px-6" ref={hotRef}>
          <div
            className={`mx-auto max-w-2xl text-center transition-all duration-700 ${
              hotVisible ? "translate-y-0 opacity-100" : "-translate-y-10 opacity-0"
            }`}
          >
            <h2 className="text-3xl font-semibold tracking-wide text-slate-800">热门文章</h2>
            <p className="mt-4 text-base text-slate-500">
              螃蟹在剥我的壳，笔记本在写我<br />漫天的我落在枫叶上雪花上，而你在想我。
            </p>
          </div>
          <div className="grid gap-8 md:grid-cols-3">
            {hotArticles.map((article, index) => (
              <article
                key={article._id}
                className={`flex flex-col overflow-hidden rounded-2xl bg-white shadow-lg transition duration-700 hover:-translate-y-1 hover:shadow-2xl ${
                  hotVisible ? "opacity-100" : "opacity-0"
                }`}
                style={{ transitionDelay: `${index * 120}ms` }}
              >
                <Link to={`/Article/${article._id}`} className="relative block h-48 overflow-hidden">
                  <div
                    className="absolute inset-0 bg-cover bg-center transition-transform duration-700 hover:scale-110"
                    style={{ backgroundImage: `url(${article.surface})` }}
                  />
                </Link>
                <div className="flex flex-1 flex-col gap-4 p-6">
                  <Link to={`/Article/${article._id}`} className="text-lg font-semibold text-slate-800 hover:text-brand">
                    {article.title}
                  </Link>
                  <div className="text-xs uppercase tracking-widest text-brand">{article.tag}</div>
                  <p className="line-clamp-3 text-sm text-slate-600">{article.content}</p>
                  <Link
                    to={`/Article/${article._id}`}
                    className="mt-auto inline-flex items-center gap-2 text-sm font-medium text-brand hover:text-brand-dark"
                  >
                    阅读更多
                    <LinkOutlined />
                  </Link>
                </div>
              </article>
            ))}
          </div>
        </div>
      </section>

      <section className="relative z-10 bg-transparent py-20">
        <div className="mx-auto flex w-full max-w-4xl flex-col gap-6 px-6 text-center">
          <div className="grid gap-6 md:grid-cols-2">
            <Link
              to="/about"
              className="rounded-2xl border border-white/40 bg-white/10 p-10 text-2xl font-semibold backdrop-blur transition hover:border-brand hover:bg-brand/80"
            >
              关于我
            </Link>
            <Link
              to="/links"
              className="rounded-2xl border border-white/40 bg-white/10 p-10 text-2xl font-semibold backdrop-blur transition hover:border-brand hover:bg-brand/80"
            >
              + 友情链接
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
};

import { Link } from "react-router-dom";

export const NotFoundPage = () => {
  return (
    <div className="flex flex-col items-center justify-center gap-6 py-20">
      <div className="relative overflow-hidden rounded-3xl border border-dashed border-brand/40 bg-white/80 px-16 py-14 text-center shadow-xl">
        <div className="absolute inset-0 -z-10 bg-gradient-to-br from-brand/10 via-sky-200/40 to-brand/5" />
        <div className="text-6xl font-bold tracking-[0.6em] text-brand">404</div>
        <p className="mt-4 text-base text-slate-600">抱歉，你访问的页面不存在。</p>
      </div>
      <p className="text-sm text-slate-500">请检查链接是否正确，或者返回首页继续浏览。</p>
      <Link
        to="/"
        className="rounded-full bg-brand px-6 py-2 text-sm font-medium text-white shadow hover:bg-brand-dark"
      >
        返回首页
      </Link>
    </div>
  );
};

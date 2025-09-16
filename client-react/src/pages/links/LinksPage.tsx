import { CheckCircleOutlined, CloseCircleOutlined } from "@ant-design/icons";
import { BubbleCanvas } from "@/components/BubbleCanvas";

const links = Array.from({ length: 16 }).map(() => ({
  name: "欢迎加入",
  href: "https://www.baidu.com",
  des: "热爱并坚持使用.欢迎分享"
}));

export const LinksPage = () => {
  return (
    <div className="space-y-10">
      <section className="relative overflow-hidden rounded-2xl shadow-lg">
        <BubbleCanvas className="h-[260px] w-full" />
        <div className="absolute inset-0 flex flex-col items-center justify-center bg-slate-900/40 text-center text-white">
          <h2 className="text-3xl font-semibold">友情链接</h2>
          <p className="mt-3 text-base">真常应物,真常得性;常应常静,常清静矣</p>
        </div>
      </section>

      <section className="space-y-6 rounded-2xl bg-white p-6 shadow-lg">
        <div className="space-y-4 text-sm leading-7 text-slate-600">
          <h3 className="border-l-4 border-brand pl-4 text-lg font-semibold text-slate-800">友链说明</h3>
          <p className="flex flex-wrap items-center gap-2">
            <span className="flex items-center gap-1 text-red-500">
              <CloseCircleOutlined /> 经常宕机
            </span>
            <span className="flex items-center gap-1 text-red-500">
              <CloseCircleOutlined /> 不合法规
            </span>
            <span className="flex items-center gap-1 text-red-500">
              <CloseCircleOutlined /> 插边球站
            </span>
            <span className="flex items-center gap-1 text-red-500">
              <CloseCircleOutlined /> 红标报毒
            </span>
            <span className="flex items-center gap-1 text-green-500">
              <CheckCircleOutlined /> 原创优先
            </span>
            <span className="flex items-center gap-1 text-green-500">
              <CheckCircleOutlined /> 技术优先
            </span>
          </p>
          <p>
            交换友链可在留言板留言. 本站链接如下：<br />
            名称：楚岚<br />网址：https://www.chulan.fun<br />图标：https://www.chulan.fun<br />描述：楚岚·一个认真的普通人
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {links.map((item, index) => (
            <a
              className="flex h-full flex-col rounded-2xl border border-slate-200 bg-slate-50 p-5 transition hover:-translate-y-1 hover:border-brand hover:shadow-lg"
              href={item.href}
              target="_blank"
              rel="noreferrer"
              key={`${item.href}-${index}`}
            >
              <div className="flex items-center gap-3">
                <span className="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-brand/10 via-sky-200/70 to-brand/40 text-lg font-semibold text-brand-dark shadow-inner">
                  {item.name.slice(0, 1)}
                </span>
                <h3 className="text-base font-semibold text-slate-800">{item.name}</h3>
              </div>
              <p className="mt-3 line-clamp-3 text-sm text-slate-600">{item.des}</p>
            </a>
          ))}
        </div>
      </section>
    </div>
  );
};

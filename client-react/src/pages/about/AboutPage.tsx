import { BubbleCanvas } from "@/components/BubbleCanvas";

export const AboutPage = () => {
  return (
    <div className="space-y-10">
      <section className="relative overflow-hidden rounded-2xl shadow-lg">
        <BubbleCanvas className="h-[260px] w-full" />
        <div className="absolute inset-0 flex flex-col items-center justify-center bg-slate-900/40 text-center text-white">
          <h2 className="text-3xl font-semibold">关于</h2>
          <p className="mt-3 text-base">真常应物,真常得性;常应常静,常清静矣</p>
        </div>
      </section>

      <section className="space-y-6 rounded-2xl bg-white p-6 shadow-lg">
        <article className="space-y-10 text-sm leading-7 text-slate-600">
          <section>
            <h3 className="border-l-4 border-brand pl-4 text-lg font-semibold text-slate-800">关于我</h3>
            <p className="mt-4">
              毕业后在长沙做在线教育课程顾问，很迷茫。2020年决定转行做开发，然沉迷撸码，日渐消瘦。
            </p>
            <p className="mt-4">可以通过以下方式联系到我：</p>
            <ul className="mt-3 space-y-2">
              <li>
                邮 箱 ：
                <a href="mailto:390226630@qq.com" className="text-brand">
                  390226630@qq.com
                </a>
              </li>
            </ul>
          </section>

          <section>
            <h3 className="border-l-4 border-brand pl-4 text-lg font-semibold text-slate-800">关于本站</h3>
            <p className="mt-4">本站建于2020年1月，主要是个人爱好写着玩。</p>
            <p className="mt-4">本站结构：</p>
            <ul className="mt-3 space-y-2">
              <li>
                前 端 ：<code className="rounded bg-slate-100 px-2 py-1">Layui + ElementUI</code>
              </li>
              <li>
                后 端 ：<code className="rounded bg-slate-100 px-2 py-1">nodeJS + MongoDB</code>
              </li>
            </ul>
            <p className="mt-4">本站采用阿里云提供的服务器ESC和存储对象OSS。</p>
          </section>

          <section>
            <h3 className="border-l-4 border-brand pl-4 text-lg font-semibold text-slate-800">关于版权</h3>
            <p className="mt-4">
              本站采用
              <a
                href="https://creativecommons.org/licenses/by-nc/4.0/deed.zh"
                target="_blank"
                rel="noreferrer"
                className="text-brand"
              >
                「署名-非商业性使用 4.0国际 (CC BY-NC 4.0)」
              </a>
              创作共享协议。只要在使用时注明出处，那么您可以对本站所有原创内容进行转载、节选、二次创作，但是您不得对其用于商业目的。
            </p>
          </section>

          <section>
            <h3 className="border-l-4 border-brand pl-4 text-lg font-semibold text-slate-800">特别说明</h3>
            <ul className="mt-3 space-y-2">
              <li>本站文章仅代表个人观点，和任何组织或个人无关。</li>
              <li>本站前端开发代码没有考虑对IE浏览器的兼容。</li>
            </ul>
            <div className="relative mt-6 h-80 w-full overflow-hidden rounded-2xl border border-brand/20 bg-white/60 shadow-inner">
              <BubbleCanvas className="absolute inset-0 h-full w-full opacity-80" height={320} />
              <div className="absolute inset-0 bg-gradient-to-br from-white/40 via-sky-100/50 to-white/10 backdrop-blur-[2px]" />
              <div className="relative flex h-full flex-col items-center justify-center gap-2 text-center">
                <span className="text-xl font-semibold text-brand-dark">你好，朋友</span>
                <p className="max-w-sm text-sm text-slate-600">
                  感谢来到楚岚的小站，愿一切美好与你不期而遇。
                </p>
              </div>
            </div>
          </section>
        </article>
      </section>
    </div>
  );
};

import { message } from "antd";
import { useCallback, useEffect, useState } from "react";
import { commitChildMessage, commitMessage, getMessageList } from "@/api";
import type { MessageItem } from "@/types/api";
import { RichTextEditor } from "@/components/RichTextEditor";
import { useAuth } from "@/contexts/AuthContext";

interface CommentWithState extends MessageItem {
  reply: {
    content: string;
    reUser: string;
    visible: boolean;
    childIndex: number | null;
  };
}

const mapComment = (item: MessageItem): CommentWithState => ({
  ...item,
  reply: {
    content: "",
    reUser: item.user.user,
    visible: false,
    childIndex: null
  }
});

const formatDateTime = (value: string) => {
  const date = new Date(value);
  const pad = (num: number) => (num < 10 ? `0${num}` : String(num));
  return `${date.getFullYear()}年${pad(date.getMonth() + 1)}月${pad(date.getDate())}日 ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
};

export const MessagePage = () => {
  const { requireLogin } = useAuth();
  const [comments, setComments] = useState<CommentWithState[]>([]);
  const [limit, setLimit] = useState(5);
  const [loading, setLoading] = useState(false);
  const [noMore, setNoMore] = useState(false);
  const [submittingMessage, setSubmittingMessage] = useState(false);
  const [submittingReplyId, setSubmittingReplyId] = useState<string | null>(null);

  const fetchComments = useCallback(async (limitValue: number) => {
    try {
      setLoading(true);
      const res = await getMessageList(0, limitValue);
      const data = res.data.data ?? [];
      setComments(data.map(mapComment));
      setNoMore(data.length < limitValue);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void fetchComments(limit);
  }, [fetchComments, limit]);

  useEffect(() => {
    const handleScroll = () => {
      const { scrollTop, scrollHeight, clientHeight } = document.documentElement;
      if (loading || noMore) {
        return;
      }
      if (scrollTop + clientHeight >= scrollHeight - 200) {
        setLimit((prev) => prev + 5);
      }
    };

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, [loading, noMore]);

  const handleMessageSubmit = useCallback(
    async (value: string) => {
      const user = await requireLogin();
      if (!user) {
        return;
      }

      setSubmittingMessage(true);
      commitMessage({ user: user._id, content: value })
        .then((res) => {
          if (res.data.code === 0) {
            message.success(res.data.msg ?? "留言成功");
            setLimit(5);
            void fetchComments(5);
          } else {
            message.error(res.data.msg ?? "留言失败，请稍后再试");
          }
        })
        .catch((error) => {
          console.error(error);
          message.error("留言失败，请稍后再试");
        })
        .finally(() => setSubmittingMessage(false));
    },
    [fetchComments, requireLogin]
  );

  const toggleReply = useCallback((parentIndex: number, childIndex?: number) => {
    setComments((prev) =>
      prev.map((item, index) => {
        if (index !== parentIndex) {
          return {
            ...item,
            reply: { ...item.reply, visible: false, content: "", childIndex: null, reUser: item.user.user }
          };
        }
        const visible = item.reply.visible && item.reply.childIndex === (childIndex ?? null);
        return {
          ...item,
          reply: {
            content: visible ? "" : item.reply.content,
            visible: !visible,
            childIndex: childIndex ?? null,
            reUser:
              childIndex !== undefined && item.children[childIndex]
                ? item.children[childIndex].user.user
                : item.user.user
          }
        };
      })
    );
  }, []);

  const updateReplyContent = useCallback((parentIndex: number, value: string) => {
    setComments((prev) =>
      prev.map((item, index) =>
        index === parentIndex
          ? { ...item, reply: { ...item.reply, content: value } }
          : item
      )
    );
  }, []);

  const handleReplySubmit = useCallback(
    async (parentIndex: number) => {
      const user = await requireLogin();
      if (!user) {
        return;
      }

      const target = comments[parentIndex];
      if (!target || !target.reply.content.trim()) {
        message.warning("请输入回复内容");
        return;
      }

      setSubmittingReplyId(target._id);
      commitChildMessage({
        parentId: target._id,
        user: user._id,
        content: target.reply.content,
        reUser: target.reply.reUser
      })
        .then((res) => {
          if (res.data.code === 0) {
            message.success(res.data.msg ?? "回复成功");
            void fetchComments(limit);
          } else {
            message.error(res.data.msg ?? "回复失败，请稍后再试");
          }
        })
        .catch((error) => {
          console.error(error);
          message.error("回复失败，请稍后再试");
        })
        .finally(() => setSubmittingReplyId(null));
    },
    [comments, fetchComments, limit, requireLogin]
  );

  return (
    <div className="space-y-10">
      <section className="space-y-4 text-center">
        <h1 className="text-3xl font-semibold text-slate-800">留言板</h1>
        <p className="text-slate-500">沟通交流，拉近你我！</p>
        <div className="mx-auto max-w-3xl">
          <RichTextEditor onSubmit={handleMessageSubmit} submitting={submittingMessage} />
        </div>
      </section>

      <section className="rounded-2xl bg-white p-6 shadow-lg">
        <ul className="space-y-6">
          {comments.map((item, index) => (
            <li key={item._id} className="space-y-4 border-b border-slate-100 pb-6 last:border-b-0">
              <div className="flex gap-4">
                <div
                  className="h-12 w-12 flex-shrink-0 rounded-full bg-cover bg-center"
                  style={{ backgroundImage: `url(${item.user.photo})` }}
                />
                <div className="flex-1 space-y-3">
                  <div className="flex items-center justify-between">
                    <span className="font-medium text-brand">{item.user.user}</span>
                    <span className="text-xs text-slate-400">{formatDateTime(item.date)}</span>
                  </div>
                  <div
                    className="rounded-xl bg-slate-50 p-4 text-sm leading-6 text-slate-600"
                    dangerouslySetInnerHTML={{ __html: item.content }}
                  />
                  <button
                    type="button"
                    className="text-xs font-medium text-brand hover:text-brand-dark"
                    onClick={() => toggleReply(index)}
                  >
                    {item.reply.visible ? "收起" : "回复"}
                  </button>
                </div>
              </div>

              <div className="space-y-3 pl-16">
                {item.children.map((child, childIndex) => (
                  <div key={`${child._id ?? childIndex}-${child.date}`} className="flex gap-3">
                    <div
                      className="h-10 w-10 flex-shrink-0 rounded-full bg-cover bg-center"
                      style={{ backgroundImage: `url(${child.user.photo})` }}
                    />
                    <div className="flex-1 space-y-1 rounded-xl bg-slate-50 p-3 text-sm leading-6 text-slate-600">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="font-medium text-brand">{child.user.user}</span>
                        <span className="text-slate-400">回复</span>
                        <span className="font-medium text-brand">{child.reUser}</span>
                      </div>
                      <div>{child.content}</div>
                      <div className="flex items-center justify-between text-xs text-slate-400">
                        <span>{formatDateTime(child.date)}</span>
                        <button
                          type="button"
                          className="font-medium text-brand hover:text-brand-dark"
                          onClick={() => toggleReply(index, childIndex)}
                        >
                          {item.reply.visible && item.reply.childIndex === childIndex ? "收起" : "回复"}
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {item.reply.visible && (
                <div className="space-y-3 rounded-2xl bg-slate-50 p-4 pl-16">
                  <textarea
                    value={item.reply.content}
                    onChange={(event) => updateReplyContent(index, event.target.value)}
                    className="h-24 w-full resize-none rounded-xl border border-slate-200 p-3 text-sm text-slate-600 focus:border-brand focus:outline-none"
                    placeholder={`回复【${item.reply.reUser}】：`}
                  />
                  <button
                    type="button"
                    className="self-end rounded-full bg-brand px-5 py-2 text-sm font-medium text-white transition hover:bg-brand-dark"
                    onClick={() => void handleReplySubmit(index)}
                    disabled={submittingReplyId === item._id}
                  >
                    {submittingReplyId === item._id ? "提交中..." : "提交"}
                  </button>
                </div>
              )}
            </li>
          ))}
        </ul>

        {loading && <p className="py-4 text-center text-sm text-slate-400">留言加载中...</p>}
        {noMore && <p className="py-4 text-center text-sm text-slate-300">没有更多留言了~</p>}
      </section>
    </div>
  );
};

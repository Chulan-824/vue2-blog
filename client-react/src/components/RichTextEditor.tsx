import { SmileOutlined } from "@ant-design/icons";
import { Button, Space } from "antd";
import { useCallback, useRef, useState } from "react";

interface RichTextEditorProps {
  onSubmit: (value: string) => void;
  submitting?: boolean;
}

export const RichTextEditor = ({ onSubmit, submitting }: RichTextEditorProps) => {
  const editorRef = useRef<HTMLDivElement | null>(null);
  const [hasContent, setHasContent] = useState(false);

  const handleInput = useCallback(() => {
    const html = editorRef.current?.innerHTML ?? "";
    setHasContent(html.trim().length > 0);
  }, []);

  const insertEmoji = useCallback((emoji: string) => {
    editorRef.current?.focus();
    document.execCommand("insertText", false, emoji);
    handleInput();
  }, [handleInput]);

  const handleSubmit = useCallback(() => {
    const html = editorRef.current?.innerHTML ?? "";
    if (!html.trim()) {
      return;
    }
    onSubmit(html.trim());
    if (editorRef.current) {
      editorRef.current.innerHTML = "";
    }
    setHasContent(false);
  }, [onSubmit]);

  return (
    <div className="space-y-4 rounded-2xl border border-slate-200 bg-slate-50 p-4">
      <Space>
        <Button
          type="text"
          icon={<SmileOutlined className="text-lg" />}
          onClick={() => insertEmoji("😊")}
        >
          表情
        </Button>
      </Space>
      <div className="relative">
        <div
          ref={editorRef}
          onInput={handleInput}
          className="min-h-[150px] rounded-xl border border-slate-200 bg-white px-4 py-3 text-left text-sm leading-6 text-slate-600 focus:outline-none"
          contentEditable
          suppressContentEditableWarning
        />
        {!hasContent && (
          <div className="pointer-events-none absolute left-5 top-3 text-sm text-slate-400">
            说点什么吧，支持插入表情哦~
          </div>
        )}
      </div>
      <div className="text-right">
        <Button type="primary" onClick={handleSubmit} loading={Boolean(submitting)}>
          提交留言
        </Button>
      </div>
    </div>
  );
};

import { PlusOutlined } from "@ant-design/icons";
import { Modal, Upload, message } from "antd";
import type { UploadChangeParam } from "antd/es/upload";
import type { UploadFile } from "antd/es/upload/interface";
import { useCallback, useMemo, useState } from "react";

interface AvatarModalProps {
  open: boolean;
  onCancel: () => void;
  onUploaded: () => Promise<void> | void;
}

const AVATAR_UPLOAD_ACTION = "http://47.96.127.142:80/upload/avatar";

export const AvatarModal = ({ open, onCancel, onUploaded }: AvatarModalProps) => {
  const [preview, setPreview] = useState<string>("");
  const [uploading, setUploading] = useState(false);

  const beforeUpload = useCallback((file: File) => {
    const isValidType = /^(image\/jpeg|image\/png|image\/gif)$/.test(file.type);
    if (!isValidType) {
      message.error("仅支持 JPG/PNG/GIF 格式的图片");
      return Upload.LIST_IGNORE;
    }

    const isLt50K = file.size / 1024 < 50;
    if (!isLt50K) {
      message.error("图片大小不能超过 50KB");
      return Upload.LIST_IGNORE;
    }

    return true;
  }, []);

  const handleChange = useCallback(
    async ({ file, fileList }: UploadChangeParam<UploadFile>) => {
      if (file.status === "uploading") {
        setUploading(true);
        return;
      }

      if (file.status === "done") {
        setUploading(false);
        if (file.originFileObj) {
          setPreview(URL.createObjectURL(file.originFileObj));
        }
        message.success("头像上传成功");
        await onUploaded();
        return;
      }

      if (file.status === "error") {
        setUploading(false);
        message.error("头像上传失败，请稍后再试");
      }

      if (fileList.length === 0) {
        setPreview("");
      }
    },
    [onUploaded]
  );

  const uploadButton = useMemo(
    () => (
      <div className="flex h-44 w-44 flex-col items-center justify-center rounded-lg border border-dashed border-slate-300 text-slate-500">
        <PlusOutlined className="text-2xl" />
        <span className="mt-2 text-sm">点击上传</span>
      </div>
    ),
    []
  );

  return (
    <Modal open={open} title="头像上传" onCancel={onCancel} footer={null} centered>
      <div className="flex w-full justify-center py-4">
        <Upload
          name="avatar"
          listType="picture-card"
          className="avatar-uploader"
          showUploadList={false}
          action={AVATAR_UPLOAD_ACTION}
          withCredentials
          beforeUpload={beforeUpload}
          onChange={handleChange}
        >
          {preview ? (
            <img src={preview} alt="avatar" className="h-44 w-44 rounded-lg object-cover" />
          ) : (
            uploadButton
          )}
        </Upload>
      </div>
      {uploading && <p className="text-center text-sm text-slate-500">上传中，请稍候...</p>}
    </Modal>
  );
};

import { LockOutlined, UserOutlined } from "@ant-design/icons";
import { Button, Form, Input, Modal, Typography, message } from "antd";
import type { AxiosResponse } from "axios";
import { useCallback, useState } from "react";

interface LoginModalProps {
  open: boolean;
  onCancel: () => void;
  onRegister: () => void;
  onSuccess: () => Promise<void> | void;
  postLogin: (options: Record<string, unknown>) => Promise<AxiosResponse<{ code: number; msg: string }>>;
}

interface LoginFormValues {
  user: string;
  pwd: string;
}

export const LoginModal = ({ open, onCancel, onRegister, onSuccess, postLogin }: LoginModalProps) => {
  const [form] = Form.useForm<LoginFormValues>();
  const [submitting, setSubmitting] = useState(false);

  const handleFinish = useCallback(
    async (values: LoginFormValues) => {
      setSubmitting(true);
      try {
        const res = await postLogin(values);
        if (res.data.code === 0) {
          message.success(res.data.msg ?? "登录成功");
          await onSuccess();
        } else {
          message.error(res.data.msg ?? "登录失败，请稍后重试");
        }
      } catch (error) {
        console.error(error);
        message.error("登录失败，请稍后重试");
      } finally {
        setSubmitting(false);
      }
    },
    [onSuccess, postLogin]
  );

  return (
    <Modal open={open} onCancel={onCancel} footer={null} centered title="登录">
      <Form
        form={form}
        layout="vertical"
        onFinish={handleFinish}
        className="mt-4"
        disabled={submitting}
      >
        <Form.Item<LoginFormValues>
          name="user"
          label="用户名"
          rules={[
            { required: true, message: "请输入用户名" },
            {
              pattern: /^[\w\u4e00-\u9fa5\uac00-\ud7ff\u0800-\u4e00\-]{2,7}$/,
              message: "2-7位，支持中日韩文字、数字、字母、_、-"
            }
          ]}
        >
          <Input prefix={<UserOutlined />} placeholder="请输入用户名" allowClear />
        </Form.Item>

        <Form.Item<LoginFormValues>
          name="pwd"
          label="密码"
          rules={[
            { required: true, message: "请输入密码" },
            {
              pattern: /^[\w<>,.?|;':"{}!@#$%^&*()/\-\[\]\\]{6,18}$/,
              message: "6-18位，支持常见符号"
            }
          ]}
        >
          <Input.Password prefix={<LockOutlined />} placeholder="请输入密码" />
        </Form.Item>

        <Button type="primary" htmlType="submit" loading={submitting} block>
          登录
        </Button>
      </Form>
      <Typography.Paragraph className="mt-4 text-right text-sm text-slate-500">
        还没有账号？
        <button
          type="button"
          onClick={onRegister}
          className="ml-2 text-brand hover:text-brand-dark"
        >
          立即注册
        </button>
      </Typography.Paragraph>
    </Modal>
  );
};

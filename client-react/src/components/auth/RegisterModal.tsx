import { LockOutlined, SafetyOutlined, UserAddOutlined } from "@ant-design/icons";
import { Button, Form, Input, Modal, Space, message } from "antd";
import type { AxiosResponse } from "axios";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type { BasicResponse, CaptchaPayload } from "@/types/api";

interface RegisterModalProps {
  open: boolean;
  onCancel: () => void;
  onLogin: () => void;
  onSuccess: () => void;
  getRegisterVCode: () => Promise<AxiosResponse<CaptchaPayload>>;
  getRegisterCheckVcode: (svgCode: string) => Promise<AxiosResponse<BasicResponse>>;
  postRegister: (options: Record<string, unknown>) => Promise<AxiosResponse<BasicResponse>>;
}

interface RegisterFormValues {
  user: string;
  pwd: string;
  checkPwd: string;
  svgCode: string;
}

export const RegisterModal = ({
  open,
  onCancel,
  onLogin,
  onSuccess,
  getRegisterVCode,
  getRegisterCheckVcode,
  postRegister
}: RegisterModalProps) => {
  const [form] = Form.useForm<RegisterFormValues>();
  const [captcha, setCaptcha] = useState<string>("");
  const [refreshCountdown, setRefreshCountdown] = useState(0);
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const resetTimer = useCallback(() => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
    }
    timerRef.current = null;
  }, []);

  const fetchCaptcha = useCallback(async () => {
    try {
      const res = await getRegisterVCode();
      setCaptcha(res.data.data);

      const seconds = Math.max(1, Math.ceil((res.data.time ?? 0) / 1000));
      setRefreshCountdown(seconds);

      resetTimer();
      timerRef.current = setInterval(() => {
        setRefreshCountdown((prev) => {
          if (prev <= 1) {
            resetTimer();
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    } catch (error) {
      console.error(error);
      message.error("验证码获取失败，请稍后再试");
    }
  }, [getRegisterVCode, resetTimer]);

  useEffect(() => {
    if (open) {
      void fetchCaptcha();
    }
    return resetTimer;
  }, [fetchCaptcha, open, resetTimer]);

  const refreshText = useMemo(() => {
    if (refreshCountdown <= 0) {
      return "刷新验证码";
    }
    return `${refreshCountdown}s后可刷新`;
  }, [refreshCountdown]);

  const handleSubmit = useCallback(
    async (values: RegisterFormValues) => {
      setSubmitting(true);
      try {
        const res = await postRegister(values);
        if (res.data.code === 0) {
          message.success(res.data.msg ?? "注册成功！");
          form.resetFields();
          onSuccess();
        } else {
          message.error(res.data.msg ?? "注册失败，请稍后再试");
        }
      } catch (error) {
        console.error(error);
        message.error("注册失败，请稍后再试");
      } finally {
        setSubmitting(false);
        void fetchCaptcha();
      }
    },
    [fetchCaptcha, form, onSuccess, postRegister]
  );

  return (
    <Modal open={open} onCancel={onCancel} footer={null} centered title="注册">
      <Form
        layout="vertical"
        form={form}
        onFinish={handleSubmit}
        disabled={submitting}
        className="mt-4"
      >
        <Form.Item<RegisterFormValues>
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
          <Input prefix={<UserAddOutlined />} placeholder="请输入用户名" allowClear />
        </Form.Item>

        <Form.Item<RegisterFormValues>
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

        <Form.Item<RegisterFormValues>
          name="checkPwd"
          label="确认密码"
          dependencies={["pwd"]}
          rules={[
            { required: true, message: "请再次输入密码" },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue("pwd") === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error("两次密码输入不一致"));
              }
            })
          ]}
        >
          <Input.Password prefix={<LockOutlined />} placeholder="请再次输入密码" />
        </Form.Item>

        <Form.Item<RegisterFormValues>
          name="svgCode"
          label="验证码"
          rules={[
            { required: true, message: "请输入验证码" },
            () => ({
              async validator(_, value) {
                if (!value) {
                  return Promise.resolve();
                }
                try {
                  const res = await getRegisterCheckVcode(value);
                  if (res.data.code === 0) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error(res.data.msg ?? "验证码错误"));
                } catch (error) {
                  console.error(error);
                  return Promise.reject(new Error("验证码校验失败"));
                }
              }
            })
          ]}
        >
          <Input prefix={<SafetyOutlined />} placeholder="请输入验证码" allowClear />
        </Form.Item>

        <div className="flex flex-col gap-2">
          <div
            className="min-h-[48px] rounded-md border border-slate-200 bg-white p-2"
            dangerouslySetInnerHTML={{ __html: captcha }}
          />
          <Button type="link" disabled={refreshCountdown > 0} onClick={() => void fetchCaptcha()} className="self-end p-0">
            {refreshText}
          </Button>
        </div>

        <Button
          type="primary"
          htmlType="submit"
          loading={submitting}
          block
          className="mt-6"
        >
          立即注册
        </Button>
      </Form>
      <Space className="mt-4 w-full justify-end text-sm text-slate-500">
        <span>已有账号？</span>
        <button type="button" onClick={onLogin} className="text-brand hover:text-brand-dark">
          去登录
        </button>
      </Space>
    </Modal>
  );
};

import {
  createContext,
  PropsWithChildren,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState
} from "react";
import { message } from "antd";
import {
  postIfLogin,
  postLogin,
  postLogout,
  postRegister,
  getRegisterVCode,
  getRegisterCheckVcode
} from "@/api";
import type { UserProfile } from "@/types/api";
import { LoginModal } from "@/components/auth/LoginModal";
import { RegisterModal } from "@/components/auth/RegisterModal";
import { AvatarModal } from "@/components/auth/AvatarModal";

interface AuthContextValue {
  user: UserProfile | null;
  loading: boolean;
  refreshUser: () => Promise<void>;
  showLogin: () => void;
  showRegister: () => void;
  showAvatar: () => void;
  hideModals: () => void;
  logout: () => Promise<void>;
  requireLogin: () => Promise<UserProfile | null>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used inside AuthProvider");
  }
  return context;
};

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const [user, setUser] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);

  const [loginVisible, setLoginVisible] = useState(false);
  const [registerVisible, setRegisterVisible] = useState(false);
  const [avatarVisible, setAvatarVisible] = useState(false);

  const refreshUser = useCallback(async () => {
    try {
      setLoading(true);
      const res = await postIfLogin();
      const profile = res.data?.userInfo ?? null;
      setUser(profile);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void refreshUser();
  }, [refreshUser]);

  const showLogin = useCallback(() => {
    setRegisterVisible(false);
    setLoginVisible(true);
  }, []);

  const showRegister = useCallback(() => {
    setLoginVisible(false);
    setRegisterVisible(true);
  }, []);

  const showAvatar = useCallback(() => {
    if (!user) {
      void requireLogin();
      return;
    }
    setAvatarVisible(true);
  }, [user]);

  const hideModals = useCallback(() => {
    setLoginVisible(false);
    setRegisterVisible(false);
    setAvatarVisible(false);
  }, []);

  const requireLogin = useCallback(async () => {
    if (user) {
      return user;
    }

    try {
      const res = await postIfLogin();
      const profile = res.data?.userInfo ?? null;
      if (profile) {
        setUser(profile);
        return profile;
      }
    } catch (error) {
      console.error(error);
    }

    message.warning("请先登录后再试试~");
    showLogin();
    return null;
  }, [showLogin, user]);

  const logout = useCallback(async () => {
    try {
      const res = await postLogout();
      if (res.data.code === 0) {
        message.success("退出登录成功");
      } else {
        message.error(res.data.msg ?? "退出失败");
      }
    } catch (error) {
      console.error(error);
      message.error("退出失败，请稍后再试");
    } finally {
      setUser(null);
      await refreshUser();
    }
  }, [refreshUser]);

  const value = useMemo<AuthContextValue>(() => ({
    user,
    loading,
    refreshUser,
    showLogin,
    showRegister,
    showAvatar,
    hideModals,
    logout,
    requireLogin
  }), [hideModals, loading, logout, refreshUser, requireLogin, showAvatar, showLogin, showRegister, user]);

  return (
    <AuthContext.Provider value={value}>
      {children}
      <LoginModal
        open={loginVisible}
        onCancel={hideModals}
        onRegister={showRegister}
        onSuccess={async () => {
          hideModals();
          await refreshUser();
        }}
        postLogin={postLogin}
      />
      <RegisterModal
        open={registerVisible}
        onCancel={hideModals}
        onLogin={showLogin}
        onSuccess={() => {
          hideModals();
          showLogin();
        }}
        getRegisterVCode={getRegisterVCode}
        getRegisterCheckVcode={getRegisterCheckVcode}
        postRegister={postRegister}
      />
      <AvatarModal
        open={avatarVisible}
        onCancel={hideModals}
        onUploaded={async () => {
          hideModals();
          await refreshUser();
        }}
      />
    </AuthContext.Provider>
  );
};

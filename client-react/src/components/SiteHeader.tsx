import {
  DownOutlined,
  LogoutOutlined,
  PictureOutlined,
  UserOutlined
} from "@ant-design/icons";
import { Avatar, Button, Dropdown, Layout, MenuProps, Space, Spin } from "antd";
import { createStyles } from "antd-style";
import { useMemo } from "react";
import { Link, useLocation } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";

const { Header } = Layout;

const useStyles = createStyles(() => ({
  header: {
    backgroundColor: "rgba(255,255,255,0.92)",
    backdropFilter: "blur(12px)",
    borderBottom: "1px solid rgba(148, 163, 184, 0.3)"
  }
}));

const navItems = [
  { key: "home", label: "首页", path: "/" },
  { key: "blog", label: "博客", path: "/blog/0" },
  { key: "message", label: "留言", path: "/message" },
  { key: "diary", label: "日记", path: "/diary" },
  { key: "links", label: "友链", path: "/links" },
  { key: "about", label: "关于", path: "/about" }
] as const;

const getActivePath = (pathname: string) => {
  if (pathname === "/") {
    return "/";
  }
  if (pathname.startsWith("/blog") || pathname.startsWith("/Article")) {
    return "/blog/0";
  }
  if (pathname.startsWith("/message")) {
    return "/message";
  }
  if (pathname.startsWith("/diary")) {
    return "/diary";
  }
  if (pathname.startsWith("/links")) {
    return "/links";
  }
  if (pathname.startsWith("/about")) {
    return "/about";
  }
  return "";
};

export const SiteHeader = () => {
  const { styles } = useStyles();
  const location = useLocation();
  const { user, loading, showLogin, showRegister, showAvatar, logout } = useAuth();

  const activePath = useMemo(() => getActivePath(location.pathname), [location.pathname]);

  const menuItems = useMemo<MenuProps["items"]>(() => {
    if (!user) {
      return [];
    }
    return [
      {
        key: "avatar",
        icon: <PictureOutlined />,
        label: "修改头像"
      },
      {
        type: "divider"
      },
      {
        key: "logout",
        icon: <LogoutOutlined />,
        label: "退出登录"
      }
    ];
  }, [user]);

  return (
    <Header
      className={`${styles.header} sticky top-0 z-50 flex h-16 items-center px-8 shadow-sm`}
    >
      <Link to="/" className="flex items-center gap-3">
        <div className="rounded-full bg-brand px-3 py-1 text-lg font-semibold text-white shadow">
          楚岚
        </div>
        <span className="font-handwriting text-3xl text-slate-700">Jack</span>
      </Link>
      <nav className="ml-10 hidden flex-1 items-center justify-start gap-6 text-base text-slate-600 lg:flex">
        {navItems.map((item) => {
          const isActive = activePath === item.path;
          return (
            <Link
              key={item.key}
              to={item.path}
              className={`relative pb-1 transition-colors duration-200 ${
                isActive ? "text-brand" : "hover:text-brand"
              }`}
            >
              {item.label}
              {isActive && (
                <span className="absolute left-0 right-0 -bottom-1 h-0.5 bg-brand" />
              )}
            </Link>
          );
        })}
      </nav>
      <div className="ml-auto flex items-center gap-4">
        {loading ? (
          <Spin size="small" />
        ) : user ? (
          <Dropdown
            menu={{
              items: menuItems,
              onClick: async ({ key }) => {
                if (key === "avatar") {
                  showAvatar();
                }
                if (key === "logout") {
                  await logout();
                }
              }
            }}
            trigger={["click"]}
          >
            <button
              type="button"
              className="flex items-center gap-2 rounded-full border border-transparent bg-white/80 px-2 py-1 shadow-sm transition hover:border-brand"
            >
              <Avatar size={40} src={user.photo} icon={<UserOutlined />} />
              <Space className="text-left">
                <span className="text-sm text-slate-600">{user.user}</span>
                <DownOutlined className="text-xs text-slate-400" />
              </Space>
            </button>
          </Dropdown>
        ) : (
          <Space>
            <Button type="primary" onClick={showLogin}>
              登录
            </Button>
            <Button onClick={showRegister}>注册</Button>
          </Space>
        )}
      </div>
    </Header>
  );
};

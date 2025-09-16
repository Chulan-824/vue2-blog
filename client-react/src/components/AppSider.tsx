import { Layout, Menu } from "antd";
import {
  HomeOutlined,
  BookOutlined,
  UserOutlined,
  SettingOutlined
} from "@ant-design/icons";

const { Sider } = Layout;

const menuItems = [
  {
    key: "dashboard",
    icon: <HomeOutlined />,
    label: "仪表盘"
  },
  {
    key: "posts",
    icon: <BookOutlined />,
    label: "文章管理"
  },
  {
    key: "profile",
    icon: <UserOutlined />,
    label: "个人中心"
  },
  {
    key: "settings",
    icon: <SettingOutlined />,
    label: "系统设置"
  }
];

export const AppSider = () => {
  return (
    <Sider
      width={220}
      breakpoint="lg"
      collapsedWidth={64}
      className="bg-white/80 backdrop-blur"
    >
      <Menu
        mode="inline"
        defaultSelectedKeys={["dashboard"]}
        items={menuItems}
        className="border-none text-base"
      />
    </Sider>
  );
};

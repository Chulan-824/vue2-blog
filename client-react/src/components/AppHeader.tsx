import { Layout, Menu, Typography } from "antd";
import { FireOutlined } from "@ant-design/icons";

const { Header } = Layout;

const items = [
  {
    key: "home",
    label: "首页"
  },
  {
    key: "articles",
    label: "文章"
  },
  {
    key: "about",
    label: "关于"
  }
];

export const AppHeader = () => {
  return (
    <Header className="sticky top-0 z-10 flex items-center justify-between bg-white/90 px-6 shadow-sm backdrop-blur">
      <div className="flex items-center gap-2">
        <FireOutlined className="text-2xl text-brand" />
        <Typography.Title level={4} className="!mb-0">
          Vue2 Blog
        </Typography.Title>
      </div>
      <Menu
        className="flex-1 justify-end border-none text-base"
        mode="horizontal"
        defaultSelectedKeys={["home"]}
        items={items}
      />
    </Header>
  );
};

import { Layout } from "antd";
import { Outlet } from "react-router-dom";
import { SiteHeader } from "@/components/SiteHeader";

const { Content } = Layout;

export const PageLayout = () => {
  return (
    <Layout className="min-h-screen bg-transparent">
      <SiteHeader />
      <Content className="mx-auto w-full max-w-6xl px-4 py-8 lg:px-8">
        <div className="rounded-2xl bg-white/85 p-6 shadow-xl backdrop-blur-sm">
          <Outlet />
        </div>
      </Content>
    </Layout>
  );
};

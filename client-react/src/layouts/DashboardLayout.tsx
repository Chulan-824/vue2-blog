import { Layout } from "antd";
import { PropsWithChildren } from "react";
import { AppHeader } from "../components/AppHeader";
import { AppSider } from "../components/AppSider";

const { Content, Footer } = Layout;

export const DashboardLayout = ({ children }: PropsWithChildren) => {
  return (
    <Layout className="min-h-screen">
      <AppHeader />
      <Layout className="bg-transparent">
        <AppSider />
        <Content className="mx-6 my-4 rounded-xl bg-white p-6 shadow-sm">
          {children}
        </Content>
      </Layout>
      <Footer className="text-center text-xs text-slate-500">
        Vue2 Blog React Client &copy; {new Date().getFullYear()} Created with Vite 7 + React 18
      </Footer>
    </Layout>
  );
};

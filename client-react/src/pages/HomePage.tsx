import { Card, Col, Row, Space, Statistic, Typography } from "antd";
import { CalendarOutlined, LikeOutlined, MessageOutlined } from "@ant-design/icons";

const stats = [
  {
    title: "本月阅读",
    value: 12450,
    suffix: "次",
    icon: <LikeOutlined className="text-lg text-brand" />
  },
  {
    title: "新增评论",
    value: 318,
    suffix: "条",
    icon: <MessageOutlined className="text-lg text-brand" />
  },
  {
    title: "发布文章",
    value: 12,
    suffix: "篇",
    icon: <CalendarOutlined className="text-lg text-brand" />
  }
];

const timeline = [
  {
    title: "深入理解React 18并发特性",
    excerpt: "并发渲染为复杂交互带来了更多的性能空间...",
    date: "2024-05-18"
  },
  {
    title: "Vite 7项目最佳实践",
    excerpt: "在企业级项目中使用Vite 7需要关注哪些细节?",
    date: "2024-05-10"
  },
  {
    title: "Tailwind CSS v4探索",
    excerpt: "一文带你了解Tailwind CSS v4的设计理念与落地方式...",
    date: "2024-04-26"
  }
];

export const HomePage = () => {
  return (
    <Space direction="vertical" size="large" className="w-full">
      <div>
        <Typography.Title level={3} className="!mb-1">
          控制台概览
        </Typography.Title>
        <Typography.Paragraph type="secondary" className="!mb-0">
          快速查看站点运行情况，掌握内容创作节奏。
        </Typography.Paragraph>
      </div>
      <Row gutter={[16, 16]}>
        {stats.map((item) => (
          <Col xs={24} md={8} key={item.title}>
            <Card className="h-full">
              <Statistic
                title={
                  <div className="flex items-center gap-2">
                    {item.icon}
                    <span>{item.title}</span>
                  </div>
                }
                value={item.value}
                suffix={item.suffix}
              />
            </Card>
          </Col>
        ))}
      </Row>
      <Card title="最新动态">
        <Space direction="vertical" size={16} className="w-full">
          {timeline.map((item) => (
            <div key={item.title} className="rounded-lg border border-slate-100 p-4">
              <Typography.Title level={5}>{item.title}</Typography.Title>
              <Typography.Paragraph type="secondary">{item.excerpt}</Typography.Paragraph>
              <Typography.Text className="text-xs text-slate-500">
                {item.date}
              </Typography.Text>
            </div>
          ))}
        </Space>
      </Card>
    </Space>
  );
};

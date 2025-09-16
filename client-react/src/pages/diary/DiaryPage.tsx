import { Card, Timeline } from "antd";

const diaries = [
  {
    date: "2020/1/5",
    txt: "这是我写的第1篇日志。",
    img: ["http://localhost:3000/img/diary/11.jpg", "http://localhost:3000/img/diary/22.jpg"]
  },
  {
    date: "2020/1/5",
    txt: "这是我写的第2篇日志。",
    img: []
  },
  {
    date: "2020/1/5",
    txt: "这是我写的第3篇日志。",
    img: []
  },
  {
    date: "2020/1/5",
    txt: "这是我写的第4篇日志。",
    img: []
  },
  {
    date: "2020/1/5",
    txt: "这是我写的第5篇日志。",
    img: ["http://localhost:3000/img/diary/bg.jpg"]
  },
  {
    date: "2020/1/5",
    txt: "这是我写的第6篇日志。",
    img: []
  }
];

export const DiaryPage = () => {
  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-white p-6 shadow-lg">
        <Timeline mode="left">
          {diaries.map((item) => (
            <Timeline.Item label={item.date} key={`${item.date}-${item.txt}`}>
              <Card bordered={false} className="bg-slate-50">
                <p className="text-sm leading-7 text-slate-600">{item.txt}</p>
                <div className="mt-3 grid gap-3">
                  {item.img.map((src) => (
                    <img src={src} alt="diary" key={src} className="w-full rounded-lg object-cover" />
                  ))}
                </div>
              </Card>
            </Timeline.Item>
          ))}
        </Timeline>
      </div>
    </div>
  );
};

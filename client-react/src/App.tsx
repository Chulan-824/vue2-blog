import { BrowserRouter, Route, Routes } from "react-router-dom";
import { AuthProvider } from "@/contexts/AuthContext";
import { PageLayout } from "@/layouts/PageLayout";
import { HomePage } from "@/pages/home/HomePage";
import { BlogPage } from "@/pages/blog/BlogPage";
import { ArticlePage } from "@/pages/article/ArticlePage";
import { MessagePage } from "@/pages/message/MessagePage";
import { DiaryPage } from "@/pages/diary/DiaryPage";
import { LinksPage } from "@/pages/links/LinksPage";
import { AboutPage } from "@/pages/about/AboutPage";
import { NotFoundPage } from "@/pages/not-found/NotFoundPage";

const App = () => {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route element={<PageLayout />}>
            <Route path="blog">
              <Route index element={<BlogPage />} />
              <Route path=":id" element={<BlogPage />} />
            </Route>
            <Route path="Article/:id" element={<ArticlePage />} />
            <Route path="message" element={<MessagePage />} />
            <Route path="diary" element={<DiaryPage />} />
            <Route path="links" element={<LinksPage />} />
            <Route path="about" element={<AboutPage />} />
            <Route path="404" element={<NotFoundPage />} />
            <Route path="*" element={<NotFoundPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
};

export default App;

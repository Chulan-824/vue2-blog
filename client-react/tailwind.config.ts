import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx,jsx,js}"],
  theme: {
    extend: {
      colors: {
        brand: {
          DEFAULT: "#1677ff",
          dark: "#0958d9"
        }
      }
    }
  },
  plugins: []
} satisfies Config;

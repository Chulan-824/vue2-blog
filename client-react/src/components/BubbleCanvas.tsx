import { useEffect, useRef } from "react";

class Bubble {
  x: number;
  y: number;
  vx: number;
  vy: number;
  r: number;
  dead: boolean;
  color: string;

  constructor(x: number, y: number) {
    this.x = x;
    this.y = y;
    this.vx = Math.random() * 3 - 1.5;
    this.vy = Math.random() * 3 - 1.5;
    this.r = Math.random() * 3 + 3;
    this.dead = false;
    this.color = `#${Math.random().toString(16).slice(2, 8)}`;
  }

  render(ctx: CanvasRenderingContext2D) {
    ctx.beginPath();
    ctx.fillStyle = this.color;
    ctx.globalCompositeOperation = "lighter";
    ctx.arc(this.x, this.y, this.r, 0, Math.PI * 2);
    ctx.fill();
  }

  update() {
    this.r *= 0.96;
    this.x += this.vx;
    this.y += this.vy;
    if (this.r < 0.01) {
      this.dead = true;
    }
  }
}

interface BubbleCanvasProps {
  height?: number;
  className?: string;
}

export const BubbleCanvas = ({ height = 260, className }: BubbleCanvasProps) => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const bubblesRef = useRef<Bubble[]>([]);
  const animationRef = useRef<number>();
  const timerRef = useRef<number>();

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) {
      return;
    }
    const ctx = canvas.getContext("2d");
    if (!ctx) {
      return;
    }

    const resize = () => {
      canvas.width = window.innerWidth;
      canvas.height = height;
    };

    const addBubble = (x: number, y: number) => {
      bubblesRef.current.push(new Bubble(x, y));
    };

    const autoGenerate = () => {
      timerRef.current = window.setInterval(() => {
        addBubble(Math.random() * canvas.width, Math.random() * canvas.height);
      }, 30);
    };

    const render = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      bubblesRef.current.forEach((bubble) => {
        bubble.render(ctx);
        bubble.update();
      });
      bubblesRef.current = bubblesRef.current.filter((bubble) => !bubble.dead);
      animationRef.current = requestAnimationFrame(render);
    };

    resize();
    autoGenerate();
    render();

    const handleResize = () => {
      resize();
    };

    const handleMousemove = (event: MouseEvent) => {
      const rect = canvas.getBoundingClientRect();
      addBubble(event.clientX - rect.left, event.clientY - rect.top);
    };

    window.addEventListener("resize", handleResize);
    canvas.addEventListener("mousemove", handleMousemove);

    return () => {
      window.removeEventListener("resize", handleResize);
      canvas.removeEventListener("mousemove", handleMousemove);
      if (timerRef.current) {
        clearInterval(timerRef.current);
      }
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [height]);

  return <canvas ref={canvasRef} className={className} height={height} />;
};

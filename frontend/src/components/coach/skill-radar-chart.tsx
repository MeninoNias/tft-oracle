import type { SkillRadar } from "@/gen/tft/v1/coach_pb";

interface SkillRadarChartProps {
  radar: SkillRadar;
  size?: number;
}

const LABELS = ["economy", "itemization", "composition", "adaptability", "consistency"];
const N = 5;
const ANGLE_OFFSET = -Math.PI / 2;

function getAngle(i: number): number {
  return ANGLE_OFFSET + (2 * Math.PI * i) / N;
}

function polarToCartesian(cx: number, cy: number, r: number, angle: number) {
  return {
    x: cx + r * Math.cos(angle),
    y: cy + r * Math.sin(angle),
  };
}

function clamp(v: number): number {
  return Math.min(Math.max(v, 0), 100);
}

export function SkillRadarChart({ radar, size = 200 }: SkillRadarChartProps) {
  const cx = size / 2;
  const cy = size / 2;
  const maxR = size * 0.38;
  const labelR = size * 0.48;

  const values: number[] = [
    radar.economy,
    radar.itemization,
    radar.composition,
    radar.adaptability,
    radar.consistency,
  ];

  const items = values.map((v, i) => ({
    value: v,
    angle: getAngle(i),
    label: LABELS[i] ?? "",
  }));

  const rings = [0.25, 0.5, 0.75, 1.0];

  // Data polygon points
  const dataPoints = items
    .map((item) => {
      const r = (clamp(item.value) / 100) * maxR;
      const p = polarToCartesian(cx, cy, r, item.angle);
      return `${p.x},${p.y}`;
    })
    .join(" ");

  return (
    <svg
      width={size}
      height={size}
      viewBox={`0 0 ${size} ${size}`}
      className="block"
    >
      {/* Grid rings */}
      {rings.map((scale) => {
        const pts = items
          .map((item) => {
            const p = polarToCartesian(cx, cy, maxR * scale, item.angle);
            return `${p.x},${p.y}`;
          })
          .join(" ");
        return (
          <polygon
            key={scale}
            points={pts}
            fill="none"
            stroke="currentColor"
            strokeWidth={0.5}
            className="text-lofi-border"
          />
        );
      })}

      {/* Axis lines */}
      {items.map((item, i) => {
        const p = polarToCartesian(cx, cy, maxR, item.angle);
        return (
          <line
            key={i}
            x1={cx}
            y1={cy}
            x2={p.x}
            y2={p.y}
            stroke="currentColor"
            strokeWidth={0.5}
            className="text-lofi-border"
          />
        );
      })}

      {/* Data polygon */}
      <polygon
        points={dataPoints}
        fill="currentColor"
        fillOpacity={0.15}
        stroke="currentColor"
        strokeWidth={1.5}
        className="text-lofi-accent"
      />

      {/* Data points */}
      {items.map((item, i) => {
        const r = (clamp(item.value) / 100) * maxR;
        const p = polarToCartesian(cx, cy, r, item.angle);
        return (
          <circle
            key={i}
            cx={p.x}
            cy={p.y}
            r={2.5}
            fill="currentColor"
            className="text-lofi-accent"
          />
        );
      })}

      {/* Labels */}
      {items.map((item) => {
        const p = polarToCartesian(cx, cy, labelR, item.angle);
        return (
          <text
            key={item.label}
            x={p.x}
            y={p.y}
            textAnchor="middle"
            dominantBaseline="middle"
            fill="currentColor"
            className="text-lofi-muted"
            fontSize={10}
          >
            {item.label}
          </text>
        );
      })}

      {/* Value labels */}
      {items.map((item, i) => {
        const r = (clamp(item.value) / 100) * maxR + 10;
        const p = polarToCartesian(cx, cy, r, item.angle);
        return (
          <text
            key={`val-${i}`}
            x={p.x}
            y={p.y}
            textAnchor="middle"
            dominantBaseline="middle"
            fill="currentColor"
            className="text-lofi-secondary"
            fontSize={9}
            fontWeight="bold"
          >
            {Math.round(item.value)}
          </text>
        );
      })}
    </svg>
  );
}

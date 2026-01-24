export default function GoldIcon({ className = '', width = 18, height = 18 }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      width={width}
      height={height}
      className={className}
      aria-hidden
    >
      <circle cx="12" cy="12" r="10" fill="#FFD700" stroke="#B8860B" strokeWidth="2" />
      <circle cx="12" cy="12" r="6" fill="none" stroke="#D4A500" strokeWidth="1" opacity="0.5" />
    </svg>
  )
}

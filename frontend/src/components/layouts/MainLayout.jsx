import Framing from '../framing/Framing'
import Stats from '../player/Stats'
import Gold from '../player/Gold'

export default function MainLayout({ children }) {
  const messages = []

  return (
    <div className="main-layout">
      <div className="header">
        <Gold />
        <Stats />
      </div>
      <div className="content">
        {children}
      </div>
      <Framing messages={messages} />
    </div>
  )
}

















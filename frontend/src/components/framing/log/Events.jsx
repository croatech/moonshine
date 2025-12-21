export default function Events({ events = [] }) {
  return (
    <div className="events">
      {events.map((event, index) => (
        <div key={index} className="event">
          <div dangerouslySetInnerHTML={{ __html: event.body }} />
        </div>
      ))}
    </div>
  )
}


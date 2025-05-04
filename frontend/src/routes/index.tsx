import { listInstruments } from "@/gen/client/expenses/v1/instrument-InstrumentService_connectquery"
import { useQuery } from "@connectrpc/connect-query"
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/")({
  component: Index,
})

function Index() {
  const { data } = useQuery(listInstruments)
  console.log("data", data)
  return (
    <div className="p-2">
      <h3>Welcome Home!</h3>
    </div>
  )
}

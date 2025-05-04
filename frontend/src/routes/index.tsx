import { useQuery } from "@connectrpc/connect-query"
import { createFileRoute } from "@tanstack/react-router"
import { listUsers } from "api/client/expenses/v1/user-UserService_connectquery"

export const Route = createFileRoute("/")({
  component: Index,
})

function Index() {
  const { data } = useQuery(listUsers)
  console.log("data", data)
  return (
    <div className="p-2">
      <h3>Welcome Home!</h3>
    </div>
  )
}

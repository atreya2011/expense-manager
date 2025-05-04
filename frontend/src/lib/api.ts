import { createPromiseClient } from "@connectrpc/connect"
import { createConnectTransport } from "@connectrpc/connect-web"
import { UserService } from "../../gen/client/expenses/v1/user_connectweb"
import { InstrumentService } from "../../gen/client/expenses/v1/instrument_connectweb"

const transport = createConnectTransport({
  baseUrl: "http://localhost:8080", // TODO: Make this configurable
})

export const userService = createPromiseClient(UserService, transport)
export const instrumentService = createPromiseClient(InstrumentService, transport)

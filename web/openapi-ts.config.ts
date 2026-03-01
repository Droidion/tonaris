import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: "../shared/openapi.json",
  output: "src/client",
  plugins: ["@tanstack/react-query"],
});

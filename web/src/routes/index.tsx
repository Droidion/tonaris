import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";

import { tonarisHelloOptions } from "../client/@tanstack/react-query.gen";

export const Route = createFileRoute("/")({
  component: Home,
});

function Home() {
  const { data, isLoading, error } = useQuery(tonarisHelloOptions());

  if (isLoading) return <p>Loading...</p>;
  if (error) return <p>Error: {String(error)}</p>;
  return <p>{data}</p>;
}

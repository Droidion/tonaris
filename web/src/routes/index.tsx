import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";

import { tonarisHelloOptions } from "../client/@tanstack/react-query.gen";

export const Route = createFileRoute("/")({
  component: Home,
});

function Home() {
  const { data, isLoading, error } = useQuery(tonarisHelloOptions());

  if (isLoading)
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Spinner className="size-8" />
      </div>
    );
  if (error) return <p>Error: {String(error)}</p>;

  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-4">
      <h1 className="scroll-m-20 text-center text-4xl font-extrabold">Tonaris</h1>
      <p>{data}</p>
      <Button variant="outline">Shadcn</Button>
    </main>
  );
}

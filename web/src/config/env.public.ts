import { z } from "zod";

const clientEnvSchema = z.object({
  VITE_API_URL: z.url(),
  VITE_CLERK_PUBLISHABLE_KEY: z.string().min(1),
});

export const clientEnv = clientEnvSchema.parse(import.meta.env);

import { z } from "zod";

const envSchema = z.object({
  API_URL: z.url(),
});

const clientEnvSchema = z.object({
  VITE_CLERK_PUBLISHABLE_KEY: z.string().min(1),
});

// Validate server environment
export const serverEnv = envSchema.parse(process.env);

// Validate client environment
export const clientEnv = clientEnvSchema.parse(import.meta.env);

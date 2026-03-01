import { z } from "zod";

const envSchema = z.object({
  DB_URL: z.url(),
});

export const serverEnv = envSchema.parse(process.env);

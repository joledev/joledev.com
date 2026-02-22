import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';

const blog = defineCollection({
  loader: glob({ pattern: '**/*.{md,mdx}', base: './src/content/blog' }),
  schema: z.object({
    title: z.string(),
    description: z.string(),
    pubDate: z.coerce.date(),
    updatedDate: z.coerce.date().optional(),
    heroImage: z.string().optional(),
    tags: z.array(z.string()),
    category: z.enum(['tutorial', 'case-study', 'opinion', 'project']),
    lang: z.enum(['es', 'en']),
    draft: z.boolean().default(false),
  }),
});

const projects = defineCollection({
  loader: glob({ pattern: '**/*.{md,mdx}', base: './src/content/projects' }),
  schema: z.object({
    title: z.string(),
    description: z.string(),
    heroImage: z.string(),
    screenshots: z.array(z.string()).optional(),
    tags: z.array(z.string()),
    category: z.enum(['web', 'system', 'integration', 'mobile', 'game', 'devops']),
    lang: z.enum(['es', 'en']),
    featured: z.boolean(),
    order: z.number(),
  }),
});

export const collections = { blog, projects };

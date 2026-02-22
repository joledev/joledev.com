import type { APIContext, GetStaticPaths } from 'astro';
import { getCollection } from 'astro:content';
import satori from 'satori';
import sharp from 'sharp';

const FONT_URL = 'https://cdn.jsdelivr.net/fontsource/fonts/dm-sans@latest/latin-700-normal.woff';

let fontData: ArrayBuffer | null = null;

async function loadFont(): Promise<ArrayBuffer> {
  if (fontData) return fontData;
  const res = await fetch(FONT_URL);
  fontData = await res.arrayBuffer();
  return fontData;
}

export const getStaticPaths: GetStaticPaths = async () => {
  const posts = await getCollection('blog');
  return posts
    .filter((p) => !p.data.draft)
    .map((post) => ({
      params: { slug: `${post.data.lang}/blog/${post.id.replace(/^(es|en)\//, '')}` },
      props: { post },
    }));
};

export async function GET({ props }: APIContext) {
  const { post } = props as { post: Awaited<ReturnType<typeof getCollection>>[number] };
  const { title, tags, pubDate } = post.data;

  const font = await loadFont();

  const dateStr = new Date(pubDate).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });

  const svg = await satori(
    {
      type: 'div',
      props: {
        style: {
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'space-between',
          padding: '60px',
          background: 'linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%)',
          fontFamily: 'DM Sans',
        },
        children: [
          {
            type: 'div',
            props: {
              style: {
                display: 'flex',
                alignItems: 'center',
                gap: '12px',
              },
              children: [
                {
                  type: 'div',
                  props: {
                    style: {
                      width: '40px',
                      height: '40px',
                      borderRadius: '8px',
                      background: '#2563eb',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      color: 'white',
                      fontSize: '18px',
                      fontWeight: 700,
                    },
                    children: 'J',
                  },
                },
                {
                  type: 'span',
                  props: {
                    style: { color: '#60a5fa', fontSize: '24px', fontWeight: 700 },
                    children: 'JoleDev',
                  },
                },
              ],
            },
          },
          {
            type: 'div',
            props: {
              style: {
                display: 'flex',
                flexDirection: 'column',
                gap: '16px',
                flex: 1,
                justifyContent: 'center',
              },
              children: [
                {
                  type: 'h1',
                  props: {
                    style: {
                      fontSize: title.length > 50 ? '40px' : '48px',
                      fontWeight: 700,
                      color: '#f1f5f9',
                      lineHeight: 1.2,
                      margin: 0,
                    },
                    children: title,
                  },
                },
              ],
            },
          },
          {
            type: 'div',
            props: {
              style: {
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
              },
              children: [
                {
                  type: 'div',
                  props: {
                    style: { display: 'flex', gap: '8px' },
                    children: (tags as string[]).slice(0, 4).map((tag: string) => ({
                      type: 'span',
                      props: {
                        style: {
                          padding: '4px 12px',
                          borderRadius: '9999px',
                          background: 'rgba(59, 130, 246, 0.2)',
                          color: '#60a5fa',
                          fontSize: '14px',
                        },
                        children: tag,
                      },
                    })),
                  },
                },
                {
                  type: 'span',
                  props: {
                    style: { color: '#94a3b8', fontSize: '16px' },
                    children: dateStr,
                  },
                },
              ],
            },
          },
        ],
      },
    },
    {
      width: 1200,
      height: 630,
      fonts: [
        {
          name: 'DM Sans',
          data: font,
          weight: 700,
          style: 'normal',
        },
      ],
    },
  );

  const png = await sharp(Buffer.from(svg)).png().toBuffer();

  return new Response(png, {
    headers: { 'Content-Type': 'image/png' },
  });
}

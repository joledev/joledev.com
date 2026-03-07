import type { APIContext } from 'astro';
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

export async function GET(_ctx: APIContext) {
  const font = await loadFont();

  const svg = await satori(
    {
      type: 'div',
      props: {
        style: {
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          padding: '60px',
          background: 'linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%)',
          fontFamily: 'DM Sans',
          gap: '24px',
        },
        children: [
          {
            type: 'div',
            props: {
              style: {
                display: 'flex',
                alignItems: 'center',
                gap: '16px',
              },
              children: [
                {
                  type: 'div',
                  props: {
                    style: {
                      width: '64px',
                      height: '64px',
                      borderRadius: '16px',
                      background: '#2563eb',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      color: 'white',
                      fontSize: '32px',
                      fontWeight: 700,
                    },
                    children: 'J',
                  },
                },
                {
                  type: 'span',
                  props: {
                    style: { color: '#f1f5f9', fontSize: '56px', fontWeight: 700 },
                    children: 'JoleDev',
                  },
                },
              ],
            },
          },
          {
            type: 'p',
            props: {
              style: {
                color: '#94a3b8',
                fontSize: '28px',
                textAlign: 'center',
                maxWidth: '800px',
                lineHeight: 1.4,
              },
              children: 'Desarrollo tecnologico a la medida de tu negocio',
            },
          },
          {
            type: 'div',
            props: {
              style: {
                display: 'flex',
                gap: '12px',
                marginTop: '16px',
              },
              children: ['Web', 'Mobile', 'APIs', 'Cloud', 'AI'].map((tag) => ({
                type: 'span',
                props: {
                  style: {
                    padding: '8px 20px',
                    borderRadius: '9999px',
                    background: 'rgba(59, 130, 246, 0.2)',
                    color: '#60a5fa',
                    fontSize: '18px',
                  },
                  children: tag,
                },
              })),
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

import es from './es.json';
import en from './en.json';

type Lang = 'es' | 'en';

const translations = { es, en } as const;

export function getLangFromUrl(url: URL): Lang {
  const [, lang] = url.pathname.split('/');
  if (lang === 'en') return 'en';
  return 'es';
}

export function useTranslations(lang: Lang) {
  return translations[lang];
}

export function getLocalizedPath(path: string, lang: Lang): string {
  // Remove any existing lang prefix
  const cleanPath = path.replace(/^\/(es|en)/, '');
  return `/${lang}${cleanPath || '/'}`;
}

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

// Route slug mapping between languages
const routeMap: Record<string, string> = {
  'agendar': 'schedule',
  'schedule': 'agendar',
  'cotizar': 'quote',
  'quote': 'cotizar',
  'proyectos': 'projects',
  'projects': 'proyectos',
};

export function getLocalizedPath(path: string, lang: Lang): string {
  // Remove existing lang prefix
  const cleanPath = path.replace(/^\/(es|en)/, '');

  // Translate route slugs
  const segments = cleanPath.split('/').filter(Boolean);
  if (segments.length > 0 && routeMap[segments[0]]) {
    segments[0] = routeMap[segments[0]];
  }

  const translated = segments.length > 0 ? `/${segments.join('/')}` : '';
  return `/${lang}${translated || '/'}`;
}

<script lang="ts">
  import { onMount } from 'svelte';
  import { animate, stagger } from 'animejs';

  interface Project {
    id: string;
    data: {
      title: string;
      description: string;
      heroImage: string;
      tags: string[];
      category: string;
      lang: string;
      featured: boolean;
      order: number;
    };
  }

  interface Props {
    projects: Project[];
    lang: 'es' | 'en';
    categories: { value: string; label: string }[];
    basePath: string;
  }

  let { projects, lang, categories, basePath }: Props = $props();

  let activeCategory = $state('all');
  let gridEl: HTMLElement;

  let filteredProjects = $derived(
    activeCategory === 'all'
      ? projects
      : projects.filter((p) => p.data.category === activeCategory),
  );

  const categoryLabels: Record<string, Record<string, string>> = {
    web: { es: 'Web', en: 'Web' },
    system: { es: 'Sistema', en: 'System' },
    integration: { es: 'Integración', en: 'Integration' },
    mobile: { es: 'Móvil', en: 'Mobile' },
    game: { es: 'Videojuego', en: 'Game' },
    devops: { es: 'DevOps', en: 'DevOps' },
  };

  // Category icons (SVG paths)
  const categoryIcons: Record<string, string> = {
    web: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z',
    mobile: 'M17 1.01L7 1c-1.1 0-2 .9-2 2v18c0 1.1.9 2 2 2h10c1.1 0 2-.9 2-2V3c0-1.1-.9-1.99-2-1.99zM17 19H7V5h10v14z',
    system: 'M20 18c1.1 0 1.99-.9 1.99-2L22 6c0-1.1-.9-2-2-2H4c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2H0v2h24v-2h-4zM4 6h16v10H4V6z',
    integration: 'M13.127 14.56l1.43-1.43 6.44 6.443L19.57 21l-6.44-6.44zm-2.254-2.253l-3.291 3.291a1 1 0 01-1.414 0l-3.536-3.535a1 1 0 010-1.415l3.291-3.29-1.06-1.06 1.413-1.415 9.193 9.193-1.414 1.414-1.06-1.061zm-.707-5.656l3.536 3.536a1 1 0 010 1.414l-.354.354 1.06 1.06 1.415-1.413a1 1 0 011.414 0l3.536 3.535a1 1 0 010 1.415l-3.536 3.535a1 1 0 01-1.414 0l-3.536-3.535a1 1 0 010-1.414l.354-.354-1.06-1.061-1.415 1.414a1 1 0 01-1.414 0L2.929 8.515a1 1 0 010-1.414l3.536-3.536a1 1 0 011.414 0z',
    game: 'M21 6H3c-1.1 0-2 .9-2 2v8c0 1.1.9 2 2 2h18c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm-10 7H8v3H6v-3H3v-2h3V8h2v3h3v2zm4.5 2c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5zm4-3c-.83 0-1.5-.67-1.5-1.5S18.67 9 19.5 9s1.5.67 1.5 1.5-.67 1.5-1.5 1.5z',
    devops: 'M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96zM10 17l-3.5-3.5 1.41-1.41L10 14.17l4.59-4.58L16 11l-6 6z',
  };

  function getCategoryLabel(cat: string): string {
    return categoryLabels[cat]?.[lang] || cat;
  }

  // Tech icon map — returns an SVG viewBox path for known technologies
  const TECH_ICONS: Record<string, { icon: string; color: string }> = {
    'React': { icon: 'M12 10.11c1.03 0 1.87.84 1.87 1.89 0 1-.84 1.85-1.87 1.85S10.13 13 10.13 12c0-1.05.84-1.89 1.87-1.89M7.37 20c.63.38 2.01-.2 3.6-1.7-.52-.59-1.03-1.23-1.51-1.9a22.7 22.7 0 01-2.4-.36c-.51 2.14-.32 3.61.31 3.96m.71-5.74l-.29-.51c-.11.29-.22.58-.29.86.27.06.57.11.88.16l-.3-.51m6.54-.76l.81-1.5-.81-1.5c-.3-.53-.62-1-.91-1.47C13.17 9 12.6 9 12 9c-.6 0-1.17 0-1.71.03-.29.47-.61.94-.91 1.47L8.57 12l.81 1.5c.3.53.62 1 .91 1.47.54.03 1.11.03 1.71.03.6 0 1.17 0 1.71-.03.29-.47.61-.94.91-1.47M12 6.78c-.19.22-.39.45-.59.72h1.18c-.2-.27-.4-.5-.59-.72m0 10.44c.19-.22.39-.45.59-.72h-1.18c.2.27.4.5.59.72M16.62 4c-.62-.38-2 .2-3.59 1.7.52.59 1.03 1.23 1.51 1.9.82.08 1.63.2 2.4.36.51-2.14.32-3.61-.32-3.96m-.7 5.74l.29.51c.11-.29.22-.58.29-.86-.27-.06-.57-.11-.88-.16l.3.51m1.45-7.05c1.47.84 1.63 3.05 1.01 5.63 2.54.75 4.37 1.99 4.37 3.68 0 1.69-1.83 2.93-4.37 3.68.62 2.58.46 4.79-1.01 5.63-1.46.84-3.45-.12-5.37-1.95-1.92 1.83-3.91 2.79-5.38 1.95-1.46-.84-1.62-3.05-1-5.63-2.54-.75-4.37-1.99-4.37-3.68 0-1.69 1.83-2.93 4.37-3.68-.62-2.58-.46-4.79 1-5.63 1.47-.84 3.46.12 5.38 1.95 1.92-1.83 3.91-2.79 5.37-1.95M17.08 12c.34.75.64 1.5.89 2.26 2.1-.63 3.28-1.53 3.28-2.26 0-.73-1.18-1.63-3.28-2.26-.25.76-.55 1.51-.89 2.26M6.92 12c-.34-.75-.64-1.5-.89-2.26-2.1.63-3.28 1.53-3.28 2.26 0 .73 1.18 1.63 3.28 2.26.25-.76.55-1.51.89-2.26m9 2.26l-.3.51c.31-.05.61-.1.88-.16-.07-.28-.18-.57-.29-.86l-.29.51m-9.82 1.12c.71 2.09 2.06 2.82 2.86 2.39.4-.22.73-.73.89-1.51a14.7 14.7 0 01-.89-.86l-.51-.89c-.82.08-1.63.2-2.35.36v.51z', color: '#61DAFB' },
    'React Native': { icon: 'M12 10.11c1.03 0 1.87.84 1.87 1.89 0 1-.84 1.85-1.87 1.85S10.13 13 10.13 12c0-1.05.84-1.89 1.87-1.89M7.37 20c.63.38 2.01-.2 3.6-1.7-.52-.59-1.03-1.23-1.51-1.9a22.7 22.7 0 01-2.4-.36c-.51 2.14-.32 3.61.31 3.96', color: '#61DAFB' },
    'TypeScript': { icon: 'M3 3h18v18H3V3zm10.71 13.1c.06.88.69 1.39 1.71 1.39.91 0 1.55-.44 1.55-1.15 0-.63-.38-.97-1.32-1.22l-.56-.15c-1.32-.36-1.98-.96-1.98-1.97 0-1.2.91-2 2.28-2 1.29 0 2.18.76 2.22 1.86h-1.08c-.06-.74-.58-1.15-1.17-1.15-.65 0-1.12.38-1.12.97 0 .56.38.88 1.22 1.1l.53.15c1.41.38 2.1.97 2.1 2.03 0 1.26-.97 2.09-2.44 2.09-1.41 0-2.38-.76-2.44-1.95h1.13zM13 11.25H7.97v1.08h1.94V19h1.13v-6.67H13v-1.08z', color: '#3178C6' },
    'Node.js': { icon: 'M12 1.85c-.27 0-.55.07-.78.2l-7.44 4.3c-.48.28-.78.8-.78 1.36v8.58c0 .56.3 1.08.78 1.36l1.95 1.12c.95.46 1.27.46 1.71.46 1.4 0 2.21-.85 2.21-2.33V8.44c0-.12-.1-.22-.22-.22H8.5c-.13 0-.23.1-.23.22v8.47c0 .66-.68 1.31-1.77.76L4.45 16.5a.26.26 0 01-.12-.21V7.71c0-.09.05-.17.12-.21l7.44-4.29c.07-.04.16-.04.23 0l7.43 4.29c.08.04.13.12.13.21v8.58c0 .08-.05.17-.13.21l-7.43 4.29c-.07.04-.15.04-.23 0l-1.88-1.12c-.06-.03-.13-.04-.2-.01-.54.24-.65.27-1.16.41-.12.04-.31.1.07.28l2.48 1.47c.24.13.5.2.78.2s.54-.07.78-.2l7.44-4.29c.48-.28.78-.8.78-1.36V7.71c0-.56-.3-1.08-.78-1.36l-7.44-4.3c-.23-.13-.5-.2-.78-.2', color: '#539E43' },
    'Go': { icon: 'M3.33 9.63c-.06 0-.1-.04-.07-.09l.49-.64c.03-.05.1-.08.16-.08h8.44c.06 0 .08.05.06.09l-.39.58c-.03.05-.1.09-.16.09l-8.53.05zm-2.01 1.66c-.06 0-.1-.04-.07-.09l.49-.64c.03-.05.1-.08.16-.08h10.79c.06 0 .09.05.07.09l-.19.55c-.01.06-.08.09-.14.09L1.32 11.29zm3.15 1.66c-.06 0-.1-.05-.07-.09l.32-.58c.03-.05.09-.09.16-.09h4.73c.06 0 .1.05.1.1l-.04.55c0 .06-.06.1-.11.1l-5.09.01zm14.78-3.21l-2.54.71c-.21.06-.22.07-.4-.15-.21-.25-.37-.42-.67-.56-.89-.44-1.76-.31-2.56.17-.95.57-1.44 1.44-1.42 2.55.02 1.09.75 1.98 1.82 2.12 1.01.14 1.85-.23 2.49-1.01.13-.16.25-.33.39-.52H12.1c-.29 0-.37-.18-.27-.42.18-.44.51-1.17.71-1.54.05-.09.17-.21.35-.21h5.26c-.02.32-.02.63-.07.94-.14.94-.49 1.81-1.05 2.58-1 1.36-2.32 2.21-3.97 2.46-1.36.21-2.63-.03-3.74-.87-1.01-.77-1.59-1.79-1.74-3.03-.18-1.5.25-2.83 1.17-3.97 1.01-1.25 2.32-2.02 3.92-2.23 1.28-.17 2.49.04 3.55.83.63.47 1.11 1.07 1.43 1.79.08.14.03.21-.14.26', color: '#00ADD8' },
    'Docker': { icon: 'M20.13 10.59c-.14-.25-.49-.42-1.01-.5.1-.61.04-1.14-.2-1.59-.24-.46-.64-.76-1.2-.9l-.24-.06-.16.2c-.2.24-.34.52-.44.82-.15.44-.06.85.24 1.21-.36.2-.94.32-1.41.38H2.61c-.31 0-.56.25-.58.56-.04.5-.02 1.01.12 1.52.16.57.44 1.05.84 1.43.44.42 1.12.75 1.93.92.94.2 1.97.24 3.05.12.85-.09 1.66-.29 2.42-.6.66-.27 1.24-.63 1.74-1.07.82-.72 1.31-1.58 1.67-2.34h.14c.88 0 1.42-.35 1.72-.66.2-.2.35-.44.44-.68l.06-.21-.14-.08zM4.35 12.5h1.26c.06 0 .11-.05.11-.11V11.2c0-.06-.05-.11-.11-.11H4.35c-.06 0-.11.05-.11.11v1.19c0 .06.05.11.11.11zm1.97 0h1.26c.06 0 .11-.05.11-.11V11.2c0-.06-.05-.11-.11-.11H6.32c-.06 0-.11.05-.11.11v1.19c0 .06.05.11.11.11zm2 0h1.26c.06 0 .11-.05.11-.11V11.2c0-.06-.05-.11-.11-.11H8.32c-.06 0-.11.05-.11.11v1.19c0 .06.05.11.11.11zm1.98 0h1.26c.06 0 .11-.05.11-.11V11.2c0-.06-.05-.11-.11-.11h-1.26c-.06 0-.11.05-.11.11v1.19c0 .06.05.11.11.11zm-3.97-1.73h1.26c.06 0 .11-.05.11-.11V9.47c0-.06-.05-.11-.11-.11H6.33c-.06 0-.11.05-.11.11v1.19c0 .06.05.11.11.11zm2 0h1.26c.06 0 .11-.05.11-.11V9.47c0-.06-.05-.11-.11-.11H8.33c-.06 0-.11.05-.11.11v1.19c0 .06.05.11.11.11zm1.98 0h1.26c.06 0 .11-.05.11-.11V9.47c0-.06-.05-.11-.11-.11h-1.26c-.06 0-.11.05-.11.11v1.19c0 .06.05.11.11.11zm0-1.73h1.26c.06 0 .11-.05.11-.11V7.74c0-.06-.05-.11-.11-.11h-1.26c-.06 0-.11.05-.11.11v1.19c0 .06.05.11.11.11zm1.99 1.73h1.26c.06 0 .11-.05.11-.11V9.47c0-.06-.05-.11-.11-.11h-1.26c-.06 0-.11.05-.11.11v1.19c0 .06.05.11.11.11z', color: '#2496ED' },
    'PostgreSQL': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 2c1.86 0 3.55.7 4.85 1.83L12 10.68 7.15 5.83A7.96 7.96 0 0112 4zm-8 8c0-1.86.7-3.55 1.83-4.85L10.68 12l-4.85 4.85A7.96 7.96 0 014 12zm8 8c-1.86 0-3.55-.7-4.85-1.83L12 13.32l4.85 4.85A7.96 7.96 0 0112 20zm6.17-3.15L13.32 12l4.85-4.85A7.96 7.96 0 0120 12c0 1.86-.7 3.55-1.83 4.85z', color: '#336791' },
    'AWS': { icon: 'M18.75 11.35a4.32 4.32 0 01-.79-.08 3.55 3.55 0 01-.73-.22l-.24-.1-.17.49c-.1.28-.24.53-.4.74-.17.22-.37.4-.61.54a2.5 2.5 0 01-.84.29 4.09 4.09 0 01-1.05.08h-.62a2.5 2.5 0 01-1.7-.69 2.32 2.32 0 01-.65-1.77V9.19h-.96v1.62a3.24 3.24 0 00.92 2.42 3.32 3.32 0 002.39.95h.62c.34 0 .67-.03.99-.1.33-.07.64-.18.93-.33.3-.16.56-.36.78-.61.22-.26.4-.56.52-.9l.15-.41.28.12c.21.08.45.15.69.2.25.05.5.07.75.07h.03v-1.27l-.28-.6z', color: '#FF9900' },
    'Terraform': { icon: 'M1 2.5l7.5 4.3v8.6L1 11.1V2.5zm8.5 4.3L17 2.5v8.6l-7.5 4.3V6.8zM17 13.7l-7.5 4.3v-8.6L17 5.1v8.6zm1-11.2l7.5 4.3v8.6L18 11.1V2.5z', color: '#7B42BC' },
    'Kubernetes': { icon: 'M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4z', color: '#326CE5' },
    'Godot': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z', color: '#478CBF' },
    'Electron': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8z', color: '#47848F' },
    'Astro': { icon: 'M16.074 16.86c-.893.496-2.093.774-3.442.774-1.35 0-2.55-.278-3.442-.773-.847-.472-1.393-1.117-1.54-1.858-.51 1.336-.32 2.846.584 4.265.123.193.256.38.398.56C10.057 21.216 11.612 22 12.632 22c1.019 0 2.574-.784 3.999-2.172.142-.18.275-.367.398-.56.904-1.42 1.094-2.93.584-4.265-.147.741-.693 1.386-1.54 1.858M8.632 2l3 5.5L14.632 2', color: '#FF5D01' },
    'Tailwind': { icon: 'M12 6c-2.67 0-4.33 1.33-5 4 1-1.33 2.17-1.83 3.5-1.5.76.19 1.31.74 1.91 1.35.98 1 2.09 2.15 4.59 2.15 2.67 0 4.33-1.33 5-4-1 1.33-2.17 1.83-3.5 1.5-.76-.19-1.3-.74-1.91-1.35C15.61 7.15 14.5 6 12 6M7 12c-2.67 0-4.33 1.33-5 4 1-1.33 2.17-1.83 3.5-1.5.76.19 1.3.74 1.91 1.35C8.39 16.85 9.5 18 12 18c2.67 0 4.33-1.33 5-4-1 1.33-2.17 1.83-3.5 1.5-.76-.19-1.3-.74-1.91-1.35C10.61 13.15 9.5 12 7 12z', color: '#06B6D4' },
    'MySQL': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#4479A1' },
    'Stripe': { icon: 'M13.976 9.15c-2.172-.806-3.356-1.426-3.356-2.409 0-.831.683-1.305 1.901-1.305 2.227 0 4.515.858 6.09 1.631l.89-5.494C18.252.975 15.697 0 12.165 0 9.667 0 7.589.654 6.104 1.872 4.56 3.147 3.757 4.992 3.757 7.218c0 4.039 2.467 5.76 6.476 7.219 2.585.92 3.445 1.574 3.445 2.583 0 .98-.84 1.545-2.354 1.545-1.875 0-4.965-.921-7.076-2.19l-.89 5.666C5.264 23.074 8.197 24 11.714 24c2.641 0 4.843-.624 6.328-1.813 1.664-1.305 2.525-3.236 2.525-5.732 0-4.128-2.524-5.851-6.591-7.305z', color: '#635BFF' },
    'PHP': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 14H9v-2H7V9h2V7h2v7zm6 0h-2v-2h-2V9h2V7h2v7z', color: '#777BB4' },
    'SQLite': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#003B57' },
    'REST API': { icon: 'M4 1h16a2 2 0 012 2v18a2 2 0 01-2 2H4a2 2 0 01-2-2V3a2 2 0 012-2zm1 4v2h14V5H5zm0 4v2h14V9H5zm0 4v2h10v-2H5z', color: '#FF6C37' },
    'Next.js': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#000000' },
    'D3.js': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#F9A03C' },
    'WebSockets': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#FF6600' },
    'RabbitMQ': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#FF6600' },
    'gRPC': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#244C5A' },
    'GitHub Actions': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#2088FF' },
    'Linux': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#FCC624' },
    'Framer Motion': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#0055FF' },
    'Vercel': { icon: 'M12 2L2 19.5h20L12 2z', color: '#000000' },
    'WhatsApp API': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#25D366' },
    'GDScript': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#478CBF' },
    'Aseprite': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#7D929E' },
    'Pixel Art': { icon: 'M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2z', color: '#E91E63' },
  };

  function getTechColor(tag: string): string {
    return TECH_ICONS[tag]?.color || 'var(--color-accent-primary)';
  }

  function setFilter(cat: string) {
    activeCategory = cat;
    // Animate cards in with anime.js
    setTimeout(() => {
      if (gridEl) {
        const cards = gridEl.querySelectorAll('.project-card');
        if (cards.length > 0) {
          animate(cards, {
            opacity: [0, 1],
            translateY: [30, 0],
            delay: stagger(60),
            duration: 500,
            easing: 'easeOutCubic',
          });
        }
      }
    }, 20);
  }

  onMount(() => {
    // Initial stagger animation
    if (gridEl) {
      const cards = gridEl.querySelectorAll('.project-card');
      animate(cards, {
        opacity: [0, 1],
        translateY: [40, 0],
        scale: [0.97, 1],
        delay: stagger(80, { start: 100 }),
        duration: 600,
        easing: 'easeOutCubic',
      });
    }
  });
</script>

<!-- Category filters -->
<div class="filters">
  <button
    class="filter-pill"
    class:active={activeCategory === 'all'}
    onclick={() => setFilter('all')}
    type="button"
  >
    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M4 8h4V4H4v4zm6 12h4v-4h-4v4zm-6 0h4v-4H4v4zm0-6h4v-4H4v4zm6 0h4v-4h-4v4zm6-10v4h4V4h-4zm-6 4h4V4h-4v4zm6 6h4v-4h-4v4zm0 6h4v-4h-4v4z"/></svg>
    {lang === 'es' ? 'Todos' : 'All'}
  </button>
  {#each categories as cat}
    <button
      class="filter-pill"
      class:active={activeCategory === cat.value}
      onclick={() => setFilter(cat.value)}
      type="button"
    >
      {#if categoryIcons[cat.value]}
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d={categoryIcons[cat.value]}/></svg>
      {/if}
      {cat.label}
    </button>
  {/each}
</div>

<!-- Results count -->
<div class="results-bar">
  <span class="results-count">
    <strong>{filteredProjects.length}</strong>
    {lang === 'es' ? 'proyectos' : 'projects'}
    {#if activeCategory !== 'all'}
      <span class="results-filter">
        {lang === 'es' ? 'en' : 'in'} <em>{categories.find(c => c.value === activeCategory)?.label}</em>
      </span>
    {/if}
  </span>
</div>

<!-- Projects grid -->
{#if filteredProjects.length === 0}
  <div class="no-projects">
    <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
    <p>{lang === 'es' ? 'No hay proyectos en esta categoría aún.' : 'No projects in this category yet.'}</p>
  </div>
{:else}
  <div class="projects-grid" bind:this={gridEl}>
    {#each filteredProjects as project, i (project.id)}
      <a
        href={`${basePath}${project.id.replace(/^(es|en)\//, '')}/`}
        class="project-card"
      >
        <div class="card-image">
          <img src={project.data.heroImage} alt={project.data.title} loading="lazy" />
          <span class="card-badge">{getCategoryLabel(project.data.category)}</span>
          <div class="card-overlay">
            <p>{project.data.description}</p>
            <span class="overlay-cta">
              {lang === 'es' ? 'Ver proyecto' : 'View project'}
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="5" y1="12" x2="19" y2="12"/><polyline points="12 5 19 12 12 19"/></svg>
            </span>
          </div>
        </div>
        <div class="card-body">
          <h3 class="card-title">{project.data.title}</h3>
          <p class="card-desc">{project.data.description}</p>
          <div class="card-tags">
            {#each project.data.tags.slice(0, 5) as tag}
              <span class="tag" style="--tag-color: {getTechColor(tag)}">
                <span class="tag-dot" style="background: {getTechColor(tag)}"></span>
                {tag}
              </span>
            {/each}
            {#if project.data.tags.length > 5}
              <span class="tag tag-more">+{project.data.tags.length - 5}</span>
            {/if}
          </div>
        </div>
      </a>
    {/each}
  </div>
{/if}

<style>
  /* ─── Filters ─── */
  .filters {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
  }

  .filter-pill {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.5rem 1rem;
    border-radius: 9999px;
    border: 1px solid var(--color-border);
    background: var(--color-bg-elevated);
    color: var(--color-text-secondary);
    font-size: 0.8125rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .filter-pill:hover {
    border-color: var(--color-accent-primary);
    color: var(--color-accent-primary);
    background: var(--color-accent-subtle);
    transform: translateY(-1px);
  }

  .filter-pill.active {
    background: var(--color-accent-primary);
    border-color: var(--color-accent-primary);
    color: #fff;
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.25);
  }

  /* ─── Results bar ─── */
  .results-bar {
    display: flex;
    align-items: center;
    margin-bottom: 1.75rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid var(--color-border);
  }

  .results-count {
    font-size: 0.875rem;
    color: var(--color-text-muted);
  }

  .results-count strong {
    color: var(--color-text-primary);
    font-weight: 700;
    font-size: 1rem;
  }

  .results-filter {
    margin-left: 0.25rem;
  }

  .results-filter em {
    color: var(--color-accent-primary);
    font-style: normal;
    font-weight: 600;
  }

  /* ─── Empty state ─── */
  .no-projects {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 4rem 1rem;
    color: var(--color-text-muted);
  }

  .no-projects p {
    font-size: 1rem;
  }

  /* ─── Grid ─── */
  .projects-grid {
    display: grid;
    gap: 1.5rem;
    grid-template-columns: 1fr;
  }

  @media (min-width: 640px) {
    .projects-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (min-width: 1024px) {
    .projects-grid {
      grid-template-columns: repeat(3, 1fr);
    }
  }

  /* ─── Card ─── */
  .project-card {
    display: flex;
    flex-direction: column;
    border-radius: 1rem;
    border: 1px solid var(--color-border);
    background: var(--color-bg-elevated);
    overflow: hidden;
    text-decoration: none;
    color: inherit;
    transition: border-color 0.25s, box-shadow 0.25s, transform 0.25s;
    opacity: 0; /* anime.js will animate in */
  }

  .project-card:hover {
    border-color: var(--color-accent-light);
    box-shadow: 0 12px 40px rgba(37, 99, 235, 0.1);
    transform: translateY(-4px);
  }

  /* ─── Card image ─── */
  .card-image {
    position: relative;
    aspect-ratio: 16 / 10;
    overflow: hidden;
  }

  .card-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.4s ease;
  }

  .project-card:hover .card-image img {
    transform: scale(1.05);
  }

  .card-badge {
    position: absolute;
    top: 0.75rem;
    left: 0.75rem;
    font-size: 0.6875rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    padding: 0.3rem 0.75rem;
    border-radius: 9999px;
    background: rgba(0, 0, 0, 0.55);
    backdrop-filter: blur(8px);
    color: #f1f5f9;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .card-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    gap: 0.75rem;
    padding: 1.25rem;
    background: linear-gradient(to top, rgba(0, 0, 0, 0.8) 0%, rgba(0, 0, 0, 0.2) 50%, transparent 100%);
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .project-card:hover .card-overlay {
    opacity: 1;
  }

  .card-overlay p {
    color: rgba(255, 255, 255, 0.9);
    font-size: 0.8125rem;
    line-height: 1.5;
  }

  .overlay-cta {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 0.8125rem;
    font-weight: 600;
    color: #fff;
    transition: gap 0.2s;
  }

  .project-card:hover .overlay-cta {
    gap: 0.625rem;
  }

  /* ─── Card body ─── */
  .card-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 1.25rem;
  }

  .card-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.0625rem;
    color: var(--color-text-primary);
    margin-bottom: 0.5rem;
    line-height: 1.3;
  }

  .card-desc {
    font-size: 0.8125rem;
    line-height: 1.5;
    color: var(--color-text-muted);
    margin-bottom: 1rem;
    flex: 1;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  /* ─── Tags ─── */
  .card-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.375rem;
    margin-top: auto;
  }

  .tag {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.6875rem;
    padding: 0.25rem 0.5rem;
    border-radius: 9999px;
    background: var(--color-bg-secondary);
    color: var(--color-text-secondary);
    font-weight: 500;
    font-family: var(--font-mono);
    letter-spacing: -0.01em;
  }

  .tag-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .tag-more {
    background: var(--color-bg-secondary);
    color: var(--color-text-muted);
    font-family: var(--font-sans);
  }

  @media (max-width: 640px) {
    .filter-pill {
      font-size: 0.75rem;
      padding: 0.375rem 0.75rem;
    }
  }
</style>

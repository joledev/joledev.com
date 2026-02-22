<script lang="ts">
  import { onMount } from 'svelte';
  import * as THREE from 'three';

  let canvas: HTMLCanvasElement;

  onMount(() => {
    const isMobile = window.innerWidth < 768;
    const isDark = () => document.documentElement.getAttribute('data-theme') === 'dark';

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(60, 1, 0.1, 100);
    camera.position.z = 20;

    const renderer = new THREE.WebGLRenderer({
      canvas,
      alpha: true,
      antialias: true,
      powerPreference: 'low-power',
    });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

    // ── Lighting ──────────────────────────────────────────
    const ambientLight = new THREE.AmbientLight(0x60a5fa, 0.4);
    scene.add(ambientLight);

    const mainLight = new THREE.DirectionalLight(0xffffff, 1.0);
    mainLight.position.set(5, 5, 8);
    scene.add(mainLight);

    const rimLight = new THREE.DirectionalLight(0x2563eb, 0.6);
    rimLight.position.set(-5, -3, -5);
    scene.add(rimLight);

    const pointLight = new THREE.PointLight(0x3b82f6, 0.8, 30);
    pointLight.position.set(0, 3, 6);
    scene.add(pointLight);

    // ── Main geometry: Icosahedron with glass material ────
    const icoGeo = new THREE.IcosahedronGeometry(isMobile ? 3.5 : 4.5, 1);

    const glassMat = new THREE.MeshPhysicalMaterial({
      color: 0x2563eb,
      metalness: 0.1,
      roughness: 0.05,
      transmission: 0.92,
      thickness: 2.0,
      ior: 1.5,
      envMapIntensity: 1.0,
      clearcoat: 1.0,
      clearcoatRoughness: 0.1,
      transparent: true,
      opacity: 0.85,
      side: THREE.DoubleSide,
    });

    const icoMesh = new THREE.Mesh(icoGeo, glassMat);
    scene.add(icoMesh);

    // Wireframe overlay
    const wireGeo = new THREE.IcosahedronGeometry(isMobile ? 3.55 : 4.55, 1);
    const wireMat = new THREE.MeshBasicMaterial({
      color: 0x60a5fa,
      wireframe: true,
      transparent: true,
      opacity: 0.3,
    });
    const wireMesh = new THREE.Mesh(wireGeo, wireMat);
    scene.add(wireMesh);

    // ── Orbiting ring (torus) ─────────────────────────────
    const torusGeo = new THREE.TorusGeometry(isMobile ? 5.5 : 7, 0.06, 16, 100);
    const torusMat = new THREE.MeshBasicMaterial({
      color: 0x3b82f6,
      transparent: true,
      opacity: 0.35,
    });
    const torus = new THREE.Mesh(torusGeo, torusMat);
    torus.rotation.x = Math.PI * 0.5;
    scene.add(torus);

    // Second ring
    const torus2Geo = new THREE.TorusGeometry(isMobile ? 6.5 : 8.5, 0.04, 16, 100);
    const torus2Mat = new THREE.MeshBasicMaterial({
      color: 0x2563eb,
      transparent: true,
      opacity: 0.2,
    });
    const torus2 = new THREE.Mesh(torus2Geo, torus2Mat);
    torus2.rotation.x = Math.PI * 0.6;
    torus2.rotation.z = Math.PI * 0.15;
    scene.add(torus2);

    // ── Small floating orbs ───────────────────────────────
    const orbCount = isMobile ? 5 : 8;
    const orbs: THREE.Mesh[] = [];
    const orbData: { angle: number; radius: number; speed: number; yOffset: number }[] = [];

    for (let i = 0; i < orbCount; i++) {
      const orbGeo = new THREE.SphereGeometry(0.15 + Math.random() * 0.15, 16, 16);
      const orbOpacity = 0.6 + Math.random() * 0.3;
      const orbMat = new THREE.MeshBasicMaterial({
        color: [0x2563eb, 0x3b82f6, 0x60a5fa][i % 3],
        transparent: true,
        opacity: orbOpacity,
      });
      orbMat.userData = { baseOpacity: orbOpacity };
      const orb = new THREE.Mesh(orbGeo, orbMat);
      orbs.push(orb);
      scene.add(orb);
      orbData.push({
        angle: (Math.PI * 2 * i) / orbCount,
        radius: (isMobile ? 5 : 6.5) + Math.random() * 2,
        speed: 0.15 + Math.random() * 0.2,
        yOffset: (Math.random() - 0.5) * 3,
      });
    }

    // ── Background particles (subtle) ─────────────────────
    const PARTICLE_COUNT = isMobile ? 40 : 80;
    const pPositions = new Float32Array(PARTICLE_COUNT * 3);
    const pColors = new Float32Array(PARTICLE_COUNT * 3);
    const pSizes = new Float32Array(PARTICLE_COUNT);

    const palette = [
      new THREE.Color('#2563EB'),
      new THREE.Color('#3B82F6'),
      new THREE.Color('#60A5FA'),
    ];

    for (let i = 0; i < PARTICLE_COUNT; i++) {
      const i3 = i * 3;
      pPositions[i3] = (Math.random() - 0.5) * 60;
      pPositions[i3 + 1] = (Math.random() - 0.5) * 40;
      pPositions[i3 + 2] = (Math.random() - 0.5) * 30 - 10;

      const color = palette[Math.floor(Math.random() * palette.length)];
      pColors[i3] = color.r;
      pColors[i3 + 1] = color.g;
      pColors[i3 + 2] = color.b;

      pSizes[i] = Math.random() * 2 + 0.5;
    }

    const pGeo = new THREE.BufferGeometry();
    pGeo.setAttribute('position', new THREE.BufferAttribute(pPositions, 3));
    pGeo.setAttribute('color', new THREE.BufferAttribute(pColors, 3));
    pGeo.setAttribute('size', new THREE.BufferAttribute(pSizes, 1));

    const pMat = new THREE.ShaderMaterial({
      vertexShader: `
        attribute float size;
        varying vec3 vColor;
        void main() {
          vColor = color;
          vec4 mvPosition = modelViewMatrix * vec4(position, 1.0);
          gl_PointSize = size * (150.0 / -mvPosition.z);
          gl_Position = projectionMatrix * mvPosition;
        }
      `,
      fragmentShader: `
        varying vec3 vColor;
        void main() {
          float d = length(gl_PointCoord - vec2(0.5));
          if (d > 0.5) discard;
          float alpha = 1.0 - smoothstep(0.2, 0.5, d);
          gl_FragColor = vec4(vColor, alpha * 0.4);
        }
      `,
      vertexColors: true,
      transparent: true,
      depthWrite: false,
    });

    const particles = new THREE.Points(pGeo, pMat);
    scene.add(particles);

    // ── Theme reactivity ────────────────────────────────────
    function applyTheme() {
      const dark = isDark();
      glassMat.opacity = dark ? 0.85 : 0.45;
      wireMat.opacity = dark ? 0.3 : 0.12;
      torusMat.opacity = dark ? 0.35 : 0.15;
      torus2Mat.opacity = dark ? 0.2 : 0.08;
      orbs.forEach((o) => {
        const mat = o.material as THREE.MeshBasicMaterial;
        mat.opacity = dark ? mat.userData.baseOpacity : mat.userData.baseOpacity * 0.4;
      });
      pMat.fragmentShader = `
        varying vec3 vColor;
        void main() {
          float d = length(gl_PointCoord - vec2(0.5));
          if (d > 0.5) discard;
          float alpha = 1.0 - smoothstep(0.2, 0.5, d);
          gl_FragColor = vec4(vColor, alpha * ${dark ? '0.4' : '0.15'});
        }
      `;
      pMat.needsUpdate = true;
    }

    const themeObserver = new MutationObserver(applyTheme);
    themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] });
    applyTheme();

    // ── Mouse tracking ────────────────────────────────────
    let mouseX = 0;
    let mouseY = 0;
    let smoothX = 0;
    let smoothY = 0;

    const onMouseMove = (e: MouseEvent) => {
      mouseX = (e.clientX / window.innerWidth - 0.5) * 2;
      mouseY = (e.clientY / window.innerHeight - 0.5) * 2;
    };

    if (!isMobile) {
      window.addEventListener('mousemove', onMouseMove, { passive: true });
    }

    // ── Resize ────────────────────────────────────────────
    let ready = false;
    const resize = () => {
      const w = canvas.clientWidth;
      const h = canvas.clientHeight;
      if (w === 0 || h === 0) return;
      ready = true;
      renderer.setSize(w, h, false);
      camera.aspect = w / h;
      camera.updateProjectionMatrix();
    };
    resize();
    window.addEventListener('resize', resize, { passive: true });

    // ── Animation loop ────────────────────────────────────
    let animId: number;
    const clock = new THREE.Clock();

    const animate = () => {
      animId = requestAnimationFrame(animate);
      if (!ready) { resize(); return; }
      const t = clock.getElapsedTime();

      // Smooth mouse interpolation
      if (!isMobile) {
        smoothX += (mouseX - smoothX) * 0.03;
        smoothY += (mouseY - smoothY) * 0.03;
      } else {
        smoothX = Math.sin(t * 0.2) * 0.3;
        smoothY = Math.cos(t * 0.15) * 0.2;
      }

      // Icosahedron: rotate + float
      icoMesh.rotation.x = t * 0.15 + smoothY * 0.5;
      icoMesh.rotation.y = t * 0.2 + smoothX * 0.5;
      icoMesh.position.y = Math.sin(t * 0.5) * 0.6;

      // Wireframe follows glass
      wireMesh.rotation.copy(icoMesh.rotation);
      wireMesh.position.copy(icoMesh.position);

      // Torus rings
      torus.rotation.z = t * 0.1;
      torus.rotation.x = Math.PI * 0.5 + smoothY * 0.3;
      torus2.rotation.z = -t * 0.08;

      // Orbiting orbs
      for (let i = 0; i < orbCount; i++) {
        const d = orbData[i];
        const angle = d.angle + t * d.speed;
        orbs[i].position.x = Math.cos(angle) * d.radius;
        orbs[i].position.z = Math.sin(angle) * d.radius * 0.4;
        orbs[i].position.y = d.yOffset + Math.sin(t * 0.8 + i) * 0.5;
      }

      // Particle field subtle rotation
      particles.rotation.y = t * 0.02 + smoothX * 0.05;
      particles.rotation.x = smoothY * 0.03;

      // Point light follows mouse subtly
      pointLight.position.x = smoothX * 4;
      pointLight.position.y = 3 - smoothY * 2;

      renderer.render(scene, camera);
    };
    animate();

    return () => {
      cancelAnimationFrame(animId);
      renderer.dispose();
      icoGeo.dispose();
      glassMat.dispose();
      wireGeo.dispose();
      wireMat.dispose();
      torusGeo.dispose();
      torusMat.dispose();
      torus2Geo.dispose();
      torus2Mat.dispose();
      orbs.forEach((o) => {
        o.geometry.dispose();
        (o.material as THREE.Material).dispose();
      });
      pGeo.dispose();
      pMat.dispose();
      themeObserver.disconnect();
      window.removeEventListener('resize', resize);
      if (!isMobile) window.removeEventListener('mousemove', onMouseMove);
    };
  });
</script>

<canvas bind:this={canvas} class="hero-canvas"></canvas>

<style>
  .hero-canvas {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    z-index: 0;
  }
</style>

---
name: frontend-vue-pinia-tailwind
description: >-
  Builds Vue 3 frontends with Pinia and Tailwind CSS in a polished, product-grade
  UI style (Vercel / v0 aesthetic). Use when the user works on Vue, SPA, dashboards,
  marketing UI, Pinia stores, Tailwind design tokens, or asks for professional
  frontend styling.
---

# Frontend: Vue 3 + Pinia + Tailwind (estilo Vercel / v0)

## Stack fija

- **Vue 3** con `<script setup>` y Composition API.
- **Vite** como bundler.
- **Pinia** para estado global; composables solo para estado local o lógica reutilizable sin persistencia de dominio.
- **Tailwind CSS** v4 si el proyecto ya lo usa; si no, la versión que declare `package.json`, sin mezclar sintaxis de versiones.
- **Vue Router** cuando haya más de una vista.

No sustituir Pinia por stores ad-hoc en `provide/inject` salvo casos muy acotados (temas, i18n).

## Arquitectura de carpetas sugerida

```
src/
  assets/
  components/
    ui/           # primitivos reutilizables (Button, Card, Input, Badge)
    layout/       # shell, nav, sidebar
  composables/
  stores/         # Pinia: un archivo por dominio (auth, user, cart, …)
  views/          # o pages/
  router/
  lib/            # clientes API, helpers sin UI
  styles/         # @tailwind; capa base si hace falta
```

- Componentes **presentacionales** pequeños; vistas **orquestan** datos y navegación.
- Llamadas HTTP en `lib/api.ts` (o por recurso) con funciones puras; las stores llaman a esas funciones y guardan estado.

## Pinia: reglas

- `defineStore` con id en kebab-case o `domainFeature` consistente con el repo.
- Estado serializable (JSON); no guardar instancias de clases ni DOM.
- Acciones async para side effects; getters solo derivación síncrona.
- Si hace falta persistencia, usar plugin oficial o localStorage en una acción dedicada, no esparcido en componentes.

## Tailwind: tokens y capas (look Vercel / v0)

Objetivo: **claro, sobrio, mucho aire**, bordes finos, tipografía legible, foco y hover evidentes pero discretos.

1. **Color**
   - Fondo: blanco casi puro / zinc-50 en claro; zinc-950 o neutral-950 en oscuro.
   - Superficies elevadas: un pelín más claras/oscuras que el fondo, no saturadas.
   - Texto primario: zinc-900 / zinc-50; secundario: zinc-600 / zinc-400.
   - Un solo **acento** (por ejemplo zinc o un azul muy contenido); reservar color fuerte para CTAs y estados.

2. **Tipografía**
   - Preferir una familial **sans geométrica** (p. ej. Inter, Geist, system-ui stack coherente).
   - Escala: `text-sm` cuerpo, `text-xs` meta; títulos con `font-semibold` o `font-medium`, evitar `font-black` salvo hero puntual.
   - `tracking-tight` en titulares grandes; `leading-relaxed` en párrafos largos.

3. **Layout y ritmo**
   - Contenedor legible: `max-w-6xl` o `max-w-7xl` con `mx-auto px-4 sm:px-6 lg:px-8`.
   - Espaciado en múltiplos de 4; secciones con `py-16` o `py-24` en landings.
   - Grids responsivos: `grid gap-6 md:grid-cols-2 lg:grid-cols-3`.

4. **Bordes y profundidad**
   - `border border-zinc-200 dark:border-zinc-800` en cards y inputs.
   - Sombras **muy suaves** (`shadow-sm`) o ninguna; preferir borde a sombra pesada.
   - `rounded-lg` o `rounded-xl` de forma **consistente** en todo el proyecto.

5. **Interacción**
   - Estados: `hover:`, `active:`, `focus-visible:ring-2 ring-offset-2` en botones y enlaces.
   - `transition-colors duration-200` (o 150) en elementos interactivos; evitar animaciones largas por defecto.
   - `disabled:opacity-50 disabled:pointer-events-none` en controles deshabilitados.

6. **Componentes UI**
   - Botones: variante default (sólido discreto), secondary (outline), ghost, destructive si aplica.
   - Inputs: altura uniforme, label arriba o flotante consistente, mensaje de error con `text-sm text-red-600`.
   - Cards: padding uniforme (`p-6`), opcional header con borde inferior sutil.

7. **Oscuro**
   - Implementar **dark mode** con `class` en `<html>` y utilidades `dark:` en componentes compartidos; no duplicar componentes por tema.

## Accesibilidad y calidad

- Contraste suficiente (WCAG AA) en texto sobre fondos.
- Botones y enlaces con foco visible; roles correctos en diálogos si se usan.
- Estados de carga y error en vistas que fetchean datos (skeleton o spinner discreto, mensaje claro).

## Checklist antes de dar por hecho un front

- [ ] Vista alineada con tokens (espaciado, radios, bordes).
- [ ] Pinia: store acotada, sin lógica de UI mezclada.
- [ ] Sin estilos inline arbitrarios que rompan el sistema; preferir clases Tailwind o `@apply` mínimo en `styles/`.
- [ ] Responsive revisado en al menos sm / md / lg.
- [ ] El front no debe contener lógica de negocio solo mostrar la información que viene del API

## Referencia ampliada

Para patrones de componentes y copy de UI, ver [reference-ui-patterns.md](reference-ui-patterns.md).
